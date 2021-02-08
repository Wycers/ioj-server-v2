package api

import (
	"bytes"
	"fmt"

	"github.com/go-resty/resty/v2"

	"github.com/infinity-oj/server-v2/pkg/models"
)

type VolumeAPI interface {
	CreateVolume() (*models.Volume, error)
	CreateDirectory(volume, directory string) (*models.Volume, error)
	CreateFile(volume, filename string, file []byte) (*models.Volume, error)
}

type volumeAPI struct {
	client *resty.Client
}

func (a *volumeAPI) CreateDirectory(volumeName, dirname string) (*models.Volume, error) {
	volume := &models.Volume{}

	_, err := a.client.R().
		SetBody(map[string]string{
			"dirname": dirname,
		}).
		SetResult(volume).
		Post(fmt.Sprintf("/volume/%s/directory", volumeName))

	if err != nil {
		return nil, err
	}

	return volume, nil
}

func (a *volumeAPI) CreateFile(volumeName, filename string, file []byte) (*models.Volume, error) {
	volume := &models.Volume{}

	_, err := a.client.R().
		SetFileReader(
			"file", filename, bytes.NewReader(file)).
		SetResult(volume).
		Post(fmt.Sprintf("/volume/%s/file", volumeName))

	if err != nil {
		return nil, err
	}

	return volume, nil
}

func (a *volumeAPI) CreateVolume() (*models.Volume, error) {
	volume := &models.Volume{}

	_, err := a.client.R().
		SetResult(volume).
		Post("/volume")
	if err != nil {
		return nil, err
	}

	return volume, nil
}

func NewVolumeAPI(client *resty.Client) VolumeAPI {
	return &volumeAPI{
		client: client,
	}
}
