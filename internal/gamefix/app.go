package gamefix

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Game is one game folder inside the main folder: <main folder>\<game name>.
type Game struct {
	Name    string   `json:"name"`
	Folder  string   `json:"folder"`
	Folders []Folder `json:"folders"` // folders copied into the game folder
}

// Folder is one folder of a game. Copy copies Source to <game folder>\<Name>,
// and GameFix.bat links Link to that copy (GameFix Copy.bat copies it there instead).
type Folder struct {
	Name   string `json:"name"`   // GameItself for the game folder itself
	Source string `json:"source"` // where it is copied from, e.g. C:\Users\me\Documents\Rockstar Games
	Link   string `json:"link"`   // where the symbolic link is made
	Copied bool   `json:"copied"` // whether the copy exists in the main folder
	Watch  bool   `json:"watch"`  // on (Enable): clients get it, and it is copied again whenever the source changes
}

// Progress is sent to the UI while a folder is being copied.
type Progress struct {
	Game   string `json:"game"`
	Folder string `json:"folder"`
	File   string `json:"file"`
	Done   int64  `json:"done"`
	Total  int64  `json:"total"`
}

// App struct
type App struct {
	ctx          context.Context
	settings     Settings
	settingsPath string // GameFix.json next to GameFix.exe

	// Watched folders are copied in the background, so the settings are
	// only used with mu held. copyMu is held while a folder is copied or
	// deleted, so only one copy runs at a time. When both are needed,
	// copyMu is locked first.
	mu      sync.Mutex
	copyMu  sync.Mutex
	watches map[string]*watch // see watchKey

	huntMu sync.Mutex

	service bool        // running as the Windows service, not the app
	log     *log.Logger // where the service writes what happened; nil in the app
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{settingsPath: exeSettingsPath()}
}

// Startup is called when the window starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.mu.Lock()
	defer a.mu.Unlock()
	a.load()
	// The app watches only while the service is not running, so watching
	// moves between them when the service is started or stopped.
	go func() {
		for range time.Tick(5 * time.Second) {
			a.mu.Lock()
			a.syncWatches()
			a.mu.Unlock()
		}
	}()
}

// GetMainFolder returns the folder that games are copied into.
func (a *App) GetMainFolder() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.settings.MainFolder
}

// ChooseMainFolder asks for the main folder and remembers it.
// It returns "" when the dialog is cancelled.
func (a *App) ChooseMainFolder() (string, error) {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "Choose the main folder",
		DefaultDirectory: a.GetMainFolder(),
	})
	if err != nil || path == "" {
		return "", err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.settings.MainFolder = filepath.Clean(path)
	if err := a.save(); err != nil {
		return "", err
	}
	a.importOldProfiles()
	a.writeBat() // also replaces the per-game scripts of older versions
	return a.settings.MainFolder, nil
}

// ListGames returns every game folder in the main folder.
func (a *App) ListGames() ([]Game, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	games := []Game{}
	if a.settings.MainFolder == "" {
		return games, nil
	}
	entries, err := os.ReadDir(a.settings.MainFolder)
	if err != nil {
		return games, err
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		folder := filepath.Join(a.settings.MainFolder, e.Name())
		games = append(games, Game{Name: e.Name(), Folder: folder, Folders: gameFolders(folder, a.folders(folder))})
	}
	sort.Slice(games, func(i, j int) bool { return strings.ToLower(games[i].Name) < strings.ToLower(games[j].Name) })
	return games, nil
}

// gameFolder is <main folder>\<name>, after checking both are usable.
func (a *App) gameFolder(name string) (string, error) {
	if a.settings.MainFolder == "" {
		return "", errors.New("choose the main folder first")
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return "", errors.New("enter a game name")
	}
	if strings.ContainsAny(name, `<>:"/\|?*`) || strings.HasSuffix(name, ".") || name == ".." {
		return "", errors.New(`a game name cannot contain < > : " / \ | ? * or end with a dot`)
	}
	return filepath.Join(a.settings.MainFolder, name), nil
}

// AddGame creates the game's own folder, <main folder>\<name>.
func (a *App) AddGame(name string) (*Game, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	folder, err := a.gameFolder(name)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(folder); err == nil {
		return nil, errors.New(`"` + filepath.Base(folder) + `" already exists`)
	}
	if err := os.MkdirAll(folder, 0o755); err != nil {
		return nil, err
	}
	return &Game{Name: filepath.Base(folder), Folder: folder, Folders: []Folder{}}, nil
}

// DeleteGame deletes <main folder>\<name> with everything in it.
func (a *App) DeleteGame(name string) error {
	a.copyMu.Lock()
	defer a.copyMu.Unlock()
	a.mu.Lock()
	defer a.mu.Unlock()
	folder, err := a.gameFolder(name)
	if err != nil {
		return err
	}
	if st, err := os.Stat(folder); err != nil || !st.IsDir() {
		return errors.New(`"` + name + `" not found`)
	}
	for f := range a.folders(folder) {
		a.stopWatch(folder, f)
	}
	if err := removeTree(folder); err != nil {
		return err
	}
	if err := a.save(); err != nil { // forgets the game's folders
		return err
	}
	return a.writeBat()
}

// removeTree deletes a folder with everything in it. Copied files keep
// their read-only attribute, which stops Windows deleting them, so that is
// cleared first.
func removeTree(folder string) error {
	filepath.WalkDir(folder, func(path string, d os.DirEntry, err error) error {
		if err == nil && d.Type().IsRegular() {
			os.Chmod(path, 0o666)
		}
		return nil
	})
	return os.RemoveAll(folder)
}

// ChooseFolder asks for a folder to add to a game.
// It returns "" when the dialog is cancelled.
func (a *App) ChooseFolder(game string) (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{Title: "Choose a folder to add to " + game})
}

// BrowseFolder asks for a folder, starting at start or the nearest folder
// above it that exists. It returns "" when the dialog is cancelled.
func (a *App) BrowseFolder(title, start string) (string, error) {
	dir := expandVars(start)
	if !filepath.IsAbs(dir) {
		dir = ""
	}
	for dir != "" {
		if st, err := os.Stat(dir); err == nil && st.IsDir() {
			break
		}
		if parent := filepath.Dir(dir); parent != dir {
			dir = parent
		} else {
			dir = ""
		}
	}
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{Title: title, DefaultDirectory: dir})
}

// SaveFolders saves the source and link path of each folder in a game and
// rewrites <main folder>\GameFix.bat. It returns the path of GameFix.bat.
func (a *App) SaveFolders(game string, changes map[string]FolderSettings) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	folder, err := a.gameFolder(game)
	if err != nil {
		return "", err
	}
	folders := a.folders(folder)
	for name, f := range changes {
		old, ok := folders[name]
		if !ok {
			continue
		}
		f.Watch = old.Watch // changed with SetWatch
		if f.Source, err = checkSource(f.Source, a.settings.MainFolder); err != nil {
			return "", errors.New(displayName(game, name) + ": " + err.Error())
		}
		if f.Link, err = cleanLink(f.Link); err != nil {
			return "", errors.New(displayName(game, name) + ": " + err.Error())
		}
		folders[name] = f
	}
	return filepath.Join(a.settings.MainFolder, BatName), a.setFolders(folder, folders)
}

// displayName is how a folder of a game is shown to the user.
func displayName(game, name string) string {
	if name == GameItself {
		return game
	}
	return name
}

// checkSource checks where a folder is copied from. "" means not set yet.
// It is written with a variable where it can, like a link path (see
// portable); Copy expands it again on this PC.
func checkSource(source, mainFolder string) (string, error) {
	source = strings.TrimSpace(source)
	if source == "" {
		return "", nil
	}
	real := filepath.Clean(expandVars(source))
	if !filepath.IsAbs(real) {
		return "", errors.New(`the source must be a full path, such as C:\Games\GTA V`)
	}
	if strings.HasSuffix(filepath.Base(real), `\`) {
		return "", errors.New("the source must be a folder, not a whole drive")
	}
	if inside(real, mainFolder) || inside(mainFolder, real) {
		return "", errors.New("the source and the main folder must not be inside each other")
	}
	return portable(real), nil
}

// AddFolder adds a folder to a game. Copy copies it from from into
// <main folder>\<game>\<folder name>, and its link path starts as from.
// A folder named like the game goes straight into the game folder:
// <main folder>\GTA V, not <main folder>\GTA V\GTA V.
func (a *App) AddFolder(game, from string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	folder, err := a.gameFolder(game)
	if err != nil {
		return err
	}
	source, err := checkSource(from, a.settings.MainFolder)
	if err != nil || source == "" {
		return err
	}
	name := filepath.Base(expandVars(source))
	if strings.EqualFold(name, filepath.Base(folder)) {
		name = GameItself
	}
	folders := a.folders(folder)
	for existing := range folders {
		if strings.EqualFold(existing, name) {
			return errors.New(`"` + displayName(game, name) + `" is already in ` + game)
		}
	}
	link, _ := cleanLink(source) // "" when the source is a folder that cannot be a link path
	folders[name] = FolderSettings{Source: source, Link: link}
	return a.setFolders(folder, folders)
}

// CopyFolder copies a folder of a game from its source, with every file and
// subfolder, into the main folder, and returns where it was copied to. The
// copy is made the same as the source: copying again updates changed files
// and deletes the ones that are no longer in the source. A watched folder is
// copied with it too, whenever its source changes.
func (a *App) CopyFolder(game, name string) (string, error) {
	a.copyMu.Lock()
	defer a.copyMu.Unlock()
	source, dest, keep, err := a.copySource(game, name)
	if err != nil {
		return "", err
	}
	total := folderSize(source)
	err = robocopy(source, dest, keep, func(file string, done int64) {
		if a.ctx != nil { // nil in tests
			runtime.EventsEmit(a.ctx, "copy-progress", Progress{Game: game, Folder: name, File: file, Done: done, Total: total})
		}
	})
	if err != nil {
		return dest, err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return dest, a.writeBat() // the folder is linked once it is copied
}

// copySource returns where a folder of a game is copied from, with its
// variables expanded, and where to. keep lists the folders inside dest that
// are not part of the source and must not be deleted: for the game folder
// itself, the copies of the game's other folders.
func (a *App) copySource(game, name string) (source, dest string, keep []string, err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	folder, err := a.gameFolder(game)
	if err != nil {
		return "", "", nil, err
	}
	folders := a.folders(folder)
	f, ok := folders[name]
	if !ok {
		return "", "", nil, errors.New(`"` + displayName(game, name) + `" is not in ` + game)
	}
	if f.Source == "" {
		return "", "", nil, errors.New(displayName(game, name) + ": set the source first")
	}
	source = expandVars(f.Source)
	if !isDir(source) {
		return "", "", nil, errors.New(displayName(game, name) + ": source folder not found: " + source)
	}
	if name == GameItself {
		for other := range folders {
			if other != GameItself {
				keep = append(keep, copyPath(folder, other))
			}
		}
	}
	return source, copyPath(folder, name), keep, nil
}

// RemoveFolder removes a folder from a game and deletes its copy. For the
// game folder itself, the copied files are kept; delete the game for that.
func (a *App) RemoveFolder(game, name string) error {
	a.copyMu.Lock()
	defer a.copyMu.Unlock()
	a.mu.Lock()
	defer a.mu.Unlock()
	folder, err := a.gameFolder(game)
	if err != nil {
		return err
	}
	folders := a.folders(folder)
	if _, ok := folders[name]; !ok {
		return errors.New(`"` + displayName(game, name) + `" is not in ` + game)
	}
	a.stopWatch(folder, name)
	if name == GameItself {
		if err := removeGameFiles(folder, folders); err != nil {
			return err
		}
	} else if err := removeTree(copyPath(folder, name)); err != nil {
		return err
	}
	delete(folders, name)
	return a.setFolders(folder, folders)
}

// removeGameFiles deletes the copy of the game folder itself: everything in
// gameFolder except the copies of the game's other folders. The game folder
// stays, so the game stays in the list.
func removeGameFiles(gameFolder string, folders map[string]FolderSettings) error {
	entries, err := os.ReadDir(gameFolder)
	if err != nil {
		return err
	}
	var failed []string
entries:
	for _, e := range entries {
		for name := range folders {
			if name != GameItself && strings.EqualFold(name, e.Name()) {
				continue entries // another folder's copy
			}
		}
		if err := removeEntry(filepath.Join(gameFolder, e.Name())); err != nil {
			failed = append(failed, err.Error())
		}
	}
	if len(failed) > 0 {
		if len(failed) > 3 {
			failed = append(failed[:3], fmt.Sprintf("and %d more", len(failed)-3))
		}
		return errors.New(strings.Join(failed, "\n") + "\nIs the game running?")
	}
	return nil
}

// removeEntry deletes a file, or a folder with everything in it. For a
// link or junction, only the link is deleted, never what it points to.
func removeEntry(path string) error {
	st, err := os.Lstat(path)
	if err != nil {
		return nil
	}
	switch {
	case st.Mode()&(os.ModeSymlink|os.ModeIrregular) != 0:
		err = os.Remove(path)
	case st.IsDir():
		err = removeTree(path)
	default:
		os.Chmod(path, 0o666) // a read-only file cannot be deleted
		err = os.Remove(path)
	}
	if err != nil {
		return fmt.Errorf("cannot delete %s: %w", path, err)
	}
	return nil
}

// key is how paths are compared: Windows names ignore case.
func key(s string) string { return strings.ToLower(s) }

// inside reports whether path is dir or somewhere below it.
func inside(path, dir string) bool {
	rel, err := filepath.Rel(dir, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, `..\`)
}

// folderSize adds up the size of every file below root.
func folderSize(root string) int64 {
	var size int64
	filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			if info, err := d.Info(); err == nil {
				size += info.Size()
			}
		}
		return nil
	})
	return size
}

// robocopy makes to the same as from with Windows' own tool: it keeps
// timestamps and attributes, copies hidden/system/read-only files, handles
// long paths and retries files that are briefly locked. Files and folders
// in to that are not in from are deleted, except the folders in keep (full
// paths inside to). progress is called for each file, at most a few times a
// second.
func robocopy(from, to string, keep []string, progress func(file string, done int64)) error {
	// /MIR: all subfolders, including empty ones, and delete what is not in
	// from. /XD: not copied or deleted. /XJ: skip junctions, which can loop.
	args := []string{from, to, "/MIR", "/COPY:DAT", "/DCOPY:T", "/XJ",
		"/R:2", "/W:1", "/NP", "/NDL", "/NJH", "/NJS", "/BYTES", "/FP"}
	if len(keep) > 0 {
		args = append(append(args, "/XD"), keep...)
	}
	cmd := exec.Command("robocopy", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000} // CREATE_NO_WINDOW
	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	var done int64
	var failed []string
	last := time.Time{}
	scan := bufio.NewScanner(out)
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
		if strings.Contains(line, "ERROR") {
			failed = append(failed, line)
			continue
		}
		if strings.HasPrefix(line, "*EXTRA") { // already in the main folder, not copied
			continue
		}
		// File lines look like: "New File <tab> 12345 <tab> C:\game\file.dat"
		fields := strings.Split(line, "\t")
		if len(fields) < 2 {
			continue
		}
		file := strings.TrimSpace(fields[len(fields)-1])
		if size, err := strconv.ParseInt(strings.TrimSpace(fields[len(fields)-2]), 10, 64); err == nil {
			done += size
		}
		if time.Since(last) > 200*time.Millisecond {
			progress(file, done)
			last = time.Now()
		}
	}
	err = cmd.Wait()
	progress("", done)
	var exit *exec.ExitError
	if errors.As(err, &exit) && exit.ExitCode() < 8 {
		return nil // 1-7 mean files were copied / extra files exist, not failure
	}
	if err != nil {
		if len(failed) > 3 {
			failed = append(failed[:3], "...")
		}
		if len(failed) == 0 {
			return errors.New("robocopy failed: " + err.Error())
		}
		return errors.New("some files could not be copied:\n" + strings.Join(failed, "\n"))
	}
	return nil
}

// OpenFolder shows a folder in Explorer.
func (a *App) OpenFolder(path string) error {
	if path == "" {
		return errors.New("no folder")
	}
	return exec.Command("explorer", path).Start()
}
