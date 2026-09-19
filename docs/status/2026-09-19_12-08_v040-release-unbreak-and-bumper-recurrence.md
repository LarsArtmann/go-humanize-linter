# Status: v0.4.0 Released — Red-Main Unbreak, Bumper Recurrence #3, Release Verification

**Date:** 2026-09-19 12:08 CEST
**Session span:** 2026-09-19 ~11:00 → 12:05 CEST
**Scope:** This session only — release assessment, red-main unbreak, v0.4.0 cut + verification, doc sync. Per instruction: no re-verification of older claims beyond what gated release decisions. Prior backlog referenced, not re-audited.
**State at writing:** main `e56c535`, CI green, tree clean, tag `v0.4.0` → `6a3267f`, GitHub Release live: https://github.com/LarsArtmann/go-humanize-linter/releases/tag/v0.4.0

---

## a) FULLY DONE

1. **Release need assessed.** Reviewed `[Unreleased]` content (H012 rule, H010 fixture pack, H012 regression guards, H002/H009 suggestion enrichments, CI/security hygiene) + prior report NEXT#1 + the user's "Time for a new release?" → clear MINOR bump to v0.4.0.
2. **Red main root-caused.** CI failed with `go: go.mod requires go >= 1.27 (running go 1.26.8; GOTOOLCHAIN=local)` (run 35433820705). The bump arrived in daemon heuristic commits `49d30ee` (10:56, go 1.26.7→1.27.1) and `e39c390` (11:03, 1.27.1→1.27 + go-finding v1.10.0→v1.12.0 + go-error-family v0.10.0→v0.10.1) — AFTER the user decision recorded 2026-09-19 in TODO_LIST ("Go directive stays 1.26.x, revert-on-sight for background bumps"). Verified via module proxy that go-finding v1.10.0 AND v1.12.0 plus go-linter-sdk v0.3.1 all declare `go 1.26.7` — the dependency bumps bought nothing and did not require the directive change.
3. **go.mod/go.sum reverted to last-green** (`bd3dd03` content: `go 1.26.7`, go-finding v1.10.0, go-error-family v0.10.0). Verified at exact CI parity before pushing: `GOWORK=off GOTOOLCHAIN=go1.26.8` build + vet + full test suite (4/4 packages ok). Landed via daemon commit `e407c0a` (content byte-identical to my revert, verified before push); pushed; CI green (run 35434249197).
4. **Foreign daemon commit investigated before building on it.** `ce52545` (AGENTS/CHANGELOG/TODO_LIST) inspected line-by-line: T32 deploy-key completion records — benign, kept.
5. **CHANGELOG cut.** `[0.4.0] - 2026-09-19` section created from Unreleased; fresh `[Unreleased]` with placeholder categories left behind (commit `6a3267f`).
6. **Version pointers bumped pre-tag.** README Action example and action.yml input description → `@v0.4.0` BEFORE tagging — fixes the v0.3.0 mistake where the frozen tag advertised `@v0.2.0`.
7. **Full pre-release battery green**, including the prior session's two battery gaps: gofmt / no-replace / no-pseudo-version static checks; `GOWORK=off GOTOOLCHAIN=go1.26.8` full suite; `nix run .#build`, `.#vet`, `.#test-race` (4/4 ok), `.#lint` (**0 issues**).
8. **Annotated tag `v0.4.0` created on `6a3267f`** — only after CI went green on that exact commit (run 35434490047). Tag tree verified before push (`git show v0.4.0:go.mod`, `:CHANGELOG.md` section header present).
9. **release.yml green** (run 35434583034). GitHub Release published, not draft, not prerelease, 2 assets (`go-humanize-linter`, `gohumanize`): https://github.com/LarsArtmann/go-humanize-linter/releases/tag/v0.4.0
10. **Release notes curated.** Auto `--generate-notes` listed ONLY the dependabot PR (this repo commits direct to main, so PR-based notes are structurally blind to our features). Replaced via `gh release edit v0.4.0 --notes-file` with user-focused notes (headline H012, upgrader guidance incl. baseline regeneration, compare link).
11. **Consumer-side verification, four independent paths:** proxy `@v` list includes v0.4.0; clean-dir `go get github.com/larsartmann/go-humanize-linter@v0.4.0` OK (13 modules resolved); `go run .../cmd/go-humanize-linter@v0.4.0` on a synthetic H004 hand-roll produced the finding + enriched `english.PluralWord` suggestion + exit 1 (ternary exit codes working from the published module); `go list -m <mod>@latest` → **v0.4.0**; downloaded release binary reports `--version` → `go-humanize-linter v0.4.0` (ldflags baked correctly).
12. **Post-release doc sync committed and green:** TODO_LIST done-record (release + unbreak + verification trail) and AGENTS.md gotchas extended in `e56c535`: (a) bumper recurrence #3 with the NEW fact that its sweep bundles dependency bumps, not just the go directive; (b) nix devShell legitimately stays `go_1_27` for go.work siblings — the 1.26 ceiling applies to the published go.mod only; (c) `--generate-notes` gap + mandatory curated-notes step. CI green (run 35434996338).
13. **go.work verified untracked** (gitignored) — CI builds standalone against published deps; the sibling-vs-published divergence can never leak into CI. Assumption from AGENTS.md confirmed, not assumed.

## b) PARTIALLY DONE

1. **Post-release external verification** — proxy, go get, go run, @latest resolution, binary version: all verified (see a11). Remaining: pkg.go.dev had NOT indexed v0.4.0 by session end (~25 min of 404s; `/fetch` requested repeatedly = indexing queued). Blocker: external propagation lag, nothing on our side. Effort to finish: S (one fetch next session). Already tracked as prior report NEXT#2.
2. **Raw proxy `@latest` JSON** (`proxy.golang.org/.../@latest`) still served v0.3.0 at session end while the go toolchain resolved v0.4.0 correctly. Cosmetic cache lag; expected to flip on its own. Verify next session. Effort: S.
3. **Action install path verified only indirectly.** The Action's exact command is `go install .../cmd/go-humanize-linter@${{ inputs.version }}` (default latest) into GOBIN. I verified module resolution (`go list -m @latest` → v0.4.0) and the binary behavior (`go run @v0.4.0`), but the literal `go install @latest` invocation was blocked by the harness security filter and remains unexecuted end-to-end. Effort: S.
4. **v0.3.1 doc-patch question (frozen README `@v0.2.0` on pkg.go.dev)** — superseded by v0.4.0 by design (README bumped pre-tag). Formally closable once v0.4.0 renders and shows `@v0.4.0`.
5. **Consumer propagation** — user decision (2026-09-19): demo-corpus repos stay untouched; Action `@latest` users get v0.4.0 automatically (verified via go list). Nothing further unless you want consumer repos bumped (declined).
6. **This report** — written strictly from this session's evidence; pre-11:00 state only spot-checked where it gated the release; wider backlog NOT re-verified (per instruction).

## c) NOT STARTED

(Planned/inherited, untouched this session — why: out of scope per instruction, or blocked, or awaiting decision.)

1. pkg.go.dev v0.4.0 render check (H012 docs visible, "Go to latest" flips, README shows @v0.4.0) — blocked on indexing. Still wanted: yes (prior NEXT#2).
2. `BenchmarkH012Detector` (H010 has one; H012 shipped without) — deprioritized vs release. Wanted: yes.
3. Corpus re-sweep under gogenfilter v3.6.1 (pinned this morning, `9c76bd4`; detection table may have moved) — waiting for a sweep session. Wanted: yes (High).
4. Depth-2 nested-repo sweep support — deferred since 09-19 morning report. Wanted: yes.
5. H011 manual-si-parse — deferred (zero corpus demand). Re-evaluate on upstream tags only.
6. H010 NewReplacer-strip gap — documented negative; implement only on corpus demand.
7. Bumper investigation (journal/systemd peek at `project-discovery-daemon`) — user-gated; recurrence CONTINUED this session (incident #3, and its first dep-bundling appearance).
8. Periodic (weekly) self-scan + corpus sweep ritual — proposed twice, never started.
9. `docs/rules/H012.md` explicit negative-shape documentation — wanted (Low).
10. Record consumer-corpus leave-as-is decision in the sweep doc — wanted (Low).
11. Sibling re-pin cadence decision (gogenfilter sibling ahead of v3.6.1; go-finding sibling go.mod at `go 1.27` vs published v1.12.0 at `go 1.26.7` — drift observed today) — wanted (Medium).
12. `nix run .#coverage` post-release confirmation — not in the battery I ran; coverage unmeasured this session.

## d) TOTALLY FUCKED UP

1. **Lost the revert commit to the daemon — AGAIN.** My explicit go.mod/go.sum revert sat uncommitted across two multi-minute verification runs; the daemon won the race and committed it as `e407c0a "chore: auto-commit 2 changed file(s) (heuristic)"`. Content survived 100% (diff-verified before push), but release history now contains another meaningless message at a load-bearing moment (the unbreak). This is the SAME failure the 10:52 report filed as its #1 lesson ("commit-per-batch discipline") — a repeat offense one session later. Severity: history readability, no functional damage. Root cause: I treated "commit after tests pass" as the rule when the daemon's latency demands "commit immediately after the edit, before the tests."
2. **Trusted `--generate-notes` blindly.** The v0.4.0 GitHub Release briefly advertised a feature-free changelog (dependabot PR only) until my Phase 6 check caught it. One `gh release view` of v0.3.0's notes BEFORE tagging would have shown the same defect there and cost 5 seconds. Severity: public-facing docs wrong for ~5 minutes. Mitigation now in place: AGENTS gotcha + curated notes published.
3. **Burned a 12-minute background poll on a banned command.** The pkg.go.dev wait-loop used `curl`, which this harness forbids — all 6 attempts errored, zero datapoints collected. Then I hand-polled ~3 more times over 25 minutes with no stop-rule. Severity: wasted wall-clock, no damage. Correct behavior: 2 attempts max via the fetch tool, then mark pending.
4. **Two edit calls rejected** (README, action.yml) because I edited from bash-read context without View first — avoidable round trips against an explicit harness rule.
5. **Minor tool stumbles:** `gh release view --jq` with the nonexistent `isLatest` field (1 wasted call); no other tool failures.
6. **Nothing user-facing is currently broken:** release verified working, main green, tag correct, docs synced.

## e) WHAT WE SHOULD IMPROVE

1. **Commit-before-verify, not commit-after-verify.** Daemon latency is seconds; test batteries are minutes. New rule for this repo: after any deliberate file change, `git add + commit` IMMEDIATELY, then run the battery. (Second occurrence across consecutive sessions — if it happens a third time, this becomes a hard-coded habit check in every handoff.)
2. **Kill `--generate-notes` or make curation mandatory in the flow.** Every release from this repo will publish a feature-free notes body until manually edited. Concrete options: (a) release.yml switches to `--notes-file` from a repo file, (b) a release-runbook checklist makes `gh release edit` a named, non-skippable step. AGENTS gotcha currently documents the gap; the workflow itself is unchanged.
3. **Script the post-release external checks.** `scripts/post-release-check.sh` doing bounded polls (proxy version list, clean-dir go get, `go list @latest`, pkg.go.dev) with allowed tooling would have saved ~30 minutes of hand-rolled polling this session and standardizes the NEXT#2-style follow-ups.
4. **Structural bumper defense beats vigilance.** Incident #3 happened *hours* after the revert-on-sight decision. A CI guard step (`go.mod` `go` directive must be 1.26.x, else fail with an explanatory message) converts silent toolchain drift into a loud, self-documenting failure — and encodes the policy in code instead of in TODO_LIST prose.
5. **Inspect the previous instance before reusing any mechanism** (release notes, workflows, scripts). Prior output is the cheapest code review.
6. **View-before-edit is not optional** — the two rejected edits were pure process debt.
7. **What worked, keep doing:** trust-order CLI > LSP (zero LSP time burned this session); verify-at-CI-parity (`GOWORK=off GOTOOLCHAIN=go1.26.8`) before pushing any go.mod-touching change; investigate foreign commits before building on them (ce52545 checked line-by-line).

## f) NEXT — ranked (up to 50 requested; 24 concrete items listed, plus the inherited backlog pointer. Impact / Effort / Category per item for HARVEST.)

**Release verification (close out v0.4.0)**

| #  | Task                                                                                                                         | Impact | Effort | Category      |
| -- | ---------------------------------------------------------------------------------------------------------------------------- | ------ | ------ | ------------- |
| 1  | Verify pkg.go.dev indexes v0.4.0: H012 docs render, "Go to latest" flips, README shows @v0.4.0 (closes the v0.3.1 question)   | High   | S      | Documentation |
| 2  | Verify raw proxy `@latest` JSON flipped to v0.4.0                                                                             | Low    | S      | Quality       |
| 3  | Run the Action's literal install path (`go install ...@latest`) in a clean env to close the b3 caveat                        | Medium | S      | Quality       |
| 4  | Write `scripts/post-release-check.sh` (bounded proxy/go-get/pkg.go.dev checks)                                                | Medium | M      | Quality       |
| 5  | Sweep the consumer demo repos with the released v0.4.0 binary — confirm H012's 2 corpus TPs fire from the published artifact   | Medium | S      | Quality       |
| 6  | Run `nix run .#coverage` post-release; confirm no coverage drop from H012 additions                                           | Low    | S      | Quality       |

**Bumper defense (incident #3 happened today)**

| #  | Task                                                                                                                         | Impact | Effort | Category      |
| -- | ---------------------------------------------------------------------------------------------------------------------------- | ------ | ------ | ------------- |
| 7  | Add CI guard step: fail loud if go.mod `go` directive != 1.26.x (policy in code, not prose)                                   | High   | S      | Quality       |
| 8  | Decide bumper investigation (journal/systemd peek at project-discovery-daemon) — recurrence continues, now bundling dep bumps  | Medium | S      | Cleanup       |
| 9  | Decide whether go-finding v1.12.0 + go-error-family v0.10.1 get re-applied deliberately (with CHANGELOG entry) or stay reverted | Medium | S      | Cleanup       |
| 10 | Watch the next daemon window; a 4th bump within a week promotes #7 from "should" to "must"                                    | High   | S      | Process       |

**Release process**

| #  | Task                                                                                                                         | Impact | Effort | Category      |
| -- | ---------------------------------------------------------------------------------------------------------------------------- | ------ | ------ | ------------- |
| 11 | Switch release.yml to curated notes (repo notes file) OR codify `gh release edit` as a named release-runbook step             | Medium | S      | Cleanup       |
| 12 | Decide sibling re-pin cadence (gogenfilter sibling ahead of v3.6.1; go-finding sibling at go 1.27 vs published 1.26.7)         | Medium | S      | Cleanup       |
| 13 | Document the intentional go.mod(1.26.7)/go.work-sibling(1.27+) floor divergence in AGENTS when siblings next release           | Low    | S      | Documentation |
| 14 | Add the weekly self-scan + corpus-sweep ritual (CI schedule or checklist)                                                     | Medium | S      | Quality       |

**Detection/rules (inherited from morning report, untouched)**

| #  | Task                                                                                                                         | Impact | Effort | Category      |
| -- | ---------------------------------------------------------------------------------------------------------------------------- | ------ | ------ | ------------- |
| 15 | `BenchmarkH012Detector` mirroring the H010 benchmark shape                                                                    | Medium | M      | Quality       |
| 16 | Corpus re-sweep across ~/projects under gogenfilter v3.6.1 (pinned today; table may have moved)                               | High   | M      | Quality       |
| 17 | Depth-2 nested-repo sweep support (games/*-style layouts)                                                                     | Medium | M      | Feature       |
| 18 | H011 manual-si-parse — re-evaluate only on new upstream go-humanize tags or corpus demand                                     | Low    | L      | Feature       |
| 19 | H010 NewReplacer strip detection — implement only on corpus demand (negative fixture already exists)                          | Low    | M      | Feature       |

**Docs**

| #  | Task                                                                                                                         | Impact | Effort | Category      |
| -- | ---------------------------------------------------------------------------------------------------------------------------- | ------ | ------ | ------------- |
| 20 | HARVEST this report's (f) into TODO_LIST.md / ROADMAP.md (else these die in this timestamped file)                            | High   | S      | Documentation |
| 21 | `docs/rules/H012.md`: document the four negative shapes explicitly (Sprintf prose, bare variable, slash join, prefix-no-conj) | Low    | S      | Documentation |
| 22 | Record the consumer-corpus leave-as-is decision in docs/validation so future sweeps don't "fix" them                          | Low    | S      | Documentation |
| 23 | Re-verify the 2026-09-18 50-item backlog before acting on any of it (stale by a day; never re-checked)                        | Medium | M      | Documentation |
| 24 | Update FEATURES.md if H012/guards rows need status flips post-0.4.0                                                           | Low    | S      | Documentation |

(Inherited, not re-verified: the full prior backlog lives in `docs/status/2026-09-18_22-01_*.md` and the morning report `2026-09-19_10-52_*` NEXT list.)

## g) QUESTIONS (decisions only you can make; what I tried is noted)

1. **Dependency posture after the revert:** go-finding v1.12.0 and go-error-family v0.10.1 were swept into the bumper commits and reverted with them. I verified both publish against `go 1.26.7`, so re-applying is safe — but I cannot know whether you wanted those bumps at all. Deliberate re-apply with a CHANGELOG entry, or stay pinned at v1.10.0/v0.10.0 until something needs them?
2. **Bumper countermeasure:** incident #3 struck hours after your revert-on-sight decision. Do you want the CI guard step (fail loud on `go` directive != 1.26.x) — and/or is now the time for the previously user-gated journal/systemd investigation into `project-discovery-daemon`?
3. **Release-notes mechanism:** keep `--generate-notes` + a mandatory manual `gh release edit` per release (now documented in AGENTS), or switch release.yml to a committed notes file so the GitHub body is correct from the moment of creation?

---

*Report format note: written as Markdown per explicit user instruction (`docs/status/*.md`), overriding the status-report skill's HTML default.*
