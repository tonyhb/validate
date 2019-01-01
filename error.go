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

	// A map that stores field with the corresponding error messages
	// This is added as a transation period not to break anybody's
	// current implementation.
	FieldFailures map[string] map[string][]string
}

func (ve *ValidationError) addFailure(field, code, msg string) {
	ve.Failures = append(ve.Failures, msg)

	// Ensure we're not assigning to a nil map
	if ve.Fields == nil {
		ve.Fields = map[string]struct{}{}
	}
	ve.Fields[field] = struct{}{}

	if ve.FieldFailures == nil {
		ve.FieldFailures = make(map[string]map[string][]string)
	}

	if ve.FieldFailures[field] == nil {
		ve.FieldFailures[field] = make(map[string][]string)
	}

	ve.FieldFailures[field][code] = append(ve.FieldFailures[field][code], msg)
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

	if ve.FieldFailures == nil {
		ve.FieldFailures = make(map[string]map[string][]string)
	}

	for f := range other.FieldFailures {
		if ve.FieldFailures[f] == nil {
			ve.FieldFailures[f] = make(map[string][]string)
		}

		for c := range other.FieldFailures[f] {
			ve.FieldFailures[f][c] = append(ve.FieldFailures[f][c], other.FieldFailures[f][c]...)
		}
	}
}
