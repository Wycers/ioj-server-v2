package redis

import (
	redis "github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/infinity-oj/server-v2/internal/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Options is  configuration of database
type Options struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Debug    bool
}

func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var err error
	o := new(Options)
	if err = v.UnmarshalKey("db", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal db option error")
	}

	logger.Info("load database options success", zap.String("url", o.URL))

	return o, err
}

// New 初始化Redis
func New(o *Options) *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:     o.Address,
		Username: o.Username,
		Password: o.Password,
	})

	return client
}

var ProviderSet = wire.NewSet(New, NewOptions)
