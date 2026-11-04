package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/sys/windows/registry"

	"raptech-portable/internal/core"
	"raptech-portable/internal/hunt"
	"raptech-portable/internal/icon"
	"raptech-portable/internal/manager"
)

// App is what the window calls. The work is in internal/manager.
type App struct {
	ctx          context.Context
	settingsPath string
	Root         string `json:"root"`

	huntMu sync.Mutex
	hunt   *hunt.Hunter
}

func NewApp() *App { return &App{} }

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	exe, _ := os.Executable()
	a.settingsPath = filepath.Join(filepath.Dir(exe), "RapTech-Portable.json")
	if b, err := os.ReadFile(a.settingsPath); err == nil {
		json.Unmarshal(b, a)
	}
	if a.Root == "" {
		if _, err := os.Stat(`G:\`); err == nil {
			a.Root = `G:\Portable`
		}
	}
	if a.Root != "" {
		go manager.FixIcons(a.Root) // launchers made before their app was there
	}
}

func (a *App) save() error {
	b, _ := json.MarshalIndent(a, "", "  ")
	return os.WriteFile(a.settingsPath, b, 0o644)
}

// GetRoot is the folder that holds the packages.
func (a *App) GetRoot() string { return a.Root }

// ChooseRoot asks for the folder that holds the packages.
func (a *App) ChooseRoot() (string, error) {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{Title: "Folder for portable apps (on the game disk)", DefaultDirectory: a.Root})
	if err != nil || dir == "" {
		return a.Root, err
	}
	a.Root = dir
	return dir, a.save()
}

// PackageView is one package with its schedule, for the window.
type PackageView struct {
	manager.Info
	ScheduleMinutes int    // 0 = off
	CCBoot          string // line for C:\CCBoot\Cmd\public.txt
	Registry        []string
	Session         [][]string
	Links           []string // link paths as written in portable.json
	Icon            string   // the launcher's icon as a data: URL, "" if none to show
	OwnIcon         bool     // a chosen icon, not the app's
	InPlace         []bool   // per folder: an existing folder used where it is
	ExeInPlace      bool     // the app starts where it is installed, not from the package
	Args            []string // arguments the launcher passes to the app
}

// List shows the packages in the root folder.
func (a *App) List() []PackageView {
	if a.Root == "" {
		return nil
	}
	var out []PackageView
	for _, info := range manager.Scan(a.Root) {
		v := PackageView{Info: info}
		if info.Error == "" {
			if on, every := manager.ScheduleOn(info.Name); on {
				v.ScheduleMinutes = int(every / time.Minute)
			}
			v.CCBoot = fmt.Sprintf(`run -winlogon "%s" --setup`, info.Launcher)
			if png, err := icon.PNG(info.Launcher); err == nil {
				v.Icon = "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
			}
			if p, err := core.Load(info.Dir); err == nil {
				v.OwnIcon = manager.HasOwnIcon(p)
				v.Registry = p.Registry
				v.ExeInPlace = p.Exe != "" && core.IsFull(p.Exe)
				v.Args = p.Args
				for _, f := range p.Folders {
					v.Links = append(v.Links, f.Link)
					v.Session = append(v.Session, f.Session)
					v.InPlace = append(v.InPlace, p.InPlace(f))
				}
			}
		}
		out = append(out, v)
	}
	return out
}

func (a *App) launcher() ([]byte, error) {
	b, err := launcherFS.ReadFile("assets/launcher.exe")
	if err != nil {
		return nil, errors.New("this build has no launcher inside; build it with build.cmd")
	}
	return b, nil
}

func (a *App) needRoot() error {
	if a.Root == "" {
		return errors.New("choose the folder for portable apps first")
	}
	return os.MkdirAll(a.Root, 0o755)
}

// RobloxPerUser reports whether Roblox on this PC is installed for one user
// (in %LOCALAPPDATA%) rather than for everyone (in Program Files).
func (a *App) RobloxPerUser() bool {
	if k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Classes\roblox-player\shell\open\command`, registry.QUERY_VALUE); err == nil {
		k.Close()
		return false
	}
	_, err := os.Stat(filepath.Join(os.Getenv("LOCALAPPDATA"), `Roblox\Versions`))
	return err == nil
}

// MakeRoblox makes the Roblox package and downloads the newest Roblox into it.
func (a *App) MakeRoblox(perUser bool) (string, error) {
	if err := a.needRoot(); err != nil {
		return "", err
	}
	l, err := a.launcher()
	if err != nil {
		return "", err
	}
	p, err := manager.CreateRoblox(a.Root, perUser, l)
	if err != nil {
		return "", err
	}
	return a.UpdateNow(p.Dir)
}

// PickExe asks for the program that starts an app.
func (a *App) PickExe() (string, error) {
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{Title: "Program that starts the app",
		Filters: []runtime.FileFilter{{DisplayName: "Programs (*.exe)", Pattern: "*.exe"}}})
}

// PickFolder asks for a folder.
func (a *App) PickFolder(title string) (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{Title: title})
}

// CreatePackage makes a package for any app: appDir is copied in, exe starts it.
func (a *App) CreatePackage(name, appDir, exe string) error {
	if err := a.needRoot(); err != nil {
		return err
	}
	l, err := a.launcher()
	if err != nil {
		return err
	}
	_, err = manager.CreatePackage(a.Root, strings.TrimSpace(name), appDir, exe, l)
	return err
}

func load(dir string) (*core.Package, error) { return core.Load(dir) }

// AddFolder adds a folder the app uses; session is a comma-separated list
// of subfolders emptied after each player.
func (a *App) AddFolder(dir, link, session string) error {
	p, err := load(dir)
	if err != nil {
		return err
	}
	var subs []string
	for _, s := range strings.Split(session, ",") {
		if s = strings.TrimSpace(s); s != "" {
			subs = append(subs, s)
		}
	}
	return manager.AddFolder(p, link, subs)
}

// RemoveFolder takes folder i out of the package (its copy is kept).
func (a *App) RemoveFolder(dir string, i int) error {
	p, err := load(dir)
	if err != nil {
		return err
	}
	return manager.RemoveFolder(p, i)
}

// AddRegistry adds a registry key and saves it now.
func (a *App) AddRegistry(dir, key string) error {
	p, err := load(dir)
	if err != nil {
		return err
	}
	return manager.AddRegistry(p, key)
}

// UpdateNow runs the package's update and returns what it logged. The
// launcher then gets the app's icon (Roblox has none before its first
// download).
func (a *App) UpdateNow(dir string) (string, error) {
	out, err := a.runLauncher(dir, "--update")
	if err == nil {
		if p, lerr := load(dir); lerr == nil {
			manager.ApplyIcon(p) // only looks; skipped while the launcher runs
		}
	}
	return out, err
}

// SetupHere links the package on this PC.
func (a *App) SetupHere(dir string) (string, error) { return a.runLauncher(dir, "--setup") }

// UnlinkHere removes the package's links on this PC.
func (a *App) UnlinkHere(dir string) (string, error) { return a.runLauncher(dir, "--unlink") }

func (a *App) runLauncher(dir, mode string) (string, error) {
	p, err := load(dir)
	if err != nil {
		return "", err
	}
	cmd := exec.Command(filepath.Join(dir, manager.LauncherName(p.Name)), mode)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	out, err := cmd.CombinedOutput()
	text := strings.TrimSpace(string(out))
	if err != nil {
		return text, fmt.Errorf("%s failed: %s", strings.TrimPrefix(mode, "--"), lastLines(text, 3))
	}
	return text, nil
}

func lastLines(s string, n int) string {
	lines := strings.Split(s, "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return strings.Join(lines, "\n")
}

// SetSchedule runs the package's update every minutes as SYSTEM; 0 turns it off.
func (a *App) SetSchedule(dir string, minutes int) error {
	p, err := load(dir)
	if err != nil {
		return err
	}
	if minutes == 0 {
		return manager.RemoveSchedule(p.Name)
	}
	if p.Update == nil {
		return errors.New("this package has no update")
	}
	return manager.SetSchedule(p.Name, filepath.Join(dir, manager.LauncherName(p.Name)), time.Duration(minutes)*time.Minute)
}

// InUse lists the programs running from the app's package and the Windows
// services installed from it, which RemoveApp stops.
func (a *App) InUse(dir string) manager.Use {
	return manager.InUse(&core.Package{Dir: dir})
}

// RemoveApp removes the app: stops its services and programs, removes its
// links, protocol handlers and update task, and deletes its folder. A
// package whose portable.json is broken can be removed too.
func (a *App) RemoveApp(dir string) error {
	p, err := load(dir)
	if err != nil {
		p = &core.Package{Dir: dir, Manifest: core.Manifest{Name: filepath.Base(dir)}}
	}
	return manager.RemovePackage(a.Root, p)
}

// FindRiot is where Riot Client is installed on this PC ("" if not found).
func (a *App) FindRiot() string {
	dir, _ := manager.FindRiot(a.Root)
	return dir
}

// MakeRiot makes the Riot package for the Riot Games folder riotDir, used
// in place (see manager.CreateRiot).
func (a *App) MakeRiot(riotDir string) error {
	if err := a.needRoot(); err != nil {
		return err
	}
	l, err := a.launcher()
	if err != nil {
		return err
	}
	_, err = manager.CreateRiot(a.Root, riotDir, l)
	return err
}

// CreateInPlace makes a package for an app that stays where it is
// installed: exe starts with args (split at spaces); nothing is copied.
func (a *App) CreateInPlace(name, exe, args string) error {
	l, err := a.launcher()
	if err != nil {
		return err
	}
	return a.createInPlace(name, exe, args, l)
}

func (a *App) createInPlace(name, exe, args string, launcher []byte) error {
	if err := a.needRoot(); err != nil {
		return err
	}
	_, err := manager.CreateInPlace(a.Root, strings.TrimSpace(name), exe, strings.Fields(args), launcher)
	return err
}

// AddFolderInPlace links link to target, an existing folder used where it
// is (nothing copied); session as in AddFolder.
func (a *App) AddFolderInPlace(dir, link, target, session string) error {
	p, err := load(dir)
	if err != nil {
		return err
	}
	var subs []string
	for _, s := range strings.Split(session, ",") {
		if s = strings.TrimSpace(s); s != "" {
			subs = append(subs, s)
		}
	}
	return manager.AddFolderInPlace(p, link, target, subs)
}

// ChangeIcon asks for a program, dll or .ico file and gives the app's
// launcher its icon. It returns false when the dialog was cancelled.
func (a *App) ChangeIcon(dir string) (bool, error) {
	start := ""
	if p, err := load(dir); err == nil && p.AppSource != "" {
		start = core.ExpandVars(p.AppSource, p.Dir)
	}
	from, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Take the icon from", DefaultDirectory: start,
		Filters: []runtime.FileFilter{{DisplayName: "Icons and programs (*.ico, *.exe, *.dll)", Pattern: "*.ico;*.exe;*.dll"}},
	})
	if err != nil || from == "" {
		return false, err
	}
	return true, a.changeIcon(dir, from)
}

func (a *App) changeIcon(dir, from string) error {
	p, err := load(dir)
	if err != nil {
		return err
	}
	return manager.SetIcon(p, from)
}

// UseAppIcon gives the launcher the app's own icon again.
func (a *App) UseAppIcon(dir string) error {
	p, err := load(dir)
	if err != nil {
		return err
	}
	return manager.UseAppIcon(p)
}

// Open shows a folder in Explorer.
func (a *App) Open(dir string) { exec.Command("explorer.exe", dir).Start() }
