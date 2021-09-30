package models

type Actuator struct {
	Model

	Name  string `json:"name" gorm:"index: name"`
	Token string `json:"-"`
	Type  string `json:"type" gorm:"index: type"`

	Creator      string `json:"creator"`
	Introduction string `json:"introduction"`
}
