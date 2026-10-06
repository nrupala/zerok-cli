# Contributing to zerok-cli

## Workflow (PR-flow discipline)

- Work happens on feature branches — **never push directly to `master`**.
- Open pull requests as **draft** first. Mark ready only when:
  - `go build ./...` and `go vet ./...` pass,
  - the CHANGELOG has an entry under `## [Unreleased]`,
  - the version was bumped (patch = fix, minor = feature, major = breaking) —
    this repo has no version file yet; use release tags,
  - the PR description states what was verified vs what was not.
- The owner merges. Merge commits reference the PR number.
- Releases are tagged `vX.Y.Z` after merge.

## Build and test

```bash
go build -o zerok ./cmd/zerok-cli   # build the CLI (from README)
go vet ./...                        # static checks
```

There is currently no automated test suite; `go test ./...` is the convention
to adopt when tests are added.

## License

MIT — keep the `// SPDX-License-Identifier: MIT` header on every Go file.
