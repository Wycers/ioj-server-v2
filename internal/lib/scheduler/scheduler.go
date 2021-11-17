package scheduler

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/infinity-oj/server-v2/internal/lib/manager"

	"github.com/infinity-oj/server-v2/internal/lib/engine/scene"

	"github.com/google/wire"

	"github.com/infinity-oj/server-v2/internal/lib/engine"

	"go.uber.org/zap"

	"github.com/infinity-oj/server-v2/pkg/models"
)

type Runtime struct {
	Blueprint  *models.Blueprint
	Problem    *models.Problem
	Submission *models.Submission
	Judgement  *models.Judgement

	graph  *engine.Graph
	result map[int]*models.Slot
}

type Scheduler struct {
	logger *zap.Logger
	mutex  *sync.Mutex

	Runtime *Runtime

	C chan int
}

func (s *Scheduler) Execute() {
	code := 0
	defer func() {
		s.C <- code
	}()
	s.logger.Debug("scheduler: execution started")

	lock := &sync.RWMutex{}

	wg := new(sync.WaitGroup)
	trigger := make(chan int32, 100)
	var n int32 = 0
	trigger <- atomic.AddInt32(&n, 1)
	for _ = range trigger {
		ids := s.Runtime.graph.Run()
		for _, block := range ids {
			var inputs models.Slots
			lock.RLock()
			for _, linkId := range block.Inputs {
				if data, ok := s.Runtime.result[linkId]; ok {
					inputs = append(inputs, data)
				} else {
					s.logger.Error("wrong process definition", zap.Int("link id", linkId))
					code = -1
					return
				}
			}
			lock.RUnlock()

			atomic.AddInt32(&n, 1)
			wg.Add(1)
			go func(block *engine.Block) {
				blockId := block.Id
				s.logger.Debug("process started", zap.Int("block id", blockId), zap.Any("inputs", inputs))

				select {
				case outputs := <-manager.Push(s.Runtime.Judgement, block, &inputs):
					s.logger.Debug("process finished normally", zap.Int("block id", blockId), zap.Any("outputs", outputs))

					if len(block.Output) != len(*outputs) {
						s.logger.Error(fmt.Sprintf("output slots mismatch, block %d expects %d but %d",
							block.Id,
							len(block.Output),
							len(*outputs),
						))
						return
					}

					lock.Lock()
					for index, output := range *outputs {
						links := s.Runtime.graph.FindLinkBySourcePort(blockId, index)
						for _, link := range links {
							s.Runtime.result[link.Id] = output
						}
					}
					lock.Unlock()
					block.Done()
					trigger <- atomic.AddInt32(&n, 1)
				case <-time.After(time.Second * 500):
					s.logger.Debug("process timeout after 500s", zap.Int("block id", blockId))
					// 其实这个时候应该是把评测挂起比较好
				}

				s.logger.Debug("process ended", zap.Int("block id", blockId))
				if atomic.AddInt32(&n, -1) == 0 {
					s.logger.Debug("pending count is 0, closing")
					close(trigger)
				}
				wg.Done()
			}(block)
		}
		if atomic.AddInt32(&n, -1) == 0 {
			s.logger.Debug("pending count is 0, closing")
			close(trigger)
		}
	}
	s.logger.Debug("scheduler: execution ended")
	wg.Wait()
}

func (s *Scheduler) OnFinish() <-chan int {
	return s.C
}

func New(logger *zap.Logger,
	problem *models.Problem, submission *models.Submission, judgement *models.Judgement,
	blueprint *models.Blueprint, programs []*models.Program,
) (*Scheduler, error) {
	blueprintId := blueprint.ID

	definition := blueprint.Definition
	// TODO: throw the mass
	if submission != nil {
		definition = strings.ReplaceAll(definition, "${userVolume}", submission.UserVolume)
		definition = strings.ReplaceAll(definition, "${account_id}", fmt.Sprintf("%d", submission.SubmitterId))
	}
	if problem != nil {
		definition = strings.ReplaceAll(definition, "${publicVolume}", problem.PublicVolume)
		definition = strings.ReplaceAll(definition, "${privateVolume}", problem.PrivateVolume)
		definition = strings.ReplaceAll(definition, "${problem_id}", problem.Name)
	}
	s := scene.NewScene(definition)
	//graph, err := engine.NewGraphByDefinition(definition)
	var bs []*scene.BlockDefinition
	for _, p := range programs {
		bs = append(bs, scene.NewBlockDefinition(p.Definition))
	}
	graph, err := engine.NewGraphByScene(bs, s)
	if err != nil {
		logger.Error("parse blueprint definition failed",
			zap.Uint64("blueprint id", blueprintId),
			zap.Error(err),
		)
		return nil, err
	}

	return &Scheduler{
		logger: logger.With(zap.String("scope", "scheduler"),
			zap.String("judgement id", judgement.Name),
		),
		mutex: &sync.Mutex{},
		C:     make(chan int, 1),
		Runtime: &Runtime{
			Problem:    problem,
			Submission: submission,
			Judgement:  judgement,
			graph:      graph,
			result:     make(map[int]*models.Slot),
		},
	}, nil
}

var ProviderSet = wire.NewSet(New)
