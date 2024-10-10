package account

import (
	"github.com/rubberduckkk/ducker/internal/domain/account"
	"github.com/rubberduckkk/ducker/internal/domain/account/entity"
	"github.com/rubberduckkk/ducker/internal/domain/account/valueobj"
)

type Service interface {
	Create(username, password, email string) (*entity.Account, error)
	Auth(username, password string) error
}

type service struct {
	accounts account.Repository
}

func New(accounts account.Repository) Service {
	return &service{accounts: accounts}
}

func (s *service) Create(username, password, email string) (*entity.Account, error) {
	info := valueobj.AccountInfo{
		Username: username,
		Password: password,
		Email:    email,
	}
	return s.accounts.Create(info)
}

func (s *service) Auth(username, password string) error {
	return s.accounts.Auth(&entity.Account{Username: username, Password: password})
}
