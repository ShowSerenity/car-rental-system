package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"showserenity.net/car-rental-system/pkg/models"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-])")

type Form struct {
	url.Values
	Errors errors
	Car    *models.Car
	Cars   []*models.Car
}

func NewCar(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
		nil,
		nil,
	}
}

func NewSnippet(data url.Values, Cars []*models.Car) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
		nil,
		Cars,
	}
}

func NewSignUp(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
		nil,
		nil,
	}
}

func NewRent(data url.Values, Car *models.Car) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
		Car,
		nil,
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters)", d))
	}
}

func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d characters)", d))
	}
}

func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors.Add(field, "This field is invalid")
	}
}
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
