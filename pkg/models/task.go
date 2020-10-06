package models

type Task struct {
	Model

	JudgementId string `json:"judgementId"`

	TaskId string `json:"taskId"`

	Type       string `json:"type"`
	Properties string `json:"properties"`
	Inputs     string `json:"inputs"`
	Outputs    string `json:"outputs"`
}
