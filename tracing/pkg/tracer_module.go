package fxsvcapp

import (
	"context"
	"errors"
	"moke-kit/tracing/tiface"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type TracerParams struct {
	fx.In

	Tracer tiface.Tracer `name:"Tracer"`
}

type TracerResult struct {
	fx.Out

	Tracer tiface.Tracer `name:"Tracer"`
}

var (
	ErrMissingTracerServiceName = errors.New("missing tracer service name")
	ErrUnsupportedTracer        = errors.New("unsupported tracer type")
)

func (f *TracerResult) Execute(
	lc fx.Lifecycle,
	l *zap.Logger,
	t SettingsParams,
) (err error) {
	if t.TraceProvider != "" && t.TraceAgentHost != "" && t.TraceAgentPort > 0 {
		if t.TraceServiceName == "" {
			return ErrMissingTracerServiceName
		} else {
			switch t.TraceProvider {
			//case datadog.Provider:
			//	f.Tracer, err = datadog.NewTracer(
			//		l,
			//		t.TraceAgentHost,
			//		t.TraceAgentPort,
			//		t.TraceServiceName,
			//		t.TraceTags...,
			//	)
			default:
				err = ErrUnsupportedTracer
			}
		}
	} else {
		//f.Tracer = noop.NewTracer()
	}

	if f.Tracer != nil {
		lc.Append(fx.Hook{
			OnStart: func(_ context.Context) error {
				return f.Tracer.Start()
			},
			OnStop: func(_ context.Context) error {
				return f.Tracer.Stop()
			},
		})
	}

	return
}

var TracerModule = fx.Provide(
	func(
		lc fx.Lifecycle,
		l *zap.Logger,
		t SettingsParams,
	) (out TracerResult, err error) {
		err = out.Execute(lc, l, t)
		return
	},
)
