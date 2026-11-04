package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestManifestAcceptsInPlacePathsButNotSystemFolders(t *testing.T) {
	good := Manifest{Name: "Riot", Exe: `G:\Online Games\Riot Games\Riot Client\RiotClientServices.exe`,
		Folders: []Folder{{Link: `C:\Riot Games`, Data: `G:\Online Games\Riot Games`}}}
	if err := good.Validate(); err != nil {
		t.Fatalf("in-place paths refused: %v", err)
	}
	for name, m := range map[string]Manifest{
		"exe in Windows":   {Name: "x", Exe: `C:\Windows\System32\cmd.exe`},
		"data is Windows":  {Name: "x", Folders: []Folder{{Link: `D:\Games\X`, Data: `C:\Windows`}}},
		"data is profile":  {Name: "x", Folders: []Folder{{Link: `D:\Games\X`, Data: `C:\Users\me`}}},
		"data is a drive":  {Name: "x", Folders: []Folder{{Link: `D:\Games\X`, Data: `G:\`}}},
		"relative escapes": {Name: "x", Folders: []Folder{{Link: `D:\Games\X`, Data: `..\Other`}}},
		"link inside data": {Name: "x", Folders: []Folder{{Link: `G:\Games\Riot\Sub`, Data: `G:\Games\Riot`}}},
		"data inside link": {Name: "x", Folders: []Folder{{Link: `G:\Games\Riot`, Data: `G:\Games\Riot\Sub`}}},
	} {
		if err := m.Validate(); err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
}

func TestInPlaceFolderIsLinkedNotOwned(t *testing.T) {
	root := t.TempDir()
	existing := filepath.Join(root, "Online Games", "Riot Games")
	mustWrite(t, filepath.Join(existing, "game.pak"), "game")
	link := filepath.Join(root, "pc", "Riot Games")
	p := &Package{Dir: filepath.Join(root, "Portable", "Riot"), Manifest: Manifest{Name: "Riot",
		Folders: []Folder{{Link: link, Data: existing}}}}

	if got := p.DataPath(p.Folders[0]); !SamePath(got, existing) {
		t.Fatalf("DataPath = %q, want the folder in place %q", got, existing)
	}
	if !p.InPlace(p.Folders[0]) {
		t.Fatal("not reported as in place")
	}
	if err := p.Setup(""); err != nil {
		t.Fatal(err)
	}
	if got := mustRead(t, filepath.Join(link, "game.pak")); got != "game" {
		t.Fatalf("link does not reach the folder in place: %q", got)
	}
	if _, err := os.Stat(filepath.Join(p.Dir, "Data")); !os.IsNotExist(err) {
		t.Fatal("something was made in the package for a folder used in place")
	}
	if err := p.Unlink(""); err != nil {
		t.Fatal(err)
	}
	if got := mustRead(t, filepath.Join(existing, "game.pak")); got != "game" {
		t.Fatal("the folder in place was changed by Unlink")
	}
}

func TestResolveExeInPlace(t *testing.T) {
	dir := t.TempDir()
	exe := filepath.Join(dir, "Riot Client", "RiotClientServices.exe")
	mustWrite(t, exe, "x")
	p := &Package{Dir: filepath.Join(t.TempDir(), "Riot")}
	got, err := p.ResolveExe(exe)
	if err != nil || !SamePath(got, exe) {
		t.Fatalf("ResolveExe = %q, %v", got, err)
	}
	got, err = p.ResolveExe(filepath.Join(dir, "*", "RiotClientServices.exe"))
	if err != nil || !strings.EqualFold(got, exe) {
		t.Fatalf("with *: %q, %v", got, err)
	}
}
