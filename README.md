# RapTech Portable

*Created by **RapTech***

Makes apps portable for CCBoot diskless PCs. Everything an app keeps on C:
(its program folder, AppData, registry keys) lives in **one package folder
next to its launcher**, on any drive: `G:\Portable\Roblox`. When the app updates
itself, the update lands in the package too.

```
G:\Portable\Roblox\
├─ Roblox.exe                 <- the launcher (start Roblox with this)
├─ portable.json              <- which folders and keys belong to the app
├─ Data\ProgramFiles\Roblox\Versions\version-…   <- the program and its updates
├─ Data\LocalAppData\Roblox\…                    <- settings (logins wiped)
└─ Registry\HKCU_Software_ROBLOX Corporation.reg
```

## What the launcher does

1. Links each folder where the app expects it to the package (NTFS junction):
   `C:\Program Files\Roblox\Versions` → `G:\Portable\Roblox\Data\ProgramFiles\Roblox\Versions`.
   A real folder in the way is renamed to `…\Versions.rp-backup`, never deleted.
2. Loads the saved registry keys. Paths inside them are stored as `%PKG%`, so the
   package works when the drive letter differs.
3. Points `roblox-player:` / `roblox:` links at the launcher, so the browser's
   Play button always starts the newest version.
4. Starts the app and waits until it and everything it started have closed.
5. Saves the registry keys back and empties the session folders (Roblox logins),
   so the next customer is not logged in as the last one.

| Command | What it does |
|---|---|
| `Roblox.exe [args]` | set up, start Roblox, wait, save |
| `Roblox.exe --setup` | links, registry and protocols only (for CCBoot at logon) |
| `Roblox.exe --update` | update the package (run on the server) |
| `Roblox.exe --unlink` | remove the links, put the original folders back |
| `Roblox.exe --save` | save the registry into the package |
| `Roblox.exe --status` | print what is linked (from a console: `start /wait`) |

Everything is logged to `RapPortable.log` in the package.

## Diskless: where updates stick

Floor PCs throw away what they write to the game disk at restart. So an update
only sticks when it runs **on the server** (or a super client):

- Roblox: `--update` downloads the newest Roblox Player straight from Roblox
  (`clientsettings.roblox.com`, `setup.rbxcdn.com`), checks every file's MD5,
  builds it in a hidden folder and moves it in when complete, then keeps the
  newest 2 versions. A version still in use is kept until next time.
- The manager's **Automatic** setting makes a Task Scheduler task
  (`\RapTech Portable\<app>`) that runs `--update` as SYSTEM.
- On the floor PCs, add `run -winlogon "G:\Portable\Roblox\Roblox.exe" --setup`
  to `C:\CCBoot\Cmd\public.txt` (with the game disk's drive letter on the clients).

## The manager (RapTech-Portable.exe)

Runs as administrator. Choose the apps folder (default `G:\Portable`), then:

- **Add Roblox**: makes the package and downloads Roblox into it. Pick whether the
  client image has Roblox installed for everyone (Program Files) or one user.
- **Add another app**: copies the app's folder into the package. Then add the
  folders it keeps elsewhere (their files are copied in, without the session
  subfolders) and its registry keys.
- Per app: **Update now**, **Automatic** updates, **Link on this PC** / **Unlink**,
  and the CCBoot line to copy.

## Build

```
build.cmd          # tests, launcher (embedded in the manager), wails build
go test ./...
```

Output: `build\bin\RapTech-Portable.exe` and `build\bin\RapPortable-Launcher.exe`.
