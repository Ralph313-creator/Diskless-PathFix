package main

import (
	"os"
	"path/filepath"
	"testing"

	"raptech-portable/internal/core"
	"raptech-portable/internal/manager"
)

// sameFolder compares folders themselves: the hunter writes full names,
// not short ones such as ADMINI~1.
func sameFolder(a, b string) bool {
	x, err1 := os.Stat(a)
	y, err2 := os.Stat(b)
	return err1 == nil && err2 == nil && os.SameFile(x, y)
}

func TestOneGlobalHuntAtATimeSkippingOurOwnFolders(t *testing.T) {
	root := t.TempDir()
	gameFixMain := t.TempDir()
	a := &App{Root: root, gameFixMain: func() string { return gameFixMain }}
	ping := filepath.Join(os.Getenv("SystemRoot"), "System32", "PING.EXE")

	first, err := a.startHunt(ping)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := a.startHunt(ping); err == nil {
		t.Error("a second hunt started while one was running")
	}
	var apps, gamefix bool
	for _, ex := range first.Excluded() {
		apps = apps || sameFolder(ex, root)
		gamefix = gamefix || sameFolder(ex, gameFixMain)
	}
	if !apps || !gamefix {
		t.Errorf("excluded %v; want the portable apps folder and the GameFix main folder", first.Excluded())
	}
	a.StopHunt()
	<-first.Done()

	inside := filepath.Join(root, "Game", "game.exe")
	os.MkdirAll(filepath.Dir(inside), 0o755)
	os.WriteFile(inside, []byte("MZ"), 0o644)
	if _, err := a.startHunt(inside); err == nil {
		t.Error("hunted a program inside the portable apps folder")
	}
}

func TestHuntAddedAndAddForAPortableApp(t *testing.T) {
	root := t.TempDir()
	p, _ := manager.CreatePackage(root, "Game", "", "", []byte("MZ"))
	saves := filepath.Join(t.TempDir(), "Saves")
	a := &App{Root: root}
	if n, err := a.AddHuntFolders(p.Dir, []string{saves, saves}); err != nil || n != 1 {
		t.Fatalf("AddHuntFolders = %d, %v; want 1", n, err)
	}
	got := a.HuntAdded(p.Dir)
	if len(got) != 1 || !core.SamePath(got[0], saves) {
		t.Fatalf("HuntAdded = %v, want [%s]", got, saves)
	}
	if got := a.HuntAdded(filepath.Join(root, "Missing")); len(got) != 0 || got == nil {
		t.Fatalf("HuntAdded of a missing app = %#v, want an empty list", got)
	}
}
