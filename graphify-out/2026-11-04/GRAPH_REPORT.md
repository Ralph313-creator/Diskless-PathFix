# Graph Report - RapTech-Portable  (2026-11-04)

## Corpus Check
- Corpus is ~10,350 words - fits in a single context window. You may not need a graph.

## Summary
- 189 nodes · 549 edges · 10 communities (9 shown, 1 thin omitted)
- Extraction: 85% EXTRACTED · 15% INFERRED · 0% AMBIGUOUS · INFERRED: 83 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- Launcher Windows Helpers
- Launcher Main Flow
- Roblox Updater Tests
- Registry Save And Load
- Link Tests
- Package Setup And Links
- Roblox Updater
- Path Rules
- Launcher Tests
- Module Root

## God Nodes (most connected - your core abstractions)
1. `MakeLink()` - 12 edges
2. `mustWrite()` - 12 edges
3. `SaveRegistry()` - 12 edges
4. `Updater` - 12 edges
5. `main()` - 11 edges
6. `LinkStatus()` - 11 edges
7. `mustRead()` - 11 edges
8. `Package` - 11 edges
9. `Run()` - 9 edges
10. `update()` - 8 edges

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

## Communities (10 total, 1 thin omitted)

### Community 0 - "Launcher Windows Helpers"
Cohesion: 0.12
Nodes (15): attachConsole(), imagePath(), otherInstances(), runElevated(), Folder, jobAccounting, Manifest, Result (+7 more)

### Community 1 - "Launcher Main Flow"
Cohesion: 0.18
Nodes (17): fail(), main(), openLog(), parseArgs(), run(), setup(), update(), Package (+9 more)

### Community 2 - "Roblox Updater Tests"
Cohesion: 0.15
Nodes (13): Created(), Prune(), assertOnly(), fakeRoblox(), read(), standardPkgs(), TestPruneKeepsNewestAndCurrent(), TestUpdateInstallsLatestVersion() (+5 more)

### Community 3 - "Registry Save And Load"
Cohesion: 0.22
Nodes (17): fromUTF16(), keyExists(), LoadRegistry(), RegFileName(), replaceFold(), rewriteReg(), runReg(), SaveRegistry() (+9 more)

### Community 4 - "Link Tests"
Cohesion: 0.30
Nodes (18): MakeLink(), mustRead(), mustWrite(), TestLinkStatusMissing(), TestMakeLinkCreatesMissingTarget(), TestMakeLinkMovesRealFolderAside(), TestMakeLinkOnMissingPathLinksIntoPackage(), TestMakeLinkRefusesToOverwriteBackup() (+10 more)

### Community 5 - "Package Setup And Links"
Cohesion: 0.14
Nodes (9): FolderStatus, LinkState, isReparse(), LinkStatus(), makeJunction(), RemoveLink(), EmptyFolder(), SetProtocol() (+1 more)

### Community 6 - "Roblox Updater"
Cohesion: 0.20
Nodes (8): cleanLeftovers(), complete(), extract(), lock(), New(), writeEntry(), pkgEntry, Updater

### Community 7 - "Path Rules"
Cohesion: 0.19
Nodes (15): CheckData(), CheckLink(), ExpandVars(), hasPrefixFold(), IsUNC(), normPath(), Portable(), protected() (+7 more)

### Community 8 - "Launcher Tests"
Cohesion: 0.32
Nodes (4): buildLauncher(), copyFile(), TestLauncherRunsPackageEndToEnd(), TestParseArgs()

## Knowledge Gaps
- **1 isolated node(s):** `raptech-portable`
  These have ≤1 connection - possible missing edges. (Counts symbols only; 14 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **1 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Updater` connect `Roblox Updater` to `Launcher Windows Helpers`, `Launcher Main Flow`, `Roblox Updater Tests`?**
  _High betweenness centrality (0.073) - this node is a cross-community bridge._
- **Are the 9 inferred relationships involving `MakeLink()` (e.g. with `.Setup()` and `IsUNC()`) actually correct?**
  _`MakeLink()` has 9 INFERRED edges - model-reasoned connections that need verification._
- **What connects `raptech-portable` to the rest of the system?**
  _1 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Launcher Windows Helpers` be split into smaller, more focused modules?**
  _Cohesion score 0.11764705882352941 - nodes in this community are weakly interconnected._
- **Why does `Package` connect `Launcher Main Flow` to `Launcher Windows Helpers`, `Link Tests`, `Package Setup And Links`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **Are the 4 inferred relationships involving `mustWrite()` (e.g. with `TestFinishSavesRegistryAndEmptiesSession()` and `TestSetupLinksAndLoadsRegistry()`) actually correct?**
  _`mustWrite()` has 4 INFERRED edges - model-reasoned connections that need verification._
- **Should `Package Setup And Links` be split into smaller, more focused modules?**
  _Cohesion score 0.1368421052631579 - nodes in this community are weakly interconnected._