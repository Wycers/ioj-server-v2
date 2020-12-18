package api

import (
	"bytes"
	"fmt"

	"github.com/go-resty/resty/v2"

	"github.com/infinity-oj/server-v2/pkg/models"
)

type VolumeAPI interface {
	CreateVolume() (*models.Volume, error)
	CreateDirectory(volume, directory string) error
	CreateFile(volume, filename string, file []byte) error
}

type volumeAPI struct {
	client *resty.Client
}

func (a *volumeAPI) CreateDirectory(volume, dirname string) error {
	_, err := a.client.R().
		SetBody(map[string]string{
			"dirname": dirname,
		}).
		Post(fmt.Sprintf("/volume/%s/directory", volume))

	if err != nil {
		return err
	}

	return nil
}

func (a *volumeAPI) CreateFile(volume, filename string, file []byte) error {

	fmt.Println(filename)

	_, err := a.client.R().
		SetFileReader(
			"file", filename, bytes.NewReader(file)).
		Post(fmt.Sprintf("/volume/%s/file", volume))

	if err != nil {
		return err
	}

	return nil
}

func (a *volumeAPI) CreateVolume() (*models.Volume, error) {
	volume := &models.Volume{}

	_, err := a.client.R().
		SetResult(volume).
		Post("/volume")
	if err != nil {
		return nil, err
	}

	// Explore response object
	//fmt.Println("Response Info:")
	//fmt.Println("  ", resp.Request.URL)
	//fmt.Println("  Error      :", err)
	//fmt.Println("  Status Code:", resp.StatusCode())
	//fmt.Println("  Status     :", resp.Status())
	//fmt.Println("  Proto      :", resp.Proto())
	//fmt.Println("  Time       :", resp.Time())
	//fmt.Println("  Received At:", resp.ReceivedAt())
	//fmt.Println("  Body       :\n", resp)
	//fmt.Println()

	return volume, nil
}

func NewVolumeAPI(client *resty.Client) VolumeAPI {
	return &volumeAPI{
		client: client,
	}
}
