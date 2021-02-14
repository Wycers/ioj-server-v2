package scheduler

import (
	"errors"
	"fmt"

	"github.com/infinity-oj/server-v2/pkg/models"
)

func File(element *TaskElement) (bool, error) {
	if element.Type != "basic/file" {
		return false, nil
	}
	fmt.Println("here we go")
	fmt.Println(element.Task.Properties)
	fmt.Println(element.Task)
	volumeName, ok := element.Task.Properties["url"]
	if !ok {
		return true, errors.New("no value")
	}
	element.Task.Outputs = models.Slots{
		&models.Slot{
			Type:  "file",
			Value: fmt.Sprintf("%s:%s", volumeName, "/"),
		},
	}
	return true, nil
}
