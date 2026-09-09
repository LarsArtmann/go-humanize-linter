# Status Report: Going Public + Private-Auth Legacy Cleanup

**Date:** 2026-09-09 02:55 CEST
**Session scope:** Assess → execute "make go-humanize-linter public" including dependency verification, visibility flip, legacy `GOPRIVATE`/deploy-key removal, and full local verification.
**Repo state at report time:** `main` at `393b7fe`, working tree clean, **4 commits unpushed**, **CI on `origin/main` RED (billing)**, repo **PUBLIC**.

---

## 0. Executive Summary

The repo was flipped to public after verifying all 4 sibling dependencies were already public — the initially-assumed "private dependency blocker" **did not exist** (the session's first, wrong assessment was corrected by the user and re-verified empirically). The entire private-auth apparatus (`GOPRIVATE`, `GONOSUMDB`, per-repo CI deploy keys) was removed from workflows, flake, docs, and tests. Everything verifies green locally, including a fresh-consumer simulation (`GOWORK=off`, no auth, public-proxy-only).

Two important discoveries during execution:

1. **CI on `origin/main` is red for a non-code reason**: GitHub Actions billing failure ("recent account payments have failed or your spending limit needs to be increased"). Jobs never start. This predates the session (last failed run 2026-09-08 19:34 UTC, commit `b98aabf`).
2. **`flake.nix` `vendorHash` had silently drifted** (deps were bumped to `go-finding` v1.8.0 etc. without updating the hash). Proven pre-existing via A/B test against `HEAD:flake.nix`, then fixed.

---

## a) FULLY DONE

| #   | Item                                                                                                                                                                                                        | Evidence                                                                                         |
| --- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------ |
| A1  | Dependency visibility verified via GitHub API — all 4 deps (`go-linter-sdk`, `go-finding`, `gogenfilter`, `go-error-family`) `private=false`                                                                | `gh api repos/*`                                                                                 |
| A2  | Public module-proxy resolution verified in clean env (no `GOPRIVATE`, no workspace, no auth) — all 4 pinned dep versions                                                                                    | `go get` in `/tmp/pub-dep-check`                                                                 |
| A3  | Repo flipped to **public**; anonymous API access confirmed                                                                                                                                                  | `gh repo edit --visibility public`; unauthenticated `api.github.com` 200                         |
| A4  | Repo metadata set: description + topics (`go`, `linter`, `golangci-lint`, `humanize`, `static-analysis`)                                                                                                    | `gh repo edit`                                                                                   |
| A5  | Anonymous `go get github.com/larsartmann/go-humanize-linter@v0.2.0` works — the module is publicly consumable                                                                                               | temp module `go get`                                                                             |
| A6  | Deploy-key machinery removed from **both** CI jobs in `ci.yml` and from `release.yml`; `GOPRIVATE` env dropped (`GOEXPERIMENT=jsonv2` kept)                                                                 | diff; YAML validated via `yaml.v3`                                                               |
| A7  | `GOPRIVATE`/`GONOSUMDB` removed from `flake.nix` (3 places), `README.md`, `CONTRIBUTING.md` (2 places), `action.yml`, `cmd/go-humanize-linter/main_test.go`, `plugin/plugin_integration_test.go` (2 places) | repo-wide grep now only hits historical docs, CHANGELOG history, and explanatory AGENTS.md notes |
| A8  | Pre-existing `vendorHash` drift diagnosed (A/B vs `HEAD:flake.nix` → identical mismatch) and fixed                                                                                                          | `nix flake check` → "all checks passed!"                                                         |
| A9  | Full local verification: `nix run .#test-race` (4/4 packages ok), `nix run .#vet` clean, `nix run .#lint` → 0 issues                                                                                        | CI-equivalent commands                                                                           |
| A10 | Fresh-user simulation: `GOWORK=off` + no `GOPRIVATE` + public-proxy-only → build + full test suite green                                                                                                    | explicit env-stripped run                                                                        |
| A11 | `gofmt` clean on all edited Go files; `.direnv` confirmed untracked (stale profile self-heals on next direnv eval)                                                                                          | `gofmt -l`, `git ls-files`                                                                       |
| A12 | AGENTS.md corrected twice: stale "GOPRIVATE REQUIRED" claims replaced with verified reality; CI gotcha updated to "machinery REMOVED (2026-09-09)" with gated follow-ups documented                         | `AGENTS.md:122-131`, `AGENTS.md:205-206`                                                         |
| A13 | Gated secret deletion designed deliberately: deleting `DEPLOY_KEY_*` before cleaned workflows land would break CI on any pre-cleanup push — documented exact commands instead                               | final report + AGENTS.md                                                                         |

## b) PARTIALLY DONE

| #  | Item                           | Done                                                     | Missing                                                                                                                                                                                                                                                                                                      |
| -- | ------------------------------ | -------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| B1 | Legacy auth removal end-to-end | Workflows/flake/docs/tests cleaned and committed locally | **4 commits unpushed**; 5 `DEPLOY_KEY_*` secrets still exist; unused deploy keys still registered in the 4 dependency repos                                                                                                                                                                                  |
| B2 | Public presence                | Repo public, description + topics set                    | No README badges (CI/Codecov/pkg.go.dev/license), no social preview, pkg.go.dev not yet indexed (still 404 ~2.5h after trigger)                                                                                                                                                                              |
| B3 | CI health                      | Workflows simplified to zero-auth                        | CI cannot validate anything while Actions billing is broken; the cleaned workflows have **never run**                                                                                                                                                                                                        |
| B4 | pkg.go.dev readiness           | Page request submitted (triggers indexing)               | Still 404. Additional unverified risk: pkg.go.dev builds docs with a **stock toolchain** — if `go-finding`'s `encoding/json/v2` usage truly requires `GOEXPERIMENT=jsonv2`, server-side doc builds may fail even after indexing (would surface as "build failed / no docs"). Needs confirmation once indexed |
| B5 | AGENTS.md accuracy             | Active sections now truthful                             | Historical `docs/status/*` + `docs/planning/*` still contain "GOPRIVATE required" claims and **names of other private projects** (left as immutable point-in-time snapshots, per docs-health philosophy — but this is now a _public_ surface)                                                                |
| B6 | Flake hygiene                  | `vendorHash` fixed for current `go.sum`                  | No CI guard added, so the same silent-drift class can recur on the next dep bump                                                                                                                                                                                                                             |

## c) NOT STARTED

- N1. Pushing the 4 local commits (requires explicit user ask; see Question 2)
- N2. Deleting the 5 `DEPLOY_KEY_*` secrets (gated on green CI, which is gated on billing)
- N3. Removing the unused read-only deploy keys from the 4 dependency repos
- N4. README public-presence work: badges, `go install` quick-start **with the `GOEXPERIMENT=jsonv2` caveat**, `golangci-lint custom` copy-paste section
- N5. Public-repo community files: issue templates, PR template, SECURITY.md, (optional) Discussions/FUNDING
- N6. CI additions: `nix flake check` step (vendorHash drift guard), anonymous-consumer build step (`GOWORK=off`), `go install @latest` smoke test
- N7. Release hardening: SHA256SUMS file, more architectures (release currently ships linux-amd64 only), optional cosign signing
- N8. Dependabot/renovate config (the vendorHash drift happened because dep bumps were manual and flake wasn't updated in the same commit)
- N9. Announcement/website/demo (website-launch territory) — not requested, listed for completeness
- N10. `.envrc`/direnv regeneration (self-heals on next shell entry; not forced during session)

## d) TOTALLY FUCKED UP

1. **The session's opening assessment was flat-out wrong.** I declared "4 private dependencies — a public repo nobody can build is worse than private" based on **stale AGENTS.md claims alone, without a single verification call**. The user had to correct me ("I think there are already public"). One `gh api` loop (10 seconds) falsified it. If unchallenged, the recommendation would have been the exact opposite of correct. This is the verify-before-claiming failure mode the memory rules explicitly warn about, and it happened _because_ the docs said "CRITICAL: GOPRIVATE REQUIRED" and I trusted the docs over reality.
2. **Missed that CI on main was red for the entire session.** I built a completion gate ("push → green CI → delete secrets") without ever running `gh run list`. The gate was unsatisfiable as stated: CI fails before any job starts due to **Actions billing**, a cause I have zero ability to fix. A status check at session start would have surfaced this and changed the plan.
3. **A corrupted edit payload was nearly applied.** One `multiedit` call to `plugin/plugin_integration_test.go` contained garbage text (`построить` — a token-glitch inside the replacement string) and had not been preceded by a proper read. The read-first guard rejected it; a `view` + clean retry fixed it. No damage, but the failure class is real: garbage-in replacement strings on files not read via the proper tool.

## e) WHAT WE SHOULD IMPROVE

1. **Verify external/state claims before asserting them** — repo visibility, CI status, tool availability. The first message of this session was wrong and the second-to-last discovery (billing) was late. Both were one-command checks.
2. **Stop writing "REQUIRED" into docs for things that are convenience.** `GOPRIVATE` was listed as CRITICAL-REQUIRED in AGENTS.md/README/CONTRIBUTING long after it stopped being necessary. Stale REQUIRED claims actively misdirect future agent sessions (proof: this session's wrong start).
3. **Catch drift classes in CI, not by luck.** `vendorHash` drifted silently on main; CI never checks `nix flake check`. Add it (F-item 9).
4. **Check CI/push state as part of pre-work, not as an afterthought** — especially before designing gated multi-step completions.
5. **Billing is a single point of failure for the whole ecosystem** — every LarsArtmann repo's CI is currently dead. Public repos get free Actions; this repo now qualifies, but the account-level failure annotation may still block runs.
6. **Session hygiene was otherwise good**: every change was verified immediately (build → tests → lint → flake check → consumer simulation), edits were exact-match, and the one dangerous tool call was caught by guards. Keep that bar.

## f) Up to 50 Things To Get Done Next

**Public consumer path (highest impact)**

1. Fix GitHub Actions billing (or confirm public-repo free tier runs) — unblocks everything below
2. Push the 4 commits; watch the cleaned zero-auth workflows actually run green
3. Delete `DEPLOY_KEY_GHL/_SDK/_FINDING/_GOGF/_ERRFAM` secrets after that green run
4. Remove the unused deploy keys from the 4 dependency repos' settings
5. Confirm pkg.go.dev indexing completes; if docs fail to build, investigate the `GOEXPERIMENT=jsonv2` server-side build limitation
6. Verify `go install github.com/larsartmann/go-humanize-linter/cmd/go-humanize-linter@latest` end-to-end anonymously; document the `GOEXPERIMENT=jsonv2` requirement in the README install section (or provide prebuilt-binary install path)
7. Test `action.yml` in a scratch public repo (the true consumer path)
8. Test the `golangci-lint custom` + module plugin flow from a stranger's machine perspective; copy-paste-ify the README section

**CI / infra hardening**
9. Add `nix flake check` step to CI (vendorHash drift guard — this exact bug sat unnoticed on main)
10. Add anonymous-consumer CI step: fresh temp module, `GOWORK=off`, `go get` + build against `@v0.2.0`
11. Add `go install @latest` smoke test to CI
12. Consider adding `nix run .#custom-lint` to CI (currently only local)
13. Dependabot for GitHub Actions versions + Go modules; make dep-bump PRs update `vendorHash` in the same change (add a CI hint/comment near `vendorHash` saying "update when go.sum changes")
14. Confirm Codecov actually receives uploads on the public repo; add badge only after verified
15. Re-check `permissions: contents: read` posture post-cleanup; workflows are now minimal — good
16. Evaluate whether the `Graph Update: go_modules` dynamic workflow should stay enabled

**Repo hygiene**
17. Check whether `reports/coverage.out` and `result/` are git-tracked; gitignore if so (public-repo bloat)
18. Force direnv re-eval so `.direnv/flake-profile-*.rc` (stale GOPRIVATE inside) regenerates
19. GitHub social preview image
20. Decide on `.crush/` session DB — ensure it's gitignored (grep showed it contains matching strings)
21. Sweep for any other stale env-var docs in `docs/validation/*`, `docs/adr/*` (grep was clean, but double-check ADRs referencing CI auth)
22. Verify the two GitHub Releases (v0.1.0, v0.2.0) have their binary assets attached post-flip
23. Add SHA256SUMS to future releases; consider multi-arch (darwin-arm64 at minimum — `go install` covers it, release assets don't)
24. Add a next-release CHANGELOG entry: "repo made public; CI no longer requires deploy keys; GOPRIVATE no longer needed"

**Docs for strangers**
25. README quick start rewrite for public audience: `go install` one-liner, `golangci-lint custom` snippet, action.yml usage
26. README badges (CI, pkg.go.dev, license, Go Report Card) — after CI is green
27. Add `docs/adr/000X-go-public-and-deauth.md` recording this decision + the stale-REQUIRED lesson
28. Decide policy for `docs/status` + `docs/planning` historical content now that the repo is public (private project names, candid language, "GOPRIVATE required" claims) — scrub, annotate as historical, or leave
29. Document the suppression namespace story (`gohumanize` vs module path) prominently — AIs write the wrong one (already in AGENTS; surface in README troubleshooting)
30. Document `GOEXPERIMENT=jsonv2` as a known adoption friction; track go-finding upstream for when json/v2 ships default
31. Expand `example_test.go` into godoc-visible examples (pkg.go.dev front page)
32. CONTRIBUTING: add "good first issue" pointers now that external PRs are possible

**Release / distribution**
33. Cut the next release (v0.2.1) from the de-authed main so `@latest` consumers get clean metadata
34. Verify release workflow's `gh release create` path works post-flip (tag freezes the workflow file — do not tag while red, see F1)
35. Consider GoReleaser vs current hand-rolled build steps
36. Consider nixpkgs/homebrew tap packaging once stable

**Community / growth**
37. Issue templates + PR template
38. SECURITY.md (private vulnerability reporting)
39. Enable GitHub Discussions (optional)
40. Submit to golangci-lint's module-plugin docs listing (upstream PR)
41. awesome-go / static-analysis list submissions
42. Demo GIF/asciinema in README (hyperframes/website-launch territory if wanted)
43. Announcement post (blog/Twitter/Reddit r/golang)
44. Set repo homepage URL when a website exists

**Strategic / sibling-ecosystem**
45. Audit the OTHER private repos that are still private — this session proved doc-rot goes unnoticed; a public-readiness sweep of siblings would catch the next `GOPRIVATE`-style staleness
46. Same billing fix applies to all sibling repos' CI — verify their workflows aren't silently dead too
47. Consider consolidating the `GOEXPERIMENT=jsonv2` requirement: it now burdens every consumer (CI env, action.yml, docs, pkg.go.dev) — weigh pinning go-finding versions that don't need it vs waiting for Go default
48. Add a `make`-free, nix-optional path check: `GOEXPERIMENT=jsonv2 go test ./...` must be the documented minimum (already true — keep it true)
49. Re-verify AGENTS.md "Critical" claims quarterly (this session is the case study for why)
50. Celebrate: the linter is installable by every Go developer on the planet as of tonight — ship the announcement when docs catch up

## g) Questions I Cannot Figure Out Myself

1. **GitHub Actions billing:** The account shows "recent payments have failed or spending limit needs to be increased" — every CI run account-wide dies before starting. Can you fix/confirm billing, and is relying on the free public-repo Actions tier acceptable for this repo going forward? I cannot see billing state or change spending limits.
2. **Push authorization:** 4 commits (workflow de-auth, vendorHash fix, docs/tests cleanup, AGENTS.md) sit local. Do you want me to push to `main` now (activating the cleaned CI once billing is fixed), or do you push yourself?
3. **Historical docs policy:** `docs/status/*` and `docs/planning/*` (now public) contain "GOPRIVATE required" claims, candid self-review language, and **names of your other private projects**. Scrub/redact them, annotate as immutable historical snapshots, or leave as-is?

---

_Verification artifacts from this session: `nix flake check` "all checks passed", `nix run .#test-race` 4/4 ok, `nix run .#lint` 0 issues, anonymous `go get @v0.2.0` success, `GOWORK=off` full-suite green. Open failures are environmental (billing, pkg.go.dev indexing), not code._
