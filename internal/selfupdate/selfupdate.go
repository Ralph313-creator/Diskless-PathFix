// Package selfupdate updates Diskless PathFix from its GitHub releases: it
// finds the newest release, downloads its exe next to the running one,
// checks its SHA-256 against the digest GitHub gives, and swaps the files.
// A running exe can be renamed but not replaced, so the old one becomes
// <exe>.old and is deleted by the next start.
package selfupdate

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	// Repo is where the releases are.
	Repo = "Ralph313-creator/Diskless-PathFix"
	// AssetName is the exe attached to each release.
	AssetName = "Diskless-PathFix.exe"
	// LatestURL is GitHub's newest-release API for Repo.
	LatestURL = "https://api.github.com/repos/" + Repo + "/releases/latest"
)

// Release is the newest release and its exe.
type Release struct {
	Version string // tag, such as v1.2.0
	URL     string // the exe's download link
	SHA256  string // hex, from GitHub's asset digest
	Notes   string
	Page    string // the release's web page
}

// Newer reports whether version latest is higher than current. Both may
// start with v; anything that is not dotted numbers is never newer.
func Newer(latest, current string) bool {
	l, ok1 := parse(latest)
	c, ok2 := parse(current)
	if !ok1 || !ok2 {
		return false
	}
	for i := 0; i < len(l) || i < len(c); i++ {
		var a, b int
		if i < len(l) {
			a = l[i]
		}
		if i < len(c) {
			b = c[i]
		}
		if a != b {
			return a > b
		}
	}
	return false
}

func parse(v string) ([]int, bool) {
	v = strings.TrimPrefix(strings.TrimSpace(v), "v")
	if v == "" {
		return nil, false
	}
	var out []int
	for _, p := range strings.Split(v, ".") {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return nil, false
		}
		out = append(out, n)
	}
	return out, true
}

func get(ctx context.Context, client *http.Client, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Diskless-PathFix-updater")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("%s: %s", url, resp.Status)
	}
	return resp, nil
}

// Latest asks GitHub (api, normally LatestURL) for the newest release.
func Latest(ctx context.Context, client *http.Client, api string) (Release, error) {
	resp, err := get(ctx, client, api)
	if err != nil {
		return Release{}, err
	}
	defer resp.Body.Close()
	var rel struct {
		Tag    string `json:"tag_name"`
		Page   string `json:"html_url"`
		Body   string `json:"body"`
		Assets []struct {
			Name   string `json:"name"`
			URL    string `json:"browser_download_url"`
			Digest string `json:"digest"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return Release{}, fmt.Errorf("reading the release: %w", err)
	}
	for _, a := range rel.Assets {
		if strings.EqualFold(a.Name, AssetName) {
			return Release{
				Version: rel.Tag,
				URL:     a.URL,
				SHA256:  strings.ToLower(strings.TrimPrefix(a.Digest, "sha256:")),
				Notes:   rel.Body,
				Page:    rel.Page,
			}, nil
		}
	}
	return Release{}, fmt.Errorf("release %s has no %s", rel.Tag, AssetName)
}

// Download saves r's exe at dst and checks its SHA-256. A file that does
// not match is deleted.
func Download(ctx context.Context, client *http.Client, r Release, dst string) error {
	if r.SHA256 == "" {
		return errors.New("the release gives no SHA-256 for its exe, so it cannot be checked")
	}
	resp, err := get(ctx, client, r.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	h := sha256.New()
	_, err = io.Copy(io.MultiWriter(f, h), resp.Body)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err == nil && hex.EncodeToString(h.Sum(nil)) != r.SHA256 {
		err = errors.New("the downloaded exe is not the one the release lists (SHA-256 differs)")
	}
	if err != nil {
		os.Remove(dst)
	}
	return err
}

// Swap puts newFile in place of exe; the old exe is kept as exe.old. If the
// new one cannot be moved in, the old one is put back.
func Swap(exe, newFile string) error {
	if _, err := os.Stat(newFile); err != nil {
		return fmt.Errorf("the new exe: %w", err)
	}
	old := exe + ".old"
	os.Remove(old) // an earlier update's; fails only if it still runs
	if _, err := os.Stat(old); err == nil {
		old = fmt.Sprintf("%s.%d.old", exe, os.Getpid())
	}
	if err := os.Rename(exe, old); err != nil {
		return fmt.Errorf("moving the running exe aside: %w", err)
	}
	if err := os.Rename(newFile, exe); err != nil {
		os.Rename(old, exe)
		return fmt.Errorf("putting the new exe in place: %w", err)
	}
	return nil
}

// Cleanup deletes what an earlier update left next to exe. Files still in
// use (an old exe the service still runs) stay until a later start.
func Cleanup(exe string) {
	os.Remove(exe + ".new")
	os.Remove(exe + ".old")
	matches, _ := filepath.Glob(exe + ".*.old")
	for _, m := range matches {
		os.Remove(m)
	}
}
