package manager

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"raptech-portable/internal/core"
)

var fakeLauncher = []byte("MZ fake launcher")

func TestCreateRobloxMachineWide(t *testing.T) {
	root := t.TempDir()
	p, err := CreateRoblox(root, false, fakeLauncher)
	if err != nil {
		t.Fatal(err)
	}
	if p.Dir != filepath.Join(root, "Roblox") {
		t.Fatalf("Dir = %s", p.Dir)
	}
	b, err := os.ReadFile(filepath.Join(p.Dir, "Roblox.exe"))
	if err != nil || string(b) != string(fakeLauncher) {
		t.Fatal("launcher not copied as Roblox.exe")
	}
	back, err := core.Load(p.Dir)
	if err != nil {
		t.Fatal(err)
	}
	if back.Update == nil || back.Update.Kind != "roblox" {
		t.Fatal("no Roblox update")
	}
	if !hasLink(back, `%ProgramFiles%\Roblox\Versions`) || !hasLink(back, `%LOCALAPPDATA%\Roblox`) {
		t.Fatalf("folders = %+v", back.Folders)
	}
	if !strings.HasPrefix(back.Exe, back.Update.Folder+`\`) {
		t.Fatalf("exe %q is not inside the updated folder %q", back.Exe, back.Update.Folder)
	}
	if !sessionHas(back, "LocalStorage") {
		t.Fatal("Roblox logins (LocalStorage) are not wiped after each player")
	}
}

func TestCreateRobloxPerUserKeepsVersionsInsideLocalAppData(t *testing.T) {
	p, err := CreateRoblox(t.TempDir(), true, fakeLauncher)
	if err != nil {
		t.Fatal(err)
	}
	if hasLink(p, `%ProgramFiles%\Roblox\Versions`) {
		t.Fatal("per-user install must not link Program Files")
	}
	if !strings.HasPrefix(p.Update.Folder, `Data\LocalAppData\Roblox\`) {
		t.Fatalf("update folder %q not under the linked LocalAppData\\Roblox", p.Update.Folder)
	}
}

func TestCreateRefusesExistingPackage(t *testing.T) {
	root := t.TempDir()
	if _, err := CreateRoblox(root, false, fakeLauncher); err != nil {
		t.Fatal(err)
	}
	if _, err := CreateRoblox(root, false, fakeLauncher); err == nil {
		t.Fatal("made a second Roblox package over the first")
	}
}

func TestCreatePackageCopiesAppFolder(t *testing.T) {
	root := t.TempDir()
	app := filepath.Join(t.TempDir(), "My Game")
	os.MkdirAll(filepath.Join(app, "bin"), 0o755)
	os.WriteFile(filepath.Join(app, "bin", "game.exe"), []byte("game"), 0o644)
	os.WriteFile(filepath.Join(app, "data.pak"), []byte("pak"), 0o644)

	p, err := CreatePackage(root, "My Game", app, filepath.Join(app, "bin", "game.exe"), fakeLauncher)
	if err != nil {
		t.Fatal(err)
	}
	if p.Exe != `App\bin\game.exe` {
		t.Fatalf("exe = %q", p.Exe)
	}
	if b, _ := os.ReadFile(filepath.Join(p.Dir, `App\data.pak`)); string(b) != "pak" {
		t.Fatal("app folder not copied")
	}
	if _, err := os.Stat(filepath.Join(p.Dir, "My Game.exe")); err != nil {
		t.Fatal("launcher not named after the app")
	}
}

func TestCreatePackageRefusesBadNames(t *testing.T) {
	for _, name := range []string{"", "a\\b", "..", "con", "x:y"} {
		if _, err := CreatePackage(t.TempDir(), name, "", "", fakeLauncher); err == nil {
			t.Errorf("name %q accepted", name)
		}
	}
}

func TestAddFolderImportsExistingAndSaves(t *testing.T) {
	root := t.TempDir()
	pc := t.TempDir()
	t.Setenv("RPTEST_PC", pc)
	os.MkdirAll(filepath.Join(pc, "Saves", "Logins"), 0o755)
	os.WriteFile(filepath.Join(pc, "Saves", "slot1.sav"), []byte("save"), 0o644)
	os.WriteFile(filepath.Join(pc, "Saves", "Logins", "token"), []byte("secret"), 0o644)

	p, err := CreatePackage(root, "Game", "", "", fakeLauncher)
	if err != nil {
		t.Fatal(err)
	}
	if err := AddFolder(p, filepath.Join(pc, "Saves"), []string{"Logins"}); err != nil {
		t.Fatal(err)
	}
	back, _ := core.Load(p.Dir)
	if len(back.Folders) != 1 {
		t.Fatalf("folders = %+v", back.Folders)
	}
	data := back.Path(back.Folders[0].Data)
	if b, _ := os.ReadFile(filepath.Join(data, "slot1.sav")); string(b) != "save" {
		t.Fatal("existing files not brought into the package")
	}
	if _, err := os.Stat(filepath.Join(data, "Logins", "token")); !os.IsNotExist(err) {
		t.Fatal("a session folder's files were copied into the package")
	}
	if b, _ := os.ReadFile(filepath.Join(pc, "Saves", "slot1.sav")); string(b) != "save" {
		t.Fatal("the original was changed")
	}
}

func TestAddFolderRefusesDuplicateAndUnsafe(t *testing.T) {
	p, _ := CreatePackage(t.TempDir(), "Game", "", "", fakeLauncher)
	if err := AddFolder(p, `C:\Windows`, nil); err == nil {
		t.Error("C:\\Windows accepted")
	}
	dir := filepath.Join(t.TempDir(), "X")
	if err := AddFolder(p, dir, nil); err != nil {
		t.Fatal(err)
	}
	if err := AddFolder(p, dir, nil); err == nil {
		t.Error("the same folder added twice")
	}
}

func TestRemoveFolderUnlinksAndKeepsData(t *testing.T) {
	p, _ := CreatePackage(t.TempDir(), "Game", "", "", fakeLauncher)
	link := filepath.Join(t.TempDir(), "Saves")
	os.MkdirAll(link, 0o755)
	os.WriteFile(filepath.Join(link, "a.sav"), []byte("a"), 0o644)
	if err := AddFolder(p, link, nil); err != nil {
		t.Fatal(err)
	}
	data := p.Path(p.Folders[0].Data)
	if err := p.Setup(""); err != nil {
		t.Fatal(err)
	}
	if err := RemoveFolder(p, 0); err != nil {
		t.Fatal(err)
	}
	back, _ := core.Load(p.Dir)
	if len(back.Folders) != 0 {
		t.Fatal("folder still in portable.json")
	}
	if st := core.LinkStatus(link, data); st != core.LinkRealFolder {
		t.Fatalf("link path is %v, want the original folder back", st)
	}
	if b, _ := os.ReadFile(filepath.Join(data, "a.sav")); string(b) != "a" {
		t.Fatal("the package's copy was deleted")
	}
	if err := RemoveFolder(back, 5); err == nil {
		t.Fatal("removing a folder that isn't there should fail")
	}
}

func TestAddRegistrySavesKeyNow(t *testing.T) {
	p, _ := CreatePackage(t.TempDir(), "Game", "", "", fakeLauncher)
	if err := AddRegistry(p, `HKCU\Software\Microsoft\Notepad`); err != nil {
		t.Fatal(err)
	}
	back, _ := core.Load(p.Dir)
	if len(back.Registry) != 1 {
		t.Fatalf("registry = %v", back.Registry)
	}
	if err := AddRegistry(back, `hkcu\software\microsoft\notepad`); err == nil {
		t.Fatal("same key added twice")
	}
	if err := AddRegistry(back, `HKCR\x`); err == nil {
		t.Fatal("bad hive accepted")
	}
}

func TestScanFindsPackages(t *testing.T) {
	root := t.TempDir()
	if _, err := CreateRoblox(root, false, fakeLauncher); err != nil {
		t.Fatal(err)
	}
	os.MkdirAll(filepath.Join(root, "NotAPackage"), 0o755)
	v := filepath.Join(root, `Roblox\Data\ProgramFiles\Roblox\Versions\version-abc`)
	os.MkdirAll(v, 0o755)
	os.WriteFile(filepath.Join(v, "RobloxPlayerBeta.exe"), []byte("x"), 0o644)

	list := Scan(root)
	if len(list) != 1 {
		t.Fatalf("found %d packages", len(list))
	}
	got := list[0]
	if got.Name != "Roblox" || got.Launcher != filepath.Join(root, "Roblox", "Roblox.exe") || got.Version != "version-abc" {
		t.Fatalf("got %+v", got)
	}
}

func hasLink(p *core.Package, link string) bool {
	for _, f := range p.Folders {
		if strings.EqualFold(f.Link, link) {
			return true
		}
	}
	return false
}

func sessionHas(p *core.Package, name string) bool {
	for _, f := range p.Folders {
		for _, s := range f.Session {
			if s == name {
				return true
			}
		}
	}
	return false
}
