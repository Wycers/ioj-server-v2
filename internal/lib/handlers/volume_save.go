package handlers

import (
	"strings"

	"github.com/infinity-oj/server-v2/internal/lib/manager"
	"github.com/pkg/errors"
	"github.com/spf13/cast"

	"github.com/infinity-oj/server-v2/internal/app/judgements"
	VS "github.com/infinity-oj/server-v2/internal/app/volumes/services"
)

type VolumeSave struct {
	jr judgements.Repository
	vs VS.Service
}

func (r *VolumeSave) IsMatched(tp string) bool {
	return tp == "volume_save"
}

func (r *VolumeSave) Work(pr *manager.ProcessRuntime) error {
	pr.Mutex.Lock()
	defer pr.Mutex.Unlock()

	judgement := pr.Judgement
	v := cast.ToString(pr.Process.Inputs[0].Value)
	if v == "" {
		return errors.New("missing volume")
	}
	judgement.Args["volume"] = v

	process := pr.Process
	vp := cast.ToString(process.Inputs[1].Value)
	tmp := strings.Split(vp, "/")

	filename := cast.ToString(process.Properties["filename"])

	nv, err := r.vs.CopyFile(tmp[0], "/", tmp[1], v, "/", filename)
	if err != nil {
		return err
	}
	judgement.Args["volume"] = nv.Name
	return nil
}

func NewVolumeSave(jr judgements.Repository, vs VS.Service) *VolumeSave {
	return &VolumeSave{jr: jr, vs: vs}
}
