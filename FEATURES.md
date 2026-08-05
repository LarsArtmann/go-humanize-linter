# Features

> Honest feature inventory by status. Verified against code on 2026-08-05.

## Status legend

| Status               | Meaning                                                      |
| -------------------- | ------------------------------------------------------------ |
| FULLY_FUNCTIONAL     | Code present AND working (tests pass or exercised).          |
| PARTIALLY_FUNCTIONAL | Ships but has known gaps, edge-case bugs, or missing pieces. |
| PLANNED              | Designed or documented but **no code exists yet**.           |
| BROKEN               | Code exists but does not work / is disabled / fails.         |

## Rules

All 9 rules are registered in `AllRules()` (`rules.go`) and `allRuleDetectors()`
(`rules.go`), and share a single per-function entry point `DetectFuncDecl()`.

| Rule | Name                  | Status           | Detects                                                                                                                              | Suggests                                |
| ---- | --------------------- | ---------------- | ------------------------------------------------------------------------------------------------------------------------------------ | --------------------------------------- |
| H001 | manual-bytes-format   | FULLY_FUNCTIONAL | KMGTPE index trick, unit slices, unit strings + div1024. Size-bucket lookup tables (switch/unit-slice without div1024) are excluded. | `humanize.Bytes` / `humanize.IBytes`    |
| H002 | manual-comma-format   | FULLY_FUNCTIONAL | mod-3, step-by-3, digit-conversion fallback                                                                                          | `humanize.Comma`                        |
| H003 | manual-reltime-format | FULLY_FUNCTIONAL | time-diff + "ago" + thresholds                                                                                                       | `humanize.RelTime` / `humanize.Time`    |
| H004 | manual-plural         | FULLY_FUNCTIONAL | `if x == 1` with string branch + string-return-type filters                                                                          | `english.Plural` / `english.PluralWord` |
| H005 | manual-si-format      | FULLY_FUNCTIONAL | division by 1000 + K/M suffix, excludes byte units                                                                                   | `humanize.SI`                           |
| H006 | manual-ftoa           | FULLY_FUNCTIONAL | nested `strings.TrimRight(strings.TrimRight(x,"0"),".")`                                                                             | `humanize.Ftoa`                         |
| H007 | manual-parse-bytes    | FULLY_FUNCTIONAL | 2+ HasSuffix/CutSuffix on byte units, map multiplier (func + package scope), aliased imports                                         | `humanize.ParseBytes`                   |
| H008 | manual-ordinal        | FULLY_FUNCTIONAL | `switch n%10`/`n%100` with st/nd/rd/th cases                                                                                         | `humanize.Ordinal`                      |
| H009 | manual-commaf         | FULLY_FUNCTIONAL | `%.Nf` Sprintf + manual comma/separator grouping loop                                                                                | `humanize.Commaf`                       |

## Interfaces

| Feature              | Status           | Notes                                                                                                                                                                                                                                                                          |
| -------------------- | ---------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| CLI binary           | FULLY_FUNCTIONAL | `--enable`, `--disable`, `--config`, `--format`, `--output`, `--quiet`, `--rules`, `--version`, `--list-files`, `--explain`, `--min-confidence`, `--verify-suppressions` (`cmd/go-humanize-linter/main.go`)                                                                    |
| Go library           | FULLY_FUNCTIONAL | `DefaultRegistry()`, `AllRules()`, `DetectFuncDecl()`, `HumanizeDetector` facade with `Run`, exported `RuleIDH001`–`H009` constants (`rules.go`)                                                                                                                               |
| golangci-lint plugin | FULLY_FUNCTIONAL | `plugin/plugin.go` using `plugin-module-register` v2 module plugin pattern. Configurable enable/disable, `minConfidence`, and `verifySuppressions` via `.golangci.yml` `linters.settings.custom.gohumanize.settings`. Diagnostics at finding position via `findingToTokenPos`. |
| GitHub Action        | FULLY_FUNCTIONAL | `action.yml` composite Action with inputs: path, enable, disable, format, version                                                                                                                                                                                              |
| Nix flake            | FULLY_FUNCTIONAL | `test`, `test-race`, `bench`, `build`, `vet`, `lint`, `coverage` apps                                                                                                                                                                                                          |
| CI workflow          | FULLY_FUNCTIONAL | test + vet + coverage (Codecov) job and golangci-lint job (`.github/workflows/ci.yml`). Versions pinned: `govulncheck@v1.6.0`, `golangci-lint v2.12.2`.                                                                                                                        |
| Release workflow     | FULLY_FUNCTIONAL | tagged-build artefacts via `.github/workflows/release.yml` (no `v0.2.0` tag yet — TODO T1)                                                                                                                                                                                     |

## Detection capabilities

| Feature                          | Status           | Notes                                                                                                                               |
| -------------------------------- | ---------------- | ----------------------------------------------------------------------------------------------------------------------------------- |
| Function-scope AST scanning      | FULLY_FUNCTIONAL | All `FuncDecl`s in non-test, non-generated `.go files                                                                               |
| Package-level var scanning       | FULLY_FUNCTIONAL | H007 detects `var x = map[string]int64{"KB": 1024}` at package scope                                                                |
| Multi-signal detection           | FULLY_FUNCTIONAL | Every rule requires 2+ corroborating signals                                                                                        |
| Word-boundary regex              | FULLY_FUNCTIONAL | Prevents false positives like "MEDIUMBLOB" matching "MB"                                                                            |
| Underscore digit normalization   | FULLY_FUNCTIONAL | `normLit()` strips `1_000_000` → `1000000`                                                                                          |
| Named constant detection         | FULLY_FUNCTIONAL | `hasConst1024` finds `const unit = 1024` patterns                                                                                   |
| Generated file skipping          | FULLY_FUNCTIONAL | `_gen.go`, `.gen.go`, `_templ.go` (`plugin/plugin.go`)                                                                              |
| `//nolint:gohumanize` directives | FULLY_FUNCTIONAL | Bare, `:all`, scoped `:H001`, comma-lists, and `//lint:ignore` syntax (`pattern_helpers.go`)                                        |
| Import-alias resolution          | FULLY_FUNCTIONAL | `buildImportAliases(file)` resolves `str "strings"` → `{"str": "strings"}`. Dot imports (`. "strings"`) are a known gap (TODO T16). |
| Map-type checking (H007)         | FULLY_FUNCTIONAL | `isByteUnitMultiplierMapLiteral` checks `map[string]int*` value type to avoid flagging `map[string]bool` lookup sets                |
| go/types integration             | PLANNED          | Pure syntactic analysis, no type info. ADR 0001 documents the trade-off.                                                            |

## Configuration

| Feature                          | Status           | Notes                                                                                                                      |
| -------------------------------- | ---------------- | -------------------------------------------------------------------------------------------------------------------------- |
| CLI `--enable`/`--disable` flags | FULLY_FUNCTIONAL | Per-rule filtering on the command line                                                                                     |
| CLI `--config` YAML file         | FULLY_FUNCTIONAL | Load enable/disable rules from `.gohumanize.yaml`. Config-first, CLI-override with set-union semantics.                    |
| CLI `--min-confidence` flag      | FULLY_FUNCTIONAL | Filter findings by confidence level: `low`, `medium`, `high`, `full` (default: `low`). Uses `finding.ByConfidenceAtLeast`. |
| CLI `--verify-suppressions` flag | FULLY_FUNCTIONAL | Detects stale `//nolint:gohumanize` directives (suppress zero findings) and misspelled linter names. Reports as `H0SUP`.   |
| CLI confidence-aware exit codes  | FULLY_FUNCTIONAL | Exit 0 = clean, exit 1 = high/full-confidence finding (must fix), exit 2 = only medium/low (triage).                       |
| Plugin enable/disable settings   | FULLY_FUNCTIONAL | `linters.settings.custom.gohumanize.settings.enable`/`disable` in `.golangci.yml`                                          |
| Plugin `minConfidence` setting   | FULLY_FUNCTIONAL | `linters.settings.custom.gohumanize.settings.minConfidence` in `.golangci.yml` — same levels as CLI `--min-confidence`     |
| Plugin `verifySuppressions`      | FULLY_FUNCTIONAL | `linters.settings.custom.gohumanize.settings.verifySuppressions` in `.golangci.yml` — same as CLI `--verify-suppressions`  |
| Scoped `//nolint` directives     | FULLY_FUNCTIONAL | Per-function suppression via `//nolint:gohumanize:H001`                                                                    |

## Validation

- H001–H009 swept against 327 Go projects in `~/projects/` (242 findings across ~80 files).
- ~0% false positive rate on H001–H006; H004 tuned from ~60% FP to ~0% FP across two iterations.
- H008: 0 findings in corpus. H009: 3 findings, all true positives.
- Import-alias-aware detection: 0 real-world hits but verified via testdata.
- Package-level var detection: 1 true positive (`clean-wizard`).
- The `--min-confidence`, `--verify-suppressions`, and `--behavior-delta` features have **not** been swept against the corpus yet (TODO T2).
- Full sweep results: `docs/validation/2026-07-31_real-world-sweep.md`.

## Test coverage

Computed via `go test ./... -cover` on 2026-08-05:

| Package                          | Coverage                                                               |
| -------------------------------- | ---------------------------------------------------------------------- |
| `go-humanize-linter` (core)      | 88.4%                                                                  |
| `cmd/go-humanize-linter` (CLI)   | 38.4%                                                                  |
| `cmd/gohumanize` (singlechecker) | 0.0% (1-liner `singlechecker.Main` wrapper — not coverable in-process) |
| `plugin`                         | 95.7%                                                                  |
