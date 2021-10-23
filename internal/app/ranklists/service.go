package ranklists

import (
	"github.com/infinity-oj/server-v2/pkg/models"
	"go.uber.org/zap"
)

type Service interface {
	GetRankList(id uint64) (*models.RankList, error)
	GetRankListsByProblem(problem *models.Problem) ([]*models.RankList, error)
}

type service struct {
	logger     *zap.Logger
	Repository Repository
}

func (s service) GetRankListsByProblem(problem *models.Problem) ([]*models.RankList, error) {
	rls, err := s.Repository.GetRankListsByProblem(problem)
	return rls, err
}

func (s service) GetRankList(id uint64) (*models.RankList, error) {
	rl, err := s.Repository.GetRankList(id)
	if err != nil {
		return nil, err
	}
	latestRecords := make(map[uint64]*models.RankListRecord)
	for i := range rl.Records {
		latestRecords[rl.Records[i].AccountID] = &rl.Records[i]
	}
	var records []models.RankListRecord
	for _, v := range latestRecords {
		records = append(records, *v)
	}
	rl.Records = records
	return rl, nil
}

func NewService(logger *zap.Logger, Repository Repository) Service {
	return &service{
		logger:     logger.With(zap.String("type", "program service")),
		Repository: Repository,
	}
}
