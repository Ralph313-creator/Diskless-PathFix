package manager

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"raptech-portable/internal/core"
)

func TestAddedListsEachLinkedFolderOnce(t *testing.T) {
	root := t.TempDir()
	app := filepath.Join(t.TempDir(), "Game")
	os.MkdirAll(app, 0o755)
	os.WriteFile(filepath.Join(app, "game.exe"), []byte("x"), 0o644)
	p, err := CreatePackage(root, "Game", app, filepath.Join(app, "game.exe"), fakeLauncher)
	if err != nil {
		t.Fatal(err)
	}
	saves := filepath.Join(t.TempDir(), "Saves")
	if err := AddFolder(p, saves, nil); err != nil {
		t.Fatal(err)
	}
	back, _ := core.Load(p.Dir)
	got := Added(back)
	// The app's own folder is a linked folder now, so it is listed once.
	if len(got) != 2 || !core.SamePath(got[0], app) || !core.SamePath(got[1], saves) {
		t.Fatalf("Added = %v, want [%s %s]", got, app, saves)
	}
}

func TestAddFoldersSkipsAddedAndReportsProblems(t *testing.T) {
	p, _ := CreatePackage(t.TempDir(), "Game", "", "", fakeLauncher)
	a := filepath.Join(t.TempDir(), "A")
	b := filepath.Join(t.TempDir(), "B")
	os.MkdirAll(a, 0o755)
	os.WriteFile(filepath.Join(a, "x"), []byte("x"), 0o644)
	if err := AddFolder(p, a, nil); err != nil {
		t.Fatal(err)
	}
	n, err := AddFolders(p, []string{a, b, `C:\Windows`})
	if n != 1 {
		t.Fatalf("added %d, want only B", n)
	}
	if err == nil || !strings.Contains(err.Error(), "Windows") {
		t.Fatalf("err = %v, want the C:\\Windows problem", err)
	}
	back, _ := core.Load(p.Dir)
	if len(back.Folders) != 2 {
		t.Fatalf("folders = %+v", back.Folders)
	}
}
