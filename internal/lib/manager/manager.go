package manager

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/google/wire"

	"github.com/infinity-oj/server-v2/pkg/models"
	"go.uber.org/zap"
)

type ProcessElement struct {
	IsLocked bool
	LockedAt time.Time

	BlockId int
	Process *models.Process
	C       chan *models.Slots
}

type ProcessManager interface {
	List()
	Push(blockId int, process *models.Process) <-chan *models.Slots
	Fetch(judgementId, processId, processType string, ignoreLock bool) *ProcessElement
	Finish(element *ProcessElement, slots *models.Slots) error
	FinishWithError(element *ProcessElement, message string) error
	Reserve(element *ProcessElement) bool
}
type manager struct {
	logger    *zap.Logger
	mutex     *sync.Mutex
	processes *list.List

	buildIns []Handler
}

func (m *manager) Reserve(element *ProcessElement) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if element.IsLocked {
		return false
	}
	element.IsLocked = true
	return true
}

func (m *manager) List() {
	fmt.Println("=== START ===")

	for te := m.processes.Front(); te != nil; te = te.Next() {
		element, ok := te.Value.(*ProcessElement)

		if !ok {
			fmt.Println(te.Value)
			continue
		}

		fmt.Printf("judgement id: %s\nprocess id:%s\ntype: %s\nlocked: %t\n\n",
			element.Process.JudgementId,
			element.Process.ProcessId,
			element.Process.Type,
			element.IsLocked,
		)
	}

	fmt.Println("==== END ====")
}

type Handler interface {
	IsMatched(tp string) bool
	Work(process *models.Process) error
}

func (m *manager) Push(blockId int, process *models.Process) (c <-chan *models.Slots) {
	m.logger.Debug("push process in processes",
		zap.String("process id", process.ProcessId),
		zap.String("process type", process.Type),
	)
	element := &ProcessElement{
		IsLocked: false,

		BlockId: blockId,
		Process: process,
		C:       make(chan *models.Slots, 1),
	}
	c = element.C

	fmt.Println("new element", element.Process)
	m.logger.Debug("consume element",
		zap.String("process id", element.Process.ProcessId),
		zap.String("process type", element.Process.Type),
	)
	for _, b := range m.buildIns {
		if b.IsMatched(process.Type) {
			if err := b.Work(process); err != nil {
				m.logger.Error("consume", zap.Error(err))
			}
			err := m.Finish(element, &element.Process.Outputs)
			if err != nil {
				m.logger.Error("consume", zap.Error(err))
			}
			return
		}
	}
	m.processes.PushBack(element)
	return
}

// Fetch returns process with specific process type.
func (m *manager) Fetch(judgementId, processId, processType string, ignoreLock bool) *ProcessElement {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for te := m.processes.Front(); te != nil; te = te.Next() {
		processElement, ok := te.Value.(*ProcessElement)

		if !ok {
			panic("internal error")
		}

		if processElement.IsLocked && time.Now().Sub(processElement.LockedAt) > 3*time.Second {
			processElement.IsLocked = false
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

		if processElement.IsLocked {
			if !ignoreLock {
				continue
			}
		}

		return processElement
	}

	return nil
}

func (m *manager) Finish(element *ProcessElement, outputs *models.Slots) error {
	m.logger.Debug("finish process ",
		zap.String("process id", element.Process.ProcessId),
		zap.String("process type", element.Process.Type),
	)

	m.remove(element)
	element.C <- outputs
	return nil
}

func (m *manager) FinishWithError(element *ProcessElement, message string) error {
	//element.runtime.Judgement.Msg = message
	//element.runtime.Judgement.Status = models.SystemError
	//element.runtime.Judgement.Score = 0
	//
	//m.FinishRuntime(element.runtime)
	return nil
}

func (m *manager) remove(element *ProcessElement) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.logger.Debug("remove process in processes",
		zap.String("process id", element.Process.ProcessId),
	)

	for te := m.processes.Front(); te != nil; te = te.Next() {
		je, ok := te.Value.(*ProcessElement)

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
