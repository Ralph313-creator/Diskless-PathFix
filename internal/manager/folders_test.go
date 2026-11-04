package manager

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"raptech-portable/internal/core"
)

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

func TestAddFolderRefusesAFolderLinkedByAnotherApp(t *testing.T) {
	root := t.TempDir()
	first := newTestPackage(t, root, "First")
	pc := filepath.Join(t.TempDir(), "Riot Games")
	os.MkdirAll(pc, 0o755)
	if err := AddFolder(first, pc, nil); err != nil {
		t.Fatal(err)
	}
	if err := first.Setup(""); err != nil {
		t.Fatal(err)
	}
	second := newTestPackage(t, root, "Second")
	err := AddFolder(second, pc, nil)
	if err == nil || !strings.Contains(err.Error(), "First") {
		t.Fatalf("err = %v, want a refusal naming the First package", err)
	}
	// Copying an app whose folder another app links is refused too.
	os.WriteFile(filepath.Join(first.DataPath(first.Folders[0]), "game.exe"), []byte("x"), 0o644)
	if _, err := CreatePackage(root, "Third", pc, filepath.Join(pc, "game.exe"), fakeLauncher); err == nil {
		t.Fatal("copied an app from a folder another app links")
	}
}
