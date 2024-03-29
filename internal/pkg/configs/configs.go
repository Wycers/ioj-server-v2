package configs

import (
	"fmt"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

// New 初始化viper
func New(path string) (*viper.Viper, error) {
	var (
		err error
		v   = viper.New()
	)
	v.AddConfigPath(".")
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err == nil {
		fmt.Printf("use config file -> %s\n", v.ConfigFileUsed())
	} else {
		return nil, err
	}

	return v, err
}

var ProviderSet = wire.NewSet(New)
