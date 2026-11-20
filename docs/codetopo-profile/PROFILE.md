# Codetopo Profile — Zerok CLI (nrupala/zerok-cli)

**What this is.** The agent-facing map of this repository, generated with
the full-language Codetopo engine (local, `~/.cargo/bin/codetopo`) so an
arriving agent can see what the repo holds — topology, trust boundaries,
money surface — before it submits or changes anything. Read this first;
the graph snapshot beside it (`snapshot.json`, `stats.json`) is the
measured evidence. This repo was **SKIPPED-FOR-LANGUAGE** in the first
estate sweep (2026-11-18): it is Go-only, and Codetopo 0.9.1 had no Go
extractor. The Go extractor landed in the full-language program
(PR #157), and this is the repo's first profile.

- **Repo:** nrupala/zerok-cli
- **Commit profiled:** `f2cf94b33b200a138cbf0db254c64bd62ce2d6fb` (main)
- **Engine:** codetopo 0.9.1 + full-language program (engine built from codetopo main `d71509d1` — all sixteen program extractors merged) · index run 2026-11-19
- **Graph:** 14 files walked · 0 failed · 146 nodes · 358 edges · audit chain valid (3 entries, head `ad8ef136…edfc`)
- **Nodes by kind:** external 66 · function 24 · variable 17 · file 14 · module 10 · doc 6 · method 6 · class 2 · package 1
- **Edges by kind:** calls 144 · reads 101 · contains 79 · imports 25 · defines 6 · writes 3

## Coverage honesty

Codetopo now extracts every language in the full-language program
(Rust, TypeScript/JavaScript, Python, Kotlin, C, Java, Go, Shell, C++,
Swift, Ruby, R, Julia, C#, PHP, Scala, Haskell, Elixir, Lua), plus
Markdown as Doc nodes. This repo's code is Go (8 `.go` files) and is
extracted at full symbol level; the 6 Markdown files are indexed as
Doc nodes. `go.mod`/`go.sum`, the SVG/PNG assets, `LICENSE` and
`NOTICE` are data, not code, and are not parsed. Symbol coverage of
the codebase is complete.

## Topology

A zero-knowledge encrypted file vault for the command line (README):
files are encrypted on the user's own machine with a password-derived
key before they are stored, so the vault directory only ever holds
ciphertext.

- `cmd/` — the Cobra commands: `init`, `import`, `verify`, `hash`
  (`root.go`, `init.go`, `import.go`, `verify.go`), with the binary
  entry point at `cmd/zerok-cli/main.go`.
- `pkg/crypto` — AES-256-GCM encryption with PBKDF2 key derivation
  (600,000 iterations, SHA-256, 32-byte random salt).
- `pkg/storage` — the content-addressed file store with SHA-256
  deduplication.

## Trust boundaries

- **Zero-knowledge by construction:** encryption happens locally with
  a key derived from the user's password; the vault directory holds
  ciphertext only. There is no account, no server, and no network
  surface in this CLI.
- Passwords enter via prompt or the `--password` flag (README); no
  credential material is stored by the tool itself.

## Money surface

**None found — explicitly.** There is no payment, billing, checkout,
card, subscription or webhook surface in this repo. Greps for
`stripe|paddle|billing|checkout|webhook|api_key` over the Go sources
return nothing.

## How an agent should approach

1. Read `README.md` and `TEST-PLAN.md` first; the crypto parameters
   above are the security surface — changing them changes the vault
   format and orphans existing vaults.
2. Regenerate this profile (below) if the repo's shape changes.

## Regeneration

`docs/codetopo-profile/refresh.sh` (in this pack) re-downloads the repo
at a named ref, re-runs the index, and rewrites `snapshot.json`,
`stats.json` and the header block inputs of this profile. Regenerate at
every material change to the codebase.
