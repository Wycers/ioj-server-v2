package repositories

import (
	"os"
	"path"

	"github.com/infinity-oj/server-v2/internal/pkg/files"
	"go.uber.org/zap"
)

type Repository interface {
	CreateDirectory(volume, directory string) error
	CreateFile(volume, fileName string, data []byte) error
	IsFileExists(volume, fileName string) bool
	FetchFile(volume, fileName string) ([]byte, error)

	ArchiveDirectory(volume, directory string) (file *os.File, err error)
}

type FileManager struct {
	logger *zap.Logger
	fm     files.FileManager
}

func (m *FileManager) ArchiveDirectory(volume, directory string) (file *os.File, err error) {
	filePath := path.Join(volume, directory)
	return m.fm.ArchiveDirectory(filePath)
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

func NewFileManager(logger *zap.Logger, fm files.FileManager) Repository {
	return &FileManager{
		logger: logger.With(zap.String("type", "Repository")),
		fm:     fm,
	}
}
