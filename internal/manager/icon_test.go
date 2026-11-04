package manager

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"raptech-portable/internal/icon"
)

var iconSource = filepath.Join(os.Getenv("SystemRoot"), "explorer.exe")

// realLauncher is this test program: a real Go exe without an icon.
func realLauncher(t *testing.T) []byte {
	t.Helper()
	self, _ := os.Executable()
	b, err := os.ReadFile(self)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func copyTo(t *testing.T, from, to string) {
	t.Helper()
	os.MkdirAll(filepath.Dir(to), 0o755)
	in, err := os.Open(from)
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()
	out, err := os.Create(to)
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()
	io.Copy(out, in)
}

func sameIcon(t *testing.T, a, b string) bool {
	t.Helper()
	x, err1 := icon.Read(a)
	y, err2 := icon.Read(b)
	if err1 != nil || err2 != nil || len(x) != len(y) {
		return false
	}
	for i := range x {
		if !bytes.Equal(x[i], y[i]) {
			return false
		}
	}
	return true
}

func TestCreatePackageGivesLauncherTheAppIcon(t *testing.T) {
	app := filepath.Join(t.TempDir(), "Game")
	copyTo(t, iconSource, filepath.Join(app, "game.exe"))
	p, err := CreatePackage(t.TempDir(), "Game", app, filepath.Join(app, "game.exe"), realLauncher(t))
	if err != nil {
		t.Fatal(err)
	}
	if !sameIcon(t, filepath.Join(p.Dir, "Game.exe"), iconSource) {
		t.Fatal("the launcher does not have the app's icon")
	}
}

func TestFixIconsOnlyFillsMissingOnes(t *testing.T) {
	root := t.TempDir()
	p, err := CreateRoblox(root, false, realLauncher(t))
	if err != nil {
		t.Fatal(err)
	}
	copyTo(t, iconSource, p.Path(`Data\ProgramFiles\Roblox\Versions\version-abc\RobloxPlayerBeta.exe`))
	if n := FixIcons(root); n != 1 {
		t.Fatalf("fixed %d, want 1", n)
	}
	if !sameIcon(t, filepath.Join(p.Dir, "Roblox.exe"), iconSource) {
		t.Fatal("the launcher did not get Roblox's icon")
	}
	if n := FixIcons(root); n != 0 {
		t.Fatalf("fixed %d the second time, want 0 (it already has one)", n)
	}
}

func TestApplyIconAfterRobloxDownload(t *testing.T) {
	p, err := CreateRoblox(t.TempDir(), false, realLauncher(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := ApplyIcon(p); err != nil {
		t.Fatalf("no Roblox yet should not be an error: %v", err)
	}
	copyTo(t, iconSource, p.Path(`Data\ProgramFiles\Roblox\Versions\version-abc\RobloxPlayerBeta.exe`))
	if err := ApplyIcon(p); err != nil {
		t.Fatal(err)
	}
	if !sameIcon(t, filepath.Join(p.Dir, "Roblox.exe"), iconSource) {
		t.Fatal("the launcher does not have Roblox's icon")
	}
}
