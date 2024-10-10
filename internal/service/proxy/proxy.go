package proxy

import (
	"net/http"

	"github.com/rubberduckkk/ducker/internal/domain/account"
)

type Service interface {
	ProxyHTTP(w http.ResponseWriter, r *http.Request)
}

type service struct {
	accounts account.Repository
}

func New(accounts account.Repository) Service {
	return &service{accounts: accounts}
}
