# RapTech Portable

Makes apps portable for CCBoot diskless PCs: an app's folders and registry
keys live in one package folder next to a launcher exe (any drive), and the
app's own updates land there too. Roblox is the first app; it can also be
updated straight from Roblox's servers.

- `internal/core` — paths and safety checks, `portable.json` manifest,
  junction links (real folder in the way -> `.rp-backup`), registry save/load
  with `%PKG%`, protocol handlers, run-and-wait with a job object.
- `internal/roblox` — downloads the newest Roblox Player version (MD5 checked,
  package map like Bloxstrap), prunes old versions.
- `cmd/launcher` — the exe copied into each package and renamed after the app.
- `internal/hunt` — hunt mode, the Hunt tab (one global tool; results go to a portable app or a GameFix game). `internal/gamefix` — the GameFix tab,
  ported from `D:\Project\Diskless-GameFix`; its own settings file
  (`RapTech-GameFix.json`) and service (`RapTechGameFixWatch`), never the
  standalone GameFix's.
- `internal/ccboot` — reads and adds `run` lines to CCBoot's command file
  `C:\CCBoot\Cmd\public.txt` (every client runs it at boot). The GameFix
  tab's footer "Add" puts `run -startup -skipupdate <main>\GameFix.bat` there.

## Rules

- Remote: private GitHub repo Ralph313-creator/RapTech-Portable (origin). Push only
  when asked. No Claude
  co-author line in commits.
- Test first (TDD). `go test ./...` must stay green; registry tests use
  `HKCU\Software\RapPortableTest` and clean up.
- On floor PCs the game disk is write-back too: updates only stick when run
  on the server (or a super client).
- This PC is the live CCBoot server and counter PC: never restart it, CCBoot or
  the NIC while customers are on.

## Build

```
build.cmd   # go test, launcher -> assets/launcher.exe (embedded), wails build
```

The manager (root package, Wails + React in frontend/) is a thin wrapper over
`internal/manager`; it runs as administrator (build/windows/wails.exe.manifest).

## graphify

This project has a knowledge graph at graphify-out/ with god nodes, community structure, and cross-file relationships.

Rules:
- For codebase questions, first run `graphify query "<question>"` when graphify-out/graph.json exists. Use `graphify path "<A>" "<B>"` for relationships and `graphify explain "<concept>"` for focused concepts. These return a scoped subgraph, usually much smaller than GRAPH_REPORT.md or raw grep output.
- Read graphify-out/GRAPH_REPORT.md only for broad architecture review or when query/path/explain do not surface enough context.
- After modifying code, run `graphify update .` to keep the graph current (AST-only, no API cost), then commit it as "Refresh the knowledge graph".
