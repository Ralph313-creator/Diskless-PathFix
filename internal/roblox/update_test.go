package roblox

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

type fakePkg struct {
	name string
	data []byte
	md5  string // "" = real MD5
}

func zipOf(t *testing.T, files map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	for name, body := range files {
		f, err := w.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		f.Write([]byte(body))
	}
	w.Close()
	return buf.Bytes()
}

// fakeRoblox serves the version API, a package manifest and the packages,
// shaped like clientsettings.roblox.com and setup.rbxcdn.com.
func fakeRoblox(t *testing.T, version string, pkgs []fakePkg) (*Updater, *int32) {
	t.Helper()
	var downloads int32
	mux := http.NewServeMux()
	mux.HandleFunc("/v2/client-version/WindowsPlayer", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"version":"0.745.0.7451137","clientVersionUpload":%q,"bootstrapperVersion":""}`, version)
	})
	var manifest strings.Builder
	manifest.WriteString("v0\r\n")
	for _, p := range pkgs {
		sum := p.md5
		if sum == "" {
			s := md5.Sum(p.data)
			sum = hex.EncodeToString(s[:])
		}
		fmt.Fprintf(&manifest, "%s\r\n%s\r\n%d\r\n%d\r\n", p.name, sum, len(p.data), len(p.data)*2)
		data := p.data
		mux.HandleFunc("/"+version+"-"+p.name, func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&downloads, 1)
			w.Write(data)
		})
	}
	mux.HandleFunc("/"+version+"-rbxPkgManifest.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(manifest.String()))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return &Updater{
		Client:     srv.Client(),
		VersionURL: srv.URL + "/v2/client-version/WindowsPlayer",
		SetupURL:   srv.URL,
	}, &downloads
}

func standardPkgs(t *testing.T) []fakePkg {
	return []fakePkg{
		{name: "RobloxApp.zip", data: zipOf(t, map[string]string{"RobloxPlayerBeta.exe": "player", "RobloxPlayerBeta.dll": "dll"})},
		// Roblox's zips use backslashes and list folders as "\Name\".
		{name: "content-configs.zip", data: zipOf(t, map[string]string{`\`: "", `\Patch\`: "", `Patch\Patch.json`: "{}"})},
		{name: "RobloxPlayerInstaller.exe", data: []byte("installer")},
	}
}

func read(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestUpdateInstallsLatestVersion(t *testing.T) {
	u, _ := fakeRoblox(t, "version-abc", standardPkgs(t))
	dir := t.TempDir()
	got, installed, err := u.Update(context.Background(), dir, 2)
	if err != nil {
		t.Fatal(err)
	}
	if got != "version-abc" || !installed {
		t.Fatalf("got %q installed=%v", got, installed)
	}
	v := filepath.Join(dir, "version-abc")
	if read(t, filepath.Join(v, "RobloxPlayerBeta.exe")) != "player" {
		t.Error("RobloxApp.zip not extracted to the version root")
	}
	if read(t, filepath.Join(v, `content\configs\Patch\Patch.json`)) != "{}" {
		t.Error(`content-configs.zip not extracted to content\configs`)
	}
	if read(t, filepath.Join(v, "RobloxPlayerInstaller.exe")) != "installer" {
		t.Error("a plain file package is not copied to the version root")
	}
	if !strings.Contains(read(t, filepath.Join(v, "AppSettings.xml")), "<ContentFolder>content</ContentFolder>") {
		t.Error("AppSettings.xml missing")
	}
	assertOnly(t, dir, "version-abc")
}

func TestUpdateWhenCurrentDownloadsNothing(t *testing.T) {
	u, downloads := fakeRoblox(t, "version-abc", standardPkgs(t))
	dir := t.TempDir()
	if _, _, err := u.Update(context.Background(), dir, 2); err != nil {
		t.Fatal(err)
	}
	atomic.StoreInt32(downloads, 0)
	got, installed, err := u.Update(context.Background(), dir, 2)
	if err != nil || got != "version-abc" || installed {
		t.Fatalf("got %q installed=%v err=%v", got, installed, err)
	}
	if n := atomic.LoadInt32(downloads); n != 0 {
		t.Fatalf("%d downloads although up to date", n)
	}
}

func TestUpdateReplacesBrokenVersionFolder(t *testing.T) {
	u, _ := fakeRoblox(t, "version-abc", standardPkgs(t))
	dir := t.TempDir()
	// A folder of that version without the player (an interrupted install).
	os.MkdirAll(filepath.Join(dir, "version-abc", "content"), 0o755)
	if _, installed, err := u.Update(context.Background(), dir, 2); err != nil || !installed {
		t.Fatalf("installed=%v err=%v", installed, err)
	}
	if read(t, filepath.Join(dir, "version-abc", "RobloxPlayerBeta.exe")) != "player" {
		t.Fatal("broken folder not replaced")
	}
}

func TestUpdateRefusesBadDownloads(t *testing.T) {
	cases := map[string][]fakePkg{
		"wrong md5": {
			{name: "RobloxApp.zip", data: zipOf(t, map[string]string{"RobloxPlayerBeta.exe": "x"}), md5: "00000000000000000000000000000000"},
		},
		"unknown package": {
			{name: "RobloxApp.zip", data: zipOf(t, map[string]string{"RobloxPlayerBeta.exe": "x"})},
			{name: "content-brandnew.zip", data: zipOf(t, map[string]string{"a": "b"})},
		},
		"path escape": {
			{name: "RobloxApp.zip", data: zipOf(t, map[string]string{"RobloxPlayerBeta.exe": "x", `..\..\evil.txt`: "x"})},
		},
		"no player": {
			{name: "content-configs.zip", data: zipOf(t, map[string]string{"a.json": "{}"})},
		},
	}
	for name, pkgs := range cases {
		t.Run(name, func(t *testing.T) {
			u, _ := fakeRoblox(t, "version-bad", pkgs)
			root := t.TempDir()
			dir := filepath.Join(root, "Versions")
			if _, _, err := u.Update(context.Background(), dir, 2); err == nil {
				t.Fatal("update accepted it")
			}
			assertOnly(t, dir)
			if _, err := os.Stat(filepath.Join(root, "evil.txt")); err == nil {
				t.Fatal("wrote outside the version folder")
			}
		})
	}
}

func TestPruneKeepsNewestAndCurrent(t *testing.T) {
	dir := t.TempDir()
	for _, v := range []string{"version-1", "version-2", "version-3", "version-4"} {
		os.MkdirAll(filepath.Join(dir, v), 0o755)
		os.WriteFile(filepath.Join(dir, v, "RobloxPlayerBeta.exe"), []byte("x"), 0o644)
		time.Sleep(20 * time.Millisecond)
	}
	os.MkdirAll(filepath.Join(dir, "SomethingElse"), 0o755)
	// version-1 is the one Roblox says is current (oldest on disk).
	Prune(dir, 2, "version-1", nil)
	assertOnly(t, dir, "SomethingElse", "version-1", "version-4")
}

func assertOnly(t *testing.T, dir string, want ...string) {
	t.Helper()
	entries, _ := os.ReadDir(dir)
	var got []string
	for _, e := range entries {
		got = append(got, e.Name())
	}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("%s holds %v, want %v", dir, got, want)
	}
}
