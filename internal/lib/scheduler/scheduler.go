package scheduler

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
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

type JudgeResult struct {
	Code    int
	Score   float64
	Message string
}

type Scheduler struct {
	logger *zap.Logger
	mutex  *sync.Mutex

	Runtime *Runtime

	C chan JudgeResult
}

func (s *Scheduler) Execute() {
	result := JudgeResult{
		Code:    0,
		Score:   0,
		Message: "",
	}
	defer func() {
		s.C <- result
	}()
	s.logger.Debug("scheduler: execution started")
	trigger := make(chan int, 100)
	pendingCnt := 0
	wg := new(sync.WaitGroup)

	trigger <- 0
	for _ = range trigger {
		s.logger.Debug("scheduler", zap.Int("pending count", pendingCnt))

		ids := s.Runtime.graph.Run()
		for _, block := range ids {
			var inputs models.Slots
			for _, linkId := range block.Inputs {
				if data, ok := s.Runtime.result[linkId]; ok {
					inputs = append(inputs, data)
				} else {
					s.logger.Error("wrong process definition", zap.Int("link id", linkId))
					result.Score = -1
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
			if newProcess.Type == "result" {
				v := newProcess.Inputs[0].Value
				score := func() float64 {
					switch i := v.(type) {
					case float64:
						return i
					case float32:
						return float64(i)
					case int64:
						return float64(i)
					case int32:
						return float64(i)
					case string:
						if s, err := strconv.ParseFloat(i, 64); err == nil {
							return s
						}
					}
					return math.NaN()
				}()
				result.Score = score
			}

			pendingCnt++
			wg.Add(1)
			go func(block *engine.Block) {
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
						fmt.Println(output)
						links := s.Runtime.graph.FindLinkBySourcePort(blockId, index)
						for _, link := range links {
							s.Runtime.result[link.Id] = output
						}
					}
					block.Done()
					trigger <- pendingCnt
				case <-time.After(time.Second * 5):
					s.logger.Debug("process timeout after 5s", zap.String("process id", newProcess.ProcessId))
					if pendingCnt == 1 {
						s.logger.Debug("pending count is 0, closing")
						close(trigger)
					}
				}
				s.logger.Debug("process ended", zap.String("process id", newProcess.ProcessId))
				pendingCnt--
				wg.Done()
			}(block)
		}
	}
	s.logger.Debug("scheduler: execution ended")
	wg.Wait()
}

func (s *Scheduler) OnFinish() <-chan JudgeResult {
	return s.C
}

func New(logger *zap.Logger,
	problem *models.Problem, submission *models.Submission, judgement *models.Judgement,
	blueprint *models.Blueprint, programs []*models.Program,
) (*Scheduler, error) {
	blueprintId := blueprint.ID

	definition := blueprint.Definition
	if submission != nil {
		definition = strings.ReplaceAll(definition, "<userVolume>", submission.UserVolume)
		definition = strings.ReplaceAll(definition, "<publicVolume>", problem.PublicVolume)
		definition = strings.ReplaceAll(definition, "<privateVolume>", problem.PrivateVolume)
	}
	//graph, err := engine.NewGraphByDefinition(definition)
	var bs []*scene.BlockDefinition
	for _, p := range programs {
		bs = append(bs, scene.NewBlockDefinition(p.Definition))
	}
	s := scene.NewScene(blueprint.Definition)
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
			zap.String("judgement id", judgement.JudgementId),
		),
		mutex: &sync.Mutex{},
		C:     make(chan JudgeResult, 1),
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
