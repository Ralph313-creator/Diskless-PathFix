package manager

import (
	"errors"
	"os"
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

// OwnIcon is the icon chosen for a package, kept in the package folder.
const OwnIcon = "icon.ico"

// HasOwnIcon reports whether an icon was chosen for the package, instead of
// the app's own.
func HasOwnIcon(p *core.Package) bool {
	_, err := os.Stat(filepath.Join(p.Dir, OwnIcon))
	return err == nil
}

// SetIcon gives the launcher the icon of from (an exe, dll or .ico) and keeps
// it in the package, so updates keep it too.
func SetIcon(p *core.Package, from string) error {
	own := filepath.Join(p.Dir, OwnIcon)
	tmp := own + ".tmp.ico"
	if err := icon.SaveICO(from, tmp); err != nil {
		os.Remove(tmp)
		return err
	}
	if err := icon.Copy(tmp, filepath.Join(p.Dir, LauncherName(p.Name))); err != nil {
		os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, own)
}

// UseAppIcon forgets the chosen icon and gives the launcher the app's again;
// when the app has none, the launcher has none either.
func UseAppIcon(p *core.Package) error {
	if err := os.Remove(filepath.Join(p.Dir, OwnIcon)); err != nil && !os.IsNotExist(err) {
		return err
	}
	if p.Exe != "" {
		if exe, err := p.ResolveExe(p.Exe); err == nil {
			if _, err := icon.Read(exe); err == nil {
				return ApplyIcon(p)
			}
		}
	}
	return icon.Remove(filepath.Join(p.Dir, LauncherName(p.Name)))
}

// ApplyIcon gives the package's launcher its icon: the chosen one, else
// the app's own. Nothing happens while the app is not there yet (Roblox
// before its first download) or has no icon.
func ApplyIcon(p *core.Package) error {
	if HasOwnIcon(p) {
		return icon.Copy(filepath.Join(p.Dir, OwnIcon), filepath.Join(p.Dir, LauncherName(p.Name)))
	}
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
