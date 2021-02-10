package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Slot struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}
type Slots []*Slot

type Task struct {
	Model

	Type        string `json:"type"`
	TaskId      string `json:"taskId"`
	JudgementId string `json:"judgementId"`

	Properties map[string]interface{} `json:"properties" gorm:"type:json"`

	Inputs  Slots `json:"inputs" gorm:"type:json"`
	Outputs Slots `json:"outputs" gorm:"type:json"`
}

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
