package manager

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc/mgr"

	"raptech-portable/internal/core"
)

func newTestPackage(t *testing.T, root, name string) *core.Package {
	t.Helper()
	p, err := CreatePackage(root, name, "", "", fakeLauncher)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRemovePackageDeletesItAndKeepsWhatLinksPointTo(t *testing.T) {
	root := t.TempDir()
	p := newTestPackage(t, root, "Game")
	// A junction inside the package to a folder outside must not be followed.
	outside := filepath.Join(t.TempDir(), "Outside")
	os.MkdirAll(outside, 0o755)
	os.WriteFile(filepath.Join(outside, "keep.txt"), []byte("keep"), 0o644)
	if err := core.MakeLink(filepath.Join(p.Dir, `Data\Inner`), outside); err != nil {
		t.Fatal(err)
	}
	// A link on "this PC" into the package is removed, the original put back.
	pc := filepath.Join(t.TempDir(), "Saves")
	os.MkdirAll(pc, 0o755)
	os.WriteFile(filepath.Join(pc, "old.sav"), []byte("old"), 0o644)
	if err := AddFolder(p, pc, nil); err != nil {
		t.Fatal(err)
	}
	if err := p.Setup(""); err != nil {
		t.Fatal(err)
	}

	if err := RemovePackage(root, p); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(p.Dir); !os.IsNotExist(err) {
		t.Fatal("the package folder is still there")
	}
	if b, _ := os.ReadFile(filepath.Join(outside, "keep.txt")); string(b) != "keep" {
		t.Fatal("a folder outside the package was deleted through a link")
	}
	if b, _ := os.ReadFile(filepath.Join(pc, "old.sav")); string(b) != "old" {
		t.Fatal("the original folder was not put back")
	}
}

func TestRemovePackageRemovesProtocolAndSchedule(t *testing.T) {
	root := t.TempDir()
	p := newTestPackage(t, root, "RPTestRemove")
	p.Protocols = []string{"rp-test-remove"}
	p.Save()
	t.Cleanup(func() { RemoveSchedule("RPTestRemove") })
	if err := p.Setup(filepath.Join(p.Dir, LauncherName(p.Name))); err != nil {
		t.Fatal(err)
	}
	if err := SetSchedule("RPTestRemove", `C:\Windows\System32\cmd.exe`, 30*time.Minute); err != nil {
		t.Fatal(err)
	}
	if err := RemovePackage(root, p); err != nil {
		t.Fatal(err)
	}
	if k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Classes\rp-test-remove`, registry.QUERY_VALUE); err == nil {
		k.Close()
		t.Error("the protocol handler is left behind")
	}
	if on, _ := ScheduleOn("RPTestRemove"); on {
		t.Error("the update task is left behind")
	}
}

func TestRemovePackageClosesItsPrograms(t *testing.T) {
	root := t.TempDir()
	p := newTestPackage(t, root, "Game")
	ping := filepath.Join(p.Dir, `App\ping.exe`)
	copyTo(t, filepath.Join(os.Getenv("SystemRoot"), "System32", "PING.EXE"), ping)
	cmd := exec.Command(ping, "-n", "60", "127.0.0.1")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer cmd.Process.Kill()
	time.Sleep(300 * time.Millisecond)

	use := InUse(p)
	if len(use.Programs) != 1 || !strings.EqualFold(use.Programs[0].Path, ping) {
		t.Fatalf("InUse = %+v, want the ping running from the package", use)
	}
	if err := RemovePackage(root, p); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(p.Dir); !os.IsNotExist(err) {
		t.Fatal("the package folder is still there")
	}
}

func TestRemovePackageRemovesItsServices(t *testing.T) {
	root := t.TempDir()
	p := newTestPackage(t, root, "Game")
	exe := filepath.Join(p.Dir, `Service\svc.exe`)
	copyTo(t, filepath.Join(os.Getenv("SystemRoot"), "System32", "PING.EXE"), exe)
	m, err := mgr.Connect()
	if err != nil {
		t.Skip("needs administrator rights:", err)
	}
	defer m.Disconnect()
	const name = "RPTestPortableSvc"
	if old, err := m.OpenService(name); err == nil {
		old.Delete()
		old.Close()
	}
	s, err := m.CreateService(name, exe, mgr.Config{StartType: mgr.StartManual, DisplayName: "RapTech Portable test"}, "-n", "1")
	if err != nil {
		t.Fatal(err)
	}
	s.Close()
	t.Cleanup(func() {
		if s, err := m.OpenService(name); err == nil {
			s.Delete()
			s.Close()
		}
	})

	use := InUse(p)
	if len(use.Services) != 1 || use.Services[0] != name {
		t.Fatalf("InUse services = %v, want %s", use.Services, name)
	}
	if err := RemovePackage(root, p); err != nil {
		t.Fatal(err)
	}
	if s, err := m.OpenService(name); err == nil {
		s.Close()
		t.Fatal("the service is still registered")
	}
}

func TestServiceExe(t *testing.T) {
	cases := map[string]string{
		`G:\Portable\STOVE\LauncherService\LauncherService.exe`: `G:\Portable\STOVE\LauncherService\LauncherService.exe`,
		`"G:\Portable\My App\svc.exe" -k run`:                   `G:\Portable\My App\svc.exe`,
		`C:\Windows\system32\svchost.exe -k netsvcs -p`:         `C:\Windows\system32\svchost.exe`,
		`C:\Program Files\Some App\svc.EXE /service`:            `C:\Program Files\Some App\svc.EXE`,
	}
	for in, want := range cases {
		if got := serviceExe(in); got != want {
			t.Errorf("serviceExe(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestRemovePackageRefusesWhatIsNotAPackage(t *testing.T) {
	root := t.TempDir()
	p := newTestPackage(t, root, "Game")
	other := &core.Package{Dir: filepath.Join(t.TempDir(), "Game"), Manifest: p.Manifest}
	os.MkdirAll(other.Dir, 0o755)
	os.WriteFile(filepath.Join(other.Dir, core.ManifestName), []byte(`{"name":"Game"}`), 0o644)
	if err := RemovePackage(root, other); err == nil {
		t.Error("removed a folder outside the apps folder")
	}
	os.Remove(filepath.Join(p.Dir, core.ManifestName))
	if err := RemovePackage(root, p); err == nil {
		t.Error("removed a folder without portable.json")
	}
	if err := RemovePackage(root, &core.Package{Dir: root, Manifest: p.Manifest}); err == nil {
		t.Error("removed the apps folder itself")
	}
}
