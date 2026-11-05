# Graph Report - RapTech-Portable  (2026-11-05)

## Corpus Check
- 74 files · ~53,435 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 5 file(s) not represented in the graph (top: (none) 2, .cmd 1, .md5 1)

## Summary
- 837 nodes · 2456 edges · 86 communities (28 shown, 58 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 303 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `7fab0ce3`
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
- TestRemovePackageRemovesProtocolAndSchedule
- expandVars
- gamefix/app.go
- main_test.go
- linkbat.go
- Package
- settings.go
- InitializeNotifications
- OnFileDropOff
- RemoveAllPendingNotifications
- WindowHide
- WindowMaximise

## God Nodes (most connected - your core abstractions)
1. `App` - 43 edges
2. `App` - 37 edges
3. `CreatePackage()` - 32 edges
4. `Package` - 31 edges
5. `GameFixTab()` - 29 edges
6. `Hunter` - 25 edges
7. `newTestApp()` - 22 edges
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
- `TestListShowsPackageWithCCBootLine()` --calls--> `CreateRoblox()`  [EXTRACTED]
  app_test.go → internal/manager/manager.go
- `TestAddFolderAndSession()` --calls--> `CreatePackage()`  [EXTRACTED]
  app_test.go → internal/manager/manager.go

## Import Cycles
- None detected.

## Communities (86 total, 58 thin omitted)

### Community 0 - "app.go"
Cohesion: 0.21
Nodes (8): sameFolder(), TestOneGlobalHuntAtATimeSkippingOurOwnFolders(), jobAccounting, WatchResult, dataFor(), serviceExe(), TestServiceExe(), TestDataFor()

### Community 1 - "icon/icon.go"
Cohesion: 0.20
Nodes (19): group, Copy(), groupCount(), groupNames(), load(), loadICO(), loadPE(), PNG() (+11 more)

### Community 2 - "testing.T"
Cohesion: 0.05
Nodes (68): TestAddFolderAndSession(), TestListShowsIconAndChosenIcon(), TestListShowsPackageWithCCBootLine(), TestListWithoutRootIsEmpty(), TestRemoveAppDeletesPackageEvenWithBrokenManifest(), AddRun(), HasRun(), RunLine() (+60 more)

### Community 3 - "GameFix.jsx"
Cohesion: 0.07
Nodes (56): App(), AppIcon(), FolderIcon(), Icon(), KeyIcon(), lastLine(), LinkIcon(), nearest() (+48 more)

### Community 4 - "links_test.go"
Cohesion: 0.06
Nodes (51): Folder, LinkState, isReparse(), LinkStatus(), makeJunction(), MakeLink(), RemoveLink(), mustRead() (+43 more)

### Community 5 - "CreatePackage"
Cohesion: 0.13
Nodes (27): TestHuntAddedAndAddForAPortableApp(), CheckLink(), ExpandVars(), hasPrefixFold(), IsUNC(), normPath(), Portable(), protected() (+19 more)

### Community 6 - "Updater"
Cohesion: 0.08
Nodes (27): fail(), main(), openLog(), run(), setup(), update(), attachConsole(), imagePath() (+19 more)

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
Cohesion: 0.09
Nodes (9): CCBoot, Folder, Game, gamefix, main, manager, PackageView, Program (+1 more)

### Community 15 - "RapTech Portable"
Cohesion: 0.25
Nodes (7): Build, Diskless: where updates stick, GameFix tab, Hunt tab: find the folders a program uses, RapTech Portable, The manager (RapTech-Portable.exe), What the launcher does

### Community 26 - "Hunter"
Cohesion: 0.07
Nodes (31): within(), find, Folder, process, State, activeProcesses(), Hunter, New() (+23 more)

### Community 35 - "frontend/package.json"
Cohesion: 0.10
Nodes (20): dependencies, react, react-dom, devDependencies, @types/react, @types/react-dom, vite, @vitejs/plugin-react (+12 more)

### Community 38 - "App"
Cohesion: 0.12
Nodes (5): lastLines(), load(), AddRegistry(), TestAddRegistrySavesKeyNow(), App

### Community 52 - "App"
Cohesion: 0.18
Nodes (5): App, CCBoot, watch, runElevated(), watchKey()

### Community 56 - "service.go"
Cohesion: 0.14
Nodes (12): NewApp(), watchService, NewApp(), installService(), installServiceAs(), logFile(), removeService(), removeServiceNamed() (+4 more)

### Community 72 - "manager/icon_test.go"
Cohesion: 0.33
Nodes (15): ApplyIcon(), FixIcons(), HasOwnIcon(), SetIcon(), copyTo(), realLauncher(), sameIcon(), TestApplyIconAfterRobloxDownload() (+7 more)

### Community 73 - "RemovePackage"
Cohesion: 0.21
Nodes (14): TestAddFolderRefusesAFolderLinkedByAnotherApp(), closePrograms(), Use, imagePath(), InUse(), RemovePackage(), removeServices(), newTestPackage() (+6 more)

### Community 74 - "TestRemovePackageRemovesProtocolAndSchedule"
Cohesion: 0.44
Nodes (7): TestRemovePackageRemovesProtocolAndSchedule(), RemoveSchedule(), ScheduleOn(), schtasks(), SetSchedule(), TestScheduleCreateQueryRemove(), TestScheduleRefusesOddTimes()

### Community 75 - "expandVars"
Cohesion: 0.33
Nodes (3): huntFolderName(), copyPath(), expandVars()

### Community 76 - "gamefix/app.go"
Cohesion: 0.12
Nodes (10): Folder, Game, Progress, displayName(), folderSize(), removeEntry(), removeGameFiles(), removeTree() (+2 more)

### Community 77 - "main_test.go"
Cohesion: 0.33
Nodes (7): parseArgs(), buildLauncher(), copyFile(), TestLauncherRunsPackageEndToEnd(), TestLauncherStartsTheAppAtItsDefaultLocation(), TestParseArgs(), options

### Community 78 - "linkbat.go"
Cohesion: 0.23
Nodes (12): batLink, cleanLink(), makeBat(), makeCopyBat(), protected(), removeOldBat(), setAsideForeign(), startBat() (+4 more)

### Community 79 - "Package"
Cohesion: 0.10
Nodes (23): Manifest, Package, Update, Load(), TestLoadReadsManifest(), TestLoadRefusesUnsafeManifests(), TestSaveThenLoadKeepsManifest(), writeManifest() (+15 more)

### Community 80 - "settings.go"
Cohesion: 0.20
Nodes (7): FolderSettings, Settings, checkSource(), inside(), readSettings(), savedVars(), TestServiceSettings()

## Knowledge Gaps
- **3 isolated node(s):** `@types/react`, `@types/react-dom`, `raptech-portable`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 143 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **58 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `App` connect `App` to `Updater`, `expandVars`, `gamefix/app.go`, `linkbat.go`, `settings.go`, `service.go`, `Hunter`?**
  _High betweenness centrality (0.039) - this node is a cross-community bridge._
- **What connects `@types/react`, `@types/react-dom`, `raptech-portable` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.05274725274725275 - nodes in this community are weakly interconnected._
- **Why does `App` connect `App` to `app.go`, `Updater`, `RemovePackage`, `TestRemovePackageRemovesProtocolAndSchedule`, `Package`, `service.go`, `Hunter`?**
  _High betweenness centrality (0.037) - this node is a cross-community bridge._
- **Should `GameFix.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.06946386946386947 - nodes in this community are weakly interconnected._
- **Why does `Hunter` connect `Hunter` to `app.go`, `App`, `Updater`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **Should `links_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.0625 - nodes in this community are weakly interconnected._