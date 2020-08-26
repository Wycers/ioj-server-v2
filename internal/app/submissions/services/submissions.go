package services

import (
	"errors"

	"go.uber.org/zap"

	problemRepository "github.com/infinity-oj/server-v2/internal/app/problems/repositories"
	"github.com/infinity-oj/server-v2/internal/app/submissions/repositories"
	"github.com/infinity-oj/server-v2/internal/pkg/models"
)

type SubmissionsService interface {
	Create(submitterID uint64, problemName string, userSpace string) (s *models.Submission, err error)
	GetSubmission(submissionId string) (s *models.Submission, err error)
	GetSubmissions(problemId string, page, pageSize int) (res []*models.Submission, err error)
}

type DefaultSubmissionService struct {
	logger               *zap.Logger
	SubmissionRepository repositories.Repository
	ProblemRepository    problemRepository.Repository

	processMap map[string]string
	idMap      map[string]int
}

func (d DefaultSubmissionService) GetSubmissions(problemId string, page, pageSize int) (res []*models.Submission, err error) {
	offset := (page - 1) * pageSize
	res, err = d.SubmissionRepository.GetSubmissions(offset, pageSize, problemId)
	if err != nil {
		d.logger.Error("get submissions", zap.String("problem id", problemId), zap.Error(err))
		return nil, err
	}
	return
}

func (d DefaultSubmissionService) GetSubmission(submissionId string) (s *models.Submission, err error) {
	s, err = d.SubmissionRepository.GetSubmission(submissionId)
	if err != nil {
		d.logger.Error("get submission", zap.String("submission id", submissionId), zap.Error(err))
		return nil, err
	}
	return
}

func (d DefaultSubmissionService) Create(submitterID uint64, problemName, userSpace string) (s *models.Submission, err error) {
	d.logger.Debug("create submission",
		zap.Uint64("submitter Id", submitterID),
		zap.String("problem name", problemName),
		zap.String("user space", userSpace),
	)
	problem, err := d.ProblemRepository.GetProblemByName(problemName)
	if err != nil {
		d.logger.Error("create submission", zap.Error(err))
		return nil, err
	}
	if problem == nil {
		d.logger.Error("create submission: unknown problem")
		return nil, errors.New("unknown problem")
	}
	d.logger.Debug("create submission",
		zap.Uint64("submitter Id", submitterID),
		zap.Uint64("problem id", problem.ID),
		zap.String("user space", userSpace),
	)
	s, err = d.SubmissionRepository.Create(submitterID, problem.ID, userSpace)
	if err != nil {
		d.logger.Error("create submission", zap.Error(err))
		return nil, err
	}
	return
}

func NewSubmissionService(
	logger *zap.Logger,
	Repository repositories.Repository,
	pRepository problemRepository.Repository,
) SubmissionsService {
	return &DefaultSubmissionService{
		logger:               logger.With(zap.String("type", "DefaultSubmissionService")),
		SubmissionRepository: Repository,
		ProblemRepository:    pRepository,
		processMap:           make(map[string]string),
		idMap:                make(map[string]int),
	}
}
