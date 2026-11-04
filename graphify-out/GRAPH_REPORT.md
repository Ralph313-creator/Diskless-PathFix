# Graph Report - RapTech-Portable  (2026-11-04)

## Corpus Check
- 51 files · ~28,440 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 5 file(s) not represented in the graph (top: (none) 2, .cmd 1, .md5 1)

## Summary
- 586 nodes · 1639 edges · 80 communities (19 shown, 61 thin omitted)
- Extraction: 89% EXTRACTED · 11% INFERRED · 0% AMBIGUOUS · INFERRED: 185 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `b64d0aab`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- app.go
- manager/icon_test.go
- update_test.go
- App.jsx
- testing.T
- App
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
- CreatePackage
- RegisterNotificationCategory
- RemoveAllDeliveredNotifications
- RemovePackage
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
- Load
- WindowIsMaximised
- WindowIsMinimised
- WindowIsNormal
- Package
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
- Scan
- main.go
- .startHunt
- InitializeNotifications
- OnFileDropOff
- RemoveAllPendingNotifications
- WindowHide
- WindowMaximise

## God Nodes (most connected - your core abstractions)
1. `App` - 33 edges
2. `Package` - 27 edges
3. `Hunter` - 25 edges
4. `CreatePackage()` - 22 edges
5. `Load()` - 20 edges
6. `AddFolder()` - 17 edges
7. `RemovePackage()` - 14 edges
8. `MakeLink()` - 13 edges
9. `SaveRegistry()` - 13 edges
10. `load()` - 12 edges

## Surprising Connections (you probably didn't know these)
- `App` --references--> `Hunter`  [EXTRACTED]
  app.go → internal/hunt/hunt.go
- `load()` --references--> `Package`  [EXTRACTED]
  app.go → internal/core/manifest.go
- `load()` --calls--> `Load()`  [EXTRACTED]
  app.go → internal/core/manifest.go
- `TestListShowsPackageWithCCBootLine()` --calls--> `CreateRoblox()`  [EXTRACTED]
  app_test.go → internal/manager/manager.go
- `TestOnlyOneHuntAtATime()` --calls--> `CreatePackage()`  [EXTRACTED]
  app_test.go → internal/manager/manager.go

## Import Cycles
- None detected.

## Communities (80 total, 61 thin omitted)

### Community 1 - "manager/icon_test.go"
Cohesion: 0.11
Nodes (23): parseArgs(), buildLauncher(), copyFile(), TestLauncherRunsPackageEndToEnd(), TestParseArgs(), group, Copy(), groupCount() (+15 more)

### Community 2 - "update_test.go"
Cohesion: 0.22
Nodes (11): assertOnly(), fakeRoblox(), read(), standardPkgs(), TestPruneKeepsNewestAndCurrent(), TestUpdateInstallsLatestVersion(), TestUpdateRefusesBadDownloads(), TestUpdateReplacesBrokenVersionFolder() (+3 more)

### Community 3 - "App.jsx"
Cohesion: 0.07
Nodes (43): dependencies, react, react-dom, devDependencies, @types/react, @types/react-dom, vite, @vitejs/plugin-react (+35 more)

### Community 4 - "testing.T"
Cohesion: 0.06
Nodes (58): TestListWithoutRootIsEmpty(), Folder, LinkState, isReparse(), LinkStatus(), makeJunction(), MakeLink(), RemoveLink() (+50 more)

### Community 5 - "App"
Cohesion: 0.12
Nodes (7): lastLines(), load(), TestAddHuntFoldersAndSession(), TestRemoveAppDeletesPackageEvenWithBrokenManifest(), SetSession(), ApplyIcon(), App

### Community 6 - "Updater"
Cohesion: 0.10
Nodes (23): fail(), main(), openLog(), run(), setup(), update(), attachConsole(), imagePath() (+15 more)

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
Cohesion: 0.17
Nodes (5): main, manager, PackageView, Program, Use

### Community 15 - "RapTech Portable"
Cohesion: 0.29
Nodes (6): Build, Diskless: where updates stick, Hunt: find the folders an app uses, RapTech Portable, The manager (RapTech-Portable.exe), What the launcher does

### Community 26 - "Hunter"
Cohesion: 0.07
Nodes (30): find, Folder, process, State, activeProcesses(), Hunter, New(), folderIn() (+22 more)

### Community 35 - "CreatePackage"
Cohesion: 0.12
Nodes (25): CheckLink(), ExpandVars(), hasPrefixFold(), normPath(), Portable(), protected(), SamePath(), TestCheckLinkRefusesSystemAndUserFolders() (+17 more)

### Community 38 - "RemovePackage"
Cohesion: 0.14
Nodes (22): LauncherName(), closePrograms(), Use, imagePath(), InUse(), RemovePackage(), removeServices(), serviceExe() (+14 more)

### Community 52 - "Load"
Cohesion: 0.20
Nodes (9): Load(), TestLoadReadsManifest(), TestLoadRefusesUnsafeManifests(), TestSaveThenLoadKeepsManifest(), writeManifest(), CheckData(), TestCheckDataStaysInsidePackage(), AddRegistry() (+1 more)

### Community 56 - "Package"
Cohesion: 0.25
Nodes (10): Manifest, Package, checkName(), CreateRoblox(), newPackage(), hasLink(), sessionHas(), TestCreateRefusesExistingPackage() (+2 more)

### Community 72 - "Scan"
Cohesion: 0.29
Nodes (7): TestListShowsPackageWithCCBootLine(), Update, Info, lastLine(), Scan(), TestScanFindsPackages(), PackageView

## Knowledge Gaps
- **3 isolated node(s):** `@types/react`, `@types/react-dom`, `raptech-portable`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 118 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **61 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `App` connect `App` to `app.go`, `manager/icon_test.go`, `CreatePackage`, `Updater`, `RemovePackage`, `Scan`, `main.go`, `.startHunt`, `Load`, `Hunter`?**
  _High betweenness centrality (0.046) - this node is a cross-community bridge._
- **What connects `@types/react`, `@types/react-dom`, `raptech-portable` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `manager/icon_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10606060606060606 - nodes in this community are weakly interconnected._
- **Why does `Hunter` connect `Hunter` to `app.go`, `.startHunt`, `testing.T`, `App`?**
  _High betweenness centrality (0.036) - this node is a cross-community bridge._
- **Should `App.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.0673758865248227 - nodes in this community are weakly interconnected._
- **Why does `Package` connect `Package` to `app.go`, `CreatePackage`, `testing.T`, `App`, `Updater`, `RemovePackage`, `Load`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.06317907444668008 - nodes in this community are weakly interconnected._