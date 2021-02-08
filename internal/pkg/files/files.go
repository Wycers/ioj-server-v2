package files

import (
	"errors"
	"os"
	"strings"

	"github.com/google/wire"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Options is log configuration struct
type Options struct {
	Type string `yaml:"type"`
	Base string `yaml:"base"`
}

type FileManager interface {
	SetBase(base string)
	GetBase() string
	CreateFile(fileName string, bytes []byte) error
	CreateDirectory(fileName string) error
	FetchFile(fileName string) ([]byte, error)
	FetchFileInfo(fileName string) (os.FileInfo, error)
	IsFileExists(fileName string) (bool, error)
	IsDirectoryExists(fileName string) (bool, error)
	GetFilesAndDirs(dirname string) ([]string, []string, error)

	ArchiveDirectory(fileName string) (file *os.File, err error)
}

func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("volumes", o); err != nil {
		return nil, err
	}

	logger.Info("load file manager configuration success")

	return o, err
}

// New for file library
func New(o *Options) (FileManager, error) {
	switch strings.ToLower(o.Type) {
	case "local":
		{
			fm := &LocalFileManager{}
			fm.SetBase(o.Base)
			return fm, nil
		}
	default:
		return nil, errors.New("unknown store type: " + o.Type)
	}
}

var ProviderSet = wire.NewSet(New, NewOptions)
