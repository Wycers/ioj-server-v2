package accounts

import (
	"errors"
	"github.com/infinity-oj/server-v2/internal/pkg/crypto"
	"github.com/infinity-oj/server-v2/internal/pkg/utils/random"
	"github.com/infinity-oj/server-v2/pkg/models"
	"go.uber.org/zap"
)

// You should never change it, otherwise all old credentials will turn invalid
const specialKey = "imf1nlTy0j"

type Service interface {
	GetAccount(name string) (account *models.Account, err error)
	GetAccountById(id uint64) (account *models.Account, err error)
	UpdateAccount(account *models.Account, nickname, email, gender, locale string) (*models.Account, error)
	CreateAccount(username, password, email string) (account *models.Account, err error)

	UpdateCredential(username, oldPassword, newPassword string) (res bool, err error)
	VerifyCredential(username, password string) (isValid bool, err error)

	GetRoleById(accountId uint64) (roles []*models.Role, err error)
}

type DefaultService struct {
	logger     *zap.Logger
	Repository Repository
}

func (s *DefaultService) UpdateCredential(username, oldPassword, newPassword string) (res bool, err error) {
	s.logger.Debug("verify credential", zap.String("username", username))
	u := new(models.Credential)
	if u, err = s.Repository.QueryCredential(username); err != nil {
		s.logger.Error("verify credential error", zap.Error(err))
		return false, err
	}
	if u == nil {
		return false, nil
	}
	hash := crypto.Sha256(u.Salt + oldPassword + specialKey)

	if hash != u.Hash {
		return false, nil
	}

	u.Salt = random.RandStringRunes(64)
	u.Hash = crypto.Sha256(u.Salt + newPassword + specialKey)

	if err := s.Repository.UpdateCredential(u); err != nil {
		return false, err
	}

	return true, nil
}

func (s *DefaultService) GetAccount(name string) (account *models.Account, err error) {
	account, err = s.Repository.GetAccountByName(name)
	return
}

func (s *DefaultService) GetAccountById(id uint64) (account *models.Account, err error) {
	account, err = s.Repository.GetAccountById(id)
	return
}

func (s *DefaultService) GetRoleById(accountId uint64) (roles []*models.Role, err error) {
	roles, err = s.Repository.GetRoles(accountId)
	return
}

func (s *DefaultService) UpdateAccount(account *models.Account, nickname, email, gender, locale string) (*models.Account, error) {
	s.logger.Debug("update account", zap.String("name", account.Name))
	account.Nickname = nickname
	account.Email = email
	account.Gender = gender
	account.Locale = locale
	err := s.Repository.UpdateAccount(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *DefaultService) CreateAccount(username, password, email string) (account *models.Account, err error) {
	s.logger.Debug("create account", zap.String("username", username))
	{
		account, err := s.GetAccount(username)
		if err != nil {
			return nil, err
		}
		if account != nil {
			return nil, errors.New("username exists")
		}
	}

	salt := random.RandStringRunes(64)
	hash := crypto.Sha256(salt + password + specialKey)
	if account, err = s.Repository.CreateAccount(username, hash, salt, email); err != nil {
		return nil, err
	}
	return
}

func (s *DefaultService) VerifyCredential(username, password string) (isValid bool, err error) {
	s.logger.Debug("verify credential", zap.String("username", username))
	u := new(models.Credential)
	if u, err = s.Repository.QueryCredential(username); err != nil {
		s.logger.Error("verify credential error", zap.Error(err))
		return false, err
	}
	if u == nil {
		return false, nil
	}
	hash := crypto.Sha256(u.Salt + password + specialKey)

	return hash == u.Hash, nil
}
func NewService(logger *zap.Logger, Repository Repository) Service {
	return &DefaultService{
		logger:     logger.With(zap.String("type", "Account Repository")),
		Repository: Repository,
	}
}
