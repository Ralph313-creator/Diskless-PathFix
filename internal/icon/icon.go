// Package icon gives a launcher the icon of the app it starts, by copying the
// app exe's first icon group (RT_GROUP_ICON and its RT_ICON images) into the
// launcher's resources.
package icon

import (
	"encoding/binary"
	"errors"
	"fmt"
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

// load reads exe's first icon group and its images.
func load(exe string) (group, [][]byte, error) {
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
