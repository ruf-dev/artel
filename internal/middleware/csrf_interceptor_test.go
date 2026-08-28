package middleware

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	csrfTestExemptMethod    = "/test.Svc/Login"
	csrfTestProtectedMethod = "/test.Svc/DoWrite"
)

// csrfPassHandler is a grpc.UnaryHandler that records whether it ran and returns a sentinel value.
func csrfPassHandler(ran *bool) grpc.UnaryHandler {
	return func(_ context.Context, _ any) (any, error) {
		*ran = true
		return "ok", nil
	}
}

func csrfTestInfo(fullMethod string) *grpc.UnaryServerInfo {
	return &grpc.UnaryServerInfo{FullMethod: fullMethod}
}

// authedViaCookieMD builds the incoming metadata a browser request carries once
// CookieToMetadataAnnotator has run: the "authenticated via cookie" marker plus the raw Cookie
// header. csrfCookie / csrfHeader are the double-submit pair — pass "" for either to omit it.
func authedViaCookieMD(csrfCookie, csrfHeader string) metadata.MD {
	pairs := []string{AuthViaCookieMarkerKey, AuthViaCookieMarkerValue}
	if csrfCookie != "" {
		pairs = append(pairs, GatewayCookieMetadataKey, CSRFCookieName+"="+csrfCookie)
	}
	if csrfHeader != "" {
		pairs = append(pairs, csrfHeaderMetadataKey, csrfHeader)
	}

	return metadata.Pairs(pairs...)
}

// TestCSRFInterceptor_NoCookieAuth_Passes: a native gRPC/CLI caller (raw Authorization header, no
// cookie marker) is never subject to the CSRF check, even on a protected method with no token.
func TestCSRFInterceptor_NoCookieAuth_Passes(t *testing.T) {
	interceptor := NewCSRFInterceptor(csrfTestExemptMethod)

	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{})

	var ran bool
	_, err := interceptor(ctx, nil, csrfTestInfo(csrfTestProtectedMethod), csrfPassHandler(&ran))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ran {
		t.Fatal("handler did not run")
	}
}

// TestCSRFInterceptor_CookieAuth_ProtectedMethod_MissingToken_Denied: the core deny-by-default —
// a cookie-authenticated write with no csrf_token pair is rejected.
func TestCSRFInterceptor_CookieAuth_ProtectedMethod_MissingToken_Denied(t *testing.T) {
	interceptor := NewCSRFInterceptor(csrfTestExemptMethod)

	ctx := metadata.NewIncomingContext(context.Background(), authedViaCookieMD("", ""))

	var ran bool
	_, err := interceptor(ctx, nil, csrfTestInfo(csrfTestProtectedMethod), csrfPassHandler(&ran))
	if status.Code(err) != codes.PermissionDenied {
		t.Fatalf("want PermissionDenied, got %v", err)
	}
	if ran {
		t.Fatal("handler ran despite missing CSRF token")
	}
}

// TestCSRFInterceptor_CookieAuth_ProtectedMethod_ValidPair_Passes: matching cookie + header lets
// a cookie-authenticated write through.
func TestCSRFInterceptor_CookieAuth_ProtectedMethod_ValidPair_Passes(t *testing.T) {
	interceptor := NewCSRFInterceptor(csrfTestExemptMethod)

	ctx := metadata.NewIncomingContext(context.Background(), authedViaCookieMD("tok-123", "tok-123"))

	var ran bool
	_, err := interceptor(ctx, nil, csrfTestInfo(csrfTestProtectedMethod), csrfPassHandler(&ran))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ran {
		t.Fatal("handler did not run despite a valid CSRF pair")
	}
}

// TestCSRFInterceptor_CookieAuth_ProtectedMethod_MismatchedPair_Denied: cookie and header both
// present but unequal is still a denial.
func TestCSRFInterceptor_CookieAuth_ProtectedMethod_MismatchedPair_Denied(t *testing.T) {
	interceptor := NewCSRFInterceptor(csrfTestExemptMethod)

	ctx := metadata.NewIncomingContext(context.Background(), authedViaCookieMD("tok-123", "tok-999"))

	var ran bool
	_, err := interceptor(ctx, nil, csrfTestInfo(csrfTestProtectedMethod), csrfPassHandler(&ran))
	if status.Code(err) != codes.PermissionDenied {
		t.Fatalf("want PermissionDenied, got %v", err)
	}
	if ran {
		t.Fatal("handler ran despite a mismatched CSRF pair")
	}
}

// TestCSRFInterceptor_CookieAuth_ExemptMethod_MissingToken_Passes guards the Login/Register fix:
// a stale/orphaned access_token cookie marks the request "authenticated via cookie" even during
// re-login, so the pre-session auth RPCs must stay reachable without a csrf_token pair.
func TestCSRFInterceptor_CookieAuth_ExemptMethod_MissingToken_Passes(t *testing.T) {
	interceptor := NewCSRFInterceptor(csrfTestExemptMethod)

	ctx := metadata.NewIncomingContext(context.Background(), authedViaCookieMD("", ""))

	var ran bool
	_, err := interceptor(ctx, nil, csrfTestInfo(csrfTestExemptMethod), csrfPassHandler(&ran))
	if err != nil {
		t.Fatalf("unexpected error on exempt method: %v", err)
	}
	if !ran {
		t.Fatal("handler did not run for an exempt method")
	}
}
