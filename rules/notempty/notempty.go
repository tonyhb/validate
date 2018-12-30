package notempty

import (
	"github.com/vizualni/govalidate/helper"
	"github.com/vizualni/govalidate/rules"
)

func init() {
	rules.Add("NotEmpty", NotEmpty)
}

// Checks whether a string is empty.
// Passes if the data is a non-empty string. Fails if the data isn't a string, or the data is an empty string.
func NotEmpty(data rules.ValidationData) (err error) {
	v, ok := helper.ToString(data.Value)
	if ok != nil {
		return rules.ErrInvalid{
			ValidationData: data,
			Failure:        "is not a string",
		}
	}
	if v == "" {
		return rules.ErrInvalid{
			ValidationData: data,
			Failure:        "is empty",
		}
	}
	return nil
}
