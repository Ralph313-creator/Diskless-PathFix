package gamefix

import (
	"path/filepath"
	"strings"
	"testing"
)

// Hunting is the Hunt tab's (internal/hunt and the main app). These test
// what GameFix adds: how found folders join a game, and which folders a
// game already has, to mark them in the Hunt tab.

func TestGameSourcesListsWhatAGameHas(t *testing.T) {
	a := newTestApp(t, t.TempDir())
	a.AddGame("Game")
	saves := makeSource(t, "Saves")
	if err := a.AddFolder("Game", saves); err != nil {
		t.Fatal(err)
	}
	got, err := a.GameSources("Game")
	if err != nil || len(got) != 1 || !strings.EqualFold(got[0], saves) {
		t.Fatalf("GameSources = %v, %v; want [%s]", got, err, saves)
	}
	if _, err := a.GameSources("Missing"); err == nil {
		t.Error("a game that does not exist was accepted")
	}
	empty, err := a.GameSources(func() string { a.AddGame("Empty"); return "Empty" }())
	if err != nil || empty == nil || len(empty) != 0 {
		t.Fatalf("GameSources of a game without folders = %#v, %v; want an empty list", empty, err)
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
