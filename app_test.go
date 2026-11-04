package main

import (
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
	want := `run -winlogon "` + filepath.Join(root, "Roblox", "Roblox.exe") + `" --setup`
	if v.CCBoot != want {
		t.Fatalf("CCBoot = %q, want %q", v.CCBoot, want)
	}
	if len(v.Links) != 2 || !strings.Contains(v.Links[1], "LOCALAPPDATA") || len(v.Session[1]) == 0 {
		t.Fatalf("links %v session %v", v.Links, v.Session)
	}
}

func TestListWithoutRootIsEmpty(t *testing.T) {
	if got := (&App{}).List(); got != nil {
		t.Fatalf("got %v", got)
	}
}
