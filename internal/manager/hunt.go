package manager

import (
	"errors"
	"fmt"

	"raptech-portable/internal/core"
)

// Added lists the folders on this PC the package already covers: each
// folder's link path (the app's own folder among them), then the folder the
// app was copied from if it is not one of those (older packages). Hunt shows
// these as added.
func Added(p *core.Package) []string {
	var out []string
	for _, f := range p.Folders {
		out = append(out, p.LinkPath(f))
	}
	if p.AppSource != "" {
		src := core.ExpandVars(p.AppSource, p.Dir)
		for _, have := range out {
			if core.SamePath(have, src) {
				return out
			}
		}
		out = append(out, src)
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
