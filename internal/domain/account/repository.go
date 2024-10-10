package account

import (
	"errors"

	"github.com/rubberduckkk/ducker/internal/domain/account/entity"
	"github.com/rubberduckkk/ducker/internal/domain/account/valueobj"
)

var (
	ErrUnsupportedOp = errors.New("unsupported operation")
	ErrUnauthorized  = errors.New("unauthorized account")
)

type Repository interface {
	Create(info valueobj.AccountInfo) (*entity.Account, error)
	Auth(*entity.Account) error
}
