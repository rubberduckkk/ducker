package proxy

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/rubberduckkk/ducker/internal/domain/account/entity"
)

func (s *service) auth(username, password string) error {
	acct := entity.NewAccount(username, password)
	if err := s.accounts.Auth(acct); err != nil {
		return err
	}
	return nil
}

func (s *service) authHTTP(req *http.Request) bool {
	if req == nil {
		return false
	}
	header := req.Header.Get("Authorization")
	if !strings.HasPrefix(header, "Basic ") {
		return false
	}

	raw := strings.TrimPrefix(header, "Basic ")
	credentials, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		logrus.WithError(err).WithField("header", header).Errorf("basic auth decode failed")
		return false
	}

	parts := strings.Split(string(credentials), ":")
	if len(parts) != 2 {
		logrus.WithField("cred", credentials).Errorf("expected 2 parts, got %d", len(parts))
		return false
	}

	username := parts[0]
	password := parts[1]
	if err = s.auth(username, password); err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"username": username,
			"password": password,
		}).Errorf("auth failed")
		return false
	}
	return true
}
