package cloudfx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/utility"
)

// AWSSettingParams is the input parameter for the AWS setting module
type AWSSettingParams struct {
	fx.In

	// Region is the AWS region
	Region string `name:"awsRegion"`
	// Key is the AWS key
	Key string `name:"awsKey"`
	// Secret is the AWS secret
	// how to get the secret: https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html
	Secret string `name:"awsSecret"`
}

// AWSSettingResult is the output result for the AWS setting module
type AWSSettingResult struct {
	fx.Out

	// Region is the AWS region
	Region string `name:"awsRegion" envconfig:"AWS_REGION" default:"us-west-2"`
	// Key is the AWS key
	Key string `name:"awsKey" envconfig:"AWS_KEY" default:""`
	// Secret is the AWS secret
	Secret string `name:"awsSecret" envconfig:"AWS_SECRET" default:"" `
}

func (g *AWSSettingResult) loadFromEnv() error {
	return utility.Load(g)
}

// AWSSettingModule is the AWS setting module
var AWSSettingModule = fx.Provide(
	func() (out AWSSettingResult, err error) {
		if err = out.loadFromEnv(); err != nil {
			return
		}
		return
	},
)
