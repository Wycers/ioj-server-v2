package services

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"

	"github.com/infinity-oj/server-v2/internal/pkg/crypto"

	"github.com/google/uuid"
	"github.com/infinity-oj/server-v2/internal/app/volumes/repositories"
	"github.com/infinity-oj/server-v2/internal/app/volumes/storages"
	"github.com/infinity-oj/server-v2/pkg/models"
	"go.uber.org/zap"
)

type Service interface {
	CreateVolume(createdBy uint64) (*models.Volume, error)

	CreateDirectory(baseVolumeName, dirname string) (*models.Volume, error)
	CreateFile(baseVolumeName, dirname, filename string, file []byte) (*models.Volume, error)
	RemoveFile(baseVolumeName, dirname, filename string) (*models.Volume, error)
	CopyFile(ov, od, of, nv, nd, nf string) (*models.Volume, error)

	GetVolume(volumeName string) (*models.Volume, error)
	GetDirectory(volumeName, dirname string) (file *os.File, err error)
	GetFile(volumeName, filename string) (*os.File, error)
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

	mp := make(map[string]*models.FileRecord)

	for i, _ := range volumes {
		cur := volumes[i]
		for _, currentRecord := range cur.FileRecords {
			fmt.Println(currentRecord)

			filepath := currentRecord.FilePath
			if _, ok := mp[filepath]; ok {
				// has previous record
				if currentRecord.Opt == "add" {
					// add a file with same file path, cover previous.
					mp[filepath] = currentRecord
				}
				if currentRecord.Opt == "del" {
					// del a file with same file path, remove previous.
					delete(mp, filepath)
				}
			} else {
				// no previous record
				if currentRecord.Opt == "add" {
					// add
					mp[filepath] = currentRecord
				}
			}
		}
	}

	fileRecords := models.FileRecords{}
	for _, v := range mp {
		fileRecords = append(fileRecords, v)
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

func (d DefaultService) CopyFile(ov, od, of, nv, nd, nf string) (*models.Volume, error) {
	oldVolume, err := d.GetVolume(ov)
	if err != nil {
		return nil, err
	}
	oldFilePath := filepath.Join(od, of)
	for _, fileRecord := range oldVolume.FileRecords {
		if fileRecord.FilePath == oldFilePath {
			newVolume, err := d.Repository.GetVolume(nv)
			if err != nil {
				return nil, err
			}
			newVolume.FileRecords = append(newVolume.FileRecords,
				&models.FileRecord{
					Opt:        "add",
					FileType:   "f",
					FilePath:   filepath.Join("/", nd, nf),
					VolumeName: fileRecord.VolumeName,
					VolumePath: fileRecord.VolumePath,
				})
			if newVolume, err = d.Repository.UpdateVolume(newVolume); err != nil {
				return nil, err
			}
			return newVolume, err
		}
	}
	return nil, errors.New("failed")
}

func (d DefaultService) CreateFile(baseVolumeName, dirname, filename string, file []byte) (*models.Volume, error) {
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
			Opt:        "add",
			FileType:   "f",
			FilePath:   filepath.Join("/", dirname, filename),
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

func (d DefaultService) RemoveFile(baseVolumeName, dirname, filename string) (*models.Volume, error) {
	baseVolume, err := d.Repository.GetVolume(baseVolumeName)
	if err != nil {
		return baseVolume, err
	}
	volume, err := d.CreateVolume(1)
	if err != nil {
		return nil, err
	}
	volume.Base = baseVolume.ID
	volume.FileRecords = models.FileRecords{
		&models.FileRecord{
			Opt:        "del",
			FileType:   "f",
			FilePath:   filepath.Join("/", dirname, filename),
			VolumeName: volume.Name,
		},
	}
	if volume, err = d.Repository.UpdateVolume(volume); err != nil {
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
			Opt:      "add",
			FilePath: dirname,
			FileType: "d",
		},
	}
	if volume, err = d.Repository.UpdateVolume(volume); err != nil {
		return nil, err
	}
	return volume, nil
}

func (d DefaultService) GetFile(volumeName, filename string) (*os.File, error) {
	volume, err := d.GetVolume(volumeName)
	if err != nil {
		return nil, err
	}
	return d.Storage.FetchFile(volume, filename)
}

func (d DefaultService) GetDirectory(volumeName, dirname string) (file *os.File, err error) {
	volume, err := d.GetVolume(volumeName)
	if err != nil {
		return nil, err
	}
	return d.Storage.FetchDirectory(volume, dirname)
}

func NewVolumeService(logger *zap.Logger, Storage storages.Storage, Repository repositories.Repository) Service {
	return &DefaultService{
		logger:     logger.With(zap.String("type", "Account Storage")),
		Storage:    Storage,
		Repository: Repository,
	}
}
