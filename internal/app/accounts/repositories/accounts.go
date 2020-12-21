package repositories

import (
	"github.com/infinity-oj/server-v2/pkg/models"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Repository interface {
	GetAccountByName(name string) (account *models.Account, err error)
	GetAccountById(id uint64) (account *models.Account, err error)
	CreateAccount(name, hash, salt, email string) (u *models.Account, err error)
	UpdateAccount(p *models.Account) (err error)

	UpdateCredential(u *models.Credential) (err error)
	QueryCredential(username string) (u *models.Credential, err error)

	GetRoles(accountId uint64) (roles []*models.Role, err error)
}

type DefaultRepository struct {
	logger *zap.Logger
	db     *gorm.DB
}

func (s *DefaultRepository) GetRoles(accountId uint64) (roles []*models.Role, err error) {
	s.logger.Debug("get roles",
		zap.Uint64("account id", accountId),
	)
	if err = s.db.Table("roles").Where("account_id = ?", accountId).Find(&roles).Error; err != nil {
		return nil, err
	}
	return
}

// CreateAccount
func (s *DefaultRepository) CreateAccount(username, hash, salt, email string) (account *models.Account, err error) {
	s.logger.Debug("create account",
		zap.String("username", username),
		zap.String("email", email),
	)

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	account = &models.Account{
		Name:     username,
		Nickname: username,
		Email:    email,
	}
	if err = tx.Create(account).Error; err != nil {
		tx.Rollback()
		s.logger.Error("create account", zap.String("username", username), zap.Error(err))
		return nil, errors.Wrapf(err, " create account %s ", username)
	}

	credential := &models.Credential{
		Username: username,
		Hash:     hash,
		Salt:     salt,
	}
	if err = tx.Create(credential).Error; err != nil {
		tx.Rollback()
		s.logger.Error("create account", zap.String("username", username), zap.Error(err))
		return nil, errors.Wrapf(err, " create user with username: %s", username)
	}

	return account, tx.Commit().Error
}

func (s *DefaultRepository) UpdateAccount(p *models.Account) (err error) {
	// TODO: find a better way...
	err = s.db.Save(&p).Error
	return
}

func (s *DefaultRepository) UpdateCredential(u *models.Credential) (err error) {
	// TODO: find a better way...
	err = s.db.Save(&u).Error
	return
}

func (s *DefaultRepository) QueryCredential(username string) (credential *models.Credential, err error) {
	credential = &models.Credential{}
	if err = s.db.Where(&models.Credential{Username: username}).Last(credential).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			s.logger.Error("Query credential failed", zap.String("username", username), zap.Error(err))
		}
		return nil, err
	}
	return
}

func (s *DefaultRepository) GetAccountByName(name string) (account *models.Account, err error) {
	account = &models.Account{}
	if err = s.db.Where(&models.Account{Name: name}).First(account).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			s.logger.Error("query account failed", zap.String("name", name), zap.Error(err))
		}
		return nil, err
	}
	return
}

func (s *DefaultRepository) GetAccountById(id uint64) (account *models.Account, err error) {
	account = &models.Account{}
	if err = s.db.First(account, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			s.logger.Error("query account failed", zap.Uint64("id", id), zap.Error(err))
		}
		return nil, err
	}
	return
}

func New(logger *zap.Logger, db *gorm.DB) Repository {
	return &DefaultRepository{
		logger: logger.With(zap.String("type", "Account Repository")),
		db:     db,
	}
}
