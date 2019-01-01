package rules

import "fmt"

// This maps all validation tags to the corresponding validation methods
var rules map[string]ValidatorFunc

func init() {
	rules = map[string]ValidatorFunc{}
}

type ValidationData struct {
	// The name of the field being validated
	Field string

	// The value of the struct field being validated
	Value interface{}

	// Arguments from the validation tags. For example, in the following
	// definition Args will will contain a single "5":
	//
	// struct {
	//     Age `validate:"GreaterThan:5"`
	// }
	//
	// Unfortunately, due to the nature of tags these will always be strings.
	Args []string

	// Code to return in case of an error. Used to match the validation error
	// a specific message code.
	Code string
}

// All validation methods must return an ErrInvalid error type if the data
// is invalid, or nil if the data is valid
type ValidatorFunc func(ValidationData) error

// Add a new validation method for a given struct tag. If a validation method
// already exists this will return an error
func Add(tag string, method ValidatorFunc) (err error) {
	if _, ok := rules[tag]; ok {
		return fmt.Errorf("Validation method for '%s' already exists", tag)
	}

	rules[tag] = method
	return
}

// Return a registered validation method for a given tag
func Get(tag string) (method ValidatorFunc, err error) {
	if m, ok := rules[tag]; !ok {
		return nil, ErrNoValidationMethod{Tag: tag}
	} else {
		return m, nil
	}
}
