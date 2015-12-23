package datastore

import (
	"net/http"

	"golang.org/x/net/context"

	"github.com/fhermeni/testr/store"
	"google.golang.org/appengine"
)

/*
//Store specifies the behavior of the datastore
type Store interface {
	//PutTestSuite save a test suite for a given repository
	//If the commit is already known, the existing suite is completed with the yet unknown tests
	//otherwise, a new suite is stored
	PutTestSuite(user, repos string, suite model.TestSuite) error
	//TestSuite returns a given TestSuite
	TestSuite(user, repos, commit string) (model.TestSuite, error)
	//TestSuites returns the last TestSuites
	TestSuites(user, repos string, last int) ([]model.TestSuite, error)
	//Auth indicates if the provided token match a user
	Auth(user, repos, token string) error
	Authorize(user, repos, token string) error
}
*/

type Store struct {
	ctx context.Context
}

type Provider struct {
	apiKey string
}

func (p Provider) NewStore(r *http.Request) store.Store {
	return Store{ctx: appengine.NewContext(r)}
}

func NewProvider() store.Provider {
	return &Provider{}
}
