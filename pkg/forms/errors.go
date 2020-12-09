package forms

// Define errors type, which we use to hold the validation error message for forms.
// The name of the form field will be used as key in this map.
type errors map[string][]string

// Add messages for a given field to the map.
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Retrieve the first error message for a given filed from the map.
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
