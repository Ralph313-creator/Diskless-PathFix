# Graph Report - Diskless-PathFix  (2026-11-06)

## Corpus Check
- 77 files · ~56,823 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 5 file(s) not represented in the graph (top: (none) 2, .cmd 1, .md5 1)

## Summary
- 845 nodes · 2497 edges · 91 communities (33 shown, 58 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 298 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `05651322`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- app.go
- icon/icon.go
- settings_test.go
- GameFix.jsx
- SaveRegistry
- AddFolder
- Updater
- runtime/package.json
- raptech-portable
- wails.json
- LinkStatus
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
- testing.T
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
- main.go
- RemoveAllDeliveredNotifications
- App
- fakeRoblox
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
- .Execute
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
- .ResolveExe
- linkbat.go
- RemoveDeliveredNotification
- gamefix/app.go
- WindowIsFullscreen
- Load
- WindowSetMinSize
- settings.go
- App.jsx
- InitializeNotifications
- OnFileDropOff
- RemoveAllPendingNotifications
- WindowHide
- WindowMaximise
- Guide.jsx
- Package
- main_test.go
- Hunt.jsx
- PublicTxt.jsx

## God Nodes (most connected - your core abstractions)
1. `App` - 41 edges
2. `App` - 34 edges
3. `CreatePackage()` - 33 edges
4. `GameFixTab()` - 32 edges
5. `Package` - 31 edges
6. `Hunter` - 25 edges
7. `newTestApp()` - 20 edges
8. `AddFolder()` - 19 edges
9. `Load()` - 18 edges
10. `load()` - 16 edges

## Surprising Connections (you probably didn't know these)
- `RapTech Portable (shown as "Diskless PathFix")` --references--> `inUserFolder()`  [INFERRED]
  CLAUDE.md → internal/gamefix/linkbat.go
- `App` --references--> `Hunter`  [EXTRACTED]
  app.go → internal/hunt/hunt.go
- `load()` --references--> `Package`  [EXTRACTED]
  app.go → internal/core/manifest.go
- `load()` --calls--> `Load()`  [EXTRACTED]
  app.go → internal/core/manifest.go
- `TestHuntAddedAndAddForAPortableApp()` --calls--> `SamePath()`  [EXTRACTED]
  app_hunt_test.go → internal/core/paths.go

## Import Cycles
- None detected.

## Communities (91 total, 58 thin omitted)

### Community 0 - "app.go"
Cohesion: 0.17
Nodes (9): sameFolder(), TestOneGlobalHuntAtATimeSkippingOurOwnFolders(), jobAccounting, WatchResult, installService(), installServiceAs(), removeService(), removeServiceNamed() (+1 more)

### Community 1 - "icon/icon.go"
Cohesion: 0.08
Nodes (47): ccbootLine(), isFile(), LinkState, Manifest, PackageView, group, RunLine(), TestRunLine() (+39 more)

### Community 2 - "settings_test.go"
Cohesion: 0.15
Nodes (24): TestCopyMakesTheCopyTheSameAsTheSource(), TestDeleteGame(), TestGamesHaveTheirOwnFolders(), TestInside(), TestCCBootLine(), TestAddHuntFoldersNamesThemApart(), TestGameSourcesListsWhatAGameHas(), portable() (+16 more)

### Community 3 - "GameFix.jsx"
Cohesion: 0.16
Nodes (26): copyPath(), folderName(), formatSize(), GameFixTab(), Icon(), IconButton(), ICONS, PathField() (+18 more)

### Community 4 - "SaveRegistry"
Cohesion: 0.13
Nodes (26): TestFinishSavesRegistryAndEmptiesSession(), testPackage(), TestSetupLinksAndLoadsRegistry(), TestStatusReportsEachFolder(), TestUnlinkPutsOriginalBack(), TestUnlinkRemovesOnlyOurProtocols(), fromUTF16(), keyExists() (+18 more)

### Community 5 - "AddFolder"
Cohesion: 0.15
Nodes (21): CheckLink(), ExpandVars(), hasPrefixFold(), normPath(), Portable(), protected(), SamePath(), TestCheckLinkRefusesSystemAndUserFolders() (+13 more)

### Community 6 - "Updater"
Cohesion: 0.09
Nodes (25): fail(), main(), openLog(), run(), setup(), update(), attachConsole(), imagePath() (+17 more)

### Community 8 - "runtime/package.json"
Cohesion: 0.12
Nodes (15): author, bugs, url, description, homepage, keywords, license, main (+7 more)

### Community 10 - "wails.json"
Cohesion: 0.18
Nodes (10): author, email, name, frontend:build, frontend:dev:serverUrl, frontend:dev:watcher, frontend:install, name (+2 more)

### Community 11 - "LinkStatus"
Cohesion: 0.16
Nodes (21): isReparse(), LinkStatus(), makeJunction(), MakeLink(), RemoveLink(), mustRead(), mustWrite(), TestLinkStatusMissing() (+13 more)

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

### Community 22 - "testing.T"
Cohesion: 0.17
Nodes (20): TestHuntAddedAndAddForAPortableApp(), TestAddFolderAndSession(), TestListShowsIconAndChosenIcon(), TestListShowsPackageWithCCBootLine(), TestListWithoutRootIsEmpty(), TestPublicTxtHasEveryLine(), TestPublicTxtWithoutGameFixBat(), TestRemoveAppDeletesPackageEvenWithBrokenManifest() (+12 more)

### Community 26 - "Hunter"
Cohesion: 0.06
Nodes (31): find, Folder, process, State, activeProcesses(), TestBatsSkipUserFoldersAsSystem(), Hunter, New() (+23 more)

### Community 35 - "frontend/package.json"
Cohesion: 0.10
Nodes (20): dependencies, react, react-dom, devDependencies, @types/react, @types/react-dom, vite, @vitejs/plugin-react (+12 more)

### Community 36 - "main.go"
Cohesion: 0.20
Nodes (3): TestServiceSettings(), TestWatchDirStops(), waitFor()

### Community 38 - "App"
Cohesion: 0.06
Nodes (27): within(), lastLines(), load(), App, closePrograms(), Use, imagePath(), InUse() (+19 more)

### Community 39 - "fakeRoblox"
Cohesion: 0.33
Nodes (10): assertOnly(), fakeRoblox(), read(), standardPkgs(), TestUpdateInstallsLatestVersion(), TestUpdateRefusesBadDownloads(), TestUpdateReplacesBrokenVersionFolder(), TestUpdateWhenCurrentDownloadsNothing() (+2 more)

### Community 52 - "App"
Cohesion: 0.12
Nodes (7): App, Game, watch, displayName(), runElevated(), serviceState(), watchKey()

### Community 56 - ".Execute"
Cohesion: 0.18
Nodes (7): NewApp(), watchService, NewApp(), logFile(), exeSettingsPath(), TestSettingsNextToExe(), main()

### Community 72 - ".ResolveExe"
Cohesion: 0.22
Nodes (7): Created(), hiddenPart(), TestResolveExeMissing(), TestResolveExePicksNewestVersionFolder(), TestStartPathGoesThroughTheDefaultLocation(), Prune(), TestPruneKeepsNewestAndCurrent()

### Community 74 - "linkbat.go"
Cohesion: 0.12
Nodes (25): Build, graphify, RapTech Portable (shown as "Diskless PathFix"), Rules, batLink, callLink(), cleanLink(), gameFolders() (+17 more)

### Community 76 - "gamefix/app.go"
Cohesion: 0.18
Nodes (8): Folder, Progress, folderSize(), removeEntry(), removeGameFiles(), removeTree(), robocopy(), TestRobocopyCopiesEverything()

### Community 78 - "Load"
Cohesion: 0.38
Nodes (5): Load(), TestLoadReadsManifest(), TestLoadRefusesUnsafeManifests(), TestSaveThenLoadKeepsManifest(), writeManifest()

### Community 80 - "settings.go"
Cohesion: 0.15
Nodes (11): FolderSettings, Settings, checkSource(), inside(), huntFolderName(), copyPath(), expandVars(), isDir() (+3 more)

### Community 81 - "App.jsx"
Cohesion: 0.22
Nodes (18): AppIcon(), FolderIcon(), Icon(), KeyIcon(), lastLine(), LinkIcon(), nearest(), NewApp() (+10 more)

### Community 88 - "Guide.jsx"
Cohesion: 0.20
Nodes (10): App(), Guide(), HelpIcon(), Icon(), store, XIcon(), container, root (+2 more)

### Community 90 - "Package"
Cohesion: 0.17
Nodes (10): Folder, Package, EmptyFolder(), RemoveProtocol(), TestEmptyFolderKeepsFolder(), SetSession(), hasLink(), sessionHas() (+2 more)

### Community 92 - "main_test.go"
Cohesion: 0.33
Nodes (7): parseArgs(), buildLauncher(), copyFile(), TestLauncherRunsPackageEndToEnd(), TestLauncherStartsTheAppAtItsDefaultLocation(), TestParseArgs(), options

### Community 93 - "Hunt.jsx"
Cohesion: 0.40
Nodes (8): Brand(), Steps(), HuntTab(), Icon(), norm(), SearchIcon(), XIcon(), EventsOn()

### Community 94 - "PublicTxt.jsx"
Cohesion: 0.57
Nodes (6): copyText(), FileIcon(), Icon(), PublicTxt(), XIcon(), ClipboardSetText()

## Knowledge Gaps
- **3 isolated node(s):** `@types/react`, `@types/react-dom`, `raptech-portable`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 143 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **58 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `App` connect `App` to `app.go`, `icon/icon.go`, `Updater`, `.Execute`, `Hunter`?**
  _High betweenness centrality (0.036) - this node is a cross-community bridge._
- **What connects `@types/react`, `@types/react-dom`, `raptech-portable` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `icon/icon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.07890122735242548 - nodes in this community are weakly interconnected._
- **Why does `App` connect `App` to `icon/icon.go`, `Updater`, `linkbat.go`, `gamefix/app.go`, `settings.go`, `.Execute`, `Hunter`?**
  _High betweenness centrality (0.034) - this node is a cross-community bridge._
- **Should `settings_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.1452991452991453 - nodes in this community are weakly interconnected._
- **Why does `Hunter` connect `Hunter` to `app.go`, `.ResolveExe`, `AddFolder`, `App`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **Should `SaveRegistry` be split into smaller, more focused modules?**
  _Cohesion score 0.1282051282051282 - nodes in this community are weakly interconnected._