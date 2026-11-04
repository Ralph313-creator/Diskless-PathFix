package core

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/sys/windows/registry"
)

func TestSetProtocolPointsAtLauncher(t *testing.T) {
	name := "rp-test-proto"
	path := `Software\Classes\` + name
	deleteTree(registry.CURRENT_USER, path)
	t.Cleanup(func() { deleteTree(registry.CURRENT_USER, path) })

	if err := SetProtocol(name, `G:\Portable\Roblox\Roblox.exe`); err != nil {
		t.Fatal(err)
	}
	k, err := registry.OpenKey(registry.CURRENT_USER, path, registry.QUERY_VALUE)
	if err != nil {
		t.Fatal(err)
	}
	defer k.Close()
	if _, _, err := k.GetStringValue("URL Protocol"); err != nil {
		t.Error("URL Protocol value missing, browsers won't use it")
	}
	c, err := registry.OpenKey(registry.CURRENT_USER, path+`\shell\open\command`, registry.QUERY_VALUE)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	if v, _, _ := c.GetStringValue(""); v != `"G:\Portable\Roblox\Roblox.exe" "%1"` {
		t.Fatalf("command = %q", v)
	}
}

func TestEmptyFolderKeepsFolder(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "LocalStorage")
	mustWrite(t, filepath.Join(dir, "RobloxCookies.dat"), "login")
	mustWrite(t, filepath.Join(dir, "sub", "x"), "y")
	if err := EmptyFolder(dir); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal("folder removed:", err)
	}
	if len(entries) != 0 {
		t.Fatalf("%d entries left", len(entries))
	}
	if err := EmptyFolder(filepath.Join(t.TempDir(), "missing")); err != nil {
		t.Fatal("missing folder should be fine:", err)
	}
}

func TestRunWaitsForGrandchildren(t *testing.T) {
	// cmd ends at once; the ping it starts lives on for about 2 seconds.
	start := time.Now()
	res, err := Run(`C:\Windows\System32\cmd.exe`, []string{"/c", "start", "/b", "ping", "-n", "3", "127.0.0.1"}, "", 0)
	if err != nil {
		t.Fatal(err)
	}
	if took := time.Since(start); took < 1500*time.Millisecond {
		t.Fatalf("returned after %v, before the grandchild ended", took)
	}
	if res.TimedOut {
		t.Fatal("reported a timeout")
	}
}

func TestRunTimeoutEndsTheTree(t *testing.T) {
	start := time.Now()
	res, err := Run(`C:\Windows\System32\cmd.exe`, []string{"/c", "ping", "-n", "30", "127.0.0.1"}, "", 700*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if !res.TimedOut {
		t.Fatal("timeout not reported")
	}
	if took := time.Since(start); took > 5*time.Second {
		t.Fatalf("took %v; the tree was not ended", took)
	}
}

func TestRunMissingExe(t *testing.T) {
	if _, err := Run(filepath.Join(t.TempDir(), "none.exe"), nil, "", 0); err == nil {
		t.Fatal("missing exe should be an error")
	}
}
