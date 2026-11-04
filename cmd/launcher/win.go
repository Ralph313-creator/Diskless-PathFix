package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

// attachConsole lets the GUI-subsystem launcher print when it was started
// from a console (cmd, a script); it stays silent when double-clicked.
func attachConsole() {
	const attachParent = ^uint32(0) // ATTACH_PARENT_PROCESS
	proc := windows.NewLazySystemDLL("kernel32.dll").NewProc("AttachConsole")
	if r, _, _ := proc.Call(uintptr(attachParent)); r == 0 {
		return
	}
	if f, err := os.OpenFile("CONOUT$", os.O_WRONLY, 0); err == nil {
		os.Stdout, os.Stderr = f, f
	}
}

type shellExecuteInfo struct {
	cbSize         uint32
	fMask          uint32
	hwnd           windows.Handle
	lpVerb         *uint16
	lpFile         *uint16
	lpParameters   *uint16
	lpDirectory    *uint16
	nShow          int32
	hInstApp       windows.Handle
	lpIDList       uintptr
	lpClass        *uint16
	hkeyClass      windows.Handle
	dwHotKey       uint32
	hIconOrMonitor windows.Handle
	hProcess       windows.Handle
}

// runElevated starts exe with args through the UAC prompt and waits.
func runElevated(exe string, args []string) error {
	const seeMaskNoCloseProcess = 0x40
	verb, _ := windows.UTF16PtrFromString("runas")
	file, _ := windows.UTF16PtrFromString(exe)
	params, _ := windows.UTF16PtrFromString(windows.ComposeCommandLine(args))
	info := shellExecuteInfo{fMask: seeMaskNoCloseProcess, lpVerb: verb, lpFile: file, lpParameters: params, nShow: windows.SW_HIDE}
	info.cbSize = uint32(unsafe.Sizeof(info))
	proc := windows.NewLazySystemDLL("shell32.dll").NewProc("ShellExecuteExW")
	if r, _, err := proc.Call(uintptr(unsafe.Pointer(&info))); r == 0 {
		return err
	}
	defer windows.CloseHandle(info.hProcess)
	windows.WaitForSingleObject(info.hProcess, windows.INFINITE)
	var code uint32
	windows.GetExitCodeProcess(info.hProcess, &code)
	if code != 0 {
		return fmt.Errorf("elevated setup ended with %d", code)
	}
	return nil
}

// otherInstances counts other running copies of this launcher exe, so only
// the last one to close saves the registry and empties session folders.
func otherInstances(self string) int {
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return 0
	}
	defer windows.CloseHandle(snap)
	base := strings.ToLower(filepath.Base(self))
	me := uint32(os.Getpid())
	n := 0
	var e windows.ProcessEntry32
	e.Size = uint32(unsafe.Sizeof(e))
	for err = windows.Process32First(snap, &e); err == nil; err = windows.Process32Next(snap, &e) {
		if e.ProcessID == me || strings.ToLower(windows.UTF16ToString(e.ExeFile[:])) != base {
			continue
		}
		if path := imagePath(e.ProcessID); path != "" && strings.EqualFold(path, self) {
			n++
		}
	}
	return n
}

func imagePath(pid uint32) string {
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return ""
	}
	defer windows.CloseHandle(h)
	buf := make([]uint16, windows.MAX_LONG_PATH)
	size := uint32(len(buf))
	if err := windows.QueryFullProcessImageName(h, 0, &buf[0], &size); err != nil {
		return ""
	}
	return windows.UTF16ToString(buf[:size])
}
