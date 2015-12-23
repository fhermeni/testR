package datastore

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"time"

	"github.com/fhermeni/testr/model"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type DsTestSuite struct {
	Sha1               string
	Owner              string
	Repo               string
	AuthorEmail        string
	AuthorDate         time.Time
	AuthorName         string
	AuthorAvatarUrl    string
	CommitterName      string
	CommitterDate      time.Time
	CommitterEmail     string
	CommitterAvatarUrl string
	Url                string
	Log                string
	Tests              []byte
}

func TestSuiteKey(c context.Context, user, repo, commit string) *datastore.Key {
	return datastore.NewKey(c, "Commit/"+user+"/"+repo+"/", commit, 0, nil)
}

func TestsToBytes(tests map[string]model.Test) []byte {
	buffer := new(bytes.Buffer)
	gzOut := gzip.NewWriter(buffer)
	jsonOut := json.NewEncoder(gzOut)
	jsonOut.Encode(tests)
	gzOut.Close()
	return buffer.Bytes()
}

func TestsFromBytes(buffer []byte) (map[string]model.Test, error) {
	tests := make(map[string]model.Test)
	r := bytes.NewReader(buffer)
	gzIn, err := gzip.NewReader(r)
	if err != nil {
		return tests, err
	}
	jsonIn := json.NewDecoder(gzIn)
	err = jsonIn.Decode(&tests)
	return tests, err
}

func ToDsTestSuite(owner, repo string, suite model.Commit) DsTestSuite {
	return DsTestSuite{
		Owner:              owner,
		Repo:               repo,
		Sha1:               suite.Sha1,
		Log:                suite.Log,
		Tests:              TestsToBytes(suite.Tests),
		AuthorName:         suite.Author.Name,
		AuthorEmail:        suite.Author.Email,
		AuthorDate:         suite.Author.Date,
		AuthorAvatarUrl:    suite.Author.AvatarUrl,
		CommitterName:      suite.Committer.Name,
		CommitterDate:      suite.Committer.Date,
		CommitterEmail:     suite.Committer.Email,
		CommitterAvatarUrl: suite.Committer.AvatarUrl,
		Url:                suite.Url,
	}
}

func toTestSuite(ds DsTestSuite) (model.Commit, error) {
	tests, err := TestsFromBytes(ds.Tests)
	if err != nil {
		return model.Commit{}, err
	}
	return model.Commit{
		Sha1:  ds.Sha1,
		Url:   ds.Url,
		Log:   ds.Log,
		Tests: tests,
		Owner: ds.Owner,
		Repo:  ds.Repo,
		Author: model.Author{
			Name:      ds.AuthorName,
			Email:     ds.AuthorEmail,
			Date:      ds.AuthorDate,
			AvatarUrl: ds.AuthorAvatarUrl,
		},
		Committer: model.Author{
			Name:      ds.CommitterName,
			Email:     ds.CommitterEmail,
			Date:      ds.CommitterDate,
			AvatarUrl: ds.CommitterAvatarUrl,
		},
	}, nil
}

func (s Store) PutTestSuite(user, repo string, suite model.Commit) error {
	//Put the commit
	ds := ToDsTestSuite(user, repo, suite)
	k := TestSuiteKey(s.ctx, user, repo, suite.Sha1)
	_, err := datastore.Put(s.ctx, k, &ds)
	return err
}

func (s Store) TestSuite(user, repo, commit string) (model.Commit, error) {
	k := TestSuiteKey(s.ctx, user, repo, commit)
	ds := DsTestSuite{}
	if err := datastore.Get(s.ctx, k, &ds); err != nil {
		return model.Commit{}, err
	}
	return toTestSuite(ds)
}

func (s Store) TestSuites(user, repo string, nb int) ([]model.Commit, error) {
	e := "Commit/" + user + "/" + repo + "/"
	q := datastore.NewQuery(e).Limit(nb).Order("CommitterDate")
	suites := make([]model.Commit, 0, 0)
	for t := q.Run(s.ctx); ; {
		var ds DsTestSuite
		_, err := t.Next(&ds)
		if err == datastore.Done {
			break
		} else if err != nil {
			return []model.Commit{}, err
		}
		ts, err := toTestSuite(ds)
		if err != nil {
			return []model.Commit{}, nil
		}
		suites = append(suites, ts)
	}
	return suites, nil
}
