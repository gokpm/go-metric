package metric

import (
	"context"
	"os"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

type Config struct {
	Name        string
	Environment string
	URL         string
}

func Setup(ctx context.Context, config *Config) error {
	httpOpts := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpointURL(config.URL),
		otlpmetrichttp.WithCompression(otlpmetrichttp.GzipCompression),
	}
	exporter, err := otlpmetrichttp.New(ctx, httpOpts...)
	if err != nil {
		return err
	}
	reader := metric.NewPeriodicReader(exporter)
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
	providerOpts := []metric.Option{
		metric.WithReader(reader),
		metric.WithResource(mergedResource),
	}
	_ = metric.NewMeterProvider(providerOpts...).Meter(config.Name)
	return nil
}
