package metric

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

type Config struct {
	Name        string
	Environment string
	URL         string
	Timeout     time.Duration
	Interval    time.Duration
}

func Init(ctx context.Context, config *Config) error {
	httpOpts := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpointURL(config.URL),
		otlpmetrichttp.WithTimeout(config.Timeout),
		otlpmetrichttp.WithCompression(otlpmetrichttp.GzipCompression),
	}
	exporter, err := otlpmetrichttp.New(ctx, httpOpts...)
	if err != nil {
		return err
	}
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	base := resource.Default()
	newResource := resource.NewWithAttributes(
		base.SchemaURL(),
		semconv.ServiceName(config.Name),
		semconv.DeploymentEnvironmentName(config.Environment),
		semconv.HostName(hostname),
	)
	mergedResource, err := resource.Merge(base, newResource)
	if err != nil {
		return err
	}
	readerOpts := []metric.PeriodicReaderOption{
		metric.WithTimeout(config.Timeout),
		metric.WithInterval(config.Interval),
	}
	reader := metric.NewPeriodicReader(exporter, readerOpts...)
	providerOpts := []metric.Option{
		metric.WithResource(mergedResource),
		metric.WithReader(reader),
	}
	provider := metric.NewMeterProvider(providerOpts...)
	otel.SetMeterProvider(provider)
	return nil
}
