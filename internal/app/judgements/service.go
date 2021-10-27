package judgements

import (
	"errors"
	"net/http"

	"github.com/infinity-oj/server-v2/internal/app/blueprints"
	"github.com/infinity-oj/server-v2/pkg/models"
	"go.uber.org/zap"
)

type Service interface {
	GetJudgement(judgementId string) (*models.Judgement, error)
	GetJudgements(accountId uint64) ([]*models.Judgement, error)
	GetJudgementPrerequisites(blueprintId uint64) (string, error)
	CreateJudgement(accountId, blueprintId uint64, args map[string]interface{}) (int, *models.Judgement, error)
	UpdateJudgement(judgementId string, status models.JudgeStatus, score float64, msg string) (*models.Judgement, error)
}

type service struct {
	logger              *zap.Logger
	repository          Repository
	blueprintRepository blueprints.Repository

	dispatcher Dispatcher
}

func (s service) GetJudgementPrerequisites(blueprintId uint64) (string, error) {
	return "upload:*.cpp,*.c,*.py,*.zip", nil
}

func (s service) UpdateJudgement(judgementId string, status models.JudgeStatus, score float64, msg string) (*models.Judgement, error) {
	s.logger.Debug("update judgement",
		zap.String("judgement id", judgementId),
		zap.String("judge status", string(status)),
		zap.String("msg", msg),
		zap.Float64("score", score),
	)

	// get judgement with judgementId
	judgement, err := s.repository.GetJudgement(judgementId)
	if err != nil {
		return nil, err
	}

	judgement.Score = score
	judgement.Status = status
	judgement.Msg = msg

	err = s.repository.Update(judgement)

	return judgement, err
}

func (s service) CreateJudgement(accountId, blueprintId uint64, args map[string]interface{}) (int, *models.Judgement, error) {
	s.logger.Debug("create judgement",
		zap.Uint64("account id", accountId),
		zap.Uint64("blueprint id", blueprintId),
		zap.Any("args", args),
	)

	//judgements, err := d.repository.GetJudgementsByAccountId(accountId)
	//if err != nil {
	//	return http.StatusInternalServerError, nil, err
	//}
	//for _, judgement := range judgements {
	//	if judgement.Status == models.Accepted || judgement.Status == models.Pending {
	//		now := time.Now()
	//		judgeTime := judgement.CreatedAt
	//		dateEquals := func(a time.Time, b time.Time) bool {
	//			y1, m1, d1 := a.Date()
	//			y2, m2, d2 := b.Date()
	//			return y1 == y2 && m1 == m2 && d1 == d2
	//		}
	//		if dateEquals(judgeTime, now) {
	//			return http.StatusForbidden, nil, errors.New("previous judgement accepted today")
	//		}
	//	}
	//}

	// get blueprint
	blueprint, err := s.blueprintRepository.GetBlueprint(blueprintId)
	if err != nil {
		s.logger.Error("create judgement, get blueprint",
			zap.Uint64("blueprint id", blueprintId),
			zap.Error(err),
		)
		return http.StatusInternalServerError, nil, err
	}
	if blueprint == nil {
		return http.StatusInternalServerError, nil, errors.New("invalid request")
	}
	s.logger.Debug("create judgement",
		zap.String("blueprint definition", blueprint.Definition),
	)

	if args == nil {
		args = map[string]interface{}{}
	}

	// create judgement
	judgement, err := s.repository.Create(blueprintId, args)
	if err != nil {
		s.logger.Error("create judgement",
			zap.Uint64("blueprint id", blueprintId),
			zap.Error(err),
		)
		return http.StatusInternalServerError, nil, err
	}
	s.logger.Debug("create judgement successfully")

	GetDispatcher().PushJudgement(judgement)

	return http.StatusOK, judgement, err
}

func (s service) GetJudgement(judgementId string) (*models.Judgement, error) {
	judgement, err := s.repository.GetJudgement(judgementId)
	return judgement, err
}

func (s service) GetJudgements(accountId uint64) ([]*models.Judgement, error) {
	judgements, err := s.repository.GetJudgementsByAccountId(accountId)
	return judgements, err
}

func NewService(
	logger *zap.Logger,
	repository Repository,
	repository2 blueprints.Repository,
	dispatcher Dispatcher,
) Service {
	pendingJudgements, err := repository.GetPendingJudgements()
	if err != nil {
		panic(err)
	}

	for _, judgement := range pendingJudgements {
		logger.Debug("restore judgement", zap.String("judgement id", judgement.Name))
		dispatcher.PushJudgement(judgement)
	}

	srv := &service{
		logger:              logger.With(zap.String("type", "Judgement Service")),
		repository:          repository,
		blueprintRepository: repository2,
		dispatcher:          dispatcher,
	}

	return srv
}
