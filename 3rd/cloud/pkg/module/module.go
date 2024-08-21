package module

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/3rd/cloud/pkg/cloudfx"
)

// AWSConfigModule is a module that provides the AWS config.
// how tou use?: https://github.com/aws/aws-sdk-go-v2
var AWSConfigModule = fx.Module("aws_config",
	cloudfx.AWSSettingModule,
	cloudfx.AWSConfigModule,
)
