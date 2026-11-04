# Graph Report - RapTech-Portable  (2026-11-04)

## Corpus Check
- 19 files · ~10,928 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 1 file(s) not represented in the graph (top: (none) 1)

## Summary
- 196 nodes · 561 edges · 12 communities (11 shown, 1 thin omitted)
- Extraction: 85% EXTRACTED · 15% INFERRED · 0% AMBIGUOUS · INFERRED: 86 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `717c3d39`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- update.go
- manifest_test.go
- fakeRoblox
- registry.go
- testing.T
- .Path
- Updater
- paths.go
- main_test.go
- raptech-portable
- links.go
- RapTech Portable

## God Nodes (most connected - your core abstractions)
1. `MakeLink()` - 12 edges
2. `mustWrite()` - 12 edges
3. `SaveRegistry()` - 12 edges
4. `Updater` - 12 edges
5. `main()` - 11 edges
6. `LinkStatus()` - 11 edges
7. `mustRead()` - 11 edges
8. `Package` - 11 edges
9. `testPackage()` - 9 edges
10. `Run()` - 9 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `attachConsole()`  [INFERRED]
  cmd/launcher/main.go → cmd/launcher/win.go
- `main()` --calls--> `SaveRegistry()`  [EXTRACTED]
  cmd/launcher/main.go → internal/core/registry.go
- `setup()` --calls--> `runElevated()`  [INFERRED]
  cmd/launcher/main.go → cmd/launcher/win.go
- `run()` --calls--> `Run()`  [EXTRACTED]
  cmd/launcher/main.go → internal/core/run.go
- `update()` --calls--> `SaveRegistry()`  [EXTRACTED]
  cmd/launcher/main.go → internal/core/registry.go

## Import Cycles
- None detected.

## Communities (12 total, 1 thin omitted)

### Community 0 - "update.go"
Cohesion: 0.12
Nodes (12): attachConsole(), imagePath(), otherInstances(), runElevated(), jobAccounting, Result, activeProcesses(), Run() (+4 more)

### Community 1 - "manifest_test.go"
Cohesion: 0.15
Nodes (18): fail(), main(), openLog(), run(), setup(), update(), Manifest, Update (+10 more)

### Community 2 - "fakeRoblox"
Cohesion: 0.21
Nodes (13): Created(), Prune(), assertOnly(), fakeRoblox(), read(), standardPkgs(), TestPruneKeepsNewestAndCurrent(), TestUpdateInstallsLatestVersion() (+5 more)

### Community 3 - "registry.go"
Cohesion: 0.22
Nodes (17): fromUTF16(), keyExists(), LoadRegistry(), RegFileName(), replaceFold(), rewriteReg(), runReg(), SaveRegistry() (+9 more)

### Community 4 - "testing.T"
Cohesion: 0.27
Nodes (21): LinkStatus(), MakeLink(), RemoveLink(), mustRead(), mustWrite(), TestLinkStatusMissing(), TestMakeLinkCreatesMissingTarget(), TestMakeLinkMovesRealFolderAside() (+13 more)

### Community 5 - ".Path"
Cohesion: 0.20
Nodes (6): Folder, FolderStatus, EmptyFolder(), RemoveProtocol(), SetProtocol(), TestSetProtocolPointsAtLauncher()

### Community 6 - "Updater"
Cohesion: 0.22
Nodes (7): cleanLeftovers(), complete(), extract(), lock(), writeEntry(), pkgEntry, Updater

### Community 7 - "paths.go"
Cohesion: 0.18
Nodes (15): CheckData(), CheckLink(), ExpandVars(), hasPrefixFold(), IsUNC(), normPath(), Portable(), protected() (+7 more)

### Community 8 - "main_test.go"
Cohesion: 0.28
Nodes (6): parseArgs(), buildLauncher(), copyFile(), TestLauncherRunsPackageEndToEnd(), TestParseArgs(), options

### Community 10 - "links.go"
Cohesion: 0.22
Nodes (3): LinkState, isReparse(), makeJunction()

### Community 11 - "RapTech Portable"
Cohesion: 0.40
Nodes (4): Build, graphify, RapTech Portable, Rules

## Knowledge Gaps
- **1 isolated node(s):** `raptech-portable`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 18 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **1 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Updater` connect `Updater` to `update.go`, `manifest_test.go`, `fakeRoblox`?**
  _High betweenness centrality (0.068) - this node is a cross-community bridge._
- **Are the 9 inferred relationships involving `MakeLink()` (e.g. with `.Setup()` and `IsUNC()`) actually correct?**
  _`MakeLink()` has 9 INFERRED edges - model-reasoned connections that need verification._
- **What connects `raptech-portable` to the rest of the system?**
  _1 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `update.go` be split into smaller, more focused modules?**
  _Cohesion score 0.11829268292682926 - nodes in this community are weakly interconnected._
- **Why does `Package` connect `manifest_test.go` to `update.go`, `testing.T`, `.Path`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **Are the 4 inferred relationships involving `mustWrite()` (e.g. with `TestFinishSavesRegistryAndEmptiesSession()` and `TestSetupLinksAndLoadsRegistry()`) actually correct?**
  _`mustWrite()` has 4 INFERRED edges - model-reasoned connections that need verification._
- **Why does `fakeRoblox()` connect `fakeRoblox` to `update.go`, `testing.T`, `Updater`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._