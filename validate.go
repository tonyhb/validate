package validate

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/tonyhb/govalidate/rules"
	_ "github.com/tonyhb/govalidate/rules/email"
	_ "github.com/tonyhb/govalidate/rules/greaterthan"
	_ "github.com/tonyhb/govalidate/rules/length"
	_ "github.com/tonyhb/govalidate/rules/lessthan"
	_ "github.com/tonyhb/govalidate/rules/maxlength"
	_ "github.com/tonyhb/govalidate/rules/minlength"
	_ "github.com/tonyhb/govalidate/rules/notempty"
	_ "github.com/tonyhb/govalidate/rules/notzero"
	_ "github.com/tonyhb/govalidate/rules/notzerotime"
	_ "github.com/tonyhb/govalidate/rules/url"
	_ "github.com/tonyhb/govalidate/rules/uuid"
)

// Takes a struct, loops through all fields and calls check on any fields that
// have a validate tag.
func Run(object interface{}, fieldsSlice ...string) error {
	pass := true // We'll override this if checking returns false
	err := ValidationError{}

	// If we have been passed a slice of fields to valiate - to check only a
	// subset of fields - change the slice into a map for O(1) lookups instead
	// of O(n).
	fields := map[string]struct{}{}
	for _, field := range fieldsSlice {
		fields[field] = struct{}{}
	}

	// Iterate through each field of the struct and validate
	value := reflect.ValueOf(object)
	typ := value.Type() // A Type's Field method returns StructFields
	for i := 0; i < value.NumField(); i++ {
		var validateTag string
		var validateError error

		if len(fields) > 0 {
			// We're only checking for a subset of fields; if this field isn't
			// included in the subset of fields to validate we can skip.
			if _, ok := fields[typ.Field(i).Name]; !ok {
				continue
			}
		}

		if validateTag = typ.Field(i).Tag.Get("validate"); validateTag == "" {
			continue
		}

		// Validate this particular field against the options in our tag
		if validateError = check(value.Field(i).Interface(), typ.Field(i).Name, validateTag); validateError == nil {
			continue
		}

		if _, ok := validateError.(rules.ErrNoValidationMethod); ok {
			return validateError
		}

		pass = false
		err.addFailure(validateError.Error())
	}

	if pass {
		return nil
	}

	return err
}

// Takes a field's value and the validation tag and applies each check
// until either a test fails or all tests pass.
func check(data interface{}, fieldName, tag string) (err error) {
	for tag != "" {
		var next string
		var args []interface{}
		i := strings.Index(tag, ",")
		if i >= 0 {
			tag, next = tag[:i], tag[i+1:]
		}

		// tag is now the method we want to call. See if it has angle brackets,
		// which indicate arguments. This is only used for numerical arguments
		// in the validation methods Length, MinLength and MaxLength.
		// This means we need to typecast to an int from a atring here.
		i = strings.Index(tag, ":")
		if i > 0 {
			var argString string
			tag, argString = tag[:i], tag[i+1:]
			a, err := strconv.Atoi(argString)
			if err != nil {
				return fmt.Errorf("Couldn't convert the argument " + argString + " to an integer")
			}
			args = append(args, a)
		}

		// Attempt to validate the data using methods registered with the rules
		// sub package
		if method, err := rules.Get(tag); err != nil {
			return err
		} else {
			var data = rules.ValidationData{
				Field: fieldName,
				Value: data,
				Args:  args,
			}
			return method(data)
		}

		// Continue with the next tag
		tag = next
	}

	return nil
}
