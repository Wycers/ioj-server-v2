package scheduler

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/infinity-oj/server-v2/internal/lib/manager"

	"github.com/infinity-oj/server-v2/internal/lib/engine/scene"

	"github.com/google/uuid"
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

	wg := new(sync.WaitGroup)
	trigger := make(chan int32, 100)
	var n int32 = 0
	trigger <- atomic.AddInt32(&n, 1)
	for _ = range trigger {
		ids := s.Runtime.graph.Run()
		for _, block := range ids {
			var inputs models.Slots
			for _, linkId := range block.Inputs {
				if data, ok := s.Runtime.result[linkId]; ok {
					inputs = append(inputs, data)
				} else {
					s.logger.Error("wrong process definition", zap.Int("link id", linkId))
					code = -1
					return
				}
			}

			atomic.AddInt32(&n, 1)
			wg.Add(1)
			go func(block *engine.Block) {
				newProcess := &models.Process{
					Model: models.Model{
						ID:        0,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						DeletedAt: nil,
					},
					Type:        block.Type,
					ProcessId:   uuid.New().String(),
					JudgementId: s.Runtime.Judgement.Name,
					Properties:  block.Properties,
					Inputs:      inputs,
					Outputs:     models.Slots{},
				}

				s.logger.Debug("process started", zap.String("process id", newProcess.ProcessId))

				select {
				case outputs := <-manager.GetManager().Push(block.Id, newProcess):
					s.logger.Debug("process finished normally", zap.String("process id", newProcess.ProcessId))

					blockId := block.Id
					if len(block.Output) != len(*outputs) {
						s.logger.Error(fmt.Sprintf("output slots mismatch, block %d expects %d but %d",
							block.Id,
							len(block.Output),
							len(*outputs),
						))
						return
					}

					for index, output := range *outputs {
						links := s.Runtime.graph.FindLinkBySourcePort(blockId, index)
						for _, link := range links {
							s.Runtime.result[link.Id] = output
						}
					}
					block.Done()
					trigger <- atomic.AddInt32(&n, 1)
				case <-time.After(time.Second * 5):
					s.logger.Debug("process timeout after 5s", zap.String("process id", newProcess.ProcessId))
				}

				s.logger.Debug("process ended", zap.String("process id", newProcess.ProcessId))
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
		definition = strings.ReplaceAll(definition, "<userVolume>", submission.UserVolume)
		definition = strings.ReplaceAll(definition, "${account_id}", fmt.Sprintf("%d", submission.SubmitterId))

	}
	if problem != nil {
		definition = strings.ReplaceAll(definition, "<publicVolume>", problem.PublicVolume)
		definition = strings.ReplaceAll(definition, "<privateVolume>", problem.PrivateVolume)
		definition = strings.ReplaceAll(definition, "${problem_id}", problem.Name)
	}
	fmt.Println(definition)
	//graph, err := engine.NewGraphByDefinition(definition)
	var bs []*scene.BlockDefinition
	for _, p := range programs {
		bs = append(bs, scene.NewBlockDefinition(p.Definition))
	}
	s := scene.NewScene(definition)
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
