package submissions

import (
	"container/list"
	"sync"

	"github.com/google/uuid"

	"github.com/infinity-oj/server-v2/pkg/models"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository interface {
	GetSubmission(submissionId string) (*models.Submission, error)
	GetSubmissionById(id uint64) (*models.Submission, error)
	GetSubmissions(offset, limit int, problemId string) ([]*models.Submission, error)
	// TODO: Improve this...
	GetSubmissionsByAccount(offset, limit int, accountId uint64) ([]*models.Submission, error)
	Create(submitterID, problemId uint64, userSpace string) (s *models.Submission, err error)
	Update(s *models.Submission) error
}

type repository struct {
	logger *zap.Logger
	db     *gorm.DB
	queue  *list.List
	mutex  *sync.Mutex
}

func (m repository) GetSubmissionsByAccount(offset, limit int, accountId uint64) (res []*models.Submission, err error) {
	if err = m.db.Model(&models.Submission{}).Where("submitter_id = ?", accountId).
		Offset(offset).Limit(limit).
		Find(&res).Error; err != nil {
		return nil, err
	}
	return
}

func (m repository) GetSubmissionById(id uint64) (*models.Submission, error) {
	submission := &models.Submission{}
	err := m.db.First(&submission, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return submission, nil
}

func (m repository) GetSubmissions(offset, limit int, problemId string) (res []*models.Submission, err error) {
	if err = m.db.Table("submissions").Where("problem_id = ?", problemId).
		Offset(offset).Limit(limit).
		Find(&res).Error; err != nil {
		return nil, err
	}
	return
}

func (m repository) GetSubmission(submissionId string) (*models.Submission, error) {
	submission := &models.Submission{}
	err := m.db.Where("submission_id = ?", submissionId).First(&submission).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return submission, nil
}

func (m repository) Create(submitterId, problemId uint64, userSpace string) (s *models.Submission, err error) {
	s = &models.Submission{
		Name:        uuid.New().String(),
		SubmitterId: submitterId,
		ProblemId:   problemId,
		UserVolume:  userSpace,
	}

	if err = m.db.Create(s).Error; err != nil {
		return nil, errors.Wrapf(err,
			" create submission with username: %d, problemID: %s, userSpace: %s",
			submitterId, problemId, userSpace,
		)
	}

	return
}

func (m repository) Update(s *models.Submission) (err error) {
	err = m.db.Save(s).Error
	return
}

func NewRepository(logger *zap.Logger, db *gorm.DB) Repository {
	return &repository{
		logger: logger.With(zap.String("type", "repository")),
		db:     db,
		queue:  list.New(),
		mutex:  &sync.Mutex{},
	}
}
