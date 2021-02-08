package storages

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/infinity-oj/server-v2/pkg/models"

	"github.com/infinity-oj/server-v2/internal/pkg/files"
	"go.uber.org/zap"
)

type Storage interface {
	CreateDirectory(volume, directory string) error
	CreateFile(volume, fileName string, data []byte) error
	IsFileExists(volume, fileName string) bool
	FetchFile(volume, fileName string) ([]byte, error)

	ArchiveDirectory(volume, directory string) (file *os.File, err error)
	ArchiveVolume(volume *models.Volume) (file *os.File, err error)
}

type FileManager struct {
	logger *zap.Logger
	fm     files.FileManager
}

func (m *FileManager) ArchiveDirectory(volume, directory string) (file *os.File, err error) {
	filePath := path.Join(volume, directory)
	return m.fm.ArchiveDirectory(filePath)
}

func (m *FileManager) ArchiveVolume(volume *models.Volume) (file *os.File, err error) {
	//filePath := path.Join(volume.Name, directory)

	zipfile, err := ioutil.TempFile("", "*.zip")
	if err != nil {
		return nil, err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	baseDir := volume.Name

	for _, fileRecord := range volume.FileRecords {
		info, err := func() (os.FileInfo, error) {
			if fileRecord.IsDir() {
				return fileRecord, nil
			} else {
				filePath := path.Join(fileRecord.VolumeName, fileRecord.VolumePath)
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

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return nil, err
		}

		if info.IsDir() {
			continue
		}

		fileBytes, err := m.FetchFile(fileRecord.VolumeName, fileRecord.VolumePath)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(writer, bytes.NewReader(fileBytes))
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	return zipfile, nil
}

func (m *FileManager) IsFileExists(volume, fileName string) bool {
	filePath := path.Join(volume, fileName)
	exist, err := m.fm.IsFileExists(filePath)
	if err != nil {
		return false
	}
	return exist
}

func (m *FileManager) FetchFile(volume, fileName string) ([]byte, error) {
	filePath := path.Join(volume, fileName)
	return m.fm.FetchFile(filePath)
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
