package regexp

import (
	"fmt"
	"regexp"

	"github.com/tonyhb/govalidate/helper"
	"github.com/tonyhb/govalidate/rules"
)

func init() {
	rules.Add("Regexp", Regexp)
}

// Validates that a string only contains alphabetic characters
func Regexp(data rules.ValidationData) (err error) {
	v, err := helper.ToString(data.Value)
	if err != nil {
		return rules.ErrInvalid{
			ValidationData: data,
			Failure:        "is not a string",
		}
	}

	// We should always be provided with a length to validate against
	if len(data.Args) == 0 {
		return fmt.Errorf("No argument found in the validation struct (eg 'Regexp:/^\\s+$/')")
	}

	// Remove the trailing slashes from our regex string. Regexps must be enclosed
	// within two "/" characters.
	re := data.Args[0]
	re = re[1 : len(re)-1]
	if regexp.MustCompile(re).MatchString(v) == false {
		return rules.ErrInvalid{
			ValidationData: data,
			Failure:        "doesn't match regular expression",
		}
	}

	return nil
}
