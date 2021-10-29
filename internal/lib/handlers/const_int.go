package handlers

import (
	"github.com/infinity-oj/server-v2/internal/lib/manager"
	"github.com/infinity-oj/server-v2/pkg/models"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

type ConstInt struct {
}

func (c *ConstInt) IsMatched(tp string) bool {
	return tp == "const_int"
}

func (c *ConstInt) Work(pr *manager.ProcessRuntime) error {
	process := pr.Process
	str, ok := process.Properties["value"]
	if !ok {
		return errors.New("no value")
	}
	process.Outputs = models.Slots{
		&models.Slot{
			Type:  "int",
			Value: cast.ToInt(str),
		},
	}
	return nil
}

func NewConstInt() *ConstInt {
	return &ConstInt{}
}
