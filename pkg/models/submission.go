package models

type Submission struct {
	Model

	Name string `json:"submissionId"`

	SubmitterId uint64 `json:"submitterId"`

	ProblemId uint64 `json:"problemId"`

	UserVolume string `json:"userVolume"`
}
