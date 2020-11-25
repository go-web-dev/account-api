package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/github.com/steevehook/account-api/logging"
	"github.com/github.com/steevehook/account-api/models"
)

// routeParam fetches params from context and converts it into julienschmidt/httprouter.Params struct
func routeParam(r *http.Request, name string) string {
	ctx := r.Context()
	psCtx := ctx.Value(httprouter.ParamsKey)
	ps, ok := psCtx.(httprouter.Params)

	if !ok {
		logging.Logger.Error("could not extract params from context")
		return ""
	}
	return ps.ByName(name)
}

// parseBody parses JSON request body
func parseBody(r *http.Request, v interface{}) error {
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logging.Logger.Error("could not read request body")
		return models.InvalidJSONError{
			Message: "could not read request body",
		}
	}

	err = json.Unmarshal(bs, v)
	switch err.(type) {
	case *json.UnsupportedTypeError, *json.UnsupportedValueError:
		return models.InvalidJSONError{
			Message: err.Error(),
		}
	}
	return nil
}
