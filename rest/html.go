package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"google.golang.org/appengine"

	"github.com/fhermeni/testr/model"
)

//Should be cached
func makeTemplate(key, path string) (*template.Template, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	cnt, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	funcMap := template.FuncMap{
		"abbrv":       abbrv,
		"statusClass": statusClass,
		"duration":    duration,
		"shorterName": shorterName,
		"pct":         pct,
		"suiteStatus": suiteStatus,
		"diff":        diff,
		"testLabels":  testsLabels,
		"testLabel":   testLabel,
	}

	return template.New(key).Funcs(funcMap).Parse(string(cnt))
}

//outHTML render a page and encache the generation timestamp for a client caching based
//on the Last-Modified header
func outHTML(page, key string, j interface{}, e error, c Context) error {
	if e != nil {
		return e
	}
	now := time.Now()
	c.w.Header().Add("Cache-Control", "public, max-age=31536000")
	c.w.Header().Add("Last-Modified", now.Format(http.TimeFormat))
	c.w.Header().Set("Content-type", "text/html; charset=utf-8")
	tpl, _ := c.tpls[page]

	clientEncache(appengine.NewContext(c.r), key, now)
	return tpl.Execute(c.w, j)
}

func suiteStatus(ts model.Commit) string {
	s := ts.Summary()
	if s.Success == s.Total {
		return "success"
	} else if s.Failures > 0 {
		return "danger"
	}
	return "default"
}

func duration(d float64) string {
	if d < 1 {
		return fmt.Sprintf("%.fms", d*1000)
	}
	return fmt.Sprintf("%.3fs", d)
}

func statusClass(tests model.Tests, test model.Test) string {
	n := test.Status
	if n < 2 {
		return "default"
	} else if n == 2 {
		//success. Too long or not ?
		if tests.Faster(test.Duration) {
			return "success"
		} else if tests.Slower(test.Duration) {
			return "warning"
		}
		return "default"
	}
	return "danger"
}

func diff(a, b int) string {
	d := b - a
	if d == 0 {
		return "-"
	}

	return fmt.Sprintf("%d", d)
}

func testLabel(test model.Commit, name string) string {
	t, ok := test.Tests[name]
	if !ok {
		return ""
	}
	if t.Status == model.SUCCESS {
		return "is:success"
	} else if t.Status == model.SKIPPED {
		return "is:skip"
	} else if t.Status == model.FAILURE {
		return "is:failure"
	}
	return ""
}

func testsLabels(tests model.Tests) string {
	success := false
	skip := false
	failure := false
	for _, t := range tests {
		if t.Status == model.SUCCESS {
			success = true
		} else if t.Status == model.SKIPPED {
			skip = true
		} else if t.Status == model.FAILURE {
			failure = true
		}
	}
	labels := ""
	if success {
		labels += " was:success"
	}
	if skip {
		labels += " was:skip"
	}
	if failure {
		labels += " was:failure"
	}
	return labels
}

func abbrv(commit string) string {
	return commit[0:6]
}

func pct(a, b int) string {
	res := float64(a) / float64(b) * float64(100)
	return fmt.Sprintf("%.1f", res)
}

func shorterName(n string) string {
	if len(n) > 30 {
		pos := strings.LastIndex(n, ".")
		end := len(n)
		return "..." + n[pos+1:end]
	}
	return n + "."
}
