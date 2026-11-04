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

The launcher carries the app's own icon: the manager copies the app exe's icon
into it when the package is made, after **Update now** (Roblox has no exe before
its first download), and at start for any launcher still without one. A launcher
that is running is skipped and gets it next time.

**Change icon…** on an app's page takes the icon from any program, dll or `.ico`
file instead. It is kept in the package as `icon.ico`, so updates keep it;
**Use app's icon** goes back to the app's own.

The manager's own icon is drawn by `build\appicon.ps1` into `build\appicon.png`;
after changing it, delete `build\windows\icon.ico` so `wails build` makes a new one,
and copy the PNG to `frontend\src\assets\logo.png` (the logo in the window).

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
  and the CCBoot line to copy. Each folder has a box for the subfolders emptied
  after each player (such as `LocalStorage`, where Roblox keeps logins).

### Apps and folders used in place (Riot)

Some apps are already on the game disk, such as Riot in `G:\Online Games\Riot Games`.
Copying them would duplicate many GB, and the app would update the copy while
its own records point at the original. So a package can **use them where they
are**:

- **Add another app → Use it where it is installed**: nothing is copied; the
  launcher starts the program from its own folder (with optional arguments).
- **Add folder → Link it to a folder that already exists**: the link points to
  that folder (tagged **in place**); nothing is copied.
- Folders and programs used in place are never deleted by **Remove app**. They
  must have the same drive letter on the floor PCs (the game disk's letter).
- A folder that is already a link, such as one another app's package made, is
  refused (with the app's name) instead of being copied or taken over.

**Add Riot** does this for Riot: it finds Riot Client from Riot's own
`%ProgramData%\Riot Games\RiotClientInstalls.json` (skipping copies inside the
apps folder), starts it in place, and copies and links Riot's small folders on C:
(`%ProgramData%\Riot Games`, `%LOCALAPPDATA%\Riot Games`). The login in
`Riot Client\Data` is never copied and is emptied after each player. Vanguard is
a boot driver and stays installed in the client image. GameFix also links
`%ProgramData%\Riot Games`; use one of the two for Riot on the floor PCs.

Don't copy Riot Client into a package with **Copy the app**: when the copy runs,
Riot points `RiotClientInstalls.json` at it.

### Hunt: find the folders an app uses

Not sure where an app keeps its settings and saves? On the app's page click
**Hunt…** and choose the **installed** program (not the launcher in the package).
RapTech Portable starts it and, while it runs, finds:

- the folders of every program it starts (and the ones those start), with the
  files they load, marked **program**
- the folders on the PC's drives where files are added, changed or deleted, with
  the number of **changes**, such as `%LOCALAPPDATA%\Riot Games`

Use the app the way players do (log in, change a setting, play a little), then
close it. Hunting stops a few seconds after its last program ends, or click
**Stop** (the app keeps running). Folders that most likely belong to Windows or
other programs are hidden under **Show more** and not ticked; folders the package
already has, and the folder the app was copied from, show as **added**. Tick the
ones the app needs and click **Add folders to the package**: their files are
copied in, and they are linked from then on.

- Close other programs while hunting: every change on the PC is counted.
- Close the app (also from the tray) before hunting. A launcher that is already
  running only passes the start on to itself, so its programs are not found.
- Hunting is ported from GameFix's hunt mode.

## Build

```
build.cmd          # tests, launcher (embedded in the manager), wails build
go test ./...
```

Output: `build\bin\RapTech-Portable.exe` and `build\bin\RapPortable-Launcher.exe`.
