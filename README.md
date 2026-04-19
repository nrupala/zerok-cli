# Zerok Vault - Go CLI

Command-line interface for Zerok Vault - Zero-knowledge encrypted storage.

## Features

- **Vault Management**: Create and manage encrypted vaults
- **Bulk Import**: Import files with encryption (50GB+ folders)
- **SHA-256 Deduplication**: Detect and skip duplicate files
- **Integrity Verification**: Verify encrypted file integrity
- **Password-based Encryption**: AES-256-GCM with PBKDF2 (600K iterations)

## Installation

```bash
# Clone the repository
git clone https://github.com/nrupala/zerok-cli.git
cd zerok-cli

# Build
go build -o zerok ./cmd/zerok-cli

# Run
./zerok --help
```

## Commands

```bash
# Initialize a new vault
zerok init --path ./vault

# Import files (with optional deduplication)
zerok import --source ./photos --path ./vault --dedup

# Verify vault integrity
zerok verify --path ./vault

# Find duplicate files
zerok hash --source ./folder
```

## License

MIT License - Copyright (c) 2026 Nrupal Akolkar