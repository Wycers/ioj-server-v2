package models

type Page struct {
	Model

	ProblemId uint64 `json:"problemId"`

	Locale string `json:"locale"`

	Title       string `json:"title" gorm:"not null"` // title
	Description string `json:"description"`

	InputFormat  string `json:"input_format"`
	OutputFormat string `json:"output_format"`
	Example      string `json:"example"`
	HintAndLimit string `json:"hint_and_limit"`
}
