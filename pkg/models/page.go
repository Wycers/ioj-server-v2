package models

type Page struct {
	Model

	ProblemId uint64 `json:"problemId" gorm:"index: problem_id"`

	Locale string `json:"locale" gorm:"index: locale"`

	Title       string `json:"title" gorm:"index: title; not null"` // title
	Description string `json:"description"`

	InputFormat  string `json:"input_format"`
	OutputFormat string `json:"output_format"`
	Example      string `json:"example"`
	HintAndLimit string `json:"hint_and_limit"`
}
