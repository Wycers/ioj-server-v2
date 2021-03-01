package processes

import (
	"github.com/infinity-oj/server-v2/pkg/models"
	"go.uber.org/zap"
)

type Service interface {
	CreateProcess(definition string) (p *models.Process, err error)
	GetProcess(id uint64) (p *models.Process, err error)
}

type service struct {
	logger     *zap.Logger
	Repository Repository
}

func (s service) CreateProcess(definition string) (p *models.Process, err error) {
	s.logger.Debug("create process",
		zap.String("definition", definition),
	)
	if p, err = s.Repository.CreateProcess(definition); err != nil {
		return p, err
	}
	return
}
func (s service) GetProcess(id uint64) (p *models.Process, err error) {
	s.logger.Debug("get process",
		zap.Uint64("id", id),
	)
	p, err = s.Repository.GetProcess(id)
	return
}

func NewService(logger *zap.Logger, Repository Repository) Service {
	return &service{
		logger:     logger.With(zap.String("type", "service")),
		Repository: Repository,
	}
}
