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

	"github.com/gstones/moke-kit/fxmain/pkg/mfx"
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

func (otel *OTelProviderResult) init(appSetting mfx.AppParams, enable bool) (err error) {
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

func initResource(appSetting mfx.AppParams) *sdkresource.Resource {
	extraResources, _ := sdkresource.New(
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
	resource, _ := sdkresource.Merge(
		sdkresource.Default(),
		extraResources,
	)
	return resource
}

func initTracerProvider(appSetting mfx.AppParams) (*sdktrace.TracerProvider, error) {
	ctx := context.Background()
	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, err
	}
	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(initResource(appSetting)),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}
func initMeterProvider(appSetting mfx.AppParams) (*sdkmetric.MeterProvider, error) {
	ctx := context.Background()
	exporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(initResource(appSetting)),
	)
	otel.SetMeterProvider(mp)
	return mp, nil
}

// CreateOTelProvider creates a OTelProvider with the given settings
func CreateOTelProvider(
	appSetting mfx.AppParams,
	sSetting SettingsParams,
) (out OTelProviderResult, err error) {
	err = out.init(appSetting, sSetting.OtelEnable)
	return
}

// OTelModule OTelModule provides OTel Tracer and Meter
var OTelModule = fx.Provide(
	func(
		appSetting mfx.AppParams,
		sSetting SettingsParams,
	) (OTelProviderResult, error) {
		return CreateOTelProvider(appSetting, sSetting)
	},
)
