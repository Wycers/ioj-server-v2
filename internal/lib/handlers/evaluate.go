package handlers

import (
	"reflect"

	"github.com/PaesslerAG/gval"
	"github.com/infinity-oj/server-v2/pkg/models"
	"github.com/pkg/errors"
)

type Evaluate struct {
}

func (e *Evaluate) IsMatched(tp string) bool {
	return tp == "basic/evaluate"
}

func (e *Evaluate) Work(process *models.Process) error {
	exp, ok := process.Properties["exp"]
	if !ok {
		return errors.New("no expression")
	}
	expStr, ok := exp.(string)
	if !ok {
		return errors.New("expression is not string")
	}

	var inputs []interface{}
	for _, v := range process.Inputs {
		inputs = append(inputs, v.Value)
	}

	value, err := gval.Evaluate(expStr, map[string]interface{}{
		"inputs": inputs,
	})
	if err != nil {
		return err
	}
	process.Outputs = models.Slots{
		{
			Type:  reflect.TypeOf(value).String(),
			Value: value,
		},
	}
	return nil
}

func NewEvaluateHandler() *Evaluate {
	return &Evaluate{}
}
