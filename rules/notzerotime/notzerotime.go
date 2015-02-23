package notzerotime

import (
	"time"

	"github.com/tonyhb/govalidate/rules"
)

func init() {
	rules.Add("NotZeroTime", NotZeroTime)
}

// Checks whether a float or int type is 0. This could mean the data is above *or* below 0.
// Fails if the data isn't a float/int type, or the data is exactly 0.
func NotZeroTime(data rules.ValidationData) error {
	if _, ok := data.Value.(time.Time); !ok {
		return rules.ErrInvalid{
			ValidationData: data,
			Failure:        "is not a Time type",
		}
	}

	if data.Value.(time.Time).Equal(time.Time{}) == true {
		return rules.ErrInvalid{
			ValidationData: data,
			Failure:        "has a zero value",
		}
	}

	return nil
}
