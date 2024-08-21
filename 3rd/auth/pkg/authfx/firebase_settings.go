package authfx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/utility"
)

// FirebaseSettingParams module params for injecting FirebaseSettings
type FirebaseSettingParams struct {
	fx.In

	KeyFile string `name:"FirebaseKeyFile"`
}

// FirebaseSettingsResult module result for exporting FirebaseSettings
type FirebaseSettingsResult struct {
	fx.Out

	KeyFile string `name:"FirebaseKeyFile" envconfig:"FIREBASE_KEY_FILE" default:"./configs/firebase/serviceAccountKey.json"`
}

// LoadFromEnv load from env
func (g *FirebaseSettingsResult) LoadFromEnv() error {
	return utility.Load(g)

}

// FirebaseSettingsModule is the Firebase settings module
var FirebaseSettingsModule = fx.Provide(
	func() (out FirebaseSettingsResult, err error) {
		err = out.LoadFromEnv()
		return
	},
)
