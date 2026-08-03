package middleware

import (
	"context"
	"net/http"

	"google.golang.org/grpc/metadata"
)

// CookieToMetadataAnnotator is registered via runtime.WithMetadata on the shared gateway mux
// (see transport.NewGatewayMux). It translates the access_token cookie into the "authorization"
// gRPC metadata key that GrpcAuthInterceptor.authWithSession already reads, so browser callers
// authenticate the same way native gRPC/CLI callers do.
//
// It only fires when the request carries neither an "Authorization" nor a
// "Grpc-Metadata-Authorization" header — an explicit header always wins over the cookie,
// deterministically, rather than relying on grpc-gateway's metadata.Join merge order.
func CookieToMetadataAnnotator(_ context.Context, r *http.Request) metadata.MD {
	if r.Header.Get("Authorization") != "" || r.Header.Get("Grpc-Metadata-Authorization") != "" {
		return nil
	}

	cookie, err := r.Cookie(AccessTokenCookieName)
	if err != nil || cookie.Value == "" {
		return nil
	}

	return metadata.Pairs(
		authHeader, cookie.Value,
		AuthViaCookieMarkerKey, AuthViaCookieMarkerValue,
	)
}

// CookieValueFromRawHeader extracts a single cookie's value out of a raw "Cookie" header
// string (as carried in the GatewayCookieMetadataKey gRPC metadata value, rather than a live
// *http.Request). Used by GrpcCSRFInterceptor and auth_api's Refresh handler, which both need
// to read a cookie from gRPC metadata rather than an http.Request.
func CookieValueFromRawHeader(rawCookieHeader, name string) (string, error) {
	header := http.Header{}
	header.Add("Cookie", rawCookieHeader)

	req := http.Request{Header: header}

	cookie, err := req.Cookie(name)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}
