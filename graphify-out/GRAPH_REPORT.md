# Graph Report - RapTech-Portable  (2026-11-06)

## Corpus Check
- 74 files · ~54,550 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 5 file(s) not represented in the graph (top: (none) 2, .cmd 1, .md5 1)

## Summary
- 847 nodes · 2494 edges · 92 communities (33 shown, 59 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 310 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `bc21e40a`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- app.go
- icon/icon_test.go
- newTestApp
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
- manager/icon_test.go
- RemovePackage
- TestRemovePackageRemovesProtocolAndSchedule
- expandVars
- gamefix/app.go
- main_test.go
- linkbat.go
- Scan
- settings.go
- testing.T
- InitializeNotifications
- OnFileDropOff
- RemoveAllPendingNotifications
- WindowHide
- WindowMaximise
- ccboot/ccboot_test.go
- update_test.go
- main.go
- CreateRoblox
- Game

## God Nodes (most connected - your core abstractions)
1. `App` - 44 edges
2. `App` - 37 edges
3. `CreatePackage()` - 32 edges
4. `Package` - 31 edges
5. `GameFixTab()` - 30 edges
6. `Hunter` - 25 edges
7. `newTestApp()` - 24 edges
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
- `TestRemoveAppDeletesPackageEvenWithBrokenManifest()` --calls--> `CreatePackage()`  [EXTRACTED]
  app_test.go → internal/manager/manager.go
- `TestListShowsIconAndChosenIcon()` --calls--> `CreatePackage()`  [EXTRACTED]
  app_test.go → internal/manager/manager.go

## Import Cycles
- None detected.

## Communities (92 total, 59 thin omitted)

### Community 0 - "app.go"
Cohesion: 0.18
Nodes (11): sameFolder(), TestOneGlobalHuntAtATimeSkippingOurOwnFolders(), jobAccounting, WatchResult, installService(), installServiceAs(), removeService(), removeServiceNamed() (+3 more)

### Community 1 - "icon/icon_test.go"
Cohesion: 0.18
Nodes (19): group, Copy(), groupCount(), groupNames(), load(), loadICO(), loadPE(), PNG() (+11 more)

### Community 2 - "newTestApp"
Cohesion: 0.18
Nodes (19): TestCopyMakesTheCopyTheSameAsTheSource(), TestDeleteGame(), TestGamesHaveTheirOwnFolders(), TestInside(), TestAddHuntFoldersNamesThemApart(), TestGameSourcesListsWhatAGameHas(), gameFolders(), portable() (+11 more)

### Community 3 - "GameFix.jsx"
Cohesion: 0.07
Nodes (57): App(), AppIcon(), FolderIcon(), Icon(), KeyIcon(), lastLine(), LinkIcon(), nearest() (+49 more)

### Community 4 - "links_test.go"
Cohesion: 0.06
Nodes (51): Folder, LinkState, isReparse(), LinkStatus(), makeJunction(), MakeLink(), RemoveLink(), mustRead() (+43 more)

### Community 5 - "CreatePackage"
Cohesion: 0.13
Nodes (28): TestHuntAddedAndAddForAPortableApp(), TestAddFolderAndSession(), Package, Load(), TestSaveThenLoadKeepsManifest(), ExpandVars(), SamePath(), TestExpandVarsReplacesPkgAndEnv() (+20 more)

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
Cohesion: 0.09
Nodes (9): CCBoot, Folder, Game, gamefix, main, manager, PackageView, Program (+1 more)

### Community 15 - "RapTech Portable"
Cohesion: 0.25
Nodes (7): Build, Diskless: where updates stick, GameFix tab, Hunt tab: find the folders a program uses, RapTech Portable, The manager (RapTech-Portable.exe), What the launcher does

### Community 26 - "Hunter"
Cohesion: 0.06
Nodes (35): find, Folder, process, State, CheckLink(), Portable(), protected(), TestCheckLinkRefusesSystemAndUserFolders() (+27 more)

### Community 35 - "frontend/package.json"
Cohesion: 0.10
Nodes (20): dependencies, react, react-dom, devDependencies, @types/react, @types/react-dom, vite, @vitejs/plugin-react (+12 more)

### Community 38 - "App"
Cohesion: 0.10
Nodes (7): within(), lastLines(), load(), SetSession(), AddRegistry(), RemoveFolder(), App

### Community 52 - "App"
Cohesion: 0.16
Nodes (6): App, CCBoot, watch, runElevated(), serviceState(), watchKey()

### Community 56 - ".Execute"
Cohesion: 0.18
Nodes (7): NewApp(), watchService, NewApp(), logFile(), exeSettingsPath(), TestSettingsNextToExe(), main()

### Community 72 - "manager/icon_test.go"
Cohesion: 0.33
Nodes (15): ApplyIcon(), FixIcons(), HasOwnIcon(), SetIcon(), copyTo(), realLauncher(), sameIcon(), TestApplyIconAfterRobloxDownload() (+7 more)

### Community 73 - "RemovePackage"
Cohesion: 0.19
Nodes (15): closePrograms(), Use, imagePath(), InUse(), RemovePackage(), removeServices(), serviceExe(), newTestPackage() (+7 more)

### Community 74 - "TestRemovePackageRemovesProtocolAndSchedule"
Cohesion: 0.44
Nodes (7): TestRemovePackageRemovesProtocolAndSchedule(), RemoveSchedule(), ScheduleOn(), schtasks(), SetSchedule(), TestScheduleCreateQueryRemove(), TestScheduleRefusesOddTimes()

### Community 75 - "expandVars"
Cohesion: 0.22
Nodes (5): FolderSettings, huntFolderName(), copyPath(), expandVars(), TestExpandVars()

### Community 76 - "gamefix/app.go"
Cohesion: 0.14
Nodes (9): Folder, Progress, displayName(), folderSize(), removeEntry(), removeGameFiles(), removeTree(), robocopy() (+1 more)

### Community 77 - "main_test.go"
Cohesion: 0.33
Nodes (7): parseArgs(), buildLauncher(), copyFile(), TestLauncherRunsPackageEndToEnd(), TestLauncherStartsTheAppAtItsDefaultLocation(), TestParseArgs(), options

### Community 78 - "linkbat.go"
Cohesion: 0.20
Nodes (14): batLink, cleanLink(), makeBat(), makeCopyBat(), protected(), removeOldBat(), setAsideForeign(), startBat() (+6 more)

### Community 79 - "Scan"
Cohesion: 0.25
Nodes (7): Manifest, Update, Info, lastLine(), Scan(), TestScanFindsPackages(), PackageView

### Community 80 - "settings.go"
Cohesion: 0.27
Nodes (6): Settings, checkSource(), inside(), readSettings(), savedVars(), TestServiceSettings()

### Community 81 - "testing.T"
Cohesion: 0.11
Nodes (21): TestListShowsIconAndChosenIcon(), TestListWithoutRootIsEmpty(), TestRemoveAppDeletesPackageEvenWithBrokenManifest(), hiddenPart(), TestLoadReadsManifest(), TestLoadRefusesUnsafeManifests(), TestResolveExeMissing(), TestResolveExePicksNewestVersionFolder() (+13 more)

### Community 87 - "ccboot/ccboot_test.go"
Cohesion: 0.22
Nodes (21): AddRun(), HasRun(), RemoveRun(), RunLine(), runs(), read(), script(), TestAddRunGoesAfterTheLastCommandAndKeepsTheRest() (+13 more)

### Community 88 - "update_test.go"
Cohesion: 0.17
Nodes (13): Created(), Prune(), assertOnly(), fakeRoblox(), read(), standardPkgs(), TestPruneKeepsNewestAndCurrent(), TestUpdateInstallsLatestVersion() (+5 more)

### Community 89 - "main.go"
Cohesion: 0.22
Nodes (3): TestWatch(), TestWatchDirStops(), waitFor()

### Community 90 - "CreateRoblox"
Cohesion: 0.33
Nodes (7): TestListShowsPackageWithCCBootLine(), CreateRoblox(), hasLink(), sessionHas(), TestCreateRefusesExistingPackage(), TestCreateRobloxMachineWide(), TestCreateRobloxPerUserKeepsVersionsInsideLocalAppData()

## Knowledge Gaps
- **3 isolated node(s):** `@types/react`, `@types/react-dom`, `raptech-portable`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 143 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **59 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `App` connect `App` to `Updater`, `expandVars`, `gamefix/app.go`, `linkbat.go`, `settings.go`, `.Execute`, `Hunter`, `Game`?**
  _High betweenness centrality (0.040) - this node is a cross-community bridge._
- **What connects `@types/react`, `@types/react-dom`, `raptech-portable` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `GameFix.jsx` be split into smaller, more focused modules?**
  _Cohesion score 0.06874717322478517 - nodes in this community are weakly interconnected._
- **Why does `App` connect `App` to `app.go`, `CreatePackage`, `Updater`, `RemovePackage`, `TestRemovePackageRemovesProtocolAndSchedule`, `Scan`, `.Execute`, `Hunter`?**
  _High betweenness centrality (0.037) - this node is a cross-community bridge._
- **Should `links_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.06349206349206349 - nodes in this community are weakly interconnected._
- **Why does `Hunter` connect `Hunter` to `app.go`, `update_test.go`, `App`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **Should `CreatePackage` be split into smaller, more focused modules?**
  _Cohesion score 0.13333333333333333 - nodes in this community are weakly interconnected._