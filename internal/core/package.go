package core

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

// RegistryDir is the folder inside a package that holds saved .reg files.
const RegistryDir = "Registry"

// Setup makes this PC use the package: links each folder, loads the saved
// registry, and points the protocols at launcher (skipped when ""). It does
// as much as it can and returns every problem.
func (p *Package) Setup(launcher string) error {
	var errs []error
	for _, f := range p.Folders {
		if err := MakeLink(p.LinkPath(f), p.DataPath(f)); err != nil {
			errs = append(errs, err)
		}
	}
	if err := LoadRegistry(p.Registry, p.Path(RegistryDir), p.Dir); err != nil {
		errs = append(errs, err)
	}
	if launcher != "" {
		for _, proto := range p.Protocols {
			if err := SetProtocol(proto, launcher); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errors.Join(errs...)
}

// Finish runs after the app has closed: saves the registry into the package
// and empties the session folders, so a login never reaches the next
// customer.
func (p *Package) Finish() error {
	var errs []error
	if err := SaveRegistry(p.Registry, p.Path(RegistryDir), p.Dir); err != nil {
		errs = append(errs, err)
	}
	for _, f := range p.Folders {
		for _, s := range f.Session {
			if err := EmptyFolder(filepath.Join(p.DataPath(f), s)); err != nil {
				errs = append(errs, fmt.Errorf("empty %s: %w", s, err))
			}
		}
	}
	return errors.Join(errs...)
}

// Unlink removes the package's links and puts back what was there before.
// Protocol handlers are removed only while they still start launcher, so
// one the app's own installer took over again is left alone.
func (p *Package) Unlink(launcher string) error {
	var errs []error
	for _, f := range p.Folders {
		if err := RemoveLink(p.LinkPath(f), p.DataPath(f)); err != nil {
			errs = append(errs, err)
		}
	}
	if launcher != "" {
		for _, proto := range p.Protocols {
			if err := RemoveProtocol(proto, launcher); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errors.Join(errs...)
}

// StartPath is the path to start exe (a program in the package) with: through
// the link of the folder it is in when that link is made, so the app sees its
// default location (such as C:\Program Files\Roblox\Versions\...) and never
// the package; otherwise the package's own path.
func (p *Package) StartPath(exe string) string {
	best, bestLen := exe, -1
	for _, f := range p.Folders {
		data := p.DataPath(f)
		rel, err := filepath.Rel(data, exe)
		if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, `..\`) || len(data) <= bestLen {
			continue
		}
		if link := p.LinkPath(f); LinkStatus(link, data) == LinkOK {
			best, bestLen = filepath.Join(link, rel), len(data)
		}
	}
	return best
}

// FolderStatus is what is at one folder's link path now.
type FolderStatus struct {
	Link  string
	Data  string
	State LinkState
}

// Status reports each folder's link.
func (p *Package) Status() []FolderStatus {
	var out []FolderStatus
	for _, f := range p.Folders {
		link, data := p.LinkPath(f), p.DataPath(f)
		out = append(out, FolderStatus{Link: link, Data: data, State: LinkStatus(link, data)})
	}
	return out
}
