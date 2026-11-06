# zerok-cli — Zerok Vault in Go

Zerok Vault CLI is a zero-knowledge encrypted file vault for the command
line: files are encrypted on your own machine with a password-derived key
before they are stored, so the vault directory only ever holds ciphertext.

The vault uses AES-256-GCM encryption with a PBKDF2 key derivation
(600,000 iterations, SHA-256, 32-byte random salt) from `pkg/crypto`, a
content-addressed file store with SHA-256 deduplication in `pkg/storage`,
and Cobra commands in `cmd/` (`init`, `import`, `verify`, `hash`).

## Install

You need Go 1.21 or newer.

```bash
git clone https://github.com/nrupala/zerok-cli.git
cd zerok-cli
go build -o zerok ./cmd/zerok-cli
```

## Quick start

```bash
# Initialize a new vault (prompts for a password, or pass --password)
./zerok init --path ./vault --password 'your-password'

# Import files into the vault, skipping duplicates by SHA-256 content hash
./zerok import --path ./vault --source ./photos --dedup --password 'your-password'

# Verify the integrity of the encrypted files in the vault
./zerok verify --path ./vault --password 'your-password'

# Find duplicate files in a folder by SHA-256 hash
./zerok hash --source ./folder
```

## Usage

Global flags (all commands): `--path` / `-p` — vault directory
(default `./vault`); `--password` / `-w` — vault password (you are
prompted if it is not supplied).

| Command | What it does | Command flags |
|---|---|---|
| `init` | Create a vault: generate a salt, derive the key, write `salt.bin`, `verifier.bin` and `vault.json`, and create the `data/` store | `--force` / `-f` — overwrite an existing vault |
| `import` | Encrypt and import a file or a whole directory tree into the vault | `--source` / `-s` (required), `--dedup` / `-d` — skip duplicate content |
| `verify` | Check the integrity of every encrypted file in the vault | `--quiet` / `-q` — only show errors |
| `hash` | Scan a directory and report groups of files with identical SHA-256 hashes | `--source` / `-s`, `--move` / `-m` — move duplicates to a directory |

## Development

- Run the tests: `go test ./...`
- Run the linter: `go vet ./...`
- The test plan (scope, suites, gates, cadence) lives in TEST-PLAN.md.
- CI runs vet, the test suite, a build and CLI round-trip tests on every
  push and pull request to `master`.

## Security

Found a vulnerability? Please report it privately — see
[SECURITY.md](SECURITY.md) for the reporting route (GitHub Private
Vulnerability Reporting). Do not open a public issue for a security
vulnerability.

## Contributing

Contributions are welcome — the process and the test policy are in
[CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT — see [LICENSE](LICENSE). Attribution for third-party components is
in [NOTICE](NOTICE).
