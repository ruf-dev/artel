package app

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	otellogGlobal "go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.redsock.ru/rerrors"
	"go.redsock.ru/toolbox/closer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (a *App) initOTel() error {
	endpoint := a.Cfg.Environment.OtelEndpoint
	if endpoint == "" {
		return nil
	}

	res, err := resource.New(context.Background(),
		resource.WithAttributes(semconv.ServiceName("artel")),
		resource.WithProcess(),
		resource.WithHost(),
	)
	if err != nil {
		return rerrors.Wrap(err, "create otel resource")
	}

	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return rerrors.Wrap(err, "create otel grpc connection")
	}

	traceExp, err := otlptracegrpc.New(context.Background(), otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return rerrors.Wrap(err, "create otlp trace exporter")
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExp),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	logExp, err := otlploggrpc.New(context.Background(), otlploggrpc.WithGRPCConn(conn))
	if err != nil {
		return rerrors.Wrap(err, "create otlp log exporter")
	}
	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExp)),
		sdklog.WithResource(res),
	)
	otellogGlobal.SetLoggerProvider(lp)

	closer.Add(func() error {
		_ = tp.Shutdown(context.Background())
		_ = lp.Shutdown(context.Background())
		return conn.Close()
	})

	return nil
}
