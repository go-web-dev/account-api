package controllers

import (
	"net/http"
)

type logoutter interface {
	Logout()
}

func logout(service logoutter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
