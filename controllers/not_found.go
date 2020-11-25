package controllers

import (
	"net/http"

	"github.com/github.com/steevehook/account-api/models"
	"github.com/github.com/steevehook/account-api/transport"
)

// NotFound represents the resource not found handler
func NotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := models.ResourceNotFoundError{
			Message: "route not found",
		}
		transport.SendHTTPError(w, err)
	})
}
