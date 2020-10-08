package services

import (
	"github.com/google/uuid"
	"github.com/infinity-oj/server-v2/internal/app/volumes/repositories"
	"github.com/infinity-oj/server-v2/pkg/models"
	"go.uber.org/zap"
)

type Service interface {
	CreateVolume() (*models.Volume, error)

	CreateFile(volumeName, filename string, file []byte) error
	CreateDirectory(volumeName, dirname string) error

	GetDirectory(volumeName, dirname string) (directories, files []string, err error)
	GetFile(volumeName, dirname, filename string) ([]byte, error)
}

type DefaultService struct {
	logger     *zap.Logger
	Repository repositories.Repository
}

func (d DefaultService) CreateVolume() (*models.Volume, error) {
	volume := uuid.New().String()
	err := d.Repository.CreateDirectory(volume, "/")
	if err != nil {
		return nil, err
	}
	return &models.Volume{Name: volume}, nil
}

func (d DefaultService) CreateFile(volumeName, filename string, file []byte) error {
	err := d.Repository.CreateFile(volumeName, filename, file)
	if err != nil {
		return err
	}
	return nil
}

func (d DefaultService) CreateDirectory(volumeName, dirname string) error {
	err := d.Repository.CreateDirectory(volumeName, dirname)
	if err != nil {
		return err
	}
	return nil
}

func (d DefaultService) GetDirectory(volumeName, dirname string) (directories, files []string, err error) {
	panic("implement me")
}

func (d DefaultService) GetFile(volumeName, dirname, filename string) ([]byte, error) {
	panic("implement me")
}

func NewVolumeService(logger *zap.Logger, Repository repositories.Repository) Service {
	return &DefaultService{
		logger:     logger.With(zap.String("type", "Account Repository")),
		Repository: Repository,
	}
}
