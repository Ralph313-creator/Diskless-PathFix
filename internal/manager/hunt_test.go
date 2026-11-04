package manager

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"raptech-portable/internal/core"
)

func TestAddedListsLinksAndAppSource(t *testing.T) {
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
	back, _ := core.Load(p.Dir) // AppSource must survive a reload
	got := Added(back)
	if len(got) != 2 || !core.SamePath(got[0], saves) || !core.SamePath(got[1], app) {
		t.Fatalf("Added = %v, want [%s %s]", got, saves, app)
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

func TestSetSessionSavesSubfolders(t *testing.T) {
	p, _ := CreatePackage(t.TempDir(), "Game", "", "", fakeLauncher)
	if err := AddFolder(p, filepath.Join(t.TempDir(), "Data"), nil); err != nil {
		t.Fatal(err)
	}
	if err := SetSession(p, 0, []string{"LocalStorage", " logs "}); err != nil {
		t.Fatal(err)
	}
	back, _ := core.Load(p.Dir)
	if got := strings.Join(back.Folders[0].Session, "|"); got != "LocalStorage|logs" {
		t.Fatalf("session = %q", got)
	}
	if err := SetSession(back, 0, []string{`..\..\Windows`}); err == nil {
		t.Fatal("a session folder outside the package folder was accepted")
	}
	if err := SetSession(back, 3, nil); err == nil {
		t.Fatal("a missing folder was accepted")
	}
}
