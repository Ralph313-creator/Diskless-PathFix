package manager

import (
	"errors"
	"fmt"
	"strings"

	"raptech-portable/internal/core"
)

// Added lists the folders on this PC the package already covers: each
// folder's link path, then the folder the app was copied from. Hunt shows
// these as added.
func Added(p *core.Package) []string {
	var out []string
	for _, f := range p.Folders {
		out = append(out, p.LinkPath(f))
	}
	if p.AppSource != "" {
		out = append(out, core.ExpandVars(p.AppSource, p.Dir))
	}
	return out
}

// AddFolders adds each folder found by hunt (see AddFolder), skipping ones
// the package already has. It returns how many it added and every problem.
func AddFolders(p *core.Package, links []string) (int, error) {
	added := 0
	var problems []error
next:
	for _, link := range links {
		for _, have := range Added(p) {
			if core.SamePath(have, core.ExpandVars(link, p.Dir)) {
				continue next
			}
		}
		if err := AddFolder(p, link, nil); err != nil {
			problems = append(problems, fmt.Errorf("%s: %w", link, err))
			continue
		}
		added++
	}
	return added, errors.Join(problems...)
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
