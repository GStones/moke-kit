package sfx

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
)

// https://github.com/open-telemetry/opentelemetry-go

// OTelProviderParams OTelProviderModule provides OTel Tracer and Meter
type OTelProviderParams struct {
	fx.In

	TracerProvider *sdktrace.TracerProvider `name:"TracerProvider"`
	MetricProvider *sdkmetric.MeterProvider `name:"MetricProvider"`
}

// OTelProviderResult OTelProviderModule provides OTel Tracer and Meter
type OTelProviderResult struct {
	fx.Out

	TracerProvider *sdktrace.TracerProvider `name:"TracerProvider"`
	MetricProvider *sdkmetric.MeterProvider `name:"MetricProvider"`
}

func (otel *OTelProviderResult) Execute() (err error) {
	otel.TracerProvider, err = initTracerProvider()
	if err != nil {
		return
	}
	otel.MetricProvider, err = initMeterProvider()
	return
}

func initResource() *sdkresource.Resource {
	extraResources, _ := sdkresource.New(
		context.Background(),
		sdkresource.WithOS(),
		sdkresource.WithProcess(),
		sdkresource.WithContainer(),
		sdkresource.WithHost(),
	)
	resource, _ := sdkresource.Merge(
		sdkresource.Default(),
		extraResources,
	)
	return resource
}

func initTracerProvider() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()
	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(initResource()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}
func initMeterProvider() (*sdkmetric.MeterProvider, error) {
	ctx := context.Background()
	exporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(initResource()),
	)
	otel.SetMeterProvider(mp)
	return mp, nil
}

var OTelModule = fx.Provide(
	func() (out OTelProviderResult, err error) {
		err = out.Execute()
		return
	},
)
