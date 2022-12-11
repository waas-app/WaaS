module github.com/waas-app/WaaS

go 1.19

require (
	github.com/allegro/bigcache/v3 v3.1.0
	github.com/coreos/go-iptables v0.6.0
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/sessions v1.2.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/improbable-eng/grpc-web v0.15.0
	github.com/miekg/dns v1.1.50
	github.com/pkg/errors v0.9.1
	github.com/place1/wg-embed v0.4.1
	github.com/spf13/cobra v1.6.1
	github.com/spf13/viper v1.7.1
	github.com/uptrace/opentelemetry-go-extra/otelgorm v0.1.17
	github.com/uptrace/opentelemetry-go-extra/otelzap v0.1.17
	github.com/urfave/negroni v1.0.0
	github.com/volatiletech/authboss/v3 v3.2.1
	go.opentelemetry.io/contrib/detectors/aws/ec2 v1.11.1
	go.opentelemetry.io/contrib/detectors/aws/ecs v1.11.1
	go.opentelemetry.io/otel v1.11.2
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.11.2
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.11.2
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.11.1
	go.opentelemetry.io/otel/sdk v1.11.2
	go.opentelemetry.io/otel/trace v1.11.2
	go.uber.org/zap v1.24.0
	golang.org/x/crypto v0.4.0
	golang.zx2c4.com/wireguard/wgctrl v0.0.0-20200609130330-bd2cb7843e1b
	google.golang.org/grpc v1.51.0
	google.golang.org/protobuf v1.28.1
	gorm.io/driver/mysql v1.4.4
	gorm.io/driver/postgres v1.4.5
	gorm.io/gorm v1.24.2
	moul.io/zapgorm2 v1.1.3
)

require (
	github.com/aws/aws-sdk-go v1.44.118 // indirect
	github.com/cenkalti/backoff/v4 v4.2.0 // indirect
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/friendsofgo/errors v0.9.2 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.13.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.12.0 // indirect
	github.com/jackc/pgx/v4 v4.17.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/klauspost/compress v1.11.7 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mdlayher/genetlink v1.0.0 // indirect
	github.com/mdlayher/netlink v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/rs/cors v1.7.0 // indirect
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.4.1 // indirect
	github.com/uptrace/opentelemetry-go-extra/otelsql v0.1.17 // indirect
	github.com/uptrace/opentelemetry-go-extra/otelutil v0.1.17 // indirect
	github.com/vishvananda/netlink v1.1.0 // indirect
	github.com/vishvananda/netns v0.0.0-20200520041808-52d707b772fe // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.11.2 // indirect
	go.opentelemetry.io/otel/metric v0.32.1 // indirect
	go.opentelemetry.io/proto/otlp v0.19.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/net v0.4.0 // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	golang.org/x/tools v0.1.12 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	golang.zx2c4.com/wireguard v0.0.20200320 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20211118181313-81c1377c94b1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	nhooyr.io/websocket v1.8.6 // indirect
)
