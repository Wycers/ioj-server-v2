package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileRecords []*FileRecord

type FileRecord struct {
	Opt      string `json:"opt"` // Add Delete
	FileType string `json:"file_type"`
	FilePath string `json:"file_path"`

	VolumeName string `json:"volume"`
	VolumePath string `json:"volumePath"`
}

type Volume struct {
	Model

	Base        uint64      `json:"base"`
	CreatedBy   uint64      `json:"created_by"`
	Name        string      `json:"name"`
	FileRecords FileRecords `json:"file_records" gorm:"type:json"`
}

func (fileRecords *FileRecords) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal json value:", value))
	}

	err := json.Unmarshal(bytes, fileRecords)
	return err
}

func (fileRecords FileRecords) Value() (driver.Value, error) {
	jsonBytes, err := json.Marshal(fileRecords)
	return string(jsonBytes), err
}

func (f FileRecord) Name() string {
	return strings.TrimPrefix(f.FilePath, string(filepath.Separator))
}

func (f FileRecord) Size() int64 {
	return 0
}

func (f FileRecord) Mode() os.FileMode {
	return 0644
}

func (f FileRecord) ModTime() time.Time {
	return time.Now()
}

func (f FileRecord) IsDir() bool {
	return f.FileType == "d"
}

func (f FileRecord) Sys() interface{} {
	return nil
}
