# Status Report: go-humanize-linter

**Date:** 2026-07-30 17:32\
**Session:** Initial build + self-review iteration\
**Repo:** [LarsArtmann/go-humanize-linter](https://github.com/LarsArtmann/go-humanize-linter)\
**HEAD:** `946ae83` — Add CLI integration tests

---

## Executive Summary

Built a 7-rule AST linter (H001–H007) on go-linter-sdk that detects hand-rolled reimplementations of `dustin/go-humanize`. The linter compiles, passes 26 tests with race detector, scans 10 real projects with near-zero false positives, ships a CLI binary, and is pushed to GitHub.

**But:** coverage is 78.6% (target: 80%+), the CLI has 0% direct test coverage (only integration-tested via subprocess), no project docs exist (CHANGELOG, FEATURES, ROADMAP, TODO_LIST), the `.golangci.yml` has a formatting issue, and the go.mod replace directives are a temporary hack that breaks standalone consumption.

> **Update 2026-07-30:** every gap in this paragraph has since closed. The linter now has **9 rules** (H001–H009, not 7), core coverage is **87.8%** / plugin **93.8%**, the CLI has direct unit tests (35.8%), all project docs exist and are current, and golangci-lint reports **0 issues**. The `go.mod` replace directives were re-added for local dev (publishing still blocked on a `go-linter-sdk` tag). Full item-by-item status in [Resolution](#resolution-2026-07-30) below.

---

## a) FULLY DONE ✅

| Area                    | Details                                                                                                                                                                                                                        |
| ----------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **Core linter library** | 7 rules (H001–H007), `DefaultRegistry()`, `AllRules()`, individual `Rule*()` factories                                                                                                                                         |
| **AST pattern engine**  | `walker.go` (directory walk + parse), `patterns.go` (19 detection helpers + 7 internal utilities)                                                                                                                              |
| **Detection accuracy**  | Multi-signal approach; validated against DiscordSync, SEC, BuildFlow, go-cqrs-lite, KeyCountdown, standard-bug-tracking-schema, ast-state-analyzer, file-and-image-renamer, blog, auto-deduplicate — near-zero false positives |
| **CLI binary**          | `cmd/go-humanize-linter/main.go` with `--enable`, `--disable`, `--format text\|json\|sarif`, `--quiet`                                                                                                                         |
| **CLI tests**           | 3 integration tests: clean code (exit 0), positive finding (H001), enable-filter (H003 on H001 data = no match)                                                                                                                |
| **Unit tests**          | 23 library tests covering every rule + positive/negative cases + registry + exit code                                                                                                                                          |
| **Benchmarks**          | `BenchmarkFullRegistry` (5.5µs/op), `BenchmarkWalkGoDir` (674ns/op)                                                                                                                                                            |
| **Testdata**            | 17 fixtures across 14 directories (positive + negative for each rule)                                                                                                                                                          |
| **Actionable findings** | Every finding carries `WithSuggestion()` showing exact humanize replacement code                                                                                                                                               |
| **Nix flake**           | `flake.nix` with test, test-race, bench, build, vet, lint, coverage apps; `flake.lock` generated                                                                                                                               |
| **CI**                  | `.github/workflows/ci.yml` with test+vet and lint jobs, GOPRIVATE + SSH auth                                                                                                                                                   |
| **Git**                 | Initialized, 10 commits, pushed to GitHub                                                                                                                                                                                      |
| **AGENTS.md**           | Project context for future AI sessions                                                                                                                                                                                         |
| **README.md**           | User-facing docs with rules table, usage examples, detection strategy                                                                                                                                                          |

---

## b) PARTIALLY DONE ⚠️

| Area                     | What's Done                                       | What's Missing                                                                                                                                                                         |
| ------------------------ | ------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Test coverage**        | 78.6% across library code                         | CLI `main.go` has 0% direct coverage (integration-tested only). Several pattern helpers below 80% (`hasConst1024`: 68%, `hasStepBy3`: 79%, `hasEqualsOneBranch`: 80%). Target is 80%+. |
| **golangci-lint config** | `.golangci.yml` exists with reasonable linter set | **1 golines formatting issue** in `rule_bytes.go:68`. Config is a simplified subset (16 linters) vs the full go-linter-sdk config (80+ linters).                                       |
| **go.mod**               | Dependencies correct, versions pinned             | **Replace directives** for go-finding and go-linter-sdk are temporary hacks because go-linter-sdk has no published tags. Breaks `go get` for standalone consumers.                     |
| **go.work**              | Exists for local dev                              | Listed in `.gitignore` (correct), but no documentation on the local-dev workflow                                                                                                       |
| **Detection coverage**   | 7 of ~12 go-humanize feature categories covered   | Missing: ordinals (0 reimplementations found, low priority), `big.Int` variants, `FormatFloat`/`FormatInteger`, `WordSeries`/`OxfordWordSeries`, `ComputeSI`/`ParseSI`                 |

---

## c) NOT STARTED ❌

| Area                       | Why It Matters                                                                                                                  |
| -------------------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| **CHANGELOG.md**           | No change history. Every LarsArtmann project has one.                                                                           |
| **FEATURES.md**            | No feature inventory. Required by docs-health convention.                                                                       |
| **TODO_LIST.md**           | No actionable task list for future work.                                                                                        |
| **ROADMAP.md**             | No long-term direction document.                                                                                                |
| **LICENSE**                | **No license file.** The repo is on GitHub with no license — default copyright applies. This blocks adoption.                   |
| **CONTRIBUTING.md**        | No contributor guide.                                                                                                           |
| **CODE_OF_CONDUCT.md**     | Missing (standard for open-source repos).                                                                                       |
| **.editorconfig**          | Missing (every sibling project has one).                                                                                        |
| **examples/** directory    | No example_test.go with runnable `Example*` functions. Go convention for discoverable API docs.                                 |
| **SECURITY.md**            | Missing.                                                                                                                        |
| **godoc**                  | No `Example*` functions; package doc exists but is minimal.                                                                     |
| **Version tag**            | No `v0.1.0` git tag. Can't be `go get`'d at a version.                                                                          |
| **golangci-lint plugin**   | No `//go:build goexperiment.jsonv2` build tag handling documented for consumers who want to use this as a golangci-lint plugin. |
| **Go vulnerability check** | No `govulncheck` in CI.                                                                                                         |
| **Releaser**               | No `.goreleaser.yml` for binary releases.                                                                                       |

---

## d) TOTALLY FUCKED UP 💥

| Problem                                    | Severity   | Status                                                                                                                                                                                                                                                                |
| ------------------------------------------ | ---------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **go.mod replace directives committed**    | **HIGH**   | These are local-dev hacks. Any consumer who `go get`s this module will get build failures because `../go-finding` and `../go-linter-sdk` don't exist on their machine. Documented as "temporary" but still committed. This is the #1 blocker for standalone adoption. |
| **CLI main.go has 0% direct coverage**     | **MEDIUM** | The `output()`, `buildRegistry()`, and `stringList` functions have no unit tests. Only tested via subprocess exec.                                                                                                                                                    |
| **golines formatting violation**           | **LOW**    | `rule_bytes.go:68` has a golines issue. The lint app would fail.                                                                                                                                                                                                      |
| **No LICENSE file**                        | **HIGH**   | Repo is legally unusable. Nobody can use, modify, or distribute the code without explicit permission.                                                                                                                                                                 |
| **H001 testdata reports "0 unit strings"** | **MEDIUM** | The KMGTPE trick testdata triggers H001 via the `div1024` signal, but `countByteUnits` returns 0 because the units are embedded in a format string `"%.1f %cB"` not as standalone unit strings. The detection works but the diagnostic message is misleading.         |

---

## e) WHAT WE SHOULD IMPROVE 🔄

### Architecture

1. **Package-level analysis** — Currently function-scope only. Package-level `var` maps (like `clean-wizard/golangcilint.go:113` with `golangciLintSizeMultiplier`) are invisible. Need a file-level or package-level scan pass for H007.

2. **go/analysis integration** — The linter uses raw `go/parser` + `ast.Inspect`. The go-finding `analysis` module provides `AnalyzerDetector` for wrapping `go/analysis.Analyzer`. A reverse adapter (linter rules → `analysis.Analyzer`) would enable integration with `go vet`, `golangci-lint` plugins, and LSP. This is the single highest-impact architectural improvement.

3. **Type information** — No `go/types` checking. All detection is syntactic. This means `strings.TrimRight` vs `myPkg.TrimRight` can't be distinguished (both match `isPackageCall("strings", "TrimRight")` if the import alias differs).

4. **Configurable thresholds** — All confidence thresholds (e.g., "2+ byte units = medium", "3+ = high") are hardcoded. Should be configurable via rule options.

5. **Suppression comments** — No `//humanize-lint:ignore H001` support. The go-finding/SDK pipeline supports suppression, but standalone CLI doesn't parse source comments.

### Code Quality

6. **`patterns.go` is 780 lines** — It's the single largest file and contains 26 functions. Should be split by concern: `patterns_bytes.go`, `patterns_comma.go`, `patterns_time.go`, etc.

7. **No error wrapping in walker** — `WalkGoDir` silently skips parse errors. Should at minimum log them.

8. **Magic numbers** — The confidence thresholds (`unitCount >= 3 → ConfidenceHigh`) are undocumented domain knowledge.

### Testing

9. **No table-driven test for detection helpers** — Each `has*` function in `patterns.go` is tested indirectly via rule tests. Direct table-driven tests for each helper would catch edge cases faster.

10. **No property-based testing** — The word-boundary regex in `countByteUnits` could be property-tested against random strings to verify no false positives.

11. **No negative testdata for H004** — There's `h003_negative` and `h005_negative` but no `h004_negative`. The `len(x) == 1` fix is tested indirectly but not with a dedicated fixture.

---

## f) Up to 50 Things to Do Next

### Critical (blocks adoption)

1. Add **LICENSE** file (MIT, matching sibling projects)
2. **Remove go.mod replace directives** once go-linter-sdk gets its first tag; until then, document the workaround clearly
3. Fix **golines formatting** in `rule_bytes.go:68`
4. Add **CHANGELOG.md** with initial release entry
5. Tag **v0.1.0** once go-linter-sdk is tagged
6. Fix H001 **diagnostic message** ("0 unit strings" when triggered by div1024 alone)

### High Impact

7. Add **golangci-lint plugin mode** (`analysis.Analyzer` wrapper) — enables integration with the entire Go linting ecosystem
8. Add **package-level `var` detection** for H007 (file-scope scan, not just function-scope)
9. Add **import alias awareness** to `isPackageCall` (resolve `s "strings"` then check `s.TrimRight`)
10. Add **example_test.go** with runnable `Example*` functions for godoc
11. Split **`patterns.go`** into focused files (`patterns_bytes.go`, `patterns_comma.go`, etc.)
12. Add **FEATURES.md**, **TODO_LIST.md**, **ROADMAP.md**
13. Add **negative testdata** for H004 (`h004_negative`)
14. Add **.editorconfig** matching sibling projects
15. Add **CONTRIBUTING.md** and **CODE_OF_CONDUCT.md**

### Coverage Gaps

16. Add unit tests for `output()` function in CLI (table-driven: text/json/sarif formats)
17. Add unit tests for `buildRegistry()` (enable/disable combinations)
18. Add unit tests for `stringList` flag type
19. Add table-driven tests for `hasConst1024` (currently 68% — missing file-level const cases)
20. Add table-driven tests for `hasStepBy3` (currently 79% — missing AssignStmt branches)
21. Add table-driven tests for `hasCommaOrSeparator` (currently 95.7% — missing WriteByte/WriteRune branches)
22. Reach **80%+ coverage** (currently 78.6%)

### Detection Improvements

23. Add **H008: Ordinal** detection (`switch n%10` returning "st"/"nd"/"rd"/"th")
24. Add **H009: WordSeries** detection (`strings.Join` with "and"/"or" conjunction on 3+ items)
25. Add **H010: FormatFloat pattern** detection (`#` format string with `,` separator)
26. Detect **big.Int variants** — `BigBytes`/`BigIBytes`/`BigComma` reimplementations
27. Detect **duration formatting** more precisely (separate from reltime: `Xd Xh Xm` without "ago")
28. Add **confidence configuration** via rule options
29. Add **suppression comment** parsing (`//humanize-lint:ignore`)
30. Improve **H003 detection** for inline `.String() + " ago"` patterns (currently missed)

### Tooling & Infrastructure

31. Add **golangci-lint full config** (match go-linter-sdk's 80+ linter set)
32. Add **govulncheck** to CI
33. Add **gosec** to CI
34. Add **.goreleaser.yml** for binary releases
35. Add **bench-baseline.sh** / **bench-check.sh** scripts (like go-finding)
36. Add **version-check.sh** script
37. Add **GitHub Actions** badge to README
38. Add **pkg.go.dev** badge to README
39. Add **dependabot.yml**
40. Add **PR template** and **issue templates**

### Documentation

41. Write **detection strategy doc** explaining multi-signal approach and confidence levels
42. Write **contributing rules** doc (how to add a new H-rule)
43. Write **rule authoring guide** with annotated example
44. Update **AGENTS.md** with CLI binary info and build instructions
45. Add **godoc examples** for `DefaultRegistry`, `AllRules`, each `Rule*` factory

### Polish

46. Add **`--version` flag** to CLI
47. Add **`--rules` flag** to list all rules and exit
48. Add **color output** for text format (optional, via `--color`)
49. Add **exit code documentation** to `--help`
50. Add **`.gitattributes`** for line-ending normalization

---

## g) Questions

1. **go-linter-sdk tagging**: go-linter-sdk currently has no published git tags. Should I tag it as `v0.1.0` myself so go-humanize-linter's go.mod replace directives can be removed? Or is there a reason it hasn't been tagged yet?

2. **golangci-lint integration**: Do you want this linter to eventually run as a golangci-lint plugin (requiring `analysis.Analyzer` wrapper), or is the standalone CLI + library form sufficient for your use case?

3. **H004 (pluralization) sensitivity**: The current H004 rule fires on any `if x == 1` returning different strings — this catches real reimplementations but also catches domain-specific conditionals that happen to check for singularity (e.g. `if count == 1 { return "single" }`). Should H004 require the branches to explicitly use string parameters (singular/plural) to reduce noise, or is the broader net acceptable?

---

## Resolution (2026-07-30)

This was the first session. The "Up to 50 Things" list drove the entire v0.1.0 → v0.2.0 build-out. Key items:

| Item                                                | Status                                                            |
| --------------------------------------------------- | ----------------------------------------------------------------- |
| #1 LICENSE, #4 CHANGELOG, #11 FEATURES/TODO/ROADMAP | done — all project docs exist and are current                     |
| #5 Tag v0.1.0                                       | done at `v0.1.0`                                                  |
| #7 golangci-lint plugin mode                        | done — `plugin/plugin.go` (93.8% coverage)                        |
| #10 `example_test.go` runnable Examples             | done                                                              |
| #11 Split `patterns.go`                             | done — split into per-rule files                                  |
| #22 Reach 80%+ coverage                             | done — core 87.8%                                                 |
| #23/#24 H008 Ordinal / H009 Commaf                  | done at `e3ef534` / `2ac66b6` (9 rules total)                     |
| #29 `//nolint:gohumanize` suppression               | done — scoped `:Hxxx` + `//lint:ignore` syntax                    |
| #46 `--version`, #47 `--rules`                      | done (+ `--explain`, `--list-files`)                              |
| Q1 H004 FP filter                                   | resolved — string-in-branch + string-return-type filters (~0% FP) |

Still open — moved to `TODO_LIST.md` / `ROADMAP.md`:

- Package-level `var` detection for H007 (→ TODO_LIST T10)
- go/types type-aware detection (→ TODO_LIST T11)
- ~~CONTRIBUTING.md rule-addition checklist~~ done at `f8ba5d6`
- Per-line diagnostics (→ ROADMAP "Precision & ergonomics")
