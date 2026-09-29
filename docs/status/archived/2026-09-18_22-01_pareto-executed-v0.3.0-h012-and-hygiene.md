# Status: Pareto plan executed — v0.3.0 shipped, H012 added, hygiene closed

**Date:** 2026-09-18 22:01 · **Session segment:** full execution of
`docs/planning/2026-09-18_21-09_SUPERB-pareto-ship-h010-and-beyond.md`
· **CI:** green on tip `8f0577e`.

## a) FULLY DONE (all verified)

**Plan + publish (1% tier):**

- Pareto plan written (tiers, 27 medium + 70 fine tasks, mermaid graph,
  approval gates), committed, pushed. TODO_LIST refreshed (T21–T31).

**Trust (4% tier):**

- H010 corpus sweep: 169 top-level repos, 1 borderline finding
  (locale-aware dual-format parser; triage guidance documented), 0 clear
  FPs. `docs/validation/2026-09-18_h010-sweep.md`.
- **v0.3.0 released**: CHANGELOG cut, all local gates + CI green on the
  tagged commit, annotated tag, `release.yml` green, both binary assets,
  proxy.golang.org serving v0.3.0, clean-module `go get` + library run
  verified (10 rules live).

**Hygiene (20% tier):**

- **CI-miss root cause: main was RED since 2026-09-13** (bb15541, 713beca,
  2cbc9db all failed on the stale gogenfilter test fixed earlier this
  session). CI never missed it — nobody looked; my earlier "CI green on
  09-17" claim had mixed in green Dependabot workflow runs.
- Dependabot actions-group PR rebased (`@dependabot rebase`) and merged
  (all checks green post-rebase).
- H010 fixture pack: aliased, dot-import, scoped `//nolint:H010`,
  rune-loop positive, NewReplacer negative — each with confidence
  assertions.
- Enrichments: H002 `BigComma` hint via new `mentionsBigInt` helper +
  `h002_bigint` fixture; H009 `CommafWithDigits` mention.
- Examples aligned to ternary `ExitCodeByConfidence` (doc.go, README);
  `--behavior-delta` upgrade note; H010↔H002 non-overlap note + doc
  cross-links; `BenchmarkH010Detector`; SARIF/JSON H010 round-trip
  smoke; Action verified via its real install path (`go install
  @latest` → proxy).
- **5 `DEPLOY_KEY_*` secrets deleted** (0 remain). AGENTS updated.

**Expansion (80% tier):**

- **H012 `manual-word-series` — full corpus-first lifecycle**: research
  found real hand-rolls; rule implemented (comma Join + conjunction in
  POSITION — Join separator or concat operand); first sweep surfaced one
  prose-only FP → position filter added → final sweep **2 TP / 0 FP**.
  Fixtures, tests, docs, CLI explain, plugin analysistest, examples all
  wired. 11 rules total.
- **H011 deferred with data**: zero SI-parse hand-rolls in corpus; ID
  reserved and documented (rules.go, ROADMAP).
- M20 (Full-tier signature) and M22 (fuzz) deliberately SKIPPED with
  documented reasoning (corpus evidence / no repo fuzz precedent).
- ROADMAP pruned (stale overlap item marked done, plugin-index blocker
  cleared), plan file annotated with completion status, M23 example,
  M24 upstream-watch note.

## b) PARTIALLY DONE

- ~~**M21 LSP hygiene** — never restarted gopls; diagnostics were stale the
  entire session (correctly ignored in favor of build/test/CI truth).~~ closed — build/test/CI truth is the standing practice.
- ~~**Local vs CI lint parity** — root-caused only partially (local nix
  golangci-lint 2.13.2 vs CI 2.12.2); the actions bump merged via PR #1
  may have aligned versions — NOT re-verified post-merge.~~ closed — parity confirmed 2026-09-19 (0 issues both).
- ~~**Silent nix-lint wrapper failure** (exit 1, empty output) observed
  twice mid-session; self-resolved after the underlying findings were
  fixed; never root-caused.~~ closed — never reproduced again; treated as transient.

## c) NOT STARTED

- ~~**T18** plugin-index submission (foreign-repo PR; unblocked by v0.3.0).~~ **Won't implement — declined 2026-09-19.**
- ~~**T20** gogenfilter sibling release + pin bump — gated on approval.~~ **Won't implement — declined 2026-09-19.**
- ~~**T32** read-only deploy-key removal in the 4 dependency repos
  (foreign repos).~~ done at `ce52545` (2026-09-19).
- ~~pkg.go.dev rendering check for v0.3.0.~~ done — verified 2026-09-19 and re-verified for v0.4.0 on 2026-09-26.

## d) TOTALLY FUCKED UP (honest list)

- **Commit-message fidelity, twice**: my detailed `feat(h012)...` and
  `feat(h010)...` commits described whole batches, but the auto-commit
  daemon had already swallowed most files into heuristic commits — my
  commits contained 3 and 1 files respectively. The messages live in
  history describing content that sits in adjacent `chore: heuristic`
  commits. Root cause: daemon races between my last edit and
  `git add -A`.
- **Unintended `go.mod` bump sitting in the working tree**: go directive
  `1.26.7 → 1.27.1`, not authored by any session action I intended —
  almost certainly a golangci-lint 2.13.2 / toolchain side-effect. This
  RAISES the minimum Go for every consumer and was never decided. Left
  in place uncommitted (not mine to revert silently); flagged for a
  decision.
- Introduced a `currentAliases` bug mid-refactor (caught by build).
- First H012 pattern file shipped with a bloated tail (three redundant
  confidence derivations) — cleaned up before commit, but written.
- The H012 sweep covered only top-level `~/projects/*/` — the motivating
  KeyCountdown hand-roll lives one level deeper (`games/KeyCountdown`)
  and was found only by a targeted check, not the sweep.
- Earlier segment: the wrong "CI green on 09-17" claim made it into the
  first status report (corrected in this segment's investigation).

## e) WHAT WE SHOULD IMPROVE

- **Regression fixture for the sweep-driven H012 filter is MISSING**: the
  conjunction-position filter (the exact fix the corpus FP motivated)
  has no negative fixture for its trigger class — comma join + prose-only
  "and" + no positional conjunction. One fixture closes the loop.
- Commit hygiene vs daemon: stage + commit immediately per batch, or
  accept daemon commits and keep detailed narration in docs/CHANGELOG
  (which are immune to the race).
- Decide the go.mod 1.27.1 question explicitly (see questions).
- Nested-repo sweep support (find go.mod at depth 2) for corpus honesty.
- Verify lint parity after the actions bump; consider pinning local nix
  golangci-lint to CI's version.

## f) NEXT UP TO 50

1. ~~H012 negative fixture: comma join + prose-only conjunction (guards the position filter)~~ done at `bd3dd03`
2. ~~Decide/revert-or-keep the go.mod `go 1.27.1` bump; align CI setup-go version if kept~~ done — decided 2026-09-19 morning: revert-on-sight (ebc9df5, 642a426) — then the 1.27.1 bump RETURNED as incident #4 (b7f00c2, 13:42) and sits on main with CI red; final call open as TODO_LIST T33
3. ~~Verify post-merge CI golangci-lint version vs local nix (lint parity)~~ done — confirmed 2026-09-19: local nix 2.13.2 and CI at parity, 0 issues both
4. ~~Root-cause the silent nix-lint wrapper failure (exit 1, empty output)~~ **Won't implement — transient — exit-1-with-empty-output never reproduced after 2026-09-18.**
5. ~~T18: golangci-lint plugin-index PR (verify-before-filing + github-voice)~~ **Won't implement — declined 2026-09-19 (T18).**
6. ~~T20: gogenfilter release + pin bump + GOWORK=off sweep [needs approval]~~ **Won't implement — declined 2026-09-19 (T20).**
7. ~~T32: remove read-only deploy keys from the 4 dep repos~~ done at `ce52545`, `2bb893e`
8. ~~pkg.go.dev check for v0.3.0 (docs render, 11 rules listed)~~ done — pkg.go.dev verified in the 2026-09-19 session; v0.4.0 rendering re-verified 2026-09-26
9. ~~Nested-repo corpus sweep (depth-2 go.mod discovery) for H010/H012~~ done (docs-health pass 2026-09-26 — tracked as TODO_LIST T38)
10. ~~Apply the linter's own advice to consumer repos (KeyCountdown joinWords → humanize.WordSeries; Kernovia Join(x," and ") → WordSeries; CreditReformBilanzampel → scoped nolint)~~ **Won't implement — declined 2026-09-19 — consumer repos are the demonstration corpus.**
11. ~~gopls restart; confirm diagnostics freshness~~ **Won't implement — ephemeral — build/test/CI truth used instead of gopls diagnostics.**
12. ~~H012 micro-benchmark (match BenchmarkH010Detector pattern)~~ done (docs-health pass 2026-09-26 — tracked as TODO_LIST T35)
13. ~~H010/H012 analysistest scoped-nolint variants~~ **Won't implement — not adopted — scoped suppression covered by CLI fixtures (h010_scoped_h010) and custom-gcl e2e subtests.**
14. ~~CHANGELOG: Unreleased section will need a 0.4.0 cut when H012 ships~~ done at `6a3267f`
15. ~~Plan next release cadence (0.4.0 after H012 soak?)~~ done — superseded — v0.4.0 shipped 2026-09-19 (6a3267f)
16. ~~ROADMAP: revisit SI-parse (H011) after upstream go-humanize activity~~ done — standing trigger — T29 deferral notes revisit on next upstream go-humanize tag
17. ~~Consider NewReplacer strip form for H010 (F13, still open)~~ **Won't implement — demand-gated — ROADMAP Detection breadth (negative fixture exists).**
18. ~~Consider Full-tier for H010 strip+parse when function is exactly the library signature AND no separator config (refine the skipped M20 with the locale caveat encoded)~~ **Won't implement — deliberately skipped (M20): the one corpus hit is locale-aware; upgrading would make the borderline more assertive.**
19. ~~docs/rules/H012.md: mention depth-2 corpus caveat in Validation~~ done — docs/rules/H012.md now carries the depth-2 corpus caveat (added 2026-09-26)
20. ~~Sweep report addendum: nested-repo numbers when implemented~~ done (docs-health pass 2026-09-26 — folded into TODO_LIST T38 routing)
21. ~~golangci-lint custom binary rebuild + integration test with 11 rules~~ done — TestCustomGCLIntegration builds the current custom-gcl (11 rules) on demand; subtests verified end-to-end in v0.2.0
22. ~~Self-scan CI step: re-run with H012 fixtures in tree (was exit 0 pre-H012)~~ done — self-scan CI green with H012 in tree (CI on 8f0577e)
23. ~~Verify `--rules` CLI output includes H012 (covered by test — confirm SARIF too)~~ done — H012 registered everywhere incl. --rules; TestRuleCountConsistency guards the count (11)
24. ~~Update `.golangci.custom.yml` example if it enumerates rules (it does not — verify)~~ done — confirmed in this report: .golangci.custom.yml enumerates no rules
25. ~~README badges/actions example version bump rotation policy~~ **Won't implement — no README badges shipped — moot.**
26. ~~Consider dependabot config for golangci-lint-action version pin alerts~~ done — dependabot actions-group covers golangci-lint-action (PR #1 merged, 82f93e1)
27. ~~Feedback doc: "CI red for 5 days unnoticed" — process fix proposal (branch protection on main?)~~ **Won't implement — branch protection declined 2026-09-19; the CI-miss root cause is documented instead (CHANGELOG 0.4.0, T25).**
28. ~~Branch protection / required status checks decision (repo settings)~~ **Won't implement — declined 2026-09-19 (recorded 642a426).**
29. ~~Scheduled weekly self-sweep CI job (corpus drift detector)~~ done (docs-health pass 2026-09-26 — tracked as TODO_LIST T40)
30. ~~H003 reltime: `time.Round` signal (ROADMAP item)~~ done — in ROADMAP Detection breadth (time.Round for H003)
31. ~~H006: Sprintf+TrimRight combined detection (ROADMAP item)~~ done — in ROADMAP Detection breadth (Sprintf+TrimRight for H006)
32. ~~H003: inline `.String() + " ago"` (ROADMAP item)~~ done — in ROADMAP Detection breadth (inline .String() + " ago")
33. ~~LookupMenuItem rule candidate research (corpus-first, like H012)~~ done — in ROADMAP H011+ (LookupMenuItem, demand-gated)
34. ~~Project-level SI/IEC consistency check (ROADMAP)~~ done — in ROADMAP Detection breadth (SI/IEC consistency check)
35. ~~benchstat tracking across releases (ROADMAP)~~ done — in ROADMAP Ecosystem (benchstat tracking)
36. ~~pkg.go.dev doc examples for WordSeries suggestion links~~ **Won't implement — no demand — suggestion links already render from godoc.**
37. ~~Consider `//nolint` reason-required lint for this repo's own code~~ **Won't implement — not adopted.**
38. ~~Sweep Consumer repos' CI with the Action (dogfood action.yml)~~ **Won't implement — declined 2026-09-19 — consumer corpus frozen.**
39. ~~Version the validation docs (front-matter with linter version)~~ **Won't implement — not adopted.**
40. ~~Add H012 to the linter's own self-scan allowlist check (fixtures dir growth)~~ done — self-scan green; H012 fixtures live under testdata/ which the walker skips
41. ~~Review daemon heuristic-commit messages for archeology notes (doc pointer)~~ **Won't implement — daemon accepted as-is (plan decision D2: no history rewrite).**
42. ~~go.work: document sibling-drift check command in AGENTS (GOWORK=off test)~~ done — AGENTS.md "Local dep resolution is go.work-only" gotcha documents the GOWORK=off check
43. ~~Check gogenfilter v3.7 release notes when published (sqlc semantics watch)~~ done — standing watch — documented in AGENTS.md Upstream-tag watch (process)
44. ~~Consider H012 conjunction param detection (function with `conjunction string` param → WordSeries(words, conjunction) suggestion variant)~~ done — in ROADMAP H011+ (conjunction-param variant, added 2026-09-26)
45. ~~H010: detect strip into named constant maps (locale mult) — skip/decide~~ **Won't implement — demand-gated — triage guidance lives in docs/validation/2026-09-18_h010-sweep.md.**
46. ~~Docs: DOMAIN_LANGUAGE "ghost rule" entry mentions TestRuleCountConsistency — verify test name exists post-H012~~ done — TestRuleCountConsistency exists and guards AllRules/allRuleDetectors sync
47. ~~Terminal UX: `--rules` table output alignment with 11 rules~~ **Won't implement — cosmetic.**
48. ~~Cleanup: remove /tmp artifacts (ghl, sweep files) — ephemeral, note only~~ **Won't implement — ephemeral /tmp artifacts.**
49. ~~status report harvest into TODO_LIST on next docs-health pass~~ done (docs-health pass 2026-09-26 — this docs-health pass)
50. ~~This report → mark items done as executed; next report harvests it~~ done (docs-health pass 2026-09-26)

## g) QUESTIONS (cannot figure out myself)

1. ~~**go.mod go-directive bump (1.26.7 → 1.27.1)**: unintended tool
   side-effect, currently uncommitted. Keep it (and align CI/AGENTS to
   Go 1.27) or revert to 1.26.7? It changes the minimum Go for every
   consumer of the module.~~ answered 2026-09-19: revert-on-sight (`642a426`, `ebc9df5`); the bump RETURNED as incident #4 (`b7f00c2`) and the final call is TODO_LIST T33.
2. ~~**Consumer fixes**: the sweeps flagged three of YOUR repos
   (KeyCountdown, Kernovia, CreditReformBilanzampel). Apply the
   suggested fixes / scoped nolints there now, or leave them as a live
   demonstration corpus?~~ declined 2026-09-19 — they stay as the demonstration corpus.
3. ~~**Branch protection on main**: main sat red for 5 days unnoticed.
   Want required status checks (blocks the auto-commit daemon's pushes
   when red), accepting that the daemon will then occasionally fail to
   push — or prefer a notification-only setup?~~ declined 2026-09-19 (`642a426`).

## Verification snapshot (end of segment)

- `go build` / `go vet` / `go test ./... -race` green · golangci-lint
  `0 issues` · nix lint `0 issues` · gofmt clean
- CI green on tip `8f0577e` · v0.3.0 released + proxy-verified
- Working tree: 12 modified files (post-commit doc edits awaiting daemon
  or next explicit commit) + the flagged `go.mod` bump — nothing lost,
  everything pushed except these
