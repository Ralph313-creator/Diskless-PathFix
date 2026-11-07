package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"

	"raptech-portable/internal/selfupdate"
)

// Version is this build's version. Keep it the same as wails.json's
// info.productVersion; a release's tag is v + Version.
const Version = "1.0.0"

// UpdateInfo is what the window shows about the newest release.
type UpdateInfo struct {
	Current string
	Latest  string
	Newer   bool
	Notes   string
	Page    string
}

var httpClient = &http.Client{Timeout: 5 * time.Minute}

// GetVersion is this build's version.
func (a *App) GetVersion() string { return Version }

// CheckUpdate asks GitHub for the newest release.
func (a *App) CheckUpdate() (UpdateInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	r, err := selfupdate.Latest(ctx, httpClient, selfupdate.LatestURL)
	if err != nil {
		return UpdateInfo{Current: Version}, err
	}
	return UpdateInfo{
		Current: Version,
		Latest:  strings.TrimPrefix(r.Version, "v"),
		Newer:   selfupdate.Newer(r.Version, Version),
		Notes:   r.Notes,
		Page:    r.Page,
	}, nil
}

// InstallUpdate downloads the newest release, puts it in place of this
// exe, restarts the GameFix service if it runs from this exe, then starts
// the new exe and closes this one. Settings next to the exe are kept.
func (a *App) InstallUpdate() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	if strings.HasSuffix(strings.ToLower(exe), "-dev.exe") {
		return errors.New("wails dev cannot update itself; use the built exe")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	r, err := selfupdate.Latest(ctx, httpClient, selfupdate.LatestURL)
	if err != nil {
		return err
	}
	if !selfupdate.Newer(r.Version, Version) {
		return errors.New("this is already the newest version")
	}
	if err := selfupdate.Download(ctx, httpClient, r, exe+".new"); err != nil {
		return err
	}
	if err := selfupdate.Swap(exe, exe+".new"); err != nil {
		return err
	}
	restartServiceOf(exe)
	if err := exec.Command(exe).Start(); err != nil {
		return errors.New("updated, but the new version did not start: open it again by hand")
	}
	if a.ctx != nil {
		runtime.Quit(a.ctx)
	}
	return nil
}

// restartServiceOf restarts the GameFix service when it runs exe, so it
// runs the new file. Anything else is left alone.
func restartServiceOf(exe string) {
	m, err := mgr.Connect()
	if err != nil {
		return
	}
	defer m.Disconnect()
	s, err := m.OpenService("RapTechGameFixWatch")
	if err != nil {
		return
	}
	defer s.Close()
	cfg, err := s.Config()
	if err != nil || !strings.Contains(strings.ToLower(cfg.BinaryPathName), strings.ToLower(exe)) {
		return
	}
	st, err := s.Query()
	if err != nil || st.State != svc.Running {
		return
	}
	s.Control(svc.Stop)
	for i := 0; i < 50; i++ {
		if st, err = s.Query(); err != nil || st.State == svc.Stopped {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	s.Start()
}
