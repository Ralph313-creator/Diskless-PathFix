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
- On the floor PCs, add `run -startup -skipupdate G:\Portable\Roblox\Roblox.exe --setup`
  to `C:\CCBoot\Cmd\public.txt` (with the game disk's drive letter on the clients).
  No quotes and no `-winlogon`: with those, public.txt stopped running on the
  clients. `-startup` runs it at logon as the customer; the launcher asks for
  admin rights itself for `C:\Program Files`.
- **public.txt** (next to the tabs) shows the lines for every app and GameFix.bat
  at once, with **Copy all** and **Open public.txt** (opens it in Notepad to
  paste into). The app never writes that file itself.
- GameFix.bat's line is `run -startup -skipupdate G:\GameFix\GameFix.bat`: CCBoot
  runs it at logon as the customer, and it asks for the admin rights mklink
  needs itself. (If run as SYSTEM it skips link paths in the customer's own
  folders and logs them as `[SKIPPED]`.)

## The manager (RapTech-Portable.exe)

Runs as administrator. For new users, a guide (**Help**, next to the tabs) opens
by itself the first time and explains each tab in plain words, and each tab
starts with numbered **Getting started** steps that tick themselves as they are
done (hidden with **Hide steps**). Both live in `frontend/src/Guide.jsx`.

Choose the apps folder (default `G:\Portable`), then:

- **Add Roblox**: makes the package and downloads Roblox into it. Pick whether the
  client image has Roblox installed for everyone (Program Files) or one user.
- **Add another app**: copies the app's folder into the package's `App` folder,
  and its usual place (such as `C:\Games\My Game`) becomes a link to that copy.
  The launcher starts the app **through its usual place**, so the app still thinks
  it is installed there and never sees the package, and everything it writes
  lands in the package. Then add the other folders it keeps on C: (their files are
  copied in, without the session subfolders) and its registry keys.
- Per app: **Update now**, **Automatic** updates, **Link on this PC** / **Unlink**,
  and the CCBoot line to copy. Each folder has a box for the subfolders emptied
  after each player (such as `LocalStorage`, where Roblox keeps logins).
- A folder that is already a link made by another app's package is refused, with
  that app's name, instead of being copied or taken over.

On this PC, linking an app's folder renames the original to `<folder>.rp-backup`
(never deleted); **Unlink** or **Remove app** puts it back.

## Hunt tab: find the folders a program uses

Hunt mode is one tool for the whole program. In the **Hunt** tab click
**Choose program…** and pick the **installed** program (any app or game).
RapTech Portable starts it and, while it runs, finds the folders of every program
it starts (marked **program**) and the folders where files change (with the
number of **changes**), such as `%LOCALAPPDATA%\Riot Games`. Use it the way
players do, then close it; hunting stops a few seconds after its last program
ends, or click **Stop** (the program keeps running). One hunt runs at a time.

Then choose **Add to**: any portable app or any GameFix game. Folders it already
has show as **already in it**; likely Windows or other programs' folders are under
**Show more**. Tick the ones it needs and add them: a portable app copies their
files in and links them; a GameFix game gets them as folders to Copy and Enable.
The portable apps folder, the GameFix main folder and RapTech Portable's own
folder are never suggested. Close other programs while hunting: every change on
the PC is counted.

## GameFix tab

The **GameFix** tab is the GameFix app (`D:\Project\Diskless-GameFix`), ported:
copy each game's folders into a main folder (one folder per game), then clients
run `GameFix.bat` from the main folder's share to link them, or
`GameFix Copy.bat` to get their own copy. **Enable** keeps a folder's copy up to
date whenever its source changes, and the **RapTech GameFix Watch** service does
that while RapTech Portable is closed. To find a game's folders, use the Hunt
tab and add them to the game there.

It never touches the standalone GameFix: its settings are
`RapTech-GameFix.json` next to `RapTech-Portable.exe`, and its service is
`RapTechGameFixWatch` (not `GameFixWatch`). The first start writes no
`GameFix.bat`, and a `GameFix.bat` or `GameFix Copy.bat` it did not make is
set aside as `GameFix.bat.before-raptech` before it writes its own.

## Build

```
build.cmd          # tests, launcher (embedded in the manager), wails build
go test ./...
```

Output: `build\bin\RapTech-Portable.exe` and `build\bin\RapPortable-Launcher.exe`.
