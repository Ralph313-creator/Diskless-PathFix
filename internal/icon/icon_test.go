package icon

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// source is a Windows program that has an icon.
var source = filepath.Join(os.Getenv("SystemRoot"), "explorer.exe")

// goExe copies this test program, a Go exe without an icon, to dir.
func goExe(t *testing.T, dir string) string {
	t.Helper()
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	in, _ := os.Open(self)
	defer in.Close()
	out := filepath.Join(dir, "launcher.exe")
	f, err := os.Create(out)
	if err != nil {
		t.Fatal(err)
	}
	io.Copy(f, in)
	f.Close()
	return out
}

func TestCopyGivesTheLauncherTheAppIcon(t *testing.T) {
	target := goExe(t, t.TempDir())
	if _, err := Read(target); !errors.Is(err, ErrNoIcon) {
		t.Fatalf("a Go exe should start without an icon, got %v", err)
	}
	want, err := Read(source)
	if err != nil {
		t.Fatal(err)
	}
	if err := Copy(source, target); err != nil {
		t.Fatal(err)
	}
	got, err := Read(target)
	if err != nil {
		t.Fatal("no icon after Copy:", err)
	}
	if len(got) != len(want) || len(got) == 0 {
		t.Fatalf("%d images, want %d", len(got), len(want))
	}
	for i := range want {
		if !bytes.Equal(got[i], want[i]) {
			t.Fatalf("image %d differs", i)
		}
	}
	// The launcher must still run.
	cmd := exec.Command(target, "-test.run=^$")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("the exe is broken after Copy: %v\n%s", err, out)
	}
}

func TestCopyTwiceKeepsOneIcon(t *testing.T) {
	target := goExe(t, t.TempDir())
	if err := Copy(source, target); err != nil {
		t.Fatal(err)
	}
	if err := Copy(source, target); err != nil {
		t.Fatal(err)
	}
	groups, err := groupCount(target)
	if err != nil || groups != 1 {
		t.Fatalf("%d icon groups (%v), want 1", groups, err)
	}
}

func TestCopyFromExeWithoutIconChangesNothing(t *testing.T) {
	dir := t.TempDir()
	target := goExe(t, dir)
	before, _ := os.ReadFile(target)
	noIcon := goExe(t, t.TempDir())
	if err := Copy(noIcon, target); !errors.Is(err, ErrNoIcon) {
		t.Fatalf("err = %v, want ErrNoIcon", err)
	}
	after, _ := os.ReadFile(target)
	if !bytes.Equal(before, after) {
		t.Fatal("the launcher was changed")
	}
}
