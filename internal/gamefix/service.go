package gamefix

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

// ServiceName is the Windows service that watches sources while the app is
// closed. It is RapTech-Portable.exe itself, started with ServiceFlag. Its
// name and flags are its own, so it never touches the standalone GameFix's
// GameFixWatch service.
const (
	ServiceName  = "RapTechGameFixWatch"
	ServiceFlag  = "--gamefix-service"
	installFlag  = "--gamefix-install-service"
	removeFlag   = "--gamefix-uninstall-service"
	serviceLog   = "RapTech-GameFix-service.log" // next to RapTech-Portable.exe
	reloadPeriod = 2 * time.Second               // how often the service looks for changed settings
)

// RunCommandLine handles the ways RapTech-Portable.exe is started for the
// GameFix tab other than as the app (as its service, or to install or
// remove it), and reports whether it did.
func RunCommandLine() bool {
	if isService, _ := svc.IsWindowsService(); isService {
		if err := svc.Run(ServiceName, &watchService{}); err != nil {
			logFile(serviceLog).Println("service failed:", err)
		}
		return true
	}
	if len(os.Args) < 2 {
		return false
	}
	var err error
	switch os.Args[1] {
	case installFlag:
		err = installService()
	case removeFlag:
		err = removeService()
	default:
		return false
	}
	if err != nil {
		// Started elevated by the app, which cannot read what it prints.
		title, _ := windows.UTF16PtrFromString("Diskless PathFix SymbolicLink service")
		text, _ := windows.UTF16PtrFromString(err.Error())
		windows.MessageBox(0, text, title, windows.MB_OK|windows.MB_ICONERROR)
		os.Exit(1)
	}
	return true
}

// watchService is the Windows service. It watches the folders that have
// Enable on in GameFix.json, like the app does while it is open, and
// reloads GameFix.json whenever it changes, so ticking or unticking Enable
// in the app starts or stops watching that folder.
type watchService struct{}

func (watchService) Execute(args []string, requests <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32) {
	status <- svc.Status{State: svc.StartPending}
	a := &App{settingsPath: exeSettingsPath(), service: true, log: logFile(serviceLog)}
	a.log.Println("started")
	var loaded time.Time
	reload := func() {
		st, err := os.Stat(a.settingsPath)
		if err != nil || st.ModTime().Equal(loaded) {
			return
		}
		settings, err := readSettings(a.settingsPath)
		if err != nil {
			return // the app may still be writing it; tried again next time
		}
		loaded = st.ModTime()
		a.mu.Lock()
		defer a.mu.Unlock()
		a.settings = settings
		a.useSavedVars()
		a.syncWatches()
		a.log.Printf("watching %d folders", len(a.watches))
	}
	reload()
	status <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

	tick := time.NewTicker(reloadPeriod)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			reload()
		case r := <-requests:
			switch r.Cmd {
			case svc.Interrogate:
				status <- r.CurrentStatus
			case svc.Stop, svc.Shutdown:
				status <- svc.Status{State: svc.StopPending}
				a.mu.Lock()
				a.settings = Settings{}
				a.syncWatches() // stops them all
				a.mu.Unlock()
				a.log.Println("stopped")
				return false, 0
			}
		}
	}
}

// logFile writes to the file name next to GameFix.exe, which is started
// again when it gets bigger than 1 MB.
func logFile(name string) *log.Logger {
	path := filepath.Join(filepath.Dir(exeSettingsPath()), name)
	if st, err := os.Stat(path); err == nil && st.Size() > 1<<20 {
		os.Remove(path)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return log.New(os.Stderr, "", log.LstdFlags)
	}
	return log.New(f, "", log.LstdFlags)
}

// useSavedVars sets the variables that sources are written with to the
// values the app saved. The service runs as the SYSTEM account, whose
// %USERPROFILE% is not the user's.
func (a *App) useSavedVars() {
	for name, value := range a.settings.Vars {
		os.Setenv(name, value)
	}
}

// installService installs and starts the service. It needs administrator
// rights.
func installService() error {
	return installServiceAs(ServiceName, "RapTech GameFix Watch",
		"Copies the folders that have Enable on in the SymbolicLink tab of Diskless PathFix into the main folder whenever their source changes, "+
			"while Diskless PathFix is closed.", ServiceFlag)
}

// removeService stops and uninstalls the service. It needs administrator
// rights.
func removeService() error {
	return removeServiceNamed(ServiceName)
}

// installServiceAs installs and starts this GameFix.exe as the service
// name, started with flag, or points the service at it and starts it if it
// is already installed. It needs administrator rights.
func installServiceAs(name, displayName, description, flag string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err == nil {
		c, err := s.Config()
		if err == nil {
			c.BinaryPathName = `"` + exe + `" ` + flag
			c.StartType = mgr.StartAutomatic
			c.Description = description
			err = s.UpdateConfig(c)
		}
		if err != nil {
			s.Close()
			return err
		}
	} else {
		s, err = m.CreateService(name, exe, mgr.Config{
			DisplayName: displayName,
			Description: description,
			StartType:   mgr.StartAutomatic,
		}, flag)
		if err != nil {
			return err
		}
	}
	defer s.Close()
	// Start again if it stops unexpectedly.
	s.SetRecoveryActions([]mgr.RecoveryAction{{Type: mgr.ServiceRestart, Delay: 5 * time.Second}}, 24*60*60)
	if st, err := s.Query(); err == nil && st.State == svc.Running {
		return nil
	}
	return s.Start()
}

// removeServiceNamed stops and uninstalls a service. It needs
// administrator rights.
func removeServiceNamed(name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return nil // not installed
	}
	defer s.Close()
	if st, err := s.Control(svc.Stop); err == nil {
		for end := time.Now().Add(15 * time.Second); st.State != svc.Stopped && time.Now().Before(end); {
			time.Sleep(300 * time.Millisecond)
			if st, err = s.Query(); err != nil {
				break
			}
		}
	}
	return s.Delete()
}

// serviceState returns "running", "stopped", "starting", "stopping" or
// "not installed". It does not need administrator rights.
func serviceState() string {
	m, err := windows.OpenSCManager(nil, nil, windows.SC_MANAGER_CONNECT)
	if err != nil {
		return "not installed"
	}
	defer windows.CloseServiceHandle(m)
	name, _ := windows.UTF16PtrFromString(ServiceName)
	h, err := windows.OpenService(m, name, windows.SERVICE_QUERY_STATUS)
	if err != nil {
		return "not installed"
	}
	defer windows.CloseServiceHandle(h)
	var st windows.SERVICE_STATUS
	if windows.QueryServiceStatus(h, &st) != nil {
		return "stopped"
	}
	switch st.CurrentState {
	case windows.SERVICE_RUNNING:
		return "running"
	case windows.SERVICE_START_PENDING, windows.SERVICE_CONTINUE_PENDING:
		return "starting"
	case windows.SERVICE_STOP_PENDING, windows.SERVICE_PAUSE_PENDING:
		return "stopping"
	}
	return "stopped"
}

// serviceRunning reports whether the service is running, and so watching
// instead of the app. Tests replace it, so they don't depend on whether the
// service is installed on this PC.
var serviceRunning = func() bool {
	s := serviceState()
	return s == "running" || s == "starting"
}

// GetServiceState returns the state of the service; see serviceState.
func (a *App) GetServiceState() string {
	return serviceState()
}

// InstallService installs and starts the service, asking for administrator
// rights. From then on it watches instead of the app, also while the app is
// closed.
func (a *App) InstallService() error {
	exe, _ := os.Executable()
	if strings.HasSuffix(strings.ToLower(exe), "-dev.exe") {
		return errors.New("install the service from RapTech-Portable.exe made by wails build, not from wails dev")
	}
	a.mu.Lock()
	err := a.save() // saves the variables the service needs; see Settings.Vars
	a.mu.Unlock()
	if err != nil {
		return err
	}
	return runElevated(installFlag)
}

// UninstallService stops and uninstalls the service, asking for
// administrator rights. The app then watches again while it is open.
func (a *App) UninstallService() error {
	return runElevated(removeFlag)
}

// runElevated runs GameFix.exe with flag as administrator and waits for it.
func runElevated(flag string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	cmd := exec.Command("powershell", "-NoProfile", "-Command",
		`try { $p = Start-Process -FilePath $env:GAMEFIX_EXE -ArgumentList $env:GAMEFIX_FLAG -Verb RunAs -Wait -PassThru; exit $p.ExitCode } catch { exit 2 }`)
	cmd.Env = append(os.Environ(), "GAMEFIX_EXE="+exe, "GAMEFIX_FLAG="+flag)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000} // CREATE_NO_WINDOW
	err = cmd.Run()
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		if exit.ExitCode() == 2 {
			return errors.New("administrator rights are needed; it was cancelled")
		}
		return errors.New("it did not work; see the message SymbolicLink showed")
	}
	return err
}
