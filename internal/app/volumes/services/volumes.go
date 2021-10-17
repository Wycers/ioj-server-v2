package services

import (
	"fmt"
	"os"
	"time"

	"github.com/infinity-oj/server-v2/internal/pkg/crypto"

	"github.com/google/uuid"
	"github.com/infinity-oj/server-v2/internal/app/volumes/repositories"
	"github.com/infinity-oj/server-v2/internal/app/volumes/storages"
	"github.com/infinity-oj/server-v2/pkg/models"
	"go.uber.org/zap"
)

type Service interface {
	CreateVolume(createdBy uint64) (*models.Volume, error)

	CreateFile(baseVolumeName, filename string, file []byte) (*models.Volume, error)
	CreateDirectory(volumeName, dirname string) (*models.Volume, error)

	GetVolume(volumeName string) (*models.Volume, error)

	DownloadDirectory(volumeName, dirname string) (file *os.File, err error)
	GetDirectory(volumeName, dirname string) (directories, files []string, err error)
	GetFile(volumeName, dirname, filename string) ([]byte, error)

	GetVolumeFileRecords(volume *models.Volume, filePath string) (*models.FileRecords, error)
}

type DefaultService struct {
	logger     *zap.Logger
	Repository repositories.Repository
	Storage    storages.Storage
}

func (d DefaultService) GetVolume(volumeName string) (*models.Volume, error) {
	volume, err := d.Repository.GetVolume(volumeName)
	if err != nil {
		return nil, err
	}
	volumes, err := d.GetVolumeChain(volume)
	if err != nil {
		return nil, err
	}
	fileRecords := models.FileRecords{}

	for i, _ := range volumes {
		cur := volumes[i]
		for _, curRecord := range cur.FileRecords {
			fmt.Println(curRecord)

			tag := -1
			for i, _ := range fileRecords {
				if fileRecords[i].FilePath == curRecord.FilePath {
					tag = i
					break
				}
			}
			if curRecord.Opt == "Add" {
				if tag == -1 {
					fileRecords = append(fileRecords, curRecord)
				} else {
					fileRecords[tag] = curRecord
				}
			} else {
				if tag != -1 {
					fileRecords = append(fileRecords[:tag], fileRecords[tag+1:]...)
				}
			}
		}
		fmt.Println()
	}
	volume.FileRecords = fileRecords

	return volume, nil
}

func (d DefaultService) GetVolumeChain(volume *models.Volume) ([]*models.Volume, error) {
	if volume.Base == 0 {
		// have no previous volume
		return []*models.Volume{volume}, nil
	}
	baseVolume, err := d.Repository.GetVolumeByID(volume.Base)
	if err != nil {
		return nil, err
	}
	previous, err := d.GetVolumeChain(baseVolume)
	if err != nil {
		return nil, err
	}
	res := append(previous, volume)
	return res, nil
}

func (d DefaultService) GetVolumeFileRecords(volume *models.Volume, filePath string) (*models.FileRecords, error) {
	//if volume.Base == 0 {
	//	// have no previous volume
	//	return volume.FileRecords, nil
	//}
	//if filePath != "" {
	//	for _, record := range *volume.FileRecords {
	//		if record.FilePath == filePath {
	//			return *
	//		}
	//	}
	//}
	//baseVolume, err := d.repository.GetVolumeByID(volume.Base)
	//if err != nil {
	//	return nil, err
	//}
	//_, err = d.GetVolumeFileRecords(baseVolume, filePath)
	//if err != nil {
	//	return nil, err
	//}
	return nil, nil
}

func (d DefaultService) CreateVolume(accountID uint64) (*models.Volume, error) {
	volumeName := uuid.New().String()
	volume, err := d.Repository.CreateVolume(nil, accountID, volumeName)
	if err != nil {
		d.logger.Error("create volume", zap.Error(err))
		return nil, err
	}

	if err := d.Storage.CreateDirectory(volumeName, "/"); err != nil {
		d.logger.Error("create directory", zap.Error(err))
		return nil, err
	}
	return volume, nil
}

func (d DefaultService) CreateFile(baseVolumeName, filename string, file []byte) (*models.Volume, error) {
	baseVolume, err := d.Repository.GetVolume(baseVolumeName)
	if err != nil {
		return baseVolume, err
	}
	volume, err := d.CreateVolume(1)
	if err != nil {
		return nil, err
	}
	volume.Base = baseVolume.ID
	volumePath := time.Now().Format("20060102150405") + crypto.Sha256Bytes(file)
	volume.FileRecords = models.FileRecords{
		&models.FileRecord{
			Opt:        "Add",
			FileType:   "f",
			FilePath:   filename,
			VolumeName: volume.Name,
			VolumePath: volumePath,
		},
	}
	if volume, err = d.Repository.UpdateVolume(volume); err != nil {
		return nil, err
	}
	if err := d.Storage.CreateFile(volume.Name, volumePath, file); err != nil {
		return nil, err
	}
	return volume, nil
}

func (d DefaultService) CreateDirectory(baseVolumeName, dirname string) (*models.Volume, error) {
	baseVolume, err := d.Repository.GetVolume(baseVolumeName)
	if err != nil {
		return baseVolume, err
	}
	volume, err := d.CreateVolume(1)
	if err != nil {
		return nil, err
	}
	volume.Base = baseVolume.ID
	if err != nil {
		return nil, err
	}
	volume.FileRecords = models.FileRecords{
		&models.FileRecord{
			Opt:      "Add",
			FilePath: dirname,
			FileType: "d",
		},
	}
	if volume, err = d.Repository.UpdateVolume(volume); err != nil {
		return nil, err
	}
	return volume, nil
}

func (d DefaultService) GetDirectory(volumeName, dirname string) (directories, files []string, err error) {
	panic("implement me")
}

func (d DefaultService) GetFile(volumeName, dirname, filename string) ([]byte, error) {
	panic("implement me")
}

func (d DefaultService) DownloadDirectory(volumeName, dirname string) (file *os.File, err error) {
	volume, err := d.GetVolume(volumeName)
	if err != nil {
		return nil, err
	}
	return d.Storage.ArchiveVolume(volume)
}

func NewVolumeService(logger *zap.Logger, Storage storages.Storage, Repository repositories.Repository) Service {
	return &DefaultService{
		logger:     logger.With(zap.String("type", "Account Storage")),
		Storage:    Storage,
		Repository: Repository,
	}
}
