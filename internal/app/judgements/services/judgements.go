package services

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/infinity-oj/server-v2/internal/pkg/crypto"

	"github.com/infinity-oj/server-v2/pkg/nodeEngine"

	"github.com/infinity-oj/server-v2/pkg/models"

	"github.com/google/uuid"
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

	CreateTask(judgementId, taskType string, blockId int, properties, inputs string) (*models.Task, error)
	GetTasks(taskType string) (task []*models.Task, err error)
	GetTask(taskId string) (task *models.Task, err error)
	UpdateTask(token, taskId, outputs string) (task *models.Task, err error)
	ReserveTask(taskId string) (token string, err error)
}

type ProcessRuntime struct {
	judgement *models.Judgement
	graph     *nodeEngine.Graph
	result    map[int][]byte
}

func (se ProcessRuntime) SetOutputs(blockId int, outputs [][]byte) error {

	block := se.graph.FindBlockById(blockId)

	if len(block.Output) != len(outputs) {

		msg := fmt.Sprintf("output slots mismatch, block %d expects %d but %d",
			blockId,
			len(block.Output),
			len(outputs),
		)
		return errors.New(msg)
	}

	for index, output := range outputs {
		links := se.graph.FindLinkBySourcePort(blockId, index)
		for _, link := range links {
			se.result[link.Id] = output
		}
	}

	block.Done()
	return nil
}

type ProcessElement struct {
	processRuntime *ProcessRuntime
	blockId        int
	taskElement    *repositories.TaskElement
}

type DefaultJudgementsService struct {
	mutex *sync.Mutex

	logger               *zap.Logger
	Repository           repositories.Repository
	processRepository    processRepository.Repository
	submissionRepository submissionRepository.Repository

	tokenMap   map[string]string
	processMap map[string]*ProcessRuntime
	queue      *list.List
}

func (d DefaultJudgementsService) UpdateJudgement(judgementId string, score int) (*models.Judgement, error) {
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

func (d DefaultJudgementsService) CreateJudgement(accountId, processId, submissionId uint64) (*models.Judgement, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.logger.Debug("create judgement",
		zap.Uint64("account id", accountId),
		zap.Uint64("processes id", processId),
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
		zap.String("submission user space", submission.UserSpace),
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

	// create task

	graph, err := nodeEngine.NewGraphByDefinition(process.Definition)
	if err != nil {
		// TODO: check while creating process
		d.logger.Error("parse process definition failed",
			zap.Uint64("submission id", submissionId),
			zap.Uint64("process id", processId),
			zap.Error(err),
		)
		return nil, err
	}

	result := make(map[int][]byte)
	element := &ProcessRuntime{
		judgement: judgement,
		graph:     graph,
		result:    result,
	}

	d.queue.PushBack(element)
	err = d.PushProcessRuntime(element)
	if err != nil {
		d.logger.Error("create judgement: initial process failed",
			zap.Uint64("submission id", submissionId),
			zap.Uint64("process id", processId),
			zap.Error(err),
		)
	}
	d.processMap[judgement.JudgementId] = element

	return judgement, nil
}

func (d DefaultJudgementsService) GetJudgement() (*models.Judgement, error) {
	panic("implement me")
}

func (d DefaultJudgementsService) GetJudgements() ([]*models.Judgement, error) {
	panic("implement me")
}

func (d DefaultJudgementsService) PushProcessRuntime(pr *ProcessRuntime) error {

	ids := pr.graph.Run()

	for _, block := range ids {
		var inputs [][]byte
		for _, linkId := range block.Inputs {
			if data, ok := pr.result[linkId]; ok {
				inputs = append(inputs, data)
			} else {
				return errors.New("wrong process definition")
			}
		}

		properties, err := json.Marshal(block.Properties)
		if err != nil {
			return err
		}

		_, _ = d.CreateTask(
			pr.judgement.JudgementId,
			block.Type,
			block.Id,
			string(properties),
			crypto.EasyEncode(inputs),
		)
	}
	return nil
}

func (d DefaultJudgementsService) CreateTask(judgementId, taskType string, blockId int, properties, inputs string) (*models.Task, error) {

	d.logger.Debug("create task",
		zap.String("judgement id", judgementId),
		zap.String("task type", taskType),
		zap.String("properties", properties),
	)

	task := &models.Task{
		JudgementId: judgementId,
		TaskId:      uuid.New().String(),
		Type:        taskType,
		Properties:  properties,
		Inputs:      inputs,
		Outputs:     "",
	}

	d.Repository.PushTaskInQueue(blockId, task)

	return task, nil
}

func (d DefaultJudgementsService) GetTasks(taskType string) (tasks []*models.Task, err error) {
	d.logger.Info("get task", zap.String("type", taskType))
	element := d.Repository.FetchTaskInQueue("*", "*", taskType)
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

func (d DefaultJudgementsService) GetTask(taskId string) (task *models.Task, err error) {
	d.logger.Info("get task",
		zap.String("task id", taskId),
	)
	element := d.Repository.FetchTaskInQueue("*", taskId, "*")
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

func (d DefaultJudgementsService) UpdateTask(token, taskId, outputsStr string) (task *models.Task, err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	taskId, ok := d.tokenMap[token]

	if !ok {
		d.logger.Debug("invalid token: no such token")
		return nil, errors.New("invalid token")
	}

	// token should only be used once
	delete(d.tokenMap, token)

	taskElement := d.Repository.FetchTaskInQueue("*", taskId, "*")
	if taskElement == nil {
		d.logger.Debug("invalid token: no such task",
			zap.String("task id", taskId),
		)
		d.Repository.UnlockTaskInQueue(taskElement)
		return nil, errors.New("invalid token")
	}

	task = taskElement.Task

	if task.TaskId != taskId {
		d.logger.Debug("task mismatch",
			zap.String("expected task id", task.TaskId),
			zap.String("actual task id", taskId),
		)
		d.Repository.UnlockTaskInQueue(taskElement)
		return nil, errors.New("task mismatch")
	}

	d.logger.Info("update task",
		zap.String("task id", taskId),
	)

	// TODO: update task
	blockId := taskElement.BlockId
	processRuntime, ok := d.processMap[taskElement.JudgementId]
	if !ok {
		d.logger.Error("missing process")
		d.Repository.UnlockTaskInQueue(taskElement)
		return nil, errors.New("missing process")
	}

	outputs, err := crypto.EasyDecode(outputsStr)
	if err != nil {
		d.logger.Debug("wrong output format")
		d.Repository.UnlockTaskInQueue(taskElement)
		return nil, errors.New("wrong output format")
	}
	err = processRuntime.SetOutputs(blockId, outputs)
	if err != nil {
		d.logger.Error("update task: set outputs failed", zap.Error(err))
		d.Repository.UnlockTaskInQueue(taskElement)
		return nil, errors.New("set outputs failed")
	}

	// update task
	//err := d.Repository.Update(element, outputs)
	//if err != nil {
	//	d.logger.Error("update task", zap.Error(err))
	//	return nil, err
	//}

	// remove from queue
	d.Repository.RemoveTaskInQueue(taskElement)

	// calculate next task
	err = d.PushProcessRuntime(processRuntime)
	if err != nil {
		d.logger.Error("create judgement: initial process failed",
			zap.String("task id", taskId),
			zap.Error(err),
		)
	}

	return task, nil
}

func (d DefaultJudgementsService) ReserveTask(taskId string) (token string, err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	taskElement := d.Repository.FetchTaskInQueue("*", taskId, "*")

	if !d.Repository.LockTaskInQueue(taskElement) {
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

func NewJudgementsService(
	logger *zap.Logger,
	Repository repositories.Repository,
	ProcessRepository processRepository.Repository,
	SubmissionRepository submissionRepository.Repository,
) JudgementsService {
	return &DefaultJudgementsService{
		mutex:                &sync.Mutex{},
		logger:               logger.With(zap.String("type", "DefaultJudgementService")),
		Repository:           Repository,
		processRepository:    ProcessRepository,
		submissionRepository: SubmissionRepository,

		tokenMap:   make(map[string]string),
		processMap: make(map[string]*ProcessRuntime),
		queue:      list.New(),
	}
}
