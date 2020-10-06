package models

type Credential struct {
	Model

	AccountId int64  `gorm:"not null"`
	Username  string `gorm:"not null; unique_index:idx1"`
	Hash      string `gorm:"not null;"`
	Salt      string `gorm:"not null;"`
}
