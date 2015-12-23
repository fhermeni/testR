package rest

import (
	"net/http"

	"github.com/fhermeni/testr/model"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func status(w http.ResponseWriter, r *http.Request, e error) {
	switch e {
	case model.ErrUserUnknown, model.ErrRepoUnknown:
		http.Error(w, e.Error(), http.StatusNotFound)
		return
	case model.ErrUserExists, model.ErrRepoExists:
		http.Error(w, e.Error(), http.StatusConflict)
		return
	case model.ErrInvalidKey:
		http.Error(w, e.Error(), http.StatusUnauthorized)
	case model.ErrPermissionDenied:
		http.Error(w, e.Error(), http.StatusForbidden)
		return
	case nil:
		return
	default:
		http.Error(w, "Internal server error. A possible bug to report", http.StatusInternalServerError)
		ctx := appengine.NewContext(r)
		log.Errorf(ctx, "Unsupported error: %s", e.Error())
	}
}
