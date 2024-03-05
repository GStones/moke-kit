package iapfx

import (
	"os"

	"github.com/awa/go-iap/appstore/api"
	"github.com/awa/go-iap/playstore"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ClientsParams is a struct that holds the parameters for the IAP clients
type ClientsParams struct {
	fx.In

	AppleClient  *api.StoreClient  `name:"appleClient"`
	GoogleClient *playstore.Client `name:"googleClient"`
}

// ClientsResult is a struct that holds the results for the IAP clients
type ClientsResult struct {
	fx.Out

	AppleClient  *api.StoreClient  `name:"appleClient"`
	GoogleClient *playstore.Client `name:"googleClient"`
}

// CreateAppleClient creates a new Apple client
func CreateAppleClient(
	sSetting SettingParams,
) (*api.StoreClient, error) {
	var privateKey []byte
	if data, err := os.ReadFile(sSetting.PrivateKeyPath); err != nil {
		return nil, err
	} else {
		privateKey = data
	}

	c := &api.StoreConfig{
		KeyContent: privateKey,       // Loads a .p8 certificate
		KeyID:      sSetting.KID,     // Your private key ID from App Store Connect (Ex: 2X9R4HXF34)
		BundleID:   sSetting.BID,     // Your app’s bundle ID
		Issuer:     sSetting.Issuer,  // Your issuer ID from the Keys page in App Store Connect (Ex: "57246542-96fe-1a63-e053-0824d011072a")
		Sandbox:    sSetting.Sandbox, // default is Production
	}
	return api.NewStoreClient(c), nil
}

// CreateGoogleClient creates a new Google client
func CreateGoogleClient(
	sSetting SettingParams,
) (*playstore.Client, error) {
	var publicKey []byte
	if data, err := os.ReadFile(sSetting.PublicKeyPath); err != nil {
		return nil, err
	} else {
		publicKey = data
	}

	client, err := playstore.New(publicKey)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// ClientsModule is a fx module that provides an IAPClient
var ClientsModule = fx.Provide(
	func(
		logger *zap.Logger,
		sSetting SettingParams,
	) (out ClientsResult, err error) {
		if sSetting.PrivateKeyPath == "" {

		} else if aClient, err := CreateAppleClient(sSetting); err != nil {
			logger.Error("Create apple client", zap.Error(err))
		} else {
			logger.Info("Create apple client success", zap.Bool("sandbox", sSetting.Sandbox))
			out.AppleClient = aClient
		}
		if sSetting.PublicKeyPath == "" {

		} else if gClient, err := CreateGoogleClient(sSetting); err != nil {
			logger.Error("create google client", zap.Error(err))
		} else {
			logger.Info("create google client success")
			out.GoogleClient = gClient
		}
		return
	},
)
