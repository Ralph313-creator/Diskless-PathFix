// Package ccboot makes the lines for CCBoot's command file, which every floor
// PC runs. The app never writes that file: the user pastes the lines into it.
package ccboot

import "strings"

// Script is CCBoot's command file on the CCBoot server.
const Script = `C:\CCBoot\Cmd\public.txt`

// When a floor PC runs a line.
const (
	// Startup: at logon, as the customer, without admin rights (CCBootClient.exe
	// -startup from HKLM Run). For what belongs to the customer: %LOCALAPPDATA%,
	// HKCU.
	Startup = "-startup"
	// Service: by the CCBoot client service, as SYSTEM, with admin rights, before
	// anyone logs on. For what needs admin rights, such as GameFix.bat's links.
	Service = "-service"
)

// RunLine is the line that makes each floor PC run cmd with args: run <when>
// -skipupdate <cmd> [args]. It has the same form as a line the user's CCBoot
// runs (run -startup -skipupdate G:\GameFix\GameFix.bat); a line with
// -winlogon and a quoted path stopped public.txt from running on the clients.
// So cmd is only quoted when it has a space.
func RunLine(when, cmd string, args ...string) string {
	if strings.ContainsAny(cmd, " \t") {
		cmd = `"` + cmd + `"`
	}
	return strings.Join(append([]string{"run", when, "-skipupdate", cmd}, args...), " ")
}
