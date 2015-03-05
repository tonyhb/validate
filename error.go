package validate

type ValidationError struct {
	Failures []string

	// Stores a list of fields that failed validation. This is useful
	// during testing: you can assert that all validation rules are
	// working as expected.
	Fields []string
}

func (ve *ValidationError) addFailure(field, msg string) {
	ve.Failures = append(ve.Failures, msg)
	ve.Fields = append(ve.Fields, field)
}

// Turn the slice of strings into one string.
func (ve ValidationError) Error() string {
	var str = "The following errors occured during validation: "
	for _, e := range ve.Failures {
		str += e + ". "
	}
	return str
}

func (ve ValidationError) Stringify() string {
	return ve.Error()
}
