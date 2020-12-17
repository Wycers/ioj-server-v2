package models

type Judgement struct {
	Model

	SubmissionId uint64
	ProcessId    uint64

	JudgementId string
	Status      JudgeStatus `sql:"type:judge_status"`
	Score       float64
}
