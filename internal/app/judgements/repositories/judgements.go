package repositories

import (
	"sync"

	"github.com/google/uuid"
	"github.com/infinity-oj/server-v2/pkg/models"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type Repository interface {
	Fetch() *models.Judgement
	Create(submissionId, processId uint64) (*models.Judgement, error)
	Update(judgement *models.Judgement) error
}

type DefaultRepository struct {
	logger *zap.Logger
	db     *gorm.DB
	mutex  *sync.Mutex
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

// 将task进行包装
//func (m DefaultRepository) WrapJudgement(judgement *models.Judgement) (*TaskElement, error) {
//
//	judgementId := judgement.JudgementId
//	tp := judgement.Type
//
//	var properties map[string]string
//	propertiesJson := judgement.Property
//	if propertiesJson != "" {
//		if err := json.Unmarshal([]byte(propertiesJson), &properties); err != nil {
//			return nil, err
//		}
//	}
//	inputs, err := crypto.EasyDecode(judgement.Inputs)
//	if err != nil {
//		return nil, err
//	}
//
//	judgementInQueue := &TaskElement{
//		Idle: false,
//
//		JudgementId: judgementId,
//		Type:        tp,
//		Properties:  properties,
//
//		Inputs:  inputs,
//		Outputs: nil,
//
//		obj: judgement,
//	}
//	return judgementInQueue, nil
//
//}

func (m DefaultRepository) Fetch() *models.Judgement {
	return nil
}

func NewJudgementRepository(logger *zap.Logger, db *gorm.DB) Repository {
	return &DefaultRepository{
		logger: logger.With(zap.String("type", "Repository")),
		db:     db,
	}
}
