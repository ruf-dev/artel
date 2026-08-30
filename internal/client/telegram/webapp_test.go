package telegram

import (
	"encoding/hex"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseAndVerifyWebAppInitData(t *testing.T) {
	const botToken = "123456:AA-test-bot-token"

	const userJSON = `{"id":42,"first_name":"Ada","last_name":"Lovelace",` +
		`"username":"ada","photo_url":"https://t.me/i/ada.jpg","language_code":"en"}`

	validInitData := signInitData(t, botToken, userJSON, time.Now())
	staleInitData := signInitData(t, botToken, userJSON, time.Now().Add(-48*time.Hour))
	tamperedInitData := tamperHash(t, validInitData)

	cases := []struct {
		name     string
		initData string
		wantErr  bool
	}{
		{"valid init data logs the telegram user in", validInitData, false},
		{"tampered hash is rejected", tamperedInitData, true},
		{"stale auth_date is rejected", staleInitData, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			claims, err := ParseAndVerifyWebAppInitData(tc.initData, botToken)

			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, int64(42), claims.Id)
			require.Equal(t, "ada", claims.Login)
			require.Equal(t, "https://t.me/i/ada.jpg", claims.Picture)
		})
	}
}

// signInitData assembles a Telegram Mini App initData query string carrying the correct HMAC
// hash for botToken, computed exactly as ParseAndVerifyWebAppInitData verifies it.
func signInitData(t *testing.T, botToken, userJSON string, authDate time.Time) string {
	t.Helper()

	values := url.Values{}
	values.Set("user", userJSON)
	values.Set("auth_date", strconv.FormatInt(authDate.Unix(), 10))
	values.Set("query_id", "AAHtest-query-id")
	values.Set("chat_instance", "-1234567890123456789")

	secretKey := hmacSha256([]byte("WebAppData"), []byte(botToken))
	hash := hex.EncodeToString(hmacSha256(secretKey, []byte(buildDataCheckString(values))))

	values.Set("hash", hash)

	return values.Encode()
}

// tamperHash flips the first character of the hash field, leaving every signed field intact.
func tamperHash(t *testing.T, initData string) string {
	t.Helper()

	values, err := url.ParseQuery(initData)
	require.NoError(t, err)

	hashBytes := []byte(values.Get("hash"))
	if hashBytes[0] == 'a' {
		hashBytes[0] = 'b'
	} else {
		hashBytes[0] = 'a'
	}

	values.Set("hash", string(hashBytes))

	return values.Encode()
}
