package module

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/3rd/agones/pkg/agonesfx"
)

var AgonesSDKModule = fx.Module("agonesSDk",
	agonesfx.AllocateSettingsModule,
	agonesfx.AgonesSDKModule,
)

var AgonesAllocateClientModule = fx.Module("agonesAllocateClient",
	agonesfx.AllocateSettingsModule,
	agonesfx.AllocateClientModule,
)
