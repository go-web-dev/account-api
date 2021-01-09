package controllers

import (
	"context"
	"net/http"

	"github.com/lestrrat-go/jwx/jwk"

	"github.com/github.com/steevehook/account-api/transport"
)

type keySetGetter interface {
	GetKeySet(ctx context.Context) (jwk.Set, error)
}

func getKeySet(service keySetGetter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keySet, err := service.GetKeySet(r.Context())
		if err != nil {
			transport.SendHTTPError(w, err)
			return
		}
		transport.SendJSON(w, http.StatusOK, keySet)
	})
}
