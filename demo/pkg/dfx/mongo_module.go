package dfx

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"moke-kit/gorm/pkg/nfx"
)

type DemoDBParams struct {
	fx.In
	DemoDatabase *mongo.Database `name:"DemoDatabase"`
}

type DemoDBResult struct {
	fx.Out
	DemoDatabase *mongo.Database `name:"DemoDatabase"`
}

func (g *DemoDBResult) Execute(
	l *zap.Logger,
	dbName string,
	mClient *mongo.Client,
) (err error) {
	g.DemoDatabase = mClient.Database(dbName)
	l.Info("OpenDbDriver", zap.String("DbName", dbName))
	return
}

var DemoDBModule = fx.Provide(
	func(
		l *zap.Logger,
		s SettingsParams,
		dsp nfx.MongoParams,
	) (out DemoDBResult, err error) {
		if err := out.Execute(l, s.DbName, dsp.MongoClient); err != nil {
			return out, err
		}
		return out, nil
	},
)
