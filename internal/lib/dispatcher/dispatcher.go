package dispatcher

import (
	"sync"

	"github.com/infinity-oj/server-v2/internal/app/judgements"

	"github.com/google/wire"

	"github.com/spf13/cast"

	"github.com/infinity-oj/server-v2/internal/app/blueprints"
	"github.com/infinity-oj/server-v2/internal/app/problems"
	"github.com/infinity-oj/server-v2/internal/app/programs"
	"github.com/infinity-oj/server-v2/internal/app/submissions"

	"github.com/infinity-oj/server-v2/internal/lib/scheduler"

	"github.com/infinity-oj/server-v2/pkg/models"
	"go.uber.org/zap"
)

type dispatcher struct {
	c      chan *models.Judgement
	logger *zap.Logger
	br     blueprints.Repository
	pr     problems.Repository
	sr     submissions.Repository
	jr     judgements.Repository
	pgr    programs.Repository
}

func (d *dispatcher) PushJudgement(judgement *models.Judgement) {
	d.c <- judgement
}

func (d *dispatcher) execute(s *scheduler.Scheduler) {
	d.logger.Debug("execute runtime",
		zap.String("judgement id", s.Runtime.Judgement.Name),
	)
	judgement := s.Runtime.Judgement
	judgement.Status = models.Running
	if err := d.jr.Update(judgement); err != nil {
		d.logger.Error("update judgement", zap.Error(err))
	}
	go s.Execute()
	code := <-s.OnFinish()
	d.logger.Debug("finish runtime",
		zap.String("judgement id", s.Runtime.Judgement.Name),
		zap.Int("return code", code),
	)
	judgement = s.Runtime.Judgement
	switch judgement.Score {
	case -1:
		judgement.Status = models.Finished
		break
	case 0:
		judgement.Status = models.WrongAnswer
		break
	case 100:
		judgement.Status = models.Accepted
		break
	default:
		judgement.Status = models.PartiallyCorrect
	}
	if err := d.jr.Update(judgement); err != nil {
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

		if submissionId := cast.ToUint64(judgement.Args["submission"]); submissionId != 0 {
			// get submission
			submission, err := d.sr.GetSubmissionById(uint64(submissionId))
			if err != nil {
				d.logger.Error("create judgement",
					zap.Uint64("submission id", uint64(submissionId)),
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

		problemId := uint64(0)
		if instances.submission != nil {
			problemId = instances.submission.ProblemId
		} else {
			problemId = cast.ToUint64(judgement.Args["problem"])
		}
		if problemId != 0 {
			// get problem
			problem, err := d.pr.GetProblemById(problemId)
			if err != nil {
				panic(err)
			}
			if problem == nil {
				d.logger.Debug("create judgement instances, problem is nil", zap.Uint64("problem id", problemId))
			}
			instances.problem = problem
		}
		d.logger.Debug("get problem", zap.Any("problem", instances.problem))

		d.logger.Debug("create judgement instances", zap.Any("instances", instances))

		programs, err := d.pgr.GetPrograms()
		if err != nil {
			// TODO
		}
		s, err := scheduler.New(d.logger,
			instances.problem, instances.submission, instances.judgement,
			instances.blueprint, programs,
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

func GetDispatcher() judgements.Dispatcher {
	if instance == nil {
		panic("init failed")
	}
	return instance
}

func New(logger *zap.Logger, pr problems.Repository, sr submissions.Repository, jr judgements.Repository,
	br blueprints.Repository, pgr programs.Repository) judgements.Dispatcher {
	once.Do(func() {
		instance = &dispatcher{
			c:      make(chan *models.Judgement),
			logger: logger.With(zap.String("scope", "dispatcher")),
			br:     br,
			pr:     pr,
			sr:     sr,
			jr:     jr,
			pgr:    pgr,
		}

		go instance.run()
	})
	return instance
}

var ProviderSet = wire.NewSet(New)
