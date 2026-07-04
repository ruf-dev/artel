package telegram

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.redsock.ru/rerrors"
	"go.redsock.ru/toolbox"

	"github.com/ruf-dev/artel/internal/utils"
)

type TgClaims struct {
	jwt.RegisteredClaims

	Id    int64  `json:"-"`
	Login string `json:"-"`

	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// telegramID unmarshals the "id" claim, which Telegram encodes inconsistently
// as either a JSON number or a JSON string.
type telegramID int64

func (t *telegramID) UnmarshalJSON(data []byte) error {
	var asInt int64

	err := json.Unmarshal(data, &asInt)
	if err == nil {
		*t = telegramID(asInt)

		return nil
	}

	var asString string

	err = json.Unmarshal(data, &asString)
	if err != nil {
		return rerrors.New(fmt.Sprintf("error reading id claim: neither a number nor a string: %s", string(data)))
	}

	parsed, err := strconv.ParseInt(asString, 10, 64)
	if err != nil {
		return rerrors.Wrap(err, fmt.Sprintf("error parsing telegram id claim %q", asString))
	}

	*t = telegramID(parsed)

	return nil
}

// UnmarshalJSON is defined explicitly so the "id" claim can be decoded via telegramID
// while every other field keeps the default struct-tag-driven unmarshaling.
func (c *TgClaims) UnmarshalJSON(data []byte) error {
	type alias TgClaims

	aux := struct {
		Id telegramID `json:"id"`
		*alias
	}{
		alias: (*alias)(c),
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return rerrors.Wrap(err, "error unmarshaling telegram claims")
	}

	c.Id = int64(aux.Id)

	return nil
}

type JWKSResponse struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	// RSA fields
	N string `json:"n,omitempty"`
	E string `json:"e,omitempty"`
	// EC / OKP fields
	Crv string   `json:"crv,omitempty"`
	X   string   `json:"x,omitempty"`
	Y   string   `json:"y,omitempty"`
	X5c []string `json:"x5c,omitempty"`
}

type TokenParser struct {
	jwksURL     string
	issuer      string
	audience    string
	jwksCache   *JWKSResponse
	cacheExpiry time.Time
}

func NewTokenParser(jwksURL, issuer, audience string) TokenParser {
	return TokenParser{
		jwksURL:  jwksURL,
		issuer:   issuer,
		audience: audience,
	}
}

func (tp *TokenParser) ParseAndVerifyIdToken(idToken string) (TgClaims, error) {
	claimsTarget := &TgClaims{}

	token, err := jwt.ParseWithClaims(idToken, claimsTarget, tp.keyFunc)
	if err != nil {
		return TgClaims{}, rerrors.Wrap(err, "error parsing id token")
	}

	if !token.Valid {
		return TgClaims{}, rerrors.New("token is invalid")
	}

	claims, ok := token.Claims.(*TgClaims)
	if !ok {
		return TgClaims{}, rerrors.New("error extracting claims")
	}

	if tp.issuer != "" && claims.Issuer != tp.issuer {
		return TgClaims{}, rerrors.New(
			fmt.Sprintf("error validating issuer: expected %s, got %s", tp.issuer, claims.Issuer),
		)
	}

	if tp.audience != "" && claims.Audience != nil {
		audienceFound := false

		for _, aud := range claims.Audience {
			if aud == tp.audience {
				audienceFound = true

				break
			}
		}

		if !audienceFound {
			return TgClaims{}, rerrors.New(fmt.Sprintf("error validating audience: expected %s", tp.audience))
		}
	}

	claims.Login = toolbox.Coalesce(claims.Name, claims.GivenName)

	if claims.Id == 0 {
		return *claims, rerrors.New("error validating claims: missing telegram id")
	}

	return *claims, nil
}

func (tp *TokenParser) keyFunc(token *jwt.Token) (interface{}, error) {
	switch token.Method.(type) {
	case *jwt.SigningMethodRSA, *jwt.SigningMethodECDSA, *jwt.SigningMethodEd25519:
	default:
		return nil, rerrors.New(fmt.Sprintf("error verifying signing method: unsupported algorithm %v", token.Header["alg"]))
	}

	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, rerrors.New("error reading kid header: not found")
	}

	jwks, err := tp.getJWKS()
	if err != nil {
		return nil, rerrors.Wrap(err, "error fetching JWKS")
	}

	jwk := findKey(jwks, kid)
	if jwk == nil {
		// kid not in cache — force a refresh and try once more (handles key rotation)
		tp.jwksCache = nil

		jwks, err = tp.getJWKS()
		if err != nil {
			return nil, rerrors.Wrap(err, "error fetching fresh JWKS")
		}

		jwk = findKey(jwks, kid)
	}

	if jwk == nil {
		return nil, rerrors.New(fmt.Sprintf("error finding key for kid: %s", kid))
	}

	return tp.jwkToPublicKey(jwk)
}

func findKey(jwks *JWKSResponse, kid string) *JWK {
	for i := range jwks.Keys {
		if jwks.Keys[i].Kid == kid {
			return &jwks.Keys[i]
		}
	}

	return nil
}

func (tp *TokenParser) getJWKS() (*JWKSResponse, error) {
	if tp.jwksCache != nil && time.Now().Before(tp.cacheExpiry) {
		return tp.jwksCache, nil
	}

	resp, err := http.Get(tp.jwksURL)
	if err != nil {
		return nil, rerrors.Wrap(err, "error fetching JWKS")
	}
	defer utils.CloseWithLog(resp.Body, "error fetching JWKS")

	var jwks JWKSResponse

	err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		return nil, rerrors.Wrap(err, "error decoding JWKS response")
	}

	tp.jwksCache = &jwks
	tp.cacheExpiry = time.Now().Add(24 * time.Hour)

	return &jwks, nil
}

func (tp *TokenParser) jwkToPublicKey(jwk *JWK) (interface{}, error) {
	switch jwk.Kty {
	case "RSA":
		return tp.jwkToRSAKey(jwk)
	case "EC":
		return tp.jwkToECKey(jwk)
	case "OKP":
		return tp.jwkToOKPKey(jwk)
	default:
		return nil, rerrors.New(fmt.Sprintf("error converting JWK: unsupported key type %s", jwk.Kty))
	}
}

func (tp *TokenParser) jwkToRSAKey(jwk *JWK) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, rerrors.Wrap(err, "error decoding RSA modulus")
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, rerrors.Wrap(err, "error decoding RSA exponent")
	}

	n := new(big.Int).SetBytes(nBytes)
	e := int(new(big.Int).SetBytes(eBytes).Int64())

	publicKey := &rsa.PublicKey{
		N: n,
		E: e,
	}

	return publicKey, nil
}

func (tp *TokenParser) jwkToECKey(jwk *JWK) (*ecdsa.PublicKey, error) {
	var curve elliptic.Curve

	switch jwk.Crv {
	case "P-256":
		curve = elliptic.P256()
	case "P-384":
		curve = elliptic.P384()
	case "P-521":
		curve = elliptic.P521()
	default:
		return nil, rerrors.New(fmt.Sprintf("error converting EC JWK: unsupported curve %s", jwk.Crv))
	}

	xBytes, err := base64.RawURLEncoding.DecodeString(jwk.X)
	if err != nil {
		return nil, rerrors.Wrap(err, "error decoding EC x coordinate")
	}

	yBytes, err := base64.RawURLEncoding.DecodeString(jwk.Y)
	if err != nil {
		return nil, rerrors.Wrap(err, "error decoding EC y coordinate")
	}

	publicKey := &ecdsa.PublicKey{
		Curve: curve,
		X:     new(big.Int).SetBytes(xBytes),
		Y:     new(big.Int).SetBytes(yBytes),
	}

	return publicKey, nil
}

func (tp *TokenParser) jwkToOKPKey(jwk *JWK) (ed25519.PublicKey, error) {
	if jwk.Crv != "Ed25519" {
		return nil, rerrors.New(fmt.Sprintf("error converting OKP JWK: unsupported curve %s", jwk.Crv))
	}

	xBytes, err := base64.RawURLEncoding.DecodeString(jwk.X)
	if err != nil {
		return nil, rerrors.Wrap(err, "error decoding Ed25519 public key")
	}

	return ed25519.PublicKey(xBytes), nil
}
