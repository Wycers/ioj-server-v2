package models

type Volume struct {
	Model

	CreatedBy uint64 `json:"created_by"`
	Name      string `json:"name"`
}
