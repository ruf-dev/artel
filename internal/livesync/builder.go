package livesync

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/url"

	"go.redsock.ru/rerrors"
)

const (
	uriBase    = "obsidian://setuplivesync?settings="
	pbkdf2Iter = 100_000
	saltSize   = 16 // bytes → 32 hex chars in output
	ivSize     = 16 // bytes → 32 hex chars in output
	aesKeyLen  = 32 // bytes, AES-256

	defaultBatchSize                          = 50
	defaultBatchesLimit                       = 50
	defaultCustomChunkSize                    = 50
	defaultConcurrencyOfReadChunksOnline      = 100
	defaultMinimumIntervalOfReadChunksOnline  = 100
	defaultSettingVersion                     = 10
	defaultNotifyThresholdOfRemoteStorageSize = 800

	byteShift3 = 24
	byteShift2 = 16
	byteShift1 = 8
)

// Config is the configuration object embedded in the Setup URI.
// Fields match exactly what the TypeScript generate_setupuri.ts produces.
type Config struct {
	CouchDBURI      string `json:"couchDB_URI"`
	CouchDBUser     string `json:"couchDB_USER"`
	CouchDBPassword string `json:"couchDB_PASSWORD"`
	CouchDBDBName   string `json:"couchDB_DBNAME"`

	// Sync behaviour — defaults mirror the official Deno script
	SyncOnStart                        bool   `json:"syncOnStart"`
	GCDelay                            int    `json:"gcDelay"`
	PeriodicReplication                bool   `json:"periodicReplication"`
	SyncOnFileOpen                     bool   `json:"syncOnFileOpen"`
	Encrypt                            bool   `json:"encrypt"`
	Passphrase                         string `json:"passphrase"` // E2EE passphrase (can differ from URI passphrase)
	UsePathObfuscation                 bool   `json:"usePathObfuscation"`
	BatchSave                          bool   `json:"batchSave"`
	BatchSize                          int    `json:"batch_size"`
	BatchesLimit                       int    `json:"batches_limit"`
	UseHistory                         bool   `json:"useHistory"`
	DisableRequestURI                  bool   `json:"disableRequestURI"`
	CustomChunkSize                    int    `json:"customChunkSize"`
	SyncAfterMerge                     bool   `json:"syncAfterMerge"`
	ConcurrencyOfReadChunksOnline      int    `json:"concurrencyOfReadChunksOnline"`
	MinimumIntervalOfReadChunksOnline  int    `json:"minimumIntervalOfReadChunksOnline"`
	HandleFilenameCaseSensitive        bool   `json:"handleFilenameCaseSensitive"`
	DoNotUseFixedRevisionForChunks     bool   `json:"doNotUseFixedRevisionForChunks"`
	SettingVersion                     int    `json:"settingVersion"`
	NotifyThresholdOfRemoteStorageSize int    `json:"notifyThresholdOfRemoteStorageSize"`
}

// DefaultConfig returns a Config pre-filled with the same defaults
// as the official generate_setupuri.ts script.
func DefaultConfig(couchURI, user, password, dbName, e2eePassphrase string) Config {
	return Config{
		CouchDBURI:      couchURI,
		CouchDBUser:     user,
		CouchDBPassword: password,
		CouchDBDBName:   dbName,

		SyncOnStart:                        true,
		GCDelay:                            0,
		PeriodicReplication:                true,
		SyncOnFileOpen:                     true,
		Encrypt:                            true,
		Passphrase:                         e2eePassphrase,
		UsePathObfuscation:                 true,
		BatchSave:                          true,
		BatchSize:                          defaultBatchSize,
		BatchesLimit:                       defaultBatchesLimit,
		UseHistory:                         true,
		DisableRequestURI:                  true,
		CustomChunkSize:                    defaultCustomChunkSize,
		SyncAfterMerge:                     false,
		ConcurrencyOfReadChunksOnline:      defaultConcurrencyOfReadChunksOnline,
		MinimumIntervalOfReadChunksOnline:  defaultMinimumIntervalOfReadChunksOnline,
		HandleFilenameCaseSensitive:        false,
		DoNotUseFixedRevisionForChunks:     false,
		SettingVersion:                     defaultSettingVersion,
		NotifyThresholdOfRemoteStorageSize: defaultNotifyThresholdOfRemoteStorageSize,
	}
}

// GenerateSetupURI creates a complete LiveSync Setup URI from a config and URI passphrase.
//
// uriPassphrase encrypts the URI itself (not the vault data).
// Use a different passphrase from the E2EE passphrase in the config.
//
// Returns the URI and any error.
func GenerateSetupURI(cfg Config, uriPassphrase string) (string, error) {
	plaintext, err := json.Marshal(cfg)
	if err != nil {
		return "", rerrors.Wrap(err, "error marshaling livesync config")
	}

	encrypted, err := encryptOctagonalWheels(string(plaintext), uriPassphrase)
	if err != nil {
		return "", rerrors.Wrap(err, "error encrypting setup uri")
	}

	return uriBase + url.QueryEscape(encrypted), nil
}

// pbkdf2Key is a minimal PBKDF2-SHA256 implementation using only Go stdlib.
// Mirrors golang.org/x/crypto/pbkdf2.Key with sha256.New.
func pbkdf2Key(password, salt []byte, iter, keyLen int) []byte {
	prf := hmac.New(sha256.New, password)
	hashLen := prf.Size()
	numBlocks := (keyLen + hashLen - 1) / hashLen

	var buf [4]byte

	dk := make([]byte, 0, numBlocks*hashLen)
	U := make([]byte, hashLen)

	for block := 1; block <= numBlocks; block++ {
		prf.Reset()
		prf.Write(salt)

		buf[0] = byte(block >> byteShift3)
		buf[1] = byte(block >> byteShift2)
		buf[2] = byte(block >> byteShift1)
		buf[3] = byte(block)
		prf.Write(buf[:4])
		dk = prf.Sum(dk)
		T := dk[len(dk)-hashLen:]
		copy(U, T)

		for n := 2; n <= iter; n++ {
			prf.Reset()
			prf.Write(U)
			U = U[:0]
			U = prf.Sum(U)

			for x := range U {
				T[x] ^= U[x]
			}
		}
	}

	return dk[:keyLen]
}

// encryptOctagonalWheels replicates the encrypt() function from
// octagonal-wheels@0.1.30/dist/encryption/encryption.js
//
// Algorithm:
//
//	digest  = SHA-256(passphrase_utf8)
//	key     = PBKDF2(digest, salt, 100_000, 32, SHA-256)  -- AES-256-GCM key
//	iv      = random 16 bytes
//	cipher  = AES-GCM-256(key, iv, plaintext_utf8)
//	output  = "%" + hex(iv) + hex(salt) + base64Std(cipher)
func encryptOctagonalWheels(plaintext, passphrase string) (string, error) {
	// 1. SHA-256 hash of the passphrase (mirrors the JS WebCrypto importKey path)
	passphraseHash := sha256.Sum256([]byte(passphrase))

	// 2. Random salt (16 bytes)
	salt := make([]byte, saltSize)

	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return "", rerrors.Wrap(err, "error generating salt")
	}

	// 3. PBKDF2 → AES-256 key (100_000 iterations, SHA-256)
	key := pbkdf2Key(passphraseHash[:], salt, pbkdf2Iter, aesKeyLen)

	// 4. Random IV (16 bytes — matches the JS 12-byte semi-static + 4-byte nonce pattern,
	//    but the output format just stores raw iv bytes; using 16 random bytes is equivalent
	//    for Go's AES-GCM which accepts any IV size via NewGCMWithNonceSize)
	iv := make([]byte, ivSize)

	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return "", rerrors.Wrap(err, "error generating iv")
	}

	// 5. AES-256-GCM encrypt
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", rerrors.Wrap(err, "error creating aes cipher for octogonal wheels")
	}

	gcm, err := cipher.NewGCMWithNonceSize(block, ivSize)
	if err != nil {
		return "", rerrors.Wrap(err, "error creating gcm")
	}

	ciphertext := gcm.Seal(nil, iv, []byte(plaintext), nil)

	// 6. Assemble: "%" + hex(iv) + hex(salt) + base64Std(ciphertext)
	b64 := base64.StdEncoding.EncodeToString(ciphertext)
	result := "%" + hex.EncodeToString(iv) + hex.EncodeToString(salt) + b64

	return result, nil
}
