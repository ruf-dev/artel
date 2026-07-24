package cryptoutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"

	"go.redsock.ru/rerrors"
)

const (
	KeySize128 = 16
	KeySize192 = 24
	KeySize256 = 32
)

// Encryptor abstracts credential encryption so the repository layer never has to
// know whether it's running with a real key or in plaintext fallback mode.
type Encryptor interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
	IsPlainText() bool
}

// AESEncryptor is the default Encryptor, backed by AES-GCM.
type AESEncryptor struct {
	key []byte
}

// NewAESEncryptor validates key and returns an Encryptor backed by AES-GCM.
func NewAESEncryptor(key []byte) (Encryptor, error) {
	if !IsValidKeySize(len(key)) {
		msg := fmt.Sprintf(
			"ENVIRONMENT_CREDS_ENCRYPTION_KEY must decode to 16, 24, or 32 bytes for AES-128/192/256; "+
				"got %d bytes. Generate one with: openssl rand -hex 32",
			len(key),
		)

		return nil, rerrors.New(msg)
	}

	e := AESEncryptor{
		key: key,
	}

	return e, nil
}

func (e AESEncryptor) Encrypt(plaintext []byte) ([]byte, error) {
	return Encrypt(e.key, plaintext)
}

func (e AESEncryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	return Decrypt(e.key, ciphertext)
}

func (e AESEncryptor) IsPlainText() bool {
	return false
}

// NoOpEncryptor is a plaintext-passthrough Encryptor, used when no encryption
// key is configured. It exists so "no key configured" degrades to a visible
// insecure mode instead of a hard startup failure.
type NoOpEncryptor struct{}

func (e NoOpEncryptor) Encrypt(plaintext []byte) ([]byte, error) {
	return plaintext, nil
}

func (e NoOpEncryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	return ciphertext, nil
}

func (e NoOpEncryptor) IsPlainText() bool {
	return true
}

// IsValidKeySize reports whether n is a valid AES key length (128/192/256-bit).
func IsValidKeySize(n int) bool {
	return n == KeySize128 || n == KeySize192 || n == KeySize256
}

func Encrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, rerrors.Wrap(err, "error creating AES cipher for encrypt function")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, rerrors.Wrap(err, "error creating GCM")
	}

	nonce := make([]byte, gcm.NonceSize())

	_, err = rand.Read(nonce)
	if err != nil {
		return nil, rerrors.Wrap(err, "error generating nonce")
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	return append(nonce, ciphertext...), nil
}

func Decrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, rerrors.Wrap(err, "error creating AES cipher for decrypt function")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, rerrors.Wrap(err, "error creating GCM")
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, rerrors.New("ciphertext too short")
	}

	nonce := ciphertext[:nonceSize]
	data := ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, rerrors.Wrap(err, "error decrypting")
	}

	return plaintext, nil
}
