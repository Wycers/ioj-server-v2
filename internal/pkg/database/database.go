package database

import (
	"github.com/google/wire"
	"github.com/infinity-oj/server-v2/pkg/models"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

// New 初始化数据库
func New(o *Options) (*gorm.DB, error) {
	var err error
	db, err := gorm.Open(postgres.Open(o.URL))
	if err != nil {
		return nil, errors.Wrap(err, "gorm open database connection error")
	}

	if o.Debug {
		db = db.Debug()
	}

	db.AutoMigrate(
		&models.Account{},
		&models.Credential{},
		&models.Problem{},
		&models.Page{},
		&models.Submission{},
		&models.Volume{},
		//&models.Role{},
		&models.Program{},
		&models.Blueprint{},
		&models.Judgement{},
		&models.Process{},
		//&models.Group{},
		//&models.UserGroupCorrelation{},
		&models.RankListRecord{},
		&models.RankListModel{},
		&models.RankList{},
	)

	return db, nil
}

var ProviderSet = wire.NewSet(New, NewOptions)
