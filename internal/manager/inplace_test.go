package manager

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"raptech-portable/internal/core"
)

// fakeRiot makes a Riot Games folder like the one on the game disk.
func fakeRiot(t *testing.T) (dir, exe string) {
	t.Helper()
	dir = filepath.Join(t.TempDir(), "Online Games", "Riot Games")
	exe = filepath.Join(dir, "Riot Client", "RiotClientServices.exe")
	copyTo(t, iconSource, exe)
	os.MkdirAll(filepath.Join(dir, "VALORANT", "live"), 0o755)
	os.WriteFile(filepath.Join(dir, "VALORANT", "live", "big.pak"), []byte("game data"), 0o644)
	return dir, exe
}

func TestCreateInPlaceCopiesNothing(t *testing.T) {
	root := t.TempDir()
	_, exe := fakeRiot(t)
	p, err := CreateInPlace(root, "Game", exe, []string{"--launch-product=valorant"}, realLauncher(t))
	if err != nil {
		t.Fatal(err)
	}
	back, _ := core.Load(p.Dir)
	if !core.SamePath(core.ExpandVars(back.Exe, back.Dir), exe) || back.Args[0] != "--launch-product=valorant" {
		t.Fatalf("exe %q args %v", back.Exe, back.Args)
	}
	if _, err := os.Stat(filepath.Join(p.Dir, "App")); !os.IsNotExist(err) {
		t.Fatal("the app was copied into the package")
	}
	if got, err := back.ResolveExe(back.Exe); err != nil || !core.SamePath(got, exe) {
		t.Fatalf("the launcher would start %q (%v)", got, err)
	}
	if !sameIcon(t, filepath.Join(p.Dir, "Game.exe"), exe) {
		t.Fatal("the launcher did not get the app's icon")
	}
}

func TestCreateInPlaceRefusesMissingOrPackagedExe(t *testing.T) {
	root := t.TempDir()
	if _, err := CreateInPlace(root, "Game", filepath.Join(t.TempDir(), "none.exe"), nil, fakeLauncher); err == nil {
		t.Error("a missing program was accepted")
	}
	inside := filepath.Join(root, "Other", "x.exe")
	copyTo(t, iconSource, inside)
	if _, err := CreateInPlace(root, "Game", inside, nil, fakeLauncher); err == nil {
		t.Error("a program inside the apps folder was accepted")
	}
}

func TestAddFolderInPlaceLinksWithoutCopying(t *testing.T) {
	root := t.TempDir()
	riot, exe := fakeRiot(t)
	p, err := CreateInPlace(root, "Game", exe, nil, fakeLauncher)
	if err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(t.TempDir(), "C", "Riot Games")
	if err := AddFolderInPlace(p, link, riot, nil); err != nil {
		t.Fatal(err)
	}
	back, _ := core.Load(p.Dir)
	if len(back.Folders) != 1 || !back.InPlace(back.Folders[0]) {
		t.Fatalf("folders = %+v", back.Folders)
	}
	if _, err := os.Stat(filepath.Join(p.Dir, "Data")); !os.IsNotExist(err) {
		t.Fatal("something was copied into the package")
	}
	if err := back.Setup(""); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(link, `VALORANT\live\big.pak`)); string(b) != "game data" {
		t.Fatal("the link does not reach the folder in place")
	}
	for _, bad := range [][2]string{
		{link, filepath.Join(t.TempDir(), "missing")},   // the folder must exist
		{filepath.Join(riot, "Sub"), riot},              // link inside the folder
		{filepath.Join(t.TempDir(), "X"), `C:\Windows`}, // never Windows
	} {
		if err := AddFolderInPlace(back, bad[0], bad[1], nil); err == nil {
			t.Errorf("AddFolderInPlace(%s, %s) accepted", bad[0], bad[1])
		}
	}
}

func TestRemovePackageNeverDeletesInPlaceFolders(t *testing.T) {
	root := t.TempDir()
	riot, exe := fakeRiot(t)
	p, _ := CreateInPlace(root, "Game", exe, nil, fakeLauncher)
	link := filepath.Join(t.TempDir(), "Riot Games")
	if err := AddFolderInPlace(p, link, riot, []string{"VALORANT"}); err != nil {
		t.Fatal(err)
	}
	p.Setup("")
	if err := RemovePackage(root, p); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(riot, `VALORANT\live\big.pak`)); string(b) != "game data" {
		t.Fatal("a folder used in place was deleted")
	}
	if _, err := os.Stat(exe); err != nil {
		t.Fatal("the program used in place was deleted")
	}
	if _, err := os.Lstat(link); !os.IsNotExist(err) {
		t.Fatal("the link was left behind")
	}
}

func TestFindRiotReadsRiotClientInstalls(t *testing.T) {
	riot, _ := fakeRiot(t)
	pd := t.TempDir()
	t.Setenv("ProgramData", pd)
	if _, ok := FindRiot(""); ok {
		t.Fatal("found Riot without RiotClientInstalls.json")
	}
	os.MkdirAll(filepath.Join(pd, "Riot Games"), 0o755)
	// Riot writes forward slashes, like this server's file.
	json := `{"rc_default": "` + strings.ReplaceAll(riot, `\`, "/") + `/Riot Client/RiotClientServices.exe"}`
	os.WriteFile(filepath.Join(pd, `Riot Games\RiotClientInstalls.json`), []byte(json), 0o644)
	got, ok := FindRiot("")
	if !ok || !core.SamePath(got, riot) {
		t.Fatalf("FindRiot = %q, %v; want %s", got, ok, riot)
	}
}

func TestFindRiotSkipsACopyInTheAppsFolder(t *testing.T) {
	riot, _ := fakeRiot(t)
	root := t.TempDir()
	copyTo(t, iconSource, filepath.Join(root, `RiotClientServices\App\RiotClientServices.exe`))
	pd := t.TempDir()
	t.Setenv("ProgramData", pd)
	os.MkdirAll(filepath.Join(pd, "Riot Games"), 0o755)
	slash := func(s string) string { return strings.ReplaceAll(s, `\`, "/") }
	// Like this server after a copy of Riot Client ran from the apps folder:
	// rc_default points at the copy, the games still at the real one.
	json := `{"associated_client": {"` + slash(riot) + `/VALORANT/live/": "` + slash(riot) + `/Riot Client/RiotClientServices.exe"},
		"rc_default": "` + slash(root) + `/RiotClientServices/App/RiotClientServices.exe"}`
	os.WriteFile(filepath.Join(pd, `Riot Games\RiotClientInstalls.json`), []byte(json), 0o644)
	got, ok := FindRiot(root)
	if !ok || !core.SamePath(got, riot) {
		t.Fatalf("FindRiot = %q, %v; want the real %s", got, ok, riot)
	}
}

func TestAddFolderRefusesAFolderLinkedByAnotherApp(t *testing.T) {
	root := t.TempDir()
	first := newTestPackage(t, root, "First")
	pc := filepath.Join(t.TempDir(), "Riot Games")
	os.MkdirAll(pc, 0o755)
	if err := AddFolder(first, pc, nil); err != nil {
		t.Fatal(err)
	}
	if err := first.Setup(""); err != nil {
		t.Fatal(err)
	}
	second := newTestPackage(t, root, "Second")
	err := AddFolder(second, pc, nil)
	if err == nil || !strings.Contains(err.Error(), "First") {
		t.Fatalf("err = %v, want a refusal naming the First package", err)
	}
	if err := AddFolderInPlace(second, pc, t.TempDir(), nil); err == nil {
		t.Fatal("AddFolderInPlace took over another app's link")
	}
}

func TestCreateRiot(t *testing.T) {
	root := t.TempDir()
	riot, exe := fakeRiot(t)
	pc := t.TempDir()
	t.Setenv("ProgramData", filepath.Join(pc, "ProgramData"))
	t.Setenv("LOCALAPPDATA", filepath.Join(pc, "Local"))
	os.MkdirAll(filepath.Join(pc, `ProgramData\Riot Games`), 0o755)
	os.WriteFile(filepath.Join(pc, `ProgramData\Riot Games\RiotClientInstalls.json`), []byte("{}"), 0o644)
	os.MkdirAll(filepath.Join(pc, `Local\Riot Games\Riot Client\Data`), 0o755)
	os.WriteFile(filepath.Join(pc, `Local\Riot Games\Riot Client\Data\RiotGamesPrivateSettings.yaml`), []byte("login"), 0o644)
	os.MkdirAll(filepath.Join(pc, `Local\Riot Games\Riot Client\Config`), 0o755)
	os.WriteFile(filepath.Join(pc, `Local\Riot Games\Riot Client\Config\RiotClientSettings.yaml`), []byte("settings"), 0o644)

	p, err := CreateRiot(root, riot, realLauncher(t))
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "Riot" || !core.SamePath(core.ExpandVars(p.Exe, p.Dir), exe) {
		t.Fatalf("name %q exe %q", p.Name, p.Exe)
	}
	var links []string
	for _, f := range p.Folders {
		links = append(links, f.Link)
		if p.InPlace(f) {
			t.Errorf("%s is used in place; only the small C: folders belong here, copied", f.Link)
		}
	}
	if strings.Join(links, "|") != `%ProgramData%\Riot Games|%LOCALAPPDATA%\Riot Games` {
		t.Fatalf("links = %v", links)
	}
	data := p.DataPath(p.Folders[1])
	if b, _ := os.ReadFile(filepath.Join(data, `Riot Client\Config\RiotClientSettings.yaml`)); string(b) != "settings" {
		t.Fatal("Riot's settings were not brought into the package")
	}
	if _, err := os.Stat(filepath.Join(data, `Riot Client\Data\RiotGamesPrivateSettings.yaml`)); !os.IsNotExist(err) {
		t.Fatal("this PC's Riot login was copied into the package")
	}
	if !sessionHas(p, `Riot Client\Data`) {
		t.Fatal("Riot logins are not emptied after each player")
	}
	if _, err := CreateRiot(t.TempDir(), t.TempDir(), fakeLauncher); err == nil {
		t.Fatal("a folder without Riot Client was accepted")
	}
}
