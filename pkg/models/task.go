package models

type SlotDescriptor struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Task struct {
	Model

	TaskId      string `json:"taskId"`
	JudgementId string `json:"judgementId"`

	Type       string `json:"type"`
	Properties string `json:"properties"`

	Inputs  string `json:"inputs"`
	Outputs string `json:"outputs"`
}
