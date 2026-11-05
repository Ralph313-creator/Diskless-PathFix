package gamefix

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	BatName     = "GameFix.bat"      // makes the symbolic links of every game; kept in the main folder
	CopyBatName = "GameFix Copy.bat" // copies every game to its link path instead of linking; kept in the main folder
	GameItself  = "."                // folder name for a folder copied straight into the game folder

	oldSetupBatName = "GameFix Client Setup.bat" // written in the main folder by an earlier version, and removed

	// BatBackup is added to a GameFix.bat or GameFix Copy.bat in the main
	// folder that this tab did not make (such as the standalone GameFix's),
	// when it is set aside before the tab writes its own.
	BatBackup = ".before-raptech"
	// ourMark is in every script this tab writes.
	ourMark = "(RapTech Portable)"
)

// batLink is one folder of a game that clients get at its link path:
// GameFix.bat links it there, and GameFix Copy.bat copies it there.
type batLink struct {
	Link   string // where the link is made
	Source string // what it points to, relative to the main folder: <game> or <game>\<folder>
}

// gameFolders lists the folders of a game, and whether each has been copied.
func gameFolders(gameFolder string, settings map[string]FolderSettings) []Folder {
	folders := []Folder{}
	for name, s := range settings {
		copied := isDir(copyPath(gameFolder, name))
		if name == GameItself {
			// The game folder always exists, so look for files that are not
			// the game's other folders.
			entries, _ := os.ReadDir(gameFolder)
			copied = false
			for _, e := range entries {
				if _, other := settings[e.Name()]; !other {
					copied = true
				}
			}
		}
		folders = append(folders, Folder{Name: name, Source: s.Source, Link: s.Link, Copied: copied, Watch: s.Watch,
			UserFolder: inUserFolder(s.Link)})
	}
	sort.Slice(folders, func(i, j int) bool { return strings.ToLower(folders[i].Name) < strings.ToLower(folders[j].Name) })
	return folders
}

// copyPath is where a folder of a game is copied to: <game folder>\<name>,
// or the game folder itself for GameItself.
func copyPath(gameFolder, name string) string {
	if name == GameItself {
		return gameFolder
	}
	return filepath.Join(gameFolder, name)
}

// links lists the folders of every game that are on (Enable in the app),
// have been copied and have a link path, which are the folders clients get
// with GameFix.bat and GameFix Copy.bat. It is called with a.mu held.
func (a *App) links() ([]batLink, error) {
	mainFolder := a.settings.MainFolder
	entries, err := os.ReadDir(mainFolder)
	if err != nil {
		return nil, err
	}
	var links []batLink
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		game := filepath.Join(mainFolder, e.Name())
		for _, f := range gameFolders(game, a.folders(game)) {
			link, err := cleanLink(f.Link)
			if err != nil || link == "" || !f.Copied || !f.Watch {
				continue
			}
			l := batLink{Link: link, Source: e.Name()}
			if f.Name != GameItself {
				l.Source += `\` + f.Name
			}
			links = append(links, l)
		}
	}
	sort.SliceStable(links, func(i, j int) bool { return strings.ToLower(links[i].Source) < strings.ToLower(links[j].Source) })
	return links, nil
}

// writeBat writes <main folder>\GameFix.bat with the links of every game,
// and <main folder>\GameFix Copy.bat, which copies them instead.
func (a *App) writeBat() error {
	mainFolder := a.settings.MainFolder
	if mainFolder == "" {
		return nil // no main folder chosen yet: nowhere to write them
	}
	removeOldBat(filepath.Join(mainFolder, oldSetupBatName))
	if entries, err := os.ReadDir(mainFolder); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				removeOldBat(filepath.Join(mainFolder, e.Name(), BatName))
			}
		}
	}
	links, err := a.links()
	if err != nil {
		return err
	}
	for name, body := range map[string]string{BatName: makeBat(links), CopyBatName: makeCopyBat(links)} {
		path := filepath.Join(mainFolder, name)
		if err := setAsideForeign(path); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			return err
		}
	}
	return nil
}

// setAsideForeign renames a script at path that this tab did not write to
// path+BatBackup (or path+BatBackup+".2", ".3"... when that is taken), so
// writing the tab's own never loses another app's links.
func setAsideForeign(path string) error {
	data, err := os.ReadFile(path)
	if err != nil || strings.Contains(string(data), ourMark) {
		return nil // missing, or the tab's own: rewritten
	}
	backup := path + BatBackup
	for i := 2; ; i++ {
		if _, err := os.Stat(backup); os.IsNotExist(err) {
			break
		}
		backup = fmt.Sprintf("%s%s.%d", path, BatBackup, i)
	}
	return os.Rename(path, backup)
}

// removeOldBat deletes a script that older versions wrote, such as the
// GameFix.bat in each game folder, but only if GameFix made it.
func removeOldBat(path string) {
	if data, err := os.ReadFile(path); err == nil && strings.HasPrefix(string(data), "@echo off\r\nrem Made by GameFix.") {
		os.Remove(path)
	}
}

var linkStart = regexp.MustCompile(`^([A-Za-z]:\\.|\\\\[^\\]+\\.|%[A-Za-z_][A-Za-z0-9_()]*%)`)

// pathVars are the variables link paths are written with, so a link works
// for whoever is logged on to the client. Longer values are more specific,
// so they are tried first: %LOCALAPPDATA% before %USERPROFILE%.
var pathVars = []string{"LOCALAPPDATA", "APPDATA", "USERPROFILE", "PUBLIC", "ProgramData", "ProgramFiles(x86)", "ProgramFiles"}

// portable writes the start of a path as a variable when it is one of the
// pathVars folders: C:\Users\me\Documents\Game -> %USERPROFILE%\Documents\Game.
func portable(path string) string {
	best, bestValue := "", ""
	for _, name := range pathVars {
		value := strings.TrimRight(os.Getenv(name), `\`)
		if value == "" || len(value) <= len(bestValue) {
			continue
		}
		if strings.EqualFold(path, value) || strings.HasPrefix(strings.ToLower(path), strings.ToLower(value)+`\`) {
			best, bestValue = name, value
		}
	}
	if best == "" {
		return path
	}
	return "%" + best + "%" + path[len(bestValue):]
}

var pathVar = regexp.MustCompile(`%([A-Za-z_][A-Za-z0-9_()]*)%`)

// expandVars replaces %NAME% with the variable's value on this PC, the way
// cmd does. Unknown variables are left as they are.
func expandVars(path string) string {
	return pathVar.ReplaceAllStringFunc(path, func(m string) string {
		if value, ok := os.LookupEnv(m[1 : len(m)-1]); ok {
			return value
		}
		return m
	})
}

// cleanLink checks a link path and writes it with a variable where it can
// (see portable). "" means the folder has no link.
func cleanLink(link string) (string, error) {
	link = strings.TrimRight(strings.ReplaceAll(strings.TrimSpace(link), "/", `\`), `\`)
	if link == "" {
		return "", nil
	}
	if strings.ContainsAny(link, `"<>|?*`) {
		return "", errors.New(`a link path cannot contain " < > | ? *`)
	}
	if !linkStart.MatchString(link) {
		return "", errors.New(`the link path must be a full path, such as C:\Users\me\Documents\Game or %USERPROFILE%\Documents\Game`)
	}
	link = portable(link)
	if protected(link) {
		return "", errors.New(link + " is a Windows or user folder. GameFix.bat and GameFix Copy.bat replace what is at the link path, so use the game's own folder inside it")
	}
	return link, nil
}

// protected reports whether link is a Windows or user folder. Those can
// never be link paths, because GameFix.bat deletes the folder at the link
// path, and GameFix Copy.bat deletes what is in it but not in the copy.
func protected(link string) bool {
	link = strings.ToLower(link)
	var patterns []string
	for _, home := range []string{`?:\users\*`, `%userprofile%`} {
		patterns = append(patterns, home)
		for _, sub := range userFolders {
			patterns = append(patterns, home+`\`+sub)
		}
	}
	for _, sub := range userFolders {
		patterns = append(patterns, `%public%\`+sub)
	}
	for _, pattern := range append(patterns, systemFolders...) {
		if ok, _ := filepath.Match(pattern, link); ok {
			return true
		}
	}
	return false
}

var userFolders = []string{
	"desktop", "documents", "downloads", "music", "pictures", "videos", "saved games", "favorites",
	"links", "contacts", "searches", "3d objects", "onedrive", "appdata", `appdata\local`,
	`appdata\locallow`, `appdata\roaming`,
}

var systemFolders = []string{
	`?:\windows`, `?:\program files`, `?:\program files (x86)`, `?:\programdata`, `?:\users`,
	`%appdata%`, `%localappdata%`, `%programdata%`, `%public%`, `%homepath%`, `%windir%`,
	`%systemroot%`, `%systemdrive%`, `%programfiles%`, `%programfiles(x86)%`, `%onedrive%`,
}

// inUserFolder reports whether a link path starts with a variable for the
// customer's own folders: %USERPROFILE%, %APPDATA%, %LOCALAPPDATA%. Those are
// only the customer's when the script runs as the customer. CCBoot runs
// GameFix.bat as SYSTEM (run -service), so the scripts skip them then (see
// userLinks). A full path such as C:\Users\RapTech\... is the same folder
// for SYSTEM, and %PUBLIC% is everyone's, so those are done.
func inUserFolder(link string) bool {
	l := strings.ToLower(strings.TrimSpace(link))
	for _, v := range []string{`%userprofile%`, `%appdata%`, `%localappdata%`, `%homedrive%%homepath%`, `%homepath%`, `%onedrive%`} {
		if l == v || strings.HasPrefix(l, v+`\`) {
			return true
		}
	}
	return false
}

// callLink is the line that does one link path with label (link or copy):
// call :<label>, or call :user<label> for a link path in a user folder.
func callLink(label string, l batLink) string {
	if inUserFolder(l.Link) {
		label = "user" + label
	}
	return `call :` + label + ` "` + l.Link + `" "%~dp0` + l.Source + `"`
}

// systemCheck sets SYSTEMRUN when the script runs as SYSTEM: CCBoot's
// service runs it that way (run -service), before anyone logs on. Then
// %USERPROFILE% and the other user folders are SYSTEM's own, not a customer's.
func systemCheck(line func(string)) {
	line("set SKIPPED=0")
	line("rem Run by CCBoot's service (run -service), it is SYSTEM, before anyone logs on: user folders are then SYSTEM's own.")
	line(`set "SYSTEMRUN="`)
	line(`if /i "%USERNAME%"=="%COMPUTERNAME%$" set SYSTEMRUN=1`)
	line(`if /i "%USERPROFILE:~-13%"=="systemprofile" set SYSTEMRUN=1`)
}

// userLinks writes :user<label>, which does a link path in a user folder
// with :<label>, or skips it when the script runs as SYSTEM.
func userLinks(line func(string), label string) {
	line(`rem :user` + label + ` "<link path in a user folder>" "<folder>"`)
	line(":user" + label)
	line("if not defined SYSTEMRUN goto " + label)
	line(`echo [SKIPPED] "%~1": it is in a user folder, and this runs as SYSTEM (CCBoot -service), before anyone logs on.`)
	line("set /a SKIPPED+=1")
	line("exit /b")
}

// makeBat writes a script that makes a symbolic link at each link path,
// pointing to the copy in the main folder, next to the script:
// mklink /d "<link>" "%~dp0<game>\<folder>". Running it from the network
// share makes the links point to the share.
func makeBat(links []batLink) string {
	var b strings.Builder
	line := func(s string) { b.WriteString(s + "\r\n") }

	startBat(line, []string{
		"rem Made by GameFix (RapTech Portable). Makes a symbolic link at each link path, pointing to the",
		"rem game folders next to this file. Run it from the main folder's network share.",
	}, "rem mklink needs administrator rights.")
	line("set OK=0")
	line("set FAILED=0")
	systemCheck(line)
	if len(links) == 0 {
		line("echo No link paths yet. Set them in the SymbolicLink tab.")
	}
	for _, l := range links {
		line(callLink("link", l))
	}
	line("echo.")
	line("echo Done: %OK% linked, %FAILED% failed, %SKIPPED% skipped.")
	line("exit /b %FAILED%")
	line("")
	userLinks(line, "link")
	line(`rem :link "<link path>" "<folder>"`)
	line(":link")
	line(`set "LINK=%~1"`)
	line(`set "SRC=%~2"`)
	line(`set "PARENT=%~dp1"`)
	line(`if not exist "%SRC%\" goto missing`)
	line("rem Remove an old link. rmdir only removes the link itself, or an empty folder.")
	line(`fsutil reparsepoint query "%LINK%" >nul 2>&1 && rmdir "%LINK%"`)
	line("rem Delete a real folder or file that is in the way, so only the link is left.")
	line(`if exist "%LINK%\" rd /s /q "%LINK%" >nul 2>&1`)
	line(`if exist "%LINK%" del /f /q /a "%LINK%" >nul 2>&1`)
	line(`if exist "%LINK%" goto inuse`)
	line(`if not exist "%PARENT%" mkdir "%PARENT%"`)
	line(`mklink /d "%LINK%" "%SRC%" >nul 2>&1 || goto failed`)
	line(`echo [OK]     "%LINK%" -^> "%SRC%"`)
	line("set /a OK+=1")
	line("exit /b")
	line(":missing")
	line(`echo [FAILED] "%LINK%": missing "%SRC%"`)
	line("set /a FAILED+=1")
	line("exit /b")
	line(":inuse")
	line(`echo [FAILED] "%LINK%": it is in the way and could not be deleted. Is the game running?`)
	line("set /a FAILED+=1")
	line("exit /b")
	line(":failed")
	line(`echo [FAILED] "%LINK%": mklink could not make the link`)
	line("set /a FAILED+=1")
	line("exit /b")
	return b.String()
}

// startBat writes the start of both scripts: what the script does, then
// getting administrator rights (whyAdmin says why), then running :main,
// which the caller writes next. CCBoot runs GameFix.bat on every client
// with nobody at the keyboard (run -service: as SYSTEM, which already has
// the rights; run by hand or with -startup it asks for them once), so the
// script never waits for a key, and what it did is kept next to it in
// <script name>.log (G:\GameFix\GameFix.log), to read on the client.
func startBat(line func(string), about []string, whyAdmin string) {
	line("@echo off")
	for _, s := range about {
		line(s)
	}
	line("chcp 65001 >nul")
	line("")
	line(whyAdmin)
	line("rem fltmc only works with them. (net session also needs the Server service, often off on diskless PCs.)")
	line("fltmc >nul 2>&1 && goto admin")
	line("rem Ask once: started again with them, it has \"elevated\" and never asks again.")
	line(`if /i "%~1"=="elevated" goto noadmin`)
	line(`set "SELF=%~f0"`)
	line(`powershell -NoProfile -Command "Start-Process -FilePath $env:SELF -ArgumentList 'elevated' -Verb RunAs"`)
	line("if errorlevel 1 goto noadmin")
	line("exit /b")
	line(":noadmin")
	line(`(echo [FAILED] %DATE% %TIME% %COMPUTERNAME% %USERNAME%: no administrator rights, so nothing was done.>"%~dpn0.log") 2>nul`)
	line("echo [FAILED] No administrator rights, so nothing was done.")
	line("timeout /t 30 >nul 2>&1")
	line("exit /b 1")
	line("")
	line(":admin")
	line(`set "LOG=%~dpn0.log"`)
	line(`(echo.>"%LOG%") 2>nul || set "LOG=%TEMP%\%~n0.log"`)
	line(`call :main %* >"%LOG%" 2>&1`)
	line("set RESULT=%errorlevel%")
	line(`type "%LOG%"`)
	line("rem The window closes by itself; a failure stays on screen longer.")
	line("if %RESULT%==0 (timeout /t 5 >nul 2>&1) else (timeout /t 30 >nul 2>&1)")
	line("exit /b %RESULT%")
	line("")
	line(":main")
	line("echo %DATE% %TIME% on %COMPUTERNAME% as %USERNAME%")
}

// makeCopyBat writes a script like makeBat's, but instead of a link it
// copies the game folder next to the script to each link path, so the
// client has its own copy that works without the share:
// robocopy "%~dp0<game>\<folder>" "<link>" /MIR. The copy at the link path
// ends up the same as the one in the main folder: changed files are
// updated and files that are not in the main folder are deleted.
func makeCopyBat(links []batLink) string {
	var b strings.Builder
	line := func(s string) { b.WriteString(s + "\r\n") }

	startBat(line, []string{
		"rem Made by GameFix (RapTech Portable). Copies the game folders next to this file to each link path,",
		"rem so each client has its own copy instead of a link. Run it from the main folder's network share.",
	}, "rem Link paths can be in folders such as Program Files, which need administrator rights.")
	line("set OK=0")
	line("set FAILED=0")
	systemCheck(line)
	if len(links) == 0 {
		line("echo No link paths yet. Set them in the SymbolicLink tab.")
	}
	for _, l := range links {
		line(callLink("copy", l))
	}
	line("echo.")
	line("echo Done: %OK% copied, %FAILED% failed, %SKIPPED% skipped.")
	line("exit /b %FAILED%")
	line("")
	userLinks(line, "copy")
	line(`rem :copy "<link path>" "<folder>"`)
	line(":copy")
	line(`set "LINK=%~1"`)
	line(`set "SRC=%~2"`)
	line(`if not exist "%SRC%\" goto missing`)
	line("rem Remove a link made by GameFix.bat first, or the copy would go into the main folder.")
	line("rem rmdir only removes the link itself, or an empty folder.")
	line(`fsutil reparsepoint query "%LINK%" >nul 2>&1 && rmdir "%LINK%"`)
	line(`if exist "%LINK%" if not exist "%LINK%\" del /f /q /a "%LINK%" >nul 2>&1`)
	line(`if exist "%LINK%" if not exist "%LINK%\" goto inuse`)
	line(`echo Copying "%SRC%" -^> "%LINK%" ...`)
	line("rem /MIR: every file and subfolder, and delete what is not in the main folder. /XJ: skip junctions, which can loop.")
	line(`robocopy "%SRC%" "%LINK%" /MIR /COPY:DAT /DCOPY:T /XJ /R:2 /W:1 /NFL /NDL /NJH /NJS /NP >nul`)
	line("if errorlevel 8 goto failed")
	line(`echo [OK]     "%LINK%"`)
	line("set /a OK+=1")
	line("exit /b")
	line(":missing")
	line(`echo [FAILED] "%LINK%": missing "%SRC%"`)
	line("set /a FAILED+=1")
	line("exit /b")
	line(":inuse")
	line(`echo [FAILED] "%LINK%": a file is in the way and could not be deleted`)
	line("set /a FAILED+=1")
	line("exit /b")
	line(":failed")
	line(`echo [FAILED] "%LINK%": some files could not be copied. Is the game running?`)
	line("set /a FAILED+=1")
	line("exit /b")
	return b.String()
}
