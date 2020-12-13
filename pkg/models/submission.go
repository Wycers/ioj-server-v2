package models

type Submission struct {
	Model

	SubmissionId string `json:"submissionId"`

	SubmitterId uint64 `json:"submitterId"`

	ProblemId uint64 `json:"problemId"`

	UserVolume string `json:"userVolume"`
}
