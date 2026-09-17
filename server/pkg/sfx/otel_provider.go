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
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/fx"
)

// https://github.com/open-telemetry/opentelemetry-go

// OTelProviderParams OTelProviderModule provides OTel Tracer and Meter
type OTelProviderParams struct {
	fx.In

	TracerProvider *sdktrace.TracerProvider `name:"TracerProvider" optional:"true"`
	MetricProvider *sdkmetric.MeterProvider `name:"MetricProvider" optional:"true"`
}

// OTelProviderResult OTelProviderModule provides OTel Tracer and Meter
type OTelProviderResult struct {
	fx.Out

	TracerProvider *sdktrace.TracerProvider `name:"TracerProvider" `
	MetricProvider *sdkmetric.MeterProvider `name:"MetricProvider"`
}

// AppIdentity is the subset of app settings used by OpenTelemetry resources.
type AppIdentity struct {
	AppName    string
	AppId      string
	Deployment string
	Version    string
}

type appIdentityParams struct {
	fx.In

	AppName    string `name:"AppName"`
	AppId      string `name:"AppId"`
	Deployment string `name:"Deployment"`
	Version    string `name:"Version"`
}

func (p appIdentityParams) identity() AppIdentity {
	return AppIdentity{
		AppName:    p.AppName,
		AppId:      p.AppId,
		Deployment: p.Deployment,
		Version:    p.Version,
	}
}

func (otel *OTelProviderResult) init(appSetting AppIdentity, enable bool) (err error) {
	if !enable {
		return
	}
	otel.TracerProvider, err = initTracerProvider(appSetting)
	if err != nil {
		return
	}
	otel.MetricProvider, err = initMeterProvider(appSetting)
	return
}

func initResource(appSetting AppIdentity) (*sdkresource.Resource, error) {
	extraResources, err := sdkresource.New(
		context.Background(),
		sdkresource.WithOS(),
		sdkresource.WithProcess(),
		sdkresource.WithContainer(),
		sdkresource.WithHost(),
		sdkresource.WithAttributes(
			semconv.ServiceName(appSetting.AppName),
			semconv.ServiceInstanceID(appSetting.AppId),
			semconv.ServiceVersion(appSetting.Version),
			semconv.ServiceNamespace(appSetting.Deployment),
		),
	)
	if err != nil {
		return sdkresource.Default(), nil
	}
	resource, err := sdkresource.Merge(
		sdkresource.Default(),
		extraResources,
	)
	if err != nil {
		return sdkresource.Default(), nil
	}
	return resource, nil
}

func initTracerProvider(appSetting AppIdentity) (*sdktrace.TracerProvider, error) {
	ctx := context.Background()
	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, err
	}
	resource, err := initResource(appSetting)
	if err != nil {
		return nil, err
	}
	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}
func initMeterProvider(appSetting AppIdentity) (*sdkmetric.MeterProvider, error) {
	ctx := context.Background()
	exporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		return nil, err
	}
	resource, err := initResource(appSetting)
	if err != nil {
		return nil, err
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(resource),
	)
	otel.SetMeterProvider(mp)
	return mp, nil
}

// CreateOTelProvider creates a OTelProvider with the given settings
func CreateOTelProvider(
	appSetting AppIdentity,
	sSetting SettingsParams,
) (out OTelProviderResult, err error) {
	err = out.init(appSetting, sSetting.OtelEnable)
	return
}

// OTelModule OTelModule provides OTel Tracer and Meter
var OTelModule = fx.Provide(
	func(
		appSetting appIdentityParams,
		sSetting SettingsParams,
	) (OTelProviderResult, error) {
		return CreateOTelProvider(appSetting.identity(), sSetting)
	},
)
