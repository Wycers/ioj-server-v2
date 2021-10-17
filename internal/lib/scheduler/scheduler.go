package scheduler

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/infinity-oj/server-v2/internal/app/processes"

	"github.com/google/uuid"
	"github.com/google/wire"

	"github.com/infinity-oj/server-v2/internal/lib/nodeEngine"

	"go.uber.org/zap"

	"github.com/infinity-oj/server-v2/pkg/models"
)

type Runtime struct {
	Blueprint  *models.Blueprint
	Problem    *models.Problem
	Submission *models.Submission
	Judgement  *models.Judgement

	graph  *nodeEngine.Graph
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
	trigger := make(chan int, 100)
	pendingCnt := 0

	trigger <- 0
	s.logger.Debug("scheduler execute blueprint")
	for _ = range trigger {
		s.logger.Debug("scheduler", zap.Int("pending count", pendingCnt))

		ids := s.Runtime.graph.Run()
		for _, block := range ids {
			var inputs models.Slots
			for _, linkId := range block.Inputs {
				if data, ok := s.Runtime.result[linkId]; ok {
					inputs = append(inputs, data)
				} else {
					s.logger.Error("wrong process definition")
					code = -1
					return
				}
			}

			newProcess := &models.Process{
				Model: models.Model{
					ID:        0,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					DeletedAt: nil,
				},
				Type:        block.Type,
				ProcessId:   uuid.New().String(),
				JudgementId: s.Runtime.Judgement.JudgementId,
				Properties:  block.Properties,
				Inputs:      inputs,
				Outputs:     models.Slots{},
			}

			go func(block *nodeEngine.Block) {
				pendingCnt++
				select {
				case outputs := <-processes.GetManager().Push(block.Id, newProcess):
					s.logger.Debug("process finished")

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

					trigger <- pendingCnt
				case <-time.After(time.Second * 5):
					s.logger.Debug("process timeout after 5s")
				}
				pendingCnt--
			}(block)
		}

		if pendingCnt == 0 {
			break
		}
	}
}

func (s *Scheduler) OnFinish() <-chan int {
	return s.C
}

func New(logger *zap.Logger,
	problem *models.Problem, submission *models.Submission, judgement *models.Judgement, blueprint *models.Blueprint,
) (*Scheduler, error) {
	blueprintId := blueprint.ID

	definition := blueprint.Definition
	if submission != nil {
		definition = strings.ReplaceAll(definition, "<userVolume>", submission.UserVolume)
		definition = strings.ReplaceAll(definition, "<publicVolume>", problem.PublicVolume)
		definition = strings.ReplaceAll(definition, "<privateVolume>", problem.PrivateVolume)
	}
	graph, err := nodeEngine.NewGraphByDefinition(definition)
	if err != nil {
		logger.Error("parse blueprint definition failed",
			zap.Uint64("blueprint id", blueprintId),
			zap.Error(err),
		)
		return nil, err
	}

	return &Scheduler{
		logger: logger.With(zap.String("scope", "scheduler"),
			zap.String("judgement id", judgement.JudgementId),
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
