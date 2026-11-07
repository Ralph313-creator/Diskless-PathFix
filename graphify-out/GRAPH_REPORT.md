# Graph Report - Diskless-PathFix  (2026-11-07)

## Corpus Check
- 83 files · ~59,660 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 5 file(s) not represented in the graph (top: (none) 2, .cmd 1, .md5 1)

## Summary
- 898 nodes · 2653 edges · 96 communities (38 shown, 58 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 308 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `5e57b333`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- app.go
- icon/icon.go
- testing.T
- GameFix.jsx
- SaveRegistry
- CreatePackage
- gamefix/app.go
- runtime/package.json
- raptech-portable
- wails.json
- links_test.go
- runtime.js
- runtime.d.ts
- models.ts
- Diskless PathFix
- EventsOnMultiple
- Package
- CanResolveFilePaths
- CheckNotificationAuthorization
- CleanupNotifications
- ClipboardGetText
- CreateRoblox
- selfupdate.go
- EventsOffAll
- Hide
- Hunter
- LogDebug
- LogError
- LogFatal
- LogInfo
- LogPrint
- LogTrace
- main.go
- OnFileDrop
- frontend/package.json
- watch_test.go
- RemoveAllDeliveredNotifications
- App
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
- Scan
- WindowGetSize
- settings.go
- WindowIsMaximised
- App
- IsNotificationAvailable
- Guide.jsx
- WindowReload
- WindowReloadApp
- WindowSetAlwaysOnTop
- WindowSetBackgroundColour
- WindowSetDarkTheme
- WindowSetLightTheme
- WindowSetMaxSize
- RapTech Portable (shown as "Diskless PathFix")
- WindowSetSize
- WindowSetSystemDefaultTheme
- WindowSetTitle
- WindowShow
- WindowUnfullscreen
- WindowUnmaximise
- WindowUnminimise
- .DataPath
- App
- linkbat.go
- RemoveDeliveredNotification
- buildLauncher
- WindowIsFullscreen
- Load
- RegisterNotificationCategory
- RemoveNotification
- App.jsx
- InitializeNotifications
- OnFileDropOff
- RemoveAllPendingNotifications
- WindowHide
- RemovePendingNotification
- WindowIsNormal
- App
- WindowMinimise
- RemovePackage
- WindowSetPosition
- context.Context
- Hunt.jsx
- PublicTxt.jsx
- WindowToggleMaximise

## God Nodes (most connected - your core abstractions)
1. `App` - 41 edges
2. `App` - 37 edges
3. `CreatePackage()` - 33 edges
4. `GameFixTab()` - 32 edges
5. `Package` - 31 edges
6. `Hunter` - 25 edges
7. `newTestApp()` - 20 edges
8. `AddFolder()` - 19 edges
9. `Load()` - 18 edges
10. `load()` - 16 edges

## Surprising Connections (you probably didn't know these)
- `RapTech Portable (shown as "Diskless PathFix")` --references--> `inUserFolder()`  [INFERRED]
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

## Communities (96 total, 58 thin omitted)

### Community 0 - "app.go"
Cohesion: 0.17
Nodes (6): sameFolder(), TestOneGlobalHuntAtATimeSkippingOurOwnFolders(), jobAccounting, WatchResult, dataFor(), TestDataFor()

### Community 1 - "icon/icon.go"
Cohesion: 0.20
Nodes (19): group, Copy(), groupCount(), groupNames(), load(), loadICO(), loadPE(), PNG() (+11 more)

### Community 2 - "testing.T"
Cohesion: 0.16
Nodes (28): TestListWithoutRootIsEmpty(), TestVersionMatchesWailsJSON(), TestCopyMakesTheCopyTheSameAsTheSource(), TestDeleteGame(), TestGamesHaveTheirOwnFolders(), TestInside(), TestCCBootLine(), TestAddHuntFoldersNamesThemApart() (+20 more)

### Community 3 - "GameFix.jsx"
Cohesion: 0.16
Nodes (26): copyPath(), folderName(), formatSize(), GameFixTab(), Icon(), IconButton(), ICONS, PathField() (+18 more)

### Community 4 - "SaveRegistry"
Cohesion: 0.14
Nodes (24): TestFinishSavesRegistryAndEmptiesSession(), testPackage(), TestSetupLinksAndLoadsRegistry(), TestStatusReportsEachFolder(), TestUnlinkRemovesOnlyOurProtocols(), fromUTF16(), keyExists(), LoadRegistry() (+16 more)

### Community 5 - "CreatePackage"
Cohesion: 0.09
Nodes (35): TestHuntAddedAndAddForAPortableApp(), TestAddFolderAndSession(), TestListShowsIconAndChosenIcon(), TestRemoveAppDeletesPackageEvenWithBrokenManifest(), CheckLink(), ExpandVars(), hasPrefixFold(), normPath() (+27 more)

### Community 6 - "gamefix/app.go"
Cohesion: 0.14
Nodes (10): Folder, Progress, displayName(), folderSize(), removeEntry(), removeGameFiles(), removeTree(), robocopy() (+2 more)

### Community 8 - "runtime/package.json"
Cohesion: 0.12
Nodes (15): author, bugs, url, description, homepage, keywords, license, main (+7 more)

### Community 10 - "wails.json"
Cohesion: 0.12
Nodes (16): author, email, name, frontend:build, frontend:dev:serverUrl, frontend:dev:watcher, frontend:install, info (+8 more)

### Community 11 - "links_test.go"
Cohesion: 0.21
Nodes (19): isReparse(), LinkStatus(), makeJunction(), MakeLink(), RemoveLink(), mustRead(), mustWrite(), TestLinkStatusMissing() (+11 more)

### Community 12 - "runtime.js"
Cohesion: 0.22
Nodes (7): EventsOff(), LogWarning(), Quit(), WindowGetPosition(), WindowIsMinimised(), WindowMaximise(), WindowSetMinSize()

### Community 13 - "runtime.d.ts"
Cohesion: 0.25
Nodes (7): EnvironmentInfo, NotificationAction, NotificationCategory, NotificationOptions, Position, Screen, Size

### Community 14 - "models.ts"
Cohesion: 0.09
Nodes (9): Folder, Game, gamefix, main, manager, PackageView, Program, UpdateInfo (+1 more)

### Community 15 - "Diskless PathFix"
Cohesion: 0.22
Nodes (8): Build, Diskless PathFix, Diskless: where updates stick, GameFix tab, Hunt tab: find the folders a program uses, Releases and the updater, The manager (RapTech-Portable.exe), What the launcher does

### Community 17 - "Package"
Cohesion: 0.35
Nodes (15): Package, ApplyIcon(), HasOwnIcon(), SetIcon(), copyTo(), realLauncher(), sameIcon(), TestApplyIconAfterRobloxDownload() (+7 more)

### Community 22 - "CreateRoblox"
Cohesion: 0.25
Nodes (9): TestListShowsPackageWithCCBootLine(), TestPublicTxtHasEveryLine(), TestPublicTxtWithoutGameFixBat(), CreateRoblox(), hasLink(), sessionHas(), TestCreateRefusesExistingPackage(), TestCreateRobloxMachineWide() (+1 more)

### Community 23 - "selfupdate.go"
Cohesion: 0.11
Nodes (20): restartServiceOf(), UpdateInfo, TestLiveRelease(), Cleanup(), Download(), get(), Latest(), Newer() (+12 more)

### Community 26 - "Hunter"
Cohesion: 0.07
Nodes (30): find, Folder, process, State, activeProcesses(), Hunter, New(), folderIn() (+22 more)

### Community 33 - "main.go"
Cohesion: 0.10
Nodes (11): NewApp(), watchService, NewApp(), installService(), installServiceAs(), logFile(), removeService(), removeServiceNamed() (+3 more)

### Community 35 - "frontend/package.json"
Cohesion: 0.10
Nodes (20): dependencies, react, react-dom, devDependencies, @types/react, @types/react-dom, vite, @vitejs/plugin-react (+12 more)

### Community 36 - "watch_test.go"
Cohesion: 0.40
Nodes (4): TestServiceSettings(), TestWatch(), TestWatchDirStops(), waitFor()

### Community 38 - "App"
Cohesion: 0.11
Nodes (5): lastLines(), load(), App, AddRegistry(), TestAddRegistrySavesKeyNow()

### Community 39 - "update_test.go"
Cohesion: 0.22
Nodes (13): Created(), Prune(), assertOnly(), fakeRoblox(), read(), standardPkgs(), TestPruneKeepsNewestAndCurrent(), TestUpdateInstallsLatestVersion() (+5 more)

### Community 50 - "Scan"
Cohesion: 0.12
Nodes (18): ccbootLine(), isFile(), LinkState, PackageView, RunLine(), TestRunLine(), FolderStatus, Info (+10 more)

### Community 52 - "settings.go"
Cohesion: 0.15
Nodes (12): FolderSettings, Settings, checkSource(), inside(), huntFolderName(), cleanLink(), expandVars(), protected() (+4 more)

### Community 54 - "App"
Cohesion: 0.13
Nodes (6): App, Game, watch, runElevated(), serviceState(), watchKey()

### Community 56 - "Guide.jsx"
Cohesion: 0.53
Nodes (5): Guide(), HelpIcon(), Icon(), store, XIcon()

### Community 64 - "RapTech Portable (shown as "Diskless PathFix")"
Cohesion: 0.40
Nodes (4): Build, graphify, RapTech Portable (shown as "Diskless PathFix"), Rules

### Community 72 - ".DataPath"
Cohesion: 0.17
Nodes (8): Folder, hiddenPart(), TestResolveExeMissing(), TestResolveExePicksNewestVersionFolder(), EmptyFolder(), RemoveProtocol(), TestEmptyFolderKeepsFolder(), TestStartPathGoesThroughTheDefaultLocation()

### Community 74 - "linkbat.go"
Cohesion: 0.18
Nodes (18): batLink, callLink(), gameFolders(), inUserFolder(), makeBat(), makeCopyBat(), removeOldBat(), setAsideForeign() (+10 more)

### Community 76 - "buildLauncher"
Cohesion: 0.67
Nodes (4): buildLauncher(), copyFile(), TestLauncherRunsPackageEndToEnd(), TestLauncherStartsTheAppAtItsDefaultLocation()

### Community 78 - "Load"
Cohesion: 0.20
Nodes (9): Manifest, Update, Load(), TestLoadReadsManifest(), TestLoadRefusesUnsafeManifests(), TestSaveThenLoadKeepsManifest(), writeManifest(), FixIcons() (+1 more)

### Community 81 - "App.jsx"
Cohesion: 0.22
Nodes (18): AppIcon(), FolderIcon(), Icon(), KeyIcon(), lastLine(), LinkIcon(), nearest(), NewApp() (+10 more)

### Community 88 - "App"
Cohesion: 0.22
Nodes (12): App(), container, root, Icon(), Update(), UpdateButton(), UpdateIcon(), useUpdate() (+4 more)

### Community 90 - "RemovePackage"
Cohesion: 0.18
Nodes (15): closePrograms(), Use, imagePath(), InUse(), RemovePackage(), removeServices(), serviceExe(), newTestPackage() (+7 more)

### Community 92 - "context.Context"
Cohesion: 0.09
Nodes (28): fail(), main(), openLog(), parseArgs(), run(), setup(), TestParseArgs(), update() (+20 more)

### Community 93 - "Hunt.jsx"
Cohesion: 0.32
Nodes (9): Brand(), Steps(), HuntTab(), Icon(), norm(), SearchIcon(), XIcon(), GetVersion() (+1 more)

### Community 94 - "PublicTxt.jsx"
Cohesion: 0.57
Nodes (6): copyText(), FileIcon(), Icon(), PublicTxt(), XIcon(), ClipboardSetText()

## Knowledge Gaps
- **3 isolated node(s):** `@types/react`, `@types/react-dom`, `raptech-portable`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 155 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **58 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `App` connect `App` to `app.go`, `main.go`, `Hunter`, `CreatePackage`, `Scan`, `selfupdate.go`, `RemovePackage`, `context.Context`?**
  _High betweenness centrality (0.037) - this node is a cross-community bridge._
- **What connects `@types/react`, `@types/react-dom`, `raptech-portable` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `SaveRegistry` be split into smaller, more focused modules?**
  _Cohesion score 0.14333333333333334 - nodes in this community are weakly interconnected._
- **Why does `App` connect `App` to `main.go`, `gamefix/app.go`, `App`, `linkbat.go`, `Scan`, `settings.go`, `context.Context`?**
  _High betweenness centrality (0.032) - this node is a cross-community bridge._
- **Should `CreatePackage` be split into smaller, more focused modules?**
  _Cohesion score 0.09009009009009009 - nodes in this community are weakly interconnected._
- **Why does `Hunter` connect `Hunter` to `app.go`, `App`, `App`, `update_test.go`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **Should `gamefix/app.go` be split into smaller, more focused modules?**
  _Cohesion score 0.13852813852813853 - nodes in this community are weakly interconnected._