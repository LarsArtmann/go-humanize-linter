# AGENTS.md — go-humanize-linter

## Overview

AST-based linter detecting hand-rolled reimplementations of `dustin/go-humanize`. Built on `go-linter-sdk` + `go-finding`. Ships as a CLI, a Go library, and a golangci-lint plugin.

## Architecture

| File                      | Responsibility                                                                                            |
| ------------------------- | --------------------------------------------------------------------------------------------------------- |
| `walker.go`               | Directory walking + Go file parsing (`WalkGoDir`, `checkFuncDecls`, `posOf`)                              |
| `pattern_helpers.go`      | Shared AST utilities (`unquoteString`, `getBasicLit`, `isPackageCall`, `makeFindingWithConfidence`, etc.) |
| `pattern_bytes.go`        | H001 AST helpers (byte units, KMGTPE, 1024 division, const detection)                                     |
| `pattern_comma.go`        | H002 AST helpers (modulo-3, step-by-3, comma separator writing)                                           |
| `pattern_time.go`         | H003 AST helpers (time.Since/Sub, time threshold comparison)                                              |
| `pattern_plural.go`       | H004 AST helpers (`== 1` branch, string-in-branch check, plural params)                                   |
| `pattern_si.go`           | H005 AST helpers (division by 1000, K/M suffix detection)                                                 |
| `pattern_ftoa.go`         | H006 AST helpers (nested TrimRight detection)                                                             |
| `pattern_parsebytes.go`   | H007 AST helpers (byte-unit suffix checks, multiplier maps)                                               |
| `pattern_ordinal.go`      | H008 AST helpers (`switch n%10/100` with st/nd/rd/th cases)                                               |
| `pattern_commaf.go`       | H009 AST helpers (`%.Nf` Sprintf + manual separator loop)                                                 |
| `rules.go`                | `DefaultRegistry()`, `AllRules()`, `DetectFuncDecl()` (shared per-fn entry point)                         |
| `rule_bytes.go`           | H001 — manual byte-size formatting                                                                        |
| `rule_comma.go`           | H002 — manual comma/thousands separator                                                                   |
| `rule_reltime.go`         | H003 — manual relative time                                                                               |
| `rule_plural.go`          | H004 — manual pluralization                                                                               |
| `rule_si.go`              | H005 — manual SI prefix (K/M)                                                                             |
| `rule_ftoa.go`            | H006 — manual float trailing-zero stripping                                                               |
| `rule_parsebytes.go`      | H007 — manual byte-size string parsing                                                                    |
| `rule_ordinal.go`         | H008 — manual ordinal formatting                                                                          |
| `rule_commaf.go`          | H009 — manual float-with-comma formatting                                                                 |
| `doc.go`                  | Package documentation                                                                                     |
| `plugin/plugin.go`        | golangci-lint plugin wrapper (`analysis.Analyzer` named `gohumanize`)                                     |
| `cmd/go-humanize-linter/` | CLI binary with `--enable`, `--disable`, `--format text\|json\|sarif`, `--quiet`, `--rules`, `--version`, `--list-files`, `--explain`  |
| `cmd/gohumanize/`         | singlechecker entry point for standalone plugin testing                                                   |

## Rule IDs

H001–H009, stable identifiers for suppression matching and filter config.

## Detection Philosophy

Each rule requires **multiple corroborating signals** in the same function:

- **H001**: byte-unit strings + division by 1024, OR "KMGTPE" index, OR unit slice, OR 3+ unit strings
- **H002**: grouping signal (mod-3 or step-by-3 or digit-conversion) + separator writing
- **H003**: time-difference computation + "ago" string + time threshold comparison
- **H004**: `if x == 1` on simple identifier (not `len(x)`) where branch contains strings AND function returns string type; OR singular/plural params
- **H005**: division by power of 1000 + K/M/G/T suffix (excludes byte units)
- **H006**: nested `strings.TrimRight(strings.TrimRight(x, "0"), ".")`
- **H007**: 2+ HasSuffix/CutSuffix/TrimSuffix on byte units OR map[string]int64 multiplier with byte-unit keys
- **H008**: `switch n%10` (or n%100) with at least 3 of 4 st/nd/rd/th case returns
- **H009**: `fmt.Sprintf("%.Nf", x)` (or `strconv.FormatFloat(x, 'f', prec, bits)` with non-zero precision) AND a manual comma/separator group loop in the same function

### H004 false-positive filters (added session 2)

H004 was originally firing on any `if x == 1` with a simple identifier — ~60% false positive rate. Two filters were added:

1. **`branchContainsString(body)`** — the `== 1` branch body must contain a string literal or string concatenation
2. **`funcReturnsString(fn)`** — the function must return `string` or `[]string` type

Together: 44 → 24 findings, ~0% FP.

## Critical: GOEXPERIMENT=jsonv2 REQUIRED

Depends on go-finding which uses `encoding/json/v2`. All `go` commands need:

- `GOEXPERIMENT=jsonv2`
- `GOPRIVATE=github.com/larsartmann/*`
- `GONOSUMDB=github.com/larsartmann/*`

Use `nix run .#test` / `nix run .#lint` (sets env automatically).

## Build Commands

```bash
nix run .#test          # go test ./... -count=1
nix run .#test-race     # go test ./... -race -count=1
nix run .#build         # go build ./...
nix run .#lint          # golangci-lint run ./...
nix run .#vet           # go vet ./...
nix run .#coverage      # go test with coverage report
```

Direct Go commands:

```bash
export GOEXPERIMENT=jsonv2 GOPRIVATE='github.com/larsartmann/*' GONOSUMDB='github.com/larsartmann/*'
go build ./...
go test ./... -race -count=1
go vet ./...
```

## Testdata

`testdata/` contains positive and negative test fixtures per rule. The walker skips `testdata/` during real scans (see `skipDirs` in `walker.go`).

## Gotchas

- **`normLit` strips underscores** — Go allows `1_000_000` digit separators. All literal comparisons use `normLit()` to normalize.
- **`byteUnitRegex` uses word boundaries** — prevents false positives like "MEDIUMBLOB" matching "MB".
- **H004 requires string return type** — functions returning `error`, `bool`, or `int` are excluded even if they have `if x == 1` with string literals.
- **H002 fallback path** — when step size is a named constant (not literal 3), the fallback uses for-loop + comma + digit-conversion as a combined signal.
- **H001 message varies by trigger** — KMGTPE index trick, unit slice, and unit-string-count each produce different messages for clarity.
- **Plugin uses `fn.Pos()`** — findings are reported at the function declaration position, consistent with the CLI. Per-line diagnostics would require changing every detector to return specific `token.Pos` values.
- **Self-detection (resolved)** — the linter's own `rule_bytes.go` would flag itself as H001 (suggestion text names byte units). It is now suppressed with a `//nolint:gohumanize` directive on `detectBytesFormat`. Run the CLI on the repo root to confirm 0 self-findings.
- **Suppression directives** — `//nolint:gohumanize` (also `//nolint`, `//nolint:all`, comma-lists) suppresses findings on a function in **both** the CLI path (`checkFuncDecls` in `walker.go`) and the plugin path (`DetectFuncDecl` in `rules.go`). The matching logic lives in `hasNoLintDirective` / `noLintMatches` in `pattern_helpers.go`.
- **Typed errors** — `WalkGoDir` and `checkFuncDecls` return `*WalkError{Dir, Err}` on walk failure; `output()` returns `*OutputError{Format, Stage, Err}`. Callers extract typed info with `errors.AsType[*WalkError](err)` / `errors.AsType[*OutputError](err)`. `Format`/`Stage` are string consts (`formatText`/`formatJSON`/`formatSARIF` and `stageRender`/`stageWrite`) — single source of truth shared by flag default, switch cases, and error fields.
- **JSON write error was swallowed (fixed)** — `output()`'s json branch previously ignored `fmt.Fprintln`'s error. Failure-path test (`TestOutput_JSONWriterFailure`) now exercises this path using `failingWriter`. The test caught the bug.
- **erraudit (historical)** — erraudit flagged out-of-scope variables (`fset`/`base`/`parseErr`/`detect`/`files`) as missing context on error wraps inside `walker.go`. The tool itself is structurally flawed for this codebase (pattern-matches identifiers without scope checking). The `//nolint:erraudit` directives in `walker.go` have been removed; the `nolint_filter` warning is filtered from `nix run .#lint` output (golanci-lint can't validate a directive for a non-registered linter). The `WalkError`/`OutputError` types satisfy `generic_return` for two of three functions; `WalkGoDir`'s signature stays `error` for v0.1.x backward compatibility.
- **Scope suppression** — `//nolint:gohumanize:H001` (per-rule) works alongside the unscoped form. Range syntax `H001-H009` is NOT supported (use individual IDs). See `ruleIDH001`–`ruleIDH009` in `pattern_helpers.go` for canonical IDs.
- **H008 / H009 thunks** — both detectors (`pattern_ordinal.go`, `pattern_commaf.go`) extract helpers to keep cyclomatic complexity under 12. The format-string walker in H009 (`walkFormatFloatVerbs`) is intentionally split rather than left as one monolithic loop so the `%.Nf` and `%<digit>f` code paths can be tested in isolation.
- **`go.mod` replace directives are local-dev only** — `go.mod` points `go-finding`, `go-linter-sdk`, and `go-error-family` at `../go-*` sibling repos. These are a development convenience, not part of any release. Publishing to the golangci-lint plugin index is blocked until `go-linter-sdk` has a tag and the replaces can be removed.
- **Analysistest gap** — `TestAnalyzerAnalysistest` (`plugin/plugin_test.go`) runs H001–H007 + clean only, even though `h008positive`/`h009positive` fixtures exist. The two new rules are not yet wired into the run call.
