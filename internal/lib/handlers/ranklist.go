package handlers

import "github.com/infinity-oj/server-v2/pkg/models"

type RankList struct {
}

func (r RankList) IsMatched(tp string) bool {
	return tp == "ranklist"
}

func (r RankList) Work(process *models.Process) error {
	panic("implement me")
}

func NewRankList() *RankList {
	return &RankList{}
}
