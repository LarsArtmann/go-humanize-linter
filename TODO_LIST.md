# TODO List

> Short- and mid-term improvement tasks for go-humanize-linter.
>
> **Open work only.** When a task is finished, delete it here and record it in
> `CHANGELOG.md`. Every item below has been verified against the code as not-yet-done.
>
> Legend — Effort: `XS` ≤15 min · `S` ≤1 h · `M` ≤4 h · `L` ≥½ day.

## Summary

| #   | Task                                                        | Tier   | Effort | Status  |
| --- | ----------------------------------------------------------- | ------ | ------ | ------- |
| T1  | Tag `v0.2.0` (code shipped; tag missing)                    | High   | XS     | ready   |
| T18 | Publish to golangci-lint plugin index                       | Low    | S      | blocked |
| T20 | Release gogenfilter (workspace sibling 30 commits ahead)    | Medium | S      | planned |

---

## Release & validation

### T1 — Tag `v0.2.0` (code shipped; tag missing) · High · _ready_

All v0.2.0 features are merged to `main` and documented in `CHANGELOG.md` under
`[0.2.0] - Unreleased`. The only git tag is `v0.1.0`. The release workflow
(`/.github/workflows/release.yml`) fires on tags.

- [ ] Tag `v0.2.0` on `main` (requires explicit user approval — never tag without it)
- [ ] Verify the release workflow fires and publishes GitHub release notes

Note: CI (`ci.yml`) currently fails at checkout — the `DEPLOY_KEY` secret is
missing or invalid (`git@github.com: Permission denied (publickey)`). Fix the
secret so the release run can go green.

---

## Distribution

### T18 — Publish to golangci-lint plugin index · Low · _blocked_

- [ ] Blocked on a tagged, `go install`-able version (v0.2.0 not tagged yet — see T1)
- [ ] Submit to the plugin index once installable

---

## Upstream

### T20 — Release gogenfilter with the new sqlc detection · Medium · _planned_

The workspace sibling `../gogenfilter` is 30 commits ahead of the pinned v3.4.0,
with significant sqlc detection changes. The `go.work` `use` entry serves the
sibling locally, so local dev runs newer detection behavior than CI (pinned
v3.4.0). Both currently agree on bare sqlc filename behavior, but the drift is
a latent split-brain.

- [ ] Tag + push a gogenfilter release (requires explicit user approval)
- [ ] `go get github.com/LarsArtmann/gogenfilter/v3@<new>` + `go mod tidy`
- [ ] Re-run the testdata sweep to confirm no detection behavior changed
