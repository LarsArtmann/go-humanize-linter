# TODO List

> Short- and mid-term improvement tasks, sorted by priority.

## Critical

- [ ] Add `//nolint:gohumanize` directive support (users need to suppress FPs)
- [ ] Self-exclusion: don't flag own source code (`rule_bytes.go` triggers H001)

## High impact

- [ ] analysistest integration test for plugin (plugin.run() at 0% coverage)
- [ ] Per-line diagnostics: detectors should report at the specific pattern location, not just func-decl
- [ ] Add `--version` flag to CLI
- [ ] Add `--rules` flag (list all rules with descriptions)
- [ ] Save validation sweep results to `docs/validation/`

## Medium impact

- [ ] Add example_test.go with runnable Examples
- [ ] Unit tests for CLI `buildRegistry()`, `output()`, `stringList`
- [ ] Table-driven tests for `hasConst1024` (68%), `hasStepBy3` (79%)
- [ ] Reach 80%+ total coverage (currently 74.1%)
- [ ] Add configurable rules in plugin mode (flags or config)
- [ ] SARIF output integration test
- [ ] GitHub Action for running linter in CI

## Low impact / future

- [ ] Package-level var detection for H007
- [ ] Add humanize.Ordinal rule (H008)
- [ ] Add humanize.Commaf rule variant (H009)
- [ ] Consider go/types for type-aware detection
- [ ] Add `--config` flag for YAML/TOML rule configuration
- [ ] Publish to golangci-lint plugin index
