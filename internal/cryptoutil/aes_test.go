package cryptoutil_test

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/ruf-dev/artel/internal/cryptoutil"
)

func TestEncryptDecryptRoundTrip(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		t.Fatal(err)
	}

	plaintext := []byte("hello, AES-256-GCM!")

	ciphertext, err := cryptoutil.Encrypt(key, plaintext)
	if err != nil {
		t.Fatal(err)
	}

	result, err := cryptoutil.Decrypt(key, ciphertext)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(plaintext, result) {
		t.Fatalf("expected %s, got %s", plaintext, result)
	}
}
