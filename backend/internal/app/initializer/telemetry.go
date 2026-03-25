package initializer

import (
	"context"

	"github.com/potibm/kasseapparat/internal/app/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func InitTelemetry(ctx context.Context, endpoint, version string) (func(), error) {
	if endpoint == "" {
		return nil, nil // No tracing if no endpoint is provided
	}

	res, err := resource.New(ctx, resource.WithAttributes(
        semconv.ServiceName(config.OtelServiceName),
		semconv.ServiceVersion(version),
    ))
    if err != nil {
        return nil, err
    }

	// traces
	traceExporter, _ := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint(endpoint))
    tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(traceExporter), sdktrace.WithResource(res))
    otel.SetTracerProvider(tp)

	// logs
	logExporter, _ := otlploggrpc.New(ctx, otlploggrpc.WithInsecure(), otlploggrpc.WithEndpoint(endpoint))
    lp := sdklog.NewLoggerProvider(sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)), sdklog.WithResource(res))
    global.SetLoggerProvider(lp)

	// metrics
	metricExporter, _ := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithInsecure(), otlpmetricgrpc.WithEndpoint(endpoint))
    mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)), 
		sdkmetric.WithResource(res),
	)
    otel.SetMeterProvider(mp)

	// Cleanup-Funktion zurückgeben
	return func() {
		_ = tp.Shutdown(ctx)
		_ = lp.Shutdown(ctx)
		_ = mp.Shutdown(ctx)
	}, nil
}