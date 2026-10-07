// Copyright (c) 2026 Nrupal Akolkar
// SPDX-License-Identifier: MIT

package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zerok-vault/zerok-cli/pkg/crypto"
	"github.com/zerok-vault/zerok-cli/pkg/storage"
	"github.com/spf13/cobra"
)

var (
	sourcePath string
	deduplicate bool
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import files into the vault",
	Long:  `Encrypts and imports files from a directory into the vault`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pwd := getPassword()
		if pwd == "" {
			return nil
		}

		path := getVaultPath()
		
		// Load vault key
		salt, err := os.ReadFile(path + "/salt.bin")
		if err != nil {
			return fmt.Errorf("vault not initialized. Run 'zerok-cli init' first")
		}
		
		key := crypto.DeriveKey(pwd, salt)

		// Initialize storage
		store, err := storage.NewFileStore(path+"/data", key)
		if err != nil {
			return err
		}

		// Handle source path
		if sourcePath == "" {
			return fmt.Errorf("source path required (--source)")
		}

		// Determine if source is file or directory
		info, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		var filesToImport []string

		if info.IsDir() {
			// Get all files from directory recursively
			filepath.Walk(sourcePath, func(p string, fi os.FileInfo, err error) error {
				if err != nil { return err }
				if !fi.IsDir() && !strings.HasPrefix(fi.Name(), ".") {
					filesToImport = append(filesToImport, p)
				}
				return nil
			})
		} else {
			filesToImport = []string{sourcePath}
		}

		// Track seen hashes for deduplication
		seenHashes := make(map[string]string) // hash -> first file path

		// Import files
		successCount := 0
		failedCount := 0
		skippedCount := 0

		for _, filePath := range filesToImport {
			// Check for duplicates
			if deduplicate {
				data, err := os.ReadFile(filePath)
				if err != nil {
					fmt.Printf("Warning: Could not read %s: %v\n", filePath, err)
					continue
				}
				
				hash := sha256.Sum256(data)
				hashStr := hex.EncodeToString(hash[:])
				
				if existingPath, exists := seenHashes[hashStr]; exists {
					fmt.Printf("Skipped (duplicate): %s (same as %s)\n", filepath.Base(filePath), filepath.Base(existingPath))
					skippedCount++
					continue
				}
				seenHashes[hashStr] = filePath
			}

			fileInfo, err := store.EncryptAndSave(filePath)
			if err != nil {
				fmt.Printf("Failed to import %s: %v\n", filePath, err)
				failedCount++
				continue
			}

			successCount++
			sizeMB := float64(fileInfo.Size) / (1024 * 1024)
			fmt.Printf("Imported: %s (%.2f MB)\n", fileInfo.Name, sizeMB)
		}

		fmt.Printf("\nImport complete: %d imported, %d failed, %d skipped\n", successCount, failedCount, skippedCount)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVarP(&sourcePath, "source", "s", "", "Source file or directory to import")
	importCmd.Flags().BoolVarP(&deduplicate, "dedup", "d", false, "Enable deduplication (SHA-256)")
	_ = importCmd.MarkFlagRequired("source")
}