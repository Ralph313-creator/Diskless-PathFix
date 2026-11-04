package main

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseArgs(t *testing.T) {
	cases := []struct {
		in       []string
		mode     string
		pkg      string
		passArgs string
	}{
		{nil, "run", "", ""},
		{[]string{"roblox-player:1+launchmode:play"}, "run", "", "roblox-player:1+launchmode:play"},
		{[]string{"--setup"}, "setup", "", ""},
		{[]string{"--pkg", `D:\P\Roblox`, "--update"}, "update", `D:\P\Roblox`, ""},
		{[]string{"--unlink"}, "unlink", "", ""},
		{[]string{"--save"}, "save", "", ""},
		{[]string{"--status"}, "status", "", ""},
		{[]string{"--app", "x"}, "run", "", "--app|x"},
	}
	for _, c := range cases {
		got := parseArgs(c.in)
		if got.mode != c.mode || got.pkg != c.pkg || strings.Join(got.pass, "|") != c.passArgs {
			t.Errorf("parseArgs(%q) = %+v, want mode=%s pkg=%q pass=%q", c.in, got, c.mode, c.pkg, c.passArgs)
		}
	}
}

// buildLauncher builds this program into dir as name.
func buildLauncher(t *testing.T, dir, name string) string {
	t.Helper()
	out := filepath.Join(dir, name)
	cmd := exec.Command("go", "build", "-o", out, ".")
	if b, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, b)
	}
	return out
}

func copyFile(t *testing.T, from, to string) {
	t.Helper()
	in, err := os.Open(from)
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()
	os.MkdirAll(filepath.Dir(to), 0o755)
	out, err := os.Create(to)
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()
	io.Copy(out, in)
}

func TestLauncherRunsPackageEndToEnd(t *testing.T) {
	root := t.TempDir()
	pkg := filepath.Join(root, "pkg")
	pc := filepath.Join(root, "pc")
	launcher := buildLauncher(t, pkg, "App.exe")
	// The "app" is cmd.exe: it writes a login file through the link, then quits.
	copyFile(t, `C:\Windows\System32\cmd.exe`, filepath.Join(pkg, `Data\Program\v1\app.exe`))
	manifest := `{
		"name": "App",
		"exe": "Data\\Program\\*\\app.exe",
		"folders": [{"link": "` + strings.ReplaceAll(pc, `\`, `\\`) + `\\AppData\\App", "data": "Data\\AppData", "session": ["LocalStorage"]}]
	}`
	os.WriteFile(filepath.Join(pkg, "portable.json"), []byte(manifest), 0o644)

	login := filepath.Join(pc, `AppData\App\LocalStorage\login.txt`)
	setting := filepath.Join(pc, `AppData\App\settings.txt`)
	cmd := exec.Command(launcher, "/c", "mkdir "+filepath.Dir(login)+" & echo x>"+login+" & echo y>"+setting)
	if b, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("launcher: %v\n%s", err, b)
	}
	if _, err := os.Stat(filepath.Join(pkg, `Data\AppData\settings.txt`)); err != nil {
		t.Fatal("the app's write did not land in the package:", err)
	}
	if _, err := os.Stat(filepath.Join(pkg, `Data\AppData\LocalStorage\login.txt`)); !os.IsNotExist(err) {
		t.Fatal("session folder not emptied after the app closed")
	}
	logText, _ := os.ReadFile(filepath.Join(pkg, LogName))
	if !strings.Contains(string(logText), "app.exe") {
		t.Fatalf("log does not mention the started program:\n%s", logText)
	}
}
