package handlers

import (
	"math"
	"strconv"

	"github.com/infinity-oj/server-v2/internal/app/judgements"
	"github.com/infinity-oj/server-v2/pkg/models"
)

type Result struct {
	jr judgements.Repository
}

func (r *Result) IsMatched(tp string) bool {
	return tp == "result"
}

func (r *Result) Work(process *models.Process) error {
	v := process.Inputs[0].Value
	score := func() float64 {
		switch i := v.(type) {
		case float64:
			return i
		case float32:
			return float64(i)
		case int64:
			return float64(i)
		case int32:
			return float64(i)
		case string:
			if s, err := strconv.ParseFloat(i, 64); err == nil {
				return s
			}
		}
		return math.NaN()
	}()
	judgement, err := r.jr.GetJudgement(process.JudgementId)
	if err != nil {
		return err
	}
	judgement.Score = score
	return r.jr.Update(judgement)
}

func NewResult(jr judgements.Repository) *Result {
	return &Result{jr: jr}
}
