package security

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
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

func TestEncryptDecryptLongText(t *testing.T) {
	key := testKey()
	// 10KB plaintext.
	plaintext := strings.Repeat("A long repeated text. ", 500)
	encrypted, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt long text: %v", err)
	}
	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("Decrypt long text: %v", err)
	}
	if decrypted != plaintext {
		t.Fatalf("long text round-trip failed, len(decrypted)=%d, len(original)=%d", len(decrypted), len(plaintext))
	}
}

func TestEncryptDecryptUnicode(t *testing.T) {
	key := testKey()
	plaintext := "日本語テスト🎉 العربية тест 中文測試"
	encrypted, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt unicode: %v", err)
	}
	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("Decrypt unicode: %v", err)
	}
	if decrypted != plaintext {
		t.Fatalf("got %q, want %q", decrypted, plaintext)
	}
}

func TestEncryptProducesUniqueOutput(t *testing.T) {
	key := testKey()
	plaintext := "same input"
	e1, _ := Encrypt(plaintext, key)
	e2, _ := Encrypt(plaintext, key)
	if e1 == e2 {
		t.Error("encrypting same plaintext should produce different ciphertexts due to random nonce")
	}
}

func TestEncryptInvalidKeyLength(t *testing.T) {
	shortKey := []byte("tooshort")
	_, err := Encrypt("test", shortKey)
	if err == nil {
		t.Error("expected error for invalid key length")
	}
}

func TestDecryptTamperedCiphertext(t *testing.T) {
	key := testKey()
	encrypted, err := Encrypt("secret data", key)
	if err != nil {
		t.Fatal(err)
	}
	// Tamper with the ciphertext by flipping a byte.
	data, _ := base64.StdEncoding.DecodeString(encrypted)
	if len(data) > 15 {
		data[15] ^= 0xFF
	}
	tampered := base64.StdEncoding.EncodeToString(data)
	_, err = Decrypt(tampered, key)
	if err == nil {
		t.Error("expected error for tampered ciphertext")
	}
}
