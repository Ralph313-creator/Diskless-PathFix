package manager

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// TaskFolder holds the update tasks in Task Scheduler.
const TaskFolder = `\RapTech Portable\`

// SetSchedule makes (or replaces) a Task Scheduler task that runs
// "<launcher> --update" as SYSTEM every interval, from 1 minute to 24 hours.
// SYSTEM can write the game disk on the server and needs no logged-on user.
func SetSchedule(name, launcher string, every time.Duration) error {
	minutes := int(every / time.Minute)
	if every%time.Minute != 0 || minutes < 1 || minutes > 1439 {
		return fmt.Errorf("update every %v: use whole minutes from 1m to 23h59m", every)
	}
	_, err := schtasks("/Create", "/TN", TaskFolder+name, "/TR", `"`+launcher+`" --update`,
		"/SC", "MINUTE", "/MO", strconv.Itoa(minutes), "/RU", "SYSTEM", "/RL", "HIGHEST", "/F")
	return err
}

// RemoveSchedule deletes the task; a missing task is fine.
func RemoveSchedule(name string) error {
	if on, _ := ScheduleOn(name); !on {
		return nil
	}
	_, err := schtasks("/Delete", "/TN", TaskFolder+name, "/F")
	return err
}

var interval = regexp.MustCompile(`<Interval>PT(?:(\d+)H)?(?:(\d+)M)?</Interval>`)

// ScheduleOn reports whether the task exists and how often it runs.
func ScheduleOn(name string) (bool, time.Duration) {
	out, err := schtasks("/Query", "/TN", TaskFolder+name, "/XML")
	if err != nil {
		return false, 0
	}
	m := interval.FindStringSubmatch(out)
	if m == nil {
		return true, 0
	}
	h, _ := strconv.Atoi(m[1])
	min, _ := strconv.Atoi(m[2])
	return true, time.Duration(h)*time.Hour + time.Duration(min)*time.Minute
}

func schtasks(args ...string) (string, error) {
	cmd := exec.Command("schtasks.exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	raw, err := cmd.CombinedOutput()
	out := strings.ReplaceAll(string(raw), "\x00", "") // /XML may come as UTF-16
	if err != nil {
		return out, fmt.Errorf("schtasks %s: %v %s", args[0], err, strings.TrimSpace(out))
	}
	return out, nil
}
