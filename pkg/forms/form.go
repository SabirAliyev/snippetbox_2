package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

// Custom form struct, which anonymously embeds a urlValues object (to hold the form data)
// and an Errors field to hold any validation errors for the form data.
type Form struct {
	url.Values
	Errors errors
}

// Initialize a custom Form struct. Notice that this takes the Form data as the parameter.
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Check that specific fields in the Form data are present
// and not blank. If any filed fails in this check add the appropriate message to the form errors.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Check that the specific field in the form contains a maximum
// number of characters, If the check fails then add the appropriate message to the for errors.
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maxidum is %d characters)", d))
	}
}

// Check that a specific field in the form matches one of set of specific permitted values.
// If the check fails then add the appropriate message to the form errors.
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
}

// Method which returns true if there are no errors.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
