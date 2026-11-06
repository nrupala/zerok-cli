# Test plan — zerok-cli

## Scope

Testing covers the cryptography package (`pkg/crypto`: PBKDF2 key
derivation, AES-256-GCM encrypt/decrypt, password verifier) with unit
tests, `go vet` over the whole module, a full build of the CLI binary,
and CLI-level tests in CI: a real round-trip that initialises a vault
in a temporary directory and imports a file into it, and a `hash`
smoke test over a duplicate-file fixture.

## Known defect — read this first (2026-11-06)

**The storage layer's decrypt path is broken, and `verify` therefore
reports every stored file as corrupted.** Proven end-to-end while
writing this plan: after a successful `init` + `import`, `verify`
reported both stored files `[CORRUPTED]`. Root cause, read in the code:
`crypto.EncryptFile` discards the GCM nonce returned by
`crypto.Encrypt`, while `storage.DecryptAndRead` treats the first 12
stored bytes as the nonce — the nonce is never stored, so decryption
can never succeed. Fixing this changes the on-disk vault format of a
security-bearing product; it is a maintainer decision (STOPPED for
this leg, recorded here and in the filing pack), not a readiness fix.
The encrypt/write side (`init`, `import`) works and is what CI tests.

## Contained build blockers fixed by this leg

The module did not compile at all when this leg started. Four small,
mechanical blockers — missing imports and the author's own half-written
lines, no design or crypto change — were fixed so CI evidence can exist:

- `cmd/import.go`: used `strings.HasPrefix` without importing `strings`.
- `cmd/verify.go`: used `os.ReadFile`/`os.ReadDir`/`os.Stat` without
  importing `os`; and a `valid` variable was declared but never used
  (`VerifyIntegrity` already returns `err == nil` as its boolean, so
  the assignment is now blank — semantics unchanged).
- `cmd/init.go`: the initialisation self-test prepared `testData` but
  never wrote it, then tried to encrypt the non-existent file (and the
  unused variable itself was a compile error). The file is now written
  before `EncryptAndSave`, completing the author's evident intent;
  `init` now succeeds end-to-end.
- `go.mod`/`go.sum`: the repo had no `go.sum` (and no indirect-require
  block); both were generated with `go mod tidy` against the declared
  dependency versions (cobra v1.8.0, x/crypto v0.21.0).

## Suites

| Suite | Covers | Run it | Runs in CI |
|---|---|---|---|
| Unit tests (`pkg/crypto/crypto_test.go`) | Encrypt→decrypt round-trip returns the plaintext; decryption with a wrong key fails; password verifier accepts the correct password and rejects a wrong one; generated salts are unique and correctly sized | `go test ./...` | Yes — CI workflow |
| CLI round-trip (CI steps) | The built binary initialises a vault (`salt.bin`, `verifier.bin` created, "Vault is ready") and imports a fixture file (an `.enc` appears in `vault/data/`) | `go build -o zerok ./cmd/zerok-cli`, then the CI step's commands | Yes — CI workflow |
| CLI smoke (CI step) | The `hash` command finds the one duplicate group in a three-file fixture (two identical, one different) | `./zerok hash --source <dir>` | Yes — CI workflow |

## CI gates

- CI workflow (`.github/workflows/ci.yml`) on every push and pull
  request to `master`: `go vet ./...`, `go test ./...`, build, CLI
  round-trip and hash smoke — all must pass.
- CodeQL workflow (`.github/workflows/codeql.yml`, Go, autobuild —
  `build-mode: none` is not supported for Go) on push, pull request
  and a weekly schedule; findings are dispositioned per the house
  alert SLA (fix, or a dismissal with the reason recorded).

## Cadence

- Every pull request: vet + the full test suite + CLI tests, and CodeQL.
- Weekly: scheduled CodeQL run.
- Every release candidate: the release-readiness run (checker +
  evidence reads), per the house Release-Readiness Standard.

## Coverage bar

No numeric floor; new functionality lands with tests per CONTRIBUTING.
