// Package icon gives a launcher the icon of the app it starts, by copying the
// app exe's first icon group (RT_GROUP_ICON and its RT_ICON images) into the
// launcher's resources.
package icon

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

// ErrNoIcon means the exe has no icon.
var ErrNoIcon = errors.New("the program has no icon")

var (
	kernel32            = windows.NewLazySystemDLL("kernel32.dll")
	enumResourceNames   = kernel32.NewProc("EnumResourceNamesW")
	beginUpdateResource = kernel32.NewProc("BeginUpdateResourceW")
	updateResource      = kernel32.NewProc("UpdateResourceW")
	endUpdateResource   = kernel32.NewProc("EndUpdateResourceW")
)

const asData = 0x2 | 0x20 // LOAD_LIBRARY_AS_DATAFILE | LOAD_LIBRARY_AS_IMAGE_RESOURCE

// The callback is made once: Windows callbacks are never freed.
var (
	enumMu    sync.Mutex
	enumNames []windows.ResourceIDOrString
	enumCB    = windows.NewCallback(func(module, typ, name, param uintptr) uintptr {
		if name <= 0xFFFF {
			enumNames = append(enumNames, windows.ResourceID(name))
		} else {
			// A string owned by Windows, read only during the call.
			enumNames = append(enumNames, windows.UTF16PtrToString(*(**uint16)(unsafe.Pointer(&name))))
		}
		return 1 // go on
	})
)

func groupNames(module windows.Handle) []windows.ResourceIDOrString {
	enumMu.Lock()
	defer enumMu.Unlock()
	enumNames = nil
	enumResourceNames.Call(uintptr(module), uintptr(windows.RT_GROUP_ICON), enumCB, 0)
	return enumNames
}

// group is one icon group: the 6-byte header, then 14 bytes per image whose
// last 2 bytes are the image's RT_ICON id.
type group []byte

func (g group) count() int {
	if len(g) < 6 {
		return 0
	}
	n := int(binary.LittleEndian.Uint16(g[4:]))
	if len(g) < 6+14*n {
		return 0
	}
	return n
}

func (g group) id(i int) uint16        { return binary.LittleEndian.Uint16(g[6+14*i+12:]) }
func (g group) setID(i int, id uint16) { binary.LittleEndian.PutUint16(g[6+14*i+12:], id) }

// load reads the icon of an .ico file, or the first icon of an exe or dll.
func load(path string) (group, [][]byte, error) {
	if strings.EqualFold(filepath.Ext(path), ".ico") {
		return loadICO(path)
	}
	return loadPE(path)
}

// loadICO reads an .ico file: a 6-byte header, then 16 bytes per image (the
// same 12 first bytes as in a group, then where the image is in the file).
func loadICO(path string) (group, [][]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	bad := fmt.Errorf("%s is not an icon file", path)
	if len(b) < 6 || binary.LittleEndian.Uint16(b[0:]) != 0 || binary.LittleEndian.Uint16(b[2:]) != 1 {
		return nil, nil, bad
	}
	n := int(binary.LittleEndian.Uint16(b[4:]))
	if n == 0 || len(b) < 6+16*n {
		return nil, nil, bad
	}
	g := make(group, 6+14*n)
	copy(g, b[:6])
	images := make([][]byte, n)
	for i := 0; i < n; i++ {
		e := b[6+16*i : 6+16*i+16]
		size := int(binary.LittleEndian.Uint32(e[8:]))
		off := int(binary.LittleEndian.Uint32(e[12:]))
		if size == 0 || off < 6+16*n || off+size > len(b) {
			return nil, nil, bad
		}
		copy(g[6+14*i:], e[:12])
		g.setID(i, uint16(i+1))
		images[i] = append([]byte{}, b[off:off+size]...)
	}
	return g, images, nil
}

// SaveICO writes the icon of from (an exe, dll or .ico) as the .ico file to.
func SaveICO(from, to string) error {
	g, images, err := load(from)
	if err != nil {
		return err
	}
	n := len(images)
	out := make([]byte, 6+16*n)
	copy(out, g[:6])
	off := len(out)
	for i, img := range images {
		e := out[6+16*i:]
		copy(e, g[6+14*i:6+14*i+12])
		binary.LittleEndian.PutUint32(e[8:], uint32(len(img)))
		binary.LittleEndian.PutUint32(e[12:], uint32(off))
		off += len(img)
	}
	for _, img := range images {
		out = append(out, img...)
	}
	return os.WriteFile(to, out, 0o644)
}

// loadPE reads exe's first icon group and its images.
func loadPE(exe string) (group, [][]byte, error) {
	module, err := windows.LoadLibraryEx(exe, 0, asData)
	if err != nil {
		return nil, nil, fmt.Errorf("open %s: %w", exe, err)
	}
	defer windows.FreeLibrary(module)
	names := groupNames(module)
	if len(names) == 0 {
		return nil, nil, ErrNoIcon
	}
	res, err := windows.FindResource(module, names[0], windows.RT_GROUP_ICON)
	if err != nil {
		return nil, nil, err
	}
	data, err := windows.LoadResourceData(module, res)
	if err != nil {
		return nil, nil, err
	}
	g := group(append([]byte{}, data...))
	if g.count() == 0 {
		return nil, nil, ErrNoIcon
	}
	images := make([][]byte, g.count())
	for i := range images {
		r, err := windows.FindResource(module, windows.ResourceID(g.id(i)), windows.RT_ICON)
		if err != nil {
			return nil, nil, fmt.Errorf("icon image %d: %w", g.id(i), err)
		}
		d, err := windows.LoadResourceData(module, r)
		if err != nil {
			return nil, nil, err
		}
		images[i] = append([]byte{}, d...)
	}
	return g, images, nil
}

// ErrNoPNG means the icon has only older bitmap images, none in PNG.
var ErrNoPNG = errors.New("the icon has no PNG image")

// PNG returns the largest PNG image of the icon of path (an exe, dll or
// .ico), to show it. Most icons since Windows Vista have a 256×256 one.
func PNG(path string) ([]byte, error) {
	_, images, err := load(path)
	if err != nil {
		return nil, err
	}
	var best []byte
	for _, img := range images {
		if len(img) > 8 && string(img[:4]) == "\x89PNG" && len(img) > len(best) {
			best = img
		}
	}
	if best == nil {
		return nil, ErrNoPNG
	}
	return best, nil
}

// Read returns the images of exe's first icon, or ErrNoIcon.
func Read(exe string) ([][]byte, error) {
	_, images, err := load(exe)
	return images, err
}

func groupCount(exe string) (int, error) {
	module, err := windows.LoadLibraryEx(exe, 0, asData)
	if err != nil {
		return 0, err
	}
	defer windows.FreeLibrary(module)
	return len(groupNames(module)), nil
}

// Remove takes the icon a launcher was given off again (its first icon group
// and the images in it). An exe without an icon is left as it is.
func Remove(exe string) error {
	module, err := windows.LoadLibraryEx(exe, 0, asData)
	if err != nil {
		return fmt.Errorf("open %s: %w", exe, err)
	}
	names := groupNames(module)
	windows.FreeLibrary(module)
	if len(names) == 0 {
		return nil
	}
	id, ok := names[0].(windows.ResourceID)
	if !ok {
		return errors.New("this program's icon was not given by RapTech Portable")
	}
	g, _, err := loadPE(exe)
	if err != nil {
		return err
	}
	name, err := windows.UTF16PtrFromString(exe)
	if err != nil {
		return err
	}
	h, _, callErr := beginUpdateResource.Call(uintptr(unsafe.Pointer(name)), 0)
	if h == 0 {
		return fmt.Errorf("open %s for writing (is it running?): %w", exe, callErr)
	}
	del := func(typ windows.ResourceID, rid uint16) error {
		if r, _, err := updateResource.Call(h, uintptr(typ), uintptr(rid), 0, 0, 0); r == 0 {
			return err
		}
		return nil
	}
	for i := 0; i < g.count(); i++ {
		if err := del(windows.RT_ICON, g.id(i)); err != nil {
			endUpdateResource.Call(h, 1)
			return fmt.Errorf("remove icon image: %w", err)
		}
	}
	if err := del(windows.RT_GROUP_ICON, uint16(id)); err != nil {
		endUpdateResource.Call(h, 1)
		return fmt.Errorf("remove icon: %w", err)
	}
	if r, _, err := endUpdateResource.Call(h, 0); r == 0 {
		return fmt.Errorf("save %s: %w", exe, err)
	}
	return nil
}

// Copy gives to (a launcher) the icon of from (the app). to is changed only
// when from has an icon; to must not be running.
func Copy(from, to string) error {
	g, images, err := load(from)
	if err != nil {
		return err
	}
	name, err := windows.UTF16PtrFromString(to)
	if err != nil {
		return err
	}
	h, _, callErr := beginUpdateResource.Call(uintptr(unsafe.Pointer(name)), 0)
	if h == 0 {
		return fmt.Errorf("open %s for writing (is it running?): %w", to, callErr)
	}
	put := func(typ windows.ResourceID, id uint16, data []byte) error {
		if r, _, err := updateResource.Call(h, uintptr(typ), uintptr(id), 0, uintptr(unsafe.Pointer(&data[0])), uintptr(len(data))); r == 0 {
			return err
		}
		return nil
	}
	for i, img := range images {
		g.setID(i, uint16(i+1))
		if err := put(windows.RT_ICON, uint16(i+1), img); err != nil {
			endUpdateResource.Call(h, 1) // discard
			return fmt.Errorf("write icon image: %w", err)
		}
	}
	if err := put(windows.RT_GROUP_ICON, 1, g); err != nil {
		endUpdateResource.Call(h, 1)
		return fmt.Errorf("write icon: %w", err)
	}
	if r, _, err := endUpdateResource.Call(h, 0); r == 0 {
		return fmt.Errorf("save %s: %w", to, err)
	}
	return nil
}
