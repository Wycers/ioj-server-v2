package judgements

import (
	"errors"
	"fmt"
	"sync"

	"github.com/spf13/cast"

	"github.com/google/uuid"
	"github.com/infinity-oj/server-v2/pkg/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository interface {
	GetJudgement(judgementId string) (*models.Judgement, error)
	GetJudgementsByAccountId(accountId uint64) ([]*models.Judgement, error)
	GetPendingJudgements() ([]*models.Judgement, error)
	Create(blueprintId uint64, args map[string]interface{}) (*models.Judgement, error)
	Update(judgement *models.Judgement) error
}

type repository struct {
	logger *zap.Logger
	db     *gorm.DB
	mutex  *sync.Mutex
}

func (m repository) GetPendingJudgements() ([]*models.Judgement, error) {
	var res []*models.Judgement
	if err := m.db.
		Model(&models.Judgement{}).
		Where("status = ?", models.Pending).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (m repository) GetJudgementsByAccountId(accountId uint64) (judgements []*models.Judgement, err error) {
	var result []*struct {
		models.Judgement
		models.Submission
	}
	if err := m.db.Model(&models.Judgement{}).
		Joins("left join submissions on judgements.submission_id = submissions.id").
		Where("submissions.submitter_id = ?", accountId).
		Order("judgements.id desc").
		Scan(&result).
		Error; err != nil {
		fmt.Println(judgements)
		return nil, err
	}
	for _, res := range result {
		judgements = append(judgements, &models.Judgement{
			Model: models.Model{
				CreatedAt: res.Judgement.CreatedAt,
			},
			BlueprintId: res.BlueprintId,
			Name:        res.Judgement.Name,
			Status:      res.Judgement.Status,
			Score:       res.Judgement.Score,
		})

	}
	return
}

func (m repository) GetJudgement(judgementId string) (*models.Judgement, error) {

	judgement := &models.Judgement{}
	if err := m.db.Where(&models.Judgement{Name: judgementId}).First(judgement).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			m.logger.Error("query account failed", zap.String("judgement id", judgementId), zap.Error(err))
		}
		return nil, err
	}
	return judgement, nil
}

func (m repository) Create(blueprintId uint64, args map[string]interface{}) (*models.Judgement, error) {
	submissionName := cast.ToString(args["submission"])
	if submissionName == "" {
		return nil, errors.New("submission is required")
	}
	submission := &models.Submission{}
	if err := m.db.First(submission, "name = ?", submissionName).Error; err != nil {
		return nil, err
	}
	judgement := &models.Judgement{
		BlueprintId: blueprintId,
		Name:        uuid.New().String(),
		Args:        args,
		Status:      models.Pending,
		Msg:         "",
		Score:       -1,
	}
	if err := m.db.Model(submission).Association("Judgements").Append(judgement); err != nil {
		return nil, err
	}
	return judgement, nil
}

func (m repository) Update(judgement *models.Judgement) error {
	err := m.db.Save(&judgement).Error
	return err
}

func NewRepository(logger *zap.Logger, db *gorm.DB) Repository {
	return &repository{
		logger: logger.With(zap.String("type", "repository")),
		db:     db,
	}
}
