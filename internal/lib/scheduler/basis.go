package scheduler

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/PaesslerAG/gval"
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

func Evaluate(element *TaskElement) (bool, error) {
	if element.Type != "basic/evaluate" {
		return false, nil
	}
	exp, ok := element.Task.Properties["exp"]
	if !ok {
		return true, errors.New("no expression")
	}
	expStr, ok := exp.(string)
	if !ok {
		return true, errors.New("expression is not string")
	}

	var inputs []interface{}
	for _, v := range element.Task.Inputs {
		inputs = append(inputs, v.Value)
	}

	value, err := gval.Evaluate(expStr, map[string]interface{}{
		"inputs": inputs,
	})
	if err != nil {
		fmt.Println(err)
	}

	element.Task.Outputs = models.Slots{
		{
			Type:  reflect.TypeOf(value).String(),
			Value: value,
		},
	}
	return true, nil
}
