package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/sys/windows/registry"
)

// testKey makes a fresh HKCU\Software\RapPortableTest\<name> and removes it
// when the test ends.
func testKey(t *testing.T) (full string, path string) {
	t.Helper()
	path = `Software\RapPortableTest\` + strings.ReplaceAll(t.Name(), "/", "_")
	t.Cleanup(func() { deleteTree(registry.CURRENT_USER, path) })
	deleteTree(registry.CURRENT_USER, path)
	return `HKCU\` + path, path
}

func deleteTree(root registry.Key, path string) {
	k, err := registry.OpenKey(root, path, registry.ALL_ACCESS)
	if err != nil {
		return
	}
	names, _ := k.ReadSubKeyNames(-1)
	k.Close()
	for _, n := range names {
		deleteTree(root, path+`\`+n)
	}
	registry.DeleteKey(root, path)
}

func TestRegistrySaveAndLoadMovesPackagePath(t *testing.T) {
	full, path := testKey(t)
	k, _, err := registry.CreateKey(registry.CURRENT_USER, path+`\Channel`, registry.ALL_ACCESS)
	if err != nil {
		t.Fatal(err)
	}
	k.SetStringValue("", `"G:\Portable\Roblox\Roblox.exe" %1`)
	k.SetStringValue("Other", `C:\Somewhere\Else`)
	k.SetStringValue("Mixed", `g:\PORTABLE\roblox\logs`)
	k.SetDWordValue("Count", 7)
	k.Close()

	dir := t.TempDir()
	if err := SaveRegistry([]string{full}, dir, `G:\Portable\Roblox`); err != nil {
		t.Fatal(err)
	}
	saved := decodeUTF16(t, filepath.Join(dir, RegFileName(full)))
	if strings.Contains(strings.ToLower(saved), `g:\\portable\\roblox`) || !strings.Contains(saved, `%PKG%`) {
		t.Fatalf("package path not replaced by %%PKG%%:\n%s", saved)
	}

	deleteTree(registry.CURRENT_USER, path)
	if err := LoadRegistry([]string{full}, dir, `E:\Apps\Roblox`); err != nil {
		t.Fatal(err)
	}
	k, err = registry.OpenKey(registry.CURRENT_USER, path+`\Channel`, registry.QUERY_VALUE)
	if err != nil {
		t.Fatal("key not restored:", err)
	}
	defer k.Close()
	if v, _, _ := k.GetStringValue(""); v != `"E:\Apps\Roblox\Roblox.exe" %1` {
		t.Errorf("default = %q", v)
	}
	if v, _, _ := k.GetStringValue("Other"); v != `C:\Somewhere\Else` {
		t.Errorf("Other = %q", v)
	}
	if v, _, _ := k.GetStringValue("Mixed"); v != `E:\Apps\Roblox\logs` {
		t.Errorf("Mixed = %q (case-insensitive match failed)", v)
	}
	if v, _, _ := k.GetIntegerValue("Count"); v != 7 {
		t.Errorf("Count = %d", v)
	}
}

func TestSaveRegistrySkipsMissingKey(t *testing.T) {
	full, _ := testKey(t)
	dir := t.TempDir()
	if err := SaveRegistry([]string{full}, dir, `G:\P`); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, RegFileName(full))); !os.IsNotExist(err) {
		t.Fatal("a file was written for a missing key")
	}
}

func TestLoadRegistryWithoutFileIsFine(t *testing.T) {
	full, _ := testKey(t)
	if err := LoadRegistry([]string{full}, t.TempDir(), `G:\P`); err != nil {
		t.Fatal(err)
	}
}

func TestRegFileNameIsSafe(t *testing.T) {
	got := RegFileName(`HKEY_CURRENT_USER\Software\ROBLOX Corporation`)
	if got != "HKCU_Software_ROBLOX Corporation.reg" {
		t.Fatalf("got %q", got)
	}
}

func decodeUTF16(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s, err := fromUTF16(raw)
	if err != nil {
		t.Fatal(err)
	}
	return s
}
