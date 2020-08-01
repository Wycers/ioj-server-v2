package database

import (
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/infinity-oj/server-v2/internal/pkg/models"
)

// Options is  configuration of database
type Options struct {
	URL   string `yaml:"url"`
	Debug bool
}

func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var err error
	o := new(Options)
	if err = v.UnmarshalKey("db", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal db option error")
	}

	logger.Info("load database options success", zap.String("url", o.URL))

	return o, err
}

// Init 初始化数据库
func New(o *Options) (*gorm.DB, error) {
	var err error
	db, err := gorm.Open("postgres", o.URL)
	if err != nil {
		return nil, errors.Wrap(err, "gorm open database connection error")
	}

	if o.Debug {
		db = db.Debug()
	}

	db.AutoMigrate(
		&models.Credential{},
		&models.Account{},
		&models.Page{},
		&models.Problem{},
		&models.Judgement{},
		&models.Submission{},
	)

	return db, nil
}

var ProviderSet = wire.NewSet(New, NewOptions)
