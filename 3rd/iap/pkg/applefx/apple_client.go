package applefx

import (
	"github.com/awa/go-iap/appstore/api"
	"github.com/awa/go-iap/playstore"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AppleClientParams struct {
	fx.In

	AppleClient  *api.StoreClient  `name:"appleClient"`
	GoogleClient *playstore.Client `name:"googleClient"`
}

type AppleClientResult struct {
	fx.Out

	AppleClient  *api.StoreClient  `name:"appleClient"`
	GoogleClient *playstore.Client `name:"googleClient"`
}

func CreateAppleClient(
	sSetting AppleClientSettingParams,
) (*api.StoreClient, error) {
	c := &api.StoreConfig{
		KeyContent: sSetting.PrivateKey, // Loads a .p8 certificate
		KeyID:      sSetting.KID,        // Your private key ID from App Store Connect (Ex: 2X9R4HXF34)
		BundleID:   sSetting.BID,        // Your app’s bundle ID
		Issuer:     sSetting.Issuer,     // Your issuer ID from the Keys page in App Store Connect (Ex: "57246542-96fe-1a63-e053-0824d011072a")
		Sandbox:    sSetting.Sandbox,    // default is Production
	}
	return api.NewStoreClient(c), nil
}

func CreateGoogleClient(
	sSetting AppleClientSettingParams,
) (*playstore.Client, error) {
	client, err := playstore.New([]byte(sSetting.PublicKey))
	if err != nil {
		return nil, err
	}
	return client, nil
}

// AppleClientModule is a fx module that provides an AppleClient
// https://github.com/awa/go-iap
var AppleClientModule = fx.Provide(
	func(
		logger *zap.Logger,
		sSetting AppleClientSettingParams,
	) (out AppleClientResult, err error) {
		if aClient, err := CreateAppleClient(sSetting); err != nil {
			logger.Error("CreateAppleClient", zap.Error(err))
		} else {
			out.AppleClient = aClient
		}
		if gClient, err := CreateGoogleClient(sSetting); err != nil {
			logger.Error("CreateGoogleClient", zap.Error(err))
		} else {
			out.GoogleClient = gClient
		}
		return
	},
)
