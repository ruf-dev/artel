package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type capturingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (w *capturingResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *capturingResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func HttpLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cw := &capturingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		reqBody, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewReader(reqBody))

		fields := log.Debug().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Any("req_headers", r.Header).
			Str("req_body", string(reqBody))

		if span := trace.SpanFromContext(r.Context()); span.SpanContext().IsValid() {
			sc := span.SpanContext()
			fields = fields.
				Str("trace_id", sc.TraceID().String()).
				Str("span_id", sc.SpanID().String())
		}

		defer func() {
			fields.Msg("incoming HTTP request")
		}()

		next.ServeHTTP(cw, r)

		fields = fields.
			Int("status", cw.statusCode).
			Any("resp_headers", cw.Header()).
			Str("resp_body", cw.body.String())
	})
}

func LogInterceptor() grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
			fields := log.Debug().
				Str("method", info.FullMethod).
				Any("request", req)

			if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
				sc := span.SpanContext()
				fields = fields.
					Str("trace_id", sc.TraceID().String()).
					Str("span_id", sc.SpanID().String())
			}

			defer func() {
				fields.Msg("incoming GRPC request")
			}()

			resp, err = handler(ctx, req)

			fields = fields.
				Err(err).
				Any("response", resp)

			return resp, err
		})
}
