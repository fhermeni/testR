package store

import (
	"net/http"

	"github.com/fhermeni/testr/model"
)

//Store specifies the behavior of the datastore
type Store interface {
	//PutTestSuite save a test suite for a given repository
	//If the commit is already known, the existing suite is completed with the yet unknown tests
	//otherwise, a new suite is stored
	PutTestSuite(owner, repo string, suite model.Commit) error
	//TestSuite returns a given TestSuite
	TestSuite(owner, repo, commit string) (model.Commit, error)
	//TestSuites returns the last TestSuites
	TestSuites(owner, repo string, last int) ([]model.Commit, error)
	//Auth indicates if the provided token match a user
	Auth(owner, repo, token string) error
	Authorize(owner, repos, token string) error

	//User management
	PutOwner(u model.Owner) error
	NewOwner(u model.Owner) error
	Owner(n string) (model.Owner, error)
	DelOwner(u model.Owner) error
}

//Provider specifies a datastore provider
type Provider interface {
	//NewStore returns a store in the context of a given request
	NewStore(req *http.Request) Store
}
