package schedulers

import (
	"container/list"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/infinity-oj/server-v2/internal/pkg/eventBus"

	"github.com/google/wire"

	"github.com/google/uuid"

	"github.com/infinity-oj/server-v2/internal/lib/nodeEngine"

	"go.uber.org/zap"

	"github.com/infinity-oj/server-v2/pkg/models"
)

type Scheduler interface {
	List()

	FinishedJudgement() chan *models.Judgement

	NewRuntime(problem *models.Problem, submission *models.Submission, judgement *models.Judgement, process *models.Process) (*Runtime, error)
	PushRuntime(runtime *Runtime) error
	FinishRuntime(runtime *Runtime)

	FetchTask(judgementId, taskId, taskType string, ignoreLock bool) *TaskElement
	FinishTask(element *TaskElement, slots *models.Slots) error
	FinishTaskWithError(element *TaskElement, message string) error
	LockTask(element *TaskElement) bool
	UnlockTask(element *TaskElement) bool
}

var pendingTasks chan *TaskElement

type Runtime struct {
	Problem    *models.Problem
	Submission *models.Submission
	Judgement  *models.Judgement

	graph  *nodeEngine.Graph
	result map[int]*models.Slot
}

type scheduler struct {
	logger *zap.Logger
	mutex  *sync.Mutex

	tasks    *list.List
	runtimes map[*models.Judgement]*Runtime

	finished chan *models.Judgement

	eventBus eventBus.Bus
}

func (s *scheduler) FinishRuntime(runtime *Runtime) {
	s.logger.Debug("finish runtime",
		zap.String("judgement id", runtime.Judgement.JudgementId),
	)

	s.finished <- runtime.Judgement

	delete(s.runtimes, runtime.Judgement)

	s.mutex.Lock()
	defer s.mutex.Unlock()
	for te := s.tasks.Front(); te != nil; te = te.Next() {
		if element, ok := te.Value.(*TaskElement); ok && element.runtime == runtime {
			s.tasks.Remove(te)
		}
	}
}

func (s *scheduler) PushRuntime(runtime *Runtime) error {
	judgement := runtime.Judgement

	if _, ok := s.runtimes[judgement]; ok {
		return errors.New("processing")
	}
	s.runtimes[judgement] = runtime

	return forward(runtime)
}

func (s *scheduler) FinishedJudgement() chan *models.Judgement {
	return s.finished
}

func (s *scheduler) release() {
	for te := s.tasks.Front(); te != nil; te = te.Next() {
		element, ok := te.Value.(*TaskElement)

		if !ok {
			fmt.Println(te.Value)
			continue
		}

		if !element.IsLocked {
			continue
		}

		fmt.Println(time.Now().Sub(element.LockedAt))

		if time.Now().Sub(element.LockedAt) > 3*time.Second {
			s.UnlockTask(element)
		}
	}
}

func (s *scheduler) releaseTimer() {
	for _ = range time.Tick(3 * time.Second) {
		s.release()
	}
	panic("should not be here")
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
			element.runtime.Judgement.JudgementId,
			element.Task.TaskId,
			element.Task.Type,
			element.IsLocked,
		)
	}

	fmt.Println("==== END ====")
}

type ProcessElement struct {
	runtime     *Runtime
	blockId     int
	taskElement *TaskElement
}

// NewRuntime create new process runtime information with judgement and process
func (s scheduler) NewRuntime(
	problem *models.Problem,
	submission *models.Submission,
	judgement *models.Judgement,
	process *models.Process,
) (*Runtime, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	submissionId := judgement.SubmissionId
	processId := process.ID

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
		return nil, err
	}

	result := make(map[int]*models.Slot)
	runtime := &Runtime{
		Problem:    problem,
		Submission: submission,
		Judgement:  judgement,
		graph:      graph,
		result:     result,
	}

	return runtime, nil
}

type TaskElement struct {
	IsLocked bool
	LockedAt time.Time

	BlockId int
	Task    *models.Task
	runtime *Runtime
}

func (s scheduler) PushTask(blockId int, task *models.Task, runtime *Runtime) {
	s.logger.Debug("push task in tasks",
		zap.String("task id", task.TaskId),
		zap.String("task type", task.Type),
	)

	element := &TaskElement{
		IsLocked: false,

		Task:    task,
		BlockId: blockId,
		runtime: runtime,
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

		if judgementId != "*" && taskElement.runtime.Judgement.JudgementId != judgementId {
			continue
		}

		if taskId != "*" && taskElement.Task.TaskId != taskId {
			continue
		}

		if taskType != "*" && taskElement.Task.Type != taskType {
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

func forward(runtime *Runtime) error {

	ids := runtime.graph.Run()

	for _, block := range ids {
		var inputs models.Slots
		for _, linkId := range block.Inputs {
			if data, ok := runtime.result[linkId]; ok {
				inputs = append(inputs, data)
			} else {
				return errors.New("wrong process definition")
			}
		}

		newTask := &models.Task{
			Model: models.Model{
				ID:        0,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: nil,
			},
			Type:        block.Type,
			TaskId:      uuid.New().String(),
			JudgementId: runtime.Judgement.JudgementId,
			Properties:  block.Properties,
			Inputs:      inputs,
			Outputs:     models.Slots{},
		}

		s.PushTask(block.Id, newTask, runtime)
	}

	return nil
}

func (s scheduler) FinishTask(element *TaskElement, outputs *models.Slots) error {
	blockId := element.BlockId
	runtime := element.runtime
	block := runtime.graph.FindBlockById(blockId)

	if len(block.Output) != len(*outputs) {
		msg := fmt.Sprintf("output slots mismatch, block %d expects %d but %d",
			blockId,
			len(block.Output),
			len(*outputs),
		)
		return errors.New(msg)
	}

	for index, output := range *outputs {
		links := runtime.graph.FindLinkBySourcePort(blockId, index)
		for _, link := range links {
			runtime.result[link.Id] = output
		}
	}

	block.Done()

	s.RemoveTask(element)

	return forward(runtime)
}

func (s *scheduler) FinishTaskWithError(element *TaskElement, message string) error {
	element.runtime.Judgement.Msg = message
	element.runtime.Judgement.Status = models.SystemError
	element.runtime.Judgement.Score = 0

	s.FinishRuntime(element.runtime)

	return nil
}

func (s *scheduler) RemoveTask(element *TaskElement) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.logger.Debug("remove task in tasks",
		zap.String("task id", element.Task.TaskId),
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

func (s *scheduler) LockTask(element *TaskElement) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.logger.Debug("judgement repository, unlock task",
		zap.String("judgement id", element.runtime.Judgement.JudgementId),
		zap.String("task id", element.Task.TaskId),
	)

	if element.IsLocked {
		s.logger.Error("lock a task that is locked",
			zap.String("judgement id", element.runtime.Judgement.JudgementId),
			zap.String("task id", element.Task.TaskId),
		)
		return false
	}
	element.IsLocked = true
	element.LockedAt = time.Now()
	return true
}

func (s *scheduler) UnlockTask(element *TaskElement) bool {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.logger.Debug("judgement repository, unlock task",
		zap.String("judgement id", element.runtime.Judgement.JudgementId),
		zap.String("task id", element.Task.TaskId),
	)

	if element.IsLocked == false {
		s.logger.Error("unlock a task that is not locked",
			zap.String("judgement id", element.runtime.Judgement.JudgementId),
			zap.String("task id", element.Task.TaskId),
		)
		return false
	}
	element.IsLocked = false
	return true
}

func (s *scheduler) consume() {
	funcs := []func(element *TaskElement) (bool, error){File, String, Evaluate}

	for {
		select {
		case element := <-pendingTasks:

			if element.Task.Type == "basic/end" {
				if score, ok := element.Task.Inputs[0].Value.(float64); !ok {
					element.runtime.Judgement.Msg = "wrong score"
					element.runtime.Judgement.Status = models.SystemError
					element.runtime.Judgement.Score = 0
				} else {
					element.runtime.Judgement.Msg = ""
					element.runtime.Judgement.Status = models.Accepted
					element.runtime.Judgement.Score = score
				}
				s.FinishRuntime(element.runtime)
				continue
			}

			fmt.Println("element", element)
			matched := false
			for _, f := range funcs {
				var err error
				matched, err = f(element)
				fmt.Println(err)
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
			fmt.Println("publish!", &s.eventBus)
			s.eventBus.Publish("task:new", element.Task)
		}
	}
}

var s *scheduler
var once sync.Once

func New(logger *zap.Logger) Scheduler {
	pendingTasks = make(chan *TaskElement, 128)
	bus := eventBus.New()

	once.Do(func() {
		s = &scheduler{
			logger:   logger.With(zap.String("type", "scheduler")),
			mutex:    &sync.Mutex{},
			tasks:    &list.List{},
			runtimes: make(map[*models.Judgement]*Runtime),
			finished: make(chan *models.Judgement, 64),
			eventBus: bus,
		}
		go s.releaseTimer()
	})

	go s.consume()

	return s
}

var ProviderSet = wire.NewSet(New)
