package rest

import (
	"net/http"
	"text/template"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/dimfeld/httptreemux"
	"github.com/fhermeni/testr/model"
	"github.com/fhermeni/testr/store"
)

func NewEndPoints(s store.Provider, assets string) (Rest, error) {
	rest := Rest{
		provider: s,
		tpls:     make(map[string]*template.Template),
		Router:   httptreemux.New(),
	}
	//Compile the templates
	tpls := []string{"htmlReport", "repos"}
	for _, t := range tpls {
		tpl, err := makeTemplate(t, assets+"/"+t+".html")
		if err != nil {
			return rest, err
		}
		rest.tpls[t] = tpl
	}
	//Repos management
	rest.Router.GET("/", wrap(s, rest.tpls, index))
	rest.Router.GET("/repos/:u/:r/", wrap(s, rest.tpls, htmlReport))

	rest.Router.GET("/api/repos/:u/:r/:c", wrap(s, rest.tpls, getSuite))
	rest.Router.GET("/api/repos/:u/:r/", wrap(s, rest.tpls, getSuites))
	rest.Router.POST("/api/repos/:u/", wrap(s, rest.tpls, newRepo))
	rest.Router.DELETE("/api/repos/:u/:r", wrap(s, rest.tpls, delRepo))
	rest.Router.POST("/api/repos/:u/:r/", wrap(s, rest.tpls, postReport))

	//User management
	rest.Router.POST("/api/users/", wrap(s, rest.tpls, newOwner))
	rest.Router.GET("/api/users/:u", wrap(s, rest.tpls, getOwner))

	return rest, nil
}

type Context struct {
	w        http.ResponseWriter
	r        *http.Request
	tpls     map[string]*template.Template
	provider store.Provider
	ps       map[string]string
}

type EndPoint func(Context) error

type Rest struct {
	provider store.Provider
	tpls     map[string]*template.Template
	Router   *httptreemux.TreeMux
}

func allowed(c Context, user, repo string) error {
	token := c.r.Header.Get("Authorization")
	store := c.provider.NewStore(c.r)
	u, err := store.Owner(user)
	if err != nil {
		ctx := appengine.NewContext(c.r)
		log.Errorf(ctx, err.Error())
		return model.ErrUserUnknown
	}
	found := false
	for _, r := range u.Repos {
		if r == repo {
			found = true
			break
		}
	}
	if !found {
		return model.ErrRepoUnknown
	}
	if token != u.ApiToken {
		return model.ErrPermissionDenied
	}
	return nil
}
func wrap(provider store.Provider, tpls map[string]*template.Template, fn EndPoint) func(http.ResponseWriter, *http.Request, map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		c := Context{
			w:        w,
			r:        r,
			ps:       ps,
			provider: provider,
			tpls:     tpls,
		}
		err := fn(c)
		status(c.w, c.r, err)
	}
}
