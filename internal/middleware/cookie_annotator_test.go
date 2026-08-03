package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ruf-dev/artel/internal/middleware"
)

func TestCookieToMetadataAnnotator_InjectsFromCookieWhenNoHeaderPresent(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/vaults/list", nil)
	req.AddCookie(&http.Cookie{Name: middleware.AccessTokenCookieName, Value: "cookie-token"})

	md := middleware.CookieToMetadataAnnotator(req.Context(), req)

	got := md.Get("authorization")
	if len(got) != 1 || got[0] != "cookie-token" {
		t.Fatalf("expected authorization=[cookie-token], got %v", got)
	}

	marker := md.Get(middleware.AuthViaCookieMarkerKey)
	if len(marker) != 1 || marker[0] != middleware.AuthViaCookieMarkerValue {
		t.Fatalf("expected %s marker to be set, got %v", middleware.AuthViaCookieMarkerKey, marker)
	}
}

func TestCookieToMetadataAnnotator_SkipsWhenAuthorizationHeaderPresent(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/vaults/list", nil)
	req.Header.Set("Authorization", "header-token")
	req.AddCookie(&http.Cookie{Name: middleware.AccessTokenCookieName, Value: "cookie-token"})

	md := middleware.CookieToMetadataAnnotator(req.Context(), req)

	if md != nil {
		t.Fatalf("expected nil metadata when Authorization header is present, got %v", md)
	}
}

func TestCookieToMetadataAnnotator_SkipsWhenGrpcMetadataAuthorizationHeaderPresent(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/vaults/list", nil)
	req.Header.Set("Grpc-Metadata-Authorization", "header-token")
	req.AddCookie(&http.Cookie{Name: middleware.AccessTokenCookieName, Value: "cookie-token"})

	md := middleware.CookieToMetadataAnnotator(req.Context(), req)

	if md != nil {
		t.Fatalf("expected nil metadata when Grpc-Metadata-Authorization header is present, got %v", md)
	}
}

func TestCookieToMetadataAnnotator_SkipsWhenNoCookiePresent(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/vaults/list", nil)

	md := middleware.CookieToMetadataAnnotator(req.Context(), req)

	if md != nil {
		t.Fatalf("expected nil metadata when no access_token cookie is present, got %v", md)
	}
}

func TestCookieToMetadataAnnotator_SkipsWhenCookieValueEmpty(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/vaults/list", nil)
	req.AddCookie(&http.Cookie{Name: middleware.AccessTokenCookieName, Value: ""})

	md := middleware.CookieToMetadataAnnotator(req.Context(), req)

	if md != nil {
		t.Fatalf("expected nil metadata when access_token cookie value is empty, got %v", md)
	}
}

func TestCookieValueFromRawHeader_ReturnsValue(t *testing.T) {
	v, err := middleware.CookieValueFromRawHeader("access_token=abc; csrf_token=def", "csrf_token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "def" {
		t.Fatalf("expected def, got %q", v)
	}
}

func TestCookieValueFromRawHeader_MissingCookieReturnsError(t *testing.T) {
	_, err := middleware.CookieValueFromRawHeader("access_token=abc", "csrf_token")
	if err == nil {
		t.Fatal("expected error for missing cookie, got nil")
	}
}
