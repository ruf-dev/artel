package cryptoutil_test

import (
	"bytes"
	"crypto/rand"
	"strings"
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

func TestIsValidKeySize(t *testing.T) {
	cases := []struct {
		name string
		n    int
		want bool
	}{
		{"zero", 0, false},
		{"aes128", 16, true},
		{"aes192", 24, true},
		{"aes256", 32, true},
		{"too short", 15, false},
		{"between valid sizes", 20, false},
		{"too long", 33, false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := cryptoutil.IsValidKeySize(tc.n)
			if got != tc.want {
				t.Fatalf("IsValidKeySize(%d) = %v, want %v", tc.n, got, tc.want)
			}
		})
	}
}

func TestDecrypt_ErrorMessageMentionsDecrypt(t *testing.T) {
	_, err := cryptoutil.Decrypt(nil, []byte("short"))
	if err == nil {
		t.Fatal("expected an error for a zero-length key")
	}

	msg := err.Error()
	if !strings.Contains(msg, "decrypt function") {
		t.Fatalf("expected error to mention 'decrypt function', got: %s", msg)
	}
	if strings.Contains(msg, "for encrypt function") {
		t.Fatalf("Decrypt error message still contains the copy-pasted Encrypt wording: %s", msg)
	}
}
