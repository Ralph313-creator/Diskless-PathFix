// Package ccboot reads and adds to CCBoot's command file, whose lines every
// diskless client runs when it boots.
package ccboot

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Script is CCBoot's command file on the CCBoot server. Its lines are
// "run [-startup|-winlogon|-service] [-skipupdate] command" and
// "file localpath remotepath"; lines starting with # are comments.
const Script = `C:\CCBoot\Cmd\public.txt`

const bom = "\uFEFF"

// RunLine is the line that runs cmd on every client, such as
// run -startup -skipupdate G:\GameFix\GameFix.bat.
func RunLine(flags, cmd string) string {
	if strings.ContainsAny(cmd, " \t") {
		cmd = `"` + cmd + `"`
	}
	return "run " + flags + " " + cmd
}

// HasRun reports whether script has a run line (not commented out) that
// runs cmd, with any flags. A missing script has none.
func HasRun(script, cmd string) (bool, error) {
	b, err := os.ReadFile(script)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	for _, line := range strings.Split(strings.TrimPrefix(string(b), bom), "\n") {
		if runs(line, cmd) {
			return true, nil
		}
	}
	return false, nil
}

// AddRun adds RunLine(flags, cmd) to script after its last command line
// (on top when it has only comments), keeping the file's own line ends.
// Nothing changes when a run line there already runs cmd; added reports
// whether the line was added. The script is made when it is missing, but
// only inside CCBoot's existing folder.
func AddRun(script, flags, cmd string) (added bool, err error) {
	b, err := os.ReadFile(script)
	if errors.Is(err, fs.ErrNotExist) {
		if _, err := os.Stat(filepath.Dir(script)); err != nil {
			return false, fmt.Errorf("CCBoot is not on this PC: %w", err)
		}
	} else if err != nil {
		return false, err
	}
	text := string(b)
	start := ""
	if strings.HasPrefix(text, bom) {
		start, text = bom, text[len(bom):]
	}
	eol := "\r\n"
	if strings.Contains(text, "\n") && !strings.Contains(text, "\r\n") {
		eol = "\n"
	}
	lines := strings.SplitAfter(text, "\n")
	last := -1
	for i, line := range lines {
		if runs(line, cmd) {
			return false, nil
		}
		if s := strings.TrimSpace(line); s != "" && !strings.HasPrefix(s, "#") {
			last = i
		}
	}
	add := RunLine(flags, cmd) + eol
	if last >= 0 && !strings.HasSuffix(lines[last], "\n") {
		lines[last] += eol // it was the file's last line
		add = RunLine(flags, cmd)
	}
	lines = append(lines[:last+1], append([]string{add}, lines[last+1:]...)...)

	// Written next to it, then moved over it, so a client booting now
	// never reads half a file.
	tmp := script + ".rp-tmp"
	if err := os.WriteFile(tmp, []byte(start+strings.Join(lines, "")), 0o644); err != nil {
		return false, err
	}
	if err := os.Rename(tmp, script); err != nil {
		os.Remove(tmp)
		return false, err
	}
	return true, nil
}

// runs reports whether line is a run line that runs cmd: after "run" and
// its -flags, the program is cmd, quoted or not.
func runs(line, cmd string) bool {
	line = strings.TrimSpace(line)
	fields := strings.Fields(line)
	if len(fields) < 2 || !strings.EqualFold(fields[0], "run") {
		return false
	}
	rest := strings.TrimSpace(line[len(fields[0]):])
	for strings.HasPrefix(rest, "-") {
		i := strings.IndexAny(rest, " \t")
		if i < 0 {
			return false
		}
		rest = strings.TrimSpace(rest[i:])
	}
	prog := rest
	if strings.HasPrefix(rest, `"`) {
		if end := strings.Index(rest[1:], `"`); end >= 0 {
			prog = rest[1 : end+1]
		}
	} else if i := strings.IndexAny(rest, " \t"); i >= 0 {
		prog = rest[:i]
	}
	return strings.EqualFold(prog, cmd) || strings.EqualFold(strings.Trim(rest, `"`), cmd)
}
