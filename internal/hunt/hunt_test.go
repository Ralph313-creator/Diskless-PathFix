package hunt

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/sys/windows"
)

func TestSuggestFolder(t *testing.T) {
	roots := defaultRoots()
	env := os.Getenv
	for path, want := range map[string]string{
		filepath.Join(env("LOCALAPPDATA"), `Riot Games\Riot Client\Logs\a.log`):          filepath.Join(env("LOCALAPPDATA"), "Riot Games"),
		filepath.Join(env("USERPROFILE"), `AppData\LocalLow\Studio\Game\save.dat`):       filepath.Join(env("USERPROFILE"), `AppData\LocalLow\Studio`),
		filepath.Join(env("USERPROFILE"), `Documents\My Games\Skyrim\Saves\x.ess`):       filepath.Join(env("USERPROFILE"), `Documents\My Games\Skyrim`),
		filepath.Join(env("USERPROFILE"), `.minecraft\options.txt`):                      filepath.Join(env("USERPROFILE"), ".minecraft"),
		filepath.Join(env("ProgramFiles"), `Roblox\Versions\version-1\RobloxPlayer.exe`): filepath.Join(env("ProgramFiles"), "Roblox"),
		`Q:\Riot Games\VALORANT\live\ShooterGame.exe`:                                    `Q:\Riot Games`,
		`Q:\Games\VALORANT\live\x.pak`:                                                   `Q:\Games\VALORANT`,
		`Q:\SteamLibrary\steamapps\common\Dota 2\game\dota.exe`:                          `Q:\SteamLibrary\steamapps\common\Dota 2`,
		env("LOCALAPPDATA"): "",
		`Q:\`:               "",
	} {
		if got := suggestFolder(path, roots); !strings.EqualFold(got, want) {
			t.Errorf("suggestFolder(%q) = %q, want %q", path, got, want)
		}
	}
}

func TestIsNoise(t *testing.T) {
	env := os.Getenv
	for folder, want := range map[string]bool{
		filepath.Join(env("SystemRoot"), "System32"):             true,
		filepath.Join(env("LOCALAPPDATA"), "Microsoft"):          true,
		filepath.Join(env("LOCALAPPDATA"), "Temp"):               true,
		filepath.Join(env("LOCALAPPDATA"), "NVIDIA Corporation"): true,
		filepath.Join(env("LOCALAPPDATA"), "RapTech Portable"):   true,
		filepath.Join(env("LOCALAPPDATA"), "Riot Games"):         false,
		filepath.Join(env("ProgramFiles"), "Roblox"):             false,
		`D:\Temple Run`: false,
	} {
		if got := isNoise(folder); got != want {
			t.Errorf("isNoise(%q) = %v, want %v", folder, got, want)
		}
	}
}

// huntFor hunts exe with args, watching only root for changes, and returns
// the states it sent. stopAfter > 0 stops it early.
func huntFor(t *testing.T, exe, args, root string, exclude []string, added []string, stopAfter time.Duration) (*Hunter, []State) {
	t.Helper()
	scanEvery, grace = 100*time.Millisecond, 500*time.Millisecond
	t.Cleanup(func() { scanEvery, grace = time.Second, 5*time.Second })
	var mu sync.Mutex
	var states []State
	h, err := New(exe, exclude, func() []string { return added }, func(st State) {
		mu.Lock()
		states = append(states, st)
		mu.Unlock()
	})
	if err != nil {
		t.Fatal(err)
	}
	h.args = args
	h.roots = append(h.roots, root)
	h.drives = []string{root}
	h.Start()
	if stopAfter > 0 {
		time.Sleep(stopAfter)
		h.Stop()
	}
	select {
	case <-h.Done():
	case <-time.After(30 * time.Second):
		h.Stop()
		t.Fatal("hunt did not stop")
	}
	mu.Lock()
	defer mu.Unlock()
	return h, states
}

func folderIn(st State, real string) *Folder {
	for i, f := range st.Folders {
		if strings.EqualFold(f.Real, real) {
			return &st.Folders[i]
		}
	}
	return nil
}

func TestHuntFindsProgramsAndChangedFolders(t *testing.T) {
	root := longPath(t.TempDir())
	saves := filepath.Join(root, "My Game", "Saves")
	skipped := filepath.Join(root, "Portable")
	already := filepath.Join(root, "Old Game")
	os.MkdirAll(already, 0o755)
	cmd := filepath.Join(os.Getenv("SystemRoot"), "System32", "cmd.exe")
	// cmd writes a save, writes into the packages folder and a folder the
	// package already has, and starts ping, which must be followed too.
	args := `/c mkdir "` + saves + `" & echo x> "` + saves + `\slot1.dat" & mkdir "` + skipped + `" & echo y> "` + skipped + `\a" & echo z> "` + already + `\b" & ping -n 2 127.0.0.1 >nul`
	_, states := huntFor(t, cmd, args, root, []string{skipped}, []string{already}, 0)

	if len(states) < 2 || !states[0].Running {
		t.Fatalf("states = %+v, want running ones first", states)
	}
	st := states[len(states)-1]
	if st.Running || st.Error != "" {
		t.Fatalf("last state: running %v, error %q", st.Running, st.Error)
	}
	for _, want := range []string{"cmd.exe", "ping.exe"} {
		if !strings.Contains(strings.ToLower(strings.Join(st.Programs, " ")), want) {
			t.Errorf("programs = %v, want %s", st.Programs, want)
		}
	}
	game := folderIn(st, filepath.Join(root, "My Game"))
	if game == nil || game.Changes == 0 || game.Noise || game.Program || game.Added {
		t.Errorf("My Game = %+v, want changes, not noise, not added", game)
	}
	if f := folderIn(st, skipped); f != nil {
		t.Errorf("the packages folder was suggested: %+v", f)
	}
	if f := folderIn(st, already); f == nil || !f.Added {
		t.Errorf("Old Game = %+v, want shown as already added", f)
	}
	system := folderIn(st, filepath.Join(os.Getenv("SystemRoot"), "System32"))
	if system == nil || !system.Program || system.Noise {
		t.Errorf("System32 = %+v, want the chosen program's folder, never noise", system)
	}
}

func TestStopLeavesTheProgramRunning(t *testing.T) {
	ping := filepath.Join(os.Getenv("SystemRoot"), "System32", "PING.EXE")
	h, states := huntFor(t, ping, "-n 30 127.0.0.1", t.TempDir(), nil, nil, 1500*time.Millisecond)
	if st := states[len(states)-1]; st.Running || st.Live != 1 {
		t.Errorf("after Stop: running %v, %d live; want stopped with ping still live", st.Running, st.Live)
	}
	for pid := range h.procs {
		if p, err := windows.OpenProcess(windows.PROCESS_TERMINATE, false, pid); err == nil {
			windows.TerminateProcess(p, 0)
			windows.CloseHandle(p)
		} else {
			t.Errorf("ping was not left running: %v", err)
		}
	}
}

func TestNewRefusesNonPrograms(t *testing.T) {
	if _, err := New(`C:\Games\notes.txt`, nil, nil, nil); err == nil {
		t.Error("a file that is not a program was accepted")
	}
}

func TestLongPath(t *testing.T) {
	long := filepath.Join(t.TempDir(), "A Long Folder Name", "Saves")
	os.MkdirAll(long, 0o755)
	buf := make([]uint16, windows.MAX_LONG_PATH)
	p, _ := windows.UTF16PtrFromString(long)
	n, err := windows.GetShortPathName(p, &buf[0], uint32(len(buf)))
	short := windows.UTF16ToString(buf[:n])
	if err != nil || !strings.Contains(short, "~") {
		t.Skipf("no short names on this drive: %q, %v", short, err)
	}
	got := longPath(short)
	if strings.Contains(got, "~") || !strings.HasSuffix(got, `\A Long Folder Name\Saves`) {
		t.Errorf("longPath(%q) = %q, want the full names", short, got)
	}
	if gone := longPath(filepath.Join(short, "gone.dat")); gone != filepath.Join(got, "gone.dat") {
		t.Errorf("longPath of a file that is gone = %q", gone)
	}
}
