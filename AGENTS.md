# AGENTS.md — go-humanize-linter

## Overview

AST-based linter detecting hand-rolled reimplementations of `dustin/go-humanize`. Built on `go-linter-sdk` + `go-finding`.

## Architecture

| File | Responsibility |
|------|---------------|
| `walker.go` | Directory walking + Go file parsing (`WalkGoDir`, `checkFuncDecls`) |
| `patterns.go` | Shared AST pattern detection helpers (all `has*` / `count*` functions) |
| `rules.go` | `DefaultRegistry()` + `AllRules()` |
| `rule_bytes.go` | H001 — manual byte-size formatting |
| `rule_comma.go` | H002 — manual comma/thousands separator |
| `rule_reltime.go` | H003 — manual relative time |
| `rule_plural.go` | H004 — manual pluralization |
| `rule_si.go` | H005 — manual SI prefix (K/M) |
| `rule_ftoa.go` | H006 — manual float trailing-zero stripping |
| `doc.go` | Package documentation |

## Rule IDs

H001–H006, stable identifiers for suppression matching and filter config.

## Detection Philosophy

Each rule requires **multiple corroborating signals** in the same function:
- H001: byte-unit strings + division by 1024, OR "KMGTPE" index, OR unit slice
- H002: grouping signal (mod-3 or step-by-3 or digit-conversion) + separator writing
- H003: time-difference computation + "ago" string + time threshold comparison
- H004: `if x == 1` on simple identifier (not `len(x)`) OR singular/plural params
- H005: division by power of 1000 + K/M/G/T suffix
- H006: nested `strings.TrimRight(strings.TrimRight(x, "0"), ".")`

## Critical: GOEXPERIMENT=jsonv2 REQUIRED

Depends on go-finding which uses `encoding/json/v2`. All `go` commands need:
- `GOEXPERIMENT=jsonv2`
- `GOPRIVATE=github.com/larsartmann/*`

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

## Testdata

`testdata/` contains positive and negative test fixtures per rule. The walker skips `testdata/` during real scans (see `skipDirs` in `walker.go`).

## Gotchas

- **`normLit` strips underscores** — Go allows `1_000_000` digit separators. All literal comparisons use `normLit()` to normalize.
- **`byteUnitRegex` uses word boundaries** — prevents false positives like "MEDIUMBLOB" matching "MB".
- **H004 requires simple identifier** — `len(x) == 1` is excluded because it's a slice-length check, not pluralization.
- **H002 fallback path** — when step size is a named constant (not literal 3), the fallback uses for-loop + comma + digit-conversion as a combined signal.
- **go.work for local dev** — uses `replace` directives in go.mod pointing to sibling directories. Not committed once tags are published.
