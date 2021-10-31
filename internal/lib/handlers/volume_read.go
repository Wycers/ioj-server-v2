package handlers

import (
	"fmt"
	"reflect"

	"github.com/infinity-oj/server-v2/internal/app/judgements"

	"github.com/infinity-oj/server-v2/internal/lib/manager"
	"github.com/infinity-oj/server-v2/pkg/models"
	"github.com/spf13/cast"

	VS "github.com/infinity-oj/server-v2/internal/app/volumes/services"
)

type VolumeRead struct {
	jr judgements.Repository
	vs VS.Service
}

func (r *VolumeRead) IsMatched(tp string) bool {
	return tp == "volume_read"
}

func (r *VolumeRead) Work(pr *manager.ProcessRuntime) error {
	//judgement := pr.Judgement
	fmt.Println("==============>", pr.Process.Inputs[0].Value)
	f := cast.ToString(pr.Process.Inputs[0].Value)
	fmt.Println("==============>", f)
	fmt.Println(reflect.TypeOf(pr.Process.Inputs[0]))

	process := pr.Process
	process.Outputs = models.Slots{
		&models.Slot{
			Type:  "file",
			Value: f,
		},
	}
	return nil
}

func NewVolumeRead(jr judgements.Repository, vs VS.Service) *VolumeRead {
	return &VolumeRead{jr: jr, vs: vs}
}
