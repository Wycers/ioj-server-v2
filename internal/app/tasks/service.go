package tasks

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/infinity-oj/server-v2/internal/lib/schedulers"
	"github.com/infinity-oj/server-v2/pkg/models"
	"go.uber.org/zap"
)

type Service interface {
	GetTasks(taskType string) (task []*models.Task, err error)
	GetTask(taskId string) (task *models.Task, err error)
	UpdateTask(taskId, warning, error string, outputs *models.Slots) (task *models.Task, err error)
	ReserveTask(taskId string) (token string, locked bool, err error)
}

type service struct {
	logger *zap.Logger

	scheduler schedulers.Scheduler
}

func (d service) GetTasks(taskType string) (tasks []*models.Task, err error) {
	d.scheduler.List()

	d.logger.Info("get task", zap.String("type", taskType))
	element := d.scheduler.FetchTask("*", "*", taskType, false)
	if element != nil {
		d.logger.Info("get tasks", zap.String("task id", element.Task.JudgementId))
		tasks = []*models.Task{
			element.Task,
		}
	} else {
		d.logger.Info("get tasks: nothing")
	}
	return
}

func (d service) GetTask(taskId string) (task *models.Task, err error) {
	d.logger.Info("get task",
		zap.String("task id", taskId),
	)
	element := d.scheduler.FetchTask("*", taskId, "*", true)
	if element != nil {
		d.logger.Info("get task",
			zap.String("judgement id", element.Task.JudgementId),
			zap.String("task id", element.Task.TaskId),
		)
		task = element.Task
	} else {
		d.logger.Debug("get tasks: nothing")
	}
	return
}

func (d service) UpdateTask(taskId, warning, error string, outputs *models.Slots) (task *models.Task, err error) {
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
		if err := d.scheduler.FinishTaskWithError(
			taskElement,
			fmt.Sprintf("warning: %s\nerror: %s\n", warning, error),
		); err != nil {
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

	err = d.scheduler.FinishTask(taskElement, outputs)

	// calculate next task
	if err != nil {
		d.logger.Error("update task: finish task failed",
			zap.String("task id", taskId),
			zap.Error(err),
		)
		d.scheduler.UnlockTask(taskElement)
		return nil, err
	}

	return task, nil
}

func (d service) ReserveTask(taskId string) (token string, locked bool, err error) {
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

func NewService(
	logger *zap.Logger,
) Service {
	s := schedulers.New(logger)

	return &service{
		logger: logger.With(zap.String("type", "Task service")),

		scheduler: s,
	}
}
