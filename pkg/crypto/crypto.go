package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	PBKDF2Iterations = 600000
	SaltBytes        = 32
	NonceBytes       = 12
	KeyBits          = 256
)

// Key Derivation
func DeriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, PBKDF2Iterations, KeyBits/8, sha256.New)
}

// GenerateSalt creates a random salt for key derivation
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltBytes)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	return salt, nil
}

// GenerateNonce creates a random nonce for encryption
func GenerateNonce() ([]byte, error) {
	nonce := make([]byte, NonceBytes)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return nonce, nil
}

// Encrypt encrypts data using AES-256-GCM
func Encrypt(key []byte, plaintext []byte) (ciphertext []byte, nonce []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce, err = GenerateNonce()
	if err != nil {
		return nil, nil, err
	}

	ciphertext = gcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

// Decrypt decrypts data using AES-256-GCM
func Decrypt(key []byte, nonce []byte, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New("decryption failed - wrong key or corrupted data")
	}

	return plaintext, nil
}

// EncryptFile encrypts a file and returns nonce + ciphertext
func EncryptFile(key []byte, data []byte) (encrypted []byte, err error) {
	ciphertext, _, err := Encrypt(key, data)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// DecryptFile decrypts a file given nonce + ciphertext
func DecryptFile(key []byte, nonce []byte, data []byte) ([]byte, error) {
	return Decrypt(key, nonce, data)
}

// Hash generates SHA-256 hash of data
func Hash(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

// HashString generates SHA-256 hash of string
func HashString(s string) string {
	h := sha256.Sum256([]byte(s))
	return base64.StdEncoding.EncodeToString(h[:])
}

// VerifyPassword verifies a password against a stored verifier
func VerifyPassword(password string, salt []byte, verifier []byte) bool {
	key := DeriveKey(password, salt)
	decrypted, err := Decrypt(key, verifier[:NonceBytes], verifier[NonceBytes:])
	if err != nil {
		return false
	}
	return string(decrypted) == "ZEROK_VAULT_VERIFIED"
}

// CreateVerifier creates an encrypted verifier for password validation
func CreateVerifier(key []byte) ([]byte, error) {
	plaintext := []byte("ZEROK_VAULT_VERIFIED")
	ciphertext, nonce, err := Encrypt(key, plaintext)
	if err != nil {
		return nil, err
	}

	// Combine nonce + ciphertext for storage
	result := make([]byte, NonceBytes+len(ciphertext))
	copy(result[:NonceBytes], nonce)
	copy(result[NonceBytes:], ciphertext)
	return result, nil
}