package processes

import (
	"errors"
	"fmt"

	manager2 "github.com/infinity-oj/server-v2/internal/lib/manager"

	"github.com/google/uuid"
	"github.com/infinity-oj/server-v2/pkg/models"
	"go.uber.org/zap"
)

type Service interface {
	GetProcesses(processType string) (process []*models.Process, err error)
	GetProcess(processId string) (process *models.Process, err error)
	UpdateProcess(processId, warning, error string, outputs *models.Slots) (process *models.Process, err error)
	ReserveProcess(processId string) (token string, locked bool, err error)
}

type service struct {
	logger *zap.Logger

	manager manager2.ProcessManager
}

func (d service) GetProcesses(processType string) (processes []*models.Process, err error) {
	d.manager.List()

	d.logger.Info("get process", zap.String("type", processType))
	element := d.manager.Fetch("*", "*", processType, false)
	if element != nil {
		d.logger.Info("get processes", zap.String("process id", element.Process.JudgementId))
		processes = []*models.Process{
			element.Process,
		}
	} else {
		d.logger.Info("get processes: nothing")
	}
	return
}

func (d service) GetProcess(processId string) (process *models.Process, err error) {
	d.logger.Info("get process",
		zap.String("process id", processId),
	)
	element := d.manager.Fetch("*", processId, "*", true)
	if element != nil {
		d.logger.Info("get process",
			zap.String("judgement id", element.Process.JudgementId),
			zap.String("process id", element.Process.ProcessId),
		)
		process = element.Process
	} else {
		d.logger.Debug("get processes: nothing")
	}
	return
}

func (d service) UpdateProcess(processId, warning, error string, outputs *models.Slots) (process *models.Process, err error) {
	processElement := d.manager.Fetch("*", processId, "*", true)
	if processElement == nil {
		d.logger.Debug("invalid token: no such process",
			zap.String("process id", processId),
		)
		//d.manager.UnlockProcess(processElement)
		return nil, errors.New("invalid token")
	}

	process = processElement.Process

	if process.ProcessId != processId {
		d.logger.Debug("process mismatch",
			zap.String("expected process id", process.ProcessId),
			zap.String("actual process id", processId),
		)
		//d.manager.UnlockProcess(processElement)
		return nil, errors.New("process mismatch")
	}

	d.logger.Info("update process",
		zap.String("process id", processId),
	)

	if error != "" {
		if err := d.manager.FinishWithError(
			processElement,
			fmt.Sprintf("warning: %s\nerror: %s\n", warning, error),
		); err != nil {
			d.logger.Error("finish process failed", zap.Error(err))
			return nil, err
		}
		return process, nil
	}

	//update process
	//err := d.repository.Update(element, outputs)
	//if err != nil {
	//	d.logger.Error("update process", zap.Error(err))
	//	return nil, err
	//}

	err = d.manager.Finish(processElement, outputs)

	// calculate next process
	if err != nil {
		d.logger.Error("update process: finish process failed",
			zap.String("process id", processId),
			zap.Error(err),
		)
		//d.manager.UnlockProcess(processElement)
		return nil, err
	}

	return process, nil
}

func (d service) ReserveProcess(processId string) (token string, locked bool, err error) {
	processElement := d.manager.Fetch("*", processId, "*", true)

	if processElement == nil {
		return "", false, errors.New("not found")
	}

	if !d.manager.Reserve(processElement) {
		return "", false, errors.New("reserved")
	}

	token = uuid.New().String()
	d.logger.Debug("reserve process",
		zap.String("process id", processId),
		zap.String("token", token),
	)

	return token, true, nil
}

func NewService(logger *zap.Logger, manager manager2.ProcessManager) Service {
	return &service{
		logger: logger.With(zap.String("type", "Process service")),

		manager: manager,
	}
}
