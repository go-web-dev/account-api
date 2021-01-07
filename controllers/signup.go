package controllers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/github.com/steevehook/account-api/logging"
	"github.com/github.com/steevehook/account-api/models"
	"github.com/github.com/steevehook/account-api/transport"
)

type signupper interface {
	Signup(models.Credentials) (models.TokenResponse, error)
}

func signup(service signupper) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logging.Logger
		var credentials models.Credentials
		err := parseBody(r, &credentials)
		if err != nil {
			logger.Error("could not parse request", zap.Error(err))
			transport.SendHTTPError(w, err)
			return
		}

		res, err := service.Signup(credentials)
		if err != nil {
			logger.Error("could not signup the current user", zap.Error(err))
			transport.SendHTTPError(w, err)
			return
		}
		logger.Info("successfully signed up the user")
		transport.SendJSON(w, http.StatusOK, res)
	})
}
