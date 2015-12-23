package datastore

import (
	"github.com/fhermeni/testr/model"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func (s Store) Auth(owner, repo, token string) error {
	return nil
}

func (s Store) NewOwner(u model.Owner) error {
	_, err := s.Owner(u.Name)
	if err == nil {
		return model.ErrUserExists
	}
	return s.PutOwner(u)
}
func (s Store) PutOwner(u model.Owner) error {
	k := OwnerKey(s.ctx, u.Name)
	_, err := datastore.Put(s.ctx, k, &u)
	return err
}

func (s Store) DelOwner(u model.Owner) error {
	k := OwnerKey(s.ctx, u.Name)
	u2, err := s.Owner(u.Name)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return model.ErrUserUnknown
		}
		return err
	}
	if u2.ApiToken != u.ApiToken {
		return model.ErrPermissionDenied
	}
	return datastore.Delete(s.ctx, k)
}

func (s Store) Owner(name string) (model.Owner, error) {
	k := OwnerKey(s.ctx, name)
	u := model.Owner{
		Repos: make([]string, 0, 0),
	}
	err := datastore.Get(s.ctx, k, &u)
	if err != nil {
		return u, model.ErrUserUnknown
	}
	return u, err
}

func (s Store) Authorize(owner, repos, token string) error {
	u, err := s.Owner(owner)
	if err != nil {
		return err
	}
	if u.ApiToken != token {
		return model.ErrInvalidKey
	}
	for _, r := range u.Repos {
		if r == repos && token == u.ApiToken {
			return nil
		}
	}
	return model.ErrPermissionDenied
}

func (s Store) Owners() ([]model.Owner, error) {
	owners := make([]model.Owner, 0, 0)
	_, err := datastore.NewQuery("Owner").GetAll(s.ctx, &owners)
	return owners, err
}

func OwnerKey(ctx context.Context, name string) *datastore.Key {
	return datastore.NewKey(ctx, "Owner", name, 0, nil)
}
