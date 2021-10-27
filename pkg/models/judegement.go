package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type Args map[string]interface{}

type Judgement struct {
	Model

	SubmissionID uint64
	BlueprintId  uint64
	Name         string
	Args         Args `gorm:"type:json"`

	Status JudgeStatus `sql:"type:judge_status" json:"status"`
	Msg    string
	Score  float64
}

func (args *Args) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal json value:", value))
	}

	err := json.Unmarshal(bytes, args)
	return err
}

func (args Args) Value() (driver.Value, error) {
	jsonBytes, err := json.Marshal(args)
	return string(jsonBytes), err
}
