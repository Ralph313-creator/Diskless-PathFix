package manager

import (
	"errors"
	"path/filepath"

	"raptech-portable/internal/core"
	"raptech-portable/internal/icon"
)

// FixIcons gives every package's launcher in root that has no icon yet the
// icon of its app, such as Roblox once its first version is downloaded. It
// returns how many it gave an icon. A launcher that is running is skipped.
func FixIcons(root string) int {
	n := 0
	for _, info := range Scan(root) {
		if info.Error != "" {
			continue
		}
		if _, err := icon.Read(info.Launcher); !errors.Is(err, icon.ErrNoIcon) {
			continue // has one, or cannot be read
		}
		p, err := core.Load(info.Dir)
		if err != nil {
			continue
		}
		if ApplyIcon(p) == nil {
			if _, err := icon.Read(info.Launcher); err == nil {
				n++
			}
		}
	}
	return n
}

// ApplyIcon gives the package's launcher the icon of the app it starts.
// Nothing happens while the app is not there yet (Roblox before its first
// download) or has no icon.
func ApplyIcon(p *core.Package) error {
	if p.Exe == "" {
		return nil
	}
	exe, err := p.ResolveExe(p.Exe)
	if err != nil {
		return nil // not downloaded or copied yet
	}
	err = icon.Copy(exe, filepath.Join(p.Dir, LauncherName(p.Name)))
	if errors.Is(err, icon.ErrNoIcon) {
		return nil
	}
	return err
}
