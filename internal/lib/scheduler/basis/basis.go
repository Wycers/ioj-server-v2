package basis

import (
	"fmt"

	"github.com/infinity-oj/server-v2/internal/lib/scheduler"
)

func File(element *scheduler.TaskElement) (bool, error) {
	if element.Type != "basic/file" {
		return false, nil
	}
	fmt.Println("here we go")
	fmt.Println(element)
	return true, nil
}
