package gamefix

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

func TestCleanLink(t *testing.T) {
	good := map[string]string{
		"":                             "",
		`C:\Users\me\Documents\Game\`:  `C:\Users\me\Documents\Game`,
		` C:/Users/me/Saves `:          `C:\Users\me\Saves`,
		`%USERPROFILE%\Documents\Game`: `%USERPROFILE%\Documents\Game`,
		`\\SERVER\share\Game`:          `\\SERVER\share\Game`,
		`D:\Games (x86)\Game\Saves`:    `D:\Games (x86)\Game\Saves`,
	}
	for in, want := range good {
		if got, err := cleanLink(in); err != nil || got != want {
			t.Errorf("cleanLink(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
	for _, ok := range []string{`C:\Users\me\Desktop\Link From Source`, `%USERPROFILE%\.minecraft`, `%APPDATA%\.minecraft`,
		`C:\Users\me\Documents\Rockstar Games`, `C:\Program Files\Steam`, `D:\Games\GTA V`} {
		if _, err := cleanLink(ok); err != nil {
			t.Errorf("cleanLink(%q): %v", ok, err)
		}
	}
	for _, bad := range []string{`C:\`, `C:`, `Documents\Game`, `C:\a"b`, `C:\what?`, `\\SERVER`,
		`C:\Windows`, `c:\users\me`, `C:\Users\me\Documents`, `C:\Users\me\AppData\Roaming\`, `%USERPROFILE%\Desktop`,
		`%appdata%`, `C:\Program Files (x86)`, `%PUBLIC%\Documents`} {
		if _, err := cleanLink(bad); err == nil {
			t.Errorf("cleanLink(%q) should fail", bad)
		}
	}
}

func TestPortable(t *testing.T) {
	t.Setenv("USERPROFILE", `C:\Users\me`)
	t.Setenv("APPDATA", `C:\Users\me\AppData\Roaming`)
	t.Setenv("LOCALAPPDATA", `C:\Users\me\AppData\Local`)
	t.Setenv("PUBLIC", `C:\Users\Public`)
	t.Setenv("ProgramData", `C:\ProgramData`)
	t.Setenv("ProgramFiles", `C:\Program Files`)
	t.Setenv("ProgramFiles(x86)", `C:\Program Files (x86)`)
	tests := map[string]string{
		`C:\Users\me\Documents\Rockstar Games`:   `%USERPROFILE%\Documents\Rockstar Games`,
		`c:\users\ME\Saved Games\Game`:           `%USERPROFILE%\Saved Games\Game`,
		`C:\Users\me\AppData\Roaming\.minecraft`: `%APPDATA%\.minecraft`,
		`C:\Users\me\AppData\Local\Riot Games`:   `%LOCALAPPDATA%\Riot Games`,
		`C:\Users\Public\Documents\Game`:         `%PUBLIC%\Documents\Game`,
		`C:\ProgramData\Game`:                    `%ProgramData%\Game`,
		`C:\Program Files (x86)\Steam`:           `%ProgramFiles(x86)%\Steam`,
		`C:\Program Files\Epic Games`:            `%ProgramFiles%\Epic Games`,
		`C:\Users\meme\Documents\Game`:           `C:\Users\meme\Documents\Game`, // another user
		`D:\Games\GTA V`:                         `D:\Games\GTA V`,
		`%USERPROFILE%\Documents\Game`:           `%USERPROFILE%\Documents\Game`,
	}
	for in, want := range tests {
		if got, err := cleanLink(in); err != nil || got != want {
			t.Errorf("cleanLink(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
	// Written with a variable, a user folder is still refused.
	for _, bad := range []string{`C:\Users\me\Documents`, `C:\Users\me\AppData\Roaming`, `C:\Users\me`} {
		if _, err := cleanLink(bad); err == nil {
			t.Errorf("cleanLink(%q) should fail", bad)
		}
	}
}

// makeSource makes a folder to add to a game, with one file in it.
func makeSource(t *testing.T, name string) string {
	src := filepath.Join(t.TempDir(), name)
	os.MkdirAll(src, 0o755)
	os.WriteFile(filepath.Join(src, "save.dat"), []byte(name), 0o644)
	return src
}

func TestAddCopyRemoveFolders(t *testing.T) {
	a := newTestApp(t, t.TempDir())
	a.AddGame("GTA V")
	game := filepath.Join(a.settings.MainFolder, "GTA V")
	bat := filepath.Join(a.settings.MainFolder, BatName)
	rockstar, noLink := makeSource(t, "Rockstar Games"), makeSource(t, "No Link")

	// Adding only remembers the source; the link path starts as the source.
	for _, src := range []string{rockstar, noLink} {
		if err := a.AddFolder("GTA V", src); err != nil {
			t.Fatal(err)
		}
	}
	if err := a.AddFolder("GTA V", rockstar); err == nil {
		t.Error("adding the same folder twice should fail")
	}
	want := []Folder{{Name: "No Link", Source: portable(noLink), Link: portable(noLink)}, {Name: "Rockstar Games", Source: portable(rockstar), Link: portable(rockstar)}}
	if got := gameFolders(game, a.folders(game)); len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("gameFolders() = %+v, want %+v", got, want)
	}
	if isDir(filepath.Join(game, "Rockstar Games")) {
		t.Error("adding a folder should not copy it")
	}

	got, err := a.SaveFolders("GTA V", map[string]FolderSettings{
		"Rockstar Games": {Source: portable(rockstar), Link: `%USERPROFILE%\Documents\Rockstar Games\`},
		"No Link":        {Source: portable(noLink), Link: ""},
		"Gone":           {Source: portable(rockstar), Link: `C:\Gone`},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got != bat {
		t.Errorf("bat = %q, want %q", got, bat)
	}
	enable(t, a, "GTA V", "Rockstar Games", "No Link")
	line := `call :link "%USERPROFILE%\Documents\Rockstar Games" "%~dp0GTA V\Rockstar Games"`
	if script, _ := os.ReadFile(bat); strings.Contains(string(script), line) {
		t.Error("GameFix.bat should not link a folder before it is copied")
	}

	// Copy copies the source into the main folder, and the folder is linked.
	for _, name := range []string{"Rockstar Games", "No Link"} {
		dest, err := a.CopyFolder("GTA V", name)
		if err != nil {
			t.Fatal(err)
		}
		if data, _ := os.ReadFile(filepath.Join(dest, "save.dat")); string(data) != name {
			t.Errorf("%s was not copied to %s", name, dest)
		}
	}
	script, _ := os.ReadFile(bat)
	if !strings.Contains(string(script), line) {
		t.Errorf("GameFix.bat is missing %s:\n%s", line, script)
	}
	if strings.Contains(string(script), "Gone") || strings.Contains(string(script), "No Link") {
		t.Errorf("GameFix.bat should only link folders with a link path:\n%s", script)
	}
	copyLine := `call :copy "%USERPROFILE%\Documents\Rockstar Games" "%~dp0GTA V\Rockstar Games"`
	if script, _ := os.ReadFile(filepath.Join(a.settings.MainFolder, CopyBatName)); !strings.Contains(string(script), copyLine) {
		t.Errorf("GameFix Copy.bat is missing %s:\n%s", copyLine, script)
	}
	if got := gameFolders(game, a.folders(game)); !got[0].Copied || !got[1].Copied {
		t.Errorf("gameFolders() = %+v, want both copied", got)
	}

	// A folder that is off is left out of both scripts, until it is on again.
	if err := a.SetWatch("GTA V", "Rockstar Games", false); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{BatName, CopyBatName} {
		if script, _ := os.ReadFile(filepath.Join(a.settings.MainFolder, name)); strings.Contains(string(script), "Rockstar") {
			t.Errorf("%s has a folder that is off:\n%s", name, script)
		}
	}
	enable(t, a, "GTA V", "Rockstar Games")
	if script, _ := os.ReadFile(bat); !strings.Contains(string(script), line) {
		t.Errorf("GameFix.bat is missing %s after turning it on again:\n%s", line, script)
	}

	for _, bad := range []map[string]FolderSettings{
		{"Rockstar Games": {Source: portable(rockstar), Link: "Documents"}},
		{"Rockstar Games": {Source: "Documents", Link: ""}},
		{"Rockstar Games": {Source: a.settings.MainFolder, Link: ""}},
	} {
		if _, err := a.SaveFolders("GTA V", bad); err == nil {
			t.Errorf("SaveFolders(%+v) should fail", bad)
		}
	}

	// Copy needs a source that exists.
	a.SaveFolders("GTA V", map[string]FolderSettings{"No Link": {Source: `C:\No Such Folder`}})
	if _, err := a.CopyFolder("GTA V", "No Link"); err == nil {
		t.Error("copying from a missing source should fail")
	}

	// A second game goes into the same GameFix.bat, and the per-game
	// GameFix.bat of older versions is removed.
	a.AddGame("Valorant")
	a.AddFolder("Valorant", makeSource(t, "Config"))
	os.WriteFile(filepath.Join(game, BatName), []byte("@echo off\r\nrem Made by GameFix. old"), 0o644)
	a.SaveFolders("Valorant", map[string]FolderSettings{"Config": {Source: a.folders(filepath.Join(a.settings.MainFolder, "Valorant"))["Config"].Source, Link: `D:\Riot\Config`}})
	enable(t, a, "Valorant", "Config")
	if _, err := a.CopyFolder("Valorant", "Config"); err != nil {
		t.Fatal(err)
	}
	script, _ = os.ReadFile(bat)
	for _, want := range []string{`"%~dp0GTA V\Rockstar Games"`, `call :link "D:\Riot\Config" "%~dp0Valorant\Config"`} {
		if !strings.Contains(string(script), want) {
			t.Errorf("GameFix.bat is missing %s:\n%s", want, script)
		}
	}
	if _, err := os.Stat(filepath.Join(game, BatName)); !os.IsNotExist(err) {
		t.Error("the old per-game GameFix.bat was not removed")
	}

	// Removing a folder deletes its copy and its link.
	if err := a.RemoveFolder("GTA V", "Rockstar Games"); err != nil {
		t.Fatal(err)
	}
	if isDir(filepath.Join(game, "Rockstar Games")) {
		t.Error("the copy of a removed folder was not deleted")
	}
	if !isDir(rockstar) {
		t.Error("removing a folder deleted its source")
	}
	if script, _ = os.ReadFile(bat); strings.Contains(string(script), "Rockstar") {
		t.Errorf("GameFix.bat still links a removed folder:\n%s", script)
	}

	// Deleting a game removes its links.
	if err := a.DeleteGame("Valorant"); err != nil {
		t.Fatal(err)
	}
	if script, _ = os.ReadFile(bat); strings.Contains(string(script), "Valorant") {
		t.Errorf("GameFix.bat still links a deleted game:\n%s", script)
	}
}

func TestFolderNamedLikeGame(t *testing.T) {
	a := newTestApp(t, filepath.Join(t.TempDir(), "main"))
	src := filepath.Join(t.TempDir(), "gta v") // case differs from the game name
	os.MkdirAll(filepath.Join(src, "x64"), 0o755)
	os.WriteFile(filepath.Join(src, "x64", "a.dat"), []byte("a"), 0o644)
	saves := makeSource(t, "Rockstar Games")
	a.AddGame("GTA V")
	game := filepath.Join(a.settings.MainFolder, "GTA V")

	if err := a.AddFolder("GTA V", src); err != nil {
		t.Fatal(err)
	}
	if err := a.AddFolder("GTA V", saves); err != nil {
		t.Fatal(err)
	}
	if _, err := a.CopyFolder("GTA V", "Rockstar Games"); err != nil {
		t.Fatal(err)
	}
	// The game folder is not copied yet, even though Rockstar Games is in it.
	want := []Folder{{Name: GameItself, Source: portable(src), Link: portable(src)}, {Name: "Rockstar Games", Source: portable(saves), Link: portable(saves), Copied: true}}
	if got := gameFolders(game, a.folders(game)); len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("gameFolders() = %+v, want %+v", got, want)
	}

	dest, err := a.CopyFolder("GTA V", GameItself)
	if err != nil {
		t.Fatal(err)
	}
	if dest != game {
		t.Errorf("dest = %q, want %q", dest, game)
	}
	if _, err := os.Stat(filepath.Join(game, "x64", "a.dat")); err != nil {
		t.Error("the game's files were not copied straight into the game folder")
	}
	if isDir(filepath.Join(game, "gta v")) {
		t.Error("the game name was added twice")
	}
	// x64 is part of the game, not a folder of its own.
	want[0].Copied = true
	if got := gameFolders(game, a.folders(game)); len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("gameFolders() = %+v, want %+v", got, want)
	}
	enable(t, a, "GTA V", GameItself)
	script, _ := os.ReadFile(filepath.Join(a.settings.MainFolder, BatName))
	if line := `call :link "` + portable(src) + `" "%~dp0GTA V"`; !strings.Contains(string(script), line) {
		t.Errorf("GameFix.bat is missing %s:\n%s", line, script)
	}

	// Removing the game folder itself deletes the game's files from the
	// main folder, including read-only ones, but not the copies of its
	// other folders.
	os.WriteFile(filepath.Join(game, "readonly.dat"), []byte("x"), 0o444)
	if err := a.RemoveFolder("GTA V", GameItself); err != nil {
		t.Fatal(err)
	}
	for _, gone := range []string{"x64", "readonly.dat"} {
		if _, err := os.Stat(filepath.Join(game, gone)); !os.IsNotExist(err) {
			t.Errorf("removing the game folder itself kept %s", gone)
		}
	}
	if _, err := os.Stat(filepath.Join(game, "Rockstar Games", "save.dat")); err != nil {
		t.Error("removing the game folder itself deleted the copy of another folder")
	}
	if !isDir(src) || !isDir(game) {
		t.Error("removing the game folder itself deleted its source, or the game")
	}
	if got := gameFolders(game, a.folders(game)); len(got) != 1 || got[0].Name != "Rockstar Games" {
		t.Errorf("gameFolders() = %+v, want only Rockstar Games", got)
	}
}

// TestBatsRunWithNobodyThere: CCBoot runs GameFix.bat on every client when
// a customer logs on (run -startup in public.txt), so the scripts must work
// with nobody at the keyboard.
func TestBatsRunWithNobodyThere(t *testing.T) {
	links := []batLink{{Link: `C:\ProgramData\Riot Games`, Source: `Riot\Riot Games`}}
	for name, body := range map[string]string{BatName: makeBat(links), CopyBatName: makeCopyBat(links)} {
		for _, line := range strings.Split(body, "\r\n") {
			if strings.EqualFold(strings.TrimSpace(line), "pause") {
				t.Errorf("%s waits for a key, so its window stays open on the customer's screen", name)
			}
			if strings.Contains(strings.ToLower(line), "net session") && !strings.HasPrefix(line, "rem ") {
				// It fails when the Server service is off, as on many
				// diskless images: then the script asks for administrator
				// rights again and again and never links anything.
				t.Errorf("%s checks for administrator rights with net session: %s", name, line)
			}
		}
		if !strings.Contains(body, "fltmc >nul 2>&1") {
			t.Errorf("%s does not check for administrator rights with fltmc", name)
		}
		if !strings.Contains(body, `-ArgumentList 'elevated'`) || !strings.Contains(body, `if /i "%~1"=="elevated" goto noadmin`) {
			t.Errorf("%s can ask for administrator rights more than once", name)
		}
	}
}

// TestBatMakesLinks runs a real GameFix.bat. mklink needs administrator
// rights, so it is skipped otherwise.
func TestBatMakesLinks(t *testing.T) {
	if exec.Command("net", "session").Run() != nil {
		t.Skip("needs administrator rights")
	}
	root := t.TempDir()
	game := filepath.Join(root, "main", "Game & Co")
	client := filepath.Join(root, "client (x86)")
	for _, name := range []string{"New", "Replace Folder", "Replace Link"} {
		os.MkdirAll(filepath.Join(game, name), 0o755)
		os.WriteFile(filepath.Join(game, name, "save.dat"), []byte(name), 0o644)
	}
	// On the client: a real folder and an old link are in the way.
	os.MkdirAll(filepath.Join(client, "Replace Folder", "sub"), 0o755)
	os.WriteFile(filepath.Join(client, "Replace Folder", "sub", "old.dat"), []byte("old"), 0o444) // read-only
	os.MkdirAll(filepath.Join(root, "elsewhere"), 0o755)
	if err := os.Symlink(filepath.Join(root, "elsewhere"), filepath.Join(client, "Replace Link")); err != nil {
		t.Fatal(err)
	}

	os.MkdirAll(filepath.Join(root, "main", "Other Game", "Saves"), 0o755)
	os.WriteFile(filepath.Join(root, "main", "Other Game", "Saves", "save.dat"), []byte("Other"), 0o644)

	links := []batLink{
		{Link: filepath.Join(client, "Deep", "New"), Source: `Game & Co\New`},
		{Link: filepath.Join(client, "Replace Folder"), Source: `Game & Co\Replace Folder`},
		{Link: filepath.Join(client, "Replace Link"), Source: `Game & Co\Replace Link`},
		{Link: filepath.Join(client, "Missing"), Source: `Game & Co\Missing`},
		{Link: filepath.Join(client, "Game"), Source: `Game & Co`},
		{Link: filepath.Join(client, "Other Saves"), Source: `Other Game\Saves`},
	}
	bat := filepath.Join(root, "main", BatName)
	os.WriteFile(bat, []byte(makeBat(links)), 0o644)
	cmd := exec.Command("cmd")
	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: `cmd /c ""` + bat + `""`} // as Explorer runs it
	out, err := cmd.CombinedOutput()
	t.Logf("%s", out)
	if exit, ok := err.(*exec.ExitError); !ok || exit.ExitCode() != 1 {
		t.Errorf("exit = %v, want 1 failure (the missing folder)", err)
	}

	for i, want := range []string{"New", "Replace Folder", "Replace Link", "", "", "Other"} {
		l := links[i]
		if want == "" {
			continue
		}
		got, err := os.ReadFile(filepath.Join(l.Link, "save.dat"))
		if err != nil || string(got) != want {
			t.Errorf("%s: reading through the link = %q, %v; want %q", l.Link, got, err, want)
		}
		if st, err := os.Lstat(l.Link); err != nil || st.Mode()&os.ModeSymlink == 0 {
			t.Errorf("%s is not a symbolic link", l.Link)
		}
	}
	if _, err := os.Stat(filepath.Join(client, "Replace Folder", "sub")); !os.IsNotExist(err) {
		t.Error("the folder that was in the way was not deleted")
	}
	if _, err := os.Stat(filepath.Join(client, "Replace Folder.bak")); !os.IsNotExist(err) {
		t.Error("a .bak folder was made")
	}
	if _, err := os.Stat(filepath.Join(root, "elsewhere")); err != nil {
		t.Error("removing the old link deleted its target")
	}
	if _, err := os.Lstat(filepath.Join(client, "Missing")); !os.IsNotExist(err) {
		t.Error("a link was made for a missing folder")
	}
	if got, err := os.ReadFile(filepath.Join(client, "Game", "New", "save.dat")); err != nil || string(got) != "New" {
		t.Errorf("the game folder link does not point to the game folder: %q, %v", got, err)
	}
	// What happened stays readable on the client after the window closes.
	log, err := os.ReadFile(filepath.Join(root, "main", "GameFix.log"))
	if err != nil || !strings.Contains(string(log), "[OK]") || !strings.Contains(string(log), "[FAILED]") {
		t.Errorf("GameFix.log = %q, %v; want the [OK] and [FAILED] lines", log, err)
	}
}

// TestCopyBatCopies runs a real GameFix Copy.bat. It asks for administrator
// rights like GameFix.bat, so it is skipped otherwise.
func TestCopyBatCopies(t *testing.T) {
	if exec.Command("net", "session").Run() != nil {
		t.Skip("needs administrator rights")
	}
	root := t.TempDir()
	game := filepath.Join(root, "main", "Game & Co")
	client := filepath.Join(root, "client (x86)")
	for _, name := range []string{"New", "Update", "Replace Link"} {
		os.MkdirAll(filepath.Join(game, name), 0o755)
		os.WriteFile(filepath.Join(game, name, "save.dat"), []byte(name), 0o644)
	}
	// On the client: an old copy with an old file and an extra one, and a
	// link made by GameFix.bat.
	os.MkdirAll(filepath.Join(client, "Update"), 0o755)
	os.WriteFile(filepath.Join(client, "Update", "save.dat"), []byte("old"), 0o644)
	os.WriteFile(filepath.Join(client, "Update", "extra.dat"), []byte("extra"), 0o444) // read-only
	if err := os.Symlink(filepath.Join(game, "Replace Link"), filepath.Join(client, "Replace Link")); err != nil {
		t.Fatal(err)
	}

	links := []batLink{
		{Link: filepath.Join(client, "Deep", "New"), Source: `Game & Co\New`},
		{Link: filepath.Join(client, "Update"), Source: `Game & Co\Update`},
		{Link: filepath.Join(client, "Replace Link"), Source: `Game & Co\Replace Link`},
		{Link: filepath.Join(client, "Missing"), Source: `Game & Co\Missing`},
	}
	bat := filepath.Join(root, "main", CopyBatName)
	os.WriteFile(bat, []byte(makeCopyBat(links)), 0o644)
	cmd := exec.Command("cmd")
	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: `cmd /c ""` + bat + `""`} // as Explorer runs it
	out, err := cmd.CombinedOutput()
	t.Logf("%s", out)
	if exit, ok := err.(*exec.ExitError); !ok || exit.ExitCode() != 1 {
		t.Errorf("exit = %v, want 1 failure (the missing folder)", err)
	}

	for _, l := range links[:3] {
		want := filepath.Base(l.Source)
		if got, err := os.ReadFile(filepath.Join(l.Link, "save.dat")); err != nil || string(got) != want {
			t.Errorf("%s: save.dat = %q, %v; want %q", l.Link, got, err, want)
		}
		if st, err := os.Lstat(l.Link); err != nil || !st.IsDir() || st.Mode()&os.ModeSymlink != 0 {
			t.Errorf("%s is not a folder of its own", l.Link)
		}
	}
	if _, err := os.Stat(filepath.Join(client, "Update", "extra.dat")); !os.IsNotExist(err) {
		t.Error("a file that is not in the main folder was kept")
	}
	if _, err := os.Stat(filepath.Join(game, "Replace Link", "save.dat")); err != nil {
		t.Error("replacing the old link changed the main folder")
	}
	if _, err := os.Lstat(filepath.Join(client, "Missing")); !os.IsNotExist(err) {
		t.Error("a folder was made for a missing folder")
	}
	log, err := os.ReadFile(filepath.Join(root, "main", "GameFix Copy.log"))
	if err != nil || !strings.Contains(string(log), "[OK]") || !strings.Contains(string(log), "[FAILED]") {
		t.Errorf("GameFix Copy.log = %q, %v; want the [OK] and [FAILED] lines", log, err)
	}
}

func TestOldSetupBatIsRemoved(t *testing.T) {
	a := newTestApp(t, t.TempDir())
	old := filepath.Join(a.settings.MainFolder, oldSetupBatName)
	os.WriteFile(old, []byte("@echo off\r\nrem Made by GameFix. Installs the GameFix Sync service on this PC."), 0o644)
	if err := a.writeBat(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Error("the GameFix Client Setup.bat of an earlier version was not removed")
	}
	// One that GameFix did not make is kept.
	os.WriteFile(old, []byte("@echo off\r\necho mine\r\n"), 0o644)
	a.writeBat()
	if _, err := os.Stat(old); err != nil {
		t.Error("a GameFix Client Setup.bat that GameFix did not make was removed")
	}
}
