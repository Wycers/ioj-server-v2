package redis

import (
	"fmt"

	redis "github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/infinity-oj/server-v2/internal/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Options is configuration of database
type Options struct {
	Address  string `mapstructure:"addr"`
	Username string
	Password string
}

func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var err error
	o := new(Options)
	if err = v.UnmarshalKey("redis", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal redis options error")
	}

	logger.Info("load redis options success", zap.String("address", o.Address))

	return o, err
}

// New 初始化Redis
func New(o *Options) *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:     o.Address,
		Username: o.Username,
		Password: o.Password,
	})

	fmt.Println(o.Address)

	return client
}

var ProviderSet = wire.NewSet(New, NewOptions)
