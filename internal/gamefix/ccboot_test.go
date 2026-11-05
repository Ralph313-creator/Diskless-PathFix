package gamefix

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCCBootLineRunsTheMainFoldersGameFixBat(t *testing.T) {
	main := filepath.Join(t.TempDir(), "GameFix")
	a := newTestApp(t, main)
	c := a.GetCCBoot()
	if want := `run -startup -skipupdate ` + filepath.Join(main, BatName); c.Line != want {
		t.Fatalf("line %q, want %q", c.Line, want)
	}
	if c.Added || c.Error != "" || c.Script != a.ccbootScript {
		t.Fatalf("got %+v, want not added yet, no error, script %s", c, a.ccbootScript)
	}
}

func TestAddToCCBootAddsTheLineOnce(t *testing.T) {
	a := newTestApp(t, filepath.Join(t.TempDir(), "GameFix"))
	os.WriteFile(a.ccbootScript, []byte("#run -startup command\r\n"), 0o644)
	for i := 0; i < 2; i++ {
		c, err := a.AddToCCBoot()
		if err != nil || !c.Added {
			t.Fatalf("try %d: %+v, %v", i+1, c, err)
		}
	}
	b, _ := os.ReadFile(a.ccbootScript)
	line := a.GetCCBoot().Line
	if got := strings.Count(string(b), line); got != 1 {
		t.Fatalf("the script has the line %d times, want once:\n%s", got, b)
	}
}

func TestCCBootWithoutCCBoot(t *testing.T) {
	a := newTestApp(t, filepath.Join(t.TempDir(), "GameFix"))
	a.ccbootScript = filepath.Join(t.TempDir(), `CCBoot\Cmd\public.txt`)
	if c := a.GetCCBoot(); c.Error == "" {
		t.Fatal("no error shown on a PC without CCBoot")
	}
	if _, err := a.AddToCCBoot(); err == nil {
		t.Fatal("added to a CCBoot that is not there")
	}
}

func TestAddToCCBootNeedsAMainFolder(t *testing.T) {
	a := newTestApp(t, filepath.Join(t.TempDir(), "GameFix"))
	a.settings.MainFolder = ""
	if _, err := a.AddToCCBoot(); err == nil {
		t.Fatal("added a line without a main folder")
	}
	if b, err := os.ReadFile(a.ccbootScript); err == nil {
		t.Fatalf("script written: %q", b)
	}
}
