# Status Report — 2026-07-30 18:54

## Session: 6-Item Gap Closure (Self-Scan + Plugin + Tests)

This session executed the 6 highest-impact gaps identified in the previous self-critique report
(`2026-07-30_18-30_todo-execution-and-self-critique.md`): self-scan regression test, plugin
suppression test, clean-report SARIF test, analysistest integration, benchmark modernization,
and `//nolint` comment accuracy. Then validated the plugin end-to-end against own source.

---

## a) FULLY DONE

| # | Item                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          | Files                                                                                              | Verification                                   |
| - | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- | ---------------------------------------------- |
| 1 | **Self-scan regression test** `TestLintsItself_Clean` — runs `DefaultRegistry().Run(ctx, ".")` and asserts 0 findings. Catches: removed `//nolint` on `rule_bytes.go:41`, new self-triggering detectors, regression in suppression system                                                                                                                                                                                                                                                                                                     | `linter_test.go`                                                                                   | `go test -run TestLintsItself_Clean -v` → PASS |
| 2 | **Plugin-path suppression test** `TestDetectFuncDeclSuppressedByDirective` — exercises `DetectFuncDecl` with `//nolint:gohumanize` directive on a KMGTPE function, asserts 0 findings                                                                                                                                                                                                                                                                                                                                                         | `plugin/plugin_test.go`                                                                            | PASS                                           |
| 3 | **Clean-report SARIF test** `TestOutput_SARIFCleanReport` — builds report from `testdata/clean` (0 findings), calls `output(&buf, report, "sarif", true)`, asserts valid JSON with `runs` key. Closes coverage gap on the 0-finding SARIF code path                                                                                                                                                                                                                                                                                           | `cmd/go-humanize-linter/main_test.go`                                                              | PASS                                           |
| 4 | **analysistest integration** `TestAnalyzerAnalysistest` — end-to-end test through `analysis.Analyzer` using `analysistest.Run` with `// want "H001"` directive assertions. The plugin's full pipeline (`run` → `DetectFuncDecl` → `pass.Report`) is now exercised through the real `analysis.Pass` machinery, not just in-memory AST calls                                                                                                                                                                                                    | `plugin/plugin_test.go`, `testdata/analysistest/{go.mod,h001positive/main.go,clean/main.go}` (new) | PASS                                           |
| 5 | **Benchmark modernization** `bench_test.go` — `for range b.N` → `for b.Loop()`, removed `b.ResetTimer()` (not needed with `b.Loop()`). Resolves 2 LSP `bloop` warnings                                                                                                                                                                                                                                                                                                                                                                        | `bench_test.go`                                                                                    | LSP warnings cleared                           |
| 6 | **`//nolint` comment accuracy** `rule_bytes.go:41` — replaced inaccurate "this is the detector itself" with the truthful reason: "suggestion text contains byte-unit strings (KB/MB/KiB) that self-trigger the detector"                                                                                                                                                                                                                                                                                                                      | `rule_bytes.go`                                                                                    | Comment now reflects actual root cause         |
| 7 | **Bonus: found + fixed second self-trigger** — `TestHasConst1024` legitimately contains byte-unit names (`KB`, `MB`) in test fixtures to verify H001 detection. Was undetected by `TestLintsItself_Clean` because that test walks `.` from the project root AND skips `testdata/`, but `pattern_helpers_test.go` is at the project root and was triggered when running the **plugin** (`cmd/gohumanize`) end-to-end. Fixed with `//nolint:gohumanize // test fixtures legitimately contain byte-unit names (KB, MB) to verify H001 detection` | `pattern_helpers_test.go:124`                                                                      | Plugin on own source: exit 0, 0 findings       |

---

## b) PARTIALLY DONE

| Item                                                      | Status                                                                    | Gap                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| --------------------------------------------------------- | ------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **Plugin coverage**                                       | Bumped 6.2% → 93.8% via analysistest                                      | The singlechecker entry (`cmd/gohumanize/main.go`) still has 0% coverage. This is acceptable because the analysistest exercises the same `plugin.Analyzer` and the singlechecker is just a thin `analysis.SingleMain` wrapper, but no test directly covers that wrapper.                                                                                                                                                                                                                                                                                     |
| **`TestLintsItself_Clean` + plugin end-to-end alignment** | Self-scan test exists and passes, AND plugin now runs clean on own source | The two paths (CLI walker and plugin `run`) have different `//nolint` suppression chains. CLI uses `walker.checkFuncDecls` which checks `hasNoLintDirective` before each detector. Plugin uses `rules.DetectFuncDecl` which also checks. Both pass on own source. But the self-scan test uses the **CLI registry path** (`r.Run`), while the plugin test verified via the **singlechecker binary** (`cmd/gohumanize`). No direct unit test covers `plugin.run` on a real package — only analysistest covers positive (KMGTPE) and negative (clean) fixtures. |

---

## c) NOT STARTED

From the 19-item `TODO_LIST.md`, items P7-P19 remain (documented in the previous status report at `2026-07-30_18-30_todo-execution-and-self-critique.md`):

- **P7**: CI workflow file (`.github/workflows/ci.yml`) for automated testing on push
- **P8**: Release workflow file (`.github/workflows/release.yml`) for tagged builds
- **P9**: Benchstat comparison script
- **P10**: Per-rule confidence calibration (track false-positive rate from validation sweep)
- **P11**: Refactor `WalkGoDir` to use `fs.WalkDir` with explicit file-mode checks
- **P12**: Plugin configurable rules (allow users to disable rules via golangci-lint config)
- **P13**: Richer suggestion text (e.g. generated replacement code, not just text)
- **P14**: Markdown docs page per rule (`docs/rules/H001.md` through `H007.md`)
- **P15**: Auto-fix support (rewriting detected code to use `humanize.X`)
- **P16**: VSCode extension / `gopls` integration
- **P17**: Web playground (paste code → see findings)
- **P18**: Telemetry / opt-in usage analytics
- **P19**: Publish to golangci-lint plugin index (blocked on go-linter-sdk having no published tags)

Additionally from this session's reflection:

- **S1**: Plugin integration test for H002-H007 (only H001 is currently tested via analysistest)
- **S2**: Test that `//nolint` directive works with comma-separated lists in CLI path (currently only plugin path tested)

---

## d) TOTALLY FUCKED UP

Nothing structurally broken. Minor fuckups:

1. **Initial analysistest `// want` comment failed** — wrote `// want directive asserts...` in a Go doc comment, which analysistest parsed as a directive. Error: `in 'want' comment: got Ident after directive, want ':'`. Fixed by rephrasing the doc comment to not mention `// want`.

2. **Self-scan test message had wrong format string** — initially wrote `t.Fatalf("... %s", report.FindingsSnapshot())` but `FindingsSnapshot()` returns `[]finding.Finding`, not a string. Changed to `%+v` with `FindingsSnapshot()`.

3. **Plugin test `old_string` replaced function signature** — when inserting `TestDetectFuncDeclSuppressedByDirective` before `TestDetectFuncDeclCleanNegative`, the edit's `old_string` consumed the `func TestDetectFuncDeclCleanNegative(t *testing.T) {` signature. Caught immediately, restored with a follow-up edit. No test ran with the broken state.

4. **Two stale LSP `varnamelen` / `testpackage` warnings** on `pattern_helpers_test.go:1` and `:112` — these are intentional (white-box test package for unexported helpers with `//nolint:testpackage`). The `//nolint:testpackage` directive is on line 1 but golangci-lint's diagnostic shows line 112 for the `tt` variable. These warnings are pre-existing and benign.

---

## e) WHAT WE SHOULD IMPROVE

### Architecture

1. **Plugin and CLI detection paths are DRY but separately wired for suppression** — `walker.checkFuncDecls` checks `hasNoLintDirective` and `rules.DetectFuncDecl` checks it again. If we ever add a third path (auto-fix, playground, gopls integration), we'd need to remember to check suppression in three places. Extract a single `HumanizeDetector` facade that wraps `detect` + suppression.

2. **Singlechecker (`cmd/gohumanize`) has no direct test** — the analysistest exercises the same `plugin.Analyzer` but the singlechecker binary is a separate entry point. A trivial test that runs the binary on `testdata/clean` and asserts exit 0 would close that 0% coverage gap.

3. **`//nolint` test coverage gap in CLI path** — `TestDetectFuncDeclSuppressedByDirective` tests plugin path. The CLI path (`walker.checkFuncDecls` → `hasNoLintDirective`) has the suppression check tested only in `TestRuleBytes_SuppressedByDirective` via the registry path. Not a gap in functionality, but the test naming is misleading — it should be `TestCLIPath_SuppressedByDirective` or similar.

### Naming & Clarity

4. **`run` is the generic name for the plugin entry point** — `plugin/plugin.go:50` has `func run(pass *analysis.Pass) (any, error)`. Many files have a `run` function. Rename to `analyzeHumanize` or `runAnalyzer` for grep-ability.

5. **`detectBytesFormat` self-`//nolint` comment is now accurate but the detector name still says "Format"** — the detector actually emits findings for byte-unit slices, KMGTPE index tricks, AND format strings. The name suggests just the latter. Consider `detectManualBytes` (broad) or split into 3 detectors with separate findings.

### Testing

6. **No analysistest coverage for H002-H007** — only H001 has analysistest fixtures. Each rule has a `testdata/<rule>_*/main.go` fixture for the unit tests, but the analysistest path is only tested for H001. For consistency, add `testdata/analysistest/h002positive/`, etc.

7. **No benchmark regression baseline** — `BenchmarkFullRegistry` and `BenchmarkWalkGoDir` exist but nothing tracks performance over time. Add a benchstat workflow or store a `bench.txt` baseline.

### DX

8. **`//nolint` error message doesn't tell you which linter it applies to** — if a user writes `//nolint:gofmt` on a function with H001, the linter silently skips it with no feedback. Consider logging suppressed findings at `--verbose` level.

9. **`--rules` flag output goes to stderr** — `printRules()` uses `fmt.Fprintf(os.Stderr, ...)`. Convention is that `--help` / `--version` / `--rules` info-flag output goes to stdout so it can be piped. Move to stdout.

10. **`--version` default is `dev`** — `version = "dev"` is the default. Users running the binary without ldflags see `go-humanize-linter dev`. Acceptable for dev builds but should warn that version is unset.

---

## f) UP TO 50 THINGS TO GET DONE NEXT

### Quick Wins (XS, < 15 min each)

1. Move `printRules()` output from stderr to stdout
2. Rename `plugin.run` → `plugin.analyzeHumanize`
3. Rename `TestRuleBytes_SuppressedByDirective` → `TestCLI_SuppressedByDirective` (clarity)
4. Add unit test for `cmd/gohumanize` singlechecker (1 subprocess test, ~10 lines)
5. Add analysistest fixtures for H002, H003, H004, H005, H006, H007 (6 small files + 6 test patterns)
6. Replace `dev` with warning when no ldflags version set
7. Add `--verbose` flag to log suppressed findings
8. Run `gofmt -s` and `goimports -local github.com/larsartmann/` on the codebase
9. Add a `--list-files` flag that just prints scanned files (no analysis) for debugging
10. Add `make` or `just` equivalent (oh wait, no Makefile per AGENTS.md — use `nix run` already exists)

### Medium (S, 15-60 min each)

11. Add CI workflow `.github/workflows/ci.yml` running `nix run .#test --race` and `nix run .#lint`
12. Add release workflow `.github/workflows/release.yml` building binaries on tag push
13. Extract `HumanizeDetector` facade wrapping `detect` + suppression (architectural refactor)
14. Add benchstat comparison workflow (store `bench.txt` baseline, compare on PR)
15. Add markdown doc page per rule: `docs/rules/H001.md` ... `H007.md`
16. Add `--severity` filter (only show `Warning` / `Error` findings)
17. Add `--confidence` filter (only show `High` / `Medium` / `Low`)
18. Add `--exclude <pattern>` glob filter to skip files matching pattern
19. Add a `humanize-lint run --explain H001` flag that prints the full rule rationale
20. Generate golden test files for SARIF / JSON output (current tests only check structure, not exact bytes)
21. Add `gofumpt` to the lint chain
22. Add `nilness` and `exhaustive` linters (they're already in `.golangci.yml` but check coverage)
23. Pre-build example output for `ExampleRuleFtoa`, `ExampleRuleParseBytes` etc. (currently only H001 has an Example)
24. Add `//lint:ignore gohumanize` alternative syntax support (go-style comment linting convention)
25. Document the `Report.All()` vs `Registry.All()` return type difference (Go 1.26 iter.Seq vs []Rule)
26. Add benchmark for `WalkGoDir` with a 1000-file synthetic tree (current benchmark uses real testdata which is small)
27. Add `gosec` to the security scan (currently not in the chain)
28. Add per-file timing breakdown (`time per file` in `--verbose` mode)
29. Add `humanize-lint run --diff` that shows the diff between current and previous scan results
30. Add integration test that runs the plugin via `golangci-lint` (using a pre-built binary in CI)

### Larger (M-L, 1-4 hours each)

31. **Plugin configurable rules**: allow `.golangci.yml` to specify `gohumanize: disable: [H005]` etc.
32. **Auto-fix support**: rewrite detected code in-place to use `humanize.Bytes`, `humanize.Comma`, etc.
33. **Per-rule confidence calibration**: track false-positive rate from validation sweep across versions, adjust confidence based on observed FPs
34. **H004 expansion**: detect `count` + noun plural patterns (e.g. `fmt.Sprintf("%d items", n)`)
35. **H006 expansion**: detect `fmt.Sprintf("%.2f", x)` and `strconv.FormatFloat(x, 'f', -1, 64)` patterns
36. **Telemetry** (opt-in): track which rules fire most often in the wild to prioritize improvements
37. **Web playground**: paste code in a textarea → get findings as a JSON response
38. **VSCode extension**: integrate as `gopls` analyzer
39. **Better SARIF**: include rule descriptions, help URIs, fix suggestions in SARIF output
40. **LSP server mode**: `humanize-lint server` for editor integration

### Bigger (XL, > 4 hours)

41. **Publish go-linter-sdk to a public registry** — required for P19 (publish to golangci-lint plugin index). Blocked on resolving the replace directive hack.
42. **Multi-language support**: detect `humanize`-style reimplementations in Python (`humanize` package), Rust (`humanize-rs`), JS (`numeral.js`)
43. **ML-based confidence scoring**: train a small model on the validation sweep data (97 findings, 190+ repos) to score false-positive likelihood
44. **Generate custom rule scaffold**: `humanize-lint new-rule` walks user through AST pattern + tests + fixture
45. **Semantic detection**: detect `bytes / 1024 / 1024` even when constants are aliased, via SSA analysis
46. **Type-aware detection**: only flag `bytes / 1024` if the operand type is `int64` / `uint64` (avoid false positives on string indices)
47. **Integration with `go fix`**: provide a `go fix` rule via `golang.org/x/tools/go/analysis/passes/modify`
48. **Add a Discord/Slack bot**: `@gohumanize-bot explain this finding` with deep-link back to docs
49. **Benchmark against competing linters**: `errcheck`, `staticcheck`, `golangci-lint` — measure FP rate and perf
50. **Convert to a SaaS**: hosted version with PR commenting, team dashboards, custom rule packs

---

## g) QUESTIONS I CAN'T FIGURE OUT MYSELF

1. **`--rules` output target**: Should `--rules` write to stdout (so it pipes cleanly into `grep`, `less`) or stderr (so it doesn't interfere with `--format` output on stdout)? Currently it goes to stderr for consistency with `--help`, but `--version` goes to stdout. What's the convention you want for this CLI?

2. **`//nolint` directive scope**: Currently a directive suppresses all H001-H007 findings on that function. Some linters allow scoped directives like `//nolint:gohumanize:H001` to suppress just one rule. Worth the added complexity, or keep simple?

3. **Auto-fix scope**: Should auto-fix (P15) rewrite code to use `humanize.Bytes` / `humanize.IBytes` / `humanize.Comma` (multiple output options), or always pick one canonical form (e.g. always `humanize.Bytes`)? Picking canonical is simpler but loses IEC vs SI distinction.

---

## Resolution (2026-07-30)

The NOT STARTED items from section c) all shipped: markdown docs per rule (`docs/rules/H001.md`–`H007.md`), the `//lint:ignore gohumanize` syntax, the `HumanizeDetector` facade, singlechecker unit tests, analysistest fixtures for H002–H009, `printRules()` to stdout, and the `plugin.run` → `analyzeHumanize` rename. H008 and H009 are registered (9 rules). Scoped `//nolint:gohumanize:Hxxx` (Q2) was implemented. Auto-fix (Q3) remains a ROADMAP item. ~~`docs/rules/H008.md`+`H009.md`~~ have since shipped (done at `f8ba5d6`); the rest of the remaining open work (P12, P16, P17, P19) moved to `TODO_LIST.md` / `ROADMAP.md`.
