package ccboot

import "testing"

func TestRunLine(t *testing.T) {
	for _, c := range []struct {
		cmd  string
		args []string
		want string
	}{
		// The same form as a line CCBoot runs: no quotes, -startup -skipupdate.
		{`G:\GameFix\GameFix.bat`, nil, `run -startup -skipupdate G:\GameFix\GameFix.bat`},
		{`G:\Portable\Roblox\Roblox.exe`, []string{"--setup"}, `run -startup -skipupdate G:\Portable\Roblox\Roblox.exe --setup`},
		// A path with a space has to be quoted.
		{`G:\Online Games\Fix\GameFix.bat`, nil, `run -startup -skipupdate "G:\Online Games\Fix\GameFix.bat"`},
	} {
		if got := RunLine(c.cmd, c.args...); got != c.want {
			t.Errorf("RunLine(%q, %q) = %q, want %q", c.cmd, c.args, got, c.want)
		}
	}
}
