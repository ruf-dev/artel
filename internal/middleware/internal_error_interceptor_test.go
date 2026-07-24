package middleware_test

import (
	"context"
	"errors"
	"testing"

	"go.redsock.ru/rerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ruf-dev/artel/internal/middleware"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

func TestRedactInternalErrors_UnclassifiedErrorIsRedacted(t *testing.T) {
	handlerErr := rerrors.Wrap(errors.New("crypto/aes: invalid key size 0"), "error creating AES cipher for encrypt function")
	handler := func(_ context.Context, _ any) (any, error) {
		return nil, handlerErr
	}
	info := &grpc.UnaryServerInfo{FullMethod: "/artel.Test/Method"}

	resp, err := middleware.RedactInternalErrors(context.Background(), nil, info, handler)

	if resp != nil {
		t.Fatalf("expected nil resp, got %v", resp)
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("expected a grpc status error")
	}
	if st.Code() != codes.Internal {
		t.Fatalf("expected codes.Internal, got %v", st.Code())
	}
	if st.Message() != "internal error" {
		t.Fatalf("expected client-facing message to be redacted, got: %q", st.Message())
	}
}

func TestRedactInternalErrors_ClassifiedUserErrorPassesThrough(t *testing.T) {
	handler := func(_ context.Context, _ any) (any, error) {
		return nil, user_errors.NotFound
	}
	info := &grpc.UnaryServerInfo{FullMethod: "/artel.Test/Method"}

	_, err := middleware.RedactInternalErrors(context.Background(), nil, info, handler)

	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("expected a grpc status error")
	}
	if st.Code() != codes.NotFound {
		t.Fatalf("expected codes.NotFound to pass through unchanged, got %v", st.Code())
	}
	if st.Message() != "not found" {
		t.Fatalf("expected sentinel message to pass through unchanged, got: %q", st.Message())
	}
}

func TestRedactInternalErrors_ClassifiedStatusErrorPassesThrough(t *testing.T) {
	handler := func(_ context.Context, _ any) (any, error) {
		return nil, status.Error(codes.PermissionDenied, "not an administrator")
	}
	info := &grpc.UnaryServerInfo{FullMethod: "/artel.Test/Method"}

	_, err := middleware.RedactInternalErrors(context.Background(), nil, info, handler)

	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("expected a grpc status error")
	}
	if st.Code() != codes.PermissionDenied {
		t.Fatalf("expected codes.PermissionDenied to pass through, got %v", st.Code())
	}
}

func TestRedactInternalErrors_NilErrorPassesThrough(t *testing.T) {
	wantResp := "ok"
	handler := func(_ context.Context, _ any) (any, error) {
		return wantResp, nil
	}
	info := &grpc.UnaryServerInfo{FullMethod: "/artel.Test/Method"}

	resp, err := middleware.RedactInternalErrors(context.Background(), nil, info, handler)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp != wantResp {
		t.Fatalf("expected resp %q, got %v", wantResp, resp)
	}
}
