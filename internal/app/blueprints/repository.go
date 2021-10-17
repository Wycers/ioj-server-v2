package blueprints

import (
	"github.com/infinity-oj/server-v2/pkg/models"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type Repository interface {
	CreateBlueprint(definition string) (p *models.Blueprint, err error)
	GetBlueprint(id uint64) (*models.Blueprint, error)
	GetBlueprints() ([]*models.Blueprint, error)
}

type repository struct {
	logger *zap.Logger
	db     *gorm.DB
}

func (m repository) GetBlueprints() (blueprints []*models.Blueprint, err error) {
	if err = m.db.Model(&models.Blueprint{}).Find(&blueprints).Error; err != nil {
		return nil, err
	}
	return blueprints, nil
}

func (m repository) GetBlueprint(id uint64) (p *models.Blueprint, err error) {
	p = &models.Blueprint{}
	if err = m.db.First(p, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			m.logger.Error("Get blueprint failed", zap.Uint64("id", id), zap.Error(err))
		}
		return nil, err
	}
	return p, nil
}

func (m repository) CreateBlueprint(definition string) (blueprint *models.Blueprint, err error) {
	blueprint = &models.Blueprint{
		Name:       "",
		Title:      "",
		Definition: definition,
	}
	if err = m.db.Create(blueprint).Error; err != nil {
		m.logger.Error("create blueprint", zap.String("definition", definition))
		return nil, err
	}

	return blueprint, nil
}

func NewRepository(logger *zap.Logger, db *gorm.DB) Repository {
	return &repository{
		logger: logger.With(zap.String("type", "blueprint repository")),
		db:     db,
	}
}
