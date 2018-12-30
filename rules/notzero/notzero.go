package notzero

import (
	"github.com/vizualni/govalidate/helper"
	"github.com/vizualni/govalidate/rules"
)

func init() {
	rules.Add("NotZero", NotZero)
}

// Checks whether a float or int type is 0. This could mean the data is above *or* below 0.
// Fails if the data isn't a float/int type, or the data is exactly 0.
func NotZero(data rules.ValidationData) error {
	v, err := helper.ToFloat64(data.Value)
	if err != nil {
		return rules.ErrInvalid{
			ValidationData: data,
			Failure:        "is not numeric",
		}
	}

	if v == 0 {
		return rules.ErrInvalid{
			ValidationData: data,
			Failure:        "is 0",
		}
	}

	return nil
}
