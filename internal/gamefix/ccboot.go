package gamefix

import (
	"path/filepath"

	"raptech-portable/internal/ccboot"
)

// CCBootLine is the line for C:\CCBoot\Cmd\public.txt that makes every floor
// PC run the main folder's GameFix.bat, for the user to copy and paste there
// (the app never writes that file). It is a -startup line (the user's choice,
// 2026-11-06, after a stretch on -service): CCBoot runs it at logon as the
// customer, so user-folder link paths are the customer's own, and
// GameFix.bat asks for admin rights itself for mklink. Empty when no main
// folder is chosen.
func (a *App) CCBootLine() string {
	main := a.GetMainFolder()
	if main == "" {
		return ""
	}
	return ccboot.RunLine(ccboot.Startup, filepath.Join(main, BatName))
}
