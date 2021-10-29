package handlers

import (
	"github.com/infinity-oj/server-v2/internal/app/judgements"
	"github.com/infinity-oj/server-v2/internal/lib/manager"
	"github.com/spf13/cast"
)

type Result struct {
	jr judgements.Repository
}

func (r *Result) IsMatched(tp string) bool {
	return tp == "result"
}

func (r *Result) Work(pr *manager.ProcessRuntime) error {
	pr.Mutex.Lock()
	defer pr.Mutex.Unlock()
	v := pr.Process.Inputs[0].Value
	score := cast.ToFloat64(v)
	pr.Judgement.Score = score
	return nil
}

func NewResult(jr judgements.Repository) *Result {
	return &Result{jr: jr}
}
