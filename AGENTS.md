# AGENTS.md — go-humanize-linter

## Overview

AST-based linter detecting hand-rolled reimplementations of `dustin/go-humanize`. Built on `go-linter-sdk` + `go-finding`. Ships as a CLI, a Go library, and a golangci-lint plugin.

## Architecture

| File                  | Responsibility                                                                      |
| --------------------- | ----------------------------------------------------------------------------------- |
| `walker.go`           | Directory walking + Go file parsing (`WalkGoDir`, `checkFuncDecls`, `posOf`)        |
| `patterns.go`         | Shared AST pattern detection helpers (all `has*` / `count*` functions)              |
| `rules.go`            | `DefaultRegistry()`, `AllRules()`, `DetectFuncDecl()` (shared per-fn entry point)   |
| `rule_bytes.go`       | H001 — manual byte-size formatting                                                  |
| `rule_comma.go`       | H002 — manual comma/thousands separator                                             |
| `rule_reltime.go`     | H003 — manual relative time                                                         |
| `rule_plural.go`      | H004 — manual pluralization                                                         |
| `rule_si.go`          | H005 — manual SI prefix (K/M)                                                       |
| `rule_ftoa.go`        | H006 — manual float trailing-zero stripping                                         |
| `rule_parsebytes.go`  | H007 — manual byte-size string parsing                                              |
| `doc.go`              | Package documentation                                                               |
| `plugin/plugin.go`    | golangci-lint plugin wrapper (`analysis.Analyzer` named `gohumanize`)               |
| `cmd/go-humanize-linter/` | CLI binary with `--enable`, `--disable`, `--format text\|json\|sarif`, `--quiet` |
| `cmd/gohumanize/`     | singlechecker entry point for standalone plugin testing                             |

## Rule IDs

H001–H007, stable identifiers for suppression matching and filter config.

## Detection Philosophy

Each rule requires **multiple corroborating signals** in the same function:

- **H001**: byte-unit strings + division by 1024, OR "KMGTPE" index, OR unit slice, OR 3+ unit strings
- **H002**: grouping signal (mod-3 or step-by-3 or digit-conversion) + separator writing
- **H003**: time-difference computation + "ago" string + time threshold comparison
- **H004**: `if x == 1` on simple identifier (not `len(x)`) where branch contains strings AND function returns string type; OR singular/plural params
- **H005**: division by power of 1000 + K/M/G/T suffix (excludes byte units)
- **H006**: nested `strings.TrimRight(strings.TrimRight(x, "0"), ".")`
- **H007**: 2+ HasSuffix/CutSuffix/TrimSuffix on byte units OR map[string]int64 multiplier with byte-unit keys

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
- **Self-detection** — the linter flags its own `rule_bytes.go` as H001 because it contains byte-unit detection patterns. This is expected.
