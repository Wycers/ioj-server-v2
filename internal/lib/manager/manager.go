package manager

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/infinity-oj/server-v2/internal/lib/engine"

	"github.com/google/uuid"
	"github.com/google/wire"

	"github.com/infinity-oj/server-v2/pkg/models"
	"go.uber.org/zap"
)

type ProcessRuntime struct {
	isLocked bool
	lockedAt time.Time
	c        chan *models.Slots
	block    *engine.Block

	Mutex     *sync.Mutex
	Judgement *models.Judgement
	Process   *models.Process
}

type ProcessManager interface {
	Push(judgement *models.Judgement, block *engine.Block, inputs *models.Slots) <-chan *models.Slots
	Fetch(judgementId, processId, processType string, ignoreLock bool) *ProcessRuntime
	Finish(element *ProcessRuntime, slots *models.Slots) error
	FinishWithError(element *ProcessRuntime, message string) error
	Reserve(element *ProcessRuntime) bool
}

type manager struct {
	logger    *zap.Logger
	mutex     *sync.Mutex
	processes *list.List

	buildIns []Handler
}

func (m *manager) Reserve(element *ProcessRuntime) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if element.isLocked {
		return false
	}
	element.isLocked = true
	return true
}

func (m *manager) List() {
	fmt.Println("=== START ===")

	for te := m.processes.Front(); te != nil; te = te.Next() {
		element, ok := te.Value.(*ProcessRuntime)

		if !ok {
			fmt.Println(te.Value)
			continue
		}

		fmt.Printf("judgement id: %s\nprocess id:%s\ntype: %s\nlocked: %t\n\n",
			element.Process.JudgementId,
			element.Process.ProcessId,
			element.Process.Type,
			element.isLocked,
		)
	}

	fmt.Println("==== END ====")
}

type Handler interface {
	IsMatched(tp string) bool
	Work(runtime *ProcessRuntime) error
}

func (m *manager) Push(judgement *models.Judgement, block *engine.Block, inputs *models.Slots) (c <-chan *models.Slots) {
	process := &models.Process{
		Model: models.Model{
			ID:        0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
		Type:        block.Type,
		ProcessId:   uuid.New().String(),
		JudgementId: judgement.Name,
		Properties:  block.Properties,
		Inputs:      *inputs,
		Outputs:     models.Slots{},
	}

	m.logger.Debug("push process in processes",
		zap.String("process id", process.ProcessId),
		zap.String("process type", process.Type),
	)
	runtime := &ProcessRuntime{
		isLocked: false,
		c:        make(chan *models.Slots, 1),
		block:    block,

		Mutex:     &sync.Mutex{},
		Judgement: judgement,
		Process:   process,
	}
	c = runtime.c

	fmt.Println("new runtime", runtime.Process)
	m.logger.Debug("consume runtime",
		zap.String("process id", runtime.Process.ProcessId),
		zap.String("process type", runtime.Process.Type),
	)
	for _, b := range m.buildIns {
		if b.IsMatched(process.Type) {
			if err := b.Work(runtime); err != nil {
				m.logger.Error("consume", zap.Error(err))
			}
			if err := m.Finish(runtime, &runtime.Process.Outputs); err != nil {
				m.logger.Error("finish", zap.Error(err))
			}
			return
		}
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.processes.PushBack(runtime)
	return
}

// Fetch returns process with specific process type.
func (m *manager) Fetch(judgementId, processId, processType string, ignoreLock bool) *ProcessRuntime {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.logger.Debug("working!", zap.String("process id", processId))

	for te := m.processes.Front(); te != nil; te = te.Next() {
		processElement, ok := te.Value.(*ProcessRuntime)

		if !ok {
			panic("internal error")
		}

		if processElement.isLocked && time.Now().Sub(processElement.lockedAt) > 1000*time.Second {
			processElement.isLocked = false
		}

		if judgementId != "*" && processElement.Process.JudgementId != judgementId {
			continue
		}

		if processId != "*" && processElement.Process.ProcessId != processId {
			continue
		}

		if processType != "*" && processElement.Process.Type != processType {
			continue
		}

		if processElement.isLocked {
			if !ignoreLock {
				continue
			}
		}

		return processElement
	}

	return nil
}

func (m *manager) Finish(element *ProcessRuntime, outputs *models.Slots) error {
	m.logger.Debug("finish process ",
		zap.String("process id", element.Process.ProcessId),
		zap.String("process type", element.Process.Type),
	)
	m.remove(element)
	element.c <- outputs
	return nil
}

func (m *manager) FinishWithError(element *ProcessRuntime, message string) error {
	//element.runtime.Judgement.Msg = message
	//element.runtime.Judgement.Status = models.SystemError
	//element.runtime.Judgement.Score = 0
	//
	//m.FinishRuntime(element.runtime)
	return nil
}

func (m *manager) remove(element *ProcessRuntime) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.logger.Debug("remove process in processes",
		zap.String("process id", element.Process.ProcessId),
	)

	for te := m.processes.Front(); te != nil; te = te.Next() {
		je, ok := te.Value.(*ProcessRuntime)

		if !ok {
			continue
		}

		if je == element {
			m.processes.Remove(te)
			break
		}
	}
}

var instance *manager
var once *sync.Once

func GetManager() ProcessManager {
	if instance == nil {
		panic("manager is nil")
	}
	return instance
}

func NewManager(logger *zap.Logger, ins []Handler) ProcessManager {
	once = &sync.Once{}
	once.Do(func() {
		instance = &manager{
			logger:    logger,
			mutex:     &sync.Mutex{},
			processes: list.New(),

			buildIns: ins,
		}
	})
	return instance
}

var ProviderSet = wire.NewSet(NewManager)
