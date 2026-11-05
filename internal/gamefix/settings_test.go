package gamefix

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// newTestApp returns an App using mainFolder, with GameFix.json in a temp folder.
func newTestApp(t *testing.T, mainFolder string) *App {
	os.MkdirAll(mainFolder, 0o755)
	serviceRunning = func() bool { return false } // the app watches, even if the service is installed here
	return &App{settings: Settings{MainFolder: mainFolder}, settingsPath: filepath.Join(t.TempDir(), SettingsName),
		ccbootScript: filepath.Join(t.TempDir(), "public.txt")} // never the real C:\CCBoot\Cmd\public.txt
}

// enable turns folders on, so clients get them, without watching their
// sources: the tests copy them by hand.
func enable(t *testing.T, a *App, game string, names ...string) {
	t.Helper()
	serviceRunning = func() bool { return true } // so the app does not watch
	for _, name := range names {
		if err := a.SetWatch(game, name, true); err != nil {
			t.Fatal(err)
		}
	}
}

func TestSettingsNextToExe(t *testing.T) {
	exe, _ := os.Executable()
	// Its own file, so the GameFix tab never shares settings with the
	// standalone GameFix.exe.
	if got, want := exeSettingsPath(), filepath.Join(filepath.Dir(exe), "RapTech-GameFix.json"); got != want {
		t.Errorf("exeSettingsPath() = %q, want %q", got, want)
	}
}

func TestOwnServiceNotTheStandaloneOne(t *testing.T) {
	if ServiceName == "GameFixWatch" || !strings.Contains(ServiceName, "RapTech") {
		t.Errorf("ServiceName = %q; it must not be the standalone GameFix's GameFixWatch", ServiceName)
	}
	for _, flag := range []string{ServiceFlag, installFlag, removeFlag} {
		if !strings.HasPrefix(flag, "--gamefix-") {
			t.Errorf("flag %q could be taken for one of the launcher's or manager's own", flag)
		}
	}
}

func TestFirstStartWritesNoBat(t *testing.T) {
	main := t.TempDir()
	foreign := "@echo off\r\nrem Made by GameFix. The standalone app's links.\r\n"
	os.WriteFile(filepath.Join(main, BatName), []byte(foreign), 0o644)
	old := defaultMainFolder
	defaultMainFolder = main
	t.Cleanup(func() { defaultMainFolder = old })
	a := &App{settingsPath: filepath.Join(t.TempDir(), SettingsName)} // no settings yet
	serviceRunning = func() bool { return true }
	a.load()
	if b, _ := os.ReadFile(filepath.Join(main, BatName)); string(b) != foreign {
		t.Fatal("the first start rewrote a GameFix.bat it did not make")
	}
	if _, err := os.Stat(filepath.Join(main, BatName+BatBackup)); !os.IsNotExist(err) {
		t.Fatal("the first start touched GameFix.bat at all")
	}
}

func TestWriteBatWithoutMainFolderWritesNothing(t *testing.T) {
	cwd, _ := os.Getwd()
	dir := t.TempDir()
	os.Chdir(dir)
	t.Cleanup(func() { os.Chdir(cwd) })
	a := &App{settingsPath: filepath.Join(t.TempDir(), SettingsName)}
	if err := a.writeBat(); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{BatName, CopyBatName} {
		if _, err := os.Stat(filepath.Join(dir, name)); !os.IsNotExist(err) {
			t.Errorf("%s was written into the current folder", name)
		}
	}
}

func TestWriteBatSetsAsideABatItDidNotMake(t *testing.T) {
	a := newTestApp(t, t.TempDir())
	main := a.settings.MainFolder
	foreign := "@echo off\r\nrem Made by GameFix. The standalone app's links.\r\n"
	foreignCopy := "@echo off\r\nrem Made by GameFix. Copies.\r\n"
	os.WriteFile(filepath.Join(main, BatName), []byte(foreign), 0o644)
	os.WriteFile(filepath.Join(main, CopyBatName), []byte(foreignCopy), 0o644)
	if err := a.writeBat(); err != nil {
		t.Fatal(err)
	}
	for name, want := range map[string]string{BatName: foreign, CopyBatName: foreignCopy} {
		if b, _ := os.ReadFile(filepath.Join(main, name+BatBackup)); string(b) != want {
			t.Errorf("%s was not set aside as %s", name, name+BatBackup)
		}
		if b, _ := os.ReadFile(filepath.Join(main, name)); !strings.Contains(string(b), "RapTech Portable") {
			t.Errorf("%s is not the tab's own now", name)
		}
	}
	// The tab's own bat is simply rewritten; the set-aside one stays as it was.
	if err := a.writeBat(); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(main, BatName+BatBackup)); string(b) != foreign {
		t.Error("the set-aside GameFix.bat was overwritten")
	}
}

func TestLoadUsesTheDefaultMainFolderTheFirstTime(t *testing.T) {
	def := t.TempDir()
	old := defaultMainFolder
	defaultMainFolder = def
	t.Cleanup(func() { defaultMainFolder = old })
	a := &App{settingsPath: filepath.Join(t.TempDir(), SettingsName)}
	serviceRunning = func() bool { return true }
	a.load()
	if a.settings.MainFolder != def {
		t.Fatalf("MainFolder = %q, want the default %q", a.settings.MainFolder, def)
	}
	// A missing default folder is not used.
	defaultMainFolder = filepath.Join(def, "missing")
	b := &App{settingsPath: filepath.Join(t.TempDir(), SettingsName)}
	b.load()
	if b.settings.MainFolder != "" {
		t.Fatalf("MainFolder = %q, want none", b.settings.MainFolder)
	}
}

func TestSettingsHoldEverything(t *testing.T) {
	a := newTestApp(t, t.TempDir())
	a.AddGame("GTA V")
	game := filepath.Join(a.settings.MainFolder, "GTA V")
	src := makeSource(t, "Saves")
	a.AddFolder("GTA V", src)
	if _, err := a.SaveFolders("GTA V", map[string]FolderSettings{"Saves": {Source: portable(src), Link: `D:\Saves`}}); err != nil {
		t.Fatal(err)
	}
	if _, err := a.CopyFolder("GTA V", "Saves"); err != nil {
		t.Fatal(err)
	}

	// Nothing but the copied folders is written into the game folder.
	entries, _ := os.ReadDir(game)
	if len(entries) != 1 {
		t.Errorf("game folder has %d entries, want only Saves", len(entries))
	}

	// A new run reads everything back from GameFix.json.
	b := &App{settingsPath: a.settingsPath}
	b.load()
	if want := (FolderSettings{Source: portable(src), Link: `D:\Saves`}); b.settings.MainFolder != a.settings.MainFolder || b.folders(game)["Saves"] != want {
		t.Errorf("loaded settings = %+v", b.settings)
	}

	// Deleting the game forgets its link paths.
	b.DeleteGame("GTA V")
	data, _ := os.ReadFile(a.settingsPath)
	if strings.Contains(string(data), "GTA V") {
		t.Errorf("GameFix.json still has the deleted game:\n%s", data)
	}
}

func TestLoadWritesLinksWithVariables(t *testing.T) {
	t.Setenv("USERPROFILE", `C:\Users\me`)
	a := newTestApp(t, t.TempDir())
	a.AddGame("GTA V")
	game := filepath.Join(a.settings.MainFolder, "GTA V")
	a.settings.Games = map[string]map[string]FolderSettings{
		game: {"Saves": {Source: `C:\Users\me\Documents\Saves`, Link: `C:\Users\me\Documents\Saves`}},
	}
	a.save()

	b := &App{settingsPath: a.settingsPath}
	b.load()
	want := FolderSettings{Source: `%USERPROFILE%\Documents\Saves`, Link: `%USERPROFILE%\Documents\Saves`}
	if got := b.folders(game)["Saves"]; got != want {
		t.Errorf("loaded Saves = %+v, want %+v", got, want)
	}
}

func TestCopyExpandsSourceVariables(t *testing.T) {
	home := t.TempDir()
	t.Setenv("USERPROFILE", home)
	src := filepath.Join(home, "AppData", "LocalLow", "LandCrab", "PEAK")
	os.MkdirAll(src, 0o755)
	os.WriteFile(filepath.Join(src, "save.dat"), []byte("peak"), 0o644)

	a := newTestApp(t, t.TempDir())
	a.AddGame("PEAK")
	if err := a.AddFolder("PEAK", src); err != nil {
		t.Fatal(err)
	}
	// LocalLow has no variable of its own, so it is written with %USERPROFILE%.
	game := filepath.Join(a.settings.MainFolder, "PEAK")
	if got, want := a.folders(game)[GameItself].Source, `%USERPROFILE%\AppData\LocalLow\LandCrab\PEAK`; got != want {
		t.Errorf("source = %q, want %q", got, want)
	}
	if _, err := a.CopyFolder("PEAK", GameItself); err != nil {
		t.Fatal(err)
	}
	if data, _ := os.ReadFile(filepath.Join(game, "save.dat")); string(data) != "peak" {
		t.Error("Copy did not expand the variable in the source")
	}
}

func TestExpandVars(t *testing.T) {
	t.Setenv("USERPROFILE", `C:\Users\me`)
	t.Setenv("ProgramFiles(x86)", `C:\Program Files (x86)`)
	tests := map[string]string{
		`%USERPROFILE%\Documents`:   `C:\Users\me\Documents`,
		`%ProgramFiles(x86)%\Steam`: `C:\Program Files (x86)\Steam`,
		`%NO_SUCH_VAR%\x`:           `%NO_SUCH_VAR%\x`,
		`D:\Games`:                  `D:\Games`,
	}
	for in, want := range tests {
		if got := expandVars(in); got != want {
			t.Errorf("expandVars(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestImportOldProfiles(t *testing.T) {
	a := newTestApp(t, t.TempDir())
	game := filepath.Join(a.settings.MainFolder, "GTA V")
	os.MkdirAll(filepath.Join(game, "Saves"), 0o755)
	os.MkdirAll(filepath.Join(game, "Config"), 0o755)
	config := makeSource(t, "Config")
	old := filepath.Join(game, "gamefix.json")
	os.WriteFile(old, []byte(`{"links": {"Saves": "D:\\Saves", "Config": "`+strings.ReplaceAll(config, `\`, `\\`)+`"}}`), 0o644)

	// Older versions kept only the link path. It is also the source when
	// it is a folder on this PC.
	a.importOldProfiles()
	if got, want := a.folders(game)["Saves"], (FolderSettings{Link: `D:\Saves`}); got != want {
		t.Errorf("imported Saves = %+v, want %+v", got, want)
	}
	if got, want := a.folders(game)["Config"], (FolderSettings{Source: portable(config), Link: portable(config)}); got != want {
		t.Errorf("imported Config = %+v, want %+v", got, want)
	}
	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Error("the old gamefix.json was not deleted")
	}
	if _, err := os.Stat(filepath.Join(a.settings.MainFolder, BatName)); err != nil {
		t.Errorf("GameFix.bat was not rewritten: %v", err)
	}
	// Imported folders are off until they are turned on.
	enable(t, a, "GTA V", "Saves")
	script, _ := os.ReadFile(filepath.Join(a.settings.MainFolder, BatName))
	if !strings.Contains(string(script), `"%~dp0GTA V\Saves"`) {
		t.Errorf("GameFix.bat is missing Saves:\n%s", script)
	}
}
