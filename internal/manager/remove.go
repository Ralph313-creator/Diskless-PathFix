package manager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"

	"raptech-portable/internal/core"
)

// Program is a running program from a package.
type Program struct {
	Name string
	PID  uint32
	Path string
}

// Use is what is using a package's files right now.
type Use struct {
	Programs []Program // running from the package, such as the app or its launcher
	Services []string  // Windows services whose program is in the package
}

// InUse lists the programs running from the package and the Windows
// services installed from it (an app such as STOVE registers its own).
func InUse(p *core.Package) Use {
	var use Use
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err == nil {
		defer windows.CloseHandle(snap)
		me := uint32(os.Getpid())
		e := windows.ProcessEntry32{Size: uint32(unsafe.Sizeof(windows.ProcessEntry32{}))}
		for err := windows.Process32First(snap, &e); err == nil; err = windows.Process32Next(snap, &e) {
			if e.ProcessID == me || e.ProcessID == 0 {
				continue
			}
			if path := imagePath(e.ProcessID); path != "" && within(path, p.Dir) {
				use.Programs = append(use.Programs, Program{Name: windows.UTF16ToString(e.ExeFile[:]), PID: e.ProcessID, Path: path})
			}
		}
	}
	m, err := mgr.Connect()
	if err != nil {
		return use // not an administrator: services can't be listed
	}
	defer m.Disconnect()
	names, _ := m.ListServices()
	for _, name := range names {
		s, err := m.OpenService(name)
		if err != nil {
			continue
		}
		cfg, err := s.Config()
		s.Close()
		if err == nil && within(serviceExe(cfg.BinaryPathName), p.Dir) {
			use.Services = append(use.Services, name)
		}
	}
	return use
}

// serviceExe is the program in a service's command line: "C:\x y\a.exe" -k
// or C:\x\a.exe -k.
func serviceExe(cmdline string) string {
	cmdline = strings.TrimSpace(cmdline)
	if strings.HasPrefix(cmdline, `"`) {
		if end := strings.Index(cmdline[1:], `"`); end >= 0 {
			return cmdline[1 : end+1]
		}
		return strings.Trim(cmdline, `"`)
	}
	if i := strings.Index(strings.ToLower(cmdline), ".exe"); i >= 0 {
		return cmdline[:i+4]
	}
	return cmdline
}

func within(path, dir string) bool {
	if path == "" {
		return false
	}
	rel, err := filepath.Rel(dir, strings.TrimPrefix(path, `\\?\`))
	return err == nil && rel != "." && rel != ".." && !strings.HasPrefix(rel, `..\`)
}

func imagePath(pid uint32) string {
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return ""
	}
	defer windows.CloseHandle(h)
	buf := make([]uint16, windows.MAX_LONG_PATH)
	size := uint32(len(buf))
	if windows.QueryFullProcessImageName(h, 0, &buf[0], &size) != nil {
		return ""
	}
	return windows.UTF16ToString(buf[:size])
}

// RemovePackage removes the app p from root: it stops and unregisters the
// Windows services installed from the package, closes the programs running
// from it, removes its links and protocol handlers on this PC (putting the
// original folders back), deletes its update task, and deletes the package
// folder. Only a folder with portable.json directly inside root is ever
// deleted, and links inside it are removed, not followed.
func RemovePackage(root string, p *core.Package) error {
	dir := filepath.Clean(p.Dir)
	if !core.SamePath(filepath.Dir(dir), root) || core.SamePath(dir, root) {
		return fmt.Errorf("%s is not an app in %s", dir, root)
	}
	if _, err := os.Stat(filepath.Join(dir, core.ManifestName)); err != nil {
		return fmt.Errorf("%s has no %s, so it is not an app folder", dir, core.ManifestName)
	}

	use := InUse(p)
	if err := removeServices(use.Services); err != nil {
		return err
	}
	if err := closePrograms(use.Programs); err != nil {
		return err
	}
	if err := p.Unlink(filepath.Join(dir, LauncherName(p.Name))); err != nil {
		return fmt.Errorf("remove its links first: %w", err)
	}
	if err := RemoveSchedule(p.Name); err != nil {
		return err
	}
	// Closed programs let go of their files a moment later.
	var err error
	for try := 0; try < 10; try++ {
		if err = os.RemoveAll(dir); err == nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("delete %s: %w", dir, err)
}

func removeServices(names []string) error {
	if len(names) == 0 {
		return nil
	}
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("services: %w", err)
	}
	defer m.Disconnect()
	var errs []error
	for _, name := range names {
		s, err := m.OpenService(name)
		if err != nil {
			errs = append(errs, fmt.Errorf("service %s: %w", name, err))
			continue
		}
		if st, err := s.Query(); err == nil && st.State != svc.Stopped {
			s.Control(svc.Stop)
			for i := 0; i < 60; i++ {
				if st, err := s.Query(); err != nil || st.State == svc.Stopped {
					break
				}
				time.Sleep(500 * time.Millisecond)
			}
		}
		if err := s.Delete(); err != nil {
			errs = append(errs, fmt.Errorf("unregister service %s: %w", name, err))
		}
		s.Close()
	}
	return errors.Join(errs...)
}

func closePrograms(list []Program) error {
	var errs []error
	for _, prog := range list {
		h, err := windows.OpenProcess(windows.PROCESS_TERMINATE|windows.SYNCHRONIZE, false, prog.PID)
		if err != nil {
			if imagePath(prog.PID) == "" {
				continue // already gone
			}
			errs = append(errs, fmt.Errorf("close %s: %w", prog.Name, err))
			continue
		}
		windows.TerminateProcess(h, 1)
		windows.WaitForSingleObject(h, 5000)
		windows.CloseHandle(h)
	}
	return errors.Join(errs...)
}
