# TODO List

> Short- and mid-term improvement tasks for go-humanize-linter.
>
> **Open work only.** When a task is finished, delete it here and record it in
> `CHANGELOG.md`. Every item below has been verified against the code as not-yet-done.
>
> Legend — Effort: `XS` ≤15 min · `S` ≤1 h · `M` ≤4 h · `L` ≥½ day.

## Summary

| #   | Task                                                         | Tier   | Effort | Status  |
| --- | ------------------------------------------------------------ | ------ | ------ | ------- |
| T1  | Tag `v0.2.0` (code shipped; tag missing)                     | High   | XS     | blocked |
| T2  | Real-world validation sweep with H008 + H009 + new detection | High   | M      | planned |
| T14 | Benchmark import-alias-aware `isPackageCall`                 | Medium | S      | planned |
| T15 | Plugin integration test through `custom-gcl` binary          | Medium | S      | planned |
| T16 | Dot-import (`. "strings"`) support for alias resolution      | Low    | S      | planned |
| T17 | H009/H002 overlap disambiguation                             | Low    | M      | planned |
| T18 | Publish to golangci-lint plugin index                        | Low    | S      | blocked |

---

## Release & validation

### T1 — Tag `v0.2.0` (code shipped; tag missing) · High · _blocked_

All v0.2.0 features are merged to `main` and documented in `CHANGELOG.md` under
`[0.2.0] - Unreleased`. The only git tag is `v0.1.0`. The release workflow
(`/.github/workflows/release.yml`) fires on tags.

- [ ] Tag `v0.2.0` on `main` (requires explicit user approval — never tag without it)
- [ ] Verify the release workflow fires and publishes GitHub release notes

### T2 — Real-world validation sweep with H008 + H009 + new detection · High · _planned_

H001–H007 were swept against 190+ Go projects
(`docs/validation/2026-07-30_real-world-sweep.md`). H008, H009, package-level
var detection, and import-alias-aware detection have never been swept against
the corpus, so the "~0% FP" claim does not yet extend to them.

- [ ] Run the linter over the 190+ project corpus with all 9 rules enabled
- [ ] Record H008/H009 finding counts and false-positive rate
- [ ] Record false-positive rate for package-level var detection and alias-aware detection
- [ ] Save to `docs/validation/`

---

## Performance

### T14 — Benchmark import-alias-aware `isPackageCall` · Medium · _planned_

The `isPackageCall` function now accepts a variadic `aliases ...map[string]string`
parameter. No performance data exists for the alias-resolution path.

- [ ] Write `BenchmarkIsPackageCall_WithAliases` and `BenchmarkIsPackageCall_WithoutAliases`
- [ ] Compare with `benchstat` to confirm no regression on the no-alias path
- [ ] Profile `buildImportAliases` for large files

---

## Testing

### T15 — Plugin integration test through `custom-gcl` binary · Medium · _planned_

The golangci-lint v2 module plugin registration is verified manually (see
`.golangci.custom.yml` and the plugin.go doc comment for the workflow) but
there is no automated integration test. Unit tests pass but don't exercise
the golangci-lint runtime discovery path.

- [ ] Write a test that builds the `custom-gcl` binary, runs it on testdata,
      and asserts findings are produced
- [ ] Gate behind `testing.Short()` skip or a build tag since it requires
      `golangci-lint custom` (network + git clone)

---

## Detection improvements

### T16 — Dot-import (`. "strings"`) support for alias resolution · Low · _planned_

The `buildImportAliases` function resolves named import aliases (e.g.,
`str "strings"`) but does not handle dot imports (`. "strings"`). Dot imports
inject all exported names into the current scope, so `HasSuffix` would be called
without a package qualifier. ADR 0001 documents this as a known gap.

- [ ] Detect dot imports in `buildImportAliases` and handle bare calls to
      `HasSuffix`, `TrimRight`, etc. in pattern helpers

### T17 — H009/H002 overlap disambiguation · Low · _planned_

H009 (manual-commaf) and H002 (manual-comma-format) can both match the same
code when it involves `%.Nf` formatting plus comma-grouping loops. Currently
both rules fire independently. A prioritization or suppression mechanism
would reduce noise.

- [ ] When H009 fires, suppress H002 on the same function (or vice versa)
- [ ] Document the precedence rule

---

## Distribution

### T18 — Publish to golangci-lint plugin index · Low · _blocked_

- [ ] Blocked on a tagged, `go install`-able version (v0.2.0 not tagged yet — see T1)
- [ ] Submit to the plugin index once installable
