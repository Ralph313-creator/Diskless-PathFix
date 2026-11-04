package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"golang.org/x/sys/windows"
)

// ManifestName is the file next to the launcher that describes a package.
const ManifestName = "portable.json"

// Manifest says where an app expects its folders and registry keys, and
// where they really live inside the package.
type Manifest struct {
	Name string `json:"name"`
	// Exe is the program to start, relative to the package. It may hold *,
	// as in Data\ProgramFiles\Roblox\Versions\*\RobloxPlayerBeta.exe; the
	// newest match wins.
	Exe  string   `json:"exe"`
	Args []string `json:"args,omitempty"`
	// Folders are linked from where the app expects them into the package.
	Folders []Folder `json:"folders,omitempty"`
	// Registry keys are imported before the app starts and saved after it
	// closes, such as HKCU\Software\ROBLOX Corporation.
	Registry []string `json:"registry,omitempty"`
	// Protocols such as roblox-player are pointed at the launcher, so links
	// from a browser start the app through the package.
	Protocols []string `json:"protocols,omitempty"`
	Update    *Update  `json:"update,omitempty"`
	// AppSource is the folder the app was copied from into App\, for
	// information (hunt shows it as already added). Never written to.
	AppSource string `json:"appSource,omitempty"`
}

// Folder is one folder the app uses.
type Folder struct {
	// Link is where the app expects the folder, with variables, such as
	// %LOCALAPPDATA%\Roblox.
	Link string `json:"link"`
	// Data is the folder inside the package, such as Data\LocalAppData\Roblox.
	Data string `json:"data"`
	// Session lists subfolders of Data emptied when the app closes, such as
	// a login store that must not pass to the next customer.
	Session []string `json:"session,omitempty"`
}

// Update says how the package keeps itself up to date.
type Update struct {
	// Kind is "roblox" (download from Roblox) or "run" (start Exe with Args
	// and let the app update itself).
	Kind string `json:"kind"`
	// Folder is the Versions folder inside the package (roblox).
	Folder string `json:"folder,omitempty"`
	// Keep is how many versions stay (roblox). 0 means 2.
	Keep int `json:"keep,omitempty"`
	// Every is how often the schedule runs it, such as "30m".
	Every   string   `json:"every,omitempty"`
	Exe     string   `json:"exe,omitempty"`
	Args    []string `json:"args,omitempty"`
	Timeout string   `json:"timeout,omitempty"`
}

// Package is a manifest and the folder it lives in.
type Package struct {
	Dir string
	Manifest
}

// Load reads and checks the manifest in dir.
func Load(dir string) (*Package, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	raw, err := os.ReadFile(filepath.Join(dir, ManifestName))
	if err != nil {
		return nil, err
	}
	p := &Package{Dir: dir}
	if err := json.Unmarshal(raw, &p.Manifest); err != nil {
		return nil, fmt.Errorf("%s: %w", ManifestName, err)
	}
	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", ManifestName, err)
	}
	return p, nil
}

// Save checks the manifest and writes it into the package folder.
func (p *Package) Save() error {
	if err := p.Validate(); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(p.Manifest, "", "  ")
	if err != nil {
		return err
	}
	tmp := filepath.Join(p.Dir, ManifestName+".tmp")
	if err := os.WriteFile(tmp, append(raw, '\n'), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Join(p.Dir, ManifestName))
}

// Validate checks every path, so nothing outside the package is written
// into and no Windows or user folder is ever moved aside.
func (m *Manifest) Validate() error {
	if strings.TrimSpace(m.Name) == "" {
		return errors.New("the package has no name")
	}
	if m.Exe != "" {
		if err := CheckPlace(m.Exe); err != nil {
			return fmt.Errorf("exe: %w", err)
		}
	}
	for i, f := range m.Folders {
		if err := CheckLink(f.Link); err != nil {
			return fmt.Errorf("folder %d: %w", i+1, err)
		}
		if err := CheckPlace(f.Data); err != nil {
			return fmt.Errorf("folder %d: %w", i+1, err)
		}
		if IsFull(f.Data) && nested(ExpandVars(f.Link, ""), ExpandVars(f.Data, "")) {
			return fmt.Errorf("folder %d: %s and %s are inside each other", i+1, f.Link, f.Data)
		}
		for _, s := range f.Session {
			if err := CheckData(s); err != nil {
				return fmt.Errorf("folder %d session: %w", i+1, err)
			}
		}
	}
	for _, key := range m.Registry {
		if _, _, err := splitKey(key); err != nil {
			return err
		}
	}
	for _, proto := range m.Protocols {
		if proto == "" || strings.ContainsAny(proto, `\/:" `) {
			return fmt.Errorf("protocol %q is not a name such as roblox-player", proto)
		}
	}
	if u := m.Update; u != nil {
		switch u.Kind {
		case "roblox":
			if err := CheckData(u.Folder); err != nil {
				return fmt.Errorf("update folder: %w", err)
			}
		case "run":
			if err := CheckPlace(u.Exe); err != nil {
				return fmt.Errorf("update exe: %w", err)
			}
		default:
			return fmt.Errorf("update kind %q is not roblox or run", u.Kind)
		}
		for _, d := range []string{u.Every, u.Timeout} {
			if d != "" {
				if _, err := time.ParseDuration(d); err != nil {
					return fmt.Errorf("update: %q is not a time such as 30m", d)
				}
			}
		}
	}
	return nil
}

// Path is a path inside the package, or a full path used in place as it is.
func (p *Package) Path(rel string) string {
	if IsFull(rel) {
		return filepath.Clean(ExpandVars(rel, p.Dir))
	}
	return filepath.Join(p.Dir, ExpandVars(rel, p.Dir))
}

// DataPath is where folder f's files really are: in the package, or the
// folder used in place.
func (p *Package) DataPath(f Folder) string { return p.Path(f.Data) }

// InPlace reports whether folder f is an existing folder used where it is
// (never copied into the package or deleted with it).
func (p *Package) InPlace(f Folder) bool { return IsFull(f.Data) }

// LinkPath is where the app expects folder f on this PC.
func (p *Package) LinkPath(f Folder) string {
	return filepath.Clean(ExpandVars(f.Link, p.Dir))
}

// ResolveExe finds the program to start. With * in the path the newest
// match wins: newest by the creation time of its folder, then of the file.
// Matches inside a folder starting with "." (a download in progress) are
// skipped.
func (p *Package) ResolveExe(rel string) (string, error) {
	path := p.Path(rel)
	if !strings.ContainsAny(rel, "*?[") {
		if _, err := os.Stat(path); err != nil {
			return "", err
		}
		return path, nil
	}
	matches, err := filepath.Glob(path)
	if err != nil {
		return "", err
	}
	type found struct {
		path      string
		dir, file time.Time
	}
	var list []found
	for _, m := range matches {
		rest, _ := filepath.Rel(p.Dir, m)
		if hiddenPart(rest) {
			continue
		}
		if info, err := os.Stat(m); err != nil || info.IsDir() {
			continue
		}
		list = append(list, found{m, Created(filepath.Dir(m)), Created(m)})
	}
	if len(list) == 0 {
		return "", fmt.Errorf("nothing matches %s", path)
	}
	sort.Slice(list, func(i, j int) bool {
		if !list[i].dir.Equal(list[j].dir) {
			return list[i].dir.After(list[j].dir)
		}
		return list[i].file.After(list[j].file)
	})
	return list[0].path, nil
}

func hiddenPart(rel string) bool {
	for _, part := range strings.Split(rel, `\`) {
		if strings.HasPrefix(part, ".") && part != "." && part != ".." {
			return true
		}
	}
	return false
}

// Created is a file or folder's creation time, or zero when it is missing.
func Created(path string) time.Time {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}
	}
	if d, ok := info.Sys().(*windows.Win32FileAttributeData); ok {
		return time.Unix(0, d.CreationTime.Nanoseconds())
	}
	return info.ModTime()
}
