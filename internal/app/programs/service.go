package programs

import (
	"github.com/infinity-oj/server-v2/pkg/models"
	"go.uber.org/zap"
)

type Service interface {
	CreateProgram(definition string) (p *models.Program, err error)
	GetProgram(id uint64) (p *models.Program, err error)
	GetPrograms() (p []*models.Program, err error)
}

type service struct {
	logger     *zap.Logger
	Repository Repository
}

func (s service) GetPrograms() (p []*models.Program, err error) {
	s.logger.Debug("get programs")
	p, err = s.Repository.GetPrograms()
	return
}

func (s service) CreateProgram(definition string) (p *models.Program, err error) {
	s.logger.Debug("create program",
		zap.String("definition", definition),
	)
	if p, err = s.Repository.CreateProgram(definition); err != nil {
		return p, err
	}
	return
}
func (s service) GetProgram(id uint64) (p *models.Program, err error) {
	s.logger.Debug("get program",
		zap.Uint64("id", id),
	)
	p, err = s.Repository.GetProgram(id)
	return
}

func NewService(logger *zap.Logger, Repository Repository) Service {
	return &service{
		logger:     logger.With(zap.String("type", "program service")),
		Repository: Repository,
	}
}
