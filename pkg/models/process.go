package models

type Process struct {
	Model

	Name        string `json:"name"`
	Title 		string `json:"title"`
	Description string `json:"description"`
	Definition  string `json:"definition"`
}
