# Graph Report - RapTech-Portable  (2026-11-05)

## Corpus Check
- 54 files · ~37,654 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 5 file(s) not represented in the graph (top: (none) 2, .cmd 1, .md5 1)

## Summary
- 647 nodes · 1946 edges · 76 communities (16 shown, 60 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 257 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `2fc98efe`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- app.go
- icon/icon.go
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
- AddFolder
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
- WindowIsMaximised
- WindowIsMinimised
- WindowIsNormal
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
- main.go
- InitializeNotifications
- OnFileDropOff
- RemoveAllPendingNotifications
- WindowHide
- WindowMaximise

## God Nodes (most connected - your core abstractions)
1. `App` - 41 edges
2. `Package` - 35 edges
3. `CreatePackage()` - 26 edges
4. `Hunter` - 25 edges
5. `Load()` - 22 edges
6. `AddFolder()` - 20 edges
7. `SamePath()` - 17 edges
8. `load()` - 16 edges
9. `ExpandVars()` - 16 edges
10. `RemovePackage()` - 15 edges

## Surprising Connections (you probably didn't know these)
- `App` --references--> `Hunter`  [EXTRACTED]
  app.go → internal/hunt/hunt.go
- `TestRemoveAppDeletesPackageEvenWithBrokenManifest()` --calls--> `CreatePackage()`  [EXTRACTED]
  app_test.go → internal/manager/manager.go
- `main()` --calls--> `Load()`  [EXTRACTED]
  cmd/launcher/main.go → internal/core/manifest.go
- `main()` --calls--> `SaveRegistry()`  [EXTRACTED]
  cmd/launcher/main.go → internal/core/registry.go
- `setup()` --references--> `Package`  [EXTRACTED]
  cmd/launcher/main.go → internal/core/manifest.go

## Import Cycles
- None detected.

## Communities (76 total, 60 thin omitted)

### Community 1 - "icon/icon.go"
Cohesion: 0.20
Nodes (19): group, Copy(), groupCount(), groupNames(), load(), loadICO(), loadPE(), PNG() (+11 more)

### Community 2 - "update_test.go"
Cohesion: 0.16
Nodes (13): Created(), Prune(), assertOnly(), fakeRoblox(), read(), standardPkgs(), TestPruneKeepsNewestAndCurrent(), TestUpdateInstallsLatestVersion() (+5 more)

### Community 3 - "App.jsx"
Cohesion: 0.06
Nodes (45): dependencies, react, react-dom, devDependencies, @types/react, @types/react-dom, vite, @vitejs/plugin-react (+37 more)

### Community 4 - "testing.T"
Cohesion: 0.06
Nodes (61): TestListWithoutRootIsEmpty(), parseArgs(), buildLauncher(), copyFile(), TestLauncherRunsPackageEndToEnd(), TestParseArgs(), Folder, LinkState (+53 more)

### Community 5 - "App"
Cohesion: 0.05
Nodes (53): inPackages(), lastLines(), load(), TestAddHuntFoldersAndSession(), TestInPlaceAppAndFolderInList(), TestListShowsIconAndChosenIcon(), TestListShowsPackageWithCCBootLine(), TestOnlyOneHuntAtATime() (+45 more)

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
Cohesion: 0.25
Nodes (7): Apps and folders used in place (Riot), Build, Diskless: where updates stick, Hunt: find the folders an app uses, RapTech Portable, The manager (RapTech-Portable.exe), What the launcher does

### Community 26 - "Hunter"
Cohesion: 0.07
Nodes (30): find, Folder, process, State, activeProcesses(), Hunter, New(), folderIn() (+22 more)

### Community 35 - "AddFolder"
Cohesion: 0.09
Nodes (43): TestManifestAcceptsInPlacePathsButNotSystemFolders(), CheckData(), CheckLink(), CheckPlace(), ExpandVars(), hasPrefixFold(), IsFull(), IsUNC() (+35 more)

### Community 38 - "RemovePackage"
Cohesion: 0.14
Nodes (22): TestRemoveAppDeletesPackageEvenWithBrokenManifest(), closePrograms(), Use, imagePath(), InUse(), RemovePackage(), removeServices(), serviceExe() (+14 more)

## Knowledge Gaps
- **3 isolated node(s):** `@types/react`, `@types/react-dom`, `raptech-portable`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 127 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **60 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `App` connect `App` to `app.go`, `Updater`, `RemovePackage`, `main.go`, `Hunter`?**
  _High betweenness centrality (0.050) - this node is a cross-community bridge._
- **What connects `@types/react`, `@types/react-dom`, `raptech-portable` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `App.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.06352941176470588 - nodes in this community are weakly interconnected._
- **Why does `Hunter` connect `Hunter` to `app.go`, `update_test.go`, `App`?**
  _High betweenness centrality (0.033) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.05765765765765766 - nodes in this community are weakly interconnected._
- **Why does `Package` connect `App` to `app.go`, `AddFolder`, `testing.T`, `RemovePackage`, `Updater`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Should `App` be split into smaller, more focused modules?**
  _Cohesion score 0.05005107252298264 - nodes in this community are weakly interconnected._