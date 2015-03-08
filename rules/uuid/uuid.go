package uuid

import (
	"regexp"

	"github.com/tonyhb/govalidate/helper"
	"github.com/tonyhb/govalidate/rules"
)

func init() {
	rules.Add("UUID", UUID)
}

// Used to check whether a string has at most N characters
// Fails if data is a string and its length is more than the specified comparator. Passes in all other cases.
func UUID(data rules.ValidationData) error {
	v, err := helper.ToString(data.Value)
	if err != nil {
		return rules.ErrInvalid{
			ValidationData: data,
			Failure:        "is not a string",
		}
	}

	if !IsUUID(v) {
		return rules.ErrInvalid{
			ValidationData: data,
			Failure:        "is an invalid UUID",
		}
	}

	return nil
}

func IsUUID(uuid string) bool {
	var hexPattern = "^(urn\\:uuid\\:)?\\{?([a-z0-9]{8})-([a-z0-9]{4})-([1-5][a-z0-9]{3})-([a-z0-9]{4})-([a-z0-9]{12})\\}?$"
	re := regexp.MustCompile(hexPattern)

	if match := re.FindStringSubmatch(uuid); match == nil {
		return false
	}
	return true
}
