package main

import (
	"encoding/json"
	"os"
	"testing"
)

// The exe's file version (wails.json) and the updater's version must agree.
func TestVersionMatchesWailsJSON(t *testing.T) {
	b, err := os.ReadFile("wails.json")
	if err != nil {
		t.Fatal(err)
	}
	var w struct {
		Info struct {
			ProductVersion string `json:"productVersion"`
		} `json:"info"`
	}
	if err := json.Unmarshal(b, &w); err != nil {
		t.Fatal(err)
	}
	if w.Info.ProductVersion != Version {
		t.Fatalf("wails.json productVersion %q, Version %q", w.Info.ProductVersion, Version)
	}
}
