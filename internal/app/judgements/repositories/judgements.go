package repositories

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/infinity-oj/server-v2/pkg/models"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"sync"
)

type Repository interface {
	GetJudgement(judgementId string) (*models.Judgement, error)
	GetJudgementsByAccountId(accountId uint64) ([]*models.Judgement, error)
	Create(submissionId, processId uint64) (*models.Judgement, error)
	Update(judgement *models.Judgement) error
}

type DefaultRepository struct {
	logger *zap.Logger
	db     *gorm.DB
	mutex  *sync.Mutex
}

func (m DefaultRepository) GetJudgementsByAccountId(accountId uint64) (judgements []*models.Judgement, err error) {
	var result []*struct {
		models.Judgement
		models.Submission
	}
	if err := m.db.Model(&models.Judgement{}).
		Joins("left join submissions on judgements.submission_id = submissions.id").
		Where("submissions.submitter_id = ?", accountId).
		Limit(5).
		Scan(&result).
		Error; err != nil {
		fmt.Println(judgements)
		return nil, err
	}
	for _, res := range result {

		judgements = append(judgements, &models.Judgement{
			Model:        models.Model{
				CreatedAt: res.Judgement.CreatedAt,
			},
			SubmissionId: res.SubmissionId,
			ProcessId:    res.ProcessId,
			JudgementId:  res.JudgementId,
			Status:       res.Status,
			Score:        res.Score,
		})

	}
	return
}

func (m DefaultRepository) GetJudgement(judgementId string) (*models.Judgement, error) {

	judgement := &models.Judgement{}
	if err := m.db.Where(&models.Judgement{JudgementId: judgementId}).First(judgement).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			m.logger.Error("query account failed", zap.String("judgement id", judgementId), zap.Error(err))
		}
		return nil, err
	}
	return judgement, nil
}

func (m DefaultRepository) Create(submissionId, processId uint64) (*models.Judgement, error) {
	judgementId := uuid.New().String()
	judgement := &models.Judgement{
		SubmissionId: submissionId,
		ProcessId:    processId,
		JudgementId:  judgementId,
		Status:       models.Pending,
		Score:        -1,
	}

	err := m.db.Save(&judgement).Error

	if err != nil {
		return nil, err
	}
	return judgement, nil
}

func (m DefaultRepository) Update(judgement *models.Judgement) error {
	err := m.db.Save(&judgement).Error
	return err
}

func NewJudgementRepository(logger *zap.Logger, db *gorm.DB) Repository {
	return &DefaultRepository{
		logger: logger.With(zap.String("type", "Repository")),
		db:     db,
	}
}
