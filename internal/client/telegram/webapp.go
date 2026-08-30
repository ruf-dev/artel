package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.redsock.ru/rerrors"
	"go.redsock.ru/toolbox"

	"github.com/ruf-dev/artel/internal/service/user_errors"
)

const webAppInitDataMaxAge = 24 * time.Hour

// webAppUser is the JSON shape of the "user" field inside a Telegram Mini App initData string.
// The id reuses telegramID because Telegram encodes it inconsistently as a number or a string.
type webAppUser struct {
	Id           telegramID `json:"id"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Username     string     `json:"username"`
	PhotoUrl     string     `json:"photo_url"`
	LanguageCode string     `json:"language_code"`
}

// ParseAndVerifyWebAppInitData validates a Telegram Mini App initData string against the bot
// token using Telegram's HMAC-SHA256 data-check-string scheme and returns the embedded Telegram
// user as TgClaims. The newer Ed25519 "signature" field is ignored — only the HMAC check runs.
func ParseAndVerifyWebAppInitData(initData string, botToken string) (*TgClaims, error) {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.InvalidTelegramToken, "error parsing init data", err.Error())
	}

	receivedHash := values.Get("hash")
	if receivedHash == "" {
		return nil, rerrors.Wrap(user_errors.InvalidTelegramToken, "init data carries no hash")
	}

	values.Del("hash")
	values.Del("signature")

	dataCheckString := buildDataCheckString(values)

	secretKey := hmacSha256([]byte("WebAppData"), []byte(botToken))
	computed := hex.EncodeToString(hmacSha256(secretKey, []byte(dataCheckString)))

	if !hmac.Equal([]byte(computed), []byte(receivedHash)) {
		return nil, rerrors.Wrap(user_errors.InvalidTelegramToken, "init data hash mismatch")
	}

	authDateUnix, err := strconv.ParseInt(values.Get("auth_date"), 10, 64)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.InvalidTelegramToken, "error parsing auth_date", err.Error())
	}

	authDate := time.Unix(authDateUnix, 0)
	if time.Since(authDate) > webAppInitDataMaxAge {
		return nil, rerrors.Wrap(user_errors.InvalidTelegramToken, "init data is older than 24h")
	}

	userRaw := values.Get("user")
	if userRaw == "" {
		return nil, rerrors.Wrap(user_errors.InvalidTelegramToken, "init data carries no user")
	}

	var user webAppUser

	err = json.Unmarshal([]byte(userRaw), &user)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.InvalidTelegramToken, "error parsing user", err.Error())
	}

	if user.Id == 0 {
		return nil, rerrors.Wrap(user_errors.InvalidTelegramToken, "init data user is missing an id")
	}

	claims := &TgClaims{
		Id:      int64(user.Id),
		Login:   toolbox.Coalesce(user.Username, user.FirstName),
		Picture: user.PhotoUrl,
	}

	return claims, nil
}

// buildDataCheckString joins every remaining key as "key=value", sorted ascending, with "\n" —
// the exact string Telegram signs.
func buildDataCheckString(values url.Values) string {
	keys := make([]string, 0, len(values))

	for key := range values {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	pairs := make([]string, 0, len(keys))

	for _, key := range keys {
		pairs = append(pairs, key+"="+values.Get(key))
	}

	return strings.Join(pairs, "\n")
}

func hmacSha256(key, message []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)

	return mac.Sum(nil)
}
