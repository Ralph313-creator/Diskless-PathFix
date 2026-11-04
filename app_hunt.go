package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"raptech-portable/internal/hunt"
	"raptech-portable/internal/manager"
)

// Hunt mode is one tool for the whole program (the Hunt tab): it hunts any
// program, and what it finds can then be added to a portable app or to a
// GameFix game.

// StartHunt asks for a program (the installed one), starts it, and finds the
// folders it uses. The window gets "hunt" events with a hunt.State. It
// returns the program, or "" when the dialog was cancelled.
func (a *App) StartHunt() (string, error) {
	exe, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   "Program to hunt (the installed one)",
		Filters: []runtime.FileFilter{{DisplayName: "Programs (*.exe)", Pattern: "*.exe"}},
	})
	if err != nil || exe == "" {
		return "", err
	}
	_, err = a.startHunt(exe)
	return exe, err
}

// startHunt hunts exe, unless a hunt is running. The portable apps folder,
// the GameFix main folder and this program's own folder are never
// suggested.
func (a *App) startHunt(exe string) (*hunt.Hunter, error) {
	a.huntMu.Lock()
	defer a.huntMu.Unlock()
	if a.hunt != nil && a.hunt.Running() {
		return nil, errors.New("a hunt is already running for " + filepath.Base(a.hunt.Exe()) + "; stop it first")
	}
	if within(exe, a.Root) {
		return nil, errors.New("choose the installed program, not one in the portable apps folder")
	}
	var exclude []string
	if a.Root != "" {
		exclude = append(exclude, a.Root)
	}
	if a.gameFixMain != nil {
		if main := a.gameFixMain(); main != "" {
			exclude = append(exclude, main)
		}
	}
	if self, err := os.Executable(); err == nil {
		exclude = append(exclude, filepath.Dir(self))
	}
	emit := func(st hunt.State) {
		if a.ctx != nil { // nil in tests
			runtime.EventsEmit(a.ctx, "hunt", st)
		}
	}
	// What a target already has is marked in the window (see HuntAdded),
	// as the target is chosen after hunting.
	h, err := hunt.New(exe, exclude, nil, emit)
	if err != nil {
		return nil, err
	}
	h.Start()
	a.hunt = h
	return h, nil
}

func within(path, dir string) bool {
	if dir == "" {
		return false
	}
	rel, err := filepath.Rel(dir, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, `..\`)
}

// StopHunt stops hunting; the program keeps running.
func (a *App) StopHunt() {
	a.huntMu.Lock()
	h := a.hunt
	a.huntMu.Unlock()
	if h != nil {
		h.Stop()
	}
}

// AddHuntFolders adds the folders ticked in the hunt to the portable app in
// dir (their files are copied in) and returns how many were added.
func (a *App) AddHuntFolders(dir string, paths []string) (int, error) {
	p, err := load(dir)
	if err != nil {
		return 0, err
	}
	return manager.AddFolders(p, paths)
}

// HuntAdded lists the folders on this PC the portable app in dir already
// covers, to mark them in the hunt's results. Never nil.
func (a *App) HuntAdded(dir string) []string {
	p, err := load(dir)
	if err != nil {
		return []string{}
	}
	return append([]string{}, manager.Added(p)...)
}
