package core

import (
	"path/filepath"
	"testing"
)

func TestStartPathGoesThroughTheDefaultLocation(t *testing.T) {
	root := t.TempDir()
	pc := filepath.Join(root, "pc", "Games", "My Game")
	p := &Package{Dir: filepath.Join(root, "pkg"), Manifest: Manifest{Name: "Game", Exe: `App\bin\game.exe`,
		Folders: []Folder{{Link: pc, Data: "App"}}}}
	mustWrite(t, filepath.Join(p.Dir, `App\bin\game.exe`), "x")
	exe, err := p.ResolveExe(p.Exe)
	if err != nil {
		t.Fatal(err)
	}

	// Not linked yet: the package's own path.
	if got := p.StartPath(exe); !SamePath(got, exe) {
		t.Fatalf("before linking: %q, want %q", got, exe)
	}
	if err := p.Setup(""); err != nil {
		t.Fatal(err)
	}
	// Linked: through the app's default location, so the app sees that.
	if got := p.StartPath(exe); !SamePath(got, filepath.Join(pc, `bin\game.exe`)) {
		t.Fatalf("after linking: %q, want it under %s", got, pc)
	}
}

func TestStartPathOutsideAnyFolderStaysInThePackage(t *testing.T) {
	root := t.TempDir()
	p := &Package{Dir: filepath.Join(root, "pkg"), Manifest: Manifest{Name: "Game",
		Folders: []Folder{{Link: filepath.Join(root, "pc", "Saves"), Data: `Data\Saves`}}}}
	exe := filepath.Join(p.Dir, `App\game.exe`)
	if got := p.StartPath(exe); !SamePath(got, exe) {
		t.Fatalf("got %q, want %q", got, exe)
	}
}
