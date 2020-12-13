package services

import (
	"errors"
	"strings"
	"sync"

	"github.com/google/uuid"

	"github.com/infinity-oj/server-v2/internal/lib/scheduler"

	"github.com/infinity-oj/server-v2/pkg/models"

	"github.com/infinity-oj/server-v2/internal/app/judgements/repositories"
	processRepository "github.com/infinity-oj/server-v2/internal/app/processes/repositories"
	submissionRepository "github.com/infinity-oj/server-v2/internal/app/submissions/repositories"
	"go.uber.org/zap"
)

type JudgementsService interface {
	GetJudgement() (*models.Judgement, error)
	GetJudgements() ([]*models.Judgement, error)
	CreateJudgement(accountId, processId, submissionId uint64) (*models.Judgement, error)
	UpdateJudgement(judgementId string, score int) (*models.Judgement, error)

	GetTasks(taskType string) (task []*models.Task, err error)
	GetTask(taskId string) (task *models.Task, err error)
	UpdateTask(token, taskId, outputs string) (task *models.Task, err error)
	ReserveTask(taskId string) (token string, err error)
}

type Service struct {
	mutex *sync.Mutex

	logger               *zap.Logger
	Repository           repositories.Repository
	processRepository    processRepository.Repository
	submissionRepository submissionRepository.Repository

	scheduler scheduler.Scheduler
	tokenMap  map[string]string
}

func (d Service) GetTasks(taskType string) (tasks []*models.Task, err error) {
	d.scheduler.List()

	d.logger.Info("get task", zap.String("type", taskType))
	element := d.scheduler.FetchTask("*", "*", taskType)
	if element != nil {
		d.logger.Info("get tasks", zap.String("judgement id", element.JudgementId))
		// TODO: use jwt
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
	element := d.scheduler.FetchTask("*", taskId, "*")
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

func (d Service) UpdateTask(token, taskId, outputs string) (task *models.Task, err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	taskId, ok := d.tokenMap[token]

	if !ok {
		d.logger.Debug("invalid token: no such token")
		return nil, errors.New("invalid token")
	}

	// token should only be used once
	delete(d.tokenMap, token)

	taskElement := d.scheduler.FetchTask("*", taskId, "*")
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

	// update task
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

func (d Service) ReserveTask(taskId string) (token string, err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	taskElement := d.scheduler.FetchTask("*", taskId, "*")

	if !d.scheduler.LockTask(taskElement) {
		return "", errors.New("participated")
	}

	token = uuid.New().String()
	d.tokenMap[token] = taskId
	d.logger.Debug("reserve task",
		zap.String("task id", taskId),
		zap.String("token", token),
	)

	return token, nil
}

func (d Service) UpdateJudgement(judgementId string, score int) (*models.Judgement, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.logger.Debug("update judgement",
		zap.String("judgement id", judgementId),
		zap.Int("judgement id", score),
	)

	// get judgement with judgementId

	// change score to `score`

	// save judgement

	return nil, nil
}

func (d Service) CreateJudgement(accountId, processId, submissionId uint64) (*models.Judgement, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.logger.Debug("create judgement",
		zap.Uint64("account id", accountId),
		zap.Uint64("process id", processId),
		zap.Uint64("submission id", submissionId),
	)

	// get process
	process, err := d.processRepository.GetProcess(processId)
	if err != nil {
		d.logger.Error("create judgement, get process",
			zap.Uint64("process id", processId),
			zap.Error(err),
		)
		return nil, err
	}
	if process == nil {
		return nil, errors.New("invalid request")
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
		return nil, err
	}
	if submission == nil {
		return nil, errors.New("invalid request")
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
		return nil, err
	}
	d.logger.Debug("create judgement successfully")

	err = d.scheduler.NewProcessRuntime(judgement, process)

	return judgement, err
}

func (d Service) GetJudgement() (*models.Judgement, error) {
	panic("implement me")
}

func (d Service) GetJudgements() ([]*models.Judgement, error) {
	panic("implement me")
}

func NewJudgementsService(
	logger *zap.Logger,
	Repository repositories.Repository,
	ProcessRepository processRepository.Repository,
	SubmissionRepository submissionRepository.Repository,
) JudgementsService {
	return &Service{
		mutex:                &sync.Mutex{},
		logger:               logger.With(zap.String("type", "DefaultJudgementService")),
		Repository:           Repository,
		processRepository:    ProcessRepository,
		submissionRepository: SubmissionRepository,

		scheduler: scheduler.New(logger),

		tokenMap: make(map[string]string),
	}
}
