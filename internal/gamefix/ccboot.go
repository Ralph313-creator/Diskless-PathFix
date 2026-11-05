package gamefix

import (
	"path/filepath"

	"raptech-portable/internal/ccboot"
)

// CCBootLine is the line for C:\CCBoot\Cmd\public.txt that makes every floor
// PC run the main folder's GameFix.bat at logon, for the user to copy and
// paste there (the app never writes that file). GameFix.bat asks for admin
// rights itself. Empty when no main folder is chosen.
func (a *App) CCBootLine() string {
	main := a.GetMainFolder()
	if main == "" {
		return ""
	}
	return ccboot.RunLine(filepath.Join(main, BatName))
}
