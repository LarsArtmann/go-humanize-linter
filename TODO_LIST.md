# TODO List

> Short- and mid-term improvement tasks for go-humanize-linter.
>
> **Open work only.** When a task is finished, delete it here and record it in
> `CHANGELOG.md`. Every item below has been verified against the code as not-yet-done.
>
> Legend — Effort: `XS` ≤15 min · `S` ≤1 h · `M` ≤4 h · `L` ≥½ day.

## Summary

| #   | Task                                                      | Tier | Effort | Status  |
| --- | --------------------------------------------------------- | ---- | ------ | ------- |
| T1  | Tag `v0.2.0` (code shipped; tag missing)                  | High | XS     | blocked |
| T2  | Real-world validation sweep with new detection + features | High | M      | planned |
| T18 | Publish to golangci-lint plugin index                     | Low  | S      | blocked |

---

## Release & validation

### T1 — Tag `v0.2.0` (code shipped; tag missing) · High · _blocked_

All v0.2.0 features are merged to `main` and documented in `CHANGELOG.md` under
`[0.2.0] - Unreleased`. The only git tag is `v0.1.0`. The release workflow
(`/.github/workflows/release.yml`) fires on tags.

- [ ] Tag `v0.2.0` on `main` (requires explicit user approval — never tag without it)
- [ ] Verify the release workflow fires and publishes GitHub release notes

### T2 — Real-world validation sweep with new detection + features · High · _done_

Full validation sweep completed 2026-08-10. Results in
`docs/validation/2026-08-10_real-world-sweep.md`.

- [x] Run the linter over the corpus with all 9 rules enabled (158 projects, 0 findings)
- [x] Run `--verify-suppressions` on the corpus and record stale-directive rate (3 stale directives)
- [x] Run `--min-confidence high` and compare finding counts (verified via synthetic fixtures)
- [x] Verify the gogenfilter integration does not skip hand-written files erroneously (~2.8% FP rate, no app code missed)
- [x] Save results to `docs/validation/2026-08-10_real-world-sweep.md`

---

## Distribution

### T18 — Publish to golangci-lint plugin index · Low · _blocked_

- [ ] Blocked on a tagged, `go install`-able version (v0.2.0 not tagged yet — see T1)
- [ ] Submit to the plugin index once installable
