package selfupdate

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestNewer(t *testing.T) {
	cases := []struct {
		latest, current string
		want            bool
	}{
		{"v1.0.1", "1.0.0", true},
		{"1.1.0", "v1.0.9", true},
		{"v2.0.0", "1.10.3", true},
		{"v1.10.0", "1.9.0", true}, // numbers, not text
		{"v1.0.0", "1.0.0", false},
		{"v1.0.0", "1.0.1", false},
		{"v1.0", "1.0.0", false},
		{"v1.0.0.1", "1.0.0", true},
		{"", "1.0.0", false},
		{"nightly", "1.0.0", false},
	}
	for _, c := range cases {
		if got := Newer(c.latest, c.current); got != c.want {
			t.Errorf("Newer(%q, %q) = %v, want %v", c.latest, c.current, got, c.want)
		}
	}
}

func sum(b []byte) string {
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

// fakeGitHub serves a releases/latest answer and the asset it names.
func fakeGitHub(t *testing.T, asset string, body []byte, digest string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	var srv *httptest.Server
	mux.HandleFunc("/repos/x/y/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"tag_name":"v1.2.0","html_url":"https://example/rel","body":"Fixed things",
			"assets":[{"name":"other.zip","browser_download_url":"%[1]s/other","digest":""},
			          {"name":%[2]q,"browser_download_url":"%[1]s/dl","digest":%[3]q}]}`,
			srv.URL, asset, digest)
	})
	mux.HandleFunc("/dl", func(w http.ResponseWriter, r *http.Request) { w.Write(body) })
	srv = httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func TestLatestPicksTheExe(t *testing.T) {
	body := []byte("new exe")
	srv := fakeGitHub(t, AssetName, body, "sha256:"+sum(body))
	r, err := Latest(context.Background(), srv.Client(), srv.URL+"/repos/x/y/releases/latest")
	if err != nil {
		t.Fatal(err)
	}
	if r.Version != "v1.2.0" || r.URL != srv.URL+"/dl" || r.SHA256 != sum(body) || r.Notes != "Fixed things" || r.Page != "https://example/rel" {
		t.Fatalf("got %+v", r)
	}
}

func TestLatestWithoutTheExeFails(t *testing.T) {
	srv := fakeGitHub(t, "something-else.exe", nil, "")
	if _, err := Latest(context.Background(), srv.Client(), srv.URL+"/repos/x/y/releases/latest"); err == nil {
		t.Fatal("want an error when the release has no " + AssetName)
	}
}

func TestLatestHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	defer srv.Close()
	if _, err := Latest(context.Background(), srv.Client(), srv.URL+"/nope"); err == nil {
		t.Fatal("want an error on 404")
	}
}

func TestDownloadChecksTheDigest(t *testing.T) {
	body := []byte("new exe")
	srv := fakeGitHub(t, AssetName, body, "sha256:"+sum(body))
	dir := t.TempDir()
	dst := filepath.Join(dir, "app.exe.new")

	good := Release{URL: srv.URL + "/dl", SHA256: sum(body)}
	if err := Download(context.Background(), srv.Client(), good, dst); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(dst); string(b) != "new exe" {
		t.Fatalf("downloaded %q", b)
	}

	bad := Release{URL: srv.URL + "/dl", SHA256: sum([]byte("other"))}
	if err := Download(context.Background(), srv.Client(), bad, dst); err == nil {
		t.Fatal("want an error when the digest does not match")
	}
	if _, err := os.Stat(dst); !os.IsNotExist(err) {
		t.Fatal("a download with the wrong digest must be deleted")
	}

	none := Release{URL: srv.URL + "/dl"}
	if err := Download(context.Background(), srv.Client(), none, dst); err == nil {
		t.Fatal("want an error when the release gives no digest")
	}
}

func TestSwapReplacesTheExeAndKeepsTheOld(t *testing.T) {
	dir := t.TempDir()
	exe := filepath.Join(dir, "app.exe")
	os.WriteFile(exe, []byte("old"), 0o644)
	os.WriteFile(exe+".old", []byte("older"), 0o644) // left from an earlier update
	os.WriteFile(exe+".new", []byte("new"), 0o644)

	if err := Swap(exe, exe+".new"); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(exe); string(b) != "new" {
		t.Fatalf("exe is %q", b)
	}
	if b, _ := os.ReadFile(exe + ".old"); string(b) != "old" {
		t.Fatalf("old is %q", b)
	}
	if _, err := os.Stat(exe + ".new"); !os.IsNotExist(err) {
		t.Fatal(".new should be gone")
	}
}

func TestSwapPutsTheOldBackWhenTheNewIsMissing(t *testing.T) {
	dir := t.TempDir()
	exe := filepath.Join(dir, "app.exe")
	os.WriteFile(exe, []byte("old"), 0o644)
	if err := Swap(exe, exe+".new"); err == nil {
		t.Fatal("want an error")
	}
	if b, _ := os.ReadFile(exe); string(b) != "old" {
		t.Fatalf("exe is %q, the old one must be put back", b)
	}
}

func TestCleanupRemovesLeftovers(t *testing.T) {
	dir := t.TempDir()
	exe := filepath.Join(dir, "app.exe")
	os.WriteFile(exe+".old", []byte("old"), 0o644)
	os.WriteFile(exe+".new", []byte("half"), 0o644)
	Cleanup(exe)
	for _, p := range []string{exe + ".old", exe + ".new"} {
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			t.Errorf("%s should be gone", p)
		}
	}
}
