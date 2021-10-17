package submissions

import (
	"errors"
	"net/http"

	"go.uber.org/zap"

	//"github.com/infinity-oj/server-v2/internal/app/judgements"
	"github.com/infinity-oj/server-v2/internal/app/problems"
	"github.com/infinity-oj/server-v2/internal/lib/scheduler"

	"github.com/infinity-oj/server-v2/pkg/models"
)

type Service interface {
	Create(submitterID uint64, problemName string, userSpace string) (code int, s *models.Submission, j *models.Judgement, err error)
	GetSubmission(submissionId string) (s *models.Submission, err error)
	GetSubmissions(problemId string, page, pageSize int) (res []*models.Submission, err error)
	GetSubmissionsByAccountId(accountId uint64, page, pageSize int) (res []*models.Submission, err error)
}

type service struct {
	logger               *zap.Logger
	SubmissionRepository Repository
	ProblemRepository    problems.Repository
	//JudgementService     judgements.Service

	scheduler scheduler.Scheduler
}

func (d service) GetSubmissionsByAccountId(accountId uint64, page, pageSize int) (res []*models.Submission, err error) {
	offset := (page - 1) * pageSize
	res, err = d.SubmissionRepository.GetSubmissionsByAccount(offset, pageSize, accountId)
	if err != nil {
		d.logger.Error("get submissions", zap.Uint64("account id", accountId), zap.Error(err))
		return nil, err
	}
	return
}

func (d service) GetSubmissions(problemId string, page, pageSize int) (res []*models.Submission, err error) {
	offset := (page - 1) * pageSize
	res, err = d.SubmissionRepository.GetSubmissions(offset, pageSize, problemId)
	if err != nil {
		d.logger.Error("get submissions", zap.String("problem id", problemId), zap.Error(err))
		return nil, err
	}
	return
}

func (d service) GetSubmission(submissionId string) (s *models.Submission, err error) {
	s, err = d.SubmissionRepository.GetSubmission(submissionId)
	if err != nil {
		d.logger.Error("get submission", zap.String("submission id", submissionId), zap.Error(err))
		return nil, err
	}
	return
}

func (d service) Create(submitterID uint64, problemName, userSpace string) (code int, s *models.Submission, j *models.Judgement, err error) {
	d.logger.Debug("create submission",
		zap.Uint64("submitter Id", submitterID),
		zap.String("problem name", problemName),
		zap.String("user space", userSpace),
	)
	problem, err := d.ProblemRepository.GetProblemByName(problemName)
	if err != nil {
		d.logger.Error("create submission", zap.Error(err))
		return http.StatusInternalServerError, nil, nil, err
	}
	if problem == nil {
		d.logger.Error("create submission: unknown problem")
		return http.StatusInternalServerError, nil, nil, errors.New("unknown problem")
	}
	d.logger.Debug("create submission",
		zap.Uint64("submitter Id", submitterID),
		zap.Uint64("problem id", problem.ID),
		zap.String("user space", userSpace),
	)
	s, err = d.SubmissionRepository.Create(submitterID, problem.ID, userSpace)
	if err != nil {
		d.logger.Error("create submission", zap.Error(err))
		return http.StatusInternalServerError, nil, nil, err
	}
	code = http.StatusOK
	//code, j, err = d.JudgementService.CreateJudgement(submitterID, problem.ProgramId, s.ID)
	return
}

func NewService(
	logger *zap.Logger,

	repository Repository,
	problemsRepository problems.Repository,
	//judgementsService judgements.Service,
) Service {
	return &service{
		logger:               logger.With(zap.String("type", "service")),
		SubmissionRepository: repository,
		ProblemRepository:    problemsRepository,
	}
}
