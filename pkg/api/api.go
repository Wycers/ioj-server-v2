package api

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/go-resty/resty/v2"
)

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

	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		if resp.StatusCode() == 401 {
			return errors.New("You need to login")
		}
		return nil
	})

	return &api{
		client: client,
	}
}
