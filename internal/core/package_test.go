package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/sys/windows/registry"
)

// testPackage makes a package whose folder links live under a temp "PC"
// folder, reached through the RPTEST_PC variable like %LOCALAPPDATA%.
func testPackage(t *testing.T) (p *Package, pc string, regKey string) {
	t.Helper()
	root := t.TempDir()
	pc = filepath.Join(root, "pc")
	t.Setenv("RPTEST_PC", pc)
	regKey, regPath := testKey(t)
	k, _, err := registry.CreateKey(registry.CURRENT_USER, regPath, registry.ALL_ACCESS)
	if err != nil {
		t.Fatal(err)
	}
	k.Close()
	p = &Package{Dir: filepath.Join(root, "pkg"), Manifest: Manifest{
		Name: "Test",
		Folders: []Folder{
			{Link: `%RPTEST_PC%\AppData\App`, Data: `Data\App`, Session: []string{"LocalStorage"}},
		},
		Registry: []string{regKey},
	}}
	return p, pc, regKey
}

func TestSetupLinksAndLoadsRegistry(t *testing.T) {
	p, pc, regKey := testPackage(t)
	mustWrite(t, filepath.Join(p.Dir, `Data\App\settings.xml`), "saved")
	// A saved registry file in the package, using %PKG%.
	_, regPath, _ := splitKey(regKey)
	k, _, _ := registry.CreateKey(registry.CURRENT_USER, regPath, registry.ALL_ACCESS)
	k.SetStringValue("Path", p.Dir+`\Data`)
	k.Close()
	if err := SaveRegistry(p.Registry, filepath.Join(p.Dir, RegistryDir), p.Dir); err != nil {
		t.Fatal(err)
	}
	deleteTree(registry.CURRENT_USER, regPath)

	if err := p.Setup(""); err != nil {
		t.Fatal(err)
	}
	if got := mustRead(t, filepath.Join(pc, `AppData\App\settings.xml`)); got != "saved" {
		t.Fatalf("link not made: %q", got)
	}
	k, err := registry.OpenKey(registry.CURRENT_USER, regPath, registry.QUERY_VALUE)
	if err != nil {
		t.Fatal("registry not loaded:", err)
	}
	defer k.Close()
	if v, _, _ := k.GetStringValue("Path"); v != p.Dir+`\Data` {
		t.Fatalf("Path = %q", v)
	}
}

func TestFinishSavesRegistryAndEmptiesSession(t *testing.T) {
	p, pc, regKey := testPackage(t)
	if err := p.Setup(""); err != nil {
		t.Fatal(err)
	}
	// The app runs: writes a login and a setting, changes the registry.
	mustWrite(t, filepath.Join(pc, `AppData\App\LocalStorage\cookies.dat`), "customer login")
	mustWrite(t, filepath.Join(pc, `AppData\App\settings.xml`), "graphics=10")
	_, regPath, _ := splitKey(regKey)
	k, _, _ := registry.CreateKey(registry.CURRENT_USER, regPath, registry.ALL_ACCESS)
	k.SetStringValue("Volume", "5")
	k.Close()

	if err := p.Finish(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(p.Dir, `Data\App\LocalStorage\cookies.dat`)); !os.IsNotExist(err) {
		t.Fatal("login kept for the next customer")
	}
	if got := mustRead(t, filepath.Join(p.Dir, `Data\App\settings.xml`)); got != "graphics=10" {
		t.Fatal("setting not kept in the package")
	}
	saved := decodeUTF16(t, filepath.Join(p.Dir, RegistryDir, RegFileName(regKey)))
	if !strings.Contains(saved, `"Volume"="5"`) {
		t.Fatalf("registry not saved:\n%s", saved)
	}
}

func TestUnlinkPutsOriginalBack(t *testing.T) {
	p, pc, _ := testPackage(t)
	mustWrite(t, filepath.Join(pc, `AppData\App\old.txt`), "original install")
	if err := p.Setup(""); err != nil {
		t.Fatal(err)
	}
	if err := p.Unlink(""); err != nil {
		t.Fatal(err)
	}
	if got := mustRead(t, filepath.Join(pc, `AppData\App\old.txt`)); got != "original install" {
		t.Fatalf("original not back: %q", got)
	}
}

func TestUnlinkRemovesOnlyOurProtocols(t *testing.T) {
	p, _, _ := testPackage(t)
	ours, theirs := "rp-test-ours", "rp-test-theirs"
	for _, name := range []string{ours, theirs} {
		path := `Software\Classes\` + name
		t.Cleanup(func() { deleteTree(registry.CURRENT_USER, path) })
	}
	p.Protocols = []string{ours, theirs}
	launcher := filepath.Join(p.Dir, "App.exe")
	if err := p.Setup(launcher); err != nil {
		t.Fatal(err)
	}
	// Something else (such as the app's own installer) took over "theirs".
	if err := SetProtocol(theirs, `C:\Other\app.exe`); err != nil {
		t.Fatal(err)
	}
	if err := p.Unlink(launcher); err != nil {
		t.Fatal(err)
	}
	if k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Classes\`+ours, registry.QUERY_VALUE); err == nil {
		k.Close()
		t.Error("our protocol handler was left behind")
	}
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Classes\`+theirs+`\shell\open\command`, registry.QUERY_VALUE)
	if err != nil {
		t.Fatal("another program's protocol handler was removed")
	}
	k.Close()
}

func TestStatusReportsEachFolder(t *testing.T) {
	p, _, _ := testPackage(t)
	if st := p.Status(); len(st) != 1 || st[0].State != LinkMissing {
		t.Fatalf("before setup: %+v", st)
	}
	if err := p.Setup(""); err != nil {
		t.Fatal(err)
	}
	if st := p.Status(); st[0].State != LinkOK {
		t.Fatalf("after setup: %+v", st)
	}
}
