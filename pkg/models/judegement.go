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

	SubmissionID uint64 `json:"-"`
	BlueprintId  uint64 `json:"blueprint_id"`
	Name         string `json:"name"`
	Args         Args   `gorm:"type:json" json:"args"`

	Status JudgeStatus `sql:"type:judge_status" json:"status"`
	Msg    string      `json:"msg"`
	Score  float64     `json:"score"`
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
