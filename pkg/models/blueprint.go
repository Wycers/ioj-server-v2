package models

type Blueprint struct {
	Model

	Name  string `json:"name" gorm:"index: name"`
	Title string `json:"title"`

	Definition string `json:"definition"`
}
