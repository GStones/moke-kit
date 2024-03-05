package iapfx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/utility"
)

// SettingParams is a struct that holds the parameters for the IAP clients
type SettingParams struct {
	fx.In

	// ---------apple client setting --------
	// Key ID for the given private key, described in App Store Connect
	KID string `name:"appleKeyId"`
	// Issuer ID for the App Store Connect team
	Issuer string `name:"appleIssuer"`
	// The bytes of the PKCS#8 private key created on App Store Connect. Keep this key safe as you can only download it once.
	PrivateKeyPath string `name:"applePrivateKey"`
	//  Application's bundle ID, e.g. com.example.testbundleid2021
	BID string `name:"appleBundleId"`
	// A boolean value that indicates whether the token is for the App Store Connect API sandbox environment
	Sandbox bool `name:"appleSandbox" `

	// ---------google store client setting --------
	// You need to prepare a public key for your Android app's in app billing
	// at https://console.developers.google.com.
	PublicKeyPath string `name:"googlePlayPublicKey"`
}

// SettingResult is a struct that holds the results for the IAP clients
type SettingResult struct {
	fx.Out
	// ---------apple client setting --------
	KID            string `name:"appleKeyId"  envconfig:"APPLE_KEY_ID" default:""`
	Issuer         string `name:"appleIssuer" envconfig:"APPLE_ISSUER" default:""`
	PrivateKeyPath string `name:"applePrivateKey" envconfig:"APPLE_PRIVATE_KEY" default:"./configs/iap/apple_key.p8"`
	BID            string `name:"appleBundleId" envconfig:"APPLE_BUNDLE_ID" default:""`
	Sandbox        bool   `name:"appleSandbox" envconfig:"APPLE_SANDBOX" default:"true"`

	// ---------google store client setting --------
	// You need to prepare a public key for your Android app's in app billing
	// at https://console.developers.google.com.
	PublicKeyPath string `name:"googlePlayPublicKey" envconfig:"GOOGLE_PLAY_PUBLIC_KEY" default:"./configs/iap/google_key.json"`
}

func (g *SettingResult) LoadFromEnv() (err error) {
	err = utility.Load(g)
	return
}

// SettingModule is a fx setting module that provides an IAPClient
var SettingModule = fx.Provide(
	func() (out SettingResult, err error) {
		err = out.LoadFromEnv()
		return
	})
