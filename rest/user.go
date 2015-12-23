package rest

import (
	"encoding/json"

	"github.com/fhermeni/testr/model"
)

func newOwner(c Context) error {
	var name string
	var err error
	if err = inJSON(&name, c); err != nil {
		return err
	}
	store := c.provider.NewStore(c.r)
	u := model.NewOwner(name)
	if err = store.NewOwner(u); err != nil {
		return err
	}
	c.w.Header().Add("Authorization", u.ApiToken)
	return outJSON(u, err, c)
}

func getOwner(c Context) error {
	name := c.ps["u"]
	store := c.provider.NewStore(c.r)
	u, err := store.Owner(name)
	if err != nil {
		return err
	}
	e := json.NewEncoder(c.w)
	return e.Encode(&u)
}
