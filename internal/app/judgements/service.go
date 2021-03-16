package judgements

import (
	"errors"
	"github.com/infinity-oj/server-v2/internal/app/problems"
	"github.com/infinity-oj/server-v2/internal/app/processes"
	"github.com/infinity-oj/server-v2/internal/app/submissions"
	"github.com/infinity-oj/server-v2/internal/lib/schedulers"
	"github.com/infinity-oj/server-v2/pkg/models"
	"go.uber.org/zap"
	"net/http"
)

type Service interface {
	GetJudgement(judgementId string) (*models.Judgement, error)
	GetJudgements(accountId uint64) ([]*models.Judgement, error)
	CreateJudgement(accountId, processId, submissionId uint64) (int, *models.Judgement, error)
	UpdateJudgement(judgementId string, status models.JudgeStatus, score float64, msg string) (*models.Judgement, error)
}

type service struct {
	logger               *zap.Logger
	Repository           Repository
	processRepository    processes.Repository
	submissionRepository submissions.Repository
	problemRepository    problems.Repository

	scheduler schedulers.Scheduler
}

func (s service) UpdateJudgement(judgementId string, status models.JudgeStatus, score float64, msg string) (*models.Judgement, error) {
	s.logger.Debug("update judgement",
		zap.String("judgement id", judgementId),
		zap.String("judge status", string(status)),
		zap.String("msg", msg),
		zap.Float64("score", score),
	)

	// get judgement with judgementId
	judgement, err := s.Repository.GetJudgement(judgementId)
	if err != nil {
		return nil, err
	}

	judgement.Score = score
	judgement.Status = status
	judgement.Msg = msg

	err = s.Repository.Update(judgement)

	return judgement, err
}

func (s service) CreateJudgement(accountId, processId, submissionId uint64) (int, *models.Judgement, error) {
	s.logger.Debug("create judgement",
		zap.Uint64("account id", accountId),
		zap.Uint64("process id", processId),
		zap.Uint64("submission id", submissionId),
	)

	//judgements, err := d.Repository.GetJudgementsByAccountId(accountId)
	//if err != nil {
	//	return http.StatusInternalServerError, nil, err
	//}
	//for _, judgement := range judgements {
	//	if judgement.Status == models.Accepted || judgement.Status == models.Pending {
	//		now := time.Now()
	//		judgeTime := judgement.CreatedAt
	//		dateEquals := func(a time.Time, b time.Time) bool {
	//			y1, m1, d1 := a.Date()
	//			y2, m2, d2 := b.Date()
	//			return y1 == y2 && m1 == m2 && d1 == d2
	//		}
	//		if dateEquals(judgeTime, now) {
	//			return http.StatusForbidden, nil, errors.New("previous judgement accepted today")
	//		}
	//	}
	//}

	// get process
	process, err := s.processRepository.GetProcess(processId)
	if err != nil {
		s.logger.Error("create judgement, get process",
			zap.Uint64("process id", processId),
			zap.Error(err),
		)
		return http.StatusInternalServerError, nil, err
	}
	if process == nil {
		return http.StatusInternalServerError, nil, errors.New("invalid request")
	}
	s.logger.Debug("create judgement",
		zap.String("process definition", process.Definition),
	)

	// get submission
	submission, err := s.submissionRepository.GetSubmissionById(submissionId)
	if err != nil {
		s.logger.Error("create judgement",
			zap.Uint64("submission id", submissionId),
			zap.Error(err),
		)
		return http.StatusInternalServerError, nil, err
	}
	if submission == nil {
		return http.StatusBadRequest, nil, errors.New("invalid request")
	}
	s.logger.Debug("create judgement",
		zap.String("submission user space", submission.UserVolume),
	)

	// create judgement
	judgement, err := s.Repository.Create(submissionId, processId)
	if err != nil {
		s.logger.Error("create judgement",
			zap.Uint64("submission id", submissionId),
			zap.Uint64("process id", processId),
			zap.Error(err),
		)
		return http.StatusInternalServerError, nil, err
	}
	s.logger.Debug("create judgement successfully")

	problem, err := s.problemRepository.GetProblemById(submission.ProblemId)
	if err != nil {
		panic(err)
	}
	r, err := s.scheduler.NewRuntime(problem, submission, judgement, process)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	err = s.scheduler.PushRuntime(r)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, judgement, err
}

func (s service) GetJudgement(judgementId string) (*models.Judgement, error) {
	judgement, err := s.Repository.GetJudgement(judgementId)
	return judgement, err
}

func (s service) GetJudgements(accountId uint64) ([]*models.Judgement, error) {
	judgements, err := s.Repository.GetJudgementsByAccountId(accountId)
	return judgements, err
}

func (s *service) FinishJudgement() {
	ch := s.scheduler.FinishedJudgement()
	for {
		judgement := <-ch
		_, _ = s.UpdateJudgement(judgement.JudgementId, judgement.Status, judgement.Score, judgement.Msg)
	}
}

func NewService(
	logger *zap.Logger,
	s schedulers.Scheduler,
	Repository Repository,
	ProblemRepository problems.Repository,
	ProcessRepository processes.Repository,
	SubmissionRepository submissions.Repository,
) Service {
	pendingJudgements, err := Repository.GetPendingJudgements()
	if err != nil {
		panic(err)
	}

	for _, judgement := range pendingJudgements {
		// get process
		process, err := ProcessRepository.GetProcess(judgement.ProcessId)
		if err != nil {
			panic(err)
		}
		if process == nil {
			continue
		}
		// get submission
		submission, err := SubmissionRepository.GetSubmissionById(judgement.SubmissionId)
		if err != nil {
			panic(err)
		}
		if submission == nil {
			continue
		}
		// get problem
		problem, err := ProblemRepository.GetProblemById(submission.ProblemId)
		if err != nil {
			panic(err)
		}
		if problem == nil {
			continue
		}

		logger.Debug("restore judgement",
			zap.String("judgement id", judgement.JudgementId),
			zap.String("submission user space", submission.UserVolume),
		)
		r, err := s.NewRuntime(problem, submission, judgement, process)
		s.PushRuntime(r)
	}

	srv := &service{
		logger:               logger.With(zap.String("type", "JudgementService")),
		Repository:           Repository,
		processRepository:    ProcessRepository,
		submissionRepository: SubmissionRepository,
		problemRepository:    ProblemRepository,

		scheduler: s,
	}
	go srv.FinishJudgement()

	return srv
}
