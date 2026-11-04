package hunt

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/sys/windows"
)

// defaultRoots are the folders whose folders are suggested: an app's folder
// is one level below the deepest one a path is in, such as
// %LOCALAPPDATA%\Riot Games for %LOCALAPPDATA%\Riot Games\Logs\x.log. A path
// in none of them is one level below its drive, such as C:\Riot Games.
func defaultRoots() []string {
	env := func(name string) string { return strings.TrimRight(os.Getenv(name), `\`) }
	var roots []string
	add := func(parts ...string) {
		if parts[0] != "" {
			roots = append(roots, filepath.Join(parts...))
		}
	}
	profile := env("USERPROFILE")
	add(profile)
	add(profile, "AppData", "LocalLow")
	add(profile, "Documents")
	add(profile, "Documents", "My Games")
	add(profile, "Saved Games")
	for _, name := range []string{"LOCALAPPDATA", "APPDATA", "PUBLIC", "ProgramData", "ProgramFiles", "ProgramFiles(x86)", "SystemRoot"} {
		add(env(name))
	}
	add(env("PUBLIC"), "Documents")
	add(env("LOCALAPPDATA"), "Programs")
	add(env("ProgramFiles(x86)"), "Steam", "steamapps", "common")
	add(env("ProgramFiles"), "Epic Games")
	return roots
}

// containerNames are folders that hold apps, rather than being one: the
// folder suggested is the one in them, such as D:\Games\VALORANT.
var containerNames = []string{"games", "steamlibrary", "steamapps", "common", "epic games", "xboxgames", "my games"}

// suggestFolder returns the folder to suggest for path: the one below the
// deepest of roots that path is in, or below its drive. It is "" for a path
// that is a root or a drive itself.
func suggestFolder(path string, roots []string) string {
	path = filepath.Clean(path)
	best := ""
	for _, r := range roots {
		if strings.EqualFold(path, r) {
			return ""
		}
		if len(r) > len(best) && inside(path, r) {
			best = r
		}
	}
	if best == "" {
		if best = filepath.VolumeName(path); best == "" {
			return ""
		}
		best += `\`
	}
	rel, err := filepath.Rel(best, path)
	if err != nil || rel == "." || strings.HasPrefix(rel, "..") {
		return ""
	}
	parts := strings.Split(rel, `\`)
	n := 1
	for n < len(parts) && slices.Contains(containerNames, strings.ToLower(parts[n-1])) {
		n++
	}
	return filepath.Join(append([]string{best}, parts[:n]...)...)
}

// noiseNames and noisePrefixes are folders that most likely belong to
// Windows or another program, not to the app.
var (
	noiseNames = []string{
		"temp", "tmp", "packages", "crashdumps", "d3dscache", "comms", "downloads", "desktop", "appdata",
		"application data", "local settings", "recovery", "perflogs", "config.msi", "onedrive", "dotnet",
		"common files", "internet explorer", "package cache", "reference assemblies", "modifiablewindowsapps",
		"connecteddevicesplatform", "system volume information", "$recycle.bin", "google", "mozilla",
		"bravesoftware", "opera software", "discord", "amd", "intel", "ati technologies", "code", ".vscode",
		"claude", ".claude", "go", "go-build", "npm", "npm-cache", "pip", "node_modules", "fontcache", "squirreltemp",
		"users", "ccboot", "ccbootclient",
	}
	noisePrefixes = []string{"microsoft", "windows", "nvidia", "gamefix", "raptech", "$windows.", "ntuser"}
)

// isNoise reports whether a folder most likely belongs to Windows or another
// program, not to the app.
func isNoise(folder string) bool {
	if root := os.Getenv("SystemRoot"); root != "" && inside(folder, root) {
		return true
	}
	name := strings.ToLower(filepath.Base(folder))
	if slices.Contains(noiseNames, name) {
		return true
	}
	for _, prefix := range noisePrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}

// inside reports whether path is dir or somewhere below it.
func inside(path, dir string) bool {
	rel, err := filepath.Rel(dir, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, `..\`)
}

func isDir(path string) bool {
	st, err := os.Stat(path)
	return err == nil && st.IsDir()
}

// longPath turns short names, such as ADMINI~1, which Windows reports some
// changes with, into the full names, as far as the path still exists.
func longPath(path string) string {
	if !strings.Contains(path, "~") {
		return path
	}
	rest := ""
	for p := path; ; {
		if short, err := windows.UTF16PtrFromString(p); err == nil {
			buf := make([]uint16, windows.MAX_LONG_PATH)
			if n, err := windows.GetLongPathName(short, &buf[0], uint32(len(buf))); err == nil && n > 0 && int(n) < len(buf) {
				return filepath.Join(windows.UTF16ToString(buf[:n]), rest)
			}
		}
		parent := filepath.Dir(p)
		if parent == p {
			return path
		}
		rest = filepath.Join(filepath.Base(p), rest)
		p = parent
	}
}

// fixedDrives lists the PC's own drives, such as C:\ and D:\.
func fixedDrives() []string {
	mask, err := windows.GetLogicalDrives()
	if err != nil {
		return nil
	}
	var drives []string
	for i := range 26 {
		if mask&(1<<i) == 0 {
			continue
		}
		root := string(rune('A'+i)) + `:\`
		if p, err := windows.UTF16PtrFromString(root); err == nil && windows.GetDriveType(p) == windows.DRIVE_FIXED {
			drives = append(drives, root)
		}
	}
	return drives
}
