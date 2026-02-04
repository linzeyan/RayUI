package security

import (
	"crypto/rand"
	"testing"
)

func testKey() []byte {
	key := make([]byte, 32)
	rand.Read(key)
	return key
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	key := testKey()
	plaintext := "Hello, RayUI! 測試中文"

	encrypted, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if encrypted == plaintext {
		t.Fatal("encrypted should differ from plaintext")
	}

	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if decrypted != plaintext {
		t.Fatalf("got %q, want %q", decrypted, plaintext)
	}
}

func TestEncryptDecryptEmpty(t *testing.T) {
	key := testKey()

	encrypted, err := Encrypt("", key)
	if err != nil {
		t.Fatalf("Encrypt empty: %v", err)
	}

	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("Decrypt empty: %v", err)
	}
	if decrypted != "" {
		t.Fatalf("got %q, want empty string", decrypted)
	}
}

func TestDecryptWrongKey(t *testing.T) {
	key1 := testKey()
	key2 := testKey()

	encrypted, err := Encrypt("secret", key1)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	_, err = Decrypt(encrypted, key2)
	if err == nil {
		t.Fatal("expected error when decrypting with wrong key")
	}
}

func TestDecryptInvalidBase64(t *testing.T) {
	key := testKey()
	_, err := Decrypt("not-valid-base64!!!", key)
	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
}

func TestDecryptTooShort(t *testing.T) {
	key := testKey()
	_, err := Decrypt("AAAA", key) // valid base64 but too short for nonce+ciphertext
	if err == nil {
		t.Fatal("expected error for too-short ciphertext")
	}
}
