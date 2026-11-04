package manager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"raptech-portable/internal/core"
)

// CreateInPlace makes root\name for an app that stays where it is, such as
// Riot on the game disk: nothing is copied, the launcher starts exe (with
// args) from its own folder, and removing the package never deletes it.
func CreateInPlace(root, name, exe string, args []string, launcher []byte) (*core.Package, error) {
	exe = filepath.Clean(strings.TrimSpace(exe))
	if info, err := os.Stat(exe); err != nil || info.IsDir() {
		return nil, fmt.Errorf("%s is not a program on this PC", exe)
	}
	if within(exe, root) {
		return nil, fmt.Errorf("%s is inside the portable apps folder; choose the installed program", exe)
	}
	if err := core.CheckPlace(exe); err != nil {
		return nil, err
	}
	p, err := newPackage(root, name, launcher)
	if err != nil {
		return nil, err
	}
	p.Exe = core.Portable(exe)
	p.Args = args
	p.AppSource = core.Portable(filepath.Dir(exe))
	if err := p.Save(); err != nil {
		os.RemoveAll(p.Dir)
		return nil, err
	}
	ApplyIcon(p) // only looks; a launcher without the app's icon still works
	return p, nil
}

// AddFolderInPlace links link (where the app expects a folder) to target, an
// existing folder used where it is: nothing is copied, and removing the
// package never deletes it. session subfolders of target are emptied after
// each player.
func AddFolderInPlace(p *core.Package, link, target string, session []string) error {
	link = core.Portable(filepath.Clean(strings.TrimSpace(link)))
	target = filepath.Clean(strings.TrimSpace(target))
	if err := core.CheckLink(link); err != nil {
		return err
	}
	if info, err := os.Stat(target); err != nil || !info.IsDir() {
		return fmt.Errorf("%s is not a folder on this PC", target)
	}
	real := core.ExpandVars(link, p.Dir)
	for _, have := range Added(p) {
		if core.SamePath(have, real) {
			return fmt.Errorf("%s is already in the package", link)
		}
	}
	if err := linkedElsewhere(real); err != nil {
		return err
	}
	p.Folders = append(p.Folders, core.Folder{Link: link, Data: core.Portable(target), Session: session})
	if err := p.Save(); err != nil {
		p.Folders = p.Folders[:len(p.Folders)-1]
		return err
	}
	return nil
}

// linkedElsewhere refuses a folder that is already a link, such as one
// made by another app's package: copying it would copy that app's files,
// and linking it would take it from that app.
func linkedElsewhere(real string) error {
	if core.LinkStatus(real, "") != core.LinkOther {
		return nil
	}
	target, _ := os.Readlink(real)
	target = strings.TrimPrefix(strings.TrimPrefix(target, `\??\`), `\\?\`)
	owner := ""
	for dir := target; dir != filepath.Dir(dir); dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, core.ManifestName)); err == nil {
			owner = filepath.Base(dir)
			break
		}
	}
	if owner != "" {
		return fmt.Errorf("%s is already linked by the %s app (to %s); remove %s first, or take the folder out of it", real, owner, target, owner)
	}
	return fmt.Errorf("%s is a link to %s; it can't be added", real, target)
}

// FindRiot reads where Riot Client is installed from Riot's own
// %ProgramData%\Riot Games\RiotClientInstalls.json, and returns its Riot
// Games folder, such as G:\Online Games\Riot Games. A Riot Client inside
// skip (the portable apps folder) is a copy and never the answer: running a
// copy makes Riot point rc_default at it, while its games still name the
// real one.
func FindRiot(skip string) (string, bool) {
	raw, err := os.ReadFile(filepath.Join(os.Getenv("ProgramData"), "Riot Games", "RiotClientInstalls.json"))
	if err != nil {
		return "", false
	}
	var installs struct {
		Default    string            `json:"rc_default"`
		Live       string            `json:"rc_live"`
		Patchlines map[string]string `json:"patchlines"`
		Games      map[string]string `json:"associated_client"`
	}
	if json.Unmarshal(raw, &installs) != nil {
		return "", false
	}
	candidates := []string{installs.Default, installs.Live}
	for _, m := range []map[string]string{installs.Patchlines, installs.Games} {
		for _, exe := range m {
			candidates = append(candidates, exe)
		}
	}
	for _, exe := range candidates {
		if exe == "" {
			continue
		}
		exe = filepath.FromSlash(exe)
		if skip != "" && within(exe, skip) {
			continue
		}
		dir := filepath.Dir(filepath.Dir(exe)) // ...\Riot Games\Riot Client\RiotClientServices.exe
		if _, err := os.Stat(filepath.Join(dir, "Riot Client", "RiotClientServices.exe")); err == nil {
			return dir, true
		}
	}
	return "", false
}

// CreateRiot makes root\Riot for the Riot Client in riotDir (such as
// G:\Online Games\Riot Games), used in place. Riot's small folders on C:
// are copied in and linked: %ProgramData%\Riot Games (where the games are)
// and %LOCALAPPDATA%\Riot Games (settings), whose login (Riot Client\Data)
// is never copied and is emptied after each player. Vanguard is not part
// of it: its driver must stay installed in the client image.
func CreateRiot(root, riotDir string, launcher []byte) (*core.Package, error) {
	exe := filepath.Join(riotDir, "Riot Client", "RiotClientServices.exe")
	if _, err := os.Stat(exe); err != nil {
		return nil, fmt.Errorf("no Riot Client in %s (looked for %s)", riotDir, exe)
	}
	p, err := CreateInPlace(root, "Riot", exe, nil, launcher)
	if err != nil {
		return nil, err
	}
	steps := []struct {
		dir     string
		session []string
	}{
		{filepath.Join(os.Getenv("ProgramData"), "Riot Games"), nil},
		{filepath.Join(os.Getenv("LOCALAPPDATA"), "Riot Games"), []string{`Riot Client\Data`}},
	}
	for _, s := range steps {
		if err := AddFolder(p, s.dir, s.session); err != nil {
			os.RemoveAll(p.Dir) // just made: nothing of the user's is in it
			return nil, err
		}
	}
	return p, nil
}
