package gamefix

import (
	"path/filepath"

	"raptech-portable/internal/ccboot"
)

// CCBootLine is the line for C:\CCBoot\Cmd\public.txt that makes every floor
// PC run the main folder's GameFix.bat, for the user to copy and paste there
// (the app never writes that file). It is a -service line: run at logon as
// the customer (-startup), GameFix.bat had no admin rights for mklink and
// could not get them by itself. As SYSTEM it skips link paths in user
// folders (see inUserFolder). Empty when no main folder is chosen.
func (a *App) CCBootLine() string {
	main := a.GetMainFolder()
	if main == "" {
		return ""
	}
	return ccboot.RunLine(ccboot.Service, filepath.Join(main, BatName))
}
