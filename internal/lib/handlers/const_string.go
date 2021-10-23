package handlers

import (
	"github.com/infinity-oj/server-v2/pkg/models"
	"github.com/pkg/errors"
)

type ConstString struct {
}

func (c *ConstString) IsMatched(tp string) bool {
	return tp == "const_string"
}

func (c *ConstString) Work(process *models.Process) error {
	str, ok := process.Properties["value"]
	if !ok {
		return errors.New("no value")
	}
	process.Outputs = models.Slots{
		&models.Slot{
			Type:  "string",
			Value: str,
		},
	}
	return nil
}

func NewConstString() *ConstString {
	return &ConstString{}
}
