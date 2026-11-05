# Graph Report - RapTech-Portable  (2026-11-06)

## Corpus Check
- 70 files · ~52,209 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 5 file(s) not represented in the graph (top: (none) 2, .cmd 1, .md5 1)

## Summary
- 805 nodes · 2355 edges · 86 communities (28 shown, 58 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 292 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `9b095d99`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- app.go
- icon/icon.go
- testing.T
- GameFix.jsx
- links_test.go
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
- App
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
- manager/icon_test.go
- RemovePackage
- .CopyFolder
- expandVars
- gamefix/app.go
- Scan
- settings.go
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
5. `GameFixTab()` - 27 edges
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
Cohesion: 0.24
Nodes (4): sameFolder(), TestOneGlobalHuntAtATimeSkippingOurOwnFolders(), jobAccounting, WatchResult

### Community 1 - "icon/icon.go"
Cohesion: 0.20
Nodes (19): group, Copy(), groupCount(), groupNames(), load(), loadICO(), loadPE(), PNG() (+11 more)

### Community 2 - "testing.T"
Cohesion: 0.07
Nodes (53): TestListWithoutRootIsEmpty(), parseArgs(), buildLauncher(), copyFile(), TestLauncherRunsPackageEndToEnd(), TestLauncherStartsTheAppAtItsDefaultLocation(), TestParseArgs(), batLink (+45 more)

### Community 3 - "GameFix.jsx"
Cohesion: 0.07
Nodes (54): App(), AppIcon(), FolderIcon(), Icon(), KeyIcon(), lastLine(), LinkIcon(), nearest() (+46 more)

### Community 4 - "links_test.go"
Cohesion: 0.06
Nodes (51): Folder, LinkState, isReparse(), LinkStatus(), makeJunction(), MakeLink(), RemoveLink(), mustRead() (+43 more)

### Community 5 - "CreatePackage"
Cohesion: 0.08
Nodes (36): TestHuntAddedAndAddForAPortableApp(), TestAddFolderAndSession(), TestListShowsIconAndChosenIcon(), TestRemoveAppDeletesPackageEvenWithBrokenManifest(), CheckLink(), ExpandVars(), hasPrefixFold(), IsUNC() (+28 more)

### Community 6 - "Updater"
Cohesion: 0.08
Nodes (25): fail(), main(), openLog(), run(), setup(), update(), attachConsole(), imagePath() (+17 more)

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
Cohesion: 0.07
Nodes (30): find, Folder, process, State, activeProcesses(), Hunter, New(), folderIn() (+22 more)

### Community 35 - "frontend/package.json"
Cohesion: 0.10
Nodes (20): dependencies, react, react-dom, devDependencies, @types/react, @types/react-dom, vite, @vitejs/plugin-react (+12 more)

### Community 38 - "App"
Cohesion: 0.13
Nodes (4): within(), lastLines(), load(), App

### Community 52 - "App"
Cohesion: 0.18
Nodes (5): App, watch, inside(), runElevated(), watchKey()

### Community 56 - "service.go"
Cohesion: 0.14
Nodes (12): NewApp(), watchService, NewApp(), installService(), installServiceAs(), logFile(), removeService(), removeServiceNamed() (+4 more)

### Community 72 - "manager/icon_test.go"
Cohesion: 0.33
Nodes (15): ApplyIcon(), FixIcons(), HasOwnIcon(), SetIcon(), copyTo(), realLauncher(), sameIcon(), TestApplyIconAfterRobloxDownload() (+7 more)

### Community 73 - "RemovePackage"
Cohesion: 0.13
Nodes (23): TestAddFolderRefusesAFolderLinkedByAnotherApp(), closePrograms(), Use, imagePath(), InUse(), RemovePackage(), removeServices(), serviceExe() (+15 more)

### Community 74 - ".CopyFolder"
Cohesion: 0.67
Nodes (3): folderSize(), robocopy(), TestRobocopyCopiesEverything()

### Community 75 - "expandVars"
Cohesion: 0.23
Nodes (5): FolderSettings, checkSource(), displayName(), huntFolderName(), expandVars()

### Community 76 - "gamefix/app.go"
Cohesion: 0.15
Nodes (6): Folder, Game, Progress, removeEntry(), removeGameFiles(), removeTree()

### Community 79 - "Scan"
Cohesion: 0.33
Nodes (6): Update, Info, lastLine(), Scan(), TestScanFindsPackages(), PackageView

### Community 80 - "settings.go"
Cohesion: 0.47
Nodes (3): Settings, readSettings(), savedVars()

### Community 88 - "update_test.go"
Cohesion: 0.17
Nodes (13): Created(), Prune(), assertOnly(), fakeRoblox(), read(), standardPkgs(), TestPruneKeepsNewestAndCurrent(), TestUpdateInstallsLatestVersion() (+5 more)

### Community 90 - "Package"
Cohesion: 0.13
Nodes (19): TestListShowsPackageWithCCBootLine(), Manifest, Package, Load(), TestLoadReadsManifest(), TestLoadRefusesUnsafeManifests(), TestSaveThenLoadKeepsManifest(), writeManifest() (+11 more)

## Knowledge Gaps
- **3 isolated node(s):** `@types/react`, `@types/react-dom`, `raptech-portable`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 140 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **58 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `App` connect `App` to `app.go`, `CreatePackage`, `Updater`, `Package`, `RemovePackage`, `Scan`, `service.go`, `Hunter`?**
  _High betweenness centrality (0.039) - this node is a cross-community bridge._
- **What connects `@types/react`, `@types/react-dom`, `raptech-portable` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.0661189358372457 - nodes in this community are weakly interconnected._
- **Why does `App` connect `App` to `testing.T`, `Updater`, `.CopyFolder`, `expandVars`, `gamefix/app.go`, `settings.go`, `service.go`, `Hunter`?**
  _High betweenness centrality (0.037) - this node is a cross-community bridge._
- **Should `GameFix.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.07093253968253968 - nodes in this community are weakly interconnected._
- **Why does `Hunter` connect `Hunter` to `app.go`, `update_test.go`, `App`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **Should `links_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.0625 - nodes in this community are weakly interconnected._