package services

import (
	"github.com/infinity-oj/server-v2/internal/app/processes/repositories"
	"github.com/infinity-oj/server-v2/pkg/models"
	"go.uber.org/zap"
)

type ProcessesService interface {
	CreateProcess(definition string) (p *models.Process, err error)
	GetProcess(id uint64) (p *models.Process, err error)
}

type DefaultProblemService struct {
	logger     *zap.Logger
	Repository repositories.Repository
}

func (s DefaultProblemService) CreateProcess(definition string) (p *models.Process, err error) {
	s.logger.Debug("create process",
		zap.String("definition", definition),
	)
	if p, err = s.Repository.CreateProcess(definition); err != nil {
		return p, err
	}
	return
}
func (s DefaultProblemService) GetProcess(id uint64) (p *models.Process, err error) {
	s.logger.Debug("get process",
		zap.Uint64("id", id),
	)
	p, err = s.Repository.GetProcess(id)
	return
}

func NewProcessService(logger *zap.Logger, Repository repositories.Repository) ProcessesService {
	return &DefaultProblemService{
		logger:     logger.With(zap.String("type", "DefaultProblemService")),
		Repository: Repository,
	}
}
