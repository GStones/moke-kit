package applefx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/utility"
)

type AppleClientSettingParams struct {
	fx.In

	// ---------apple client setting --------
	// Key ID for the given private key, described in App Store Connect
	KID string `name:"appleKeyId"`
	// Issuer ID for the App Store Connect team
	Issuer string `name:"appleIssuer"`
	// The bytes of the PKCS#8 private key created on App Store Connect. Keep this key safe as you can only download it once.
	PrivateKey []byte `name:"applePrivateKey"`
	//  Application's bundle ID, e.g. com.example.testbundleid2021
	BID string `name:"appleBundleId"`
	// A boolean value that indicates whether the token is for the App Store Connect API sandbox environment
	Sandbox bool `name:"appleSandbox" `

	// ---------google store client setting --------
	// You need to prepare a public key for your Android app's in app billing
	// at https://console.developers.google.com.
	PublicKey string `name:"googlePlayPublicKey"`
}

type AppleClientSettingResult struct {
	fx.Out
	// ---------apple client setting --------
	KID        string `name:"appleKeyId"  envconfig:"APPLE_KEY_ID" default:""`
	Issuer     string `name:"appleIssuer" envconfig:"APPLE_ISSUER" default:""`
	PrivateKey []byte `name:"applePrivateKey" envconfig:"APPLE_PRIVATE_KEY" default:""`
	BID        string `name:"appleBundleId" envconfig:"APPLE_BUNDLE_ID" default:""`
	Sandbox    bool   `name:"appleSandbox" envconfig:"APPLE_SANDBOX" default:"true"`

	// ---------google store client setting --------
	// You need to prepare a public key for your Android app's in app billing
	// at https://console.developers.google.com.
	PublicKey string `name:"googlePlayPublicKey" envconfig:"GOOGLE_PLAY_PUBLIC_KEY" default:""`
}

func (g *AppleClientSettingResult) LoadFromEnv() (err error) {
	err = utility.Load(g)
	return
}

var AppleClientSettingModule = fx.Provide(
	func() (out AppleClientSettingResult, err error) {
		err = out.LoadFromEnv()
		return
	})
