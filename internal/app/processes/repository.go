package processes

import (
	"github.com/infinity-oj/server-v2/pkg/models"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type Repository interface {
	CreateProcess(definition string) (p *models.Process, err error)
	GetProcess(id uint64) (*models.Process, error)
}

type repository struct {
	logger *zap.Logger
	db     *gorm.DB
}

func (m repository) GetProcess(id uint64) (p *models.Process, err error) {
	p = &models.Process{}
	if err = m.db.First(p, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			m.logger.Error("Get process failed", zap.Uint64("id", id), zap.Error(err))
		}
		return nil, err
	}
	return p, nil
}

func (m repository) CreateProcess(definition string) (process *models.Process, err error) {
	process = &models.Process{
		//FileIoInputName:  "",
		//FileIoOutputName: "",
		Definition: definition,
	}
	if err = m.db.Create(process).Error; err != nil {
		m.logger.Error("create process", zap.String("definition", definition))
		return nil, err
	}

	return process, nil
}

func NewRepository(logger *zap.Logger, db *gorm.DB) Repository {
	return &repository{
		logger: logger.With(zap.String("type", "Repository")),
		db:     db,
	}
}
