# Status Report — 2026-09-19 12:17 CEST — T32 Deploy-Key Removal Session

**Session scope:** T32 — remove the read-only CI deploy keys from the 4 dependency
repos (go-linter-sdk, go-finding, gogenfilter, go-error-family), plus the doc
updates and self-review fixes that fell out of it. Per instruction, this report
covers ONLY this session's run and what it directly noticed — no unrelated
research. The v0.4.0 work observed in git history was made by ANOTHER session;
reported as observation, not investigated.

---

## What This Session Did (chronological)

1. **T32 execution, end-to-end** (github-voice: none needed; pure `gh api` work):
   - Verified `gh` auth (LarsArtmann, `repo` + `admin:org` scopes).
   - Listed deploy keys in all 4 dep repos: **exactly one key each**,
     `ci-deploy-<repo>-readonly (auto 2026-09-02)`, all `read_only: true`,
     IDs `162081806` / `162081807` / `162081810` / `162081811`.
   - Pre-deletion safety reasoning: single key per repo + title matching the
     documented per-repo CI machinery ⇒ no other consumer could depend on them;
     deploy keys are re-addable (reversible operation).
   - Deleted all 4 (API success × 4).
   - Re-listed: **0 deploy keys remain in all 4 repos** (verified via API).
2. **Doc updates:** TODO_LIST.md (T32 → Done, detail section removed),
   AGENTS.md Gotchas bullet (keys deleted 2026-09-19), CHANGELOG.md
   (security-hygiene entry under Unreleased → Changed). Committed by the
   auto-commit daemon as `ce52545`.
3. **Self-review caught a split brain I created:** AGENTS.md had a SECOND
   deploy-key mention (the "going public" paragraph) still saying
   "Remaining cleanup (pending): the 5 `DEPLOY_KEY_*` secrets can be deleted…"
   — stale since T31 (secrets deleted 09-18). My first edit sweep updated only
   the Gotchas bullet. Fixed the stale paragraph; verified with
   `grep -n pending AGENTS.md` → zero matches. Committed as `2bb893e`.
4. **Post-completion external verifications:**
   - `gh secret list` on go-humanize-linter → **empty** (T31 independently
     re-confirmed; zero secrets remain).
   - Local leftover check: the three global `url.*.insteadOf` git rewrites map
     https→SSH via Lars's **personal** `id_ed25519` (still present in
     `~/.ssh`) — they do NOT reference the deleted deploy keys. No local
     breakage from the deletions; nothing to clean.
   - CI on main: latest visible run (12:13 CEST) **success**; the v0.4.0
     Release workflow run (11:23 CEST) **success**.

---

## a) FULLY DONE

| # | Item | Evidence |
|---|------|----------|
| 1 | T32: all 4 dep-repo read-only deploy keys removed | API re-list = `0` keys × 4 repos |
| 2 | T32 doc trail: TODO_LIST (Done), AGENTS.md gotcha, CHANGELOG entry | commit `ce52545` (3 files) |
| 3 | AGENTS.md split-brain fixed ("pending" claim removed) | commit `2bb893e`; `grep pending` = 0 matches |
| 4 | T31 external re-verification: zero `DEPLOY_KEY_*` secrets | `gh secret list` = empty |
| 5 | Local-machine impact check: no insteadOf/ssh references to deleted keys | `git config --global --get-regexp` + `ls ~/.ssh` |
| 6 | Going-public auth cleanup arc now COMPLETE (workflows 09-09 → secrets 09-18 → dep keys 09-19) | all three legs verified |

## b) PARTIALLY DONE

1. **First-pass doc consistency** — I updated 1 of 2 AGENTS.md deploy-key
   mentions and declared done. The second (line ~138) said "pending" while the
   first said "deleted". Caught in self-review minutes later, fixed, verified —
   but the miss happened. Now fully consistent.
2. **Section (f) → TODO_LIST harvest** — this report's next-task list is NOT yet
   harvested into TODO_LIST.md/ROADMAP.md. Deliberate: TODO_LIST only admits
   work "verified against the code as not-yet-done", and most section-(f) items
   are unverified brainstorm needing your triage first (docs-health HARVEST
   would apply routing rigor — awaiting your go-ahead per "then wait").

## c) NOT STARTED (session-scoped leftovers)

1. Harvest of section (f) below into TODO_LIST.md / ROADMAP.md (needs user
   triage of the brainstorm items).
2. Nothing else — T32 itself has no remaining checkboxes; the task is closed
   end-to-end.

## d) TOTALLY FUCKED UP

**Nothing.** No code was touched this session (docs + GitHub settings only);
main's CI is green; all deletions verified zeroed-out and were reversible in
any case.

Honest near-misses worth recording (not fucked up, but instructive):

1. **The split brain I created and fixed** (AGENTS.md "pending" vs "deleted").
   If uncaught, a future session would have re-pursued already-completed
   cleanup — exactly the "status reports are point-in-time" trap the repo has
   been burned by before.
2. **Noticed, not caused by this session (not researched, per instruction):**
   - Commit `e407c0a` (11:14 CEST): go.mod + go.sum changed by a tool/other
     session, swept by the daemon. Sits between my docs commit and the v0.4.0
     release commits — likely deliberate release prep by the other session;
     CI green, so flagging for awareness only. Per the AGENTS gotcha, unintended
     `go`-directive/dependency bumps are a known hazard class.
   - `ff1f074` and the three v0.4.0 commits (`6a3267f`, `e56c535`, tag `v0.4.0`
     released 11:23 CEST): all from another session. My CHANGELOG entry landed
     under `[Unreleased]` after the v0.4.0 section — no conflict.

## e) WHAT WE SHOULD IMPROVE

1. **Post-edit mention sweep:** after any consistency edit, grep the WHOLE
   file (and ideally docs/) for the changed concept before declaring done.
   Line 217 was fixed; line 138 was 80 lines away and missed. Mandatory step,
   not optional polish.
2. **Surface decision tensions out loud:** TODO_LIST records "no foreign-repo
   actions for now" (2026-09-19 decline of T18/T20), and T32 IS a foreign-repo
   action. I proceeded correctly on your explicit assignment but never stated
   the tension while working. Silence about known conflicts is how split
   brains start.
3. **Export evidence before destructive ops:** I captured key metadata in the
   transcript, but a one-line JSON export to a scratch file before deletion
   would make audits trivial. Cheap insurance, skipped.
4. **Batch pre-flight checks up front:** `gh secret list` and local
   insteadOf/ssh checks ran AFTER the deletions as "verification". Cheaper and
   safer to batch them before.
5. **CHANGELOG scope policy is undefined:** today's entry documents foreign-repo
   settings hygiene in the linter's user-facing CHANGELOG. Defensible, but
   nobody decided that policy — it just happened. Needs a ruling (see g2).

## f) Up to 50 Things We Should Get Done Next

**Tag legend:** `[V]` verified-today evidence · `[K]` known-open from loaded
project context (not re-verified this session) · `[D]` declined by your
2026-09-19 decision (listed for completeness, act only if reversed) · `[B]`
brainstorm — needs verification/routing before entering TODO_LIST.

### Immediate / decisions

1. `[V]` Triage this section (f): route verified items into TODO_LIST, the rest into ROADMAP (docs-health HARVEST).
2. `[V]` Rule on the foreign-repo standing policy (question g1).
3. `[V]` Rule on CHANGELOG scope for foreign-repo hygiene (question g2).
4. `[V]` Rule on status-report archive/prune policy (question g3).

### Project items known-open (from context)

5. `[K]` T29 H011 manual-si-parse — deferred; revisit trigger = next dustin/go-humanize upstream tag.
6. `[K]` gogenfilter sibling (go.work) vs pinned v3.6.0 parity check (`GOWORK=off go test ./...` gate) — documented drift risk.
7. `[K]` Upstream-tag watch: diff the next go-humanize tag, enumerate new public API → new rule candidates.
8. `[K]` Consumer-corpus research: hand-rolled `BytesN`/`IBytesN` (v1.1.0 surface) — potential new rule, zero demand confirmed so far.
9. `[K]` Remove gogenfilter legacy-suffix fallback IF upstream ever adds `_gen.go` patterns (conditional).
10. `[K]` Plan GOEXPERIMENT=jsonv2 removal once encoding/json/v2 is default in a Go release.
11. `[K]` Verify v0.4.0 propagation: pkg.go.dev listing + README install/tag references (release happened 11:23 today).
12. `[K]` Verify `--explain` / `--rules` list H012 (registered-everywhere claim, cheap self-check).
13. `[K]` RULE-doc cross-link audit: H002↔H010 exists; check H012's doc cross-links consistently.
14. `[K]` ADR existence check for the H012 conjunction-position filter decision (AGENTS documents it; ADR file unverified).
15. `[K]` ROADMAP review: strike-or-keep the deferred-H011 entry at next upstream tag.

### Process / CI

16. `[K]` CI guard from the "scanner must prove it scanned" lesson: assert measured-input > 0 in the custom-gcl integration test.
17. `[K]` Baseline regeneration runbook as a checked procedure for new-rule releases (README mentions it; not a checklist).
18. `[K]` Exit-2 triage runbook: document the policy for medium/low-confidence findings in CI.
19. `[K]` Pre-tag ritual: "never tag while CI red" is documented but unenforced (branch protection declined) — make it a checklist step in go-release flow.
20. `[K]` Auto-commit daemon "heuristic" commits blur history (visible again today); evaluate commit-message conventions or daemon tuning — policy call, affects all repos.
21. `[K]` Periodic org-wide deploy-key inventory script (this session's check, made reusable) so key hygiene is one command.
22. `[K]` Dependency sweep: check for newer go-linter-sdk / go-finding releases (v0.2.0+ required; latest state unknown to this session).

### Testing / quality

23. `[B]` SARIF output schema validation (never explicitly validated).
24. `[B]` `--save-baseline` + `--behavior-delta` end-to-end integration test.
25. `[B]` Fuzz `normLit` / `byteUnitRegex` edges (digit separators, unicode, case).
26. `[B]` Walker performance benchmark on a large tree (single-read design — measure it).
27. `[B]` Verify release.yml test step includes `-race` (nix has test-race; release flow unknown).
28. `[B]` `--format json` machine-contract test (schema snapshot to catch breaking changes).
29. `[B]` Institutionalize FP-source fixtures: every fixed FP class becomes a permanent negative fixture (H012 did this; make it the rule).
30. `[B]` Confidence-calibration study: sample real findings, measure tier accuracy (Full/High/Medium actually meaning what we claim).
31. `[B]` Windows CI leg (walker path handling; dev platform is linux-only today).

### Docs

32. `[K]` FEATURES.md freshness: confirm H012 (and v0.4.0 changes) are listed.
33. `[B]` docs/DOMAIN_LANGUAGE.md existence/content check (referenced by the docs framework; presence unverified).
34. `[B]` CONTRIBUTING spot-check: no stale private-dep/deploy-key instructions remain (cleaned 09-09; one grep to confirm).
35. `[K]` docs-health ANNOTATE pass over the oldest status reports once their items complete.
36. `[B]` docs/rules/ README or index documenting the H011 reservation (currently only in AGENTS.md).

### Ecosystem (declined by default — recorded only)

37. `[D]` T18 golangci-lint plugin-index submission.
38. `[D]` T20 gogenfilter release.
39. `[D]` Fix consumer repos (KeyCountdown, Kernovia, CreditReformBilanzampel) — they stay as the demonstration corpus.
40. `[D]` Mirror the deploy-key removal note into each dependency repo's own docs.

### Deeper improvements (brainstorm)

41. `[B]` `--min-confidence` default policy: `low` shows everything; consider `medium` as CI-friendly default (opt-in to noise).
42. `[B]` `--verify-suppressions` on by default in CI (currently opt-in) — stale-suppression rot is exactly what it catches.
43. `[B]` Baseline JSON format versioning for forward compatibility.
44. `[B]` Suppression-namespace ergonomics: H0SUP flags wrong-namespace directives; consider also accepting module-path aliases with a warning instead of silently no-oping upstream.
45. `[B]` Rule-ID stability contract documented for downstream (which IDs may disappear, e.g. reserved H011).
46. `[B]` README floor note: ParseComma/ParseCommaf suggestions require go-humanize v1.1.0+ (consumers on v1.0.x can't follow them).
47. `[B]` Corpus re-sweep schedule: H012 swept same-day, H010 on 09-18 — define a cadence for all-rule sweeps after rule edits.
48. `[B]` Ritual item: `nix run .#test` + `nix run .#lint` before the next release (not run this session — docs-only, listed as process hygiene).
49. `[B]` Consider annotating (not unifying) the CLI-vs-plugin H0SUP position difference in docs/rules or an ADR, so nobody "fixes" it without the context AGENTS.md warns about.
50. `[B]` Review whether `[Unreleased]` CHANGELOG entries should auto-flow into the next release section via a cut script (v0.4.0's cut was manual today).

## g) Questions I Can NOT Figure Out Myself

1. **Foreign-repo standing rule:** TODO_LIST records your 2026-09-19 decline
   ("no foreign-repo actions for now", re T18/T20), yet T32 — also a
   foreign-repo action — was assigned today and I executed it. Going forward:
   does an explicit task assignment always override the blanket decline, or
   should I re-confirm with you before EVERY foreign-repo touch (settings,
   issues, PRs, releases)?
2. **Scope policies, two related rulings:** (a) should docs/status reports
   keep accumulating forever or be pruned/archived after annotation?
   (b) should the linter's CHANGELOG record foreign-repo-only hygiene (like
   today's deploy-key entry) or stay strictly about the linter artifact?
3. **Personal git rewrites:** your global git config rewrites ALL https GitHub
   URLs to SSH (`git@github.com:`) via your personal `id_ed25519`. Everything
   involved is now public and anonymously fetchable over https. Keep the
   rewrites as-is (personal preference, push still needs SSH anyway), or
   narrow them so pure fetch workflows no longer depend on your key?

---

**Point-in-time snapshot.** Written 2026-09-19 12:17 CEST · session covered
T32 only · v0.4.0 observed as released (11:23 CEST) by another session.
