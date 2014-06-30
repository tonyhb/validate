package validate

type ValidateErrors []string

// Turn the slice of strings into one string.
func (this ValidateErrors) Stringify() string {
	var str = "The following errors occured during validation: "
	for _, e := range this {
		str += e + ". "
	}
	return str
}
