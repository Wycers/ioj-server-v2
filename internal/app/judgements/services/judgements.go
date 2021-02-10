package services

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/infinity-oj/server-v2/internal/lib/scheduler"

	"github.com/infinity-oj/server-v2/pkg/models"

	"github.com/infinity-oj/server-v2/internal/app/judgements/repositories"
	processRepository "github.com/infinity-oj/server-v2/internal/app/processes/repositories"
	submissionRepository "github.com/infinity-oj/server-v2/internal/app/submissions/repositories"
	"go.uber.org/zap"
)

type JudgementsService interface {
	GetJudgement(judgementId string) (*models.Judgement, error)
	GetJudgements(accountId uint64) ([]*models.Judgement, error)
	CreateJudgement(accountId, processId, submissionId uint64) (int, *models.Judgement, error)
	UpdateJudgement(judgementId string, status models.JudgeStatus, score float64, msg string) (*models.Judgement, error)

	GetTasks(taskType string) (task []*models.Task, err error)
	GetTask(taskId string) (task *models.Task, err error)
	UpdateTask(taskId, outputs, warning, error string) (task *models.Task, err error)
	ReserveTask(taskId string) (token string, locked bool, err error)
}

type Service struct {
	logger               *zap.Logger
	Repository           repositories.Repository
	processRepository    processRepository.Repository
	submissionRepository submissionRepository.Repository

	scheduler scheduler.Scheduler
}

func (d Service) GetTasks(taskType string) (tasks []*models.Task, err error) {
	d.scheduler.List()

	for {
		if element := d.scheduler.FetchTask("*", "*", "basic/end", true); element != nil {
			if score, ok := element.Task.Inputs[0].Value.(float64); !ok {
				d.logger.Error("wrong score", zap.Error(err))
			} else {
				if _, err := d.UpdateJudgement(element.JudgementId, models.Accepted, score, ""); err != nil {
					return nil, err
				}
			}
			err = d.scheduler.FinishTask(element, []string{})
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	d.logger.Info("get task", zap.String("type", taskType))
	element := d.scheduler.FetchTask("*", "*", taskType, false)
	if element != nil {
		d.logger.Info("get tasks", zap.String("judgement id", element.JudgementId))
		tasks = []*models.Task{
			element.Task,
		}
	} else {
		d.logger.Info("get tasks: nothing")
	}
	return
}

func (d Service) GetTask(taskId string) (task *models.Task, err error) {
	d.logger.Info("get task",
		zap.String("task id", taskId),
	)
	element := d.scheduler.FetchTask("*", taskId, "*", true)
	if element != nil {
		d.logger.Info("get task",
			zap.String("judgement id", element.JudgementId),
			zap.String("task id", element.Task.TaskId),
		)
		task = element.Task
	} else {
		d.logger.Debug("get tasks: nothing")
	}
	return
}

func (d Service) UpdateTask(taskId, outputs, warning, error string) (task *models.Task, err error) {
	taskElement := d.scheduler.FetchTask("*", taskId, "*", true)
	if taskElement == nil {
		d.logger.Debug("invalid token: no such task",
			zap.String("task id", taskId),
		)
		d.scheduler.UnlockTask(taskElement)
		return nil, errors.New("invalid token")
	}

	task = taskElement.Task

	if task.TaskId != taskId {
		d.logger.Debug("task mismatch",
			zap.String("expected task id", task.TaskId),
			zap.String("actual task id", taskId),
		)
		d.scheduler.UnlockTask(taskElement)
		return nil, errors.New("task mismatch")
	}

	d.logger.Info("update task",
		zap.String("task id", taskId),
	)

	if error != "" {
		d.scheduler.RemoveTask(taskElement)
		_, err := d.UpdateJudgement(taskElement.JudgementId, models.SystemError, 0, fmt.Sprintf("warning: %s\nerror: %s\n", warning, error))
		if err != nil {
			d.logger.Error("finish task failed", zap.Error(err))
			return nil, err
		}
		return task, nil
	}

	//update task
	//err := d.Repository.Update(element, outputs)
	//if err != nil {
	//	d.logger.Error("update task", zap.Error(err))
	//	return nil, err
	//}

	err = d.scheduler.FinishTask(taskElement, strings.Split(outputs, ","))

	// calculate next task
	if err != nil {
		d.logger.Error("create judgement: initial process failed",
			zap.String("task id", taskId),
			zap.Error(err),
		)
	}

	return task, nil
}

func (d Service) ReserveTask(taskId string) (token string, locked bool, err error) {
	taskElement := d.scheduler.FetchTask("*", taskId, "*", true)

	if taskElement == nil {
		return "", false, errors.New("not found")
	}

	if !d.scheduler.LockTask(taskElement) {
		return "", false, errors.New("participated")
	}

	token = uuid.New().String()
	d.logger.Debug("reserve task",
		zap.String("task id", taskId),
		zap.String("token", token),
	)

	return token, true, nil
}

func (d Service) UpdateJudgement(judgementId string, status models.JudgeStatus, score float64, msg string) (*models.Judgement, error) {
	d.logger.Debug("update judgement",
		zap.String("judgement id", judgementId),
		zap.String("judge status", string(status)),
		zap.String("msg", msg),
		zap.Float64("score", score),
	)

	// get judgement with judgementId
	judgement, err := d.Repository.GetJudgement(judgementId)
	if err != nil {
		return nil, err
	}

	judgement.Score = score
	judgement.Status = status
	judgement.Msg = msg

	err = d.Repository.Update(judgement)

	return judgement, err
}

func (d Service) CreateJudgement(accountId, processId, submissionId uint64) (int, *models.Judgement, error) {
	d.logger.Debug("create judgement",
		zap.Uint64("account id", accountId),
		zap.Uint64("process id", processId),
		zap.Uint64("submission id", submissionId),
	)

	judgements, err := d.Repository.GetJudgementsByAccountId(accountId)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	for _, judgement := range judgements {
		if judgement.Status == models.Accepted || judgement.Status == models.Pending {
			now := time.Now()
			judgeTime := judgement.CreatedAt
			dateEquals := func(a time.Time, b time.Time) bool {
				y1, m1, d1 := a.Date()
				y2, m2, d2 := b.Date()
				return y1 == y2 && m1 == m2 && d1 == d2
			}
			if dateEquals(judgeTime, now) {
				return http.StatusForbidden, nil, errors.New("previous judgement accepted today")
			}
		}
	}

	// get process
	process, err := d.processRepository.GetProcess(processId)
	if err != nil {
		d.logger.Error("create judgement, get process",
			zap.Uint64("process id", processId),
			zap.Error(err),
		)
		return http.StatusInternalServerError, nil, err
	}
	if process == nil {
		return http.StatusInternalServerError, nil, errors.New("invalid request")
	}
	d.logger.Debug("create judgement",
		zap.String("process definition", process.Definition),
	)

	// get submission
	submission, err := d.submissionRepository.GetSubmissionById(submissionId)
	if err != nil {
		d.logger.Error("create judgement",
			zap.Uint64("submission id", submissionId),
			zap.Error(err),
		)
		return http.StatusInternalServerError, nil, err
	}
	if submission == nil {
		return http.StatusBadRequest, nil, errors.New("invalid request")
	}
	d.logger.Debug("create judgement",
		zap.String("submission user space", submission.UserVolume),
	)

	// create judgement
	judgement, err := d.Repository.Create(submissionId, processId)
	if err != nil {
		d.logger.Error("create judgement",
			zap.Uint64("submission id", submissionId),
			zap.Uint64("process id", processId),
			zap.Error(err),
		)
		return http.StatusInternalServerError, nil, err
	}
	d.logger.Debug("create judgement successfully")

	err = d.scheduler.NewProcessRuntime(submission, judgement, process)

	return http.StatusOK, judgement, err
}

func (d Service) GetJudgement(judgementId string) (*models.Judgement, error) {
	judgement, err := d.Repository.GetJudgement(judgementId)
	return judgement, err
}

func (d Service) GetJudgements(accountId uint64) ([]*models.Judgement, error) {
	judgements, err := d.Repository.GetJudgementsByAccountId(accountId)
	return judgements, err
}

func NewJudgementsService(
	logger *zap.Logger,
	Repository repositories.Repository,
	ProcessRepository processRepository.Repository,
	SubmissionRepository submissionRepository.Repository,
) JudgementsService {
	s := scheduler.New(logger)

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
		logger.Debug("restore judgement",
			zap.String("judgement id", judgement.JudgementId),
			zap.String("submission user space", submission.UserVolume),
		)
		err = s.NewProcessRuntime(submission, judgement, process)
	}

	return &Service{
		logger:               logger.With(zap.String("type", "JudgementService")),
		Repository:           Repository,
		processRepository:    ProcessRepository,
		submissionRepository: SubmissionRepository,

		scheduler: s,
	}
}
