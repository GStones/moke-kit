package module

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/3rd/iap/pkg/iapfx"
)

// IAPModule is a fx module that provides an IAPClient
// https://github.com/awa/go-iap
var IAPModule = fx.Module("iap",
	iapfx.SettingModule,
	iapfx.ClientsModule,
)
