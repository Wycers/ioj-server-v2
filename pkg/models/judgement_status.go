package models

import "database/sql/driver"

type JudgeStatus string

const (
	Pending JudgeStatus = "Pending"

	PartiallyCorrect JudgeStatus = "PartiallyCorrect"
	WrongAnswer      JudgeStatus = "WrongAnswer"
	Accepted         JudgeStatus = "Accepted"

	SystemError         JudgeStatus = "SystemError"
	JudgementFailed     JudgeStatus = "JudgementFailed"
	CompilationError    JudgeStatus = "CompilationError"
	FileError           JudgeStatus = "FileError"
	RuntimeError        JudgeStatus = "RuntimeError"
	TimeLimitExceeded   JudgeStatus = "TimeLimitExceeded"
	MemoryLimitExceeded JudgeStatus = "MemoryLimitExceeded"
	OutputLimitExceeded JudgeStatus = "OutputLimitExceeded"
	InvalidInteraction  JudgeStatus = "InvalidInteraction"

	ConfigurationError JudgeStatus = "ConfigurationError"
	Canceled           JudgeStatus = "Canceled"
)

func (p *JudgeStatus) Scan(value interface{}) error {
	*p = JudgeStatus(value.([]byte))
	return nil
}

func (p JudgeStatus) Value() (driver.Value, error) {
	return string(p), nil
}

func (JudgeStatus) TableName() string {
	return "judge_status"
}
