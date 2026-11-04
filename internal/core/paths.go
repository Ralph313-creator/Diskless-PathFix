// Package core makes a package folder work on this PC: links, registry,
// protocol handlers, and starting the app.
package core

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// pathVars are the variables link paths are written with, so a package works
// for whoever is logged on. Longer values are more specific, so they win:
// %LOCALAPPDATA% before %USERPROFILE%. (From GameFix linkbat.go.)
var pathVars = []string{"LOCALAPPDATA", "APPDATA", "USERPROFILE", "PUBLIC", "ProgramData", "ProgramFiles(x86)", "ProgramFiles"}

// Portable writes the start of a path as a variable when it is one of the
// pathVars folders: C:\Users\me\AppData\Local\Roblox -> %LOCALAPPDATA%\Roblox.
func Portable(path string) string {
	best, bestValue := "", ""
	for _, name := range pathVars {
		value := strings.TrimRight(os.Getenv(name), `\`)
		if value == "" || len(value) <= len(bestValue) {
			continue
		}
		if strings.EqualFold(path, value) || hasPrefixFold(path, value+`\`) {
			best, bestValue = name, value
		}
	}
	if best == "" {
		return path
	}
	return "%" + best + "%" + path[len(bestValue):]
}

var pathVar = regexp.MustCompile(`%([A-Za-z_][A-Za-z0-9_()]*)%`)

// ExpandVars replaces %NAME% with the variable's value, the way cmd does.
// %PKG% is the package folder. Unknown variables are left as they are.
func ExpandVars(path, pkg string) string {
	return pathVar.ReplaceAllStringFunc(path, func(m string) string {
		name := m[1 : len(m)-1]
		if strings.EqualFold(name, "PKG") && pkg != "" {
			return pkg
		}
		if value, ok := os.LookupEnv(name); ok {
			return value
		}
		return m
	})
}

func hasPrefixFold(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[:len(prefix)], prefix)
}

// SamePath reports whether two Windows paths name the same place, ignoring
// case, a trailing backslash and the \\?\ and \??\ prefixes.
func SamePath(a, b string) bool {
	return strings.EqualFold(normPath(a), normPath(b))
}

func normPath(p string) string {
	p = strings.ReplaceAll(p, "/", `\`)
	for _, prefix := range []string{`\\?\UNC\`, `\??\UNC\`} {
		if hasPrefixFold(p, prefix) {
			return strings.TrimRight(filepath.Clean(`\\`+p[len(prefix):]), `\`)
		}
	}
	for _, prefix := range []string{`\\?\`, `\??\`} {
		if hasPrefixFold(p, prefix) {
			p = p[len(prefix):]
		}
	}
	return strings.TrimRight(filepath.Clean(p), `\`)
}

// IsUNC reports whether p is a network path such as \\SERVER\share\x.
func IsUNC(p string) bool {
	return strings.HasPrefix(normPath(p), `\\`)
}

var absPath = regexp.MustCompile(`^([A-Za-z]:\\.|\\\\[^\\]+\\[^\\]+)`)

// CheckLink checks a link path (with variables) before anything is put
// there. It must be a full path, and never a Windows or user folder itself,
// because a real folder at the link path is moved aside.
func CheckLink(link string) error {
	link = strings.TrimRight(strings.ReplaceAll(strings.TrimSpace(link), "/", `\`), `\`)
	if link == "" {
		return errors.New("the link path is empty")
	}
	if strings.ContainsAny(link, `"<>|?*`) {
		return errors.New(`a link path cannot contain " < > | ? *`)
	}
	if protected(Portable(link)) || protected(Portable(ExpandVars(link, ""))) {
		return errors.New(link + " is a Windows or user folder; use the app's own folder inside it")
	}
	if !absPath.MatchString(ExpandVars(link, "")) {
		return errors.New(link + ` is not a full path, such as %LOCALAPPDATA%\Roblox`)
	}
	return nil
}

// CheckData checks a path inside the package: relative, and staying inside.
func CheckData(rel string) error {
	rel = strings.ReplaceAll(strings.TrimSpace(rel), "/", `\`)
	if rel == "" {
		return errors.New("the package path is empty")
	}
	if filepath.IsAbs(rel) || filepath.VolumeName(rel) != "" || strings.HasPrefix(rel, `\`) {
		return errors.New(rel + ` must be inside the package, such as Data\LocalAppData\Roblox`)
	}
	for _, part := range strings.Split(rel, `\`) {
		if part == ".." {
			return errors.New(rel + " must stay inside the package")
		}
	}
	return nil
}

// IsFull reports whether a path in portable.json is a full path (a folder
// or program used in place, outside the package), not one inside it.
func IsFull(path string) bool {
	return absPath.MatchString(ExpandVars(strings.ReplaceAll(strings.TrimSpace(path), "/", `\`), ""))
}

// CheckPlace checks a folder or program the package uses: inside the
// package (see CheckData), or a full path used in place, such as
// G:\Online Games\Riot Games, which must not be in Windows or be a Windows,
// user or drive-root folder.
func CheckPlace(path string) error {
	if !IsFull(path) {
		return CheckData(path)
	}
	full := strings.TrimRight(filepath.Clean(ExpandVars(path, "")), `\`)
	if root := os.Getenv("SystemRoot"); root != "" && (SamePath(full, root) || hasPrefixFold(full, strings.TrimRight(root, `\`)+`\`)) {
		return errors.New(path + " is in Windows")
	}
	if protected(Portable(full)) || protected(full) {
		return errors.New(path + " is a Windows or user folder, or a whole drive")
	}
	return nil
}

// nested reports whether a and b are the same folder or one is inside the other.
func nested(a, b string) bool {
	a, b = normPath(a), normPath(b)
	return strings.EqualFold(a, b) || hasPrefixFold(a, b+`\`) || hasPrefixFold(b, a+`\`)
}

// protected reports whether link is a Windows or user folder (from GameFix).
func protected(link string) bool {
	link = strings.ToLower(link)
	var patterns []string
	for _, home := range []string{`?:\users\*`, `%userprofile%`} {
		patterns = append(patterns, home)
		for _, sub := range userFolders {
			patterns = append(patterns, home+`\`+sub)
		}
	}
	for _, sub := range userFolders {
		patterns = append(patterns, `%public%\`+sub)
	}
	for _, pattern := range append(patterns, systemFolders...) {
		if ok, _ := filepath.Match(pattern, link); ok {
			return true
		}
	}
	return false
}

var userFolders = []string{
	"desktop", "documents", "downloads", "music", "pictures", "videos", "saved games", "favorites",
	"links", "contacts", "searches", "3d objects", "onedrive", "appdata", `appdata\local`,
	`appdata\locallow`, `appdata\roaming`,
}

var systemFolders = []string{
	`?:`, `?:\windows`, `?:\windows\*`, `?:\program files`, `?:\program files (x86)`, `?:\programdata`, `?:\users`,
	`%appdata%`, `%localappdata%`, `%programdata%`, `%public%`, `%homepath%`, `%windir%`,
	`%systemroot%`, `%systemdrive%`, `%programfiles%`, `%programfiles(x86)%`, `%onedrive%`,
}
