package handlers

import (
	"github.com/infinity-oj/server-v2/pkg/models"
	"github.com/pkg/errors"
)

type File struct {
}

func (f File) IsMatched(tp string) bool {
	return tp == "basic/file"
}

func (f File) Work(process *models.Process) error {
	url, ok := process.Properties["url"]
	if !ok {
		return errors.New("no value")
	}
	process.Outputs = models.Slots{
		&models.Slot{
			Type:  "file",
			Value: url,
		},
	}
	return nil
}

func NewFileHandler() *File {
	return &File{}
}
