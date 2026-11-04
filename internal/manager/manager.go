// Package manager makes and lists packages; the Wails app is a thin wrapper.
package manager

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"raptech-portable/internal/core"
)

// LauncherName is the launcher's file name inside a package.
func LauncherName(name string) string { return name + ".exe" }

var reserved = regexp.MustCompile(`(?i)^(con|prn|aux|nul|com[0-9]|lpt[0-9])$`)

func checkName(name string) error {
	if strings.TrimSpace(name) == "" || name == "." || name == ".." ||
		strings.ContainsAny(name, `\/:*?"<>|`) || reserved.MatchString(name) {
		return fmt.Errorf("%q cannot be a package name", name)
	}
	return nil
}

func newPackage(root, name string, launcher []byte) (*core.Package, error) {
	if err := checkName(name); err != nil {
		return nil, err
	}
	dir := filepath.Join(root, name)
	if _, err := os.Stat(dir); err == nil {
		return nil, fmt.Errorf("%s already exists", dir)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(dir, LauncherName(name)), launcher, 0o755); err != nil {
		return nil, err
	}
	return &core.Package{Dir: dir, Manifest: core.Manifest{Name: name}}, nil
}

// CreateRoblox makes root\Roblox for Roblox Player. perUser is a Roblox
// installed for one user (%LOCALAPPDATA%\Roblox\Versions); otherwise it is
// installed for all users (%ProgramFiles%\Roblox\Versions), like this
// server. Logins, logs and downloads are emptied after each player.
func CreateRoblox(root string, perUser bool, launcher []byte) (*core.Package, error) {
	p, err := newPackage(root, "Roblox", launcher)
	if err != nil {
		return nil, err
	}
	local := core.Folder{Link: `%LOCALAPPDATA%\Roblox`, Data: `Data\LocalAppData\Roblox`, Session: []string{"LocalStorage", "logs", "Downloads"}}
	versions := `Data\ProgramFiles\Roblox\Versions`
	if perUser {
		// Versions sits inside the linked %LOCALAPPDATA%\Roblox.
		versions = `Data\LocalAppData\Roblox\Versions`
		p.Folders = []core.Folder{local}
	} else {
		p.Folders = []core.Folder{{Link: `%ProgramFiles%\Roblox\Versions`, Data: versions}, local}
	}
	p.Exe = versions + `\*\RobloxPlayerBeta.exe`
	p.Registry = []string{`HKCU\Software\ROBLOX Corporation`}
	p.Protocols = []string{"roblox-player", "roblox"}
	p.Update = &core.Update{Kind: "roblox", Folder: versions, Keep: 2, Every: "30m"}
	return p, p.Save()
}

// CreatePackage makes root\name. With appDir, that folder is copied into
// the package's App folder and exe (inside appDir) is what starts.
func CreatePackage(root, name, appDir, exe string, launcher []byte) (*core.Package, error) {
	if err := checkName(name); err != nil {
		return nil, err
	}
	var rel string
	if appDir != "" {
		var err error
		rel, err = filepath.Rel(appDir, exe)
		if err != nil || strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
			return nil, fmt.Errorf("%s is not inside %s", exe, appDir)
		}
	}
	p, err := newPackage(root, name, launcher)
	if err != nil {
		return nil, err
	}
	if appDir != "" {
		if err := Copy(appDir, filepath.Join(p.Dir, "App"), nil); err != nil {
			return nil, err
		}
		p.Exe = `App\` + rel
	}
	return p, p.Save()
}

// AddFolder adds a folder the app uses. link is where the app expects it;
// what is there now is copied into the package, without the session
// subfolders (such as logins), and the original is left as it is.
func AddFolder(p *core.Package, link string, session []string) error {
	link = core.Portable(filepath.Clean(strings.TrimSpace(link)))
	if err := core.CheckLink(link); err != nil {
		return err
	}
	real := core.ExpandVars(link, p.Dir)
	for _, f := range p.Folders {
		if core.SamePath(core.ExpandVars(f.Link, p.Dir), real) {
			return fmt.Errorf("%s is already in the package", link)
		}
	}
	f := core.Folder{Link: link, Data: dataFor(link), Session: session}
	if info, err := os.Stat(real); err == nil && info.IsDir() {
		if err := Copy(real, p.Path(f.Data), session); err != nil {
			return err
		}
	}
	p.Folders = append(p.Folders, f)
	return p.Save()
}

// dataFor is where a link path's folder goes inside the package:
// %LOCALAPPDATA%\Roblox -> Data\LocalAppData\Roblox, D:\Games\X -> Data\D\Games\X.
func dataFor(link string) string {
	if strings.HasPrefix(link, `\\`) {
		return `Data\Network\` + strings.TrimLeft(link, `\`)
	}
	if strings.HasPrefix(link, "%") {
		if end := strings.Index(link[1:], "%"); end > 0 {
			name := strings.NewReplacer("(", "", ")", "").Replace(link[1 : end+1])
			for _, pretty := range []string{"LocalAppData", "AppData", "UserProfile", "Public", "ProgramData", "ProgramFilesx86", "ProgramFiles"} {
				if strings.EqualFold(name, pretty) {
					name = pretty
				}
			}
			return `Data\` + name + link[end+2:]
		}
	}
	return `Data\` + strings.Replace(link, `:`, "", 1)
}

// Copy makes dest the same as src (robocopy /MIR) without the skip
// subfolders. src is never changed.
func Copy(src, dest string, skip []string) error {
	args := []string{src, dest, "/MIR", "/COPY:DAT", "/DCOPY:DAT", "/R:1", "/W:1", "/NFL", "/NDL", "/NJH", "/NJS", "/NP"}
	if len(skip) > 0 {
		args = append(args, "/XD")
		for _, s := range skip {
			args = append(args, filepath.Join(src, s))
		}
	}
	cmd := exec.Command("robocopy.exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	out, err := cmd.CombinedOutput()
	var exit *exec.ExitError
	if errors.As(err, &exit) && exit.ExitCode() < 8 {
		return nil // 1-7 mean files were copied or extra ones removed
	}
	if err != nil {
		return fmt.Errorf("copy %s: %v %s", src, err, strings.TrimSpace(string(out)))
	}
	return nil
}

// Info is one package as the manager shows it.
type Info struct {
	Name     string
	Dir      string
	Launcher string
	Version  string // newest version folder for Roblox, else ""
	Folders  []core.FolderStatus
	Update   *core.Update
	LastLog  string
	Error    string
}

// Scan lists the packages in root's subfolders.
func Scan(root string) []Info {
	entries, _ := os.ReadDir(root)
	var out []Info
	for _, e := range entries {
		dir := filepath.Join(root, e.Name())
		if _, err := os.Stat(filepath.Join(dir, core.ManifestName)); err != nil {
			continue
		}
		p, err := core.Load(dir)
		if err != nil {
			out = append(out, Info{Name: e.Name(), Dir: dir, Error: err.Error()})
			continue
		}
		info := Info{Name: p.Name, Dir: dir, Launcher: filepath.Join(dir, LauncherName(p.Name)),
			Folders: p.Status(), Update: p.Update, LastLog: lastLine(filepath.Join(dir, "RapPortable.log"))}
		if p.Update != nil && p.Update.Kind == "roblox" && p.Exe != "" {
			if exe, err := p.ResolveExe(p.Exe); err == nil {
				info.Version = filepath.Base(filepath.Dir(exe))
			}
		}
		out = append(out, info)
	}
	return out
}

func lastLine(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	if len(b) > 4096 {
		b = b[len(b)-4096:]
	}
	lines := strings.Split(strings.TrimSpace(string(b)), "\n")
	return strings.TrimSpace(lines[len(lines)-1])
}
