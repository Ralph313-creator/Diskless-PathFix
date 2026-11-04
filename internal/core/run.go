package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

// SetProtocol makes links such as roblox-player:... start exe, for this
// user (HKCU wins over HKLM), so they never hit an old version's path.
func SetProtocol(name, exe string) error {
	base := `Software\Classes\` + name
	set := func(path, value, data string) error {
		k, _, err := registry.CreateKey(registry.CURRENT_USER, path, registry.SET_VALUE)
		if err != nil {
			return fmt.Errorf("protocol %s: %w", name, err)
		}
		defer k.Close()
		return k.SetStringValue(value, data)
	}
	return errors.Join(
		set(base, "", "URL:"+name+" Protocol"),
		set(base, "URL Protocol", ""),
		set(base+`\DefaultIcon`, "", `"`+exe+`",0`),
		set(base+`\shell\open\command`, "", `"`+exe+`" "%1"`),
	)
}

// EmptyFolder deletes everything inside dir but keeps dir. A missing dir is
// fine.
func EmptyFolder(dir string) error {
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	var errs []error
	for _, e := range entries {
		if err := os.RemoveAll(filepath.Join(dir, e.Name())); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// Result is how a Run ended.
type Result struct {
	ExitCode uint32 // of the first program
	TimedOut bool   // the timeout ended every program
}

// Run starts exe and waits until it and every program it started have
// ended, using a job object (launchers often start the real app and quit).
// With a timeout above 0, all of them are ended when it passes.
func Run(exe string, args []string, dir string, timeout time.Duration) (Result, error) {
	var res Result
	if _, err := os.Stat(exe); err != nil {
		return res, err
	}
	if dir == "" {
		dir = filepath.Dir(exe)
	}
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return res, err
	}
	defer windows.CloseHandle(job)

	cmdline, err := windows.UTF16PtrFromString(windows.ComposeCommandLine(append([]string{exe}, args...)))
	if err != nil {
		return res, err
	}
	app, _ := windows.UTF16PtrFromString(exe)
	wdir, _ := windows.UTF16PtrFromString(dir)
	si := &windows.StartupInfo{Cb: uint32(unsafe.Sizeof(windows.StartupInfo{}))}
	var pi windows.ProcessInformation
	if err := windows.CreateProcess(app, cmdline, nil, nil, false, windows.CREATE_SUSPENDED, nil, wdir, si, &pi); err != nil {
		return res, fmt.Errorf("start %s: %w", exe, err)
	}
	defer windows.CloseHandle(pi.Process)
	defer windows.CloseHandle(pi.Thread)
	if err := windows.AssignProcessToJobObject(job, pi.Process); err != nil {
		windows.ResumeThread(pi.Thread)
		return res, fmt.Errorf("watch %s: %w", exe, err)
	}
	if _, err := windows.ResumeThread(pi.Thread); err != nil {
		windows.TerminateProcess(pi.Process, 1)
		return res, err
	}

	var deadline time.Time
	if timeout > 0 {
		deadline = time.Now().Add(timeout)
	}
	for activeProcesses(job) > 0 {
		if !deadline.IsZero() && time.Now().After(deadline) {
			res.TimedOut = true
			windows.TerminateJobObject(job, 1)
			windows.WaitForSingleObject(pi.Process, 5000)
			break
		}
		time.Sleep(250 * time.Millisecond)
	}
	windows.GetExitCodeProcess(pi.Process, &res.ExitCode)
	return res, nil
}

type jobAccounting struct {
	TotalUserTime, TotalKernelTime, ThisPeriodUserTime, ThisPeriodKernelTime int64
	TotalPageFaults, TotalProcesses, ActiveProcesses, TotalTerminated        uint32
}

func activeProcesses(job windows.Handle) uint32 {
	var info jobAccounting
	err := windows.QueryInformationJobObject(job, windows.JobObjectBasicAccountingInformation,
		uintptr(unsafe.Pointer(&info)), uint32(unsafe.Sizeof(info)), nil)
	if err != nil {
		return 0
	}
	return info.ActiveProcesses
}
