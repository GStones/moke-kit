package agonesfx

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gstones/moke-kit/3rd/agones/internal/sdk"
	"github.com/gstones/moke-kit/fxmain/pkg/mfx"
	"github.com/gstones/moke-kit/server/siface"
	"github.com/gstones/moke-kit/utility"
)

// Agones SDK module
// https://agones.dev/site/docs/guides/client-sdks/go/

type SDKParams struct {
	fx.In

	SDK siface.IAgones `name:"AgonesSDK" `
}

type SDKResult struct {
	fx.Out

	SDK siface.IAgones `name:"AgonesSDK" `
}

// CreateAgones deploy local agones sdk server
func CreateAgones() (siface.IAgones, error) {
	a := &sdk.Agones{}
	if err := a.Init(); err != nil {
		return nil, err
	}
	return a, nil
}

// CreateMock create an empty agones sdk,it will not do anything
// This is used for local/debugging
func CreateMock() (siface.IAgones, error) {
	a := &sdk.Mock{}
	if err := a.Init(); err != nil {
		return nil, err
	}
	return a, nil
}

// AgonesSDKModule is a fx module that provides an Agones sdk(deployment==prod)/mock Agones sdk(deployment!=prod)
var AgonesSDKModule = fx.Provide(
	func(
		l *zap.Logger,
		appSetting mfx.AppParams,
	) (out SDKResult, err error) {
		if deploy := utility.ParseDeployments(appSetting.Deployment); deploy.IsProd() {
			a, err := CreateAgones()
			if err != nil {
				return out, err
			}
			out.SDK = a
		} else {
			a, err := CreateMock()
			if err != nil {
				return out, err
			}
			out.SDK = a
		}
		return
	},
)
