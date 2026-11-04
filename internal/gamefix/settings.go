package gamefix

import (
	"bytes"
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
)

// SettingsName is the file next to RapTech-Portable.exe that holds
// everything the GameFix tab remembers. It is not the standalone GameFix's
// GameFix.json, so the two never share settings.
const SettingsName = "RapTech-GameFix.json"

// defaultMainFolder is used the first time, when it exists. Tests replace it.
var defaultMainFolder = `G:\GameFix`

// Settings are remembered between runs, in SettingsName next to the exe.
type Settings struct {
	MainFolder string `json:"mainFolder"`
	// Games holds the folders of each game:
	// game folder -> folder name (GameItself for the game folder) -> source and link path.
	Games map[string]map[string]FolderSettings `json:"games"`
	// Vars are the values of the variables sources are written with, such
	// as USERPROFILE, for the user running the app. The service runs as
	// SYSTEM, so it expands sources with these instead of its own.
	Vars map[string]string `json:"vars,omitempty"`
}

// FolderSettings is what the app remembers about one folder of a game.
type FolderSettings struct {
	Source string `json:"source"`          // where Copy copies it from, on this PC; may use %VARIABLES%
	Link   string `json:"link"`            // where GameFix.bat makes the symbolic link
	Watch  bool   `json:"watch,omitempty"` // on (Enable in the app): clients get it, and it is copied again whenever the source changes
}

// UnmarshalJSON also reads older versions, which kept only the link path.
// That was where the folder was copied from, unless it was changed, so it
// is taken as the source when it is a folder on this PC.
func (f *FolderSettings) UnmarshalJSON(data []byte) error {
	var link string
	if json.Unmarshal(data, &link) == nil {
		*f = FolderSettings{Link: link}
		if isDir(expandVars(link)) {
			f.Source = link
		}
		return nil
	}
	type plain FolderSettings
	return json.Unmarshal(data, (*plain)(f))
}

func isDir(path string) bool {
	st, err := os.Stat(path)
	return err == nil && st.IsDir()
}

// exeSettingsPath is GameFix.json in the folder GameFix.exe lives in.
func exeSettingsPath() string {
	exe, err := os.Executable()
	if err != nil {
		return SettingsName
	}
	return filepath.Join(filepath.Dir(exe), SettingsName)
}

// load reads the settings. The first time, the main folder is
// defaultMainFolder when it exists. It also takes the link paths from the
// gamefix.json files that old GameFix versions kept in each game folder.
func (a *App) load() {
	var err error
	a.settings, err = readSettings(a.settingsPath)
	// The first time, the tab has no games yet: it must not rewrite a
	// GameFix.bat that is already in the main folder.
	firstStart := os.IsNotExist(err)
	if a.settings.MainFolder == "" && isDir(defaultMainFolder) {
		a.settings.MainFolder = defaultMainFolder
	}
	a.importOldProfiles()
	a.syncWatches()

	// Sources and link paths saved before they were written with variables
	// are updated, and so are the variables the service uses.
	changed := !maps.Equal(a.settings.Vars, savedVars(a.settings))
	for _, folders := range a.settings.Games {
		for name, f := range folders {
			old := f
			if link, err := cleanLink(f.Link); err == nil {
				f.Link = link
			}
			if source, err := checkSource(f.Source, a.settings.MainFolder); err == nil {
				f.Source = source
			}
			if f != old {
				folders[name] = f
				changed = true
			}
		}
	}
	if changed {
		a.save()
		if !firstStart {
			a.writeBat()
		}
	}
}

// readSettings reads a GameFix.json. It fails for a file that is only
// partly written, so the service does not take that as no settings.
func readSettings(path string) (Settings, error) {
	var s Settings
	data, err := os.ReadFile(path)
	if err != nil {
		return s, err
	}
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf")) // BOM added by Notepad and PowerShell
	return s, json.Unmarshal(data, &s)
}

// savedVars returns this PC's values of the pathVars and of any other
// variable a source is written with; see Settings.Vars.
func savedVars(s Settings) map[string]string {
	names := append([]string{}, pathVars...)
	for _, folders := range s.Games {
		for _, f := range folders {
			for _, m := range pathVar.FindAllStringSubmatch(f.Source, -1) {
				names = append(names, m[1])
			}
		}
	}
	vars := map[string]string{}
	for _, name := range names {
		if value, ok := os.LookupEnv(name); ok {
			vars[name] = value
		}
	}
	return vars
}

// save writes GameFix.json, and starts or stops watching sources to match.
// Games that were deleted from the main folder are forgotten.
func (a *App) save() error {
	for folder := range a.settings.Games {
		if _, err := os.Stat(folder); os.IsNotExist(err) && inside(folder, a.settings.MainFolder) {
			delete(a.settings.Games, folder)
		}
	}
	a.syncWatches()
	a.settings.Vars = savedVars(a.settings)
	data, err := json.MarshalIndent(a.settings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(a.settingsPath, data, 0o644)
}

// folders returns a copy of the folders of a game.
func (a *App) folders(gameFolder string) map[string]FolderSettings {
	folders := map[string]FolderSettings{}
	for name, f := range a.settings.Games[gameFolder] {
		folders[name] = f
	}
	return folders
}

// setFolders saves the folders of a game and rewrites GameFix.bat. Link
// paths that cleanLink refuses are cleared.
func (a *App) setFolders(gameFolder string, folders map[string]FolderSettings) error {
	for name, f := range folders {
		f.Link, _ = cleanLink(f.Link)
		folders[name] = f
	}
	if a.settings.Games == nil {
		a.settings.Games = map[string]map[string]FolderSettings{}
	}
	a.settings.Games[gameFolder] = folders
	if err := a.save(); err != nil {
		return err
	}
	return a.writeBat()
}

// importOldProfiles moves the gamefix.json files that older versions kept
// in each game folder into GameFix.json, and deletes them.
func (a *App) importOldProfiles() {
	entries, err := os.ReadDir(a.settings.MainFolder)
	if err != nil {
		return
	}
	found := false
	for _, e := range entries {
		old := filepath.Join(a.settings.MainFolder, e.Name(), "gamefix.json")
		data, err := os.ReadFile(old)
		if err != nil {
			continue
		}
		var p struct {
			Links map[string]FolderSettings `json:"links"`
		}
		if json.Unmarshal(data, &p) != nil {
			continue
		}
		for name, f := range p.Links {
			f.Link, _ = cleanLink(f.Link)
			f.Source, _ = checkSource(f.Source, a.settings.MainFolder)
			p.Links[name] = f
		}
		game := filepath.Join(a.settings.MainFolder, e.Name())
		if _, known := a.settings.Games[game]; !known {
			if a.settings.Games == nil {
				a.settings.Games = map[string]map[string]FolderSettings{}
			}
			a.settings.Games[game] = p.Links
		}
		os.Remove(old)
		found = true
	}
	if found {
		a.save()
		a.writeBat()
	}
}
