# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2026-07-30

### Added

- **H008** (manual-ordinal): Detects `switch n%10/100` with st/nd/rd/th cases. Suggests `humanize.Ordinal`.
- **H009** (manual-commaf): Detects `%.Nf` + manual separator grouping. Suggests `humanize.Commaf`.
- Scoped `//nolint:gohumanize:H001` directives — each rule can be suppressed individually.
- Go-style `//lint:ignore gohumanize` alternative syntax.
- `HumanizeDetector` facade with `Run` and `RunOverPackage` methods.
- `--list-files <dir>` CLI flag for debugging walker scope.
- `--explain Hxxx` CLI flag — prints a one-paragraph rationale for any rule.
- `TestSinglechecker_CleanCode` + `TestSinglechecker_PositiveFinding` — closes 0% coverage gap on the singlechecker binary.
- Analysistest fixtures for H002-H009 (was: H001 only).
- Markdown docs per rule: `docs/rules/H001.md` through `H007.md`.
- Release workflow at `.github/workflows/release.yml` (tagged builds).
- CI coverage reporting via Codecov.
- `--version` stderr warning when built without ldflags (dev builds).

### Changed

- `plugin.run` → `plugin.analyzeHumanize` (better grep-ability).
- `printRules()` now writes to stdout (was: stderr) so it pipes cleanly.
- Refactored suppression parser to support both `//nolint:` and `//lint:ignore` flavours with colon-scoped rule IDs.
- `gofumpt` + `goimports -local github.com/larsartmann/` formatting applied.
- `TestRuleBytes_SuppressedByDirective` → `TestCLI_SuppressedByDirective` (clarity).
- Coverage: plugin 93.8% (unchanged), core 87.9% (new detection code has uncovered paths).

## [0.1.0] - 2026-07-30

### Added

- **H001** (manual-bytes-format): Detects hand-rolled byte-size formatting — KMGTPE index trick, []string unit slices, 2+ byte-unit strings with division by 1024. Suggests `humanize.Bytes` / `humanize.IBytes`.
- **H002** (manual-comma-format): Detects thousands-separator insertion loops (mod-3, step-by-3, digit-conversion fallback). Suggests `humanize.Comma`.
- **H003** (manual-reltime-format): Detects relative-time formatting with `time.Since`/`time.Sub` + "ago" strings + time thresholds. Suggests `humanize.RelTime` / `humanize.Time`.
- **H004** (manual-plural): Detects English pluralization via `if x == 1` switches and singular/plural parameter pairs. Suggests `humanize.Plural` / `humanize.PluralWord`. Two-stage false-positive filter: (1) branch must contain string literals, (2) function must return string type.
- **H005** (manual-si-format): Detects SI-prefix formatting (division by 1000/1e6 + K/M suffix). Suggests `humanize.SI`.
- **H006** (manual-ftoa): Detects nested `strings.TrimRight(strings.TrimRight(x, "0"), ".")` for trailing-zero stripping. Suggests `humanize.Ftoa`.
- **H007** (manual-parse-bytes): Detects hand-rolled byte-size string parsing (multiple HasSuffix/CutSuffix on byte units, or map[string]int64 multiplier lookups). Suggests `humanize.ParseBytes`.
- CLI binary (`cmd/go-humanize-linter/`) with `--enable`, `--disable`, `--format text|json|sarif`, `--quiet` flags.
- golangci-lint plugin wrapper (`plugin/`) exposing all rules as `analysis.Analyzer` named `gohumanize`.
- Standalone singlechecker (`cmd/gohumanize/`) for testing the plugin path without golangci-lint.
- `DetectFuncDecl()` shared per-function entry point used by both CLI and plugin.
- `flake.nix` with test, test-race, bench, build, vet, lint, and coverage apps.
- CI workflow with test+vet and lint jobs.
- Negative testdata for H001, H003, H004, and H005.
- 26 library tests + 3 CLI integration tests + 3 plugin tests, all passing with race detector.

### Changed

- H001 message now varies by trigger: KMGTPE index trick, unit string slice, or unit-string count.
- patterns.go split into 8 focused files: pattern_bytes.go, pattern_comma.go, pattern_time.go, pattern_plural.go, pattern_si.go, pattern_ftoa.go, pattern_parsebytes.go, pattern_helpers.go.

### Removed

- Temporary `replace` directives in go.mod (resolved after go-linter-sdk v0.1.0 tag).
