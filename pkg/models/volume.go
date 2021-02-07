package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type FileRecords []*FileRecord

type FileRecord struct {
	Opt      string `json:"opt"` // Add Delete Modify
	FilePath string `json:"file_path"`
	FileType string `json:"file_type"`
}

type Volume struct {
	Model

	Base        uint64       `json:"base"`
	CreatedBy   uint64       `json:"created_by"`
	Name        string       `json:"name"`
	FileRecords *FileRecords `json:"file_records" gorm:"type:json"`
}

func (fileRecord *FileRecords) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal json value:", value))
	}

	err := json.Unmarshal(bytes, fileRecord)
	return err
}

func (fileRecord *FileRecords) Value() (driver.Value, error) {
	jsonBytes, err := json.Marshal(fileRecord)
	return string(jsonBytes), err
}
