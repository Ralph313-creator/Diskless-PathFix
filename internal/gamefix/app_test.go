package gamefix

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRobocopyCopiesEverything(t *testing.T) {
	src := filepath.Join(t.TempDir(), "My Game")
	dst := filepath.Join(t.TempDir(), "My Game")
	files := map[string]string{
		"game.exe":             "exe",
		`data\level1.pak`:      "0123456789",
		`data\saves\slot1.sav`: "save",
	}
	for name, body := range files {
		path := filepath.Join(src, name)
		os.MkdirAll(filepath.Dir(path), 0o755)
		os.WriteFile(path, []byte(body), 0o644)
	}
	os.MkdirAll(filepath.Join(src, "empty"), 0o755)

	var done int64
	if err := robocopy(src, dst, nil, func(_ string, d int64) { done = d }); err != nil {
		t.Fatal(err)
	}
	for name, body := range files {
		got, err := os.ReadFile(filepath.Join(dst, name))
		if err != nil || string(got) != body {
			t.Errorf("%s = %q, %v; want %q", name, got, err, body)
		}
	}
	if st, err := os.Stat(filepath.Join(dst, "empty")); err != nil || !st.IsDir() {
		t.Error("empty folder was not copied")
	}
	if want := folderSize(src); done != want {
		t.Errorf("progress done = %d, want %d", done, want)
	}
}

func TestCopyMakesTheCopyTheSameAsTheSource(t *testing.T) {
	a := newTestApp(t, t.TempDir())
	a.AddGame("GTA V")
	game := filepath.Join(a.settings.MainFolder, "GTA V")
	src := filepath.Join(t.TempDir(), "GTA V") // the game folder itself
	for _, name := range []string{"game.exe", `x64\old.rpf`, `x64\keep.rpf`} {
		os.MkdirAll(filepath.Dir(filepath.Join(src, name)), 0o755)
		os.WriteFile(filepath.Join(src, name), []byte(name), 0o644)
	}
	saves := makeSource(t, "Rockstar Games")
	a.AddFolder("GTA V", src)
	a.AddFolder("GTA V", saves)
	for _, name := range []string{GameItself, "Rockstar Games"} {
		if _, err := a.CopyFolder("GTA V", name); err != nil {
			t.Fatal(err)
		}
	}

	// Deleted from the sources, including a read-only file and a folder.
	os.Remove(filepath.Join(src, `x64\old.rpf`))
	os.WriteFile(filepath.Join(game, "extra.txt"), []byte("only in the copy"), 0o444)
	os.Remove(filepath.Join(saves, "save.dat"))
	os.WriteFile(filepath.Join(saves, "new.dat"), []byte("new"), 0o644)
	for _, name := range []string{GameItself, "Rockstar Games"} {
		if _, err := a.CopyFolder("GTA V", name); err != nil {
			t.Fatal(err)
		}
	}

	for _, gone := range []string{`x64\old.rpf`, "extra.txt", `Rockstar Games\save.dat`} {
		if _, err := os.Stat(filepath.Join(game, gone)); !os.IsNotExist(err) {
			t.Errorf("%s is not in the source, but was kept in the copy", gone)
		}
	}
	// The game's other folder is inside the game folder, but not part of
	// its source, so copying the game folder must not delete it.
	for _, kept := range []string{"game.exe", `x64\keep.rpf`, `Rockstar Games\new.dat`} {
		if _, err := os.Stat(filepath.Join(game, kept)); err != nil {
			t.Errorf("%s was deleted from the copy", kept)
		}
	}
	if _, err := a.CopyFolder("GTA V", GameItself); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(game, `Rockstar Games\new.dat`)); err != nil {
		t.Error("copying the game folder deleted the copy of its other folder")
	}
}

func TestGamesHaveTheirOwnFolders(t *testing.T) {
	a := newTestApp(t, t.TempDir())
	for _, bad := range []string{"", "  ", `a\b`, "a:b", "what?", "dots.", ".."} {
		if _, err := a.AddGame(bad); err == nil {
			t.Errorf("AddGame(%q) should fail", bad)
		}
	}
	if _, err := a.AddGame("GTA V"); err != nil {
		t.Fatal(err)
	}
	if _, err := a.AddGame("GTA V"); err == nil {
		t.Error("adding the same game twice should fail")
	}
	if err := a.AddFolder("GTA V", makeSource(t, "Rockstar Games")); err != nil {
		t.Fatal(err)
	}

	games, err := a.ListGames()
	if err != nil {
		t.Fatal(err)
	}
	if len(games) != 1 || games[0].Name != "GTA V" || len(games[0].Folders) != 1 || games[0].Folders[0].Name != "Rockstar Games" {
		t.Errorf("ListGames() = %+v", games)
	}
}

func TestDeleteGame(t *testing.T) {
	a := newTestApp(t, t.TempDir())
	a.AddGame("GTA V")
	a.AddGame("Valorant")
	file := filepath.Join(a.settings.MainFolder, "GTA V", "Rockstar Games", "save.dat")
	os.MkdirAll(filepath.Dir(file), 0o755)
	os.WriteFile(file, []byte("save"), 0o444) // read-only, like many copied game files

	if err := a.DeleteGame("GTA V"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(a.settings.MainFolder, "GTA V")); !os.IsNotExist(err) {
		t.Error("game folder still exists")
	}
	if _, err := os.Stat(filepath.Join(a.settings.MainFolder, "Valorant")); err != nil {
		t.Error("other game was deleted too")
	}
	for _, bad := range []string{"", "..", `..\x`, "GTA V"} {
		if err := a.DeleteGame(bad); err == nil {
			t.Errorf("DeleteGame(%q) should fail", bad)
		}
	}
}

func TestInside(t *testing.T) {
	tests := []struct {
		path, dir string
		want      bool
	}{
		{`D:\Games`, `D:\Games`, true},
		{`D:\Games\GTA V`, `D:\Games`, true},
		{`D:\Games2`, `D:\Games`, false},
		{`D:\Other`, `D:\Games`, false},
		{`E:\Games`, `D:\Games`, false},
	}
	for _, tt := range tests {
		if got := inside(tt.path, tt.dir); got != tt.want {
			t.Errorf("inside(%q, %q) = %v, want %v", tt.path, tt.dir, got, tt.want)
		}
	}
}
