# Features

> Honest feature inventory by status.

## Status legend

| Status | Meaning |
| ------ | ------- |
| DONE | Works as intended, exercised by tests and real-world validation |
| PARTIALLY DONE | Ships but has known gaps |
| PLANNED | Designed or documented but not yet implemented |

## Rules

| Rule | Name | Status | Notes |
| ---- | ---- | ------ | ----- |
| H001 | manual-bytes-format | DONE | Detects KMGTPE trick, unit slices, unit strings + div1024. Validated: 27 real findings. |
| H002 | manual-comma-format | DONE | Detects mod-3, step-3, digit-conversion fallback. Validated: 7 real findings. |
| H003 | manual-reltime-format | DONE | Requires time-diff + "ago" + thresholds. Validated: 25 real findings. |
| H004 | manual-plural | DONE | String-in-branch + string-return-type filters. Validated: 24 real findings, ~0% FP. |
| H005 | manual-si-format | DONE | Division by 1000 + K/M suffix, excludes byte units. Validated: 10 real findings. |
| H006 | manual-ftoa | DONE | Nested TrimRight pattern. Validated: 2 real findings. |
| H007 | manual-parse-bytes | DONE | HasSuffix/CutSuffix or map multiplier. Validated: 2 real findings. |

## Interfaces

| Feature | Status | Notes |
| ------- | ------ | ----- |
| CLI binary | DONE | `--enable`, `--disable`, `--format text\|json\|sarif`, `--quiet` |
| Go library | DONE | `DefaultRegistry()`, `AllRules()`, `DetectFuncDecl()` |
| golangci-lint plugin | PARTIALLY DONE | `plugin/plugin.go` works but reports at func-decl position, not per-line |
| Nix flake | DONE | test, test-race, bench, build, vet, lint, coverage |
| CI workflow | DONE | test+vet job, lint job |

## Detection capabilities

| Feature | Status | Notes |
| ------- | ------ | ----- |
| Function-scope AST scanning | DONE | All FuncDecls in non-test, non-generated .go files |
| Multi-signal detection | DONE | Every rule requires 2+ corroborating signals |
| Word-boundary regex | DONE | Prevents false positives like "MEDIUMBLOB" matching "MB" |
| Underscore digit normalization | DONE | `normLit()` strips `1_000_000` → `1000000` |
| Named constant detection | DONE | `hasConst1024` finds `const unit = 1024` patterns |
| Generated file skipping | DONE | `_gen.go`, `.gen.go`, `_templ.go` |
| Package-level var detection | PLANNED | Only FuncDecl scope is scanned currently |
| go/types integration | PLANNED | Pure syntactic analysis, no type info |
| `//nolint:gohumanize` directives | PLANNED | No suppression mechanism yet |

## Validation

- Validated against 190+ Go projects in `~/projects/`
- 97 total findings across ~30 projects
- ~0% false positive rate on 6 of 7 rules
- H004 tuned from ~60% FP to ~0% FP across two iterations
