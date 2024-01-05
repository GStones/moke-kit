package sfx

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/server/internal/agones"
	"github.com/gstones/moke-kit/server/siface"
)

type AgonesParams struct {
	fx.In
	Agones siface.IAgones `name:"Agones" `
}

type AgonesResult struct {
	fx.Out
	Agones siface.IAgones `name:"Agones" `
}

func CreateAgones() (siface.IAgones, error) {
	a := &agones.Agones{}
	if err := a.Init(); err != nil {
		return nil, err
	}
	return a, nil
}

func CreateDefault() (siface.IAgones, error) {
	a := &agones.Default{}
	if err := a.Init(); err != nil {
		return nil, err
	}
	return a, nil
}

// AgonesModule is a fx module that provides an Agones sdk
// Requires: deploy local agones sdk server:
// https://agones.dev/site/docs/guides/client-sdks/local/
var AgonesModule = fx.Provide(
	func(l *zap.Logger) (out AgonesResult, err error) {
		a, err := CreateAgones()
		if err != nil {
			return out, err
		}
		out.Agones = a
		return
	},
)

// AgonesDefaultModule is a fx module that provides an empty Agones sdk
var AgonesDefaultModule = fx.Provide(
	func(l *zap.Logger) (out AgonesResult, err error) {
		a, err := CreateDefault()
		if err != nil {
			return out, err
		}
		out.Agones = a
		return
	})
