# Graph Report - RapTech-Portable  (2026-11-06)

## Corpus Check
- 71 files · ~54,100 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 5 file(s) not represented in the graph (top: (none) 2, .cmd 1, .md5 1)

## Summary
- 813 nodes · 2382 edges · 86 communities (28 shown, 58 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 292 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `69055ffc`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- app.go
- icon/icon_test.go
- settings_test.go
- GameFix.jsx
- SaveRegistry
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
- ClipboardSetText
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
- RegisterNotificationCategory
- RemoveAllDeliveredNotifications
- testing.T
- RemoveNotification
- RemoveNotificationCategory
- RemovePendingNotification
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
- WindowIsNormal
- service.go
- WindowReload
- WindowReloadApp
- WindowSetAlwaysOnTop
- WindowSetBackgroundColour
- WindowSetDarkTheme
- WindowSetLightTheme
- WindowSetMaxSize
- WindowSetPosition
- WindowSetSize
- WindowSetSystemDefaultTheme
- WindowSetTitle
- WindowShow
- WindowUnfullscreen
- WindowUnmaximise
- WindowUnminimise
- RemovePackage
- App
- linkbat.go
- gamefix/app.go
- Load
- IsUNC
- main.go
- InitializeNotifications
- OnFileDropOff
- RemoveAllPendingNotifications
- WindowHide
- WindowMaximise
- update_test.go
- Package

## God Nodes (most connected - your core abstractions)
1. `App` - 40 edges
2. `App` - 37 edges
3. `CreatePackage()` - 32 edges
4. `Package` - 31 edges
5. `GameFixTab()` - 28 edges
6. `Hunter` - 25 edges
7. `newTestApp()` - 19 edges
8. `AddFolder()` - 19 edges
9. `Load()` - 18 edges
10. `load()` - 16 edges

## Surprising Connections (you probably didn't know these)
- `App` --references--> `Hunter`  [EXTRACTED]
  app.go → internal/hunt/hunt.go
- `load()` --references--> `Package`  [EXTRACTED]
  app.go → internal/core/manifest.go
- `load()` --calls--> `Load()`  [EXTRACTED]
  app.go → internal/core/manifest.go
- `main()` --calls--> `Load()`  [EXTRACTED]
  cmd/launcher/main.go → internal/core/manifest.go
- `main()` --calls--> `SaveRegistry()`  [EXTRACTED]
  cmd/launcher/main.go → internal/core/registry.go

## Import Cycles
- None detected.

## Communities (86 total, 58 thin omitted)

### Community 0 - "app.go"
Cohesion: 0.20
Nodes (4): sameFolder(), TestOneGlobalHuntAtATimeSkippingOurOwnFolders(), jobAccounting, WatchResult

### Community 1 - "icon/icon_test.go"
Cohesion: 0.19
Nodes (18): group, Copy(), groupCount(), groupNames(), load(), loadICO(), loadPE(), Read() (+10 more)

### Community 2 - "settings_test.go"
Cohesion: 0.14
Nodes (25): TestCopyMakesTheCopyTheSameAsTheSource(), TestDeleteGame(), TestGamesHaveTheirOwnFolders(), TestInside(), TestAddHuntFoldersNamesThemApart(), TestGameSourcesListsWhatAGameHas(), gameFolders(), portable() (+17 more)

### Community 3 - "GameFix.jsx"
Cohesion: 0.07
Nodes (60): App(), AppIcon(), FolderIcon(), Icon(), KeyIcon(), lastLine(), LinkIcon(), nearest() (+52 more)

### Community 4 - "SaveRegistry"
Cohesion: 0.09
Nodes (24): Folder, LinkState, FolderStatus, fromUTF16(), keyExists(), LoadRegistry(), RegFileName(), replaceFold() (+16 more)

### Community 5 - "CreatePackage"
Cohesion: 0.12
Nodes (28): TestHuntAddedAndAddForAPortableApp(), TestAddFolderAndSession(), TestListShowsIconAndChosenIcon(), TestRemoveAppDeletesPackageEvenWithBrokenManifest(), CheckLink(), ExpandVars(), Portable(), protected() (+20 more)

### Community 6 - "Updater"
Cohesion: 0.09
Nodes (26): fail(), main(), openLog(), parseArgs(), run(), setup(), TestParseArgs(), update() (+18 more)

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
Cohesion: 0.20
Nodes (8): Environment(), IsNotificationAvailable(), Quit(), RemoveDeliveredNotification(), WindowIsFullscreen(), WindowMinimise(), WindowSetMinSize(), WindowToggleMaximise()

### Community 13 - "runtime.d.ts"
Cohesion: 0.25
Nodes (7): EnvironmentInfo, NotificationAction, NotificationCategory, NotificationOptions, Position, Screen, Size

### Community 14 - "models.ts"
Cohesion: 0.10
Nodes (8): Folder, Game, gamefix, main, manager, PackageView, Program, Use

### Community 15 - "RapTech Portable"
Cohesion: 0.25
Nodes (7): Build, Diskless: where updates stick, GameFix tab, Hunt tab: find the folders a program uses, RapTech Portable, The manager (RapTech-Portable.exe), What the launcher does

### Community 26 - "Hunter"
Cohesion: 0.06
Nodes (30): find, Folder, process, State, activeProcesses(), Hunter, New(), folderIn() (+22 more)

### Community 35 - "frontend/package.json"
Cohesion: 0.10
Nodes (20): dependencies, react, react-dom, devDependencies, @types/react, @types/react-dom, vite, @vitejs/plugin-react (+12 more)

### Community 38 - "testing.T"
Cohesion: 0.12
Nodes (34): TestListWithoutRootIsEmpty(), buildLauncher(), copyFile(), TestLauncherRunsPackageEndToEnd(), TestLauncherStartsTheAppAtItsDefaultLocation(), isReparse(), LinkStatus(), makeJunction() (+26 more)

### Community 52 - "App"
Cohesion: 0.12
Nodes (12): App, watch, checkSource(), displayName(), inside(), removeTree(), huntFolderName(), cleanLink() (+4 more)

### Community 56 - "service.go"
Cohesion: 0.11
Nodes (14): NewApp(), watchService, NewApp(), installService(), installServiceAs(), logFile(), removeService(), removeServiceNamed() (+6 more)

### Community 72 - "RemovePackage"
Cohesion: 0.18
Nodes (16): TestAddFolderRefusesAFolderLinkedByAnotherApp(), closePrograms(), Use, imagePath(), InUse(), RemovePackage(), removeServices(), serviceExe() (+8 more)

### Community 73 - "App"
Cohesion: 0.07
Nodes (16): ccbootLine(), within(), lastLines(), load(), Update, PNG(), Info, TestRemovePackageRemovesProtocolAndSchedule() (+8 more)

### Community 74 - "linkbat.go"
Cohesion: 0.20
Nodes (13): batLink, makeBat(), makeCopyBat(), protected(), removeOldBat(), setAsideForeign(), startBat(), TestBatMakesLinks() (+5 more)

### Community 76 - "gamefix/app.go"
Cohesion: 0.15
Nodes (8): Folder, Game, Progress, folderSize(), removeEntry(), removeGameFiles(), robocopy(), TestRobocopyCopiesEverything()

### Community 78 - "Load"
Cohesion: 0.10
Nodes (18): Manifest, Load(), TestLoadReadsManifest(), TestLoadRefusesUnsafeManifests(), TestSaveThenLoadKeepsManifest(), writeManifest(), CheckData(), TestCheckDataStaysInsidePackage() (+10 more)

### Community 79 - "IsUNC"
Cohesion: 0.50
Nodes (4): hasPrefixFold(), IsUNC(), normPath(), TestIsUNC()

### Community 80 - "main.go"
Cohesion: 0.14
Nodes (7): FolderSettings, Settings, readSettings(), savedVars(), TestServiceSettings(), TestWatchDirStops(), waitFor()

### Community 88 - "update_test.go"
Cohesion: 0.17
Nodes (13): Created(), Prune(), assertOnly(), fakeRoblox(), read(), standardPkgs(), TestPruneKeepsNewestAndCurrent(), TestUpdateInstallsLatestVersion() (+5 more)

### Community 90 - "Package"
Cohesion: 0.21
Nodes (23): TestListShowsPackageWithCCBootLine(), Package, ApplyIcon(), HasOwnIcon(), SetIcon(), copyTo(), realLauncher(), sameIcon() (+15 more)

## Knowledge Gaps
- **3 isolated node(s):** `@types/react`, `@types/react-dom`, `raptech-portable`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 140 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **58 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `App` connect `App` to `app.go`, `Updater`, `RemovePackage`, `Load`, `service.go`, `Hunter`?**
  _High betweenness centrality (0.038) - this node is a cross-community bridge._
- **What connects `@types/react`, `@types/react-dom`, `raptech-portable` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `settings_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.13793103448275862 - nodes in this community are weakly interconnected._
- **Why does `App` connect `App` to `Updater`, `linkbat.go`, `gamefix/app.go`, `main.go`, `service.go`, `Hunter`?**
  _High betweenness centrality (0.037) - this node is a cross-community bridge._
- **Should `GameFix.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.0676056338028169 - nodes in this community are weakly interconnected._
- **Why does `Hunter` connect `Hunter` to `app.go`, `App`, `update_test.go`, `CreatePackage`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **Should `SaveRegistry` be split into smaller, more focused modules?**
  _Cohesion score 0.09090909090909091 - nodes in this community are weakly interconnected._