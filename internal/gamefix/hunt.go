package gamefix

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Hunt mode is the Hunt tab, one tool for the whole program. This file is
// what it needs from GameFix: the folders a game already has, and adding
// the folders it found to a game.

// GameSources lists the folders on this PC that a game's folders are copied
// from, to mark them in the Hunt tab. Never nil.
func (a *App) GameSources(game string) ([]string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	gameFolder, err := a.gameFolder(game)
	if err != nil {
		return nil, err
	}
	if !isDir(gameFolder) {
		return nil, fmt.Errorf("there is no game %q", game)
	}
	sources := []string{}
	for _, f := range a.folders(gameFolder) {
		if f.Source != "" {
			sources = append(sources, expandVars(f.Source))
		}
	}
	return sources, nil
}

// AddHuntFolders adds folders found by hunt mode to a game, and returns how
// many it added. Folders the game already has are skipped. Each is named
// like AddFolder names it; if that name is taken, where it is is added,
// such as "Roblox (LOCALAPPDATA)".
func (a *App) AddHuntFolders(game string, paths []string) (int, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	gameFolder, err := a.gameFolder(game)
	if err != nil {
		return 0, err
	}
	folders := a.folders(gameFolder)
	added := 0
	var problems []string
paths:
	for _, path := range paths {
		source, err := checkSource(path, a.settings.MainFolder)
		if err != nil || source == "" {
			problems = append(problems, fmt.Sprintf("%s: %v", path, err))
			continue
		}
		for _, f := range folders {
			if strings.EqualFold(expandVars(f.Source), expandVars(source)) {
				continue paths // already added
			}
		}
		link, _ := cleanLink(source) // "" when the source is a folder that cannot be a link path
		folders[huntFolderName(gameFolder, folders, source)] = FolderSettings{Source: source, Link: link}
		added++
	}
	if added > 0 {
		if err := a.setFolders(gameFolder, folders); err != nil {
			return 0, err
		}
	}
	if len(problems) > 0 {
		return added, errors.New(strings.Join(problems, "\n"))
	}
	return added, nil
}

// huntFolderName names a folder added by hunt mode: like AddFolder, by its
// own name, or GameItself if it is named like the game. When that name is
// taken, or a folder with it is already in the game's files, where it is is
// added, such as "Roblox (LOCALAPPDATA)" or "Saves (D)".
func huntFolderName(gameFolder string, folders map[string]FolderSettings, source string) string {
	base := filepath.Base(expandVars(source))
	free := func(name string) bool {
		for existing := range folders {
			if strings.EqualFold(existing, name) {
				return false
			}
		}
		_, err := os.Stat(copyPath(gameFolder, name))
		return name == GameItself || os.IsNotExist(err)
	}
	if strings.EqualFold(base, filepath.Base(gameFolder)) {
		if free(GameItself) {
			return GameItself
		}
		// Named apart below, as the game folder is also shown with the game's name.
	} else if free(base) {
		return base
	}
	where := strings.TrimSuffix(filepath.VolumeName(expandVars(source)), ":")
	if m := pathVar.FindStringSubmatch(source); m != nil && strings.HasPrefix(source, "%") {
		where = m[1]
	}
	name := base + " (" + where + ")"
	for i := 2; !free(name); i++ {
		name = fmt.Sprintf("%s (%s %d)", base, where, i)
	}
	return name
}
