package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/rs/zerolog"
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
		ctx := enrichCtxLogger(r.Context())
		r = r.WithContext(ctx)

		cw := &capturingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		reqBody, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewReader(reqBody))

		fields := log.Ctx(ctx).Debug().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Any("req_headers", r.Header).
			Str("req_body", string(reqBody))

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
			ctx = enrichCtxLogger(ctx)

			resp, err = handler(ctx, req)

			var event *zerolog.Event
			if err != nil {
				event = log.Ctx(ctx).Error()
			} else {
				event = log.Ctx(ctx).Debug()
			}

			event.
				Str("method", info.FullMethod).
				Any("request", req).
				Err(err).
				Any("response", resp).
				Msg("incoming GRPC request")

			return resp, err
		},
	)
}

func enrichCtxLogger(ctx context.Context) context.Context {
	logger := log.Logger

	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		sc := span.SpanContext()
		logger = logger.With().
			Str("trace_id", sc.TraceID().String()).
			Str("span_id", sc.SpanID().String()).
			Logger()
	}

	return logger.WithContext(ctx)
}
