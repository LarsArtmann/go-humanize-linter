# TODO List

> Short- and mid-term improvement tasks for go-humanize-linter.
>
> Every item below is broken into **atomic sub-tasks** (~5–12 min each) so work
> can be picked up, verified, and committed independently. Items are sorted by a
> **Priority** score = `Impact × CustomerValue ÷ Effort` (higher = do first).
>
> Legend — Effort: `XS` ≤15 min · `S` ≤30 min · `M` ≤2 h · `L` ≥½ day.
> Status: `planned` · `in-progress` · `done`.

## Priority summary

| P#  | Item                                                              | Tier            | Impact | Cust.Value | Effort | Priority | Status   |
| --- | ----------------------------------------------------------------- | --------------- | :----: | :--------: | :----: | :------: | -------- |
| P1  | `//nolint:gohumanize` suppression directive                       | Critical        |   5    |     5      |   M    | **4.2**  | done     |
| P2  | `--version` flag                                                  | High / polish   |   2    |     4      |  XS    | **8.0**  | done     |
| P3  | `--rules` flag                                                    | High / polish   |   3    |     4      |  XS    | **12.0** | done     |
| P4  | CLI unit tests (`buildRegistry`/`output`/`stringList`)            | Medium          |   4    |     3      |   S    | **6.0**  | done     |
| P5  | Self-exclusion: don't flag own source                             | Critical        |   2    |     2      |  XS    | **8.0**  | done     |
| P6  | `example_test.go` runnable Examples                              | Medium          |   3    |     4      |   S    | **6.0**  | done     |
| P7  | analysistest integration test for plugin (6.2% → up)              | High            |   4    |     3      |   M    | **3.0**  | planned  |
| P8  | SARIF output integration test                                     | Medium          |   2    |     2      |  XS    | **4.0**  | done     |
| P9  | Save validation sweep to `docs/validation/`                       | High            |   2    |     3      |  XS    | **6.0**  | done     |
| P10 | Table-driven tests for `hasConst1024` / `hasStepBy3`              | Medium          |   2    |     2      |   S    | **4.0**  | done     |
| P11 | Reach 80%+ coverage (reframe: CLI + plugin)                       | Medium          |   4    |     2      |   S    | **4.0**  | partial  |
| P12 | Configurable rules in plugin mode (flags)                         | Medium          |   3    |     3      |   M    | **2.3**  | planned  |
| P13 | GitHub Action for running linter in CI                            | Medium          |   3    |     3      |   S    | **3.0**  | planned  |
| P14 | H008 — `humanize.Ordinal` rule                                    | Low / future    |   3    |     2      |   M    | **1.5**  | planned  |
| P15 | H009 — `humanize.Commaf` rule variant                             | Low / future    |   3    |     2      |   M    | **1.5**  | planned  |
| P16 | Package-level `var` detection for H007                            | Low / future    |   2    |     1      |   M    | **0.8**  | planned  |
| P17 | go/types type-aware detection                                     | Low / future    |   5    |     3      |   L    | **0.9**  | planned  |
| P18 | `--config` flag for YAML/TOML rule configuration                  | Low / future    |   2    |     2      |   M    | **1.0**  | planned  |
| P19 | Publish to golangci-lint plugin index                             | Low / future    |   2    |     2      |   S    | **2.0**  | blocked  |

> Notes on the legacy "74.1% coverage" entry: the core library is now at **88.4%**.
> The real coverage gap is the **CLI (0% direct)** and **plugin (6.2%)**, which is
> why P4, P7, and P11 are grouped together.

---

## P1 — `//nolint:gohumanize` suppression directive  · Critical · _done_

Users hit false positives and need an escape hatch. Honour the standard
`//nolint` convention used across the Go ecosystem.

- [x] `hasNoLintDirective(file, fn, "gohumanize")` helper — scans comment groups overlapping `fn.Pos()` for `//nolint:gohumanize` or `//nolint:all`
- [x] Wire suppression into `checkFuncDecls` (CLI path)
- [x] Wire suppression into `DetectFuncDecl` / plugin `run()` path
- [x] testdata `h001_suppressed` → expect 0 findings
- [x] Unit test for the directive helper
- [x] README suppression section

## P2 — `--version` flag  · Polish · _done_

- [x] `version` var (overridable via `-ldflags`) + `--version`/`-v` flag
- [x] Print `go-humanize-linter <version>` and exit 0
- [x] CLI test asserting version output

## P3 — `--rules` flag  · Polish · _done_

- [x] `--rules` flag prints ID · name · severity · description table and exits 0
- [x] CLI test asserting all 7 rules appear

## P4 — CLI unit tests (`buildRegistry` / `output` / `stringList`)  · Medium · _done_

- [x] Refactor `output()` to take `io.Writer` (testable without subprocess)
- [x] `buildRegistry`: enable-only, disable-only, mixed, unknown-id cases
- [x] `output`: text / json / sarif to a buffer
- [x] `stringList` `Set`/`String` round-trip

## P5 — Self-exclusion: don't flag own source  · Critical · _done_

The linter flags its own `rule_bytes.go` (H001). Resolved cleanly by annotating
the self-matching detectors with the suppression directive from P1 — no special
"skip own path" hack needed.

- [x] Annotate self-flagging detectors in `rule_bytes.go` with `//nolint:gohumanize`
- [x] Self-scan test: running the CLI on the repo's own source yields 0 findings

## P6 — `example_test.go` runnable Examples  · Medium · _done_

- [x] `ExampleDefaultRegistry`, `ExampleAllRules`, `ExampleDetectFuncDecl`
- [x] `go test ./... -run Example` passes (verifiable in `go doc`)

## P7 — analysistest integration test for plugin  · High · _planned_

`plugin.run()` is at 6.2% coverage; only `Analyzer` metadata is asserted today.

- [ ] Create `testdata/src/` analysistest layout (positive + clean packages)
- [ ] `analysistest.Run` with `plugin.Analyzer` + expected diagnostics
- [ ] Assert generated-file skipping (`_gen.go`) inside the harness

## P8 — SARIF output integration test  · Medium · _done_

- [x] CLI test: `--format sarif` emits valid JSON containing `"runs"`

## P9 — Save validation sweep to `docs/validation/`  · High · _done_

- [x] `docs/validation/2026-07-30_real-world-sweep.md` with the 10-project results

## P10 — Table-driven tests for `hasConst1024` / `hasStepBy3`  · Medium · _done_

- [x] `pattern_helpers_test.go` table tests: file-level const, in-fn const, MUL chains, `i+=3`, `i=i+3`

## P11 — Reach 80%+ coverage (reframe: CLI + plugin)  · Medium · _partial_

Core is 88.4%. Remaining work folded into P4 (CLI) and P7 (plugin).

- [x] Core library ≥ 80% (currently 88.4%)
- [ ] CLI direct coverage ≥ 60% (P4 adds text/json/sarif + registry; remaining is `main()`)
- [ ] Plugin coverage ≥ 50% (blocked on P7 analysistest)

## P12 — Configurable rules in plugin mode  · Medium · _planned_

- [ ] Add `enable` / `disable` string flags to `plugin.Analyzer.Flags`
- [ ] Filter `DetectFuncDecl` results by configured rule set in `run()`
- [ ] analysistest covering enable/disable behaviour

## P13 — GitHub Action for running linter in CI  · Medium · _planned_

- [ ] Reusable composite `action.yml` (input: `path`, `enable`, `disable`, `format`)
- [ ] README "Use in GitHub Actions" section

## P14 — H008 — `humanize.Ordinal` rule  · Future · _planned_

- [ ] Detect `switch n%10` returning `"st"`/`"nd"`/`"rd"`/`"th"`
- [ ] `rule_ordinal.go` + `pattern_ordinal.go` + testdata (positive + negative)
- [ ] Register in `AllRules()` / `DefaultRegistry()` / `DetectFuncDecl()`

## P15 — H009 — `humanize.Commaf` rule variant  · Future · _planned_

- [ ] Detect `%f`/`FormatFloat` + manual `.` / `,` separator grouping
- [ ] rule + pattern + testdata + registration

## P16 — Package-level `var` detection for H007  · Future · _planned_

- [ ] File/package-scope scan for `map[string]int64` with byte-unit keys
- [ ] Hook into `checkFuncDecls` as a separate decl-kind pass

## P17 — go/types type-aware detection  · Future · _planned_

- [ ] Resolve import aliases before `isPackageCall` (e.g. `s "strings"`)
- [ ] Decide: full `go/types` or lightweight `pass.TypesInfo` in plugin path only
- [ ] Benchmark impact on scan speed

## P18 — `--config` flag for YAML/TOML rule configuration  · Future · _planned_

- [ ] Define config schema (enabled/disabled, thresholds)
- [ ] Loader + `--config` flag + precedence over CLI flags

## P19 — Publish to golangci-lint plugin index  · Future · _blocked_

- [ ] Blocked on go-linter-sdk first tag (so `go.mod` replace directives can be removed)
- [ ] Submit to the plugin index once `v0.1.0` is gettable
