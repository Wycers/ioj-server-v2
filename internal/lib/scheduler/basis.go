package scheduler

import (
	"fmt"
)

func File(element *TaskElement) (bool, error) {
	if element.Type != "basic/file" {
		return false, nil
	}
	fmt.Println("here we go")
	fmt.Println(element.Task.Properties)
	fmt.Println(element.Task)
	return true, nil
}
