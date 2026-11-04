package main

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"raptech-portable/internal/manager"
)

func TestListShowsPackageWithCCBootLine(t *testing.T) {
	root := t.TempDir()
	if _, err := manager.CreateRoblox(root, false, []byte("MZ")); err != nil {
		t.Fatal(err)
	}
	a := &App{Root: root}
	list := a.List()
	if len(list) != 1 {
		t.Fatalf("listed %d packages", len(list))
	}
	v := list[0]
	want := `run -winlogon "` + filepath.Join(root, "Roblox", "Roblox.exe") + `" --setup`
	if v.CCBoot != want {
		t.Fatalf("CCBoot = %q, want %q", v.CCBoot, want)
	}
	if len(v.Links) != 2 || !strings.Contains(v.Links[1], "LOCALAPPDATA") || len(v.Session[1]) == 0 {
		t.Fatalf("links %v session %v", v.Links, v.Session)
	}
}

func TestOnlyOneHuntAtATime(t *testing.T) {
	root := t.TempDir()
	p, err := manager.CreatePackage(root, "Game", "", "", []byte("MZ"))
	if err != nil {
		t.Fatal(err)
	}
	a := &App{Root: root}
	ping := filepath.Join(os.Getenv("SystemRoot"), "System32", "PING.EXE")
	first, err := a.startHunt(p.Dir, ping)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := a.startHunt(p.Dir, ping); err == nil {
		t.Error("a second hunt started while one was running")
	}
	if !slices.ContainsFunc(first.Excluded(), func(s string) bool { return strings.EqualFold(s, root) }) {
		t.Errorf("excluded %v, want the packages folder %s", first.Excluded(), root)
	}
	a.StopHunt()
	<-first.Done()
	if h, err := a.startHunt(p.Dir, ping); err != nil {
		t.Errorf("no new hunt after the first stopped: %v", err)
	} else {
		a.StopHunt()
		<-h.Done()
	}
}

func TestAddHuntFoldersAndSession(t *testing.T) {
	root := t.TempDir()
	p, _ := manager.CreatePackage(root, "Game", "", "", []byte("MZ"))
	a := &App{Root: root}
	found := filepath.Join(t.TempDir(), "Saves")
	n, err := a.AddHuntFolders(p.Dir, []string{found, found})
	if err != nil || n != 1 {
		t.Fatalf("AddHuntFolders = %d, %v; want 1 added", n, err)
	}
	if err := a.SetSession(p.Dir, 0, "LocalStorage, logs"); err != nil {
		t.Fatal(err)
	}
	v := a.List()[0]
	if strings.Join(v.Session[0], "|") != "LocalStorage|logs" {
		t.Fatalf("session = %v", v.Session[0])
	}
}

func TestRemoveAppDeletesPackageEvenWithBrokenManifest(t *testing.T) {
	root := t.TempDir()
	good, _ := manager.CreatePackage(root, "Good", "", "", []byte("MZ"))
	bad, _ := manager.CreatePackage(root, "Bad", "", "", []byte("MZ"))
	os.WriteFile(filepath.Join(bad.Dir, "portable.json"), []byte("{broken"), 0o644)
	a := &App{Root: root}
	if use := a.InUse(good.Dir); len(use.Programs) != 0 {
		t.Fatalf("InUse = %+v", use)
	}
	for _, dir := range []string{good.Dir, bad.Dir} {
		if err := a.RemoveApp(dir); err != nil {
			t.Fatalf("RemoveApp(%s): %v", dir, err)
		}
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			t.Fatalf("%s still there", dir)
		}
	}
	if len(a.List()) != 0 {
		t.Fatal("removed apps are still listed")
	}
}

func TestListWithoutRootIsEmpty(t *testing.T) {
	if got := (&App{}).List(); got != nil {
		t.Fatalf("got %v", got)
	}
}
