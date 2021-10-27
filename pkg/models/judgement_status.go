package models

import "database/sql/driver"

type JudgeStatus string

const (
	Pending  JudgeStatus = "Pending"
	Running  JudgeStatus = "Running"
	Canceled JudgeStatus = "Canceled"
	Finished JudgeStatus = "Finished"

	PartiallyCorrect JudgeStatus = "PartiallyCorrect"
	WrongAnswer      JudgeStatus = "WrongAnswer"
	Accepted         JudgeStatus = "Accepted"

	TimeLimitExceeded   JudgeStatus = "TimeLimitExceeded"
	MemoryLimitExceeded JudgeStatus = "MemoryLimitExceeded"
	OutputLimitExceeded JudgeStatus = "OutputLimitExceeded"
	RuntimeError        JudgeStatus = "RuntimeError"
	FileError           JudgeStatus = "FileError"

	SystemError        JudgeStatus = "SystemError"
	JudgementFailed    JudgeStatus = "JudgementFailed"
	CompilationError   JudgeStatus = "CompilationError"
	ConfigurationError JudgeStatus = "ConfigurationError"
	InvalidInteraction JudgeStatus = "InvalidInteraction"
)

func (p *JudgeStatus) Scan(value interface{}) error {
	*p = JudgeStatus(value.(string))
	return nil
}

func (p JudgeStatus) Value() (driver.Value, error) {
	return string(p), nil
}

func (JudgeStatus) TableName() string {
	return "judge_status"
}
