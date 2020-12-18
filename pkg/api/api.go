package api

import (
	"fmt"
	"net/http"

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

type API interface {
	SetHostUrl(hostUrl string)
	SetCookieJar(Jar http.CookieJar)

	NewAccountAPI() AccountAPI
	NewVolumeAPI() VolumeAPI
	NewJudgementAPI() JudgementAPI
	NewSubmissionAPI() SubmissionAPI
}

type api struct {
	client *resty.Client
}

func (a api) SetHostUrl(hostUrl string) {
	a.client.SetHostURL(hostUrl)
}

func (a api) SetCookieJar(Jar http.CookieJar) {
	a.client.SetCookieJar(Jar)
}

func (a api) NewAccountAPI() AccountAPI {
	return NewAccountAPI(a.client)
}

func (a api) NewVolumeAPI() VolumeAPI {
	return NewVolumeAPI(a.client)
}

func (a api) NewJudgementAPI() JudgementAPI {
	return NewJudgementAPI(a.client)
}

func (a api) NewSubmissionAPI() SubmissionAPI {
	return NewSubmissionAPI(a.client)
}

func New() API {
	client := resty.New()
	return &api{
		client: client,
	}
}
