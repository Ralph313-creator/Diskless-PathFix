# Graph Report - RapTech-Portable  (2026-11-05)

## Corpus Check
- 47 files · ~29,888 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 5 file(s) not represented in the graph (top: (none) 2, .cmd 1, .md5 1)

## Summary
- 548 nodes · 1543 edges · 75 communities (17 shown, 58 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 191 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `b080ee21`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- app.go
- icon/icon_test.go
- update_test.go
- App.jsx
- links_test.go
- testing.T
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
- SaveRegistry
- LogDebug
- LogError
- LogFatal
- LogInfo
- LogPrint
- LogTrace
- LogWarning
- OnFileDrop
- main_test.go
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
- IsNotificationAvailable
- WindowIsMaximised
- WindowIsMinimised
- WindowIsNormal
- Quit
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
- RemoveDeliveredNotification
- WindowIsFullscreen
- WindowMinimise

## God Nodes (most connected - your core abstractions)
1. `App` - 30 edges
2. `Package` - 29 edges
3. `CreatePackage()` - 29 edges
4. `Load()` - 16 edges
5. `AddFolder()` - 16 edges
6. `load()` - 14 edges
7. `LinkStatus()` - 14 edges
8. `RemovePackage()` - 14 edges
9. `MakeLink()` - 13 edges
10. `mustWrite()` - 13 edges

## Surprising Connections (you probably didn't know these)
- `load()` --references--> `Package`  [EXTRACTED]
  app.go → internal/core/manifest.go
- `load()` --calls--> `Load()`  [EXTRACTED]
  app.go → internal/core/manifest.go
- `TestAddFolderAndSession()` --calls--> `CreatePackage()`  [EXTRACTED]
  app_test.go → internal/manager/manager.go
- `TestRemoveAppDeletesPackageEvenWithBrokenManifest()` --calls--> `CreatePackage()`  [EXTRACTED]
  app_test.go → internal/manager/manager.go
- `TestListShowsIconAndChosenIcon()` --calls--> `CreatePackage()`  [EXTRACTED]
  app_test.go → internal/manager/manager.go

## Import Cycles
- None detected.

## Communities (75 total, 58 thin omitted)

### Community 1 - "icon/icon_test.go"
Cohesion: 0.16
Nodes (20): group, activeProcesses(), Copy(), groupCount(), groupNames(), load(), loadICO(), loadPE() (+12 more)

### Community 2 - "update_test.go"
Cohesion: 0.16
Nodes (13): Created(), Prune(), assertOnly(), fakeRoblox(), read(), standardPkgs(), TestPruneKeepsNewestAndCurrent(), TestUpdateInstallsLatestVersion() (+5 more)

### Community 3 - "App.jsx"
Cohesion: 0.06
Nodes (41): dependencies, react, react-dom, devDependencies, @types/react, @types/react-dom, vite, @vitejs/plugin-react (+33 more)

### Community 4 - "links_test.go"
Cohesion: 0.11
Nodes (27): Folder, LinkState, Manifest, isReparse(), LinkStatus(), makeJunction(), MakeLink(), RemoveLink() (+19 more)

### Community 5 - "testing.T"
Cohesion: 0.06
Nodes (82): TestListShowsPackageWithCCBootLine(), TestListWithoutRootIsEmpty(), Package, hiddenPart(), Load(), TestLoadReadsManifest(), TestLoadRefusesUnsafeManifests(), TestResolveExeMissing() (+74 more)

### Community 6 - "Updater"
Cohesion: 0.09
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
Nodes (8): Environment(), InitializeNotifications(), OnFileDropOff(), RemoveAllPendingNotifications(), WindowHide(), WindowMaximise(), WindowSetMinSize(), WindowToggleMaximise()

### Community 13 - "runtime.d.ts"
Cohesion: 0.25
Nodes (7): EnvironmentInfo, NotificationAction, NotificationCategory, NotificationOptions, Position, Screen, Size

### Community 14 - "models.ts"
Cohesion: 0.17
Nodes (5): main, manager, PackageView, Program, Use

### Community 15 - "RapTech Portable"
Cohesion: 0.33
Nodes (5): Build, Diskless: where updates stick, RapTech Portable, The manager (RapTech-Portable.exe), What the launcher does

### Community 16 - "EventsOnMultiple"
Cohesion: 0.67
Nodes (3): EventsOn(), EventsOnce(), EventsOnMultiple()

### Community 26 - "SaveRegistry"
Cohesion: 0.14
Nodes (24): TestFinishSavesRegistryAndEmptiesSession(), testPackage(), TestSetupLinksAndLoadsRegistry(), TestStatusReportsEachFolder(), TestUnlinkRemovesOnlyOurProtocols(), fromUTF16(), keyExists(), LoadRegistry() (+16 more)

### Community 35 - "main_test.go"
Cohesion: 0.33
Nodes (7): parseArgs(), buildLauncher(), copyFile(), TestLauncherRunsPackageEndToEnd(), TestLauncherStartsTheAppAtItsDefaultLocation(), TestParseArgs(), options

### Community 38 - "App"
Cohesion: 0.07
Nodes (19): lastLines(), load(), NewApp(), TestAddFolderAndSession(), TestListShowsIconAndChosenIcon(), TestRemoveAppDeletesPackageEvenWithBrokenManifest(), PNG(), closePrograms() (+11 more)

## Knowledge Gaps
- **3 isolated node(s):** `@types/react`, `@types/react-dom`, `raptech-portable`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 113 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **58 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `App` connect `App` to `app.go`, `testing.T`, `Updater`?**
  _High betweenness centrality (0.038) - this node is a cross-community bridge._
- **Are the 17 inferred relationships involving `CreatePackage()` (e.g. with `TestAddFolderRefusesAFolderLinkedByAnotherApp()` and `TestSetSessionSavesSubfolders()`) actually correct?**
  _`CreatePackage()` has 17 INFERRED edges - model-reasoned connections that need verification._
- **What connects `@types/react`, `@types/react-dom`, `raptech-portable` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `App.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.06475485661424607 - nodes in this community are weakly interconnected._
- **Why does `Package` connect `testing.T` to `app.go`, `links_test.go`, `App`, `Updater`, `SaveRegistry`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Should `links_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10526315789473684 - nodes in this community are weakly interconnected._
- **Why does `NewApp()` connect `App` to `app.go`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._