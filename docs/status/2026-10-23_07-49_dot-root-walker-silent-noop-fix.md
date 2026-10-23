# Status Report: Dot-Root Walker Fix (`go-humanize-linter .` scanned nothing)

- **Date:** 2026-10-23 07:49
- **Session scope:** diagnose `go-humanize-linter .` → 0 findings on darknet-monitor; fix; test; changelog. This report covers THIS session only.
- **Verdict:** Root cause found and fixed in-repo. Fix is NOT yet on PATH (nix snapshot) and NOT released.

## The Bug

`skipDirName()` (walker.go) skips any directory whose basename starts with `.` — including the walk root itself when invoked as `.` or `./`. `filepath.SkipDir` on the first callback aborted the walk before a single file was visited: `go-humanize-linter .` = scan of **zero files** = 0 findings, exit 0, no warning. Every other root spelling (`internal`, `./internal/cli`, absolute paths) worked, which masked the bug indefinitely. `--list-files .` confirmed: 0 files.

Trigger case: `~/projects/darknet-monitor/internal/cli/humanize.go` — a textbook H001 (`KiB/MiB/GiB` const chain + `%.1f GiB` Sprintfs + switch), invisible under dot-root, correctly reported under absolute path.

**Fix:** walker.go — walk root entry (`path == dir`) is exempt from all skip rules; the caller explicitly named it, so it gets scanned. Subdirectory skip rules unchanged.

---

## a) FULLY DONE

1. **Diagnosis** — isolated to dot-root via A/B testing (`.` → 0 findings, `internal` → 1, absolute → 1, `./` → 0); confirmed mechanism with `--list-files .` → 0 files.
2. **Verified PATH binary currency** — store path carries commit `52439af` = HEAD, ruling out the known stale-binary gotcha before blaming it.
3. **Walker fix** (walker.go:92) — root-entry skip exemption with explanatory comment.
4. **Regression test** `TestWalkGoDir_DotRootScansEverything` — covers `.` and `./`, asserts hidden subdir (`.hidden/`) still skipped under a dot root, uses `t.Chdir` (correctly NOT `t.Parallel` — Chdir forbids it).
5. **Full suite green** — `nix run .#test`, all 4 packages ok, including plugin tests.
6. **Lint green** — `nix run .#lint`, clean.
7. **CHANGELOG entry** — Unreleased/Fixed, full root-cause narrative, regression guard named.
8. **End-to-end verification** — fresh build (`go build -o /tmp/ghl`) on darknet-monitor with `.`: reports `internal/cli/humanize.go:6 [H001]`, exit 1. Correct suggestion (`humanize.IBytes` for KiB/MiB/GiB).

## b) PARTIALLY DONE

1. **Deployment** — fix lives in the repo only. The PATH binary (`/nix/store/...-52439af...`) predates it; every dot-root invocation stays a silent no-op until the next system switch. Flagged to user, not actionable by me.
2. **Test coverage** — unit-level done. Missing: CLI-level e2e dot-root test (run the actual binary with `.` via `t.Chdir`), `..` root case, trailing-slash variants beyond `./`.
3. **Blast-radius audit** — NOT performed (see c). The bug plausibly invalidated past "clean" results anywhere `.` was used.

## c) NOT STARTED

1. **action.yml audit** — does the GitHub Action invoke the binary with `.`? If yes, every consumer CI run ever was a silent no-op. Highest-priority open item; one file read.
2. **CI self-scan step audit** — what path form does the repo's own self-scan use? (The self-detection gotcha implies it finds things, so probably not dot — unverified.)
3. **README invocation audit** — does any doc recommend `go-humanize-linter .`?
4. **CLI "scanned 0 files" guard** — a stderr warning when a scan finds no Go files would have exposed this bug class forever.
5. **Release** — fix is Unreleased; see question 1 (release is entangled with T33).
6. **AGENTS.md gotcha update** — PATH-binary gotcha could note the dot-root divergence reason until next switch.
7. **darknet-monitor remediation** — the now-visible H001 (`humanizeBytes` → `humanize.IBytes`) is unfixed in that repo; also worth re-scanning darknet-monitor fully now that dot-root works.
8. **`nix run .#test-race` / `#vet` / `#coverage` / `nix flake check`** — not run this session (only `#test` + `#lint`).

## d) TOTALLY FUCKED UP

1. **Edit collateral** — my first walker_test.go edit used the NEXT test's doc-comment first line as its anchor and consumed it, corrupting the file. Caught on sight, repaired immediately. Lesson: anchor with unique trailing context, not a following comment line.
2. **Bogus fixture** — first test run failed because I wrote `sub/skipped.go` and expected it skipped (`sub` matches no skip rule — sloppy fixture authoring). Wasted a full test cycle. Removed.
3. **Diagnosis order** — I deep-read `rule_bytes.go`, `pattern_bytes.go`, `hasConst1024`, `specHas1024Value` suspecting const-chain resolution failure, when the correct first move was the trivial A/B (`.` vs absolute) that isolated the bug in one command. Code-reading before reproduction = wasted analysis. The detection logic was never broken.

## e) WHAT WE SHOULD IMPROVE

1. **Reproduce before reading** — cheap input-variation experiments (path spellings!) beat source deep-dives. Do the A/B first.
2. **Silent-empty-scan guard** — a linter that reports "0 findings" on 0 files scanned is lying by omission. Ship the warning.
3. **Scan observability** — default text output should state files scanned; `--list-files` exists but nobody runs it reflexively.
4. **CLI e2e tests for invocation forms** — the entire CLI test suite passed while the flagship invocation was a no-op. Test the user-facing entry shapes (`.`), not just library calls with tempdir paths.
5. **Blast-radius reflex** — when a silent-failure bug is root-caused, immediately audit where the broken path form is used (action.yml, CI, docs) in the same session. I stopped at the fix + changelog.

## f) NEXT — up to 50 tasks

1. Read action.yml; verify invocation path form; if `.` — fix + note in CHANGELOG that Action results were affected
2. Check CI self-scan step path form (workflows)
3. Grep README/docs for `go-humanize-linter .` recommendations; correct if misleading
4. Add CLI e2e dot-root test (exec binary with `.` in temp cwd)
5. Add `..` root case to walker tests
6. Add root-exemption note to WalkGoDir doc comment (function-level, not just inline)
7. Extend skip-rules test: `vendor`, `testdata`, `node_modules`, `__debug` basenames (currently only hidden/.bak/archived/forks covered)
8. Add invariant test: `WalkGoDir(dir)` ≡ `WalkGoDir(abs(dir))` file-set equivalence
9. Add "scanned N files" to default text output (observability)
10. Add stderr warning + distinct behavior when 0 Go files found (decide: warning only vs new exit code)
11. Decide symlink policy for walker (WalkDir doesn't follow symlinks — document or handle)
12. Run `nix run .#test-race`
13. Run `nix run .#vet`
14. Run `nix run .#coverage` — walker_test coverage delta
15. Run `nix flake check`
16. Update AGENTS.md PATH-binary gotcha: dot-root bug is a second divergence reason until next system switch
17. Resolve T33 (go 1.27.1 revert vs adopt) — blocks any release
18. Decide release vehicle for the dot-root fix (see question 1): hotfix branch from v0.4.0 vs Unreleased
19. If releasing: curate notes via `gh release edit <tag> --notes-file` (repo commits direct to main; auto-notes miss everything)
20. Post-switch verification task: confirm PATH binary scans `.` correctly after next nixos switch
21. Re-scan darknet-monitor fully with fresh binary; fix remaining findings
22. Fix darknet-monitor `humanizeBytes` → `humanize.IBytes(uint64(b))` (their repo)
23. Check darknet-monitor for other humanize hand-rolls now visible
24. Re-verify 2026-10-08 fleet sweep methodology used absolute paths (so results stand)
25. Search public GitHub for Action/binary consumers whose "clean" results may be false (gh search)
26. Consider `--strict-files` flag: error when file count is 0 (opt-in for CI)
27. Add CHANGELOG "Action users affected" note if action.yml audit confirms dot usage
28. Update FEATURES.md if it documents walker behavior
29. Add TODO_LIST entry for release-blocking chain (dot-fix → T33 → release)
30. Review whether `IsGeneratedFile` behaves identically for relative vs absolute paths (worked in tests; confirm intent)
31. Consider golden-file test for text output including the files-scanned line (pairs with task 9)
32. Plugin path sanity note: unaffected (uses `pass.Files`, never walks) — one-line AGENTS.md confirmation
33. Sweep testdata/ for any fixture invoking WalkGoDir with hidden-ish roots (gap coverage)
34. Consider testing root spelled with redundant separators (`././.`) — WalkDir semantics edge
35. Document walker root-exemption semantics in AGENTS.md testdata section (one line)
36. If hotfix branch chosen: cherry-pick walker fix + test + changelog onto v0.4.0 tag; verify CI pins pass on 1.26.7
37. Monitor main CI: currently red since 2026-09-19 (T33) — this session's commits inherit that redness
38. Confirm auto-commit daemon swept this session's changes (expected behavior per AGENTS.md)
39. After release: update AGENTS.md stale-binary triage with new store-path commit example
40. Consider upstream test in go-linter-sdk? No — bug was repo-local; skip (rejected, kept for honesty)
41. Ask darknet-monitor owner (Lars) whether IBytes or Bytes semantics wanted before fixing (KiB = IEC = IBytes — likely no question needed)
42. Add `walker_test.go` case: root is a literal `vendor` dir passed explicitly → now scans (root exemption) — assert intended semantics
43. Check `.golangci.custom.yml` example for path-form issues (plugin path immune, verify claim)
44. Consider a `--fail-empty` alias bikeshed only if task 10 lands warning-first
45. Re-read this report after switch+release; archive under docs/status/archived/ when resolved
46. Keep task list honest: 40-45 include rejected/trivial items — prune in next docs-health pass

## g) QUESTIONS (cannot figure out myself)

1. **Release vehicle:** hotfix v0.4.1 branched from the v0.4.0 tag (Go 1.26.7 base, cherry-pick only this fix — consumers unaffected by the go bump), or wait for T33 resolution and ship Unreleased as v0.5.0 (which currently carries the go 1.27.1 requirement)? Every day without a release, Action/CLI consumers scanning with `.` get silent no-ops.
2. **Output change appetite:** OK to add a "scanned N files" line (and a 0-files warning to stderr) to default text output? Scripts that parse stdout could notice; json/sarif unaffected.
3. **Private consumers:** do you know of any private repos/CI using the binary or Action whose historical "0 findings" I should re-audit once the fix ships? (I can search public GitHub, but private usage is invisible to me.)

---

**Session files touched:** `walker.go` (fix), `walker_test.go` (regression test + repaired collateral), `CHANGELOG.md` (Unreleased/Fixed entry). Tests + lint green. Now waiting for instructions.
