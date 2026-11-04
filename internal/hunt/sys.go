package hunt

import (
	"encoding/binary"
	"errors"
	"runtime"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

// process is one running process.
type process struct {
	pid, parent uint32
	name        string // its program's file name
}

func processList() []process {
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil
	}
	defer windows.CloseHandle(snap)
	var list []process
	e := windows.ProcessEntry32{Size: uint32(unsafe.Sizeof(windows.ProcessEntry32{}))}
	for err := windows.Process32First(snap, &e); err == nil; err = windows.Process32Next(snap, &e) {
		list = append(list, process{e.ProcessID, e.ParentProcessID, windows.UTF16ToString(e.ExeFile[:])})
	}
	return list
}

func processCreated(pid uint32) (windows.Filetime, bool) {
	p, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return windows.Filetime{}, false
	}
	defer windows.CloseHandle(p)
	var created, exit, kernel, user windows.Filetime
	if windows.GetProcessTimes(p, &created, &exit, &kernel, &user) != nil {
		return windows.Filetime{}, false
	}
	return created, true
}

func processPath(pid uint32) string {
	p, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return ""
	}
	defer windows.CloseHandle(p)
	return imagePath(p)
}

func imagePath(p windows.Handle) string {
	buf := make([]uint16, windows.MAX_LONG_PATH)
	size := uint32(len(buf))
	if windows.QueryFullProcessImageName(p, 0, &buf[0], &size) != nil {
		return ""
	}
	return windows.UTF16ToString(buf[:size])
}

// moduleFiles lists the files a process has loaded, such as its DLLs.
func moduleFiles(p windows.Handle) []string {
	var mods [2048]windows.Handle
	var needed uint32
	size := uint32(unsafe.Sizeof(mods[0]))
	if windows.EnumProcessModulesEx(p, &mods[0], uint32(len(mods))*size, &needed, windows.LIST_MODULES_ALL) != nil {
		return nil
	}
	buf := make([]uint16, windows.MAX_LONG_PATH)
	var files []string
	for _, m := range mods[:min(int(needed/size), len(mods))] {
		if windows.GetModuleFileNameEx(p, m, &buf[0], uint32(len(buf))) == nil {
			files = append(files, windows.UTF16ToString(buf))
		}
	}
	return files
}

// watchDir calls changed each time a file or folder in dir, or below it, is
// added, deleted, renamed or written, until stop is closed, with the paths
// that changed, relative to dir (nil when more changed at once than Windows
// could keep track of).
func watchDir(dir string, stop <-chan struct{}, changed func(names []string)) error {
	path, err := windows.UTF16PtrFromString(dir)
	if err != nil {
		return err
	}
	h, err := windows.CreateFile(path, windows.FILE_LIST_DIRECTORY,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE, nil,
		windows.OPEN_EXISTING, windows.FILE_FLAG_BACKUP_SEMANTICS|windows.FILE_FLAG_OVERLAPPED, 0)
	if err != nil {
		return err
	}
	defer windows.CloseHandle(h)
	done, err := windows.CreateEvent(nil, 1, 0, nil)
	if err != nil {
		return err
	}
	defer windows.CloseHandle(done)
	stopped, err := windows.CreateEvent(nil, 1, 0, nil)
	if err != nil {
		return err
	}
	defer windows.CloseHandle(stopped)
	returned := make(chan struct{})
	defer close(returned)
	go func() {
		select {
		case <-stop:
			windows.SetEvent(stopped)
		case <-returned:
		}
	}()

	// Windows writes into buf and ov after ReadDirectoryChanges returns, so
	// they are pinned: otherwise they could move with the goroutine's stack.
	buf := make([]byte, 256*1024) // big, as whole drives are watched
	ov := new(windows.Overlapped)
	var pin runtime.Pinner
	pin.Pin(&buf[0])
	pin.Pin(ov)
	defer pin.Unpin()
	ov.HEvent = done
	const filter = windows.FILE_NOTIFY_CHANGE_FILE_NAME | windows.FILE_NOTIFY_CHANGE_DIR_NAME |
		windows.FILE_NOTIFY_CHANGE_SIZE | windows.FILE_NOTIFY_CHANGE_LAST_WRITE
	for {
		windows.ResetEvent(done)
		if err := windows.ReadDirectoryChanges(h, &buf[0], uint32(len(buf)), true, filter, nil, ov, 0); err != nil {
			return err
		}
		event, err := windows.WaitForMultipleObjects([]windows.Handle{done, stopped}, false, windows.INFINITE)
		var n uint32
		if err != nil || event != windows.WAIT_OBJECT_0 {
			windows.CancelIoEx(h, ov)
			windows.GetOverlappedResult(h, ov, &n, true)
			return err
		}
		if err := windows.GetOverlappedResult(h, ov, &n, false); err != nil {
			return err
		}
		if n == 0 && !isDir(dir) {
			return errors.New(dir + " is gone")
		}
		changed(changedNames(buf[:n]))
	}
}

// changedNames reads the FILE_NOTIFY_INFORMATION records ReadDirectoryChanges wrote.
func changedNames(buf []byte) []string {
	var names []string
	for len(buf) >= 12 {
		next := int(binary.LittleEndian.Uint32(buf[0:]))
		size := int(binary.LittleEndian.Uint32(buf[8:]))
		if 12+size > len(buf) {
			break
		}
		words := make([]uint16, size/2)
		for i := range words {
			words[i] = binary.LittleEndian.Uint16(buf[12+2*i:])
		}
		names = append(names, string(utf16.Decode(words)))
		if next <= 0 || next > len(buf) {
			break
		}
		buf = buf[next:]
	}
	return names
}
