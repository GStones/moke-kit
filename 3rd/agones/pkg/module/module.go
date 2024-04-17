package module

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/3rd/agones/pkg/agonesfx"
)

// AgonesSDKModule is a module that provides the Agones go client SDK.
// https://agones.dev/site/docs/guides/client-sdks/
var AgonesSDKModule = fx.Module("agonesSDk",
	agonesfx.SettingsModule,
	agonesfx.AgonesSDKModule,
)

// AgonesAllocateClientModule is a module that provides Agones allocate mTls grpc client.
var AgonesAllocateClientModule = fx.Module("agonesAllocateClient",
	agonesfx.SettingsModule,
	agonesfx.AllocateClientModule,
)
