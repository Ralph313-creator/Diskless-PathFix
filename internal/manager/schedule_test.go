package manager

import (
	"strings"
	"testing"
	"time"
)

func TestDataFor(t *testing.T) {
	cases := map[string]string{
		`%LOCALAPPDATA%\Roblox`:       `Data\LocalAppData\Roblox`,
		`%ProgramFiles(x86)%\Steam`:   `Data\ProgramFilesx86\Steam`,
		`%USERPROFILE%\Saved Games\X`: `Data\UserProfile\Saved Games\X`,
		`D:\Games\GTA V`:              `Data\D\Games\GTA V`,
		`\\SERVER\share\Game`:         `Data\Network\SERVER\share\Game`,
	}
	for in, want := range cases {
		if got := dataFor(in); got != want {
			t.Errorf("dataFor(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestScheduleCreateQueryRemove(t *testing.T) {
	name := "Test-" + strings.ReplaceAll(t.Name(), "/", "_")
	t.Cleanup(func() { RemoveSchedule(name) })
	RemoveSchedule(name)

	if on, _ := ScheduleOn(name); on {
		t.Fatal("task exists before it was made")
	}
	if err := SetSchedule(name, `C:\Windows\System32\cmd.exe`, 45*time.Minute); err != nil {
		t.Fatal(err)
	}
	on, every := ScheduleOn(name)
	if !on || every != 45*time.Minute {
		t.Fatalf("on=%v every=%v, want on every 45m", on, every)
	}
	if err := RemoveSchedule(name); err != nil {
		t.Fatal(err)
	}
	if on, _ := ScheduleOn(name); on {
		t.Fatal("task still there after removing")
	}
}

func TestScheduleRefusesOddTimes(t *testing.T) {
	for _, d := range []time.Duration{0, 30 * time.Second, 25 * time.Hour} {
		if err := SetSchedule("Test-Odd", `C:\Windows\System32\cmd.exe`, d); err == nil {
			RemoveSchedule("Test-Odd")
			t.Errorf("%v accepted", d)
		}
	}
}
