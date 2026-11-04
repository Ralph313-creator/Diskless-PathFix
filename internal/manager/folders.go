package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"raptech-portable/internal/core"
)

// linkedElsewhere refuses a folder that is already a link, such as one made
// by another app's package: copying it would copy that app's files, and
// linking it would take it from that app.
func linkedElsewhere(real string) error {
	if core.LinkStatus(real, "") != core.LinkOther {
		return nil
	}
	target, _ := os.Readlink(real)
	target = strings.TrimPrefix(strings.TrimPrefix(target, `\??\`), `\\?\`)
	for dir := target; dir != filepath.Dir(dir); dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, core.ManifestName)); err == nil {
			owner := filepath.Base(dir)
			return fmt.Errorf("%s is already linked by the %s app (to %s); remove %s first, or take the folder out of it", real, owner, target, owner)
		}
	}
	return fmt.Errorf("%s is a link to %s; it can't be added", real, target)
}

// SetSession sets which subfolders of folder i are emptied after each
// player, such as LocalStorage for logins.
func SetSession(p *core.Package, i int, subs []string) error {
	if i < 0 || i >= len(p.Folders) {
		return fmt.Errorf("the package has no folder %d", i+1)
	}
	var clean []string
	for _, s := range subs {
		if s = strings.TrimSpace(s); s != "" {
			clean = append(clean, s)
		}
	}
	old := p.Folders[i].Session
	p.Folders[i].Session = clean
	if err := p.Save(); err != nil {
		p.Folders[i].Session = old
		return err
	}
	return nil
}
