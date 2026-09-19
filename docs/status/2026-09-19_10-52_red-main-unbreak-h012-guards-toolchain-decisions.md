# Status: Red-Main Unbreak, H012 Position-Filter Guards, Toolchain Policy

**Date:** 2026-09-19 10:52 CEST
**Session span:** 2026-09-18 ~22:00 → 2026-09-19 10:50 CEST (overnight blocked on user decisions)
**Scope:** This session only — handoff state verification, red-main fix, H012 regression guards, verification battery, user decisions. Not a project-wide audit.
**State at writing:** HEAD `642a426`, main CI green, working tree clean, `go` directive 1.26.7, v0.3.0 is the released version (H012 + guards sit in `Unreleased`).

---

## a) FULLY DONE

1. **Handoff state verified, red main discovered and root-caused.** The auto-commit daemon had swept an unintended `go 1.26.7 → 1.27.1` directive bump into heuristic commit `888f3c3`. Every CI job failed (`go-version: "1.26"` + `GOTOOLCHAIN: local` in both workflow jobs → `go: go.mod requires go >= 1.27.1`); local builds failed identically (go 1.26.7, `GOTOOLCHAIN=local` from `~/.config/go/env`). Q1 from the prior session ("keep or revert?") dissolved — the bump broke everything.
2. **go.mod restored to 1.26.7 — twice.** First restore: explicit fix commit `ebc9df5` (pushed, CI green). The directive was re-bumped in the working tree mid-session by unidentified background hygiene (go 1.27.0/1.27.1 toolchains are cached in GOMODCACHE; `flake.lock` churns in parallel); second restore folded into `b404990`. Directive has held at 1.26.7 since.
3. **CI green on every commit since the fix**, each run checked: `ebc9df5`, daemon commits `084986c`/`15b04d6`/`9c76bd4`, `b404990`, `642a426`.
4. **H012 conjunction-position filter regression-guarded** (the top unprotected item from the handoff):
   - `testdata/h012_negative/main.go`: three new negatives — Sprintf-format prose "and" next to a comma join (the original bug-tracking-schema FP shape), bare-variable prose ("cause and effect"), prefix-join-without-conjunction.
   - `pattern_word_series_test.go` (white-box, follows the `pattern_generated_test.go` `//nolint:testpackage` convention): `TestCollectWordSeriesEvidence_ConjunctionPosition` with the mutant-killing pair — prose must NOT set `conjunctionLiteral` (with `anyCommaJoin=true` asserted so the discrimination has teeth), Join separator and concat operand MUST set it.
   - Fixture-level test comment updated (`TestRuleWordSeries_Negative`).
5. **Verification battery (all green):** workspace `go test ./... -count=1`; `GOWORK=off go test ./... -count=1` (published-deps behavior); `nix run .#lint` → **0 issues**, which also proved the 5 lingering LSP diagnostics (mnd ×2, wsl_v5, godoclint, golines) are stale — local nix golangci-lint 2.13.2 and CI 2.12.2 are at parity on this code.
6. **pkg.go.dev v0.3.0 verified:** published Sep 18, H010 + examples rendered, docs complete. H012 correctly absent (post-dates the tag). The tag's frozen README shows the Action example as `@v0.2.0`; current main already says `@v0.3.0` — self-heals at the next release.
7. **AGENTS.md updated** (`b404990`): new gotcha "Go-directive bumps break CI instantly" (mechanism, incident timeline, revert-on-sight rule) + H012 bullet now documents the position filter and its guards; testdata inventory entry expanded.
8. **User decisions collected (4/4 answered) and recorded** in `TODO_LIST.md` (`642a426`): keep Go 1.26 with revert-on-sight; leave consumer repos as demonstration corpus; no branch protection; no foreign-repo actions (T18 plugin-index PR and T20 gogenfilter release declined — sections removed, decision dated; T32 deploy-key cleanup given its detail section).
9. **CHANGELOG Unreleased line for the H012 guards** added at report time (the on-sight doc-fix permission fired late — see e).

## b) PARTIALLY DONE

1. **"Commit the dirty files" handoff item** — resolved, but imperfectly: my H012 fixture + white-box test landed inside daemon heuristic commits `084986c` / `15b04d6` (content fully intact and CI-verified; messages meaningless). My explicit commits cover the go.mod fix, docs, and decisions.
2. **Verification battery breadth** — build + tests + lint done; `go vet` and race mode (`nix run .#vet`, `nix run .#test-race`) were NOT run this session.
3. **go.mod stability** — managed by reversion, not root-caused. The bumper is unidentified (prime suspect: `project-discovery-daemon`, a nix-store binary with no discoverable config). User chose this posture deliberately; recurrence risk stands.

## c) NOT STARTED

1. H012 micro-benchmark (H010 has `BenchmarkH010Detector`; H012 has none).
2. Nested/depth-2 repo sweep support — corpus method scans `~/projects/*` top-level only; `games/KeyCountdown`-style layouts need targeting by hand.
3. v0.4.0 release — H012 + H010 fixture pack + guards + enrichments all sit in `Unreleased`; nothing ships to users until a tag is cut.
4. T32 — remove read-only deploy keys from the 4 dependency repos (foreign GitHub settings; ready, untouched).
5. Corpus re-sweep after this morning's gogenfilter v3.6.1 pin (`9c76bd4`) — detection behavior could shift; not re-validated beyond the test suite.
6. Bumper investigation — declined, offer stands.

## d) TOTALLY FUCKED UP

Nothing is left broken. Permanent blemishes, honestly filed:

1. **History readability damage (irreversible without forbidden history rewrite):** two real feature commits (H012 fixtures, white-box test) carry `chore: auto-commit N changed file(s) (heuristic)` messages. The handoff explicitly warned "stage+commit immediately per batch" — I batched edits across code + docs before committing and lost the race twice.
2. **Assumption in place of verification:** this morning's `9c76bd4` appeared between my commits; I pushed past it assuming "probably the daemon's flake.lock churn." It was actually a deliberate, well-messaged gogenfilter v3.6.1 bump (Windows path fix). Harmless outcome, wrong process — I only learned the content while writing this report.
3. **Wasted round trips:** first `question` call rejected on schema (missing per-question `description`); two `multiedit` calls rejected on stale file mtime because I didn't re-View after learning the daemon mutates files. Three avoidable tool failures.

## e) WHAT WE SHOULD IMPROVE

1. **Commit-per-batch discipline.** After tests pass on a code batch, commit it before touching any other file. The daemon races every edit window; docs can follow in their own commit.
2. **Re-verify after any pause.** Overnight question-blocks and long waits are windows for foreign commits (user, daemon, dependabot). `git log` before building on top — never assume commit authorship from the push fast-forward.
3. **CHANGELOG on merge, not on report.** The Unreleased entry for the guards should have been written when the tests went green; it took a status-report prompt (standing on-sight permission existed).
4. **Close the battery:** `go vet` and race mode belong in every verification claim; "tests + lint" is not the full AGENTS command set.
5. **Trust order: CLI > LSP.** Five stale LSP warnings consumed attention that one `nix run .#lint` settled. The lesson was already in AGENTS; applying it up front would have saved a cycle.
6. **The go-directive war is systemic.** Revert-on-sight works but does not scale; if bumps recur weekly, the cheap next step is a systemd/journal peek at `project-discovery-daemon` (user-gated).

## f) NEXT — ranked, realistic (24 items; the prior 50-item backlog lives in `docs/status/2026-09-18_22-01_*.md` and was not re-verified this session)

**Release track**

1. Decide v0.4.0 timing; when cutting: CHANGELOG `Unreleased → 0.4.0`, full go-release phases (tag → release.yml → proxy verify → clean-module `go get` → action.yml bump).
2. After 0.4.0: verify pkg.go.dev renders H012 + guards and flips "Go to latest".
3. Decide on a doc-only v0.3.1 for the frozen tag's README `@v0.2.0` Action example — or let 0.4.0 supersede.
   **Detection/rules track**
4. `BenchmarkH012Detector` (mirror the H010 benchmark shape).
5. Corpus re-sweep under gogenfilter v3.6.1 (detection table may have moved; 169-repo loop script exists).
6. Nested-repo (depth-2) sweep support; re-run sweep including `games/*`.
7. H011 manual-si-parse — keep deferred; re-evaluate only on new upstream go-humanize tags or corpus demand.
8. H010 documented gap: NewReplacer strip form stays undetected (negative fixture exists) — implement only if corpus shows demand.
   **Verification track**
9. Run `nix run .#vet` and `nix run .#test-race` to close this session's battery gaps.
10. Add a periodic (weekly?) self-scan + corpus sweep ritual so red-main-drift and detection regressions surface faster than 5 days.
    **Hygiene/infra track**
11. T32: remove read-only deploy keys from go-linter-sdk, go-finding, gogenfilter, go-error-family.
12. Watch for go.mod re-bumps; if frequency grows, revisit bumper investigation with the user.
13. Decide re-pin cadence for gogenfilter (sibling drifts ahead of pins; v3.6.1 pinned today, sibling carries more).
14. Consider documenting the auto-commit daemon + question-block pause protocol in AGENTS so future sessions re-verify by default (partially done via the gotcha).
15. LSP stale-diagnostic note: if the 5 warnings persist across sessions, an LSP restart or cache clear is warranted (not a code problem — nix lint is the source of truth).
    **Docs track**
16. After 0.4.0: refresh README rules table example outputs if example tests changed.
17. Keep `docs/rules/H012.md` in sync if the position filter ever gains forms (it documents the current four subtest cases only implicitly).
18. Record the consumer-repo corpus decision (leave-as-is) in the sweep doc so future sweeps don't "fix" them by accident.

## g) QUESTIONS (cannot be figured out from the repo)

1. **v0.4.0 timing:** cut now so H012 + guards reach users (Action installs `@latest`), or accumulate more first?
2. **v0.3.1 doc-patch:** worth a patch release just to fix the frozen v0.3.0 README's `@v0.2.0` Action example on pkg.go.dev, or fold into 0.4.0?
3. **Sibling-bump protocol:** your manual dependency bumps (like this morning's gogenfilter v3.6.1) — should they auto-trigger a corpus re-sweep from my side, or only on explicit request?
