package validate

type ValidationError struct {
	Failures []string

	// Stores a list of fields that failed validation. This is useful
	// during testing: you can assert that all validation rules are
	// working as expected.
	//
	// This is a map because maps can be compared via reflect.DeepEqual even if
	// its contents are ordered differently. Slices can not. This makes testing
	// more reliable: you can do this:
	//
	//   if ! reflect.DeepEqual(ExpectedFields, err.Fields) {
	//       t.Fatalf()
	//   }
	//
	// instead of ordering your ExpectedField slice in the same manner as
	// returned by ValidationError.
	//
	// @see http://play.golang.org/p/MhB4tDJVCz
	Fields map[string]struct{}
}

func (ve *ValidationError) addFailure(field, msg string) {
	ve.Failures = append(ve.Failures, msg)

	// Ensure we're not assigning to a nil map
	if ve.Fields == nil {
		ve.Fields = map[string]struct{}{}
	}
	ve.Fields[field] = struct{}{}
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

// Merge validation errors together. This is used with recursion when validating
// anonymous structs.
func (ve *ValidationError) Merge(other ValidationError) {
	for _, v := range other.Failures {
		ve.Failures = append(ve.Failures, v)
	}
	for f, v := range other.Fields {
		if ve.Fields == nil {
			ve.Fields = map[string]struct{}{}
		}
		ve.Fields[f] = v
	}
}
