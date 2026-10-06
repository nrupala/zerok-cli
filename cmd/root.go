// Copyright (c) 2026 Nrupal Akolkar
// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	vaultPath string
	password  string
)

var rootCmd = &cobra.Command{
	Use:   "zerok-cli",
	Short: "Zerok Vault CLI - Zero-knowledge encrypted file management",
	Long: `Zerok Vault CLI - Command-line tool for managing encrypted vaults

Features:
- Mass import of files with encryption
- SHA-256 deduplication
- Integrity verification
- Vault management

Examples:
  zerok-cli init --path ./vault --password mysecret
  zerok-cli import --path ./vault --source ~/photos
  zerok-cli verify --path ./vault
  zerok-cli hash --path ./vault --dedup`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&vaultPath, "path", "p", "./vault", "Path to vault directory")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "w", "", "Vault password (prompted if not provided)")
}

func getVaultPath() string {
	if vaultPath == "" {
		cwd, _ := os.Getwd()
		return cwd + "/vault"
	}
	return vaultPath
}

func getPassword() string {
	if password != "" {
		return password
	}
	fmt.Print("Enter vault password: ")
	var pwd string
	fmt.Scanln(&pwd)
	return pwd
}

func ensureVaultDir() error {
	path := getVaultPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}