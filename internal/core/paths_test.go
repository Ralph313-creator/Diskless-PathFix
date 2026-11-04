package core

import "testing"

func TestPortableUsesMostSpecificVariable(t *testing.T) {
	t.Setenv("USERPROFILE", `C:\Users\me`)
	t.Setenv("LOCALAPPDATA", `C:\Users\me\AppData\Local`)
	t.Setenv("ProgramFiles", `C:\Program Files`)
	cases := map[string]string{
		`C:\Users\me\AppData\Local\Roblox`: `%LOCALAPPDATA%\Roblox`,
		`c:\users\me\Documents\Game`:       `%USERPROFILE%\Documents\Game`,
		`C:\Program Files\Roblox\Versions`: `%ProgramFiles%\Roblox\Versions`,
		`C:\Program FilesX\Thing`:          `C:\Program FilesX\Thing`,
		`D:\Games\GTA V`:                   `D:\Games\GTA V`,
	}
	for in, want := range cases {
		if got := Portable(in); got != want {
			t.Errorf("Portable(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestExpandVarsReplacesPkgAndEnv(t *testing.T) {
	t.Setenv("LOCALAPPDATA", `C:\Users\me\AppData\Local`)
	got := ExpandVars(`%LOCALAPPDATA%\Roblox;%pkg%\Data;%NOPE_NOT_SET%`, `G:\Portable\Roblox`)
	want := `C:\Users\me\AppData\Local\Roblox;G:\Portable\Roblox\Data;%NOPE_NOT_SET%`
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestSamePathIgnoresCasePrefixAndSlash(t *testing.T) {
	yes := [][2]string{
		{`G:\Portable\Roblox\`, `g:\portable\roblox`},
		{`\??\G:\Portable\Roblox`, `G:\Portable\Roblox`},
		{`\\?\G:\Portable`, `G:/Portable`},
		{`\\?\UNC\SERVER\share\x`, `\\server\SHARE\x`},
	}
	for _, c := range yes {
		if !SamePath(c[0], c[1]) {
			t.Errorf("SamePath(%q, %q) = false, want true", c[0], c[1])
		}
	}
	if SamePath(`G:\Portable\Roblox`, `G:\Portable\RobloxX`) {
		t.Error("different folders reported as the same")
	}
}

func TestIsUNC(t *testing.T) {
	if !IsUNC(`\\SERVER\share\Roblox`) || !IsUNC(`\\?\UNC\SERVER\share`) {
		t.Error("network path not detected")
	}
	if IsUNC(`G:\Portable`) || IsUNC(`\\?\G:\Portable`) {
		t.Error("local path reported as network")
	}
}

func TestCheckLinkRefusesSystemAndUserFolders(t *testing.T) {
	t.Setenv("USERPROFILE", `C:\Users\me`)
	t.Setenv("LOCALAPPDATA", `C:\Users\me\AppData\Local`)
	t.Setenv("ProgramFiles", `C:\Program Files`)
	bad := []string{
		`%LOCALAPPDATA%`, `C:\Users\me\AppData\Local`, `C:\Program Files`, `C:\Windows`,
		`C:\Windows\System32`, `C:\Users\me\Documents`, `C:\`, `C:`, `Roblox`, `relative\path`,
		`C:\a"b`, ``,
	}
	for _, link := range bad {
		if err := CheckLink(link); err == nil {
			t.Errorf("CheckLink(%q) accepted, want refused", link)
		}
	}
	good := []string{`%LOCALAPPDATA%\Roblox`, `C:\Program Files\Roblox\Versions`, `D:\Games\X`, `\\SERVER\share\x`}
	for _, link := range good {
		if err := CheckLink(link); err != nil {
			t.Errorf("CheckLink(%q) = %v, want ok", link, err)
		}
	}
}

func TestCheckDataStaysInsidePackage(t *testing.T) {
	for _, rel := range []string{``, `..\x`, `Data\..\..\x`, `C:\x`, `\x`, `G:x`} {
		if err := CheckData(rel); err == nil {
			t.Errorf("CheckData(%q) accepted, want refused", rel)
		}
	}
	for _, rel := range []string{`Data\LocalAppData\Roblox`, `Data/x`, `Roblox.exe`} {
		if err := CheckData(rel); err != nil {
			t.Errorf("CheckData(%q) = %v, want ok", rel, err)
		}
	}
}
