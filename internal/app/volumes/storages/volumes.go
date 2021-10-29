package storages

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/infinity-oj/server-v2/pkg/models"

	"github.com/infinity-oj/server-v2/internal/pkg/files"
	"go.uber.org/zap"
)

type Storage interface {
	CreateDirectory(volume, directory string) error
	CreateFile(volume, fileName string, data []byte) error
	IsFileExists(volume, fileName string) bool

	FetchFile(volume *models.Volume, fileName string) (*os.File, error)
	FetchDirectory(volume *models.Volume, directory string) (file *os.File, err error)
}

type FileManager struct {
	logger *zap.Logger
	fm     files.FileManager
}

func (m *FileManager) FetchFile(volume *models.Volume, fileName string) (*os.File, error) {
	tmpFile, err := ioutil.TempFile("", volume.Name+"*")
	if err != nil {
		m.logger.Error("new temp file error", zap.Error(err))
		return nil, err
	}
	defer tmpFile.Close()

	for _, fileRecord := range volume.FileRecords {
		m.logger.Debug("file record", zap.Any("fr", fileRecord), zap.Any("fn", fileName))
		if fileRecord.FilePath != fileName {
			continue
		}

		fileBytes, err := m.fm.FetchFile(filepath.Join(fileRecord.VolumeName, fileRecord.VolumePath))
		if err != nil {
			return nil, err
		}
		err = ioutil.WriteFile(tmpFile.Name(), fileBytes, 0644)
		if err != nil {
			m.logger.Error("write file error", zap.Error(err))
			return nil, err
		}
		return tmpFile, nil
	}
	return nil, errors.New("not found")
}

func (m *FileManager) FetchDirectory(volume *models.Volume, directory string) (file *os.File, err error) {
	zipFile, err := ioutil.TempFile("", "*.zip")
	if err != nil {
		return nil, err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	filePath := path.Join(volume.Name, directory)
	baseDir := filepath.Base(filePath)

	for _, fileRecord := range volume.FileRecords {
		fmt.Println(fileRecord)

		info, err := func() (os.FileInfo, error) {
			if fileRecord.IsDir() {
				return fileRecord, nil
			} else {
				filePath := path.Join(fileRecord.VolumeName, fileRecord.VolumePath)
				filePath = filepath.ToSlash(filePath)
				info, err := m.fm.FetchFileInfo(filePath)
				if err != nil {
					return nil, err
				}
				return info, err
			}
		}()
		if err != nil {
			return nil, err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return nil, err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, fileRecord.FilePath)
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		header.Name = filepath.ToSlash(header.Name)

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return nil, err
		}

		if info.IsDir() {
			continue
		}
		fileBytes, err := m.fm.FetchFile(filepath.Join(fileRecord.VolumeName, fileRecord.VolumePath))
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(writer, bytes.NewReader(fileBytes))
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	return zipFile, nil
}

func (m *FileManager) IsFileExists(volume, fileName string) bool {
	filePath := path.Join(volume, fileName)
	exist, err := m.fm.IsFileExists(filePath)
	if err != nil {
		return false
	}
	return exist
}

func (m *FileManager) CreateDirectory(volume, directory string) error {
	filePath := path.Join(volume, directory)
	return m.fm.CreateDirectory(filePath)
}

func (m *FileManager) CreateFile(volume, fileName string, data []byte) error {
	filePath := path.Join(volume, fileName)
	return m.fm.CreateFile(filePath, data)
}

func NewFileManager(logger *zap.Logger, fm files.FileManager) Storage {
	return &FileManager{
		logger: logger.With(zap.String("type", "Storage")),
		fm:     fm,
	}
}
