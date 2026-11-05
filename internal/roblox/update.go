// Package roblox downloads Roblox Player versions straight from Roblox into
// a Versions folder, the way the official installer lays them out.
package roblox

import (
	"archive/zip"
	"bufio"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"golang.org/x/sys/windows"

	"raptech-portable/internal/core"
)

const (
	DefaultVersionURL = "https://clientsettings.roblox.com/v2/client-version/WindowsPlayer"
	DefaultSetupURL   = "https://setup.rbxcdn.com"
	PlayerExe         = "RobloxPlayerBeta.exe"
)

// packageDirs is where each Roblox package zip is extracted inside a
// version folder.
var packageDirs = map[string]string{
	"RobloxApp.zip":                     ``,
	"redist.zip":                        ``,
	"WebView2.zip":                      ``,
	"shaders.zip":                       `shaders`,
	"ssl.zip":                           `ssl`,
	"WebView2RuntimeInstaller.zip":      `WebView2RuntimeInstaller`,
	"content-avatar.zip":                `content\avatar`,
	"content-configs.zip":               `content\configs`,
	"content-fonts.zip":                 `content\fonts`,
	"content-sky.zip":                   `content\sky`,
	"content-sounds.zip":                `content\sounds`,
	"content-textures2.zip":             `content\textures`,
	"content-models.zip":                `content\models`,
	"content-platform-fonts.zip":        `PlatformContent\pc\fonts`,
	"content-platform-dictionaries.zip": `PlatformContent\pc\shared_compression_dictionaries`,
	"content-terrain.zip":               `PlatformContent\pc\terrain`,
	"content-textures3.zip":             `PlatformContent\pc\textures`,
	"extracontent-luapackages.zip":      `ExtraContent\LuaPackages`,
	"extracontent-translations.zip":     `ExtraContent\translations`,
	"extracontent-models.zip":           `ExtraContent\models`,
	"extracontent-textures.zip":         `ExtraContent\textures`,
	"extracontent-places.zip":           `ExtraContent\places`,
}

const appSettings = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\r\n<Settings>\r\n\t<ContentFolder>content</ContentFolder>\r\n\t<BaseUrl>http://www.roblox.com</BaseUrl>\r\n</Settings>\r\n"

// Updater keeps a Versions folder on the newest Roblox Player.
type Updater struct {
	Client     *http.Client
	VersionURL string
	SetupURL   string
	Log        func(format string, a ...any)
}

// New is an Updater for the real Roblox servers.
func New(log func(string, ...any)) *Updater {
	return &Updater{Client: &http.Client{Timeout: 30 * time.Minute}, VersionURL: DefaultVersionURL, SetupURL: DefaultSetupURL, Log: log}
}

func (u *Updater) logf(format string, a ...any) {
	if u.Log != nil {
		u.Log(format, a...)
	}
}

type pkgEntry struct {
	name string
	md5  string
}

// Update installs the newest version into dir when it is missing, then
// prunes old versions down to keep. It returns the newest version and
// whether it was installed now.
func (u *Updater) Update(ctx context.Context, dir string, keep int) (string, bool, error) {
	version, err := u.Latest(ctx)
	if err != nil {
		return "", false, err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return version, false, err
	}
	cleanLeftovers(dir, u.logf)
	if complete(filepath.Join(dir, version)) {
		Prune(dir, keep, version, u.logf)
		return version, false, nil
	}
	unlock, err := lock(dir)
	if err != nil {
		return version, false, err
	}
	defer unlock()

	u.logf("installing %s", version)
	if err := u.install(ctx, dir, version); err != nil {
		return version, false, err
	}
	u.logf("installed %s", version)
	Prune(dir, keep, version, u.logf)
	return version, true, nil
}

// Latest asks Roblox for the newest Windows Player version, such as
// version-5a21b87bbc2641cf.
func (u *Updater) Latest(ctx context.Context) (string, error) {
	body, err := u.get(ctx, u.VersionURL)
	if err != nil {
		return "", fmt.Errorf("ask Roblox for the version: %w", err)
	}
	defer body.Close()
	var v struct {
		ClientVersionUpload string `json:"clientVersionUpload"`
	}
	if err := json.NewDecoder(body).Decode(&v); err != nil {
		return "", fmt.Errorf("read Roblox version: %w", err)
	}
	if !strings.HasPrefix(v.ClientVersionUpload, "version-") || strings.ContainsAny(v.ClientVersionUpload, `\/.:`) {
		return "", fmt.Errorf("Roblox sent an odd version %q", v.ClientVersionUpload)
	}
	return v.ClientVersionUpload, nil
}

func (u *Updater) get(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := u.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("%s: %s", url, resp.Status)
	}
	return resp.Body, nil
}

func (u *Updater) install(ctx context.Context, dir, version string) (err error) {
	entries, err := u.manifest(ctx, version)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if _, known := packageDirs[e.name]; !known && strings.HasSuffix(strings.ToLower(e.name), ".zip") {
			return fmt.Errorf("Roblox added a package this updater does not know (%s); update Diskless PathFix", e.name)
		}
		if strings.ContainsAny(e.name, `\/:`) {
			return fmt.Errorf("odd package name %q", e.name)
		}
	}

	tmp := filepath.Join(dir, ".rp-tmp-"+version)
	os.RemoveAll(tmp)
	defer func() {
		if err != nil {
			os.RemoveAll(tmp)
		}
	}()
	dl := filepath.Join(tmp, ".download")
	if err := os.MkdirAll(dl, 0o755); err != nil {
		return err
	}
	for _, e := range entries {
		file := filepath.Join(dl, e.name)
		if err := u.download(ctx, version, e, file); err != nil {
			return err
		}
		sub, isZip := packageDirs[e.name]
		if !isZip {
			if err := os.Rename(file, filepath.Join(tmp, e.name)); err != nil {
				return err
			}
			continue
		}
		if err := extract(file, filepath.Join(tmp, sub)); err != nil {
			return fmt.Errorf("%s: %w", e.name, err)
		}
		os.Remove(file)
	}
	os.RemoveAll(dl)
	if _, err := os.Stat(filepath.Join(tmp, PlayerExe)); err != nil {
		return fmt.Errorf("%s has no %s", version, PlayerExe)
	}
	if err := os.WriteFile(filepath.Join(tmp, "AppSettings.xml"), []byte(appSettings), 0o644); err != nil {
		return err
	}
	final := filepath.Join(dir, version)
	if _, err := os.Stat(final); err == nil {
		// An unfinished folder of this version is in the way.
		old := filepath.Join(dir, ".rp-old-"+version)
		if err := os.Rename(final, old); err != nil {
			return fmt.Errorf("replace unfinished %s (in use?): %w", version, err)
		}
		os.RemoveAll(old)
	}
	return os.Rename(tmp, final)
}

func (u *Updater) manifest(ctx context.Context, version string) ([]pkgEntry, error) {
	body, err := u.get(ctx, u.SetupURL+"/"+version+"-rbxPkgManifest.txt")
	if err != nil {
		return nil, fmt.Errorf("package list: %w", err)
	}
	defer body.Close()
	var lines []string
	sc := bufio.NewScanner(body)
	for sc.Scan() {
		if line := strings.TrimSpace(sc.Text()); line != "" {
			lines = append(lines, line)
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	if len(lines) == 0 || lines[0] != "v0" || (len(lines)-1)%4 != 0 {
		return nil, errors.New("the package list is not in the v0 format")
	}
	var out []pkgEntry
	for i := 1; i < len(lines); i += 4 {
		out = append(out, pkgEntry{name: lines[i], md5: strings.ToLower(lines[i+1])})
	}
	return out, nil
}

// download saves one package to file and checks its MD5, trying 3 times.
func (u *Updater) download(ctx context.Context, version string, e pkgEntry, file string) error {
	var last error
	for try := 1; try <= 3; try++ {
		if try > 1 {
			u.logf("retry %s (%v)", e.name, last)
		}
		last = u.downloadOnce(ctx, version, e, file)
		if last == nil || ctx.Err() != nil {
			return last
		}
	}
	return last
}

func (u *Updater) downloadOnce(ctx context.Context, version string, e pkgEntry, file string) error {
	body, err := u.get(ctx, u.SetupURL+"/"+version+"-"+e.name)
	if err != nil {
		return err
	}
	defer body.Close()
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	h := md5.New()
	_, err = io.Copy(io.MultiWriter(f, h), body)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		return fmt.Errorf("download %s: %w", e.name, err)
	}
	if got := hex.EncodeToString(h.Sum(nil)); got != e.md5 {
		return fmt.Errorf("%s is damaged (MD5 %s, want %s)", e.name, got, e.md5)
	}
	return nil
}

// extract unzips file into dest. Roblox zips use backslashes and list
// folders as "\Name\"; entries that would land outside dest are refused.
func extract(file, dest string) error {
	r, err := zip.OpenReader(file)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		name := strings.ReplaceAll(f.Name, "/", `\`)
		if strings.HasSuffix(name, `\`) {
			continue // folder entry; folders are made for files below
		}
		name = strings.TrimLeft(name, `\`)
		if err := core.CheckData(name); err != nil {
			return fmt.Errorf("unsafe path %q in zip", f.Name)
		}
		out := filepath.Join(dest, name)
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return err
		}
		if err := writeEntry(f, out); err != nil {
			return err
		}
	}
	return nil
}

func writeEntry(f *zip.File, out string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	w, err := os.Create(out)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, rc)
	if cerr := w.Close(); err == nil {
		err = cerr
	}
	return err
}

// complete reports whether a version folder has the player and settings.
func complete(dir string) bool {
	for _, f := range []string{PlayerExe, "AppSettings.xml"} {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			return false
		}
	}
	return true
}

// Prune deletes version-* folders beyond the keep newest (by creation
// time), never current. A folder in use cannot be renamed, so it is left
// for next time instead of being half deleted.
func Prune(dir string, keep int, current string, log func(string, ...any)) {
	if keep < 1 {
		keep = 2
	}
	entries, _ := os.ReadDir(dir)
	var versions []string
	for _, e := range entries {
		if e.IsDir() && strings.HasPrefix(e.Name(), "version-") {
			versions = append(versions, e.Name())
		}
	}
	sort.Slice(versions, func(i, j int) bool {
		return core.Created(filepath.Join(dir, versions[i])).After(core.Created(filepath.Join(dir, versions[j])))
	})
	kept := 0
	if current != "" {
		kept = 1 // current always stays and counts
	}
	for _, v := range versions {
		if v == current {
			continue
		}
		if kept < keep {
			kept++
			continue
		}
		old := filepath.Join(dir, ".rp-old-"+v)
		if err := os.Rename(filepath.Join(dir, v), old); err != nil {
			if log != nil {
				log("keep %s for now (in use): %v", v, err)
			}
			continue
		}
		os.RemoveAll(old)
		if log != nil {
			log("removed old %s", v)
		}
	}
}

// cleanLeftovers removes what an interrupted update or prune left behind.
func cleanLeftovers(dir string, log func(string, ...any)) {
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".rp-tmp-") || strings.HasPrefix(e.Name(), ".rp-old-") {
			if err := os.RemoveAll(filepath.Join(dir, e.Name())); err == nil {
				log("cleaned %s", e.Name())
			}
		}
	}
}

// lock makes sure only one update runs on dir. The lock file is deleted when
// the handle closes, also when the process dies.
func lock(dir string) (func(), error) {
	name, err := windows.UTF16PtrFromString(filepath.Join(dir, ".rp-update.lock"))
	if err != nil {
		return nil, err
	}
	h, err := windows.CreateFile(name, windows.GENERIC_WRITE, 0, nil, windows.CREATE_ALWAYS,
		windows.FILE_ATTRIBUTE_HIDDEN|windows.FILE_FLAG_DELETE_ON_CLOSE, 0)
	if err != nil {
		return nil, errors.New("another update of this folder is running")
	}
	return func() { windows.CloseHandle(h) }, nil
}
