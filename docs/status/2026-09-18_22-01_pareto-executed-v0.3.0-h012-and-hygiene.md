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

- **M21 LSP hygiene** — never restarted gopls; diagnostics were stale the
  entire session (correctly ignored in favor of build/test/CI truth).
- **Local vs CI lint parity** — root-caused only partially (local nix
  golangci-lint 2.13.2 vs CI 2.12.2); the actions bump merged via PR #1
  may have aligned versions — NOT re-verified post-merge.
- **Silent nix-lint wrapper failure** (exit 1, empty output) observed
  twice mid-session; self-resolved after the underlying findings were
  fixed; never root-caused.

## c) NOT STARTED

- **T18** plugin-index submission (foreign-repo PR; unblocked by v0.3.0).
- **T20** gogenfilter sibling release + pin bump — gated on approval.
- **T32** read-only deploy-key removal in the 4 dependency repos
  (foreign repos).
- pkg.go.dev rendering check for v0.3.0.

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

1. H012 negative fixture: comma join + prose-only conjunction (guards the position filter)
2. Decide/revert-or-keep the go.mod `go 1.27.1` bump; align CI setup-go version if kept
3. Verify post-merge CI golangci-lint version vs local nix (lint parity)
4. Root-cause the silent nix-lint wrapper failure (exit 1, empty output)
5. T18: golangci-lint plugin-index PR (verify-before-filing + github-voice)
6. T20: gogenfilter release + pin bump + GOWORK=off sweep [needs approval]
7. T32: remove read-only deploy keys from the 4 dep repos
8. pkg.go.dev check for v0.3.0 (docs render, 11 rules listed)
9. Nested-repo corpus sweep (depth-2 go.mod discovery) for H010/H012
10. Apply the linter's own advice to consumer repos (KeyCountdown joinWords → humanize.WordSeries; Kernovia Join(x," and ") → WordSeries; CreditReformBilanzampel → scoped nolint)
11. gopls restart; confirm diagnostics freshness
12. H012 micro-benchmark (match BenchmarkH010Detector pattern)
13. H010/H012 analysistest scoped-nolint variants
14. CHANGELOG: Unreleased section will need a 0.4.0 cut when H012 ships
15. Plan next release cadence (0.4.0 after H012 soak?)
16. ROADMAP: revisit SI-parse (H011) after upstream go-humanize activity
17. Consider NewReplacer strip form for H010 (F13, still open)
18. Consider Full-tier for H010 strip+parse when function is exactly the library signature AND no separator config (refine the skipped M20 with the locale caveat encoded)
19. docs/rules/H012.md: mention depth-2 corpus caveat in Validation
20. Sweep report addendum: nested-repo numbers when implemented
21. golangci-lint custom binary rebuild + integration test with 11 rules
22. Self-scan CI step: re-run with H012 fixtures in tree (was exit 0 pre-H012)
23. Verify `--rules` CLI output includes H012 (covered by test — confirm SARIF too)
24. Update `.golangci.custom.yml` example if it enumerates rules (it does not — verify)
25. README badges/actions example version bump rotation policy
26. Consider dependabot config for golangci-lint-action version pin alerts
27. Feedback doc: "CI red for 5 days unnoticed" — process fix proposal (branch protection on main?)
28. Branch protection / required status checks decision (repo settings)
29. Scheduled weekly self-sweep CI job (corpus drift detector)
30. H003 reltime: `time.Round` signal (ROADMAP item)
31. H006: Sprintf+TrimRight combined detection (ROADMAP item)
32. H003: inline `.String() + " ago"` (ROADMAP item)
33. LookupMenuItem rule candidate research (corpus-first, like H012)
34. Project-level SI/IEC consistency check (ROADMAP)
35. benchstat tracking across releases (ROADMAP)
36. pkg.go.dev doc examples for WordSeries suggestion links
37. Consider `//nolint` reason-required lint for this repo's own code
38. Sweep Consumer repos' CI with the Action (dogfood action.yml)
39. Version the validation docs (front-matter with linter version)
40. Add H012 to the linter's own self-scan allowlist check (fixtures dir growth)
41. Review daemon heuristic-commit messages for archeology notes (doc pointer)
42. go.work: document sibling-drift check command in AGENTS (GOWORK=off test)
43. Check gogenfilter v3.7 release notes when published (sqlc semantics watch)
44. Consider H012 conjunction param detection (function with `conjunction string` param → WordSeries(words, conjunction) suggestion variant)
45. H010: detect strip into named constant maps (locale mult) — skip/decide
46. Docs: DOMAIN_LANGUAGE "ghost rule" entry mentions TestRuleCountConsistency — verify test name exists post-H012
47. Terminal UX: `--rules` table output alignment with 11 rules
48. Cleanup: remove /tmp artifacts (ghl, sweep files) — ephemeral, note only
49. status report harvest into TODO_LIST on next docs-health pass
50. This report → mark items done as executed; next report harvests it

## g) QUESTIONS (cannot figure out myself)

1. **go.mod go-directive bump (1.26.7 → 1.27.1)**: unintended tool
   side-effect, currently uncommitted. Keep it (and align CI/AGENTS to
   Go 1.27) or revert to 1.26.7? It changes the minimum Go for every
   consumer of the module.
2. **Consumer fixes**: the sweeps flagged three of YOUR repos
   (KeyCountdown, Kernovia, CreditReformBilanzampel). Apply the
   suggested fixes / scoped nolints there now, or leave them as a live
   demonstration corpus?
3. **Branch protection on main**: main sat red for 5 days unnoticed.
   Want required status checks (blocks the auto-commit daemon's pushes
   when red), accepting that the daemon will then occasionally fail to
   push — or prefer a notification-only setup?

## Verification snapshot (end of segment)

- `go build` / `go vet` / `go test ./... -race` green · golangci-lint
  `0 issues` · nix lint `0 issues` · gofmt clean
- CI green on tip `8f0577e` · v0.3.0 released + proxy-verified
- Working tree: 12 modified files (post-commit doc edits awaiting daemon
  or next explicit commit) + the flagged `go.mod` bump — nothing lost,
  everything pushed except these
