package core

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unicode/utf16"

	"golang.org/x/sys/windows/registry"
)

// PkgToken stands for the package folder inside saved .reg files, so a
// package saved on G: works when it shows up as another drive.
const PkgToken = "%PKG%"

// hives maps the names a key may start with to the short name reg.exe uses.
var hives = map[string]string{
	"HKCU": "HKCU", "HKEY_CURRENT_USER": "HKCU",
	"HKLM": "HKLM", "HKEY_LOCAL_MACHINE": "HKLM",
}

// splitKey splits HKCU\Software\X into its hive (HKCU or HKLM) and path.
func splitKey(key string) (hive, path string, err error) {
	head, rest, _ := strings.Cut(strings.Trim(key, `\ `), `\`)
	hive, ok := hives[strings.ToUpper(head)]
	if !ok || strings.TrimSpace(rest) == "" {
		return "", "", fmt.Errorf(`registry key %q must be HKCU\... or HKLM\...`, key)
	}
	return hive, rest, nil
}

// RegFileName is the file a key is saved in: HKCU_Software_ROBLOX Corporation.reg.
func RegFileName(key string) string {
	hive, path, err := splitKey(key)
	if err != nil {
		return "invalid.reg"
	}
	name := hive + "_" + strings.ReplaceAll(path, `\`, "_")
	name = strings.Map(func(r rune) rune {
		if strings.ContainsRune(`<>:"/\|?*`, r) || r < 32 {
			return '_'
		}
		return r
	}, name)
	return name + ".reg"
}

func keyExists(key string) bool {
	hive, path, err := splitKey(key)
	if err != nil {
		return false
	}
	root := registry.CURRENT_USER
	if hive == "HKLM" {
		root = registry.LOCAL_MACHINE
	}
	k, err := registry.OpenKey(root, path, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	k.Close()
	return true
}

// SaveRegistry exports each key into dir, writing the package folder pkg as
// %PKG%. A key that does not exist is skipped and its old file kept.
func SaveRegistry(keys []string, dir, pkg string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	var errs []error
	for _, key := range keys {
		if !keyExists(key) {
			continue
		}
		hive, path, _ := splitKey(key)
		final := filepath.Join(dir, RegFileName(key))
		tmp := final + ".tmp"
		if err := runReg("export", hive+`\`+path, tmp, "/y"); err != nil {
			errs = append(errs, fmt.Errorf("save %s: %w", key, err))
			continue
		}
		raw, err := os.ReadFile(tmp)
		if err == nil {
			raw, err = rewriteReg(raw, pkg, PkgToken)
		}
		if err == nil {
			err = os.WriteFile(tmp, raw, 0o644)
		}
		if err == nil {
			err = os.Rename(tmp, final)
		}
		if err != nil {
			os.Remove(tmp)
			errs = append(errs, fmt.Errorf("save %s: %w", key, err))
		}
	}
	return errors.Join(errs...)
}

// LoadRegistry imports each key's saved file from dir, writing %PKG% as the
// package folder pkg. A key with no file is skipped.
func LoadRegistry(keys []string, dir, pkg string) error {
	var errs []error
	for _, key := range keys {
		raw, err := os.ReadFile(filepath.Join(dir, RegFileName(key)))
		if os.IsNotExist(err) {
			continue
		}
		if err == nil {
			raw, err = rewriteReg(raw, PkgToken, pkg)
		}
		if err != nil {
			errs = append(errs, fmt.Errorf("load %s: %w", key, err))
			continue
		}
		tmp, err := os.CreateTemp("", "rp-*.reg")
		if err != nil {
			errs = append(errs, err)
			continue
		}
		tmp.Write(raw)
		tmp.Close()
		if err := runReg("import", tmp.Name()); err != nil {
			errs = append(errs, fmt.Errorf("load %s: %w", key, err))
		}
		os.Remove(tmp.Name())
	}
	return errors.Join(errs...)
}

func runReg(args ...string) error {
	cmd := exec.Command("reg.exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000} // CREATE_NO_WINDOW
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// rewriteReg replaces from with to in a .reg file (UTF-16 or ANSI), ignoring
// case. In value lines backslashes are doubled, so both forms are replaced
// there; key lines ([HKEY_...]) use single backslashes. Paths stored as hex
// (REG_EXPAND_SZ, REG_MULTI_SZ) are left as they are.
func rewriteReg(raw []byte, from, to string) ([]byte, error) {
	text, err := fromUTF16(raw)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "[") {
			lines[i] = replaceFold(line, from, to)
			continue
		}
		lines[i] = replaceFold(line, strings.ReplaceAll(from, `\`, `\\`), strings.ReplaceAll(to, `\`, `\\`))
	}
	return toUTF16(strings.Join(lines, "\n")), nil
}

// replaceFold is strings.ReplaceAll ignoring case. If lower-casing changes
// the length (rare non-ASCII), it matches exactly instead.
func replaceFold(s, old, new string) string {
	lower, lowerOld := strings.ToLower(s), strings.ToLower(old)
	if old == "" || len(lower) != len(s) || len(lowerOld) != len(old) {
		return strings.ReplaceAll(s, old, new)
	}
	var b strings.Builder
	for {
		i := strings.Index(lower, lowerOld)
		if i < 0 {
			b.WriteString(s)
			return b.String()
		}
		b.WriteString(s[:i])
		b.WriteString(new)
		s, lower = s[i+len(old):], lower[i+len(old):]
	}
}

func fromUTF16(raw []byte) (string, error) {
	if len(raw) >= 2 && raw[0] == 0xFF && raw[1] == 0xFE {
		raw = raw[2:]
		if len(raw)%2 != 0 {
			return "", errors.New("broken UTF-16 file")
		}
		u := make([]uint16, len(raw)/2)
		for i := range u {
			u[i] = binary.LittleEndian.Uint16(raw[2*i:])
		}
		return string(utf16.Decode(u)), nil
	}
	return string(bytes.TrimPrefix(raw, []byte{0xEF, 0xBB, 0xBF})), nil
}

func toUTF16(s string) []byte {
	u := utf16.Encode([]rune(s))
	out := make([]byte, 2+2*len(u))
	out[0], out[1] = 0xFF, 0xFE
	for i, c := range u {
		binary.LittleEndian.PutUint16(out[2+2*i:], c)
	}
	return out
}
