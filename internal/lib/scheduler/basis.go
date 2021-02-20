package scheduler

import (
	"errors"

	"github.com/infinity-oj/server-v2/pkg/models"
)

func File(element *TaskElement) (bool, error) {
	if element.Type != "basic/file" {
		return false, nil
	}
	url, ok := element.Task.Properties["url"]
	if !ok {
		return true, errors.New("no value")
	}
	element.Task.Outputs = models.Slots{
		&models.Slot{
			Type:  "file",
			Value: url,
		},
	}
	return true, nil
}

func String(element *TaskElement) (bool, error) {
	if element.Type != "basic/string" {
		return false, nil
	}
	str, ok := element.Task.Properties["value"]
	if !ok {
		return true, errors.New("no value")
	}
	element.Task.Outputs = models.Slots{
		&models.Slot{
			Type:  "string",
			Value: str,
		},
	}
	return true, nil
}
