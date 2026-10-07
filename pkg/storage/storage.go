// Copyright (c) 2026 Nrupal Akolkar
// SPDX-License-Identifier: MIT

package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/zerok-vault/zerok-cli/pkg/crypto"
)

type FileStore struct {
	BasePath string
	Key      []byte
}

func NewFileStore(basePath string, key []byte) (*FileStore, error) {
	// Ensure base directory exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, err
	}
	return &FileStore{BasePath: basePath, Key: key}, nil
}

// FileInfo holds metadata about a stored file
type FileInfo struct {
	ID        string
	Name      string
	Size      int64
	MimeType  string
	Hash      string
	CreatedAt int64
}

// EncryptAndSave encrypts a file and saves it to the store
func (s *FileStore) EncryptAndSave(sourcePath string) (*FileInfo, error) {
	// Read source file
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read source file: %w", err)
	}

	// Generate file ID based on content hash
	hash := sha256.Sum256(data)
	fileID := hex.EncodeToString(hash[:16])

	// Check if file already exists (deduplication)
	existingPath := filepath.Join(s.BasePath, fileID+".enc")
	if _, err := os.Stat(existingPath); err == nil {
		// File already exists, return existing info
		return &FileInfo{
			ID:        fileID,
			Name:      filepath.Base(sourcePath),
			Size:      int64(len(data)),
			Hash:      hex.EncodeToString(hash[:]),
			CreatedAt: 0,
		}, nil
	}

	// Encrypt the data
	encrypted, err := crypto.EncryptFile(s.Key, data)
	if err != nil {
		return nil, fmt.Errorf("encryption failed: %w", err)
	}

	// Save encrypted data
	if err := os.WriteFile(existingPath, encrypted, 0600); err != nil {
		return nil, fmt.Errorf("failed to write encrypted file: %w", err)
	}

	return &FileInfo{
		ID:        fileID,
		Name:      filepath.Base(sourcePath),
		Size:      int64(len(data)),
		Hash:      hex.EncodeToString(hash[:]),
		CreatedAt: 0,
	}, nil
}

// DecryptAndRead decrypts a file and returns its contents
func (s *FileStore) DecryptAndRead(fileID string) ([]byte, error) {
	encryptedPath := filepath.Join(s.BasePath, fileID+".enc")
	
	encrypted, err := os.ReadFile(encryptedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read encrypted file: %w", err)
	}

	// For now, we store nonce + ciphertext together
	// Extract nonce (first 12 bytes) and ciphertext
	if len(encrypted) < crypto.NonceBytes {
		return nil, errors.New("invalid encrypted file format")
	}

	nonce := encrypted[:crypto.NonceBytes]
	ciphertext := encrypted[crypto.NonceBytes:]

	// Decrypt
	plaintext, err := crypto.DecryptFile(s.Key, nonce, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return plaintext, nil
}

// Delete removes an encrypted file from the store
func (s *FileStore) Delete(fileID string) error {
	path := filepath.Join(s.BasePath, fileID+".enc")
	return os.Remove(path)
}

// ListFiles lists all encrypted files in the store
func (s *FileStore) ListFiles() ([]string, error) {
	entries, err := os.ReadDir(s.BasePath)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".enc" {
			files = append(files, e.Name()[:len(e.Name())-4])
		}
	}
	return files, nil
}

// HashFiles generates hashes for deduplication detection
func HashFiles(paths []string) (map[string][]string, error) {
	hashes := make(map[string][]string)

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		h := sha256.Sum256(data)
		hashStr := hex.EncodeToString(h[:])

		hashes[hashStr] = append(hashes[hashStr], path)
	}

	// Return only duplicates
	duplicates := make(map[string][]string)
	for h, paths := range hashes {
		if len(paths) > 1 {
			duplicates[h] = paths
		}
	}

	return duplicates, nil
}

// ImportDirectory recursively imports a directory
func (s *FileStore) ImportDirectory(dirPath string) ([]FileInfo, error) {
	var results []FileInfo

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip hidden files
		if filepath.Base(path)[0] == '.' {
			return nil
		}

		// Import file
		fileInfo, err := s.EncryptAndSave(path)
		if err != nil {
			fmt.Printf("Warning: Failed to import %s: %v\n", path, err)
			return nil // Continue with other files
		}

		results = append(results, *fileInfo)
		fmt.Printf("Imported: %s -> %s\n", path, fileInfo.ID)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

// CopyFile copies a file (for backup/restore operations)
func CopyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	return err
}

// VerifyIntegrity checks file integrity
func (s *FileStore) VerifyIntegrity(fileID string) (bool, error) {
	_, err := s.DecryptAndRead(fileID)
	return err == nil, err
}