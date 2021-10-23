package handlers

import (
	"github.com/infinity-oj/server-v2/internal/app/accounts"

	"github.com/infinity-oj/server-v2/internal/app/ranklists"
	"github.com/infinity-oj/server-v2/pkg/models"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

type RankList struct {
	rr ranklists.Repository
	ar accounts.Repository
}

func (r RankList) IsMatched(tp string) bool {
	return tp == "ranklist"
}

func (r RankList) Work(process *models.Process) error {
	iRankListID, ok := process.Properties["ranklistID"]
	if !ok {
		return errors.New("ranklist id is not found")
	}
	rankListID := cast.ToUint64(iRankListID)
	rl, err := r.rr.GetRankList(rankListID)
	if err != nil {
		return err
	}
	if rl == nil {
		return errors.New("ranklist is not found")
	}

	iAccountID, ok := process.Properties["accountID"]
	if !ok {
		return errors.New("account id is not found")
	}
	accountID := cast.ToUint64(iAccountID)

	account, err := r.ar.GetAccountById(accountID)
	if err != nil {
		return err
	}
	if rl == nil {
		return errors.New("account is not found")
	}

	iKey, ok := process.Properties["metric_key"]
	if !ok {
		return errors.New("key is not found")
	}
	key := cast.ToString(iKey)

	value := process.Inputs[0].Value

	_, err = r.rr.CreateRankListRecord(rl, account, key, value)

	return err
}

func NewRankList(rr ranklists.Repository, ar accounts.Repository) *RankList {
	return &RankList{
		rr: rr,
		ar: ar,
	}
}
