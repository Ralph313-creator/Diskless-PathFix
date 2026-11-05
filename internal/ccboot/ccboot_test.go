package ccboot

import (
	"os"
	"path/filepath"
	"testing"
)

const bat = `G:\GameFix\GameFix.bat`

func script(t *testing.T, text string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "public.txt")
	if text != "" {
		if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return path
}

func read(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestRunLine(t *testing.T) {
	if got, want := RunLine("-startup -skipupdate", bat), `run -startup -skipupdate G:\GameFix\GameFix.bat`; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	if got, want := RunLine("-startup", `G:\Game Fix\GameFix.bat`), `run -startup "G:\Game Fix\GameFix.bat"`; got != want {
		t.Fatalf("a path with a space is not quoted: got %q, want %q", got, want)
	}
}

func TestHasRunFindsTheCommandWhateverItsFlagsOrQuotes(t *testing.T) {
	for _, text := range []string{
		"run -startup -skipupdate G:\\GameFix\\GameFix.bat\r\n",
		"run -winlogon \"g:\\gamefix\\GAMEFIX.BAT\"\r\n",
		"  RUN   -startup   G:\\GameFix\\GameFix.bat",
	} {
		has, err := HasRun(script(t, text), bat)
		if err != nil || !has {
			t.Fatalf("%q: has=%v err=%v, want true", text, has, err)
		}
	}
}

func TestHasRunIgnoresCommentsAndOtherCommands(t *testing.T) {
	for _, text := range []string{
		"#run -startup -skipupdate G:\\GameFix\\GameFix.bat\r\n",
		"run -startup -skipupdate G:\\GameFix\\GameFix Copy.bat\r\n",
		"run -winlogon \"G:\\Portable\\STOVE\\STOVE.exe\" --setup\r\n",
		"file G:\\GameFix\\GameFix.bat C:\\x.bat\r\n",
	} {
		has, err := HasRun(script(t, text), bat)
		if err != nil || has {
			t.Fatalf("%q: has=%v err=%v, want false", text, has, err)
		}
	}
}

func TestHasRunWithoutScript(t *testing.T) {
	has, err := HasRun(script(t, ""), bat)
	if err != nil || has {
		t.Fatalf("has=%v err=%v, want false and no error", has, err)
	}
}

func TestAddRunGoesAfterTheLastCommandAndKeepsTheRest(t *testing.T) {
	path := script(t, "run -winlogon \"G:\\Portable\\STOVE\\STOVE.exe\" --setup\r\n"+
		"\r\n#run -service command\r\n#run -startup command\r\n")
	added, err := AddRun(path, "-startup -skipupdate", bat)
	if err != nil || !added {
		t.Fatalf("added=%v err=%v", added, err)
	}
	want := "run -winlogon \"G:\\Portable\\STOVE\\STOVE.exe\" --setup\r\n" +
		"run -startup -skipupdate G:\\GameFix\\GameFix.bat\r\n" +
		"\r\n#run -service command\r\n#run -startup command\r\n"
	if got := read(t, path); got != want {
		t.Fatalf("got\n%q\nwant\n%q", got, want)
	}
}

func TestAddRunTwiceAddsOnce(t *testing.T) {
	path := script(t, "run -startup -skipupdate G:\\GameFix\\GameFix.bat\r\n")
	added, err := AddRun(path, "-startup -skipupdate", bat)
	if err != nil || added {
		t.Fatalf("added=%v err=%v, want false: it is already there", added, err)
	}
	if got := read(t, path); got != "run -startup -skipupdate G:\\GameFix\\GameFix.bat\r\n" {
		t.Fatalf("script changed: %q", got)
	}
}

func TestAddRunOnlyComments(t *testing.T) {
	path := script(t, "\xEF\xBB\xBF#run -startup command\n#file localpath remotepath")
	if _, err := AddRun(path, "-startup -skipupdate", bat); err != nil {
		t.Fatal(err)
	}
	want := "\xEF\xBB\xBFrun -startup -skipupdate G:\\GameFix\\GameFix.bat\n#run -startup command\n#file localpath remotepath"
	if got := read(t, path); got != want {
		t.Fatalf("got %q, want %q (on top, after the BOM, with the file's own line ends)", got, want)
	}
}

func TestAddRunLastLineWithoutLineEnd(t *testing.T) {
	path := script(t, "run -startup a.bat")
	if _, err := AddRun(path, "-startup -skipupdate", bat); err != nil {
		t.Fatal(err)
	}
	if got, want := read(t, path), "run -startup a.bat\r\nrun -startup -skipupdate G:\\GameFix\\GameFix.bat"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestAddRunMakesAMissingScript(t *testing.T) {
	path := script(t, "")
	if _, err := AddRun(path, "-startup -skipupdate", bat); err != nil {
		t.Fatal(err)
	}
	if got, want := read(t, path), "run -startup -skipupdate G:\\GameFix\\GameFix.bat\r\n"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestAddRunWithoutCCBoot(t *testing.T) {
	path := filepath.Join(t.TempDir(), `CCBoot\Cmd\public.txt`)
	if _, err := AddRun(path, "-startup -skipupdate", bat); err == nil {
		t.Fatal("made CCBoot's folder on a PC without CCBoot")
	}
}
