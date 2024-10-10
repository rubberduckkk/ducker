package account

import (
	"github.com/rubberduckkk/ducker/internal/domain/account"
	"github.com/rubberduckkk/ducker/internal/domain/account/entity"
	"github.com/rubberduckkk/ducker/internal/domain/account/valueobj"
	"github.com/rubberduckkk/ducker/internal/infra/config"
	"github.com/rubberduckkk/ducker/pkg/hash"
)

type accountRepo struct {
	cfg *config.Config
}

func NewRepo(cfg *config.Config) account.Repository {
	return &accountRepo{
		cfg: cfg,
	}
}

func (a *accountRepo) Create(info valueobj.AccountInfo) (*entity.Account, error) {
	return nil, account.ErrUnsupportedOp
}

func (a *accountRepo) Auth(e *entity.Account) error {
	for _, pass := range a.cfg.Account.Passes {
		if err := hash.VerifyPassword(pass, e.Password); err == nil {
			return nil
		}
	}

	return account.ErrUnauthorized
}
