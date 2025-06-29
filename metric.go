package metric

import (
	"context"
	"os"
	"time"

	ometric "go.opentelemetry.io/otel/metric"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

var ok bool
var provider *metric.MeterProvider

type Config struct {
	Ok          bool
	Name        string
	Environment string
	URL         string
	DialTimeout time.Duration
}

func Setup(config *Config) (ometric.Meter, error) {
	ok = config.Ok
	if !ok {
		return nil, nil
	}
	httpOpts := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpointURL(config.URL),
		otlpmetrichttp.WithCompression(otlpmetrichttp.GzipCompression),
	}
	ctx, cancel := context.WithTimeout(context.TODO(), config.DialTimeout)
	defer cancel()
	exporter, err := otlpmetrichttp.New(ctx, httpOpts...)
	if err != nil {
		return nil, err
	}
	reader := metric.NewPeriodicReader(exporter)
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
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
		return nil, err
	}
	providerOpts := []metric.Option{
		metric.WithReader(reader),
		metric.WithResource(mergedResource),
	}
	provider = metric.NewMeterProvider(providerOpts...)
	return provider.Meter(config.Name), nil
}

func Shutdown(ctx context.Context) error {
	if !ok {
		return nil
	}
	err := provider.ForceFlush(ctx)
	if err != nil {
		return err
	}
	err = provider.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
