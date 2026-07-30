# Features

> Honest feature inventory by status. Verified against code on 2026-07-30.

## Status legend

| Status               | Meaning                                                              |
| -------------------- | -------------------------------------------------------------------- |
| FULLY_FUNCTIONAL     | Code present AND working (tests pass or exercised).                  |
| PARTIALLY_FUNCTIONAL | Ships but has known gaps, edge-case bugs, or missing pieces.         |
| PLANNED              | Designed or documented but **no code exists yet**.                   |
| BROKEN               | Code exists but does not work / is disabled / fails.                 |

## Rules

All 9 rules are registered in `AllRules()` (`rules.go:27`) and `allRuleDetectors()`
(`rules.go:158`), and share a single per-function entry point `DetectFuncDecl()`.

| Rule | Name                 | Status           | Detects                                                           | Suggests                                  |
| ---- | -------------------- | ---------------- | ---------------------------------------------------------------- | ----------------------------------------- |
| H001 | manual-bytes-format  | FULLY_FUNCTIONAL | KMGTPE index trick, unit slices, unit strings + div1024           | `humanize.Bytes` / `humanize.IBytes`      |
| H002 | manual-comma-format  | FULLY_FUNCTIONAL | mod-3, step-by-3, digit-conversion fallback                       | `humanize.Comma`                          |
| H003 | manual-reltime-format| FULLY_FUNCTIONAL | time-diff + "ago" + thresholds                                    | `humanize.RelTime` / `humanize.Time`      |
| H004 | manual-plural        | FULLY_FUNCTIONAL | `if x == 1` with string branch + string-return-type filters       | `humanize.Plural` / `humanize.PluralWord` |
| H005 | manual-si-format     | FULLY_FUNCTIONAL | division by 1000 + K/M suffix, excludes byte units               | `humanize.SI`                             |
| H006 | manual-ftoa          | FULLY_FUNCTIONAL | nested `strings.TrimRight(strings.TrimRight(x,"0"),".")`         | `humanize.Ftoa`                           |
| H007 | manual-parse-bytes   | FULLY_FUNCTIONAL | 2+ HasSuffix/CutSuffix on byte units or map multiplier           | `humanize.ParseBytes`                     |
| H008 | manual-ordinal       | FULLY_FUNCTIONAL | `switch n%10`/`n%100` with st/nd/rd/th cases                      | `humanize.Ordinal`                        |
| H009 | manual-commaf        | FULLY_FUNCTIONAL | `%.Nf` Sprintf + manual comma/separator grouping loop            | `humanize.Commaf`                         |

## Interfaces

| Feature              | Status               | Notes                                                                                |
| -------------------- | -------------------- | ------------------------------------------------------------------------------------ |
| CLI binary           | FULLY_FUNCTIONAL     | `--enable`, `--disable`, `--format text\|json\|sarif`, `--quiet`, `--rules`, `--version`, `--list-files`, `--explain` (`cmd/go-humanize-linter/main.go`) |
| Go library           | FULLY_FUNCTIONAL     | `DefaultRegistry()`, `AllRules()`, `DetectFuncDecl()`, `HumanizeDetector` facade with `Run` / `RunOverPackage` (`rules.go`) |
| golangci-lint plugin | PARTIALLY_FUNCTIONAL | `plugin/plugin.go` works (93.8% coverage) but reports at the func-decl position, not per-line. Rule enable/disable from `.golangci.yml` is not yet supported (see TODO T11). |
| Nix flake            | FULLY_FUNCTIONAL     | `test`, `test-race`, `bench`, `build`, `vet`, `lint`, `coverage` apps                |
| CI workflow          | FULLY_FUNCTIONAL     | test + vet + coverage (Codecov) job and golangci-lint job (`.github/workflows/ci.yml`) |
| Release workflow     | FULLY_FUNCTIONAL     | tagged-build artefacts via `.github/workflows/release.yml` (no `v0.2.0` tag yet — TODO T6) |

## Detection capabilities

| Feature                          | Status           | Notes                                                                |
| -------------------------------- | ---------------- | -------------------------------------------------------------------- |
| Function-scope AST scanning      | FULLY_FUNCTIONAL | All `FuncDecl`s in non-test, non-generated `.go files               |
| Multi-signal detection           | FULLY_FUNCTIONAL | Every rule requires 2+ corroborating signals                         |
| Word-boundary regex              | FULLY_FUNCTIONAL | Prevents false positives like "MEDIUMBLOB" matching "MB"             |
| Underscore digit normalization   | FULLY_FUNCTIONAL | `normLit()` strips `1_000_000` → `1000000`                           |
| Named constant detection         | FULLY_FUNCTIONAL | `hasConst1024` finds `const unit = 1024` patterns                    |
| Generated file skipping          | FULLY_FUNCTIONAL | `_gen.go`, `.gen.go`, `_templ.go` (`plugin/plugin.go:80`)           |
| `//nolint:gohumanize` directives | FULLY_FUNCTIONAL | Bare, `:all`, scoped `:H001`, comma-lists, and `//lint:ignore` syntax (`pattern_helpers.go`) |
| Package-level `var` detection    | PLANNED          | Only `FuncDecl` scope is scanned (TODO T13)                          |
| go/types integration             | PLANNED          | Pure syntactic analysis, no type info (TODO T14)                     |

## Validation

- H001–H007 swept against 190+ Go projects in `~/projects/` (97 findings across ~30 projects).
- ~0% false positive rate on H001–H006; H004 tuned from ~60% FP to ~0% FP across two iterations.
- H008 and H009 have **not** been swept against the corpus yet (TODO T7).
- Full sweep results: `docs/validation/2026-07-30_real-world-sweep.md`.

## Test coverage

Computed via `go test ./... -cover` on 2026-07-30:

| Package                  | Coverage |
| ------------------------ | -------- |
| `go-humanize-linter` (core) | 87.8% |
| `cmd/go-humanize-linter` (CLI) | 35.8% |
| `cmd/gohumanize` (singlechecker) | 0.0% (1-liner `singlechecker.Main` wrapper) |
| `plugin` | 93.8% |
