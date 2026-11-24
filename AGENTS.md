# Agent Instructions

## Repository purpose

- Go library that turns media titles into slugs for matching. Zaparoo Core indexes media and resolves token titles with it, and other Zaparoo clients port or call it to match titles the same way.
- Slug output is a compatibility contract: Core stores slugs in its media database and resolves written tokens against them. Treat any change to `Slugify` output as breaking and call it out.
- `filetitle` and `browseprefix` output feeds `Slugify` during indexing, so the same contract applies to them.

## Commands

- Tests with race detection: `task test`
- Non-mutating lint: `task lint`
- Apply lint and formatting fixes deliberately: `task lint-fix`
- Fuzz targets: `task fuzz` (`FUZZ_TIME` overrides the 30s default)
- Benchmarks: `task bench`

## Rules

- Use the Go version declared in `go.mod`.
- Add or update tests for every behavior change, with a regression test that fails before a bug fix. Call `t.Parallel()` in test functions.
- Do not add dependencies without discussion.
- Library code does not log or print.
- Comments explain why, not what.
