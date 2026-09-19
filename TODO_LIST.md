# TODO List

> Short- and mid-term improvement tasks for go-humanize-linter.
>
> **Open work only.** When a task is finished, delete it here and record it in
> `CHANGELOG.md`. Every item below has been verified against the code as not-yet-done.
>
> Legend — Effort: `XS` ≤15 min · `S` ≤1 h · `M` ≤4 h · `L` ≥½ day.

## Summary

| #   | Task                                                          | Tier   | Effort | Status                                              |
| --- | ------------------------------------------------------------- | ------ | ------ | --------------------------------------------------- |
| T29 | H011 rule: manual SI-string parsing (ParseSI)                 | Medium | M      | deferred — zero corpus demand (2026-09-18 research) |

Declined 2026-09-19 (user decision): T18 plugin-index PR and T20
gogenfilter release — no foreign-repo actions for now. Consumer repos
(KeyCountdown, Kernovia, CreditReformBilanzampel) stay unfixed as the
demonstration corpus; main stays unprotected; Go directive stays 1.26.x
(revert-on-sight for background bumps).

Done 2026-09-18 (recorded in CHANGELOG): T21 push+CI, T22 H010 sweep,
T23 v0.3.0 release, T24 H010 fixture pack, T25 CI-miss investigation
(main was red since 09-13, unnoticed), T26 dependabot PR (merged #1),
T27 exit-code examples, T28 suggestion enrichments, T30 H012 rule,
T31 DEPLOY_KEY_* secrets deleted.

Done 2026-09-19: H012 conjunction-position regression guards (negative
fixtures + unit discrimination test), go-directive incident root-caused
and documented (AGENTS.md Gotchas), lint parity confirmed (local nix
0 issues, CI green), pkg.go.dev v0.3.0 verified, T32 read-only deploy
keys removed from all 4 dep repos (zero keys remain).

Done 2026-09-19: v0.4.0 released (H012 rule, H010 fixture pack, H012
guards, H002/H009 suggestion enrichments, CI/security hygiene). A second
background toolchain bump (go 1.27 + go-finding v1.12.0 sweep) was
reverted per the revert-on-sight decision before tagging; unbreak
verified with GOWORK=off GOTOOLCHAIN=go1.26.8 build/vet/test, then full
nix battery (build/vet/test-race/lint). Release verified: proxy indexed,
clean-dir go get + go run works, binary reports v0.4.0, curated release
notes written (auto-generated notes only listed the dependabot PR).
