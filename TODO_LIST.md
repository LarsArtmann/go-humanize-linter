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

### T2 — Real-world validation sweep with new detection + features · High · _planned_

H001–H009 were swept against 327 Go projects
(`docs/validation/2026-07-31_real-world-sweep.md`), but the new features
(`--verify-suppressions`, `--min-confidence`, plugin confidence filtering, the
H001 size-bucket filter, import-alias-aware detection, dot-import support, and
the gogenfilter generated-file integration) have not been validated against the
corpus with these features active.

- [ ] Run the linter over the 327-project corpus with all 9 rules enabled
- [ ] Run `--verify-suppressions` on the corpus and record stale-directive rate
- [ ] Run `--min-confidence high` and compare finding counts
- [ ] Verify the gogenfilter integration does not skip hand-written files erroneously
- [ ] Save results to `docs/validation/`

---

## Distribution

### T18 — Publish to golangci-lint plugin index · Low · _blocked_

- [ ] Blocked on a tagged, `go install`-able version (v0.2.0 not tagged yet — see T1)
- [ ] Submit to the plugin index once installable
