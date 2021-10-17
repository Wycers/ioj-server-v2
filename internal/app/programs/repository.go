package programs

import (
	"github.com/infinity-oj/server-v2/pkg/models"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type Repository interface {
	CreateProgram(definition string) (p *models.Program, err error)
	GetProgram(id uint64) (*models.Program, error)
	GetPrograms() ([]*models.Program, error)
}

type repository struct {
	logger *zap.Logger
	db     *gorm.DB
}

func (m repository) GetPrograms() (programs []*models.Program, err error) {
	if err = m.db.Model(&models.Program{}).Find(&programs).Error; err != nil {
		return nil, err
	}
	return programs, nil
}

func (m repository) GetProgram(id uint64) (p *models.Program, err error) {
	p = &models.Program{}
	if err = m.db.First(p, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			m.logger.Error("Get program failed", zap.Uint64("id", id), zap.Error(err))
		}
		return nil, err
	}
	return p, nil
}

func (m repository) CreateProgram(definition string) (program *models.Program, err error) {
	program = &models.Program{
		//FileIoInputName:  "",
		//FileIoOutputName: "",
		Definition: definition,
	}
	if err = m.db.Create(program).Error; err != nil {
		m.logger.Error("create program", zap.String("definition", definition))
		return nil, err
	}

	return program, nil
}

func NewRepository(logger *zap.Logger, db *gorm.DB) Repository {
	return &repository{
		logger: logger.With(zap.String("type", "Program repository")),
		db:     db,
	}
}
