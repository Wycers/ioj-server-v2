package blueprints

import (
	"github.com/infinity-oj/server-v2/pkg/models"
	"go.uber.org/zap"
)

type Service interface {
	CreateBlueprint(definition string) (p *models.Blueprint, err error)
	GetBlueprint(id uint64) (p *models.Blueprint, err error)
	GetBlueprints() (p []*models.Blueprint, err error)
}

type service struct {
	logger     *zap.Logger
	Repository Repository
}

func (s service) GetBlueprints() (p []*models.Blueprint, err error) {
	s.logger.Debug("get blueprints")
	p, err = s.Repository.GetBlueprints()
	return
}

func (s service) CreateBlueprint(definition string) (p *models.Blueprint, err error) {
	s.logger.Debug("create blueprint",
		zap.String("definition", definition),
	)
	if p, err = s.Repository.CreateBlueprint(definition); err != nil {
		return p, err
	}
	return
}
func (s service) GetBlueprint(id uint64) (p *models.Blueprint, err error) {
	s.logger.Debug("get blueprint",
		zap.Uint64("id", id),
	)
	p, err = s.Repository.GetBlueprint(id)
	return
}

func NewService(logger *zap.Logger, Repository Repository) Service {
	return &service{
		logger:     logger.With(zap.String("type", "service")),
		Repository: Repository,
	}
}
