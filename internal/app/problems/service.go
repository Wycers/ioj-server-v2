package problems

import (
	"github.com/infinity-oj/server-v2/pkg/models"
	"go.uber.org/zap"
)

type Service interface {
	CreateProblem(name, title string) (p *models.Problem, err error)
	UpdateProblem(p *models.Problem, name, title, publicVolume, privateVolume string) (*models.Problem, error)
	GetProblemById(id uint64) (p *models.Problem, err error)
	GetProblemByName(name string) (p *models.Problem, err error)
	GetProblems(page, pageSize int) (res []*models.Problem, err error)
}

type service struct {
	logger     *zap.Logger
	Repository Repository
}

func (s service) GetProblemById(id uint64) (p *models.Problem, err error) {
	p, err = s.Repository.GetProblemById(id)
	return
}

func (s service) UpdateProblem(p *models.Problem, name, title, publicVolume, privateVolume string) (*models.Problem, error) {
	p.Name = name
	p.Title = title
	p.PublicVolume = publicVolume
	p.PrivateVolume = privateVolume
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

func (s service) GetProblems(page, pageSize int) (res []*models.Problem, err error) {
	offset := (page - 1) * pageSize
	res, err = s.Repository.GetProblems(offset, pageSize)
	if err != nil {
		s.logger.Error("get problems", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.Error(err))
		return nil, err
	}
	return
}

func (s service) GetProblemByName(name string) (p *models.Problem, err error) {
	p, err = s.Repository.GetProblemByName(name)
	return
}

func (s service) CreateProblem(name, title string) (p *models.Problem, err error) {
	if p, err = s.Repository.CreateProblem(name, title); err != nil {
		return p, err
	}
	return
}

func NewService(logger *zap.Logger, Repository Repository) Service {
	return &service{
		logger:     logger.With(zap.String("type", "ProblemService")),
		Repository: Repository,
	}
}
