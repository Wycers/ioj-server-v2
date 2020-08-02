package models

type Problem struct {
	Model

	Name      string `json:"name" gorm:"unique_index:idx2"`
	Title     string `json:"title"`
	ProcessId uint64
}
