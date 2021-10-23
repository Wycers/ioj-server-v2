package ranklists

import (
	"github.com/infinity-oj/server-v2/pkg/models"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	CreateRankList() (p *models.RankList, err error)
	CreateRankListModel(id uint64) (*models.RankListModel, error)
	CreateRankListRecord(rl *models.RankList, account *models.Account, key string, value interface{}) (*models.RankListRecord, error)
	GetRankList(id uint64) (*models.RankList, error)
	GetRankListsByProblem(problem *models.Problem) ([]*models.RankList, error)
}

type repository struct {
	logger *zap.Logger
	db     *gorm.DB
}

func (m repository) GetRankListsByProblem(problem *models.Problem) ([]*models.RankList, error) {
	var rl []*models.RankList
	if err := m.db.Model(&models.RankList{}).Where("problem_id = ?", problem.ID).Find(&rl).Error; err != nil {
		return nil, err
	}
	return rl, nil
}

func (m repository) CreateRankListModel(id uint64) (*models.RankListModel, error) {
	rm := &models.RankListModel{
		RankListID: 0,
		Key:        "key",
		Priority:   1,
		Order:      "dec",
	}
	if err := m.db.Create(&rm).Error; err != nil {
		return nil, err
	}
	return rm, nil
}

func (m repository) CreateRankListRecord(rl *models.RankList, account *models.Account, key string, value interface{}) (*models.RankListRecord, error) {
	rlc := &models.RankListRecord{
		Key:   key,
		Value: cast.ToFloat64(value),
	}
	if err := m.db.Create(rlc).Error; err != nil {
		m.logger.Error("create ranklist record error",
			zap.Any("ranklist", rl), zap.Any("record", rlc), zap.Error(err))
		return nil, err
	}
	if err := m.db.Model(rlc).Association("Account").Append(account); err != nil {
		m.logger.Error("create ranklist record error",
			zap.Any("ranklist", rl), zap.Any("record", rlc), zap.Error(err))
		return nil, err
	}
	if err := m.db.Model(rl).Association("Records").Append(rlc); err != nil {
		m.logger.Error("create ranklist record error",
			zap.Any("ranklist", rl), zap.Any("record", rlc), zap.Error(err))
		return nil, err
	}
	return rlc, nil
}

func (m repository) GetRankList(id uint64) (*models.RankList, error) {
	rl := &models.RankList{}
	if err := m.db.Model(rl).Preload("Records.Account").Preload(clause.Associations).First(rl, id).Error; err != nil {
		return nil, err
	}
	return rl, nil
}

func (m repository) CreateRankList() (rl *models.RankList, err error) {
	rl = &models.RankList{
		Model:     models.Model{},
		ProblemID: 0,
		Name:      "",
		Title:     "",
		Models: []models.RankListModel{
			{
				Key:      "score",
				Priority: 0,
				Order:    "inc",
			},
		},
		Records: nil,
	}
	return rl, m.db.Create(rl).Error
}

func NewRepository(logger *zap.Logger, db *gorm.DB) Repository {
	return &repository{
		logger: logger.With(zap.String("type", " repository")),
		db:     db,
	}
}
