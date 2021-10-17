package processes

import (
	"container/list"
	"fmt"
	"sync"
	"time"

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

	pendingProcesses chan *ProcessElement
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

func (m *manager) Push(blockId int, process *models.Process) <-chan *models.Slots {
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

	m.pendingProcesses <- element
	return element.C
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

func (m *manager) consume() {
	for {
		select {
		case element := <-m.pendingProcesses:

			//if element.Process.Type == "basic/end" {
			//	if score, ok := element.Process.Inputs[0].Value.(float64); !ok {
			//		element.runtime.Judgement.Msg = "wrong score"
			//		element.runtime.Judgement.Status = models.SystemError
			//		element.runtime.Judgement.Score = 0
			//	} else {
			//		element.runtime.Judgement.Msg = ""
			//		element.runtime.Judgement.Status = models.Accepted
			//		element.runtime.Judgement.Score = score
			//	}
			//	m.FinishRuntime(element.runtime)
			//	continue
			//}

			fmt.Println("element", element)
			matched := false
			for _, f := range []func(element *ProcessElement) (bool, error){File, String, Evaluate} {
				var err error
				matched, err = f(element)
				fmt.Println(err)
				if matched {
					break
				}
			}
			if matched {
				err := m.Finish(element, &element.Process.Outputs)
				if err != nil {
					fmt.Println(err)
					continue
				}
				continue
			}

			m.processes.PushBack(element)
		}
	}
}

var instance ProcessManager
var once *sync.Once

func GetManager() ProcessManager {
	if instance == nil {
		panic("manage is nil")
	}
	return instance
}

func NewManager(logger *zap.Logger) ProcessManager {
	once = &sync.Once{}
	once.Do(func() {
		instance = &manager{
			logger:           logger,
			mutex:            &sync.Mutex{},
			processes:        list.New(),
			pendingProcesses: make(chan *ProcessElement, 128),
		}
	})
	return instance
}
