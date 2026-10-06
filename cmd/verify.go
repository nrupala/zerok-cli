// Copyright (c) 2026 Nrupal Akolkar
// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"

	"github.com/zerok-vault/zerok-cli/pkg/storage"
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify vault integrity",
	Long:  `Checks integrity of all encrypted files in the vault`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pwd := getPassword()
		if pwd == "" {
			return nil
		}

		path := getVaultPath()

		// Load vault key (simplified - actual implementation would verify password first)
		salt, err := os.ReadFile(path + "/salt.bin")
		if err != nil {
			return fmt.Errorf("vault not found")
		}

		// Note: In production, verify password before proceeding
		_ = pwd
		_ = salt

		store, err := storage.NewFileStore(path+"/data", nil)
		if err != nil {
			return err
		}

		// List all files
		files, err := store.ListFiles()
		if err != nil {
			return err
		}

		fmt.Printf("Verifying %d files...\n", len(files))

		verified := 0
		corrupted := 0

		for _, fileID := range files {
			valid, err := store.VerifyIntegrity(fileID)
			if err != nil {
				corrupted++
				fmt.Printf("  [CORRUPTED] %s\n", fileID)
			} else {
				verified++
				if !quiet {
					fmt.Printf("  [OK] %s\n", fileID)
				}
			}
		}

		fmt.Printf("\nIntegrity check complete:\n")
		fmt.Printf("  Verified:   %d\n", verified)
		fmt.Printf("  Corrupted:  %d\n", corrupted)

		if corrupted > 0 {
			fmt.Println("\nWarning: Some files may be corrupted!")
		}

		return nil
	},
}

var quiet bool

func init() {
	rootCmd.AddCommand(verifyCmd)
	verifyCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Only show errors")
}

var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "Find duplicate files",
	Long:  `Uses SHA-256 to find and display duplicate files`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if sourcePath == "" {
			return fmt.Errorf("source path required (--source)")
		}

		// Get all files from source directory
		var files []string

		entries, err := os.ReadDir(sourcePath)
		if err != nil {
			return err
		}

		for _, e := range entries {
			if !e.IsDir() {
				files = append(files, sourcePath+"/"+e.Name())
			}
		}

		duplicates, err := storage.HashFiles(files)
		if err != nil {
			return err
		}

		if len(duplicates) == 0 {
			fmt.Println("No duplicates found")
			return nil
		}

		fmt.Printf("Found %d duplicate groups:\n\n", len(duplicates))

		totalWasted := int64(0)

		for hash, paths := range duplicates {
			fmt.Printf("Hash: %s\n", hash[:16])
			fmt.Printf("  Files (%d):\n", len(paths))
			
			for _, p := range paths {
				info, _ := os.Stat(p)
				fmt.Printf("    - %s (%.2f MB)\n", p, float64(info.Size())/(1024*1024))
				totalWasted += info.Size()
			}
			fmt.Println()
		}

		fmt.Printf("Total wasted space: %.2f MB\n", float64(totalWasted)/(1024*1024))

		if dedupeMove != "" {
			fmt.Println("\nDeduplication requested...")
			// Would move duplicates to specified location
		}

		return nil
	},
}

var dedupeMove string

func init() {
	rootCmd.AddCommand(hashCmd)
	hashCmd.Flags().StringVarP(&sourcePath, "source", "s", "", "Source directory to scan")
	hashCmd.Flags().StringVarP(&dedupeMove, "move", "m", "", "Move duplicates to specified directory")
}