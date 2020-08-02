package services

import (
	"github.com/infinity-oj/server-v2/internal/app/problems/repositories"
	"github.com/infinity-oj/server-v2/internal/pkg/models"
	"go.uber.org/zap"
)

type ProblemsService interface {
	CreateProblem(name, title string) (p *models.Problem, err error)
	UpdateProblem(p *models.Problem, name string, title string) (*models.Problem, error)
	GetProblemById(id uint64) (p *models.Problem, err error)
	GetProblemByName(name string) (p *models.Problem, err error)
	GetProblems(page, pageSize int) (res []*models.Problem, err error)
}

type DefaultProblemService struct {
	logger     *zap.Logger
	Repository repositories.Repository
}

func (s DefaultProblemService) GetProblemById(id uint64) (p *models.Problem, err error) {
	p, err = s.Repository.GetProblemById(id)
	return
}

func (s DefaultProblemService) UpdateProblem(p *models.Problem, name string, title string) (*models.Problem, error) {
	p.Name = name
	p.Title = title
	if err := s.Repository.UpdateProblem(p); err != nil {
		s.logger.Error("update problem",
			zap.String("name", p.Name),
			zap.String("title", p.Title),
			zap.String("new name", name),
			zap.String("new title", title),
			zap.Error(err))
		return nil, err
	}
	return p, nil
}

func (s DefaultProblemService) GetProblems(page, pageSize int) (res []*models.Problem, err error) {
	offset := (page - 1) * pageSize
	res, err = s.Repository.GetProblems(offset, pageSize)
	if err != nil {
		s.logger.Error("get problems", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.Error(err))
		return nil, err
	}
	return
}

func (s DefaultProblemService) GetProblemByName(name string) (p *models.Problem, err error) {
	p, err = s.Repository.GetProblemByName(name)
	return
}

func (s DefaultProblemService) CreateProblem(name, title string) (p *models.Problem, err error) {
	if p, err = s.Repository.CreateProblem(name, title); err != nil {
		return p, err
	}
	return
}

func NewProblemService(logger *zap.Logger, Repository repositories.Repository) ProblemsService {
	return &DefaultProblemService{
		logger:     logger.With(zap.String("type", "DefaultProblemService")),
		Repository: Repository,
	}
}
