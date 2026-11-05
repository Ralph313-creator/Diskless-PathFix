# Graph Report - RapTech-Portable  (2026-11-06)

## Corpus Check
- 76 files · ~55,571 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 5 file(s) not represented in the graph (top: (none) 2, .cmd 1, .md5 1)

## Summary
- 835 nodes · 2464 edges · 87 communities (29 shown, 58 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 294 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `b9c93423`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- app.go
- icon/icon_test.go
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
- CreateRoblox
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
- expandVars
- RemoveAllDeliveredNotifications
- TestRemovePackageRemovesProtocolAndSchedule
- RunLine
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
- WindowGetPosition
- WindowGetSize
- App
- WindowIsMaximised
- WindowIsMinimised
- IsNotificationAvailable
- service.go
- WindowReload
- WindowReloadApp
- WindowSetAlwaysOnTop
- WindowSetBackgroundColour
- WindowSetDarkTheme
- WindowSetLightTheme
- WindowSetMaxSize
- Quit
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
- RemoveDeliveredNotification
- gamefix/app.go
- WindowIsFullscreen
- Load
- WindowSetMinSize
- settings.go
- InitializeNotifications
- OnFileDropOff
- RemoveAllPendingNotifications
- WindowHide
- WindowMaximise
- Package

## God Nodes (most connected - your core abstractions)
1. `App` - 41 edges
2. `App` - 39 edges
3. `CreatePackage()` - 33 edges
4. `GameFixTab()` - 31 edges
5. `Package` - 31 edges
6. `Hunter` - 25 edges
7. `newTestApp()` - 20 edges
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
- `TestPublicTxtHasEveryLine()` --calls--> `CreatePackage()`  [EXTRACTED]
  app_test.go → internal/manager/manager.go
- `main()` --calls--> `Load()`  [EXTRACTED]
  cmd/launcher/main.go → internal/core/manifest.go

## Import Cycles
- None detected.

## Communities (87 total, 58 thin omitted)

### Community 0 - "app.go"
Cohesion: 0.20
Nodes (4): sameFolder(), TestOneGlobalHuntAtATimeSkippingOurOwnFolders(), jobAccounting, WatchResult

### Community 1 - "icon/icon_test.go"
Cohesion: 0.19
Nodes (18): group, Copy(), groupCount(), groupNames(), load(), loadICO(), loadPE(), Read() (+10 more)

### Community 2 - "testing.T"
Cohesion: 0.05
Nodes (62): TestListWithoutRootIsEmpty(), parseArgs(), buildLauncher(), copyFile(), TestLauncherRunsPackageEndToEnd(), TestLauncherStartsTheAppAtItsDefaultLocation(), TestParseArgs(), Created() (+54 more)

### Community 3 - "GameFix.jsx"
Cohesion: 0.06
Nodes (67): App(), AppIcon(), FolderIcon(), Icon(), KeyIcon(), lastLine(), LinkIcon(), nearest() (+59 more)

### Community 4 - "links_test.go"
Cohesion: 0.07
Nodes (49): Folder, LinkState, isReparse(), LinkStatus(), makeJunction(), MakeLink(), RemoveLink(), mustRead() (+41 more)

### Community 5 - "CreatePackage"
Cohesion: 0.09
Nodes (33): TestHuntAddedAndAddForAPortableApp(), TestAddFolderAndSession(), TestListShowsIconAndChosenIcon(), TestRemoveAppDeletesPackageEvenWithBrokenManifest(), CheckLink(), hasPrefixFold(), IsUNC(), normPath() (+25 more)

### Community 6 - "Updater"
Cohesion: 0.08
Nodes (28): fail(), main(), openLog(), run(), setup(), update(), attachConsole(), imagePath() (+20 more)

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
Cohesion: 0.22
Nodes (7): RegisterNotificationCategory(), RemoveNotification(), RemovePendingNotification(), WindowIsNormal(), WindowMinimise(), WindowSetPosition(), WindowToggleMaximise()

### Community 13 - "runtime.d.ts"
Cohesion: 0.25
Nodes (7): EnvironmentInfo, NotificationAction, NotificationCategory, NotificationOptions, Position, Screen, Size

### Community 14 - "models.ts"
Cohesion: 0.10
Nodes (8): Folder, Game, gamefix, main, manager, PackageView, Program, Use

### Community 15 - "RapTech Portable"
Cohesion: 0.25
Nodes (7): Build, Diskless: where updates stick, GameFix tab, Hunt tab: find the folders a program uses, RapTech Portable, The manager (RapTech-Portable.exe), What the launcher does

### Community 22 - "CreateRoblox"
Cohesion: 0.25
Nodes (9): TestListShowsPackageWithCCBootLine(), TestPublicTxtHasEveryLine(), TestPublicTxtWithoutGameFixBat(), CreateRoblox(), hasLink(), sessionHas(), TestCreateRefusesExistingPackage(), TestCreateRobloxMachineWide() (+1 more)

### Community 26 - "Hunter"
Cohesion: 0.08
Nodes (27): find, Folder, process, State, Hunter, New(), folderIn(), huntFor() (+19 more)

### Community 35 - "frontend/package.json"
Cohesion: 0.10
Nodes (20): dependencies, react, react-dom, devDependencies, @types/react, @types/react-dom, vite, @vitejs/plugin-react (+12 more)

### Community 36 - "expandVars"
Cohesion: 0.31
Nodes (5): checkSource(), inside(), huntFolderName(), expandVars(), watchKey()

### Community 38 - "TestRemovePackageRemovesProtocolAndSchedule"
Cohesion: 0.44
Nodes (7): TestRemovePackageRemovesProtocolAndSchedule(), RemoveSchedule(), ScheduleOn(), schtasks(), SetSchedule(), TestScheduleCreateQueryRemove(), TestScheduleRefusesOddTimes()

### Community 39 - "RunLine"
Cohesion: 0.40
Nodes (4): ccbootLine(), isFile(), RunLine(), TestRunLine()

### Community 52 - "App"
Cohesion: 0.17
Nodes (6): App, watch, displayName(), removeTree(), copyPath(), runElevated()

### Community 56 - "service.go"
Cohesion: 0.14
Nodes (12): NewApp(), watchService, NewApp(), installService(), installServiceAs(), logFile(), removeService(), removeServiceNamed() (+4 more)

### Community 72 - "RemovePackage"
Cohesion: 0.19
Nodes (15): closePrograms(), Use, imagePath(), InUse(), RemovePackage(), removeServices(), serviceExe(), newTestPackage() (+7 more)

### Community 73 - "App"
Cohesion: 0.08
Nodes (10): within(), lastLines(), load(), SetSession(), TestSetSessionSavesSubfolders(), AddRegistry(), RemoveFolder(), TestAddRegistrySavesKeyNow() (+2 more)

### Community 74 - "linkbat.go"
Cohesion: 0.28
Nodes (9): batLink, gameFolders(), makeBat(), makeCopyBat(), protected(), removeOldBat(), setAsideForeign(), startBat() (+1 more)

### Community 76 - "gamefix/app.go"
Cohesion: 0.15
Nodes (8): Folder, Game, Progress, folderSize(), removeEntry(), removeGameFiles(), robocopy(), TestRobocopyCopiesEverything()

### Community 78 - "Load"
Cohesion: 0.15
Nodes (12): Manifest, Update, Load(), TestSaveThenLoadKeepsManifest(), PNG(), FixIcons(), Info, lastLine() (+4 more)

### Community 80 - "settings.go"
Cohesion: 0.22
Nodes (6): FolderSettings, Settings, cleanLink(), isDir(), readSettings(), savedVars()

### Community 90 - "Package"
Cohesion: 0.35
Nodes (15): Package, ApplyIcon(), HasOwnIcon(), SetIcon(), copyTo(), realLauncher(), sameIcon(), TestApplyIconAfterRobloxDownload() (+7 more)

## Knowledge Gaps
- **3 isolated node(s):** `@types/react`, `@types/react-dom`, `raptech-portable`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 143 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **58 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `App` connect `App` to `app.go`, `CreatePackage`, `Updater`, `RunLine`, `RemovePackage`, `TestRemovePackageRemovesProtocolAndSchedule`, `Load`, `service.go`, `Hunter`?**
  _High betweenness centrality (0.041) - this node is a cross-community bridge._
- **What connects `@types/react`, `@types/react-dom`, `raptech-portable` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.052100840336134456 - nodes in this community are weakly interconnected._
- **Why does `App` connect `App` to `expandVars`, `Updater`, `RunLine`, `App`, `linkbat.go`, `gamefix/app.go`, `settings.go`, `service.go`?**
  _High betweenness centrality (0.035) - this node is a cross-community bridge._
- **Should `GameFix.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.06426484907497566 - nodes in this community are weakly interconnected._
- **Why does `Hunter` connect `Hunter` to `app.go`, `App`, `testing.T`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **Should `links_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.06721311475409836 - nodes in this community are weakly interconnected._