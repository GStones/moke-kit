package cloudfx

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/utility"
)

// AWSConfigParams module params for injecting AWSConfig
// https://github.com/aws/aws-sdk-go-v2
type AWSConfigParams struct {
	fx.In
	// Config aws config object
	Config aws.Config `name:"awsConfig"`
}

type AWSConfigResult struct {
	fx.Out
	// Config aws config object
	Config aws.Config `name:"awsConfig"`
}

func (g *AWSConfigResult) loadFromEnv() error {
	return utility.Load(g)
}

// AWSConfigModule is the AWS config module
// how tou use?: https://github.com/aws/aws-sdk-go-v2
var AWSConfigModule = fx.Provide(
	func(
		setting AWSSettingParams,
	) (out AWSConfigResult, err error) {
		out.Config, err = config.LoadDefaultConfig(context.Background(),
			config.WithRegion(setting.Region),
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(setting.Key, setting.Secret, ""),
			),
		)
		return
	},
)
