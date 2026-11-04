package core

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unicode/utf16"

	"golang.org/x/sys/windows"
)

// BackupSuffix is added to a real folder that was in the way of a link.
const BackupSuffix = ".rp-backup"

// LinkState is what is at a link path now.
type LinkState int

const (
	LinkMissing    LinkState = iota // nothing there
	LinkOK                          // a link to the package folder
	LinkOther                       // a link to somewhere else
	LinkRealFolder                  // a real folder, such as the app's own install
	LinkRealFile                    // a file
)

func (s LinkState) String() string {
	return [...]string{"missing", "ok", "other link", "real folder", "real file"}[s]
}

// LinkStatus says what is at link, compared with the package folder target.
func LinkStatus(link, target string) LinkState {
	info, err := os.Lstat(link)
	if err != nil {
		return LinkMissing
	}
	if isReparse(info) {
		if dest, err := os.Readlink(link); err == nil && SamePath(dest, target) {
			return LinkOK
		}
		return LinkOther
	}
	if info.IsDir() {
		return LinkRealFolder
	}
	return LinkRealFile
}

func isReparse(info os.FileInfo) bool {
	if d, ok := info.Sys().(*syscall.Win32FileAttributeData); ok {
		return d.FileAttributes&windows.FILE_ATTRIBUTE_REPARSE_POINT != 0
	}
	return info.Mode()&os.ModeSymlink != 0
}

// MakeLink makes link point to target, a folder in the package, creating it
// when missing. A real folder or file in the way is renamed to
// link+BackupSuffix, never deleted; if that backup already exists nothing is
// changed. A link to somewhere else is replaced (only the link is removed).
// A local target gets a junction, a network target a symbolic link.
func MakeLink(link, target string) error {
	link, target = filepath.Clean(link), filepath.Clean(target)
	if err := os.MkdirAll(target, 0o755); err != nil {
		return err
	}
	switch LinkStatus(link, target) {
	case LinkOK:
		return nil
	case LinkOther:
		if err := os.Remove(link); err != nil {
			return fmt.Errorf("remove old link %s: %w", link, err)
		}
	case LinkRealFolder, LinkRealFile:
		backup := link + BackupSuffix
		if _, err := os.Lstat(backup); err == nil {
			return fmt.Errorf("%s is in the way and %s already exists; move one of them yourself", link, backup)
		}
		if err := os.Rename(link, backup); err != nil {
			return fmt.Errorf("move %s aside (is the app running?): %w", link, err)
		}
	}
	if err := os.MkdirAll(filepath.Dir(link), 0o755); err != nil {
		return err
	}
	if IsUNC(target) {
		return os.Symlink(target, link)
	}
	return makeJunction(link, target)
}

// RemoveLink removes link if it is a link (never a real folder) and puts a
// backup made by MakeLink back.
func RemoveLink(link, target string) error {
	link = filepath.Clean(link)
	switch LinkStatus(link, target) {
	case LinkOK, LinkOther:
		if err := os.Remove(link); err != nil {
			return fmt.Errorf("remove link %s: %w", link, err)
		}
	case LinkRealFolder, LinkRealFile:
		return nil
	}
	backup := link + BackupSuffix
	if _, err := os.Lstat(backup); err == nil {
		return os.Rename(backup, link)
	}
	return nil
}

// makeJunction makes an NTFS junction (mount point) at link pointing to the
// local folder target. Junctions need no administrator rights.
func makeJunction(link, target string) error {
	target, err := filepath.Abs(target)
	if err != nil {
		return err
	}
	if err := os.Mkdir(link, 0o755); err != nil {
		return err
	}
	ok := false
	defer func() {
		if !ok {
			os.Remove(link)
		}
	}()
	sub := utf16.Encode([]rune(`\??\` + target))
	printName := utf16.Encode([]rune(target))
	// Path buffer: substitute name, NUL, print name, NUL (UTF-16).
	pathBytes := 2 * (len(sub) + 1 + len(printName) + 1)
	buf := make([]byte, 8+8+pathBytes)
	binary.LittleEndian.PutUint32(buf[0:], windows.IO_REPARSE_TAG_MOUNT_POINT)
	binary.LittleEndian.PutUint16(buf[4:], uint16(8+pathBytes))
	binary.LittleEndian.PutUint16(buf[8:], 0)                         // substitute offset
	binary.LittleEndian.PutUint16(buf[10:], uint16(2*len(sub)))       // substitute length
	binary.LittleEndian.PutUint16(buf[12:], uint16(2*(len(sub)+1)))   // print offset
	binary.LittleEndian.PutUint16(buf[14:], uint16(2*len(printName))) // print length
	pos := 16
	for _, u := range append(append(append(sub, 0), printName...), 0) {
		binary.LittleEndian.PutUint16(buf[pos:], u)
		pos += 2
	}

	name, err := windows.UTF16PtrFromString(link)
	if err != nil {
		return err
	}
	h, err := windows.CreateFile(name, windows.GENERIC_WRITE, 0, nil, windows.OPEN_EXISTING,
		windows.FILE_FLAG_BACKUP_SEMANTICS|windows.FILE_FLAG_OPEN_REPARSE_POINT, 0)
	if err != nil {
		return fmt.Errorf("open %s: %w", link, err)
	}
	defer windows.CloseHandle(h)
	var n uint32
	if err := windows.DeviceIoControl(h, windows.FSCTL_SET_REPARSE_POINT, &buf[0], uint32(len(buf)), nil, 0, &n, nil); err != nil {
		return fmt.Errorf("make junction %s: %w", link, err)
	}
	ok = true
	return nil
}
