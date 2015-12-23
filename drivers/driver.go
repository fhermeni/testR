package driver

import (
	"errors"

	"github.com/fhermeni/testr/drivers/surefire"
	"github.com/fhermeni/testr/model"
)

//Driver specifies a conversion of a vendor specific testing report to our pivot model
type Driver interface {
	Testify() ([]model.Test, error)
}

//Testify is a factory that call the right driver depending on kind.
func Testify(kind string, cnt string) (map[string]model.Test, error) {
	switch kind {

	case "surefire":
		s, err := surefire.New(cnt)
		if err != nil {
			return map[string]model.Test{}, err
		}
		return s.Testify()
	}
	return map[string]model.Test{}, errors.New("No driver to convert '" + kind + "' reports")
}
