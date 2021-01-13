package models

type Group struct {
	Model

	Name        string `json:"name"`
	Description string `json:"description"`
}
