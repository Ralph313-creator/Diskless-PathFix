# Graph Report - RapTech-Portable  (2026-11-04)

## Corpus Check
- 51 files · ~34,103 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 5 file(s) not represented in the graph (top: (none) 2, .cmd 1, .md5 1)

## Summary
- 610 nodes · 1744 edges · 82 communities (21 shown, 61 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 210 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `d857cb8d`
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
- Package
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
- CreateRoblox
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
- Info
- main.go
- CreatePackage
- InitializeNotifications
- OnFileDropOff
- RemoveAllPendingNotifications
- WindowHide
- WindowMaximise
- load
- CheckLink

## God Nodes (most connected - your core abstractions)
1. `App` - 36 edges
2. `Package` - 30 edges
3. `CreatePackage()` - 26 edges
4. `Hunter` - 25 edges
5. `Load()` - 20 edges
6. `AddFolder()` - 17 edges
7. `load()` - 15 edges
8. `RemovePackage()` - 14 edges
9. `Package()` - 13 edges
10. `MakeLink()` - 13 edges

## Surprising Connections (you probably didn't know these)
- `App` --references--> `Hunter`  [EXTRACTED]
  app.go → internal/hunt/hunt.go
- `load()` --references--> `Package`  [EXTRACTED]
  app.go → internal/core/manifest.go
- `load()` --calls--> `Load()`  [EXTRACTED]
  app.go → internal/core/manifest.go
- `TestAddHuntFoldersAndSession()` --calls--> `CreatePackage()`  [EXTRACTED]
  app_test.go → internal/manager/manager.go
- `TestRemoveAppDeletesPackageEvenWithBrokenManifest()` --calls--> `CreatePackage()`  [EXTRACTED]
  app_test.go → internal/manager/manager.go

## Import Cycles
- None detected.

## Communities (82 total, 61 thin omitted)

### Community 1 - "manager/icon_test.go"
Cohesion: 0.11
Nodes (36): group, Copy(), groupCount(), groupNames(), load(), loadICO(), loadPE(), Read() (+28 more)

### Community 2 - "update_test.go"
Cohesion: 0.16
Nodes (13): Created(), Prune(), assertOnly(), fakeRoblox(), read(), standardPkgs(), TestPruneKeepsNewestAndCurrent(), TestUpdateInstallsLatestVersion() (+5 more)

### Community 3 - "App.jsx"
Cohesion: 0.06
Nodes (44): dependencies, react, react-dom, devDependencies, @types/react, @types/react-dom, vite, @vitejs/plugin-react (+36 more)

### Community 4 - "testing.T"
Cohesion: 0.08
Nodes (51): TestListWithoutRootIsEmpty(), buildLauncher(), copyFile(), TestLauncherRunsPackageEndToEnd(), TestParseArgs(), isReparse(), LinkStatus(), makeJunction() (+43 more)

### Community 6 - "Updater"
Cohesion: 0.08
Nodes (29): fail(), main(), openLog(), parseArgs(), run(), setup(), update(), attachConsole() (+21 more)

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
Nodes (28): find, Folder, process, State, Hunter, New(), folderIn(), huntFor() (+20 more)

### Community 35 - "Package"
Cohesion: 0.15
Nodes (17): Folder, Package, hiddenPart(), TestResolveExeMissing(), TestResolveExePicksNewestVersionFolder(), ExpandVars(), SamePath(), TestExpandVarsReplacesPkgAndEnv() (+9 more)

### Community 38 - "RemovePackage"
Cohesion: 0.14
Nodes (22): closePrograms(), Use, imagePath(), InUse(), RemovePackage(), removeServices(), serviceExe(), newTestPackage() (+14 more)

### Community 52 - "Load"
Cohesion: 0.24
Nodes (5): inPackages(), Manifest, Load(), TestSaveThenLoadKeepsManifest(), TestAddFolderImportsExistingAndSaves()

### Community 56 - "CreateRoblox"
Cohesion: 0.33
Nodes (7): TestListShowsPackageWithCCBootLine(), CreateRoblox(), hasLink(), sessionHas(), TestCreateRefusesExistingPackage(), TestCreateRobloxMachineWide(), TestCreateRobloxPerUserKeepsVersionsInsideLocalAppData()

### Community 72 - "Info"
Cohesion: 0.33
Nodes (5): LinkState, Update, FolderStatus, Info, PackageView

### Community 74 - "CreatePackage"
Cohesion: 0.18
Nodes (10): TestOnlyOneHuntAtATime(), checkName(), Copy(), CreatePackage(), newPackage(), RemoveFolder(), TestAddFolderRefusesDuplicateAndUnsafe(), TestCreatePackageCopiesAppFolder() (+2 more)

### Community 80 - "load"
Cohesion: 0.14
Nodes (9): load(), TestAddHuntFoldersAndSession(), TestListShowsIconAndChosenIcon(), TestRemoveAppDeletesPackageEvenWithBrokenManifest(), PNG(), SetSession(), TestSetSessionSavesSubfolders(), AddRegistry() (+1 more)

### Community 81 - "CheckLink"
Cohesion: 0.22
Nodes (9): CheckLink(), hasPrefixFold(), IsUNC(), normPath(), Portable(), protected(), TestCheckLinkRefusesSystemAndUserFolders(), TestIsUNC() (+1 more)

## Knowledge Gaps
- **3 isolated node(s):** `@types/react`, `@types/react-dom`, `raptech-portable`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 122 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **61 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `App` connect `App` to `app.go`, `manager/icon_test.go`, `Updater`, `RemovePackage`, `main.go`, `CreatePackage`, `load`, `Load`, `Hunter`?**
  _High betweenness centrality (0.048) - this node is a cross-community bridge._
- **What connects `@types/react`, `@types/react-dom`, `raptech-portable` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `manager/icon_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.11295681063122924 - nodes in this community are weakly interconnected._
- **Why does `Hunter` connect `Hunter` to `app.go`, `update_test.go`, `Load`, `App`?**
  _High betweenness centrality (0.035) - this node is a cross-community bridge._
- **Should `App.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.06448979591836734 - nodes in this community are weakly interconnected._
- **Why does `Package` connect `Package` to `app.go`, `manager/icon_test.go`, `testing.T`, `Updater`, `RemovePackage`, `CreatePackage`, `load`, `Load`, `CreateRoblox`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.08182349503214495 - nodes in this community are weakly interconnected._