package problems

import (
	"github.com/infinity-oj/server-v2/pkg/models"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Repository interface {
	CreateProblem(name, title string) (p *models.Problem, err error)
	UpdateProblem(p *models.Problem) error
	CreatePage(problemId uint64, locale, title, description string) (p *models.Page, err error)

	GetProblemById(id uint64) (*models.Problem, error)
	GetProblemByName(name string) (p *models.Problem, err error)
	GetProblems(offset, limit int) (p []*models.Problem, err error)

	CountProblems() (count int64)
}

type DefaultRepository struct {
	logger *zap.Logger
	db     *gorm.DB
}

func (m DefaultRepository) GetProblemById(id uint64) (p *models.Problem, err error) {
	p = &models.Problem{}
	if err = m.db.First(p, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			m.logger.Error("Query problem failed", zap.Uint64("id", id), zap.Error(err))
		}
		return nil, err
	}
	return p, nil
}

func (m DefaultRepository) UpdateProblem(p *models.Problem) error {
	return m.db.Save(p).Error
}

func (m DefaultRepository) CountProblems() (count int64) {
	m.db.Table("problems").Count(&count)
	return
}

func (m DefaultRepository) GetProblems(offset, limit int) (p []*models.Problem, err error) {
	if err = m.db.Table("problems").Limit(limit).Offset(offset).Find(&p).Error; err != nil {
		return nil, err
	}
	return
}

func (m DefaultRepository) GetProblemByName(name string) (p *models.Problem, err error) {
	p = &models.Problem{}
	if err = m.db.Where(&models.Problem{Name: name}).First(p).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			m.logger.Error("Query problem failed", zap.String("name", name), zap.Error(err))
		}
		return nil, err
	}
	return
}

func (m DefaultRepository) CreateProblem(name, title string) (problem *models.Problem, err error) {
	problem = &models.Problem{
		Name:  name,
		Title: title,
	}
	if err = m.db.Create(problem).Error; err != nil {
		m.logger.Error("create problem", zap.String("name", name))
		return nil, err
	}

	return problem, nil
}

func (m DefaultRepository) CreatePage(problemId uint64, locale, title, description string) (page *models.Page, err error) {
	m.logger.Debug("create page",
		zap.Uint64("problemId", problemId),
		zap.String("locale", locale),
	)
	page = &models.Page{
		ProblemId:   problemId,
		Locale:      locale,
		Title:       title,
		Description: description,
	}
	if err = m.db.Create(page).Error; err != nil {
		m.logger.Error("create page",
			zap.Uint64("problemId", problemId),
			zap.String("locale", locale),
			zap.Error(err),
		)
		return nil, errors.Wrapf(err, "create problem page with title: %s", title)
	}
	return page, nil
}

func NewRepository(logger *zap.Logger, db *gorm.DB) Repository {
	return &DefaultRepository{
		logger: logger.With(zap.String("type", "Repository")),
		db:     db,
	}
}
