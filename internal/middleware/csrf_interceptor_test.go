package middleware_test

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/ruf-dev/artel/internal/middleware"
)

func handlerOK(_ context.Context, _ any) (any, error) {
	return "ok", nil
}

func TestCSRFInterceptor_PassesThroughWhenNotAuthenticatedViaCookie(t *testing.T) {
	interceptor := middleware.NewCSRFInterceptor()
	info := &grpc.UnaryServerInfo{FullMethod: "/artel.Vaults/CreateVault"}

	// No AuthViaCookieMarkerKey metadata at all — a native gRPC/header caller.
	md := metadata.Pairs("authorization", "some-token")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	resp, err := interceptor(ctx, nil, info, handlerOK)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "ok" {
		t.Fatalf("expected handler response to pass through, got %v", resp)
	}
}

func TestCSRFInterceptor_PassesThroughForExemptMethod(t *testing.T) {
	interceptor := middleware.NewCSRFInterceptor("/artel.Vaults/ListVaults")
	info := &grpc.UnaryServerInfo{FullMethod: "/artel.Vaults/ListVaults"}

	md := metadata.Pairs(middleware.AuthViaCookieMarkerKey, middleware.AuthViaCookieMarkerValue)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	resp, err := interceptor(ctx, nil, info, handlerOK)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "ok" {
		t.Fatalf("expected handler response to pass through, got %v", resp)
	}
}

func TestCSRFInterceptor_DeniesWhenCookieAuthenticatedAndTokenMissing(t *testing.T) {
	interceptor := middleware.NewCSRFInterceptor()
	info := &grpc.UnaryServerInfo{FullMethod: "/artel.Vaults/CreateVault"}

	md := metadata.Pairs(middleware.AuthViaCookieMarkerKey, middleware.AuthViaCookieMarkerValue)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	_, err := interceptor(ctx, nil, info, handlerOK)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("expected a grpc status error")
	}
	if st.Code() != codes.PermissionDenied {
		t.Fatalf("expected codes.PermissionDenied, got %v", st.Code())
	}
}

func TestCSRFInterceptor_DeniesWhenTokenMismatch(t *testing.T) {
	interceptor := middleware.NewCSRFInterceptor()
	info := &grpc.UnaryServerInfo{FullMethod: "/artel.Vaults/CreateVault"}

	md := metadata.Pairs(
		middleware.AuthViaCookieMarkerKey, middleware.AuthViaCookieMarkerValue,
		middleware.GatewayCookieMetadataKey, "csrf_token=correct-value",
		"x-csrf-token", "wrong-value",
	)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	_, err := interceptor(ctx, nil, info, handlerOK)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("expected a grpc status error")
	}
	if st.Code() != codes.PermissionDenied {
		t.Fatalf("expected codes.PermissionDenied, got %v", st.Code())
	}
}

func TestCSRFInterceptor_AllowsWhenTokenMatches(t *testing.T) {
	interceptor := middleware.NewCSRFInterceptor()
	info := &grpc.UnaryServerInfo{FullMethod: "/artel.Vaults/CreateVault"}

	md := metadata.Pairs(
		middleware.AuthViaCookieMarkerKey, middleware.AuthViaCookieMarkerValue,
		middleware.GatewayCookieMetadataKey, "csrf_token=matching-value",
		"x-csrf-token", "matching-value",
	)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	resp, err := interceptor(ctx, nil, info, handlerOK)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "ok" {
		t.Fatalf("expected handler response to pass through, got %v", resp)
	}
}

func TestCSRFInterceptor_MissingIncomingMetadataPassesThrough(t *testing.T) {
	interceptor := middleware.NewCSRFInterceptor()
	info := &grpc.UnaryServerInfo{FullMethod: "/artel.Vaults/CreateVault"}

	resp, err := interceptor(context.Background(), nil, info, handlerOK)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "ok" {
		t.Fatalf("expected handler response to pass through, got %v", resp)
	}
}
