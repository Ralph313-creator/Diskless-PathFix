# Graph Report - RapTech-Portable  (2026-11-06)

## Corpus Check
- 76 files · ~56,774 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 5 file(s) not represented in the graph (top: (none) 2, .cmd 1, .md5 1)

## Summary
- 842 nodes · 2487 edges · 95 communities (37 shown, 58 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 298 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `37d8c6ca`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- app.go
- icon/icon.go
- settings_test.go
- GameFix.jsx
- LinkStatus
- CreatePackage
- Updater
- runtime/package.json
- raptech-portable
- wails.json
- RapTech Portable
- runtime.js
- runtime.d.ts
- models.ts
- RapTech Portable
- EventsOnMultiple
- BrowserOpenURL
- CanResolveFilePaths
- CheckNotificationAuthorization
- CleanupNotifications
- ClipboardGetText
- CreateRoblox
- EventsOff
- EventsOffAll
- Hide
- Hunter
- LogDebug
- LogError
- LogFatal
- LogInfo
- LogPrint
- LogTrace
- LogWarning
- OnFileDrop
- frontend/package.json
- testing.T
- RemoveAllDeliveredNotifications
- LauncherName
- update_test.go
- RemoveNotificationCategory
- Environment
- RequestNotificationAuthorization
- ResolveFilePaths
- ScreenGetAll
- SendNotification
- SendNotificationWithActions
- Show
- WindowCenter
- WindowFullscreen
- WindowGetPosition
- WindowGetSize
- App
- WindowIsMaximised
- WindowIsMinimised
- IsNotificationAvailable
- service.go
- WindowReload
- WindowReloadApp
- WindowSetAlwaysOnTop
- WindowSetBackgroundColour
- WindowSetDarkTheme
- WindowSetLightTheme
- WindowSetMaxSize
- Quit
- WindowSetSize
- WindowSetSystemDefaultTheme
- WindowSetTitle
- WindowShow
- WindowUnfullscreen
- WindowUnmaximise
- WindowUnminimise
- InUse
- App
- linkbat.go
- RemoveDeliveredNotification
- gamefix/app.go
- WindowIsFullscreen
- Load
- WindowSetMinSize
- settings.go
- App.jsx
- InitializeNotifications
- OnFileDropOff
- RemoveAllPendingNotifications
- WindowHide
- WindowMaximise
- SaveRegistry
- Guide.jsx
- .DataPath
- Package
- main.go
- main_test.go
- Hunt.jsx
- PublicTxt.jsx

## God Nodes (most connected - your core abstractions)
1. `App` - 41 edges
2. `App` - 39 edges
3. `CreatePackage()` - 33 edges
4. `GameFixTab()` - 31 edges
5. `Package` - 31 edges
6. `Hunter` - 25 edges
7. `newTestApp()` - 20 edges
8. `AddFolder()` - 19 edges
9. `Load()` - 18 edges
10. `load()` - 16 edges

## Surprising Connections (you probably didn't know these)
- `RapTech Portable` --references--> `inUserFolder()`  [INFERRED]
  CLAUDE.md → internal/gamefix/linkbat.go
- `App` --references--> `Hunter`  [EXTRACTED]
  app.go → internal/hunt/hunt.go
- `load()` --references--> `Package`  [EXTRACTED]
  app.go → internal/core/manifest.go
- `load()` --calls--> `Load()`  [EXTRACTED]
  app.go → internal/core/manifest.go
- `TestPublicTxtHasEveryLine()` --calls--> `CreatePackage()`  [EXTRACTED]
  app_test.go → internal/manager/manager.go

## Import Cycles
- None detected.

## Communities (95 total, 58 thin omitted)

### Community 0 - "app.go"
Cohesion: 0.19
Nodes (10): sameFolder(), TestOneGlobalHuntAtATimeSkippingOurOwnFolders(), jobAccounting, Result, WatchResult, activeProcesses(), Run(), TestRunMissingExe() (+2 more)

### Community 1 - "icon/icon.go"
Cohesion: 0.20
Nodes (19): group, Copy(), groupCount(), groupNames(), load(), loadICO(), loadPE(), PNG() (+11 more)

### Community 2 - "settings_test.go"
Cohesion: 0.15
Nodes (24): TestCopyMakesTheCopyTheSameAsTheSource(), TestDeleteGame(), TestGamesHaveTheirOwnFolders(), TestInside(), TestCCBootLine(), TestAddHuntFoldersNamesThemApart(), TestGameSourcesListsWhatAGameHas(), portable() (+16 more)

### Community 3 - "GameFix.jsx"
Cohesion: 0.16
Nodes (26): copyPath(), folderName(), formatSize(), GameFixTab(), Icon(), IconButton(), ICONS, PathField() (+18 more)

### Community 4 - "LinkStatus"
Cohesion: 0.14
Nodes (23): isReparse(), LinkStatus(), makeJunction(), MakeLink(), RemoveLink(), mustRead(), mustWrite(), TestLinkStatusMissing() (+15 more)

### Community 5 - "CreatePackage"
Cohesion: 0.12
Nodes (29): TestHuntAddedAndAddForAPortableApp(), TestAddFolderAndSession(), TestListShowsIconAndChosenIcon(), TestRemoveAppDeletesPackageEvenWithBrokenManifest(), CheckLink(), ExpandVars(), hasPrefixFold(), normPath() (+21 more)

### Community 6 - "Updater"
Cohesion: 0.10
Nodes (20): fail(), main(), openLog(), run(), setup(), update(), attachConsole(), imagePath() (+12 more)

### Community 8 - "runtime/package.json"
Cohesion: 0.12
Nodes (15): author, bugs, url, description, homepage, keywords, license, main (+7 more)

### Community 10 - "wails.json"
Cohesion: 0.18
Nodes (10): author, email, name, frontend:build, frontend:dev:serverUrl, frontend:dev:watcher, frontend:install, name (+2 more)

### Community 11 - "RapTech Portable"
Cohesion: 0.40
Nodes (4): Build, graphify, RapTech Portable, Rules

### Community 12 - "runtime.js"
Cohesion: 0.22
Nodes (7): RegisterNotificationCategory(), RemoveNotification(), RemovePendingNotification(), WindowIsNormal(), WindowMinimise(), WindowSetPosition(), WindowToggleMaximise()

### Community 13 - "runtime.d.ts"
Cohesion: 0.25
Nodes (7): EnvironmentInfo, NotificationAction, NotificationCategory, NotificationOptions, Position, Screen, Size

### Community 14 - "models.ts"
Cohesion: 0.10
Nodes (8): Folder, Game, gamefix, main, manager, PackageView, Program, Use

### Community 15 - "RapTech Portable"
Cohesion: 0.25
Nodes (7): Build, Diskless: where updates stick, GameFix tab, Hunt tab: find the folders a program uses, RapTech Portable, The manager (RapTech-Portable.exe), What the launcher does

### Community 22 - "CreateRoblox"
Cohesion: 0.20
Nodes (11): TestListShowsPackageWithCCBootLine(), TestPublicTxtHasEveryLine(), TestPublicTxtWithoutGameFixBat(), checkName(), CreateRoblox(), newPackage(), hasLink(), sessionHas() (+3 more)

### Community 26 - "Hunter"
Cohesion: 0.07
Nodes (30): find, Folder, process, State, TestBatsSkipUserFoldersAsSystem(), Hunter, New(), folderIn() (+22 more)

### Community 35 - "frontend/package.json"
Cohesion: 0.10
Nodes (20): dependencies, react, react-dom, devDependencies, @types/react, @types/react-dom, vite, @vitejs/plugin-react (+12 more)

### Community 36 - "testing.T"
Cohesion: 0.14
Nodes (22): TestListWithoutRootIsEmpty(), hiddenPart(), TestResolveExeMissing(), TestResolveExePicksNewestVersionFolder(), TestFinishSavesRegistryAndEmptiesSession(), testPackage(), TestSetupLinksAndLoadsRegistry(), TestStatusReportsEachFolder() (+14 more)

### Community 38 - "LauncherName"
Cohesion: 0.20
Nodes (14): FixIcons(), Info, lastLine(), LauncherName(), Scan(), TestScanFindsPackages(), TestRemovePackageRemovesProtocolAndSchedule(), RemoveSchedule() (+6 more)

### Community 39 - "update_test.go"
Cohesion: 0.16
Nodes (13): Created(), Prune(), assertOnly(), fakeRoblox(), read(), standardPkgs(), TestPruneKeepsNewestAndCurrent(), TestUpdateInstallsLatestVersion() (+5 more)

### Community 52 - "App"
Cohesion: 0.14
Nodes (6): App, Game, watch, displayName(), runElevated(), watchKey()

### Community 56 - "service.go"
Cohesion: 0.13
Nodes (13): NewApp(), watchService, NewApp(), installService(), installServiceAs(), logFile(), removeService(), removeServiceNamed() (+5 more)

### Community 72 - "InUse"
Cohesion: 0.28
Nodes (8): closePrograms(), Use, imagePath(), InUse(), serviceExe(), TestServiceExe(), within(), Program

### Community 73 - "App"
Cohesion: 0.08
Nodes (10): ccbootLine(), within(), isFile(), lastLines(), load(), RunLine(), TestRunLine(), AddRegistry() (+2 more)

### Community 74 - "linkbat.go"
Cohesion: 0.16
Nodes (21): batLink, callLink(), cleanLink(), gameFolders(), inUserFolder(), makeBat(), makeCopyBat(), protected() (+13 more)

### Community 76 - "gamefix/app.go"
Cohesion: 0.18
Nodes (8): Folder, Progress, folderSize(), removeEntry(), removeGameFiles(), removeTree(), robocopy(), TestRobocopyCopiesEverything()

### Community 78 - "Load"
Cohesion: 0.16
Nodes (10): Manifest, Update, Load(), TestLoadReadsManifest(), TestLoadRefusesUnsafeManifests(), TestSaveThenLoadKeepsManifest(), writeManifest(), SetSession() (+2 more)

### Community 80 - "settings.go"
Cohesion: 0.16
Nodes (11): FolderSettings, Settings, checkSource(), inside(), huntFolderName(), copyPath(), expandVars(), isDir() (+3 more)

### Community 81 - "App.jsx"
Cohesion: 0.22
Nodes (18): AppIcon(), FolderIcon(), Icon(), KeyIcon(), lastLine(), LinkIcon(), nearest(), NewApp() (+10 more)

### Community 87 - "SaveRegistry"
Cohesion: 0.19
Nodes (16): fromUTF16(), keyExists(), LoadRegistry(), RegFileName(), replaceFold(), rewriteReg(), runReg(), SaveRegistry() (+8 more)

### Community 88 - "Guide.jsx"
Cohesion: 0.20
Nodes (10): App(), Guide(), HelpIcon(), Icon(), store, XIcon(), container, root (+2 more)

### Community 89 - ".DataPath"
Cohesion: 0.23
Nodes (4): Folder, LinkState, FolderStatus, RemoveProtocol()

### Community 90 - "Package"
Cohesion: 0.37
Nodes (14): Package, ApplyIcon(), HasOwnIcon(), SetIcon(), copyTo(), realLauncher(), sameIcon(), TestApplyIconAfterRobloxDownload() (+6 more)

### Community 91 - "main.go"
Cohesion: 0.20
Nodes (3): TestServiceSettings(), TestWatchDirStops(), waitFor()

### Community 92 - "main_test.go"
Cohesion: 0.33
Nodes (7): parseArgs(), buildLauncher(), copyFile(), TestLauncherRunsPackageEndToEnd(), TestLauncherStartsTheAppAtItsDefaultLocation(), TestParseArgs(), options

### Community 93 - "Hunt.jsx"
Cohesion: 0.50
Nodes (7): Steps(), HuntTab(), Icon(), norm(), SearchIcon(), XIcon(), EventsOn()

### Community 94 - "PublicTxt.jsx"
Cohesion: 0.57
Nodes (6): copyText(), FileIcon(), Icon(), PublicTxt(), XIcon(), ClipboardSetText()

## Knowledge Gaps
- **3 isolated node(s):** `@types/react`, `@types/react-dom`, `raptech-portable`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 143 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **58 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `App` connect `App` to `app.go`, `Updater`, `LauncherName`, `InUse`, `Load`, `service.go`, `Hunter`?**
  _High betweenness centrality (0.041) - this node is a cross-community bridge._
- **What connects `@types/react`, `@types/react-dom`, `raptech-portable` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `settings_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.1452991452991453 - nodes in this community are weakly interconnected._
- **Why does `App` connect `App` to `Updater`, `App`, `linkbat.go`, `gamefix/app.go`, `settings.go`, `service.go`, `Hunter`?**
  _High betweenness centrality (0.035) - this node is a cross-community bridge._
- **Should `LinkStatus` be split into smaller, more focused modules?**
  _Cohesion score 0.13768115942028986 - nodes in this community are weakly interconnected._
- **Why does `Hunter` connect `Hunter` to `app.go`, `App`, `update_test.go`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **Should `CreatePackage` be split into smaller, more focused modules?**
  _Cohesion score 0.11576354679802955 - nodes in this community are weakly interconnected._