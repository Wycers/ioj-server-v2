package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Slot struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Task struct {
	Model

	TaskId      string `json:"taskId"`
	JudgementId string `json:"judgementId"`

	Type       string `json:"type"`
	Properties string `json:"properties"`

	Inputs  Slots `json:"inputs" gorm:"type:json"`
	Outputs Slots `json:"outputs" gorm:"type:json"`
}

type Slots []*Slot

func (slots *Slots) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal json value:", value))
	}

	err := json.Unmarshal(bytes, slots)
	return err
}

func (slots Slots) Value() (driver.Value, error) {
	jsonBytes, err := json.Marshal(slots)
	return string(jsonBytes), err
}
