package selfupdate

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

// TestLiveRelease reads the real newest release and downloads its exe.
// It needs the internet, so it only runs with DPF_LIVE=1.
func TestLiveRelease(t *testing.T) {
	if os.Getenv("DPF_LIVE") != "1" {
		t.Skip("set DPF_LIVE=1 to check the real GitHub release")
	}
	r, err := Latest(context.Background(), http.DefaultClient, LatestURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("newest %s, sha256 %s", r.Version, r.SHA256)
	if err := Download(context.Background(), http.DefaultClient, r, filepath.Join(t.TempDir(), AssetName)); err != nil {
		t.Fatal(err)
	}
}
