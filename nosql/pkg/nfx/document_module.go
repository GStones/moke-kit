package nfx

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"moke-kit/nosql/document/diface"
)

type DocumentStoreParams struct {
	fx.In
	DocumentStoreProvider diface.DocumentStoreProvider `name:"DocumentStoreProvider"`
}

type DocumentStoreResult struct {
	fx.Out
	DocumentStoreProvider diface.DocumentStoreProvider `name:"DocumentStoreProvider"`
}

func (g *DocumentStoreResult) NewDocument(
	lc fx.Lifecycle,
	l *zap.Logger,
	n SettingsParams,
) (err error) {
	//if n.DocumentStoreUrl != "" {
	//	if u, e := url.Parse(n.DocumentStoreUrl); e != nil {
	//		err = e
	//	} else {
	//		switch u.Scheme {
	//		case "couchbase":
	//			username := u.User.Username()
	//			if username == "" {
	//				username = n.NoSqlUser
	//			}
	//
	//			password, set := u.User.Password()
	//			if !set {
	//				password = n.NoSqlPassword
	//			}
	//
	//			cfg := couchbase.ClusterConfig{
	//				ConnUrl:  fmt.Sprintf("couchbase://%s", u.Host),
	//				Username: username,
	//				Password: password,
	//			}
	//			l.Info("Connect to couchbase", zap.String("url", cfg.ConnUrl))
	//			g.DocumentStoreProvider, err = couchbase.NewDocumentStoreProvider(cfg, l)
	//		case "badger":
	//			path := u.Path
	//			if runtime.GOOS == "windows" {
	//				// u.Path always has leading /, which on Windows does not play well with drive letters.
	//				if len(path) > 2 && path[0] == '/' && path[2] == ':' {
	//					path = path[1:]
	//				}
	//			}
	//			g.DocumentStoreProvider, err = badger.NewDocumentStoreProvider(path, n.GCInterval, l)
	//		case "mongodb":
	//			username := u.User.Username()
	//			if username == "" {
	//				username = n.NoSqlUser
	//			}
	//
	//			password, set := u.User.Password()
	//			if !set {
	//				password = n.NoSqlPassword
	//			}
	//
	//			cfg := mongodb.ClusterConfig{
	//				ConnUrl:  fmt.Sprintf("mongodb://%s", u.Host),
	//				Username: username,
	//				Password: password,
	//			}
	//			l.Info("Connect to mongodb", zap.String("url", cfg.ConnUrl))
	//			g.DocumentStoreProvider, err = mongodb.NewDocumentStoreProvider(cfg, l)
	//
	//		case "test":
	//			l.Info("Connect to test", zap.String("url", "test"))
	//			g.DocumentStoreProvider, err = mock.NewDocumentStoreProvider()
	//
	//		default:
	//			return errors.ErrInvalidNosqlURL
	//		}
	//	}
	//} else {
	//	return errors.ErrMissingNosqlURL
	//}

	if g.DocumentStoreProvider != nil {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return g.DocumentStoreProvider.Shutdown()
			},
		})
	}
	return
}

var DocumentStoreModule = fx.Provide(
	func(
		lc fx.Lifecycle,
		l *zap.Logger,
		n SettingsParams,
	) (dOut DocumentStoreResult, err error) {
		err = dOut.NewDocument(lc, l, n)
		return
	},
)
