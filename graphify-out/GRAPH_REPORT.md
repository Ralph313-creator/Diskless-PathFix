# Graph Report - RapTech-Portable  (2026-11-05)

## Corpus Check
- 68 files · ~51,514 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 5 file(s) not represented in the graph (top: (none) 2, .cmd 1, .md5 1)

## Summary
- 799 nodes · 2344 edges · 87 communities (29 shown, 58 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 285 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `77c86103`
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
- .Execute
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
- Package
- Load
- SamePath
- cleanLink
- .gameFolder
- expandVars
- linkbat.go
- .State
- settings.go
- CreateRoblox
- InitializeNotifications
- OnFileDropOff
- RemoveAllPendingNotifications
- WindowHide
- WindowMaximise

## God Nodes (most connected - your core abstractions)
1. `App` - 44 edges
2. `App` - 36 edges
3. `Package` - 31 edges
4. `CreatePackage()` - 31 edges
5. `GameFixTab()` - 30 edges
6. `Hunter` - 28 edges
7. `Load()` - 20 edges
8. `AddFolder()` - 19 edges
9. `newTestApp()` - 18 edges
10. `load()` - 15 edges

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

## Communities (87 total, 58 thin omitted)

### Community 0 - "app.go"
Cohesion: 0.16
Nodes (8): jobAccounting, Folder, Progress, WatchResult, dataFor(), serviceExe(), TestServiceExe(), TestDataFor()

### Community 1 - "icon/icon.go"
Cohesion: 0.12
Nodes (36): group, Copy(), groupCount(), groupNames(), load(), loadICO(), loadPE(), PNG() (+28 more)

### Community 2 - "testing.T"
Cohesion: 0.06
Nodes (57): TestListWithoutRootIsEmpty(), parseArgs(), buildLauncher(), copyFile(), TestLauncherRunsPackageEndToEnd(), TestLauncherStartsTheAppAtItsDefaultLocation(), TestParseArgs(), TestLoadReadsManifest() (+49 more)

### Community 3 - "GameFix.jsx"
Cohesion: 0.08
Nodes (55): App(), AppIcon(), FolderIcon(), HuntCard(), Icon(), KeyIcon(), lastLine(), LinkIcon() (+47 more)

### Community 4 - "links_test.go"
Cohesion: 0.06
Nodes (53): Folder, LinkState, isReparse(), LinkStatus(), makeJunction(), MakeLink(), RemoveLink(), mustRead() (+45 more)

### Community 5 - "CreatePackage"
Cohesion: 0.13
Nodes (18): TestAddFolderAndSession(), TestListShowsIconAndChosenIcon(), TestRemoveAppDeletesPackageEvenWithBrokenManifest(), linkedElsewhere(), SetSession(), TestAddFolderRefusesAFolderLinkedByAnotherApp(), TestSetSessionSavesSubfolders(), AddFolder() (+10 more)

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
Nodes (8): Environment(), IsNotificationAvailable(), Quit(), RemoveDeliveredNotification(), WindowIsFullscreen(), WindowMinimise(), WindowSetMinSize(), WindowToggleMaximise()

### Community 13 - "runtime.d.ts"
Cohesion: 0.25
Nodes (7): EnvironmentInfo, NotificationAction, NotificationCategory, NotificationOptions, Position, Screen, Size

### Community 14 - "models.ts"
Cohesion: 0.10
Nodes (8): Folder, Game, gamefix, main, manager, PackageView, Program, Use

### Community 15 - "RapTech Portable"
Cohesion: 0.25
Nodes (7): Build, Diskless: where updates stick, GameFix tab, Hunt: find the folders an app uses, RapTech Portable, The manager (RapTech-Portable.exe), What the launcher does

### Community 26 - "Hunter"
Cohesion: 0.05
Nodes (39): HuntState, find, Folder, process, Created(), activeProcesses(), Hunter, State (+31 more)

### Community 35 - "frontend/package.json"
Cohesion: 0.10
Nodes (20): dependencies, react, react-dom, devDependencies, @types/react, @types/react-dom, vite, @vitejs/plugin-react (+12 more)

### Community 38 - "App"
Cohesion: 0.12
Nodes (5): lastLines(), load(), AddRegistry(), TestAddRegistrySavesKeyNow(), App

### Community 52 - "App"
Cohesion: 0.13
Nodes (5): App, watch, runElevated(), serviceState(), watchKey()

### Community 56 - ".Execute"
Cohesion: 0.14
Nodes (11): NewApp(), watchService, NewApp(), installService(), installServiceAs(), logFile(), removeService(), removeServiceNamed() (+3 more)

### Community 72 - "Package"
Cohesion: 0.25
Nodes (13): Package, closePrograms(), Use, imagePath(), InUse(), RemovePackage(), removeServices(), newTestPackage() (+5 more)

### Community 73 - "Load"
Cohesion: 0.22
Nodes (6): inPackages(), Manifest, Update, Load(), TestSaveThenLoadKeepsManifest(), TestAddFolderImportsExistingAndSaves()

### Community 74 - "SamePath"
Cohesion: 0.25
Nodes (10): ExpandVars(), SamePath(), TestExpandVarsReplacesPkgAndEnv(), TestSamePathIgnoresCasePrefixAndSlash(), TestStartPathOutsideAnyFolderStaysInThePackage(), Added(), AddFolders(), TestAddedListsEachLinkedFolderOnce() (+2 more)

### Community 75 - "cleanLink"
Cohesion: 0.38
Nodes (4): checkSource(), displayName(), inside(), cleanLink()

### Community 76 - ".gameFolder"
Cohesion: 0.24
Nodes (4): Game, removeEntry(), removeGameFiles(), removeTree()

### Community 77 - "expandVars"
Cohesion: 0.27
Nodes (6): folderSize(), robocopy(), TestRobocopyCopiesEverything(), huntFolderName(), copyPath(), expandVars()

### Community 78 - "linkbat.go"
Cohesion: 0.38
Nodes (7): batLink, makeBat(), makeCopyBat(), protected(), removeOldBat(), setAsideForeign(), startBat()

### Community 79 - ".State"
Cohesion: 0.22
Nodes (9): CheckLink(), hasPrefixFold(), IsUNC(), normPath(), Portable(), protected(), TestCheckLinkRefusesSystemAndUserFolders(), TestIsUNC() (+1 more)

### Community 80 - "settings.go"
Cohesion: 0.31
Nodes (4): FolderSettings, Settings, readSettings(), savedVars()

### Community 81 - "CreateRoblox"
Cohesion: 0.29
Nodes (8): TestListShowsPackageWithCCBootLine(), CreateRoblox(), hasLink(), sessionHas(), TestCreateRefusesExistingPackage(), TestCreateRobloxMachineWide(), TestCreateRobloxPerUserKeepsVersionsInsideLocalAppData(), TestScanFindsPackages()

## Knowledge Gaps
- **3 isolated node(s):** `@types/react`, `@types/react-dom`, `raptech-portable`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 137 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **58 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `App` connect `App` to `app.go`, `Updater`, `cleanLink`, `.gameFolder`, `expandVars`, `linkbat.go`, `settings.go`, `.Execute`, `Hunter`?**
  _High betweenness centrality (0.045) - this node is a cross-community bridge._
- **What connects `@types/react`, `@types/react-dom`, `raptech-portable` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `icon/icon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.12367864693446089 - nodes in this community are weakly interconnected._
- **Why does `Hunter` connect `Hunter` to `app.go`, `App`, `Load`, `.State`, `App`?**
  _High betweenness centrality (0.038) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.0635814889336016 - nodes in this community are weakly interconnected._
- **Why does `App` connect `App` to `app.go`, `icon/icon.go`, `CreatePackage`, `Updater`, `Package`, `Load`, `SamePath`, `.Execute`, `Hunter`?**
  _High betweenness centrality (0.037) - this node is a cross-community bridge._
- **Should `GameFix.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.07667900581702802 - nodes in this community are weakly interconnected._