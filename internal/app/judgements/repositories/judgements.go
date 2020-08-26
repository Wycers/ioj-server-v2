package repositories

import (
	"container/list"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/infinity-oj/server-v2/internal/pkg/models"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type TaskElement struct {
	IsLocked bool

	JudgementId string
	BlockId     int

	TaskId string
	Type   string

	Task *models.Task
}

type Repository interface {
	List()
	Fetch() *models.Judgement
	Create(submissionId, processId uint64) (*models.Judgement, error)
	Update(judgement *models.Judgement) error

	PushTaskInQueue(blockId int, task *models.Task)
	FetchTaskInQueue(judgementId, taskId, taskType string) *TaskElement
	RemoveTaskInQueue(element *TaskElement)
	LockTaskInQueue(element *TaskElement) bool
	UnlockTaskInQueue(element *TaskElement) bool
}

type DefaultRepository struct {
	logger *zap.Logger
	db     *gorm.DB
	mutex  *sync.Mutex
	queue  *list.List // tasks list
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

func (m DefaultRepository) List() {
	fmt.Println("=== START ===")

	for te := m.queue.Front(); te != nil; te = te.Next() {
		element, ok := te.Value.(*TaskElement)

		if !ok {
			fmt.Println(te.Value)
			continue
		}

		fmt.Printf("judgement id: %s\ntask id:%s\ntype: %s\n locked: %t\n\n",
			element.JudgementId,
			element.TaskId,
			element.Type,
			element.IsLocked,
		)
	}

	fmt.Println("==== END ====")
}

// FetchJudgementInQueue returns task with specific task type.
func (m DefaultRepository) FetchTaskInQueue(judgementId, taskId, taskType string) *TaskElement {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for te := m.queue.Front(); te != nil; te = te.Next() {
		taskElement, ok := te.Value.(*TaskElement)

		if !ok {
			fmt.Println(te.Value)
			panic("something wrong")
		}

		if judgementId != "*" && taskElement.JudgementId != judgementId {
			continue
		}

		if taskId != "*" && taskElement.TaskId != taskId {
			continue
		}

		if taskType != "*" && taskElement.Type != taskType {
			continue
		}

		return taskElement
	}

	return nil
}

func (m DefaultRepository) LockTaskInQueue(element *TaskElement) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.logger.Debug("judgement repository, unlock task",
		zap.String("judgement id", element.JudgementId),
		zap.String("task id", element.TaskId),
	)

	if element.IsLocked {
		m.logger.Error("lock a task that is locked",
			zap.String("judgement id", element.JudgementId),
			zap.String("task id", element.TaskId),
		)
		return false
	}
	element.IsLocked = true
	return true
}

func (m DefaultRepository) UnlockTaskInQueue(element *TaskElement) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.logger.Debug("judgement repository, unlock task",
		zap.String("judgement id", element.JudgementId),
		zap.String("task id", element.TaskId),
	)

	if element.IsLocked == false {
		m.logger.Error("unlock a task that is not locked",
			zap.String("judgement id", element.JudgementId),
			zap.String("task id", element.TaskId),
		)
		return false
	}
	element.IsLocked = false
	return true
}

func (m DefaultRepository) PushTaskInQueue(blockId int, task *models.Task) {
	m.mutex.Lock()
	m.mutex.Unlock()

	m.logger.Debug("push task in queue",
		zap.String("task id", task.TaskId),
		zap.String("task type", task.Type),
	)

	element := &TaskElement{
		IsLocked:    false,
		JudgementId: task.JudgementId,
		BlockId:     blockId,
		TaskId:      task.TaskId,
		Type:        task.Type,
		Task:        task,
	}

	m.queue.PushBack(element)
}

func (m DefaultRepository) RemoveTaskInQueue(element *TaskElement) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.logger.Debug("remove task in queue",
		zap.String("task id", element.TaskId),
	)

	for te := m.queue.Front(); te != nil; te = te.Next() {
		je, ok := te.Value.(*TaskElement)

		if !ok {
			continue
		}

		if je == element {
			m.queue.Remove(te)
			break
		}
	}
}

func NewJudgementRepository(logger *zap.Logger, db *gorm.DB) Repository {
	return &DefaultRepository{
		logger: logger.With(zap.String("type", "Repository")),
		db:     db,
		mutex:  &sync.Mutex{},
		queue:  list.New(),
	}
}
