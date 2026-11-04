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

func TestIcoFileRoundTrip(t *testing.T) {
	want, err := Read(source)
	if err != nil {
		t.Fatal(err)
	}
	ico := filepath.Join(t.TempDir(), "app.ico")
	if err := SaveICO(source, ico); err != nil {
		t.Fatal(err)
	}
	// An .ico file works as a source like an exe does.
	target := goExe(t, t.TempDir())
	if err := Copy(ico, target); err != nil {
		t.Fatal(err)
	}
	got, err := Read(target)
	if err != nil || len(got) != len(want) {
		t.Fatalf("%d images (%v), want %d", len(got), err, len(want))
	}
	for i := range want {
		if !bytes.Equal(got[i], want[i]) {
			t.Fatalf("image %d differs after going through an .ico file", i)
		}
	}
	// Windows itself can read the .ico file.
	cmd := exec.Command("powershell", "-NoProfile", "-Command",
		`Add-Type -AssemblyName System.Drawing; $i = New-Object System.Drawing.Icon '`+ico+`'; $i.Width`)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Windows cannot read the .ico: %v %s", err, out)
	}
}

func TestBrokenIcoIsRefused(t *testing.T) {
	bad := filepath.Join(t.TempDir(), "bad.ico")
	os.WriteFile(bad, []byte("not an icon at all"), 0o644)
	target := goExe(t, t.TempDir())
	before, _ := os.ReadFile(target)
	if err := Copy(bad, target); err == nil {
		t.Fatal("a broken .ico was accepted")
	}
	after, _ := os.ReadFile(target)
	if !bytes.Equal(before, after) {
		t.Fatal("the launcher was changed")
	}
}

func TestRemoveTakesTheIconOff(t *testing.T) {
	target := goExe(t, t.TempDir())
	if err := Copy(source, target); err != nil {
		t.Fatal(err)
	}
	if err := Remove(target); err != nil {
		t.Fatal(err)
	}
	if _, err := Read(target); !errors.Is(err, ErrNoIcon) {
		t.Fatalf("Read after Remove: %v, want ErrNoIcon", err)
	}
	if err := Remove(target); err != nil {
		t.Fatalf("Remove on an exe without an icon: %v", err)
	}
	if out, err := exec.Command(target, "-test.run=^$").CombinedOutput(); err != nil {
		t.Fatalf("the exe is broken after Remove: %v\n%s", err, out)
	}
}

func TestPNGPicksTheLargestPNGImage(t *testing.T) {
	png, err := PNG(source)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(png, []byte("\x89PNG")) {
		t.Fatal("not a PNG")
	}
	images, _ := Read(source)
	for _, img := range images {
		if bytes.HasPrefix(img, []byte("\x89PNG")) && len(img) > len(png) {
			t.Fatal("a bigger PNG image was skipped")
		}
	}
	if _, err := PNG(goExe(t, t.TempDir())); !errors.Is(err, ErrNoIcon) {
		t.Fatalf("exe without icon: %v, want ErrNoIcon", err)
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
