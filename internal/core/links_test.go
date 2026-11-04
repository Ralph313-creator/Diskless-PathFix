package core

import (
	"os"
	"path/filepath"
	"testing"
)

func mustWrite(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func mustRead(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestMakeLinkOnMissingPathLinksIntoPackage(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "pkg", "Data", "Roblox")
	link := filepath.Join(root, "user", "AppData", "Roblox")
	mustWrite(t, filepath.Join(target, "settings.xml"), "pkg")

	if err := MakeLink(link, target); err != nil {
		t.Fatal(err)
	}
	if got := mustRead(t, filepath.Join(link, "settings.xml")); got != "pkg" {
		t.Fatalf("read through link = %q", got)
	}
	// Writing through the link lands in the package (how updates end up there).
	mustWrite(t, filepath.Join(link, "Versions", "version-new", "x.exe"), "new")
	if got := mustRead(t, filepath.Join(target, "Versions", "version-new", "x.exe")); got != "new" {
		t.Fatalf("write through link did not reach the package: %q", got)
	}
	if st := LinkStatus(link, target); st != LinkOK {
		t.Fatalf("status = %v, want LinkOK", st)
	}
}

func TestMakeLinkCreatesMissingTarget(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "pkg", "Data", "New")
	link := filepath.Join(root, "user", "New")
	if err := MakeLink(link, target); err != nil {
		t.Fatal(err)
	}
	if info, err := os.Stat(target); err != nil || !info.IsDir() {
		t.Fatalf("target folder not made: %v", err)
	}
}

func TestMakeLinkMovesRealFolderAside(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "pkg", "Roblox")
	link := filepath.Join(root, "user", "Roblox")
	mustWrite(t, filepath.Join(target, "a.txt"), "pkg")
	mustWrite(t, filepath.Join(link, "old.txt"), "old install")

	if LinkStatus(link, target) != LinkRealFolder {
		t.Fatal("real folder not seen")
	}
	if err := MakeLink(link, target); err != nil {
		t.Fatal(err)
	}
	if got := mustRead(t, filepath.Join(link+BackupSuffix, "old.txt")); got != "old install" {
		t.Fatalf("backup = %q", got)
	}
	if got := mustRead(t, filepath.Join(link, "a.txt")); got != "pkg" {
		t.Fatalf("link = %q", got)
	}
}

func TestMakeLinkRefusesToOverwriteBackup(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "pkg", "Roblox")
	link := filepath.Join(root, "user", "Roblox")
	mustWrite(t, filepath.Join(link, "new.txt"), "second")
	mustWrite(t, filepath.Join(link+BackupSuffix, "first.txt"), "first")

	if err := MakeLink(link, target); err == nil {
		t.Fatal("MakeLink replaced a folder although a backup already exists")
	}
	if got := mustRead(t, filepath.Join(link, "new.txt")); got != "second" {
		t.Fatal("the real folder was touched")
	}
	if got := mustRead(t, filepath.Join(link+BackupSuffix, "first.txt")); got != "first" {
		t.Fatal("the old backup was touched")
	}
}

func TestMakeLinkReplacesLinkToOtherPlace(t *testing.T) {
	root := t.TempDir()
	other := filepath.Join(root, "other")
	target := filepath.Join(root, "pkg")
	link := filepath.Join(root, "user", "App")
	mustWrite(t, filepath.Join(other, "keep.txt"), "other data")
	mustWrite(t, filepath.Join(target, "a.txt"), "pkg")
	if err := MakeLink(link, other); err != nil {
		t.Fatal(err)
	}
	if LinkStatus(link, target) != LinkOther {
		t.Fatal("link to another place not seen")
	}
	if err := MakeLink(link, target); err != nil {
		t.Fatal(err)
	}
	if got := mustRead(t, filepath.Join(link, "a.txt")); got != "pkg" {
		t.Fatalf("link = %q", got)
	}
	if got := mustRead(t, filepath.Join(other, "keep.txt")); got != "other data" {
		t.Fatal("the old link's target was changed")
	}
	if _, err := os.Stat(link + BackupSuffix); !os.IsNotExist(err) {
		t.Fatal("a link must not be backed up")
	}
}

func TestMakeLinkTwiceIsNoChange(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "pkg")
	link := filepath.Join(root, "user", "App")
	if err := MakeLink(link, target); err != nil {
		t.Fatal(err)
	}
	if err := MakeLink(link, target); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(link + BackupSuffix); !os.IsNotExist(err) {
		t.Fatal("second MakeLink made a backup")
	}
}

func TestRemoveLinkRestoresBackupAndKeepsData(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "pkg")
	link := filepath.Join(root, "user", "App")
	mustWrite(t, filepath.Join(target, "a.txt"), "pkg")
	mustWrite(t, filepath.Join(link, "old.txt"), "old")
	if err := MakeLink(link, target); err != nil {
		t.Fatal(err)
	}
	if err := RemoveLink(link, target); err != nil {
		t.Fatal(err)
	}
	if got := mustRead(t, filepath.Join(link, "old.txt")); got != "old" {
		t.Fatalf("backup not restored: %q", got)
	}
	if got := mustRead(t, filepath.Join(target, "a.txt")); got != "pkg" {
		t.Fatal("package data lost")
	}
}

func TestRemoveLinkLeavesRealFolderAlone(t *testing.T) {
	root := t.TempDir()
	link := filepath.Join(root, "App")
	mustWrite(t, filepath.Join(link, "x.txt"), "real")
	if err := RemoveLink(link, filepath.Join(root, "pkg")); err != nil {
		t.Fatal(err)
	}
	if got := mustRead(t, filepath.Join(link, "x.txt")); got != "real" {
		t.Fatal("a real folder was removed")
	}
}

func TestLinkStatusMissing(t *testing.T) {
	root := t.TempDir()
	if st := LinkStatus(filepath.Join(root, "nothing"), root); st != LinkMissing {
		t.Fatalf("status = %v, want LinkMissing", st)
	}
}
