package middleware

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestRequestIsSecure(t *testing.T) {
	tests := []struct {
		name           string
		tls            bool
		forwardedProto string
		want           bool
	}{
		{
			name: "tls set, no header",
			tls:  true,
			want: true,
		},
		{
			name:           "forwarded proto https, no tls",
			forwardedProto: "https",
			want:           true,
		},
		{
			name:           "forwarded proto mixed case Https, no tls",
			forwardedProto: "Https",
			want:           true,
		},
		{
			name:           "forwarded proto http overrides tls",
			tls:            true,
			forwardedProto: "http",
			want:           false,
		},
		{
			name: "neither tls nor header",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.tls {
				r.TLS = &tls.ConnectionState{}
			}
			if tt.forwardedProto != "" {
				r.Header.Set("X-Forwarded-Proto", tt.forwardedProto)
			}

			got := requestIsSecure(r)
			if got != tt.want {
				t.Fatalf("requestIsSecure() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequestSchemeAnnotator(t *testing.T) {
	insecureReq := httptest.NewRequest(http.MethodGet, "/", nil)

	secureReq := httptest.NewRequest(http.MethodGet, "/", nil)
	secureReq.TLS = &tls.ConnectionState{}

	tests := []struct {
		name string
		req  *http.Request
		want metadata.MD
	}{
		{
			name: "insecure request returns nil metadata",
			req:  insecureReq,
			want: nil,
		},
		{
			name: "secure request returns request-secure metadata",
			req:  secureReq,
			want: metadata.Pairs(RequestSecureKey, RequestSecureValue),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RequestSchemeAnnotator(context.Background(), tt.req)

			if len(got) != len(tt.want) {
				t.Fatalf("RequestSchemeAnnotator() = %v, want %v", got, tt.want)
			}

			for k, v := range tt.want {
				if !equalStringSlices(got.Get(k), v) {
					t.Fatalf("RequestSchemeAnnotator() key %q = %v, want %v", k, got.Get(k), v)
				}
			}
		})
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
