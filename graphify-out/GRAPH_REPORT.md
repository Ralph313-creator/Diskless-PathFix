# Graph Report - RapTech-Portable  (2026-11-04)

## Corpus Check
- 45 files · ~24,791 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 5 file(s) not represented in the graph (top: (none) 2, .cmd 1, .md5 1)

## Summary
- 526 nodes · 1387 edges · 80 communities (19 shown, 61 thin omitted)
- Extraction: 89% EXTRACTED · 11% INFERRED · 0% AMBIGUOUS · INFERRED: 151 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `67287be7`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- app.go
- CreateRoblox
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
- PackageView
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
- Load
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
- Package
- WindowIsMaximised
- WindowIsMinimised
- WindowIsNormal
- Run
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
- load
- main.go
- IsUNC
- InitializeNotifications
- OnFileDropOff
- RemoveAllPendingNotifications
- WindowHide
- WindowMaximise

## God Nodes (most connected - your core abstractions)
1. `App` - 31 edges
2. `Hunter` - 25 edges
3. `Package` - 23 edges
4. `Load()` - 19 edges
5. `CreatePackage()` - 18 edges
6. `AddFolder()` - 16 edges
7. `SaveRegistry()` - 13 edges
8. `LinkStatus()` - 12 edges
9. `MakeLink()` - 12 edges
10. `mustWrite()` - 12 edges

## Surprising Connections (you probably didn't know these)
- `App` --references--> `Hunter`  [EXTRACTED]
  app.go → internal/hunt/hunt.go
- `load()` --references--> `Package`  [EXTRACTED]
  app.go → internal/core/manifest.go
- `load()` --calls--> `Load()`  [EXTRACTED]
  app.go → internal/core/manifest.go
- `TestOnlyOneHuntAtATime()` --calls--> `CreatePackage()`  [EXTRACTED]
  app_test.go → internal/manager/manager.go
- `TestAddHuntFoldersAndSession()` --calls--> `CreatePackage()`  [EXTRACTED]
  app_test.go → internal/manager/manager.go

## Import Cycles
- None detected.

## Communities (80 total, 61 thin omitted)

### Community 0 - "app.go"
Cohesion: 0.18
Nodes (9): attachConsole(), imagePath(), otherInstances(), runElevated(), jobAccounting, process, processList(), dataFor() (+1 more)

### Community 1 - "CreateRoblox"
Cohesion: 0.19
Nodes (12): TestListShowsPackageWithCCBootLine(), Update, checkName(), CreateRoblox(), Info, lastLine(), LauncherName(), newPackage() (+4 more)

### Community 2 - "update_test.go"
Cohesion: 0.13
Nodes (16): Created(), hiddenPart(), TestResolveExeMissing(), TestResolveExePicksNewestVersionFolder(), Prune(), assertOnly(), fakeRoblox(), read() (+8 more)

### Community 3 - "App.jsx"
Cohesion: 0.07
Nodes (42): dependencies, react, react-dom, devDependencies, @types/react, @types/react-dom, vite, @vitejs/plugin-react (+34 more)

### Community 4 - "testing.T"
Cohesion: 0.07
Nodes (54): TestListWithoutRootIsEmpty(), parseArgs(), buildLauncher(), copyFile(), TestLauncherRunsPackageEndToEnd(), TestParseArgs(), Folder, LinkState (+46 more)

### Community 6 - "Updater"
Cohesion: 0.18
Nodes (9): CheckData(), TestCheckDataStaysInsidePackage(), cleanLeftovers(), complete(), extract(), lock(), writeEntry(), pkgEntry (+1 more)

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

### Community 15 - "RapTech Portable"
Cohesion: 0.29
Nodes (6): Build, Diskless: where updates stick, Hunt: find the folders an app uses, RapTech Portable, The manager (RapTech-Portable.exe), What the launcher does

### Community 26 - "Hunter"
Cohesion: 0.07
Nodes (28): find, Folder, State, activeProcesses(), Hunter, New(), folderIn(), huntFor() (+20 more)

### Community 35 - "CreatePackage"
Cohesion: 0.15
Nodes (22): CheckLink(), ExpandVars(), Portable(), protected(), SamePath(), TestCheckLinkRefusesSystemAndUserFolders(), TestExpandVarsReplacesPkgAndEnv(), TestPortableUsesMostSpecificVariable() (+14 more)

### Community 38 - "Load"
Cohesion: 0.18
Nodes (8): inPackages(), TestOnlyOneHuntAtATime(), Manifest, Load(), TestLoadReadsManifest(), TestLoadRefusesUnsafeManifests(), TestSaveThenLoadKeepsManifest(), writeManifest()

### Community 52 - "Package"
Cohesion: 0.26
Nodes (12): fail(), main(), openLog(), run(), setup(), update(), Package, hasLink() (+4 more)

### Community 56 - "Run"
Cohesion: 0.23
Nodes (11): Result, Run(), TestRunMissingExe(), TestRunTimeoutEndsTheTree(), TestRunWaitsForGrandchildren(), RemoveSchedule(), ScheduleOn(), schtasks() (+3 more)

### Community 72 - "load"
Cohesion: 0.20
Nodes (6): load(), TestAddHuntFoldersAndSession(), AddRegistry(), RemoveFolder(), TestAddRegistrySavesKeyNow(), TestRemoveFolderUnlinksAndKeepsData()

### Community 74 - "IsUNC"
Cohesion: 0.50
Nodes (4): hasPrefixFold(), IsUNC(), normPath(), TestIsUNC()

## Knowledge Gaps
- **3 isolated node(s):** `@types/react`, `@types/react-dom`, `raptech-portable`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 110 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **61 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `App` connect `App` to `app.go`, `CreateRoblox`, `Updater`, `Load`, `load`, `main.go`, `Run`, `Hunter`?**
  _High betweenness centrality (0.051) - this node is a cross-community bridge._
- **What connects `@types/react`, `@types/react-dom`, `raptech-portable` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `update_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.12615384615384614 - nodes in this community are weakly interconnected._
- **Why does `Hunter` connect `Hunter` to `app.go`, `update_test.go`, `App`, `Load`?**
  _High betweenness centrality (0.040) - this node is a cross-community bridge._
- **Should `App.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.06845513413506013 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.06874717322478517 - nodes in this community are weakly interconnected._
- **Should `App.js` be split into smaller, more focused modules?**
  _Cohesion score 0.09523809523809523 - nodes in this community are weakly interconnected._