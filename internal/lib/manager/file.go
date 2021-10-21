package manager

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/PaesslerAG/gval"
	"github.com/infinity-oj/server-v2/pkg/models"
)

func FileHandler(element *ProcessElement) (bool, error) {
	if element.Process.Type != "basic/file" {
		return false, nil
	}
	url, ok := element.Process.Properties["url"]
	if !ok {
		return true, errors.New("no value")
	}
	element.Process.Outputs = models.Slots{
		&models.Slot{
			Type:  "file",
			Value: url,
		},
	}
	return true, nil
}

func String(element *ProcessElement) (bool, error) {
	if element.Process.Type != "basic/string" {
		return false, nil
	}
	str, ok := element.Process.Properties["value"]
	if !ok {
		return true, errors.New("no value")
	}
	element.Process.Outputs = models.Slots{
		&models.Slot{
			Type:  "string",
			Value: str,
		},
	}
	return true, nil
}
func ConstString(element *ProcessElement) (bool, error) {
	if element.Process.Type != "const_string" {
		return false, nil
	}
	str, ok := element.Process.Properties["value"]
	fmt.Println(element.Process.Properties)
	if !ok {
		return true, errors.New("no value")
	}
	element.Process.Outputs = models.Slots{
		&models.Slot{
			Type:  "string",
			Value: str,
		},
	}
	return true, nil
}

func Evaluate(element *ProcessElement) (bool, error) {
	if element.Process.Type != "basic/evaluate" {
		return false, nil
	}
	exp, ok := element.Process.Properties["exp"]
	if !ok {
		return true, errors.New("no expression")
	}
	expStr, ok := exp.(string)
	if !ok {
		return true, errors.New("expression is not string")
	}

	var inputs []interface{}
	for _, v := range element.Process.Inputs {
		inputs = append(inputs, v.Value)
	}

	value, err := gval.Evaluate(expStr, map[string]interface{}{
		"inputs": inputs,
	})
	if err != nil {
		fmt.Println(err)
	}

	element.Process.Outputs = models.Slots{
		{
			Type:  reflect.TypeOf(value).String(),
			Value: value,
		},
	}
	return true, nil
}
