// Package hunt finds the folders an app uses (ported from GameFix's hunt
// mode). It starts the app's program, and while that runs it notes the
// folders of every program it starts and of the files they load, such as
// C:\Program Files\Roblox, and the folders where files are added, changed or
// deleted on the PC's drives, such as %LOCALAPPDATA%\Roblox. Changes are not
// tied to a program, so other programs' changes are found too; the ones that
// most likely are not the app are marked as noise.
package hunt

import (
	"errors"
	"fmt"
	"maps"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sys/windows"

	"raptech-portable/internal/core"
)

const errorElevationRequired = syscall.Errno(740) // ERROR_ELEVATION_REQUIRED

var (
	scanEvery = time.Second     // how often the programs are looked at
	grace     = 5 * time.Second // hunting stops this long after the last program ended
	waitRun   = 2 * time.Minute // how long to wait for a program that asks for administrator rights
)

// Folder is a folder the hunt suggests adding to the package.
type Folder struct {
	Path    string `json:"path"`    // with a variable where it can, such as %LOCALAPPDATA%\Roblox
	Real    string `json:"real"`    // the folder on this PC
	Program bool   `json:"program"` // a program of the app was in it, or loaded a file from it
	Changes int    `json:"changes"` // files added, changed or deleted in it while hunting
	Example string `json:"example"` // one of those programs or files
	Noise   bool   `json:"noise"`   // most likely Windows or another program, not the app
	Added   bool   `json:"added"`   // the package already has it
}

// State is what the hunt has found so far; it is sent while hunting and
// when it stops.
type State struct {
	ID       int64    `json:"id"` // when it started; each hunt has its own
	Exe      string   `json:"exe"`
	Running  bool     `json:"running"`
	Programs []string `json:"programs"` // the programs it started, and the ones they started
	Live     int      `json:"live"`     // how many of those still run
	Folders  []Folder `json:"folders"`  // program folders first, then by changes
	Error    string   `json:"error"`
}

// Hunter is one hunt.
type Hunter struct {
	exe      string
	args     string   // command line arguments, for tests
	roots    []string // see suggestFolder
	drives   []string // watched for changes
	exclude  []string // never suggested, such as the packages folder
	added    func() []string
	emit     func(State)
	start    time.Time
	stop     chan struct{}
	stopOnce sync.Once
	done     chan struct{}

	mu       sync.Mutex
	running  bool
	byPath   bool                        // started as administrator: its process is found by its path
	procs    map[uint32]windows.Filetime // the app's processes, with when each started
	live     int
	programs []string
	seen     map[string]bool
	finds    map[string]*find
	err      string
}

type find struct {
	real    string
	program bool
	changes int
	example string
}

// New makes a hunt of the program exe. exclude are folders never suggested
// (and everything below them); added returns the folders the package
// already has, shown as added; emit gets the state while hunting.
func New(exe string, exclude []string, added func() []string, emit func(State)) (*Hunter, error) {
	if !strings.EqualFold(filepath.Ext(exe), ".exe") {
		return nil, errors.New("choose a program (.exe)")
	}
	if added == nil {
		added = func() []string { return nil }
	}
	if emit == nil {
		emit = func(State) {}
	}
	return &Hunter{
		exe: filepath.Clean(exe), roots: defaultRoots(), drives: fixedDrives(), exclude: exclude,
		added: added, emit: emit, stop: make(chan struct{}), done: make(chan struct{}),
		procs: map[uint32]windows.Filetime{}, seen: map[string]bool{}, finds: map[string]*find{},
	}, nil
}

// Start starts the program and hunts in the background.
func (h *Hunter) Start() {
	h.mu.Lock()
	h.start, h.running = time.Now(), true
	h.mu.Unlock()
	go h.run()
}

// Stop stops hunting. The app keeps running.
func (h *Hunter) Stop() { h.stopOnce.Do(func() { close(h.stop) }) }

// Done is closed when the hunt has stopped.
func (h *Hunter) Done() <-chan struct{} { return h.done }

// Running reports whether it is still hunting.
func (h *Hunter) Running() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.running
}

// Exe is the program being hunted.
func (h *Hunter) Exe() string { return h.exe }

// Excluded are the folders never suggested.
func (h *Hunter) Excluded() []string { return slices.Clone(h.exclude) }

func (h *Hunter) run() {
	defer close(h.done)
	// Found paths are full names, so these must be too.
	for i := range h.roots {
		h.roots[i] = longPath(h.roots[i])
	}
	for i := range h.exclude {
		h.exclude[i] = longPath(h.exclude[i])
	}
	watching := make(chan struct{})
	defer close(watching)
	for _, d := range h.drives {
		go WatchDir(d, watching, func(names []string) {
			for _, name := range names {
				h.note(filepath.Join(d, name), false)
			}
		})
	}
	time.Sleep(300 * time.Millisecond) // so the drives are watched before the program writes
	if err := h.launch(); err != nil {
		h.finish(fmt.Sprintf("cannot start %s: %v", filepath.Base(h.exe), err))
		return
	}
	h.emit(h.State())
	last := time.Now()
	tick := time.NewTicker(scanEvery)
	defer tick.Stop()
	for {
		select {
		case <-h.stop:
			h.finish("")
			return
		case <-tick.C:
		}
		live := h.scan()
		h.mu.Lock()
		found := len(h.procs) > 0
		h.mu.Unlock()
		switch {
		case live > 0:
			last = time.Now()
		case found && time.Since(last) > grace:
			h.finish("")
			return
		case !found && time.Since(last) > waitRun:
			h.finish(filepath.Base(h.exe) + " did not start")
			return
		}
		h.emit(h.State())
	}
}

func (h *Hunter) finish(problem string) {
	h.mu.Lock()
	h.running, h.err = false, problem
	h.mu.Unlock()
	h.emit(h.State())
}

// launch starts the program from its own folder, as apps expect. One that
// asks for administrator rights is started the way Explorer starts it; its
// process is then found by its path.
func (h *Hunter) launch() error {
	cmd := exec.Command(h.exe)
	cmd.Dir = filepath.Dir(h.exe)
	if h.args != "" {
		cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: `"` + h.exe + `" ` + h.args}
	}
	err := cmd.Start()
	if err == nil {
		pid := uint32(cmd.Process.Pid)
		created, _ := processCreated(pid)
		h.mu.Lock()
		h.procs[pid] = created
		h.mu.Unlock()
		cmd.Process.Release()
		return nil
	}
	if !errors.Is(err, errorElevationRequired) {
		return err
	}
	h.mu.Lock()
	h.byPath = true
	h.mu.Unlock()
	verb, _ := windows.UTF16PtrFromString("runas")
	file, _ := windows.UTF16PtrFromString(h.exe)
	dir, _ := windows.UTF16PtrFromString(filepath.Dir(h.exe))
	var args *uint16
	if h.args != "" {
		args, _ = windows.UTF16PtrFromString(h.args)
	}
	return windows.ShellExecute(0, verb, file, args, dir, windows.SW_SHOWNORMAL)
}

// scan finds the processes the app's processes started since the last scan,
// notes what each running one runs and loads, and returns how many run.
func (h *Hunter) scan() int {
	list := processList()
	running := map[uint32]bool{}
	for _, p := range list {
		running[p.pid] = true
	}
	h.mu.Lock()
	known := maps.Clone(h.procs)
	byPath := h.byPath
	h.mu.Unlock()

	// A child may be listed before its parent, so look until none is added.
	for added := true; added; {
		added = false
		for _, p := range list {
			if _, ok := known[p.pid]; ok {
				continue
			}
			parentCreated, isChild := known[p.parent]
			if !isChild && !(byPath && strings.EqualFold(p.name, filepath.Base(h.exe))) {
				continue
			}
			created, ok := processCreated(p.pid)
			if !ok || created.Nanoseconds() < h.start.UnixNano() {
				continue
			}
			if isChild && running[p.parent] {
				// The parent's number may now be another program's.
				if c, ok := processCreated(p.parent); ok && c != parentCreated {
					continue
				}
			}
			if !isChild && !strings.EqualFold(processPath(p.pid), h.exe) {
				continue
			}
			known[p.pid] = created
			added = true
		}
	}

	live := 0
	for pid, created := range known {
		if !running[pid] {
			continue // ended; kept, so the processes it started are still found
		}
		if !h.inspect(pid, created) {
			delete(known, pid) // its number is now another program's
			continue
		}
		live++
	}
	h.mu.Lock()
	h.procs, h.live = known, live
	h.mu.Unlock()
	return live
}

// inspect notes the program a process runs and the files it has loaded. It
// returns false if the process is not the one that started at created.
func (h *Hunter) inspect(pid uint32, created windows.Filetime) bool {
	p, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION|windows.PROCESS_VM_READ, false, pid)
	if err != nil {
		p, err = windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	}
	if err != nil {
		return true // such as anti-cheat, which cannot be looked into: still the app's
	}
	defer windows.CloseHandle(p)
	var c, exit, kernel, user windows.Filetime
	if windows.GetProcessTimes(p, &c, &exit, &kernel, &user) == nil && c != created {
		return false
	}
	if exe := imagePath(p); exe != "" {
		h.noteProgram(exe, true)
	}
	for _, file := range moduleFiles(p) {
		h.noteProgram(file, false)
	}
	return true
}

func (h *Hunter) noteProgram(path string, exe bool) {
	path = strings.TrimPrefix(path, `\\?\`)
	h.mu.Lock()
	if h.seen[strings.ToLower(path)] {
		h.mu.Unlock()
		return
	}
	h.seen[strings.ToLower(path)] = true
	if exe {
		h.programs = append(h.programs, filepath.Base(path))
	}
	h.mu.Unlock()
	h.note(path, true)
}

// note counts path, a program file or a file that changed, in the folder
// suggested for it.
func (h *Hunter) note(path string, program bool) {
	path = longPath(path)
	for _, ex := range h.exclude {
		if inside(path, ex) {
			return
		}
	}
	folder := suggestFolder(path, h.roots)
	if folder == "" || strings.EqualFold(folder, path) && !isDir(folder) {
		return // such as a file right in a user's folder, like NTUSER.DAT
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	f := h.finds[strings.ToLower(folder)]
	if f == nil {
		f = &find{real: folder}
		h.finds[strings.ToLower(folder)] = f
	}
	if program {
		if !f.program {
			f.program, f.example = true, path
		}
	} else {
		f.changes++
		if f.example == "" {
			f.example = path
		}
	}
}

// State is what has been found so far.
func (h *Hunter) State() State {
	added := h.added()
	isAdded := func(real string) bool {
		for _, a := range added {
			if core.SamePath(a, real) {
				return true
			}
		}
		return false
	}
	exeFolder := suggestFolder(h.exe, h.roots)
	h.mu.Lock()
	defer h.mu.Unlock()
	st := State{
		ID: h.start.UnixMilli(), Exe: h.exe, Running: h.running,
		Programs: slices.Clone(h.programs), Live: h.live, Error: h.err, Folders: []Folder{},
	}
	if st.Programs == nil {
		st.Programs = []string{}
	}
	for _, f := range h.finds {
		st.Folders = append(st.Folders, Folder{
			Path: core.Portable(f.real), Real: f.real, Program: f.program, Changes: f.changes, Example: f.example,
			Noise: !strings.EqualFold(f.real, exeFolder) && (isNoise(f.real) || core.CheckLink(f.real) != nil),
			Added: isAdded(f.real),
		})
	}
	sort.Slice(st.Folders, func(i, j int) bool {
		a, b := st.Folders[i], st.Folders[j]
		if a.Program != b.Program {
			return a.Program
		}
		if a.Changes != b.Changes {
			return a.Changes > b.Changes
		}
		return strings.ToLower(a.Path) < strings.ToLower(b.Path)
	})
	return st
}
