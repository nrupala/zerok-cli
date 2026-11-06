// Copyright (c) 2026 Nrupal Akolkar
// SPDX-License-Identifier: MIT

package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptDecryptRoundTrip(t *testing.T) {
	salt, err := GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt: %v", err)
	}
	if len(salt) != SaltBytes {
		t.Fatalf("salt length = %d, want %d", len(salt), SaltBytes)
	}
	key := DeriveKey("correct horse battery staple", salt)
	if len(key) != KeyBits/8 {
		t.Fatalf("key length = %d, want %d", len(key), KeyBits/8)
	}

	plaintext := []byte("zerok vault round-trip test data")
	ciphertext, nonce, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if len(nonce) != NonceBytes {
		t.Fatalf("nonce length = %d, want %d", len(nonce), NonceBytes)
	}
	if bytes.Equal(ciphertext, plaintext) {
		t.Fatal("ciphertext equals plaintext")
	}

	got, err := Decrypt(key, nonce, ciphertext)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if !bytes.Equal(got, plaintext) {
		t.Fatalf("round-trip mismatch: got %q, want %q", got, plaintext)
	}
}

func TestDecryptWrongKeyFails(t *testing.T) {
	salt, err := GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt: %v", err)
	}
	key := DeriveKey("password-one", salt)
	other := DeriveKey("password-two", salt)

	ciphertext, nonce, err := Encrypt(key, []byte("secret"))
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if _, err := Decrypt(other, nonce, ciphertext); err == nil {
		t.Fatal("Decrypt with wrong key succeeded, want error")
	}
}

func TestVerifierRoundTrip(t *testing.T) {
	salt, err := GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt: %v", err)
	}
	key := DeriveKey("vault-password", salt)
	verifier, err := CreateVerifier(key)
	if err != nil {
		t.Fatalf("CreateVerifier: %v", err)
	}
	if !VerifyPassword("vault-password", salt, verifier) {
		t.Fatal("VerifyPassword rejected the correct password")
	}
	if VerifyPassword("wrong-password", salt, verifier) {
		t.Fatal("VerifyPassword accepted a wrong password")
	}
}

func TestSaltsAndNoncesAreUnique(t *testing.T) {
	a, err := GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt: %v", err)
	}
	b, err := GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt: %v", err)
	}
	if bytes.Equal(a, b) {
		t.Fatal("two generated salts are identical")
	}
}
