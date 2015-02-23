package email

import (
	"regexp"

	"github.com/tonyhb/govalidate/helper"
	"github.com/tonyhb/govalidate/rules"
)

func init() {
	rules.Add("Email", Email)
}

func Email(data rules.ValidationData) (err error) {
	v, ok := helper.ToString(data.Value)
	if ok != nil {
		return rules.ErrInvalid{
			ValidationData: data,
			Failure:        "is not a string",
		}
	}

	if regexp.MustCompile(`(?i)[A-Z0-9!#$%&'*+/=?^_{|}~-]+(?:\.[A-Z0-9!#$%&'*+/=?^_{|}~-]+)*@(?:[A-Z0-9](?:[A-Z0-9-]*[A-Z0-9])?\.)+[A-Z0-9](?:[A-Z0-9-]*[A-Z0-9])?`).MatchString(v) {
		return
	}

	return rules.ErrInvalid{
		ValidationData: data,
		Failure:        "is not a valid email address",
	}
}
