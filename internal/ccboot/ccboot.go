// Package ccboot makes the lines for CCBoot's command file, which every floor
// PC runs. The app never writes that file: the user pastes the lines into it.
package ccboot

import "strings"

// Script is CCBoot's command file on the CCBoot server.
const Script = `C:\CCBoot\Cmd\public.txt`

// RunLine is the line that makes each floor PC run cmd with args at logon:
// run -startup -skipupdate <cmd> [args]. It has the same form as a line the
// user's CCBoot runs (run -startup -skipupdate G:\GameFix\GameFix.bat); a
// line with -winlogon and a quoted path stopped public.txt from running on
// the clients. So cmd is only quoted when it has a space.
func RunLine(cmd string, args ...string) string {
	if strings.ContainsAny(cmd, " \t") {
		cmd = `"` + cmd + `"`
	}
	return strings.Join(append([]string{"run -startup -skipupdate", cmd}, args...), " ")
}
