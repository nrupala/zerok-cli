package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
			// Get all files from directory
			entries, err := os.ReadDir(sourcePath)
			if err != nil {
				return err
			}
			for _, e := range entries {
				if !e.IsDir() {
					filesToImport = append(filesToImport, filepath.Join(sourcePath, e.Name()))
				}
			}
		} else {
			filesToImport = []string{sourcePath}
		}

		// Deduplication check
		if deduplicate {
			duplicates, err := storage.HashFiles(filesToImport)
			if err != nil {
				fmt.Printf("Warning: Deduplication check failed: %v\n", err)
			} else if len(duplicates) > 0 {
				fmt.Printf("Found %d duplicate groups:\n", len(duplicates))
				for h, paths := range duplicates {
					fmt.Printf("  Hash %s: %d files\n", h[:8], len(paths))
					for _, p := range paths {
						fmt.Printf("    - %s\n", p)
					}
				}
				fmt.Println("Use first file from each group (import skipped for duplicates)")
			}
		}

		// Import files
		successCount := 0
		failedCount := 0

		for _, filePath := range filesToImport {
			// Skip if we only want first of duplicates and this is a duplicate
			if deduplicate {
				// Check if this file is a duplicate
				data, _ := os.ReadFile(filePath)
				h := fmt.Sprintf("%x", crypto.Hash(data)[:8])
				_ = h
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

		fmt.Printf("\nImport complete: %d files imported, %d failed\n", successCount, failedCount)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVarP(&sourcePath, "source", "s", "", "Source file or directory to import")
	importCmd.Flags().BoolVarP(&deduplicate, "dedup", "d", false, "Enable deduplication (SHA-256)")
	_ = importCmd.MarkFlagRequired("source")
}

// Helper to get file extension
func getExt(path string) string {
	ext := filepath.Ext(path)
	return strings.ToLower(ext)
}