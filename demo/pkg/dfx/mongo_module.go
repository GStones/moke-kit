package dfx

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"moke-kit/nosql/document/diface"
	"moke-kit/nosql/pkg/nfx"
)

type DemoDBParams struct {
	fx.In
	GameServerStore diface.ICollection `name:"GameServerStore"`
}

type DemoDBResult struct {
	fx.Out
	GameServerStore diface.ICollection `name:"GameServerStore"`
}

func (g *DemoDBResult) Execute(
	l *zap.Logger,
	dbName string,
	dp diface.IDocumentDb,
) (err error) {
	g.GameServerStore, err = dp.OpenDbDriver(dbName)
	l.Info("OpenDbDriver", zap.String("DbName", dbName))
	return
}

var DemoDBModule = fx.Provide(
	func(
		l *zap.Logger,
		s SettingsParams,
		dsp nfx.DocumentStoreParams,
	) (out DemoDBResult, err error) {
		if err := out.Execute(l, s.DbName, dsp.DriverProvider); err != nil {
			return out, err
		}
		return out, nil
	},
)
