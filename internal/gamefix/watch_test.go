package gamefix

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"
)

// waitFor reports whether ok becomes true within a few seconds.
func waitFor(ok func() bool) bool {
	for end := time.Now().Add(5 * time.Second); time.Now().Before(end); time.Sleep(50 * time.Millisecond) {
		if ok() {
			return true
		}
	}
	return ok()
}

func TestWatch(t *testing.T) {
	settleTime = 100 * time.Millisecond
	a := newTestApp(t, t.TempDir())
	a.AddGame("GTA V")
	src := makeSource(t, "Rockstar Games")
	a.AddFolder("GTA V", src)
	dest := filepath.Join(a.settings.MainFolder, "GTA V", "Rockstar Games")
	copied := func(name, data string) func() bool {
		return func() bool {
			got, err := os.ReadFile(filepath.Join(dest, name))
			return err == nil && string(got) == data
		}
	}

	if err := a.SetWatch("GTA V", "Nope", true); err == nil {
		t.Error("watching a folder that is not in the game should fail")
	}
	if err := a.SetWatch("GTA V", "Rockstar Games", true); err != nil {
		t.Fatal(err)
	}
	defer a.SetWatch("GTA V", "Rockstar Games", false)
	// Turning it on copies what is there.
	if !waitFor(copied("save.dat", "Rockstar Games")) {
		t.Fatal("the folder was not copied when watching started")
	}

	// A new file, a changed file and a new subfolder are copied.
	os.WriteFile(filepath.Join(src, "save.dat"), []byte("changed"), 0o644)
	os.MkdirAll(filepath.Join(src, "Profiles"), 0o755)
	os.WriteFile(filepath.Join(src, "Profiles", "new.dat"), []byte("new"), 0o644)
	if !waitFor(copied("save.dat", "changed")) || !waitFor(copied(`Profiles\new.dat`, "new")) {
		t.Error("changes to the source were not copied")
	}

	// Saving the folders keeps watching.
	a.SaveFolders("GTA V", map[string]FolderSettings{"Rockstar Games": {Source: portable(src), Link: `D:\Saves`}})
	if f := a.folders(filepath.Join(a.settings.MainFolder, "GTA V"))["Rockstar Games"]; !f.Watch {
		t.Error("saving the folders turned watching off")
	}
	os.WriteFile(filepath.Join(src, "after save.dat"), []byte("x"), 0o644)
	if !waitFor(copied("after save.dat", "x")) {
		t.Error("changes were not copied after saving")
	}

	// Once turned off, changes are not copied.
	if err := a.SetWatch("GTA V", "Rockstar Games", false); err != nil {
		t.Fatal(err)
	}
	if len(a.watches) != 0 {
		t.Errorf("still watching %d folders", len(a.watches))
	}
	os.WriteFile(filepath.Join(src, "late.dat"), []byte("late"), 0o644)
	time.Sleep(4 * settleTime)
	if copied("late.dat", "late")() {
		t.Error("a change was copied after watching was turned off")
	}

	// Removing a watched folder stops watching it.
	a.SetWatch("GTA V", "Rockstar Games", true)
	if err := a.RemoveFolder("GTA V", "Rockstar Games"); err != nil {
		t.Fatal(err)
	}
	if len(a.watches) != 0 {
		t.Error("a removed folder is still watched")
	}
}

func TestServiceSettings(t *testing.T) {
	t.Setenv("USERPROFILE", `C:\Users\me`)
	t.Setenv("GAMES", `D:\Games`)
	a := newTestApp(t, t.TempDir())
	a.AddGame("GTA V")
	if err := a.setFolders(filepath.Join(a.settings.MainFolder, "GTA V"), map[string]FolderSettings{
		"Saves": {Source: `%GAMES%\GTA V\Saves`},
	}); err != nil {
		t.Fatal(err)
	}
	// The service reads the variables the app saved, not its own.
	s, err := readSettings(a.settingsPath)
	if err != nil {
		t.Fatal(err)
	}
	if s.Vars["USERPROFILE"] != `C:\Users\me` || s.Vars["GAMES"] != `D:\Games` {
		t.Errorf("Vars = %v", s.Vars)
	}
	// A file saved with a BOM, such as by Notepad, can be read.
	data, _ := os.ReadFile(a.settingsPath)
	os.WriteFile(a.settingsPath, append([]byte("\xef\xbb\xbf"), data...), 0o644)
	if s, err := readSettings(a.settingsPath); err != nil || s.MainFolder != a.settings.MainFolder {
		t.Errorf("reading GameFix.json with a BOM: %v", err)
	}
	// A file the app is still writing is not taken as no settings.
	os.WriteFile(a.settingsPath, data[:len(data)/2], 0o644)
	if _, err := readSettings(a.settingsPath); err == nil {
		t.Error("reading a partly written GameFix.json should fail")
	}
}

func TestWatchDirStops(t *testing.T) {
	dir := t.TempDir()
	stop := make(chan struct{})
	changes := make(chan []string, 10)
	result := make(chan error)
	go func() { result <- watchDir(dir, stop, func(names []string) { changes <- names }) }()
	time.Sleep(100 * time.Millisecond)
	os.MkdirAll(filepath.Join(dir, "Saves"), 0o755)
	os.WriteFile(filepath.Join(dir, "Saves", "slot ü.dat"), []byte("a"), 0o644)
	select {
	case names := <-changes:
		for len(names) < 2 {
			select {
			case more := <-changes:
				names = append(names, more...)
			case <-time.After(time.Second):
				t.Fatalf("changed names = %q, want Saves and its file", names)
			}
		}
		if !slices.Contains(names, "Saves") || !slices.Contains(names, filepath.Join("Saves", "slot ü.dat")) {
			t.Errorf("changed names = %q, want Saves and its file", names)
		}
	case <-time.After(5 * time.Second):
		t.Error("no change was reported")
	}
	close(stop)
	select {
	case err := <-result:
		if err != nil {
			t.Errorf("watchDir() = %v after stopping", err)
		}
	case <-time.After(5 * time.Second):
		t.Error("watchDir did not stop")
	}
	if err := watchDir(filepath.Join(dir, "missing"), stop, func([]string) {}); err == nil {
		t.Error("watching a missing folder should fail")
	}
}
