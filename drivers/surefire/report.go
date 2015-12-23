package surefire

//package surefire is a driver for maven surefire test reports
import (
	"encoding/xml"
	"io"
	"strings"

	"github.com/fhermeni/testr/model"
)

//Property store an environment property
type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

//Failure stores data related to a test failure
type Failure struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Output  string `xml:",chardata"`
}

//TestCase denotes well, a testcase
type TestCase struct {
	//Name is the name of the test
	Name string `xml:"name,attr"`
	//ClassName is the class providing the test
	ClassName string `xml:"classname,attr"`
	//Time is the test duration in seconds
	Time float64 `xml:"time,attr"`
	//Failure is an optional data denoting the failure message
	Failure *Failure `xml:"failure,omitempty"`
}

//TestSuite denotes a collection of surefire tests
type TestSuite struct {
	Name       string     `xml:"name,attr"`
	Time       float64    `xml:"time,attr"`
	Tests      int        `xml:"tests,attr"`
	Errors     int        `xml:"errors,attr"`
	Skipped    int        `xml:"skipped,attr"`
	Failures   int        `xml:"failures,attr"`
	Properties []Property `xml:"properties>property"`
	TestCases  []TestCase `xml:"testcase"`
}

//Decode the surefire XML report
func Decode(in io.Reader) (TestSuite, error) {
	d := xml.NewDecoder(in)
	var s TestSuite
	err := d.Decode(&s)
	return s, err
}

//Surefire driver
type Surefire struct {
	Content TestSuite
}

//New makes a new surefire driver
func New(cnt string) (Surefire, error) {
	s := Surefire{}
	c, err := Decode(strings.NewReader(cnt))
	s.Content = c
	return s, err
}

//Testify convert the surefire test to the pivot model
//the resulting testname is the concatenation of the classname and the test method
//upon failure, the output message if the surefire message then the output, separated by two \n
func (s *Surefire) Testify() (map[string]model.Test, error) {
	tests := make(map[string]model.Test)
	for _, r := range s.Content.TestCases {
		t := model.Test{
			Name:     r.ClassName + "#" + r.Name,
			Duration: r.Time,
			Status:   model.SUCCESS,
		}
		if r.Failure != nil {
			t.Status = model.FAILURE
			t.Output = r.Failure.Message + "\n\n" + r.Failure.Output
		}
		tests[t.Name] = t
	}
	return tests, nil
}
