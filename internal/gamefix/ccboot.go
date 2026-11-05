package gamefix

import (
	"errors"
	"os"
	"path/filepath"

	"raptech-portable/internal/ccboot"
)

// ccbootFlags runs GameFix.bat when each client starts.
const ccbootFlags = "-startup -skipupdate"

// CCBoot is GameFix.bat's line in CCBoot's command file, which every
// client runs when it boots.
type CCBoot struct {
	Line   string `json:"line"`   // run -startup -skipupdate G:\GameFix\GameFix.bat
	Script string `json:"script"` // C:\CCBoot\Cmd\public.txt
	Added  bool   `json:"added"`  // the script already runs GameFix.bat
	Error  string `json:"error"`  // why it can't be added, such as no CCBoot on this PC
}

// GetCCBoot shows GameFix.bat's line and whether CCBoot's script has it.
func (a *App) GetCCBoot() CCBoot {
	return a.ccbootState(a.GetMainFolder())
}

// AddToCCBoot adds GameFix.bat's line to CCBoot's script, once.
func (a *App) AddToCCBoot() (CCBoot, error) {
	main := a.GetMainFolder()
	if main == "" {
		return a.ccbootState(main), errors.New("choose a main folder first")
	}
	_, err := ccboot.AddRun(a.ccbootScript, ccbootFlags, filepath.Join(main, BatName))
	return a.ccbootState(main), err
}

// RemoveFromCCBoot takes GameFix.bat's line out of CCBoot's script, so
// clients no longer run it when they start.
func (a *App) RemoveFromCCBoot() (CCBoot, error) {
	main := a.GetMainFolder()
	if main == "" {
		return a.ccbootState(main), errors.New("choose a main folder first")
	}
	_, err := ccboot.RemoveRun(a.ccbootScript, filepath.Join(main, BatName))
	return a.ccbootState(main), err
}

func (a *App) ccbootState(main string) CCBoot {
	c := CCBoot{Script: a.ccbootScript}
	if main == "" {
		c.Error = "Choose a main folder first."
		return c
	}
	bat := filepath.Join(main, BatName)
	c.Line = ccboot.RunLine(ccbootFlags, bat)
	if _, err := os.Stat(filepath.Dir(a.ccbootScript)); err != nil {
		c.Error = "CCBoot is not on this PC (no " + filepath.Dir(a.ccbootScript) + ")."
		return c
	}
	has, err := ccboot.HasRun(a.ccbootScript, bat)
	if err != nil {
		c.Error = err.Error()
	}
	c.Added = has
	return c
}
