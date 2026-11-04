// The launcher sits in a package folder next to portable.json, renamed after
// the app (such as Roblox.exe). It links the app's folders into the package,
// loads its registry, starts it, and saves everything back when it closes.
//
//	Roblox.exe [args]      set up, start the app with args, wait, save
//	Roblox.exe --setup     links, registry and protocols only (CCBoot logon)
//	Roblox.exe --update    update the package (run on the server)
//	Roblox.exe --save      save the registry into the package
//	Roblox.exe --unlink    remove the links, put the original folders back
//	Roblox.exe --status    print what is linked
//
// --pkg <folder> before the mode uses another package folder.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/sys/windows"

	"raptech-portable/internal/core"
	"raptech-portable/internal/roblox"
)

// LogName is the launcher's log inside the package.
const LogName = "RapPortable.log"

type options struct {
	mode string   // run, setup, update, save, unlink, status
	pkg  string   // package folder; "" = the launcher's folder
	pass []string // arguments for the app
}

func parseArgs(args []string) options {
	o := options{mode: "run"}
	if len(args) >= 2 && args[0] == "--pkg" {
		o.pkg, args = args[1], args[2:]
	}
	if len(args) > 0 {
		switch args[0] {
		case "--setup", "--update", "--save", "--unlink", "--status":
			o.mode, args = strings.TrimPrefix(args[0], "--"), args[1:]
		}
	}
	if len(args) > 0 {
		o.pass = args
	}
	return o
}

func main() {
	attachConsole()
	o := parseArgs(os.Args[1:])
	self, err := os.Executable()
	if err != nil {
		fail("", err)
	}
	if o.pkg == "" {
		o.pkg = filepath.Dir(self)
	}
	logf := openLog(o.pkg)
	p, err := core.Load(o.pkg)
	if err != nil {
		logf("cannot read the package: %v", err)
		fail(o.mode, err)
	}
	logf("%s %s %q", p.Name, o.mode, o.pass)

	switch o.mode {
	case "setup":
		err = setup(p, self, logf)
	case "run":
		err = run(p, self, o.pass, logf)
	case "update":
		err = update(p, self, logf)
	case "save":
		err = core.SaveRegistry(p.Registry, p.Path(core.RegistryDir), p.Dir)
	case "unlink":
		err = p.Unlink(self)
	case "status":
		for _, s := range p.Status() {
			fmt.Printf("%-12s %s -> %s\n", s.State, s.Link, s.Data)
		}
	}
	if err != nil {
		logf("%s failed: %v", o.mode, err)
		fail(o.mode, err)
	}
	logf("%s done", o.mode)
}

// setup links the package. When a link or key needs administrator rights
// (such as C:\Program Files or HKLM) it asks for them once and runs
// --setup again elevated.
func setup(p *core.Package, self string, logf func(string, ...any)) error {
	err := p.Setup(self)
	if err == nil {
		return nil
	}
	if errors.Is(err, windows.ERROR_ACCESS_DENIED) && !windows.GetCurrentProcessToken().IsElevated() {
		logf("needs administrator rights: %v", err)
		if eerr := runElevated(self, []string{"--pkg", p.Dir, "--setup"}); eerr != nil {
			return fmt.Errorf("%w (and elevating failed: %v)", err, eerr)
		}
		return p.Setup(self) // check again; now only real problems remain
	}
	return err
}

func run(p *core.Package, self string, pass []string, logf func(string, ...any)) error {
	if err := setup(p, self, logf); err != nil {
		logf("setup problems (starting anyway): %v", err)
	}
	if p.Exe == "" {
		return errors.New("portable.json has no exe")
	}
	exe, err := p.ResolveExe(p.Exe)
	if err != nil {
		return err
	}
	// Through the app's default location, so it never sees the package.
	exe = p.StartPath(exe)
	logf("start %s", exe)
	res, err := core.Run(exe, append(append([]string{}, p.Args...), pass...), "", 0)
	if err != nil {
		return err
	}
	logf("closed (exit %d)", res.ExitCode)
	if n := otherInstances(self); n > 0 {
		logf("%d other copies still running; saving later", n)
		return nil
	}
	return p.Finish()
}

func update(p *core.Package, self string, logf func(string, ...any)) error {
	u := p.Update
	if u == nil {
		return errors.New("portable.json has no update")
	}
	switch u.Kind {
	case "roblox":
		ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
		defer cancel()
		v, installed, err := roblox.New(logf).Update(ctx, p.Path(u.Folder), u.Keep)
		if err == nil {
			logf("Roblox %s (installed now: %v)", v, installed)
		}
		return err
	case "run":
		timeout := 30 * time.Minute
		if u.Timeout != "" {
			timeout, _ = time.ParseDuration(u.Timeout)
		}
		if err := setup(p, self, logf); err != nil {
			logf("setup problems (updating anyway): %v", err)
		}
		res, err := core.Run(p.Path(u.Exe), u.Args, "", timeout)
		if err == nil {
			logf("update program ended (exit %d, timed out %v)", res.ExitCode, res.TimedOut)
			err = core.SaveRegistry(p.Registry, p.Path(core.RegistryDir), p.Dir)
		}
		return err
	}
	return fmt.Errorf("unknown update kind %q", u.Kind)
}

// openLog appends to the package's log, or %TEMP% when the package folder
// cannot be written. Over 1 MB it starts again, keeping one old log.
func openLog(dir string) func(string, ...any) {
	path := filepath.Join(dir, LogName)
	if info, err := os.Stat(path); err == nil && info.Size() > 1<<20 {
		os.Rename(path, path+".old")
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		f, err = os.OpenFile(filepath.Join(os.TempDir(), LogName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	}
	host, _ := os.Hostname()
	return func(format string, a ...any) {
		line := fmt.Sprintf("%s %s [%d] %s\r\n", time.Now().Format("2006-01-02 15:04:05"), host, os.Getpid(), fmt.Sprintf(format, a...))
		if err == nil {
			f.WriteString(line)
		}
		os.Stdout.WriteString(line)
	}
}

// fail shows the problem: a message box when the app was started normally
// (no console), plus the exit code for scripts.
func fail(mode string, err error) {
	if mode == "run" || mode == "" {
		title, _ := windows.UTF16PtrFromString("RapTech Portable")
		text, _ := windows.UTF16PtrFromString(err.Error())
		windows.MessageBox(0, text, title, windows.MB_ICONERROR|windows.MB_OK)
	} else {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(1)
}
