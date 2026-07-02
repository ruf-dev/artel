package telegram_test

import (
	"encoding/json"
	"testing"

	"github.com/ruf-dev/artel/internal/client/telegram"
)

func TestTgClaims_UnmarshalJSON_IdAsNumber(t *testing.T) {
	raw := []byte(`{"id":123456789,"name":"John"}`)

	var claims telegram.TgClaims
	err := json.Unmarshal(raw, &claims)
	if err != nil {
		t.Fatal(err)
	}

	if claims.Id != 123456789 {
		t.Fatalf("expected id 123456789, got %d", claims.Id)
	}
}

func TestTgClaims_UnmarshalJSON_IdAsString(t *testing.T) {
	raw := []byte(`{"id":"123456789","name":"John"}`)

	var claims telegram.TgClaims
	err := json.Unmarshal(raw, &claims)
	if err != nil {
		t.Fatal(err)
	}

	if claims.Id != 123456789 {
		t.Fatalf("expected id 123456789, got %d", claims.Id)
	}
}

func TestTgClaims_UnmarshalJSON_IdAbsent(t *testing.T) {
	raw := []byte(`{"name":"John"}`)

	var claims telegram.TgClaims
	err := json.Unmarshal(raw, &claims)
	if err != nil {
		t.Fatal(err)
	}

	if claims.Id != 0 {
		t.Fatalf("expected id 0 when absent, got %d", claims.Id)
	}
}
