// Copyright (c) 2026 Nrupal Akolkar
// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"os"

	"github.com/zerok-vault/zerok-cli/pkg/crypto"
	"github.com/zerok-vault/zerok-cli/pkg/storage"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new vault",
	Long:  `Creates a new encrypted vault with the specified password`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pwd := getPassword()
		if pwd == "" {
			return nil // Cancelled or empty
		}

		// Generate salt and derive key
		salt, err := crypto.GenerateSalt()
		if err != nil {
			return err
		}

		key := crypto.DeriveKey(pwd, salt)

		// Create verifier for password validation
		verifier, err := crypto.CreateVerifier(key)
		if err != nil {
			return err
		}

		// Create vault directory structure
		path := getVaultPath()
		if err := os.MkdirAll(path+"/data", 0755); err != nil {
			return err
		}

		// Write salt (for future key derivation)
		os.WriteFile(path+"/salt.bin", salt, 0600)
		os.WriteFile(path+"/verifier.bin", verifier, 0600)

		// Write metadata JSON
		vaultJSON := fmt.Sprintf(`{"created":%d,"version":"1.0","encryption":"AES-256-GCM"}`, os.Stderr.Fd())
		os.WriteFile(path+"/vault.json", []byte(vaultJSON), 0644)

		// Initialize file store
		store, err := storage.NewFileStore(path+"/data", key)
		if err != nil {
			return err
		}

		// Test encryption with a small sample
		testData := []byte("Zerok Vault CLI initialized - encryption test")
		encrypted, err := store.EncryptAndSave(path + "/.init_test")
		if err != nil {
			return fmt.Errorf("encryption test failed: %w", err)
		}
		_ = encrypted

		// Clean up test
		os.Remove(path + "/.init_test")

		fmt.Printf("Vault initialized at: %s\n", path)
		fmt.Println("Encryption: AES-256-GCM with PBKDF2 (600K iterations)")
		fmt.Println("Vault is ready for use!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolP("force", "f", false, "Overwrite existing vault")
}