package gamefix

import (
	"path/filepath"
	"testing"
)

// TestCCBootLine: the line the user copies into C:\CCBoot\Cmd\public.txt,
// so every floor PC runs the main folder's GameFix.bat at logon.
func TestCCBootLine(t *testing.T) {
	main := filepath.Join(t.TempDir(), "GameFix")
	a := newTestApp(t, main)
	if want := `run -startup -skipupdate ` + filepath.Join(main, BatName); a.CCBootLine() != want {
		t.Fatalf("CCBootLine() = %q, want %q", a.CCBootLine(), want)
	}
	a.settings.MainFolder = ""
	if got := a.CCBootLine(); got != "" {
		t.Fatalf("no main folder: CCBootLine() = %q, want nothing", got)
	}
}
