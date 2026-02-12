package utils

import (
	"testing"
)

func TestDESEncrypt_Success(t *testing.T) {
	key := "12345678"
	plaintext := "hello world"

	ciphertext, err := DESEncrypt(plaintext, key)
	if err != nil {
		t.Fatalf("DESEncrypt failed: %v", err)
	}

	if ciphertext == "" {
		t.Error("Expected non-empty ciphertext")
	}
}

func TestDESDecrypt_Success(t *testing.T) {
	key := "12345678"
	plaintext := "hello world"

	// First encrypt
	ciphertext, err := DESEncrypt(plaintext, key)
	if err != nil {
		t.Fatalf("DESEncrypt failed: %v", err)
	}

	// Then decrypt
	decrypted, err := DESDecrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("DESDecrypt failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Expected %s, got %s", plaintext, decrypted)
	}
}

func TestDESEncryptDecrypt_RoundTrip(t *testing.T) {
	testCases := []struct {
		name      string
		plaintext string
		key       string
	}{
		{"simple text", "test", "abcd1234"},
		{"long text", "this is a longer text for testing", "key12345"},
		{"empty text", "", "testkey1"},
		{"special chars", "!@#$%^&*()", "pass1234"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encrypted, err := DESEncrypt(tc.plaintext, tc.key)
			if err != nil {
				t.Fatalf("Encryption failed: %v", err)
			}

			decrypted, err := DESDecrypt(encrypted, tc.key)
			if err != nil {
				t.Fatalf("Decryption failed: %v", err)
			}

			if decrypted != tc.plaintext {
				t.Errorf("Round trip failed: expected %q, got %q", tc.plaintext, decrypted)
			}
		})
	}
}

func TestDESEncrypt_InvalidKeyLength(t *testing.T) {
	testCases := []struct {
		name string
		key  string
	}{
		{"too short", "short"},
		{"too long", "toolongkey123"},
		{"empty", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := DESEncrypt("test", tc.key)
			if err == nil {
				t.Error("Expected error for invalid key length, got nil")
			}
		})
	}
}

func TestDESDecrypt_InvalidKeyLength(t *testing.T) {
	_, err := DESDecrypt("someciphertext", "short")
	if err == nil {
		t.Error("Expected error for invalid key length, got nil")
	}
}

func TestDESDecrypt_InvalidBase64(t *testing.T) {
	key := "12345678"
	_, err := DESDecrypt("not-valid-base64!", key)
	if err == nil {
		t.Error("Expected error for invalid base64, got nil")
	}
}

func TestPkcs5Padding(t *testing.T) {
	testCases := []struct {
		name      string
		data      []byte
		blockSize int
		expected  int
	}{
		{"exact block", []byte("12345678"), 8, 16},
		{"partial block", []byte("123"), 8, 8},
		{"empty", []byte(""), 8, 8},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			padded := pkcs5Padding(tc.data, tc.blockSize)
			if len(padded) != tc.expected {
				t.Errorf("Expected length %d, got %d", tc.expected, len(padded))
			}
		})
	}
}

func TestPkcs5UnPadding(t *testing.T) {
	testCases := []struct {
		name     string
		data     []byte
		expected []byte
	}{
		{"valid padding", []byte{1, 2, 3, 4, 5, 5, 5, 5, 5}, []byte{1, 2, 3, 4}},
		{"single byte padding", []byte{1, 2, 3, 1}, []byte{1, 2, 3}},
		{"empty", []byte{}, []byte{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := pkcs5UnPadding(tc.data)
			if len(result) != len(tc.expected) {
				t.Errorf("Expected length %d, got %d", len(tc.expected), len(result))
			}
		})
	}
}
