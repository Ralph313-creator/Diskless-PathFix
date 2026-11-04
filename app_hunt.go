package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"raptech-portable/internal/core"
	"raptech-portable/internal/hunt"
	"raptech-portable/internal/manager"
)

// StartHunt asks for the program that starts the app (the installed one,
// not the package's launcher), starts it, and finds the folders it uses.
// The window gets "hunt" events with a hunt.State. It returns the program,
// or "" when the dialog was cancelled.
func (a *App) StartHunt(dir string) (string, error) {
	start := ""
	if p, err := core.Load(dir); err == nil && p.AppSource != "" {
		start = core.ExpandVars(p.AppSource, p.Dir)
	}
	exe, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "Program that starts the app (the installed one)",
		DefaultDirectory: start,
		Filters:          []runtime.FileFilter{{DisplayName: "Programs (*.exe)", Pattern: "*.exe"}},
	})
	if err != nil || exe == "" {
		return "", err
	}
	_, err = a.startHunt(dir, exe)
	return exe, err
}

// startHunt hunts exe for the package in dir, unless a hunt is running.
func (a *App) startHunt(dir, exe string) (*hunt.Hunter, error) {
	a.huntMu.Lock()
	defer a.huntMu.Unlock()
	if a.hunt != nil && a.hunt.Running() {
		return nil, errors.New("a hunt is already running for " + filepath.Base(a.hunt.Exe()) + "; stop it first")
	}
	if inPackages(exe, a.Root) {
		return nil, errors.New("choose the installed program, not one in the portable apps folder")
	}
	exclude := []string{a.Root}
	if self, err := os.Executable(); err == nil {
		exclude = append(exclude, filepath.Dir(self))
	}
	added := func() []string {
		if p, err := core.Load(dir); err == nil {
			return manager.Added(p)
		}
		return nil
	}
	emit := func(st hunt.State) {
		if a.ctx != nil { // nil in tests
			runtime.EventsEmit(a.ctx, "hunt", st)
		}
	}
	h, err := hunt.New(exe, exclude, added, emit)
	if err != nil {
		return nil, err
	}
	h.Start()
	a.hunt = h
	return h, nil
}

func inPackages(path, root string) bool {
	if root == "" {
		return false
	}
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, `..\`)
}

// StopHunt stops hunting; the app keeps running.
func (a *App) StopHunt() {
	a.huntMu.Lock()
	h := a.hunt
	a.huntMu.Unlock()
	if h != nil {
		h.Stop()
	}
}

// AddHuntFolders adds the folders ticked in the hunt to the package (their
// files are copied in) and returns how many were added.
func (a *App) AddHuntFolders(dir string, paths []string) (int, error) {
	p, err := load(dir)
	if err != nil {
		return 0, err
	}
	return manager.AddFolders(p, paths)
}
