package models

type Actuator struct {
	Model

	Name         string `json:"name" gorm:"index: name"`
	Introduction string `json:"introduction"`
	Token        string `json:"-"`
	Creator      string `json:"creator"`
}
