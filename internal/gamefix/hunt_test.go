package gamefix

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The hunting itself (programs followed, changed folders found, noise, Stop
// leaving the game running) is tested in internal/hunt. These test what
// GameFix adds: one hunt at a time, and how found folders join a game.

func TestOnlyOneHuntAtATime(t *testing.T) {
	a := newTestApp(t, t.TempDir())
	a.AddGame("Game")
	ping := filepath.Join(os.Getenv("SystemRoot"), "System32", "PING.EXE")
	first, err := a.newHunter("Game", ping)
	if err != nil {
		t.Fatal(err)
	}
	if err := a.runHunt("Game", first); err != nil {
		t.Fatal(err)
	}
	second, _ := a.newHunter("Game", ping)
	if err := a.runHunt("Game", second); err == nil || !strings.Contains(err.Error(), "Game") {
		t.Errorf("a second hunt started while one was running (err %v)", err)
	}
	a.StopHunt()
	<-first.Done()
	// The hunter writes folders with their full names, not ADMINI~1, so
	// compare the folders themselves.
	got, err1 := os.Stat(first.Excluded()[0])
	want, err2 := os.Stat(a.settings.MainFolder)
	if err1 != nil || err2 != nil || !os.SameFile(got, want) {
		t.Errorf("excluded %v, want the main folder %s first", first.Excluded(), a.settings.MainFolder)
	}
	if _, err := a.newHunter("Game", `C:\Games\notes.txt`); err == nil {
		t.Error("hunting a file that is not a program should fail")
	}
}

func TestAddHuntFoldersNamesThemApart(t *testing.T) {
	a := newTestApp(t, t.TempDir())
	a.AddGame("Roblox")
	installed := makeSource(t, "Roblox")
	if err := a.AddFolder("Roblox", installed); err != nil {
		t.Fatal(err)
	}
	data, saves := makeSource(t, "Roblox"), makeSource(t, "Saves")
	n, err := a.AddHuntFolders("Roblox", []string{installed, data, saves})
	if err != nil || n != 2 {
		t.Fatalf("AddHuntFolders = %d, %v; want 2 added", n, err)
	}
	folders := a.folders(filepath.Join(a.settings.MainFolder, "Roblox"))
	if _, ok := folders["Saves"]; !ok {
		t.Errorf("folders = %v, want Saves", folders)
	}
	found := false
	for name, f := range folders {
		if strings.HasPrefix(name, "Roblox (") && expandVars(f.Source) == data {
			found = true
		}
	}
	if !found || folders[GameItself].Source != portable(installed) {
		t.Errorf("folders = %+v, want the game folder and a second Roblox named apart", folders)
	}
}
