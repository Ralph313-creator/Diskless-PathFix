package ccboot

import "testing"

func TestRunLine(t *testing.T) {
	for _, c := range []struct {
		when string
		cmd  string
		args []string
		want string
	}{
		// The same form as a line CCBoot runs: no quotes.
		{Startup, `G:\Portable\Roblox\Roblox.exe`, []string{"--setup"}, `run -startup -skipupdate G:\Portable\Roblox\Roblox.exe --setup`},
		// -service: CCBoot's service runs it as SYSTEM, with admin rights.
		{Service, `G:\GameFix\GameFix.bat`, nil, `run -service -skipupdate G:\GameFix\GameFix.bat`},
		// A path with a space has to be quoted.
		{Service, `G:\Online Games\Fix\GameFix.bat`, nil, `run -service -skipupdate "G:\Online Games\Fix\GameFix.bat"`},
	} {
		if got := RunLine(c.when, c.cmd, c.args...); got != c.want {
			t.Errorf("RunLine(%q, %q, %q) = %q, want %q", c.when, c.cmd, c.args, got, c.want)
		}
	}
}
