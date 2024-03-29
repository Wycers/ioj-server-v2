package handlers

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/infinity-oj/server-v2/internal/app/judgements"

	"github.com/infinity-oj/server-v2/internal/lib/manager"
	"github.com/infinity-oj/server-v2/pkg/models"
	"github.com/spf13/cast"

	VR "github.com/infinity-oj/server-v2/internal/app/volumes/repositories"
	VS "github.com/infinity-oj/server-v2/internal/app/volumes/storages"
)

type VolumeFetch struct {
	jr judgements.Repository
	vr VR.Repository
	vs VS.Storage
}

func (r *VolumeFetch) IsMatched(tp string) bool {
	return tp == "volume_fetch"
}

func (r *VolumeFetch) Work(pr *manager.ProcessRuntime) error {
	process := pr.Process
	vp := cast.ToString(process.Inputs[0].Value)
	tmp := strings.Split(vp, "/")

	volumeName := tmp[0]
	fileName := filepath.Join("/", tmp[1])

	volume, err := r.vr.GetVolume(volumeName)
	if err != nil {
		return err
	}

	file, err := r.vs.FetchFile(volume, fileName)
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return err
	}

	process.Outputs = models.Slots{
		&models.Slot{
			Type:  "bytes",
			Value: string(bytes),
		},
	}
	return nil
}

func NewVolumeFetch(jr judgements.Repository, vr VR.Repository, vs VS.Storage) *VolumeFetch {
	return &VolumeFetch{jr: jr, vr: vr, vs: vs}
}
