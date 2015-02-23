package minlength

import (
	"fmt"
	"strconv"

	"github.com/tonyhb/govalidate/helper"
	"github.com/tonyhb/govalidate/rules"
)

func init() {
	rules.Add("MinLength", MinLength)
}

// Used to check whether a string has at least N characters
// Fails if data is a string and its length is less than the specified comparator. Passes in all other cases.
func MinLength(data rules.ValidationData) error {
	v, err := helper.ToString(data.Value)
	if err != nil {
		return rules.ErrInvalid{
			ValidationData: data,
			Failure:        "is not a string",
		}
	}

	// We should always be provided with a length to validate against
	if len(data.Args) == 0 {
		return fmt.Errorf("No argument found in the validation struct (eg 'MinLength:5')")
	}

	// Typecast our argument and test
	var min int
	if min, err = strconv.Atoi(data.Args[0]); err != nil {
		return err
	}

	if len(v) < min {
		return rules.ErrInvalid{
			ValidationData: data,
			Failure:        fmt.Sprintf("is too short; it must be at least %d characters long", min),
		}

	}

	return nil
}
