package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func writeManifest(t *testing.T, dir, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, ManifestName), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLoadReadsManifest(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
		"name": "Roblox",
		"exe": "Data\\Versions\\*\\RobloxPlayerBeta.exe",
		"folders": [{"link": "%LOCALAPPDATA%\\Roblox", "data": "Data\\LocalAppData\\Roblox", "session": ["LocalStorage"]}],
		"registry": ["HKCU\\Software\\ROBLOX Corporation"],
		"protocols": ["roblox-player"],
		"update": {"kind": "roblox", "folder": "Data\\Versions", "keep": 2, "every": "30m"}
	}`)
	p, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "Roblox" || len(p.Folders) != 1 || p.Folders[0].Session[0] != "LocalStorage" || p.Update.Keep != 2 {
		t.Fatalf("manifest not read right: %+v", p.Manifest)
	}
	if p.Dir != dir {
		t.Fatalf("Dir = %q, want %q", p.Dir, dir)
	}
}

func TestLoadRefusesUnsafeManifests(t *testing.T) {
	cases := map[string]string{
		"no name":          `{"name": ""}`,
		"escaping data":    `{"name": "x", "folders": [{"link": "D:\\Games\\X", "data": "..\\..\\Windows"}]}`,
		"system link":      `{"name": "x", "folders": [{"link": "C:\\Windows", "data": "Data\\W"}]}`,
		"escaping exe":     `{"name": "x", "exe": "C:\\Windows\\notepad.exe"}`,
		"escaping session": `{"name": "x", "folders": [{"link": "D:\\Games\\X", "data": "Data\\X", "session": ["..\\..\\.."]}]}`,
		"bad registry":     `{"name": "x", "registry": ["HKXX\\Software\\A"]}`,
		"bad protocol":     `{"name": "x", "protocols": ["a\\b"]}`,
		"bad update kind":  `{"name": "x", "update": {"kind": "magic"}}`,
		"bad every":        `{"name": "x", "update": {"kind": "roblox", "folder": "Data\\V", "every": "soon"}}`,
		"broken json":      `{"name": `,
	}
	for name, body := range cases {
		dir := t.TempDir()
		writeManifest(t, dir, body)
		if _, err := Load(dir); err == nil {
			t.Errorf("%s: Load accepted it", name)
		}
	}
}

func TestSaveThenLoadKeepsManifest(t *testing.T) {
	dir := t.TempDir()
	p := &Package{Dir: dir, Manifest: Manifest{Name: "App", Exe: "App.exe", Registry: []string{`HKCU\Software\App`}}}
	if err := p.Save(); err != nil {
		t.Fatal(err)
	}
	back, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if back.Name != "App" || back.Exe != "App.exe" || back.Registry[0] != `HKCU\Software\App` {
		t.Fatalf("got %+v", back.Manifest)
	}
}

func TestResolveExePicksNewestVersionFolder(t *testing.T) {
	dir := t.TempDir()
	mk := func(version string) {
		d := filepath.Join(dir, "Data", "Versions", version)
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(d, "RobloxPlayerBeta.exe"), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		time.Sleep(20 * time.Millisecond) // creation times must differ
	}
	mk("version-old")
	mk("version-new")
	mk(".rp-tmp-version-newest") // a download in progress is never picked
	p := &Package{Dir: dir}
	got, err := p.ResolveExe(`Data\Versions\*\RobloxPlayerBeta.exe`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, `\version-new\`) {
		t.Fatalf("got %q, want the version-new exe", got)
	}
}

func TestResolveExeMissing(t *testing.T) {
	p := &Package{Dir: t.TempDir()}
	if _, err := p.ResolveExe(`Data\*\none.exe`); err == nil {
		t.Error("no match should be an error")
	}
	if _, err := p.ResolveExe(`none.exe`); err == nil {
		t.Error("missing exe should be an error")
	}
}
