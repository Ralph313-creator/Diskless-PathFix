# Graph Report - Diskless-PathFix  (2026-11-07)

## Corpus Check
- 82 files · ~59,576 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 5 file(s) not represented in the graph (top: (none) 2, .cmd 1, .md5 1)

## Summary
- 896 nodes · 2642 edges · 96 communities (39 shown, 57 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 306 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `8a399cd4`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- app.go
- icon/icon.go
- testing.T
- GameFix.jsx
- links_test.go
- CreatePackage
- context.Context
- runtime/package.json
- raptech-portable
- wails.json
- IsUNC
- runtime.js
- runtime.d.ts
- models.ts
- Diskless PathFix
- EventsOnMultiple
- manager/icon_test.go
- CanResolveFilePaths
- CheckNotificationAuthorization
- CleanupNotifications
- ClipboardGetText
- CreateRoblox
- selfupdate_test.go
- EventsOffAll
- Hide
- Hunter
- LogDebug
- LogError
- LogFatal
- LogInfo
- LogPrint
- LogTrace
- Latest
- OnFileDrop
- frontend/package.json
- gameFolders
- RemoveAllDeliveredNotifications
- App
- update_test.go
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
- TestRemovePackageRemovesProtocolAndSchedule
- WindowGetSize
- App
- WindowIsMaximised
- .PublicTxt
- IsNotificationAvailable
- Guide.jsx
- WindowReload
- WindowReloadApp
- WindowSetAlwaysOnTop
- WindowSetBackgroundColour
- WindowSetDarkTheme
- WindowSetLightTheme
- WindowSetMaxSize
- RapTech Portable (shown as "Diskless PathFix")
- WindowSetSize
- WindowSetSystemDefaultTheme
- WindowSetTitle
- WindowShow
- WindowUnfullscreen
- WindowUnmaximise
- WindowUnminimise
- .ResolveExe
- Run
- linkbat.go
- RemoveDeliveredNotification
- CheckData
- WindowIsFullscreen
- Load
- RegisterNotificationCategory
- RemoveNotification
- App.jsx
- InitializeNotifications
- OnFileDropOff
- RemoveAllPendingNotifications
- WindowHide
- RemovePendingNotification
- WindowIsNormal
- App
- WindowMinimise
- Package
- WindowSetPosition
- main
- Hunt.jsx
- PublicTxt.jsx
- WindowToggleMaximise

## God Nodes (most connected - your core abstractions)
1. `App` - 41 edges
2. `App` - 37 edges
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
- `TestPublicTxtHasEveryLine()` --calls--> `CreatePackage()`  [EXTRACTED]
  app_test.go → internal/manager/manager.go

## Import Cycles
- None detected.

## Communities (96 total, 57 thin omitted)

### Community 0 - "app.go"
Cohesion: 0.13
Nodes (8): sameFolder(), TestOneGlobalHuntAtATimeSkippingOurOwnFolders(), jobAccounting, Folder, Progress, WatchResult, dataFor(), TestDataFor()

### Community 1 - "icon/icon.go"
Cohesion: 0.20
Nodes (19): group, Copy(), groupCount(), groupNames(), load(), loadICO(), loadPE(), PNG() (+11 more)

### Community 2 - "testing.T"
Cohesion: 0.14
Nodes (26): TestAddFolderAndSession(), TestListShowsIconAndChosenIcon(), TestListWithoutRootIsEmpty(), TestVersionMatchesWailsJSON(), buildLauncher(), copyFile(), TestLauncherRunsPackageEndToEnd(), TestLauncherStartsTheAppAtItsDefaultLocation() (+18 more)

### Community 3 - "GameFix.jsx"
Cohesion: 0.16
Nodes (26): copyPath(), folderName(), formatSize(), GameFixTab(), Icon(), IconButton(), ICONS, PathField() (+18 more)

### Community 4 - "links_test.go"
Cohesion: 0.07
Nodes (47): Folder, LinkState, isReparse(), LinkStatus(), makeJunction(), MakeLink(), RemoveLink(), mustRead() (+39 more)

### Community 5 - "CreatePackage"
Cohesion: 0.10
Nodes (33): TestHuntAddedAndAddForAPortableApp(), TestRemoveAppDeletesPackageEvenWithBrokenManifest(), CheckLink(), ExpandVars(), Portable(), protected(), SamePath(), TestCheckLinkRefusesSystemAndUserFolders() (+25 more)

### Community 6 - "context.Context"
Cohesion: 0.32
Nodes (5): cleanLeftovers(), complete(), lock(), pkgEntry, Updater

### Community 8 - "runtime/package.json"
Cohesion: 0.12
Nodes (15): author, bugs, url, description, homepage, keywords, license, main (+7 more)

### Community 10 - "wails.json"
Cohesion: 0.12
Nodes (16): author, email, name, frontend:build, frontend:dev:serverUrl, frontend:dev:watcher, frontend:install, info (+8 more)

### Community 11 - "IsUNC"
Cohesion: 0.50
Nodes (4): hasPrefixFold(), IsUNC(), normPath(), TestIsUNC()

### Community 12 - "runtime.js"
Cohesion: 0.22
Nodes (7): EventsOff(), LogWarning(), Quit(), WindowGetPosition(), WindowIsMinimised(), WindowMaximise(), WindowSetMinSize()

### Community 13 - "runtime.d.ts"
Cohesion: 0.25
Nodes (7): EnvironmentInfo, NotificationAction, NotificationCategory, NotificationOptions, Position, Screen, Size

### Community 14 - "models.ts"
Cohesion: 0.09
Nodes (9): Folder, Game, gamefix, main, manager, PackageView, Program, UpdateInfo (+1 more)

### Community 15 - "Diskless PathFix"
Cohesion: 0.22
Nodes (8): Build, Diskless PathFix, Diskless: where updates stick, GameFix tab, Hunt tab: find the folders a program uses, Releases and the updater, The manager (RapTech-Portable.exe), What the launcher does

### Community 17 - "manager/icon_test.go"
Cohesion: 0.33
Nodes (15): ApplyIcon(), FixIcons(), HasOwnIcon(), SetIcon(), copyTo(), realLauncher(), sameIcon(), TestApplyIconAfterRobloxDownload() (+7 more)

### Community 22 - "CreateRoblox"
Cohesion: 0.25
Nodes (9): TestListShowsPackageWithCCBootLine(), TestPublicTxtHasEveryLine(), TestPublicTxtWithoutGameFixBat(), CreateRoblox(), hasLink(), sessionHas(), TestCreateRefusesExistingPackage(), TestCreateRobloxMachineWide() (+1 more)

### Community 23 - "selfupdate_test.go"
Cohesion: 0.19
Nodes (8): Swap(), fakeGitHub(), TestDownloadChecksTheDigest(), TestLatestHTTPError(), TestLatestPicksTheExe(), TestLatestWithoutTheExeFails(), TestSwapPutsTheOldBackWhenTheNewIsMissing(), TestSwapReplacesTheExeAndKeepsTheOld()

### Community 26 - "Hunter"
Cohesion: 0.06
Nodes (32): within(), find, Folder, process, State, activeProcesses(), TestBatsSkipUserFoldersAsSystem(), Hunter (+24 more)

### Community 33 - "Latest"
Cohesion: 0.22
Nodes (9): restartServiceOf(), UpdateInfo, Download(), get(), Latest(), Newer(), parse(), TestNewer() (+1 more)

### Community 35 - "frontend/package.json"
Cohesion: 0.10
Nodes (20): dependencies, react, react-dom, devDependencies, @types/react, @types/react-dom, vite, @vitejs/plugin-react (+12 more)

### Community 36 - "gameFolders"
Cohesion: 0.24
Nodes (12): TestAddHuntFoldersNamesThemApart(), TestGameSourcesListsWhatAGameHas(), gameFolders(), portable(), makeSource(), TestAddCopyRemoveFolders(), TestFolderNamedLikeGame(), isDir() (+4 more)

### Community 38 - "App"
Cohesion: 0.09
Nodes (8): lastLines(), load(), App, SetSession(), TestSetSessionSavesSubfolders(), Cleanup(), TestCleanupRemovesLeftovers(), App

### Community 39 - "update_test.go"
Cohesion: 0.29
Nodes (11): assertOnly(), fakeRoblox(), read(), standardPkgs(), TestPruneKeepsNewestAndCurrent(), TestUpdateInstallsLatestVersion(), TestUpdateRefusesBadDownloads(), TestUpdateReplacesBrokenVersionFolder() (+3 more)

### Community 50 - "TestRemovePackageRemovesProtocolAndSchedule"
Cohesion: 0.44
Nodes (7): TestRemovePackageRemovesProtocolAndSchedule(), RemoveSchedule(), ScheduleOn(), schtasks(), SetSchedule(), TestScheduleCreateQueryRemove(), TestScheduleRefusesOddTimes()

### Community 52 - "App"
Cohesion: 0.05
Nodes (34): NewApp(), App, FolderSettings, Game, Settings, watch, watchService, checkSource() (+26 more)

### Community 54 - ".PublicTxt"
Cohesion: 0.40
Nodes (4): ccbootLine(), isFile(), RunLine(), TestRunLine()

### Community 56 - "Guide.jsx"
Cohesion: 0.53
Nodes (5): Guide(), HelpIcon(), Icon(), store, XIcon()

### Community 64 - "RapTech Portable (shown as "Diskless PathFix")"
Cohesion: 0.40
Nodes (4): Build, graphify, RapTech Portable (shown as "Diskless PathFix"), Rules

### Community 72 - ".ResolveExe"
Cohesion: 0.25
Nodes (6): Created(), hiddenPart(), TestResolveExeMissing(), TestResolveExePicksNewestVersionFolder(), TestStartPathGoesThroughTheDefaultLocation(), Prune()

### Community 73 - "Run"
Cohesion: 0.40
Nodes (5): Result, Run(), TestRunMissingExe(), TestRunTimeoutEndsTheTree(), TestRunWaitsForGrandchildren()

### Community 74 - "linkbat.go"
Cohesion: 0.17
Nodes (20): batLink, callLink(), cleanLink(), inUserFolder(), makeBat(), makeCopyBat(), protected(), removeOldBat() (+12 more)

### Community 76 - "CheckData"
Cohesion: 0.40
Nodes (4): CheckData(), TestCheckDataStaysInsidePackage(), extract(), writeEntry()

### Community 78 - "Load"
Cohesion: 0.13
Nodes (15): Manifest, PackageView, Update, Load(), TestLoadReadsManifest(), TestLoadRefusesUnsafeManifests(), TestSaveThenLoadKeepsManifest(), writeManifest() (+7 more)

### Community 81 - "App.jsx"
Cohesion: 0.22
Nodes (18): AppIcon(), FolderIcon(), Icon(), KeyIcon(), lastLine(), LinkIcon(), nearest(), NewApp() (+10 more)

### Community 88 - "App"
Cohesion: 0.22
Nodes (12): App(), container, root, Icon(), Update(), UpdateButton(), UpdateIcon(), useUpdate() (+4 more)

### Community 90 - "Package"
Cohesion: 0.23
Nodes (12): Package, closePrograms(), Use, imagePath(), InUse(), RemovePackage(), removeServices(), serviceExe() (+4 more)

### Community 92 - "main"
Cohesion: 0.18
Nodes (14): fail(), main(), openLog(), parseArgs(), run(), setup(), TestParseArgs(), update() (+6 more)

### Community 93 - "Hunt.jsx"
Cohesion: 0.32
Nodes (9): Brand(), Steps(), HuntTab(), Icon(), norm(), SearchIcon(), XIcon(), GetVersion() (+1 more)

### Community 94 - "PublicTxt.jsx"
Cohesion: 0.57
Nodes (6): copyText(), FileIcon(), Icon(), PublicTxt(), XIcon(), ClipboardSetText()

## Knowledge Gaps
- **3 isolated node(s):** `@types/react`, `@types/react-dom`, `raptech-portable`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 155 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **57 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `App` connect `App` to `app.go`, `Latest`, `Hunter`, `CreatePackage`, `context.Context`, `Load`, `TestRemovePackageRemovesProtocolAndSchedule`, `App`, `.PublicTxt`, `Package`?**
  _High betweenness centrality (0.037) - this node is a cross-community bridge._
- **What connects `@types/react`, `@types/react-dom`, `raptech-portable` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `app.go` be split into smaller, more focused modules?**
  _Cohesion score 0.13283395755305868 - nodes in this community are weakly interconnected._
- **Why does `App` connect `App` to `app.go`, `context.Context`, `linkbat.go`, `.PublicTxt`, `Hunter`?**
  _High betweenness centrality (0.032) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.1354679802955665 - nodes in this community are weakly interconnected._
- **Why does `Hunter` connect `Hunter` to `app.go`, `.ResolveExe`, `App`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **Should `links_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.0701344243132671 - nodes in this community are weakly interconnected._