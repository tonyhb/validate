package lessthan

import (
	"fmt"
	"strconv"

	"github.com/tonyhb/govalidate/helper"
	"github.com/tonyhb/govalidate/rules"
)

func init() {
	rules.Add("LessThan", LessThan)
}

// Passes if the data is a float/int and is less than the specified integer.
// Note that this is *not* a less than or equals check, and this only comapres
// floats/ints to a predefined integer specified in your tag.
// Fails if the data is not a float/int or the data is less than or equals the comparator
func LessThan(data rules.ValidationData) error {
	v, err := helper.ToFloat64(data.Value)
	if err != nil {
		return rules.ErrInvalid{
			ValidationData: data,
			Failure:        "is not numeric",
		}
	}

	// We should always be provided with a length to validate against
	if len(data.Args) == 0 {
		return fmt.Errorf("No argument found in the validation struct (eg 'LessThan:5')")
	}

	// Typecast our argument and test
	var max float64
	if max, err = strconv.ParseFloat(data.Args[0], 64); err != nil {
		return err
	}

	if v > max {
		return rules.ErrInvalid{
			ValidationData: data,
			Failure:        fmt.Sprintf("must be less than %d", max),
		}
	}

	return nil
}
