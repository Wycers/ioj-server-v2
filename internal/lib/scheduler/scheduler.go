package scheduler

import (
	"container/list"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"

	"github.com/infinity-oj/server-v2/internal/lib/nodeEngine"

	"go.uber.org/zap"

	"github.com/infinity-oj/server-v2/pkg/models"
)

type Scheduler interface {
	List()

	NewProcessRuntime(problem *models.Problem, submission *models.Submission, judgement *models.Judgement, process *models.Process) error

	PushTask(blockId int, task *models.Task)
	FetchTask(judgementId, taskId, taskType string, ignoreLock bool) *TaskElement
	FinishTask(element *TaskElement, slots *models.Slots) error
	RemoveTask(element *TaskElement)
	LockTask(element *TaskElement) bool
	UnlockTask(element *TaskElement) bool
}

var pendingTasks chan *TaskElement

type processRuntime struct {
	problem    *models.Problem
	submission *models.Submission
	judgement  *models.Judgement

	graph  *nodeEngine.Graph
	result map[int]*models.Slot
}

type scheduler struct {
	logger *zap.Logger
	mutex  *sync.Mutex

	tasks     *list.List
	processes map[string]*processRuntime
}

func (s scheduler) List() {
	fmt.Println("=== START ===")

	for te := s.tasks.Front(); te != nil; te = te.Next() {
		element, ok := te.Value.(*TaskElement)

		if !ok {
			fmt.Println(te.Value)
			continue
		}

		fmt.Printf("judgement id: %s\ntask id:%s\ntype: %s\nlocked: %t\n\n",
			element.JudgementId,
			element.TaskId,
			element.Type,
			element.IsLocked,
		)
	}

	fmt.Println("==== END ====")
}

type ProcessElement struct {
	processRuntime *processRuntime
	blockId        int
	taskElement    *TaskElement
}

// NewProcessRuntime create new process runtime information with judgement and process
func (s scheduler) NewProcessRuntime(problem *models.Problem, submission *models.Submission, judgement *models.Judgement, process *models.Process) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	submissionId := judgement.SubmissionId
	processId := process.ID
	judgementId := judgement.JudgementId

	definition := process.Definition
	definition = strings.ReplaceAll(definition, "<userVolume>", submission.UserVolume)
	definition = strings.ReplaceAll(definition, "<publicVolume>", problem.PublicVolume)
	definition = strings.ReplaceAll(definition, "<privateVolume>", problem.PrivateVolume)
	graph, err := nodeEngine.NewGraphByDefinition(definition)
	if err != nil {
		s.logger.Error("parse process definition failed",
			zap.Uint64("submission id", submissionId),
			zap.Uint64("process id", processId),
			zap.Error(err),
		)
		return nil
	}

	result := make(map[int]*models.Slot)
	pr := &processRuntime{
		submission: submission,
		judgement:  judgement,
		graph:      graph,
		result:     result,
	}

	s.processes[judgementId] = pr

	return forward(pr)
}

type TaskElement struct {
	IsLocked bool

	JudgementId string
	BlockId     int

	TaskId string
	Type   string

	Task *models.Task
}

func (s scheduler) PushTask(blockId int, task *models.Task) {
	s.logger.Debug("push task in tasks",
		zap.String("task id", task.TaskId),
		zap.String("task type", task.Type),
	)

	element := &TaskElement{
		Task: task,

		IsLocked:    false,
		JudgementId: task.JudgementId,
		BlockId:     blockId,
		TaskId:      task.TaskId,
		Type:        task.Type,
	}

	pendingTasks <- element
}

// FetchTask returns task with specific task type.
func (s scheduler) FetchTask(judgementId, taskId, taskType string, ignoreLock bool) *TaskElement {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for te := s.tasks.Front(); te != nil; te = te.Next() {
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

		if taskElement.IsLocked {
			if !ignoreLock {
				continue
			}
		}

		return taskElement
	}

	return nil
}

func forward(pr *processRuntime) error {

	ids := pr.graph.Run()

	for _, block := range ids {
		var inputs models.Slots
		for _, linkId := range block.Inputs {
			if data, ok := pr.result[linkId]; ok {
				inputs = append(inputs, data)
			} else {
				return errors.New("wrong process definition")
			}
		}

		newTask := &models.Task{
			JudgementId: pr.judgement.JudgementId,
			TaskId:      uuid.New().String(),
			Type:        block.Type,
			Properties:  block.Properties,
			Inputs:      inputs,
			Outputs:     models.Slots{},
		}

		s.PushTask(block.Id, newTask)
	}

	return nil
}

func (s scheduler) FinishTask(element *TaskElement, outputs *models.Slots) error {
	blockId := element.BlockId
	pr, ok := s.processes[element.JudgementId]
	if !ok {
		s.logger.Error("missing process")
		s.UnlockTask(element)
		return errors.New("missing process")
	}
	block := pr.graph.FindBlockById(blockId)

	if len(block.Output) != len(*outputs) {
		msg := fmt.Sprintf("output slots mismatch, block %d expects %d but %d",
			blockId,
			len(block.Output),
			len(*outputs),
		)
		return errors.New(msg)
	}

	for index, output := range *outputs {
		links := pr.graph.FindLinkBySourcePort(blockId, index)
		for _, link := range links {
			pr.result[link.Id] = output
		}
	}

	block.Done()

	s.RemoveTask(element)

	return forward(pr)
}
func (s scheduler) RemoveTask(element *TaskElement) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.logger.Debug("remove task in tasks",
		zap.String("task id", element.TaskId),
	)

	for te := s.tasks.Front(); te != nil; te = te.Next() {
		je, ok := te.Value.(*TaskElement)

		if !ok {
			continue
		}

		if je == element {
			s.tasks.Remove(te)
			break
		}
	}
}

func (s scheduler) LockTask(element *TaskElement) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.logger.Debug("judgement repository, unlock task",
		zap.String("judgement id", element.JudgementId),
		zap.String("task id", element.TaskId),
	)

	if element.IsLocked {
		s.logger.Error("lock a task that is locked",
			zap.String("judgement id", element.JudgementId),
			zap.String("task id", element.TaskId),
		)
		return false
	}
	element.IsLocked = true
	return true
}

func (s scheduler) UnlockTask(element *TaskElement) bool {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.logger.Debug("judgement repository, unlock task",
		zap.String("judgement id", element.JudgementId),
		zap.String("task id", element.TaskId),
	)

	if element.IsLocked == false {
		s.logger.Error("unlock a task that is not locked",
			zap.String("judgement id", element.JudgementId),
			zap.String("task id", element.TaskId),
		)
		return false
	}
	element.IsLocked = false
	return true
}

var s *scheduler
var once sync.Once

func New(logger *zap.Logger) Scheduler {
	pendingTasks = make(chan *TaskElement, 128)

	funcs := []func(element *TaskElement) (bool, error){File}

	once.Do(func() {
		s = &scheduler{
			logger.With(zap.String("type", "scheduler")),
			&sync.Mutex{},
			&list.List{},
			make(map[string]*processRuntime),
		}

		go func() {
			for {
				element := <-pendingTasks
				matched := false
				for _, f := range funcs {
					matched, _ = f(element)
					if matched {
						break
					}
				}
				if matched {
					err := s.FinishTask(element, &element.Task.Outputs)
					if err != nil {
						fmt.Println(err)
						continue
					}
					continue
				}

				s.tasks.PushBack(element)
			}
		}()
	})

	return s
}
