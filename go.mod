module github.com/liluoliluoli/gnboot

go 1.23

toolchain go1.23.3

require (
	github.com/go-cinch/common/constant v1.0.5
	github.com/go-cinch/common/i18n v1.0.6
	github.com/go-cinch/common/id v1.0.4
	github.com/go-cinch/common/log v1.2.0
	github.com/go-cinch/common/middleware/i18n v1.0.5
	github.com/go-cinch/common/middleware/logging v1.0.0
	github.com/go-cinch/common/middleware/tenant v1.0.1
	github.com/go-cinch/common/middleware/trace v1.0.2
	github.com/go-cinch/common/plugins/gorm/filter v1.0.1
	github.com/go-cinch/common/plugins/gorm/log v1.0.4
	github.com/go-cinch/common/plugins/gorm/tenant v1.0.1
	github.com/go-cinch/common/plugins/k8s/pod v1.0.0
	github.com/go-cinch/common/plugins/kratos/config/env v1.0.0
	github.com/go-cinch/common/plugins/kratos/encoding/yml v1.0.1
	github.com/go-cinch/common/utils v1.0.4
	github.com/go-cinch/common/worker v1.0.7
	github.com/go-kratos/kratos/v2 v2.8.4
	github.com/google/uuid v1.6.0
	github.com/google/wire v0.6.0
	github.com/liluoliluoli/gnboot/api v0.0.0
	github.com/pkg/errors v0.9.1
	github.com/redis/go-redis/v9 v9.2.1
	github.com/samber/lo v1.47.0
	go.opentelemetry.io/otel v1.34.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.30.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.30.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.30.0
	go.opentelemetry.io/otel/sdk v1.30.0
	go.opentelemetry.io/otel/trace v1.34.0
	golang.org/x/text v0.21.0
	google.golang.org/grpc v1.67.3
	google.golang.org/protobuf v1.36.1
	gorm.io/gen v0.3.22
	gorm.io/gorm v1.25.2
	gorm.io/plugin/dbresolver v1.4.1
)

require (
	github.com/biter777/countries v1.7.5
	github.com/gojek/heimdall/v7 v7.0.3
	github.com/golang-jwt/jwt/v5 v5.1.0
)

require (
	cel.dev/expr v0.16.1 // indirect
	dario.cat/mergo v1.0.0 // indirect
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cncf/xds/go v0.0.0-20240905190251-b4127c9b8d78 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/envoyproxy/go-control-plane v0.13.1 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-cinch/common/migrate v1.0.4 // indirect
	github.com/go-cinch/common/nx v1.0.4 // indirect
	github.com/go-gorp/gorp/v3 v3.1.0 // indirect
	github.com/go-kratos/aegis v0.2.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-playground/form/v4 v4.2.1 // indirect
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/gojek/valkyrie v0.0.0-20180215180059-6aee720afcdf // indirect
	github.com/golang-module/carbon/v2 v2.3.12 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/gnostic v0.7.0 // indirect
	github.com/google/gnostic-models v0.6.9-0.20230804172637-c7be7c783f49 // indirect
	github.com/gorhill/cronexpr v0.0.0-20180427100037-88b0669f7d75 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.22.0 // indirect
	github.com/hibiken/asynq v0.24.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/lufia/plan9stats v0.0.0-20230326075908-cb1d2100619a // indirect
	github.com/nicksnyder/go-i18n/v2 v2.2.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20221212215047-62379fc7944b // indirect
	github.com/r3labs/diff/v3 v3.0.1 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/rubenv/sql-migrate v1.5.1 // indirect
	github.com/shirou/gopsutil/v3 v3.23.6 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/sony/sonyflake v1.1.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/tklauser/go-sysconf v0.3.11 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.5 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.3 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/metric v1.34.0 // indirect
	go.opentelemetry.io/proto/otlp v1.3.1 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/time v0.8.0 // indirect
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241209162323-e6fa225c2576 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241223144023-3abc09e42ca8 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/datatypes v1.2.0 // indirect
	gorm.io/driver/mysql v1.5.1 // indirect
	gorm.io/hints v1.1.2 // indirect
)

replace github.com/liluoliluoli/gnboot/api => ./api
