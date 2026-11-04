package gamefix

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"raptech-portable/internal/hunt"
)

// Hunt mode finds the folders a game uses. The hunting itself (following
// the game's programs, watching the drives for changes, telling noise
// apart) is internal/hunt, shared with the Portable tab; this file ties it
// to a game in the main folder.

// HuntState is sent to the window as "gamefix-hunt" events while hunting,
// and when it stops: what internal/hunt found, and for which game.
type HuntState struct {
	hunt.State
	Game string `json:"game"`
}

// StartHunt asks for the program that starts a game, starts it, and finds
// the folders it uses until it ends or StopHunt is called. It returns the
// program, or "" if the dialog was cancelled.
func (a *App) StartHunt(game string) (string, error) {
	exe, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   "Choose the program that starts " + game,
		Filters: []runtime.FileFilter{{DisplayName: "Programs (*.exe)", Pattern: "*.exe"}},
	})
	if err != nil || exe == "" {
		return "", err
	}
	h, err := a.newHunter(game, exe)
	if err != nil {
		return "", err
	}
	return exe, a.runHunt(game, h)
}

// StopHunt stops hunting. The game keeps running.
func (a *App) StopHunt() {
	a.huntMu.Lock()
	h := a.hunt
	a.huntMu.Unlock()
	if h != nil {
		h.Stop()
	}
}

// newHunter makes a hunt of a game's program, which runHunt starts. The
// main folder and this program's own folder are never suggested; the
// sources the game already has show as added.
func (a *App) newHunter(game, exe string) (*hunt.Hunter, error) {
	a.mu.Lock()
	gameFolder, err := a.gameFolder(game)
	mainFolder := a.settings.MainFolder
	a.mu.Unlock()
	if err != nil {
		return nil, err
	}
	var exclude []string
	if mainFolder != "" {
		exclude = append(exclude, mainFolder)
	}
	if self, err := os.Executable(); err == nil {
		exclude = append(exclude, filepath.Dir(self))
	}
	added := func() []string {
		a.mu.Lock()
		defer a.mu.Unlock()
		var sources []string
		for _, f := range a.folders(gameFolder) {
			if f.Source != "" {
				sources = append(sources, expandVars(f.Source))
			}
		}
		return sources
	}
	name := filepath.Base(gameFolder)
	emit := func(st hunt.State) {
		if a.ctx != nil { // nil in tests
			runtime.EventsEmit(a.ctx, "gamefix-hunt", HuntState{State: st, Game: name})
		}
	}
	return hunt.New(exe, exclude, added, emit)
}

// runHunt starts h for game, unless another hunt is running.
func (a *App) runHunt(game string, h *hunt.Hunter) error {
	a.huntMu.Lock()
	defer a.huntMu.Unlock()
	if a.hunt != nil && a.hunt.Running() {
		return errors.New("hunt mode is already running for " + a.huntGame + ". Stop it first")
	}
	a.hunt, a.huntGame = h, game
	h.Start()
	return nil
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
