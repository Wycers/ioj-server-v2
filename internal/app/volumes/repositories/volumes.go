package repositories

import (
	"github.com/infinity-oj/server-v2/pkg/models"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type Repository interface {
	CreateVolume(baseVolume *models.Volume, accountId uint64, volumeName string) (*models.Volume, error)
	UpdateVolume(volume *models.Volume) (*models.Volume, error)
	GetVolume(volumeName string) (*models.Volume, error)
	GetVolumeByID(volumeID uint64) (*models.Volume, error)
}

type repository struct {
	logger *zap.Logger
	db     *gorm.DB
}

func (r repository) UpdateVolume(volume *models.Volume) (*models.Volume, error) {
	if err := r.db.Save(&volume).Error; err != nil {
		return nil, err
	}
	return volume, nil
}

func (r repository) GetVolumeByID(volumeID uint64) (*models.Volume, error) {
	volume := &models.Volume{}
	if err := r.db.Where("id = ?", volumeID).Limit(1).Find(volume).Error; err != nil {
		return nil, err
	}
	return volume, nil
}

func (r repository) GetVolume(volumeName string) (*models.Volume, error) {
	volume := &models.Volume{}
	if err := r.db.Where("name = ?", volumeName).Limit(1).Find(volume).Error; err != nil {
		return nil, err
	}
	return volume, nil
}

func (r repository) CreateVolume(baseVolume *models.Volume, accountId uint64, volumeName string) (*models.Volume, error) {
	var baseVolumeID uint64 = 0
	if baseVolume != nil {
		baseVolumeID = baseVolume.ID
	}
	volume := &models.Volume{
		Base:        baseVolumeID,
		CreatedBy:   accountId,
		Name:        volumeName,
		FileRecords: models.FileRecords{},
	}

	if err := r.db.Create(volume).Error; err != nil {
		r.logger.Error("create volume", zap.String("name", volumeName), zap.Error(err))
		return nil, err
	}

	return volume, nil
}

func NewRepository(logger *zap.Logger, db *gorm.DB) Repository {
	return &repository{
		logger: logger.With(zap.String("type", "repository")),
		db:     db,
	}
}
