# Status Report — 2026-08-05 04:28 CEST — Linter Improvements to Prevent AI Mistakes

**Date:** 2026-08-05 04:28 CEST  
**Session scope:** Analyze 29 sibling-project status reports from today's `go-humanize-linter` remediation sweep and design concrete improvements to the linter so it prevents (or at least surfaces) the recurring AI mistakes observed.  
**Branch:** `main` (working tree clean; 9 commits ahead of `origin/main`, none of them from this session)  
**Author:** Crush

---

## TL;DR

Read all 29 go-humanize-linter-related status reports written today across the project corpus (4,893 lines). The linter is doing its core job — flagging hand-rolled reimplementations of `dustin/go-humanize` — but it is **silent on the failure modes that actually hurt quality**: wrong API substitution, stale suppression directives, behavior-changing output format drift, false positives that trigger cargo-culted suppressions, and no CI enforcement. This session produced an improvement plan, **no source code changes yet**. Seven sibling projects are currently in broken or degraded states that a smarter linter (or better tooling around the linter) could have prevented.

---

## a) FULLY DONE

1. **Located and read every go-humanize status report from today.**
   - 29 files across `~/projects/` and `~/forks/mindwalk/`.
   - 4,893 total lines of status-report text.
   - Confirmed with `find` + `wc -l` (29 files, 4,893 lines).

2. **Catalogued every project state.**
   - 23 projects reported "0 findings" after remediation.
   - 1 project never finished the fix (`golangci-lint-auto-configure`).
   - 4 projects introduced real regressions during the remediation (`BuildFlow`, `file-and-image-renamer`, `invoices`, `auto-deduplicate` pre-existing but blocked).
   - 1 project fixed a linter bug (`AI-Speed-Test`, not yet pushed upstream).

3. **Extracted recurring AI mistake patterns.**
   - Wrong API substitution (`humanize.SIWithDigits` used as a pluralizer).
   - Wrong suppression syntax (`//nolint:go-humanize-linter/H003` instead of `//nolint:gohumanize:H003`).
   - Behavior-changing output without warning (`KB` → `KiB`, `1.0 GB` → `1.1 GB`).
   - False positives causing cargo-culted suppressions (`strings.Join(args, " ")`, size-bucket switches).
   - No CI / no pinned binary / no verification of the linter itself.
   - Linter hints using non-existent APIs (`humanize.Plural` does not exist; `english.PluralWord` does).
   - Low-confidence findings treated as high-confidence work orders (H005 `durationMs` at 0.75 confidence).
   - No cross-project consistency check (SI vs IEC, duplicated helpers).

4. **Designed concrete linter improvements.**
   - Add a `--verify-suppressions` mode that reports `//nolint:gohumanize` directives that suppress zero findings.
   - Warn on unknown linter names in `//nolint` directives (e.g. `go-humanize-linter`).
   - Improve rule suggestion text to use exact, real API paths (`github.com/dustin/go-humanize/english.PluralWord`).
   - Fix H002/H009 false positives on `strings.Join` with space separator (already partially patched in `AI-Speed-Test` sessions; needs to be upstreamed / generalized).
   - Add an H001 false-positive filter for table-driven size-bucket constants (not byte-formatting output).
   - Add per-statement suppression support so suppressions can live on the actual suspicious line, not just the function doc-comment.
   - Add confidence-floor exit-code policy (medium → exit 2, high+ → exit 1).
   - Add a `--explain` note or `--behavior-delta` flag that warns when a fix changes visible output (SI ↔ IEC).
   - Add a `--verify-config` mode that validates config files and unknown rule IDs.

5. **Documented the 7 currently-broken sibling-project states.**
   - `BuildFlow`: `execution/plural.go` does not compile (`humanize.SIWithDigits` used as pluralizer).
   - `file-and-image-renamer`: commit `151baba` ships a wrong `//nolint` directive; linter re-fires on that commit.
   - `golangci-lint-auto-configure`: H004 finding never fixed; guessed wrong module path.
   - `auto-deduplicate`: pre-existing build errors (`undefined: v2`, `undefined: duplicate.Duplicate`) block verification.
   - `invoices`: pre-existing depguard config broken since `ab78c65`; 50+ false positives after humanize import.
   - `mr-sync`: 3 commits unpushed; `vendorHash` bump needed.
   - `emeet-pixyd`: final `vendorHash` uncommitted in working tree.
   - `AI-Speed-Test` linter repo: 2 linter-fix commits unpushed upstream.

---

## b) PARTIALLY DONE

1. **Improvement design is concrete but not prioritized.**
   - The 8 improvement areas are well-scoped, but no Pareto ranking has been applied to decide which to ship first.
   - Some improvements overlap (suppression verification + unknown-linter-name warning could share the same AST pass).

2. **Source-code reading of the linter is partial.**
   - Read `rules.go`, `pattern_helpers.go` (suppression parsing), `rule_bytes.go`, and `cmd/go-humanize-linter/main.go`.
   - Did not yet read every rule/pattern file (`rule_comma.go`, `rule_commaf.go`, `rule_plural.go`, `rule_reltime.go`, `rule_si.go`, `rule_ftoa.go`, `rule_parsebytes.go`, `rule_ordinal.go` and all `pattern_*.go` helpers).
   - Did not yet inspect test coverage to know where to add regression tests.

3. **Linter bug fix from `AI-Speed-Test` is identified but not merged into this repo.**
   - The `pattern_comma.go:62` fix (removing `" "` from separator list) lives only in `~/projects/AI-Speed-Test/docs/status/` and presumably in the linter repo working tree ahead of `origin/main`.
   - It needs to be verified, merged, and released before the other improvements.

---

## c) NOT STARTED

1. **No source code changes in this repo.** No files in `go-humanize-linter/` were modified in this session.
2. **No tests written** for any proposed improvement.
3. **No implementation of `--verify-suppressions`** or the unknown-linter-name warning.
4. **No implementation of improved suggestion text** (`english.PluralWord`, etc.).
5. **No implementation of H001 size-bucket false-positive filter.**
6. **No implementation of per-statement suppression support.**
7. **No implementation of confidence-floor exit codes.**
8. **No implementation of `--behavior-delta` / output-change warnings.**
9. **No implementation of `--verify-config`.**
10. **No decision on which improvement to ship first.**
11. **No `AGENTS.md` update** documenting the recurring AI mistakes.
12. **No `TODO_LIST.md` / `FEATURES.md` / `CHANGELOG.md` updates.**
13. **No push of the `AI-Speed-Test` linter fix commits** to `origin/main`.
14. **No verification of the 7 broken sibling projects** — identified but not fixed.
15. **No design for how the linter would detect "wrong humanize API usage"** (e.g. `humanize.SIWithDigits` for pluralization).

---

## d) TOTALLY FUCKED UP

**Nothing was broken in this repo during this session** — it was analysis-only.

However, the analysis revealed that **the linter itself contributed to several downstream breakages**, which is the central problem this session is meant to address:

| # | Project | What the linter (or its absence) caused | Severity |
|---|---------|----------------------------------------|----------|
| 1 | `BuildFlow` | Linter suggested replacing pluralization; AI used `humanize.SIWithDigits` because the actual API (`english.PluralWord`) was not obvious. Now `execution/` and `internal/cli/` do not compile. | Critical |
| 2 | `file-and-image-renamer` | Linter has no way to validate suppression syntax. AI committed `//nolint:go-humanize-linter/H003` (wrong namespace). The directive is a no-op; the finding will re-fire for anyone on that commit. | Medium |
| 3 | `golangci-lint-auto-configure` | Linter hint text says `humanize.Plural`, which does not exist. AI guessed `github.com/larsartmann/go-humanize` (nonexistent) and gave up. Finding still reported. | Medium |
| 4 | `AI-Speed-Test` | Linter H002/H009 falsely flagged `strings.Join(args, " ")` as a comma-separator. AI first cargo-culted a suppression config instead of reading the source. | Low (fixed, not pushed) |
| 5 | `DiscordSync` | H001 false-positive on a size-bucket lookup table. AI added `//nolint:gochecknoglobals` for a read-only table rather than the linter learning the difference. | Low |
| 6 | `KeyCountdown` | H005 fired at 0.75 confidence on `Nanoseconds()/1e6`. AI fixed it, but the rule message overstates the match ("1.5K, 2.3M" vs. actual ms conversion). | Low |
| 7 | `mr-sync`, `emeet-pixyd`, `invoices` | Linter fixes required `go.mod` / `vendorHash` / `depguard` changes the linter cannot see. AI shipped source changes without verifying the build pipeline. | Medium |

**Honest process failures in this session:**

- **Did not read every rule file before designing improvements.** The proposals are grounded in status-report symptoms, but the exact implementation points (which function, which line) are not yet mapped.
- **Did not run the current linter** against the sibling projects to reproduce any finding. The analysis is report-driven, not observation-driven.
- **Did not push the already-fixed `pattern_comma.go` bug.** The `AI-Speed-Test` sessions produced a real fix; this session should have upstreamed it before expanding scope.

---

## e) WHAT WE SHOULD IMPROVE

### 1. Make suppression directives machine-verifiable

**Problem:** AIs write `//nolint:go-humanize-linter/H003` because the module path is `github.com/larsartmann/go-humanize-linter`. The actual analyzer name is `gohumanize`. There is no feedback when the directive is misspelled.

**Solution:**
- Add `//nolint:verify` or `--verify-suppressions` mode that reports directives referencing an unknown linter name.
- Add a `--verify-suppressions` run that reports any `//nolint:gohumanize[:Hxxx]` directive that does not actually suppress a finding in the current run (stale or mis-scoped).

### 2. Stop suggesting APIs that do not exist

**Problem:** H004 says "use humanize.Plural / humanize.PluralWord". `humanize.Plural` does not exist in `github.com/dustin/go-humanize`. The real function is `english.PluralWord` in the `github.com/dustin/go-humanize/english` sub-package.

**Solution:**
- Update H004 suggestion text to: `humanize.Plural` is not exposed; use `github.com/dustin/go-humanize/english.PluralWord(n, singular, plural)` or `english.Plural(n, singular, plural)` depending on whether the count is rendered separately.
- Add a doc/rule example showing both shapes.

### 3. Reduce false positives that train AIs to suppress blindly

**Problem:** H002/H009 flagged `strings.Join(args, " ")` because `isSeparatorLiteral` accepted space as a comma-like separator. H001 flags table-driven size buckets (`small`, `medium`, `large` thresholds) because it sees unit strings.

**Solution:**
- Merge the `AI-Speed-Test` `pattern_comma.go` fix (drop `" "` from separator list) and add regression tests.
- Add an H001 false-positive filter: if the function contains a `switch` over thresholds with labels like `small/medium/large` and no actual byte-formatting output string, skip it.

### 4. Surface behavior-changing output as a warning, not just a style finding

**Problem:** Replacing `formatBytes` with `humanize.IBytes` changes `"1.0 KB"` to `"1.0 KiB"`. Replacing with `humanize.Bytes` changes `"1.0 GB"` to `"1.1 GB"` for `1<<30` bytes. The linter reports H001 but does not warn that the fix changes user-visible strings.

**Solution:**
- Add a `--behavior-delta` or `--strict-compat` mode that emits an extra warning when the suggested replacement's output would differ from common hand-rolled implementations.
- Update H001 suggestion text to explicitly call out the SI vs IEC choice.

### 5. Add confidence-aware exit codes

**Problem:** AIs treat every finding as a work order. Low-confidence H005 matches (0.75) caused AIs to refactor code that was not actually a hand-rolled SI-prefix formatter.

**Solution:**
- Add `--min-confidence low|medium|high|full` (default `low`).
- Return exit 2 for medium-confidence findings, exit 1 for high/full findings, so CI can distinguish "please triage" from "must fix".

### 6. Support per-statement suppression

**Problem:** Some functions mix legitimate custom status copy with a single hand-rolled time-formatting block. AIs must suppress the whole function (`//nolint:gohumanize:H003` on the function line) because they cannot suppress the specific `time.Since` call.

**Solution:**
- Allow `//nolint:gohumanize:H003` on the line immediately above the suspicious expression (not just the function declaration).
- This requires detectors to return a more specific `token.Pos` than `fn.Pos()`.

### 7. Provide a `--verify-config` mode

**Problem:** AIs create `.go-humanize-linter.yml` files with `--disable H002,H009` instead of fixing the underlying false positive.

**Solution:**
- Add `--verify-config` that validates rule IDs and emits warnings for disabled rules that correspond to known fixed false positives.

### 8. Add a project-level consistency check

**Problem:** `auto-deduplicate` ended up with both `humanize.Bytes` (SI) in one package and a custom `FileSize.String()` switch in another, producing inconsistent output.

**Solution:**
- Add a new rule or mode (`--consistency`) that flags multiple byte-formatting helpers in the same module using different conventions (SI vs IEC).
- This is a meta-rule, not an AST rule; it may belong in a separate tool or a linter plugin.

---

## f) UP TO 50 THINGS TO GET DONE NEXT

Prioritized by Pareto impact (prevent the most AI damage with the least effort first):

### Immediate (today — high leverage, small change)

1. **Push the `AI-Speed-Test` `pattern_comma.go` fix** to `origin/main` (2 commits ahead).
2. **Update H004 suggestion text** to reference `github.com/dustin/go-humanize/english.PluralWord` and `english.Plural`.
3. **Add a regression test** for `//nolint` unknown linter names (`go-humanize-linter` should warn).
4. **Add `--verify-suppressions` flag** to the CLI that reports directives suppressing zero findings.
5. **Add a test** for `--verify-suppressions` on the `file-and-image-renamer` wrong-directive case.
6. **Fix the `BuildFlow` compilation bug** (replace `humanize.SIWithDigits` with `english.PluralWord`) — this is a downstream obligation caused by the linter.
7. **Fix the `file-and-image-renamer` wrong directive** by committing the corrected `//nolint:gohumanize:H003`.
8. **Fix `golangci-lint-auto-configure` H004** using `english.PluralWord`.
9. **Add H001 size-bucket false-positive filter** with a regression test.
10. **Update rule explanations** (`--explain H004`) to show exact import path and example call.

### Short-term (this week)

11. Add `--min-confidence` flag and confidence-aware exit codes.
12. Add per-statement suppression support (return specific `token.Pos` from detectors).
13. Add `--behavior-delta` / `--strict-compat` warning for output-changing replacements.
14. Add `--verify-config` mode for `.go-humanize-linter.yml`.
15. Update `AGENTS.md` with a new "Preventing AI Mistakes" section summarizing the 8 improvements.
16. Update `CHANGELOG.md` under `[Unreleased]` for the H002/H009 false-positive fix and H004 suggestion correction.
17. Update `FEATURES.md` to mention suppression verification and confidence-aware exit codes.
18. Update `TODO_LIST.md` with the 7 broken downstream projects as blockers.
19. Add a `docs/adr/0002-suppression-verification.md` explaining why suppression verification is necessary.
20. Add a `docs/adr/0003-confidence-exit-codes.md` for confidence-aware CI behavior.
21. Run `nix run .#test` and `nix run .#test-race` after each change.
22. Run `nix run .#lint` and fix the 3 pre-existing findings (`funlen`, `godox`, `golines`).
23. Add testdata fixtures for every new false-positive class.
24. Add an integration test that runs the linter against a synthetic repo containing each AI mistake pattern.
25. Add a `docs/status/2026-08-05_04-28_linter-improvements-followup.md` once implementation starts.

### Medium-term (next sprint)

26. Implement the project-level consistency check (SI vs IEC, duplicated helpers).
27. Add a `go-humanize-linter doctor` subcommand that audits downstream project health.
28. Add a `--fix` mode that proposes exact replacements (with choice of SI/IEC for H001).
29. Improve H005 detection to distinguish true SI-prefix output from ms conversion.
30. Add a rule for obviously-wrong humanize API usage (e.g. `humanize.SIWithDigits` where the result is used as a plural suffix).
31. Add corpus-wide regression test: run the linter against all sibling projects and assert no new false positives.
32. Add a `mise`/`task` target that runs the linter against all sibling projects.
33. Pin the linter binary location in sibling projects (replace `/tmp/go-humanize-linter` with a Nix flake app or committed tool).
34. Add CI integration examples to `action.yml` and `.github/workflows/ci.yml`.
35. Document the correct `//nolint:gohumanize[:Hxxx]` syntax in every rule doc (`docs/rules/H001.md`–`H009.md`).
36. Add a `--list-suppressions` mode that shows all active `//nolint:gohumanize` directives.
37. Add a `--stats` mode that reports rule hit counts and confidence distribution.
38. Add a `golangci-lint` custom-lint verification step to `nix run .#custom-lint`.
39. Investigate why H003 reports function position instead of the `time.Since` call site; fix if low-effort.
40. Add a test that the CLI exits 0 with zero findings (regression for the reported exit-code bug).

### Long-term / architectural

41. Refactor detectors to return a `token.Pos` + message pair instead of a whole-function finding.
42. Split `pattern_helpers.go` into `suppression.go`, `finding_builder.go`, and `pattern_imports.go`.
43. Consider a second-order rule: "function contains both `humanize.Bytes` and `humanize.IBytes` callers — pick one".
44. Consider an import-graph rule: "module depends on `dustin/go-humanize` but still has H001-H009 findings".
45. Add a "migration audit" report format that shows before/after output strings for every H001 finding.
46. Add a public `SuggestFix(ruleID, signals)` API so the CLI and plugin can share fix suggestions.
47. Evaluate whether the linter should ship as a `golangci-lint` bundled linter via a PR upstream instead of a module plugin.
48. Add a sibling-project CI matrix that runs the latest linter against representative repos nightly.
49. Write a case-study doc (`docs/case-studies/ai-mistakes-2026-08-05.md`) using the 29 reports as evidence.
50. Schedule a brutal-self-review of the linter once the improvements are merged.

---

## g) UP TO 3 QUESTIONS I CANNOT FIGURE OUT MYSELF

1. **Which improvement should ship first?** The fixes fall into two buckets: (a) linter-core fixes that prevent false positives and wrong suggestions (H002/H009, H004 suggestion text, suppression verification), and (b) downstream cleanup of the 7 broken projects. The latter are consequences of the linter's current gaps; fixing the linter first prevents more damage, but the downstream projects are broken *now*. Should I prioritize upstream linter fixes or stop-the-bleeding in sibling projects?

2. **Should the linter detect "wrong humanize API usage" as a new rule?** For example, `humanize.SIWithDigits(float64(count), 0, "")` used as a plural suffix is semantically wrong but is not a hand-rolled reimplementation of humanize. Adding a rule for misused humanize APIs expands the linter's scope beyond its current charter. Is that in scope, or should that live in a separate "humanize-usage" linter?

3. **How opinionated should the linter be about SI vs IEC?** The 29 reports show AIs chose both `humanize.Bytes` and `humanize.IBytes` for the same codebase. Should the linter (a) remain neutral and just flag manual implementations, (b) enforce one convention project-wide, or (c) add a project-level config setting (`default_units: si|iec`) that suggestions respect? The right answer depends on whether you want the linter to make product decisions or stay purely technical.
