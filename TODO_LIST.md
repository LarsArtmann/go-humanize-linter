# TODO List

> Short- and mid-term improvement tasks for go-humanize-linter.
>
> **Open work only.** When a task is finished, delete it here and record it in
> `CHANGELOG.md`. Every item below has been verified against the code as not-yet-done.
>
> Legend — Effort: `XS` ≤15 min · `S` ≤1 h · `M` ≤4 h · `L` ≥½ day.

## Summary

| #   | Task                                                                      | Tier   | Effort | Status   |
| --- | ------------------------------------------------------------------------- | ------ | ------ | -------- |
| T1  | Tag `v0.2.0` (code shipped; tag missing)                                  | High   | XS     | blocked  |
| T2  | Re-run real-world validation sweep with H008 + H009                       | High   | M      | planned  |
| T3  | Test the `--output` / `writeReport` flag                                  | High   | XS     | planned  |
| T4  | Document `--output` flag in README                                        | Medium | XS     | planned  |
| T5  | Pin `govulncheck` + `golangci-lint` versions in CI                        | High   | XS     | planned  |
| T6  | Export `RuleIDH001`–`RuleIDH009` constants                                | Medium | S      | planned  |
| T7  | Replace `nix run .#lint` `grep -v` filter with `.golangci.yml` plugin reg | High   | M      | planned  |
| T8  | GitHub Action composite `action.yml`                                      | Medium | S      | planned  |
| T9  | Configurable rules in plugin mode via `Analyzer.Flags`                    | Medium | M      | planned  |
| T10 | Package-level `var` detection for H007                                    | Low    | M      | planned  |
| T11 | go/types type-aware detection                                             | Low    | L      | planned  |
| T12 | `--config` flag for YAML/TOML rule configuration                          | Low    | M      | planned  |
| T13 | Publish to golangci-lint plugin index                                     | Low    | S      | blocked  |

---

## Release & validation

### T1 — Tag `v0.2.0` (code shipped; tag missing) · High · _blocked_

H008, H009, scoped suppression, `HumanizeDetector`, `--explain`, `--list-files`,
`--output`, and the adoption unlock (replace directives removed) are all merged to
`main` and documented in `CHANGELOG.md` under `[0.2.0]`. The only git tag is
`v0.1.0`. The release workflow (`.github/workflows/release.yml`) fires on tags, so
the 0.2.0 release artefacts have never been built.

- [ ] Tag `v0.2.0` on `main` (requires explicit user approval — never tag without it)
- [ ] Verify the release workflow fires and publishes GitHub release notes

### T2 — Re-run real-world validation sweep with H008 + H009 · High · _planned_

H001–H007 were swept against 190+ Go projects
(`docs/validation/2026-07-30_real-world-sweep.md`). H008 and H009 have never been
swept against the corpus, so the "~0% FP" claim does not yet extend to them.

- [ ] Run the linter over the 190+ project corpus with all 9 rules enabled
- [ ] Record H008/H009 finding counts and false-positive rate
- [ ] Save to `docs/validation/`

---

## Documentation

### T4 — Document `--output` flag in README · Medium · _planned_

The `--output <file>` flag shipped in v0.2.0 (`cmd/go-humanize-linter/main.go:72`)
but the README "Usage" section does not mention it. Users won't know it exists.

- [ ] Add `--output <file>` to the README CLI usage examples

---

## Test gaps

### T3 — Test the `--output` / `writeReport` flag · High · _planned_

The `--output <file>` flag (`writeReport`, `cmd/go-humanize-linter/main.go:146`) has
**no test**. The flag works (verified manually) but there is no regression guard for
file creation, content correctness, or cleanup on error. The `output()` renderer
itself has failure-path tests (`TestOutput_JSONWriterFailure`), but the
file-writing wrapper does not.

- [ ] `TestWriteReport_*` — file creation, content correctness, exit code, cleanup on error

---

## Tooling & CI

### T5 — Pin `govulncheck` + `golangci-lint` versions in CI · High · _planned_

The CI workflow installs `govulncheck` (`@latest`) and `golangci-lint`
(`version: latest`), which is non-reproducible: a breaking upstream release can flip
CI from green to red with no code change.

- [ ] Pin both tools to specific versions in `.github/workflows/ci.yml`

### T6 — Export `RuleIDH001`–`RuleIDH009` constants · Medium · _planned_

The CLI defines local `h001`–`h009` constants (`cmd/go-humanize-linter/main.go`) to
avoid `goconst` warnings, duplicating the rule IDs that live in the core package.
Exporting `RuleIDH001`–`RuleIDH009` from `humanizelint` lets the CLI import a single
source of truth.

- [ ] Export the constants from the core package
- [ ] Import them in `main.go` and `plugin/plugin.go` instead of local copies

### T7 — Replace `nix run .#lint` `grep -v` filter with `.golangci.yml` plugin registration · High · _planned_

`flake.nix:160` silences the `Found unknown linters in //nolint directives:
gohumanize` warning with a `grep -v`. This hides a real signal from the dev shell
while a raw `golangci-lint run` (as the CI lint job does) still emits it. The proper
fix is to register the project's own analyzer via `.golangci.yml`'s `plugins:` map.

- [ ] Register `gohumanize` as a golangci-lint plugin in `.golangci.yml`
- [ ] Remove the `grep -v` band-aid from `flake.nix`

### T8 — GitHub Action composite `action.yml` · Medium · _planned_

- [ ] Reusable composite `action.yml` (inputs: `path`, `enable`, `disable`, `format`)
- [ ] README "Use in GitHub Actions" section

---

## Larger / future

### T9 — Configurable rules in plugin mode via `Analyzer.Flags` · Medium · _planned_

The CLI supports `--enable`/`--disable`; the golangci-lint plugin does not. Users
cannot disable a rule from their `.golangci.yml` without `//nolint` (scoped
directives cover the common case, but config-level control is still missing).

- [ ] Add `enable` / `disable` string flags to `plugin.Analyzer.Flags`
- [ ] Filter `DetectFuncDecl` results by configured rule set in `analyzeHumanize`

### T10 — Package-level `var` detection for H007 · Low · _planned_

Only `*ast.FuncDecl` scope is scanned. A `var multiplier = map[string]int64{"KB": 1024}`
at package scope is invisible to H007.

- [ ] File/package-scope scan for `map[string]int64` with byte-unit keys
- [ ] Hook into `checkFuncDecls` as a separate decl-kind pass

### T11 — go/types type-aware detection · Low · _planned_

Detection is purely syntactic — import aliases (`s "strings"`) and typed values
are not resolved. Type info would cut false negatives on generic / aliased code.

- [ ] Resolve import aliases before `isPackageCall`
- [ ] Decide: full `go/types` or lightweight `pass.TypesInfo` in plugin path only
- [ ] Benchmark impact on scan speed

### T12 — `--config` flag for YAML/TOML rule configuration · Low · _planned_

- [ ] Define config schema (enabled/disabled, thresholds)
- [ ] Loader + `--config` flag + precedence over CLI flags

### T13 — Publish to golangci-lint plugin index · Low · _blocked_

- [ ] Blocked on a tagged, `go install`-able version (the replace directives are gone
      and deps resolve from tags, but `v0.2.0` is not tagged yet — see T1; also depends
      on `go-linter-sdk` being publicly `go install`-able)
- [ ] Submit to the plugin index once installable
