package rest

import (
	"compress/gzip"
	"encoding/json"
	"net/http"
	"sync"
)

// Create a Pool that contains previously used Writers and
// can create new ones if we run out.
var zippers = sync.Pool{New: func() interface{} {
	return gzip.NewWriter(nil)
}}

func outJSON(j interface{}, e error, c Context) error {
	if e != nil {
		return e
	}
	c.w.Header().Set("Content-type", "application/json; charset=utf-8")
	enc := json.NewEncoder(c.w)
	return enc.Encode(j)
}

func inJSON(j interface{}, c Context) error {
	dec := json.NewDecoder(c.r.Body)
	err := dec.Decode(&j)
	if err != nil {
		http.Error(c.w, "Bad json: "+err.Error(), http.StatusBadRequest)
	}
	return err
}
