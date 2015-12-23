package model

import (
	"crypto/rand"
	"errors"
	"log"
	"strings"
	"time"
)

const (
	UNKOWN = iota
	SKIPPED
	SUCCESS
	FAILURE
)

var (
	ErrUserUnknown      = errors.New("Unknown owner")
	ErrRepoUnknown      = errors.New("Unknown repository")
	ErrUserExists       = errors.New("The owner already exists")
	ErrRepoExists       = errors.New("The repository already exists")
	ErrPermissionDenied = errors.New("Permission denied")
	ErrInvalidKey       = errors.New("invalid key")
)

//Author denotes a person involved in a commit
type Author struct {
	//Name is the author fullname
	Name string
	//Email of the author
	Email string
	//Date the action occurred
	Date time.Time
	//AvatarUrl is a link to the author avatar img
	AvatarUrl string
}

type Commit struct {
	Owner     string
	Repo      string
	Sha1      string
	Author    Author
	Committer Author
	Url       string
	Tests     map[string]Test
	Log       string
}

type Owner struct {
	Name     string
	ApiToken string `json:"-"`
	Repos    []string
}

func NewOwner(name string) Owner {
	return Owner{
		Name:     name,
		Repos:    make([]string, 0, 0),
		ApiToken: string(randomBytes(16)),
	}
}

func randomBytes(s int) []byte {
	alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, s)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return bytes
}

func (c Commit) ShortLog() string {
	newLine := strings.Index(c.Log, "\n")
	dot := strings.Index(c.Log, ".")
	pos := 80
	if len(c.Log) < 80 {
		pos = len(c.Log)
	}
	if dot > 0 && dot < pos {
		pos = dot
	}
	if newLine > 0 && newLine < pos {
		pos = newLine
	}
	log.Println(pos)
	if pos < 80 {
		return c.Log[0:pos]
	}
	return c.Log[0:pos] + "..."
}

type Test struct {
	Name     string
	Status   int
	Duration float64
	Output   string
}

type TestSuite struct {
	Commit Commit
	Tests  map[string]Test
}

type Report struct {
	Owner   string
	Repo    string
	Commits []Commit
	Tests   map[string]Tests
	Last    Commit
	Prev    *Commit
}

type Tests []Test

func (tests Tests) Avg() float64 {
	avg := 0.0
	nb := 0
	for _, t := range tests {
		if t.Status == SUCCESS {
			nb++
			avg += t.Duration
		}
	}
	return avg / float64(nb)
}

func (tests Commit) Fullfill(ts Commit) {
	//Complete an existing one
	for n, ts := range ts.Tests {
		if _, ok := tests.Tests[n]; !ok {
			tests.Tests[n] = ts
		}
	}
}

func (tests Tests) Stddev() float64 {
	avg := tests.Avg()
	stddev := 0.0
	nb := 0
	for _, t := range tests {
		if t.Status == SUCCESS {
			stddev += ((t.Duration - avg) * (t.Duration - avg))
			nb++
		}
	}
	return stddev / float64(nb)
}

type Summary struct {
	Total    int
	Success  int
	Slowdown int
	Speedup  int
	Failures int
	Skip     int
}

//Summary returns some numbers
//total, success, slowdown, speedup, failures, skip
func (ts Commit) Summary() Summary {
	s := Summary{}

	for _, t := range ts.Tests {
		if t.Status == SUCCESS {
			s.Success++
		} else if t.Status == FAILURE {
			s.Failures++
		} else if t.Status == SKIPPED {
			s.Skip++
		} else {
			log.Println("unknown")
		}
		s.Total++
	}
	return s
}

func (s Summary) Pct(kind string) int {
	switch kind {
	case "success":
		return s.Success / s.Total * 100
	case "failures":
		return s.Failures / s.Total * 100
	case "skip":
		return s.Skip / s.Total * 100
	case "slowdown":
		return s.Slowdown / s.Success * 100
	case "speedup":
		return s.Speedup / s.Success * 100
	}
	return s.Success / s.Total * 100
}

func (tests Tests) Slower(d float64) bool {
	stddev := tests.Stddev()
	if stddev < 0.001 {
		return false
	}
	return d > tests.Avg()-stddev
}

func (tests Tests) Faster(d float64) bool {
	stddev := tests.Stddev()
	if stddev < 0.001 {
		return false
	}
	return d < tests.Avg()+stddev
}

func Reportify(owner, repo string, suites []Commit) Report {
	r := Report{
		Owner: owner,
		Repo:  repo,
		Tests: make(map[string]Tests),
	}
	commits := make([]Commit, 0, 0)
	for i, ts := range suites {
		commits = append(commits, ts)
		for _, t := range ts.Tests {
			n := t.Name
			trend, ok := r.Tests[n]
			if !ok {
				trend = make(Tests, len(suites), len(suites))
			}
			trend[i] = t
			r.Tests[n] = trend
		}
	}
	if len(suites) > 0 {
		r.Last = suites[len(suites)-1]
	}
	r.Commits = commits
	if len(suites) >= 2 {
		r.Prev = &suites[len(suites)-2]
	}
	return r
}
