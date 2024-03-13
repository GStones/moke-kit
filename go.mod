module github.com/gstones/moke-kit

go 1.21

require (
	agones.dev/agones v1.38.0
	github.com/abiosoft/ishell v2.0.0+incompatible
	github.com/awa/go-iap v1.26.5
	github.com/duke-git/lancet/v2 v2.2.9
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.0.1
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1
	github.com/gstones/zinx v0.1.2
	github.com/hashicorp/vault/api v1.11.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/nats-io/nats.go v1.32.0
	github.com/pkg/errors v0.9.1
	github.com/redis/go-redis/v9 v9.4.0
	github.com/soheilhy/cmux v0.1.5
	github.com/spf13/cobra v1.8.0
	go.mongodb.org/mongo-driver v1.13.1
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.47.0
	go.opentelemetry.io/otel v1.23.1
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.23.1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.23.1
	go.opentelemetry.io/otel/sdk v1.23.1
	go.opentelemetry.io/otel/sdk/metric v1.23.1
	go.opentelemetry.io/otel/trace v1.23.1
	go.uber.org/atomic v1.11.0
	go.uber.org/automaxprocs v1.5.3
	go.uber.org/fx v1.20.1
	go.uber.org/ratelimit v0.3.0
	go.uber.org/zap v1.26.0
	google.golang.org/genproto/googleapis/api v0.0.0-20240125205218-1f4bbc51befe
	google.golang.org/grpc v1.61.0
	google.golang.org/protobuf v1.32.0
	gopkg.in/fsnotify.v1 v1.4.7
	gorm.io/driver/sqlite v1.5.4
	gorm.io/gorm v1.25.6
)

require (
	cloud.google.com/go/compute v1.23.4 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	github.com/abiosoft/readline v0.0.0-20180607040430-155bce2042db // indirect
	github.com/benbjohnson/clock v1.3.5 // indirect
	github.com/cenkalti/backoff/v3 v3.2.2 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/chzyer/logex v1.1.10 // indirect
	github.com/chzyer/test v0.0.0-20180213035817-a1ea475d72b1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/flynn-archive/go-shlex v0.0.0-20150515145356-3f9db97f8568 // indirect
	github.com/go-jose/go-jose/v3 v3.0.3 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-test/deep v1.1.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.5.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.5 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-secure-stdlib/parseutil v0.1.8 // indirect
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2 // indirect
	github.com/hashicorp/go-sockaddr v1.0.6 // indirect
	github.com/hashicorp/hcl v1.0.1-vault-5 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.17.5 // indirect
	github.com/klauspost/cpuid/v2 v2.2.6 // indirect
	github.com/klauspost/reedsolomon v1.12.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/templexxx/cpufeat v0.0.0-20180724012125-cef66df7f161 // indirect
	github.com/templexxx/xor v0.0.0-20191217153810-f85b25db303b // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/xtaci/kcp-go v5.4.20+incompatible // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.47.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.23.1 // indirect
	go.opentelemetry.io/otel/metric v1.23.1 // indirect
	go.opentelemetry.io/proto/otlp v1.1.0 // indirect
	go.uber.org/dig v1.17.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.19.0 // indirect
	golang.org/x/exp v0.0.0-20240119083558-1b970713d09a // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/oauth2 v0.16.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/api v0.161.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto v0.0.0-20240125205218-1f4bbc51befe // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240125205218-1f4bbc51befe // indirect
)
