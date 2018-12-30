package email

import (
	"regexp"

	"github.com/vizualni/govalidate/helper"
	"github.com/vizualni/govalidate/rules"
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

	if IsEmail(v) {
		return
	}

	return rules.ErrInvalid{
		ValidationData: data,
		Failure:        "is not a valid email address",
	}
}

func IsEmail(str string) bool {
	return regexp.MustCompile(`(?i)[A-Z0-9!#$%&'*+/=?^_{|}~-]+(?:\.[A-Z0-9!#$%&'*+/=?^_{|}~-]+)*@(?:[A-Z0-9](?:[A-Z0-9-]*[A-Z0-9])?\.)+[A-Z0-9](?:[A-Z0-9-]*[A-Z0-9])?`).MatchString(str)
}
