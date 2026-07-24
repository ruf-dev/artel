package cryptoutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"go.redsock.ru/rerrors"
)

const (
	KeySize128 = 16
	KeySize192 = 24
	KeySize256 = 32
)

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
