package validate

type ValidationError struct {
	Failures []string
}

func (ve *ValidationError) addFailure(msg string) {
	ve.Failures = append(ve.Failures, msg)
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
