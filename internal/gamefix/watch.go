package gamefix

import (
	"errors"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"raptech-portable/internal/hunt"
)

// settleTime is how long a source has to stay unchanged before it is copied,
// so a game that writes a save in several steps is copied once, when done.
var settleTime = 2 * time.Second

// retryTime is how long to wait before watching a source again that could
// not be watched, such as one that was deleted or is on a drive that is gone.
var retryTime = 10 * time.Second

// watch is one folder of a game whose source is being watched.
type watch struct {
	game   string // game name
	name   string // folder name
	source string // the source with its variables expanded
	stop   chan struct{}
}

// WatchResult is sent to the UI after a watched folder was copied.
type WatchResult struct {
	Game   string `json:"game"`
	Folder string `json:"folder"`
	Error  string `json:"error"` // "" when it worked
}

// SetWatch turns a folder of a game on or off; the app calls it Enable.
// While it is on, clients get the folder (see links), and the folder is
// copied again whenever its source changes, by the app while it is open or
// by the service. A folder without a source can be on, but is not watched.
func (a *App) SetWatch(game, name string, on bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	folder, err := a.gameFolder(game)
	if err != nil {
		return err
	}
	folders := a.folders(folder)
	f, ok := folders[name]
	if !ok {
		return errors.New(`"` + displayName(game, name) + `" is not in ` + game)
	}
	f.Watch = on
	folders[name] = f
	return a.setFolders(folder, folders)
}

// watchKey names a folder of a game in App.watches.
func watchKey(gameFolder, name string) string {
	return gameFolder + "\x00" + name
}

// syncWatches starts watching each folder that has Watch on, and stops
// watching the rest. A folder whose source changed is watched again from its
// new source. It is called whenever the settings are saved, with a.mu held.
// While the service is running, it watches instead of the app.
func (a *App) syncWatches() {
	want := map[string]*watch{}
	appWatches := a.service || !serviceRunning()
	for gameFolder, folders := range a.settings.Games {
		if !appWatches || !inside(gameFolder, a.settings.MainFolder) {
			continue
		}
		for name, f := range folders {
			if f.Watch && f.Source != "" {
				want[watchKey(gameFolder, name)] = &watch{game: filepath.Base(gameFolder), name: name, source: expandVars(f.Source)}
			}
		}
	}
	if a.watches == nil {
		a.watches = map[string]*watch{}
	}
	for key, w := range a.watches {
		if n, ok := want[key]; !ok || n.source != w.source {
			close(w.stop)
			delete(a.watches, key)
		}
	}
	for key, w := range want {
		if _, ok := a.watches[key]; !ok {
			w.stop = make(chan struct{})
			a.watches[key] = w
			go a.runWatch(w)
		}
	}
}

// stopWatch stops watching a folder of a game, with a.mu held.
func (a *App) stopWatch(gameFolder, name string) {
	if w, ok := a.watches[watchKey(gameFolder, name)]; ok {
		close(w.stop)
		delete(a.watches, watchKey(gameFolder, name))
	}
}

// runWatch copies a watched folder once to catch up on changes made while it
// was not watched, then again each time its source changes, until it is
// stopped.
func (a *App) runWatch(w *watch) {
	changes := make(chan struct{}, 1)
	changed := func() {
		select {
		case changes <- struct{}{}:
		default: // a copy is already due
		}
	}
	changed()
	go func() {
		reported := false
		for {
			err := watchDir(w.source, w.stop, func([]string) { changed() })
			if err == nil {
				return // stopped
			}
			if !reported { // once, not every retry
				a.emitWatch(w, "cannot watch the source: "+err.Error()+". Trying again every 10 seconds.")
				reported = true
			}
			select {
			case <-w.stop:
				return
			case <-time.After(retryTime):
			}
		}
	}()

	for {
		select {
		case <-w.stop:
			return
		case <-changes:
		}
		// Wait until the source stops changing.
		for settled := false; !settled; {
			select {
			case <-w.stop:
				return
			case <-changes:
			case <-time.After(settleTime):
				settled = true
			}
		}
		_, err := a.CopyFolder(w.game, w.name)
		select {
		case <-w.stop:
			return // removed or turned off while copying
		default:
		}
		if err != nil {
			a.emitWatch(w, err.Error())
		} else {
			a.emitWatch(w, "")
		}
	}
}

// emitWatch tells the UI that a watched folder was copied, or failed with
// problem.
func (a *App) emitWatch(w *watch, problem string) {
	if a.log != nil {
		if problem == "" {
			a.log.Printf("copied %s", displayName(w.game, w.name))
		} else {
			a.log.Printf("%s: %s", displayName(w.game, w.name), problem)
		}
	}
	if a.ctx != nil { // nil in tests and in the service
		runtime.EventsEmit(a.ctx, "watch-copied", WatchResult{Game: w.game, Folder: w.name, Error: problem})
	}
}

// watchDir is the drive watcher shared with hunt mode: it calls changed each
// time something in dir, or below it, changes, until stop is closed.
var watchDir = hunt.WatchDir
