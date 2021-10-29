package handlers

import (
	"github.com/infinity-oj/server-v2/internal/lib/manager"
	"github.com/infinity-oj/server-v2/pkg/models"
	"github.com/spf13/cast"

	"github.com/infinity-oj/server-v2/internal/app/judgements"
	VS "github.com/infinity-oj/server-v2/internal/app/volumes/services"
)

type VolumeCreate struct {
	jr judgements.Repository
	vs VS.Service
}

func (r *VolumeCreate) IsMatched(tp string) bool {
	return tp == "volume"
}

func (r *VolumeCreate) Work(pr *manager.ProcessRuntime) error {
	judgement := pr.Judgement
	v := cast.ToString(judgement.Args["volume"])
	if v == "" {
		volume, err := r.vs.CreateVolume(0)
		if err != nil {
			return err
		}
		v = volume.Name
	}
	judgement.Args["volume"] = v

	process := pr.Process
	process.Outputs = models.Slots{
		&models.Slot{
			Type:  "volume",
			Value: v,
		},
	}
	return nil
}

func NewVolumeCreate(jr judgements.Repository, vs VS.Service) *VolumeCreate {
	return &VolumeCreate{jr: jr, vs: vs}
}
