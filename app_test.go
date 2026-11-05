package main

import (
	"os"
	"path/filepath"
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
	// Same form as a line CCBoot runs (run -startup -skipupdate G:\GameFix\GameFix.bat):
	// -winlogon and quotes stopped public.txt from running on the clients.
	want := `run -startup -skipupdate ` + filepath.Join(root, "Roblox", "Roblox.exe") + ` --setup`
	if v.CCBoot != want {
		t.Fatalf("CCBoot = %q, want %q", v.CCBoot, want)
	}
	if len(v.Links) != 2 || !strings.Contains(v.Links[1], "LOCALAPPDATA") || len(v.Session[1]) == 0 {
		t.Fatalf("links %v session %v", v.Links, v.Session)
	}
}

func TestAddFolderAndSession(t *testing.T) {
	root := t.TempDir()
	p, _ := manager.CreatePackage(root, "Game", "", "", []byte("MZ"))
	a := &App{Root: root}
	if err := a.AddFolder(p.Dir, filepath.Join(t.TempDir(), "Saves"), ""); err != nil {
		t.Fatal(err)
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

func TestListShowsIconAndChosenIcon(t *testing.T) {
	root := t.TempDir()
	self, _ := os.Executable()
	launcher, _ := os.ReadFile(self) // a real exe without an icon
	p, err := manager.CreatePackage(root, "Game", "", "", launcher)
	if err != nil {
		t.Fatal(err)
	}
	a := &App{Root: root}
	if v := a.List()[0]; v.Icon != "" || v.OwnIcon {
		t.Fatalf("before: icon %q own %v, want none", v.Icon[:min(20, len(v.Icon))], v.OwnIcon)
	}
	if err := a.changeIcon(p.Dir, filepath.Join(os.Getenv("SystemRoot"), "explorer.exe")); err != nil {
		t.Fatal(err)
	}
	v := a.List()[0]
	if !strings.HasPrefix(v.Icon, "data:image/png;base64,") || !v.OwnIcon {
		t.Fatalf("after: icon %q own %v", v.Icon[:min(30, len(v.Icon))], v.OwnIcon)
	}
	if err := a.UseAppIcon(p.Dir); err != nil {
		t.Fatal(err)
	}
	if v := a.List()[0]; v.OwnIcon || v.Icon != "" {
		t.Fatal("still the chosen icon after Use app's icon (the app has none)")
	}
}

func TestListWithoutRootIsEmpty(t *testing.T) {
	if got := (&App{}).List(); got != nil {
		t.Fatalf("got %v", got)
	}
}
