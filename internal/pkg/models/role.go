package models

type Role struct {
	Model

	AccountId uint64 `json:"submitterId"`
	Name      string `json:"name" gorm:"unique_index:account_name"`
}
