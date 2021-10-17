package judgements

import (
	"sync"

	"github.com/infinity-oj/server-v2/internal/lib/scheduler"

	"go.uber.org/zap"

	"github.com/infinity-oj/server-v2/internal/app/blueprints"

	"github.com/infinity-oj/server-v2/internal/app/problems"
	"github.com/infinity-oj/server-v2/internal/app/submissions"

	"github.com/infinity-oj/server-v2/pkg/models"
)

type Dispatcher interface {
	PushJudgement(judgement *models.Judgement)
}

type dispatcher struct {
	c      chan *models.Judgement
	logger *zap.Logger
	br     blueprints.Repository
	pr     problems.Repository
	sr     submissions.Repository
	jr     Repository
}

func (d *dispatcher) PushJudgement(judgement *models.Judgement) {
	d.c <- judgement
}

func (d *dispatcher) execute(s *scheduler.Scheduler) {
	d.logger.Debug("execute runtime",
		zap.String("judgement id", s.Runtime.Judgement.JudgementId),
	)
	go s.Execute()
	code := <-s.OnFinish()
	d.logger.Debug("finish runtime",
		zap.String("judgement id", s.Runtime.Judgement.JudgementId),
		zap.Int("return code", code),
	)
	judgement := s.Runtime.Judgement
	judgement.Status = models.Accepted
	err := d.jr.Update(judgement)
	if err != nil {
		d.logger.Error("update judgement", zap.Error(err))
	}
}

func (d *dispatcher) run() {
	for judgement := range d.c {
		instances := &(struct {
			blueprint  *models.Blueprint
			judgement  *models.Judgement
			problem    *models.Problem
			submission *models.Submission
		}{
			judgement: judgement,
		})
		d.logger.Debug("get judgement", zap.Any("judgement", judgement))

		// get blueprint
		blueprint, err := d.br.GetBlueprint(judgement.BlueprintId)
		if err != nil {
			panic(err)
		}
		if blueprint == nil {
			continue
		}
		instances.blueprint = blueprint
		d.logger.Debug("get blueprint", zap.Any("blueprint", blueprint))

		submissionId, ok := judgement.Args["submission"].(uint64)
		if ok {
			// get submission
			submission, err := d.sr.GetSubmissionById(submissionId)
			if err != nil {
				d.logger.Error("create judgement",
					zap.Uint64("submission id", submissionId),
					zap.Error(err),
				)
				continue
			}
			if submission == nil {
				d.logger.Debug("create judgement",
					zap.String("submission user space", submission.UserVolume),
				)
			}

			instances.submission = submission
		}
		d.logger.Debug("get submission", zap.Any("submission", instances.submission))

		problemId, ok := judgement.Args["problem"].(float64)
		if ok {
			// get problem
			problem, err := d.pr.GetProblemById(uint64(problemId))
			if err != nil {
				panic(err)
			}
			if problem == nil {
				d.logger.Debug("create judgement instances, problem is nil")
			}
			instances.problem = problem
		}
		d.logger.Debug("get problem", zap.Any("problem", instances.problem))

		d.logger.Debug("create judgement instances", zap.Any("instances", instances))

		s, err := scheduler.New(d.logger,
			instances.problem, instances.submission, instances.judgement, instances.blueprint,
		)
		if err != nil {
			d.logger.Error("create scheduler error", zap.Error(err))
			continue
		}

		go d.execute(s)
	}
}

var instance *dispatcher
var once sync.Once

func GetDispatcher() Dispatcher {
	if instance == nil {
		panic("init failed")
	}
	return instance
}

func InitDispatcher(logger *zap.Logger, br blueprints.Repository, pr problems.Repository, sr submissions.Repository, jr Repository) Dispatcher {
	once.Do(func() {
		instance = &dispatcher{
			c:      make(chan *models.Judgement),
			logger: logger.With(zap.String("scope", "dispatcher")),
			br:     br,
			pr:     pr,
			sr:     sr,
			jr:     jr,
		}

		go instance.run()
	})
	return instance
}
