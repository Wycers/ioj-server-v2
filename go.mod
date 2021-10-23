module github.com/infinity-oj/server-v2

go 1.16

require (
	github.com/PaesslerAG/gval v1.1.0
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/asaskevich/EventBus v0.0.0-20200907212545-49d423059eef
	github.com/gin-contrib/pprof v1.3.0
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-contrib/zap v0.0.1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.2.0
	github.com/go-redis/redis/v8 v8.11.3
	github.com/go-resty/resty/v2 v2.5.0
	github.com/google/uuid v1.2.0
	github.com/google/wire v0.5.0
	github.com/gorilla/websocket v1.4.2
	github.com/jinzhu/copier v0.2.5
	github.com/json-iterator/go v1.1.10
	github.com/opentracing-contrib/go-gin v0.0.0-20201220185307-1dd2273433a4
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.9.0
	github.com/spf13/viper v1.7.1
	github.com/swaggo/swag v1.7.0
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.4.0+incompatible
	go.uber.org/zap v1.16.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/HdrHistogram/hdrhistogram-go v1.0.1 // indirect
	github.com/spf13/cast v1.3.0
	gorm.io/driver/postgres v1.1.2
	gorm.io/gorm v1.21.16
)
