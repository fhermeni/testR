package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fhermeni/testr/drivers"
	"github.com/fhermeni/testr/model"
	"google.golang.org/appengine"
)

func htmlReport(c Context) error {
	owner, _ := c.ps["u"]
	repo, _ := c.ps["r"]
	ctx := appengine.NewContext(c.r)
	if clientCached(ctx, c.w, c.r, owner+"/"+repo) {
		return nil
	}
	store := c.provider.NewStore(c.r)
	suites, err := store.TestSuites(owner, repo, 10)
	if err != nil {
		return err
	}
	report := model.Reportify(owner, repo, suites)
	return outHTML("htmlReport", owner+"/"+repo, report, err, c)
}

type RepoSummary struct {
	Owner     model.Owner
	Summaries map[string]*model.Summary
}

func getSuite(c Context) error {
	name := c.ps["u"]
	repo, _ := c.ps["r"]
	commit, _ := c.ps["c"]
	store := c.provider.NewStore(c.r)
	ts, err := store.TestSuite(name, repo, commit)
	return outJSON(&ts, err, c)
}

func getSuites(c Context) error {
	name := c.ps["u"]
	repo, _ := c.ps["r"]
	amount := c.r.URL.Query().Get("last")
	nb, err := strconv.Atoi(amount)
	if err != nil || nb <= 0 {
		nb = 10
	}
	store := c.provider.NewStore(c.r)
	ts, err := store.TestSuites(name, repo, nb)
	return outJSON(&ts, err, c)
}

func newRepo(c Context) error {
	name := c.ps["u"]
	d := json.NewDecoder(c.r.Body)
	var repo string
	if err := d.Decode(&repo); err != nil {
		return err
	}

	if err := allowed(c, name, repo); err != nil && err != model.ErrRepoUnknown {
		return err
	}
	store := c.provider.NewStore(c.r)

	u, err := store.Owner(name)
	if err != nil {
		return err
	}
	for _, r := range u.Repos {
		if r == repo {
			return model.ErrRepoExists
		}
	}
	repos := u.Repos
	repos = append(repos, repo)
	u.Repos = repos
	if err := store.PutOwner(u); err != nil {
		return err
	}
	c.w.WriteHeader(http.StatusCreated)
	return nil
}

func delRepo(c Context) error {
	name := c.ps["u"]
	repo, _ := c.ps["r"]
	store := c.provider.NewStore(c.r)

	u, err := store.Owner(name)
	if err != nil {
		return err
	}
	repos := make([]string, 0, 0)
	for _, r := range u.Repos {
		if r != repo {
			repos = append(repos, r)
		}
	}
	u.Repos = repos
	return store.PutOwner(u)
}

type Deposit struct {
	Commit  model.Commit
	Content string
	Type    string
}

func postReport(c Context) error {
	owner, _ := c.ps["u"]
	repo, _ := c.ps["r"]
	if err := allowed(c, owner, repo); err != nil {
		return err
	}
	rawReport := Deposit{}
	if err := inJSON(&rawReport, c); err != nil {
		return err
	}
	st := c.provider.NewStore(c.r)
	tests, err := driver.Testify(rawReport.Type, string(rawReport.Content))
	if err != nil {
		return err
	}
	co := rawReport.Commit
	co.Tests = tests
	co.Owner = owner
	co.Repo = repo
	if in, err := st.TestSuite(owner, repo, co.Sha1); err == nil {
		in.Fullfill(co)
		co = in
	}
	if err := st.PutTestSuite(owner, repo, co); err != nil {
		return err
	}
	c.w.WriteHeader(http.StatusCreated)
	clearCache(appengine.NewContext(c.r), owner+"/"+repo)
	return nil
}
