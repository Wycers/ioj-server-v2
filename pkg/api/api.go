package api

import (
	"fmt"

	cookiejar "github.com/juju/persistent-cookiejar"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

// Options is log configuration struct
type Options struct {
	Url string `yaml:"url"`
}

var Jar, _ = cookiejar.New(nil)

func NewOptions(v *viper.Viper) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)

	o.Url = fmt.Sprintf("%s/api/v1", v.Get("host").(string))

	fmt.Printf("Host: %s\n", o.Url)

	return o, err
}

func NewClient(options *Options) *resty.Client {
	client := resty.New()

	client.SetHostURL(options.Url)
	client.SetCookieJar(Jar)
	return client
}
