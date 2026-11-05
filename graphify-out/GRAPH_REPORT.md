# Graph Report - RapTech-Portable  (2026-11-06)

## Corpus Check
- 70 files · ~52,373 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 5 file(s) not represented in the graph (top: (none) 2, .cmd 1, .md5 1)

## Summary
- 806 nodes · 2357 edges · 88 communities (29 shown, 59 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 292 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `a819d1a3`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- app.go
- icon/icon.go
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
- LinkStatus
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
- testing.T
- App
- linkbat.go
- main.go
- gamefix/app.go
- main_test.go
- .Validate
- IsUNC
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

## Communities (88 total, 59 thin omitted)

### Community 0 - "app.go"
Cohesion: 0.23
Nodes (4): sameFolder(), TestOneGlobalHuntAtATimeSkippingOurOwnFolders(), jobAccounting, WatchResult

### Community 1 - "icon/icon.go"
Cohesion: 0.10
Nodes (41): ccbootLine(), group, Update, Copy(), groupCount(), groupNames(), load(), loadICO() (+33 more)

### Community 2 - "settings_test.go"
Cohesion: 0.14
Nodes (24): TestCopyMakesTheCopyTheSameAsTheSource(), TestDeleteGame(), TestGamesHaveTheirOwnFolders(), TestInside(), TestAddHuntFoldersNamesThemApart(), TestGameSourcesListsWhatAGameHas(), gameFolders(), portable() (+16 more)

### Community 3 - "GameFix.jsx"
Cohesion: 0.07
Nodes (54): App(), AppIcon(), FolderIcon(), Icon(), KeyIcon(), lastLine(), LinkIcon(), nearest() (+46 more)

### Community 4 - "SaveRegistry"
Cohesion: 0.15
Nodes (23): TestFinishSavesRegistryAndEmptiesSession(), testPackage(), TestSetupLinksAndLoadsRegistry(), TestUnlinkRemovesOnlyOurProtocols(), fromUTF16(), keyExists(), LoadRegistry(), RegFileName() (+15 more)

### Community 5 - "CreatePackage"
Cohesion: 0.13
Nodes (30): TestHuntAddedAndAddForAPortableApp(), TestAddFolderAndSession(), TestListShowsIconAndChosenIcon(), TestRemoveAppDeletesPackageEvenWithBrokenManifest(), Load(), CheckLink(), ExpandVars(), Portable() (+22 more)

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
Cohesion: 0.10
Nodes (8): Folder, Game, gamefix, main, manager, PackageView, Program, Use

### Community 15 - "RapTech Portable"
Cohesion: 0.25
Nodes (7): Build, Diskless: where updates stick, GameFix tab, Hunt tab: find the folders a program uses, RapTech Portable, The manager (RapTech-Portable.exe), What the launcher does

### Community 26 - "Hunter"
Cohesion: 0.07
Nodes (28): find, Folder, process, State, Hunter, New(), folderIn(), huntFor() (+20 more)

### Community 35 - "frontend/package.json"
Cohesion: 0.10
Nodes (20): dependencies, react, react-dom, devDependencies, @types/react, @types/react-dom, vite, @vitejs/plugin-react (+12 more)

### Community 38 - "LinkStatus"
Cohesion: 0.18
Nodes (18): LinkState, isReparse(), LinkStatus(), makeJunction(), MakeLink(), RemoveLink(), mustRead(), mustWrite() (+10 more)

### Community 52 - "App"
Cohesion: 0.11
Nodes (10): App, Game, watch, displayName(), huntFolderName(), copyPath(), expandVars(), runElevated() (+2 more)

### Community 56 - "service.go"
Cohesion: 0.13
Nodes (13): NewApp(), watchService, NewApp(), installService(), installServiceAs(), logFile(), removeService(), removeServiceNamed() (+5 more)

### Community 72 - "testing.T"
Cohesion: 0.13
Nodes (17): TestListWithoutRootIsEmpty(), TestLinkStatusMissing(), TestMakeLinkCreatesMissingTarget(), TestMakeLinkTwiceIsNoChange(), hiddenPart(), TestLoadReadsManifest(), TestLoadRefusesUnsafeManifests(), TestResolveExeMissing() (+9 more)

### Community 73 - "App"
Cohesion: 0.06
Nodes (26): within(), lastLines(), load(), closePrograms(), Use, imagePath(), InUse(), RemovePackage() (+18 more)

### Community 74 - "linkbat.go"
Cohesion: 0.21
Nodes (14): batLink, cleanLink(), makeBat(), makeCopyBat(), protected(), removeOldBat(), setAsideForeign(), startBat() (+6 more)

### Community 75 - "main.go"
Cohesion: 0.22
Nodes (3): TestWatch(), TestWatchDirStops(), waitFor()

### Community 76 - "gamefix/app.go"
Cohesion: 0.17
Nodes (8): Folder, Progress, folderSize(), removeEntry(), removeGameFiles(), removeTree(), robocopy(), TestRobocopyCopiesEverything()

### Community 77 - "main_test.go"
Cohesion: 0.33
Nodes (7): parseArgs(), buildLauncher(), copyFile(), TestLauncherRunsPackageEndToEnd(), TestLauncherStartsTheAppAtItsDefaultLocation(), TestParseArgs(), options

### Community 79 - "IsUNC"
Cohesion: 0.50
Nodes (4): hasPrefixFold(), IsUNC(), normPath(), TestIsUNC()

### Community 80 - "settings.go"
Cohesion: 0.20
Nodes (7): FolderSettings, Settings, checkSource(), inside(), readSettings(), savedVars(), TestServiceSettings()

### Community 88 - "update_test.go"
Cohesion: 0.17
Nodes (13): Created(), Prune(), assertOnly(), fakeRoblox(), read(), standardPkgs(), TestPruneKeepsNewestAndCurrent(), TestUpdateInstallsLatestVersion() (+5 more)

### Community 90 - "Package"
Cohesion: 0.14
Nodes (14): TestListShowsPackageWithCCBootLine(), Folder, Package, RemoveProtocol(), SetSession(), checkName(), CreateRoblox(), newPackage() (+6 more)

## Knowledge Gaps
- **3 isolated node(s):** `@types/react`, `@types/react-dom`, `raptech-portable`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 140 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **59 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `App` connect `App` to `app.go`, `icon/icon.go`, `Updater`, `service.go`, `Hunter`?**
  _High betweenness centrality (0.039) - this node is a cross-community bridge._
- **What connects `@types/react`, `@types/react-dom`, `raptech-portable` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `icon/icon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10459183673469388 - nodes in this community are weakly interconnected._
- **Why does `App` connect `App` to `Updater`, `linkbat.go`, `gamefix/app.go`, `settings.go`, `service.go`, `Hunter`?**
  _High betweenness centrality (0.037) - this node is a cross-community bridge._
- **Should `settings_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.1354679802955665 - nodes in this community are weakly interconnected._
- **Why does `Hunter` connect `Hunter` to `app.go`, `App`, `update_test.go`, `CreatePackage`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **Should `GameFix.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.07093253968253968 - nodes in this community are weakly interconnected._