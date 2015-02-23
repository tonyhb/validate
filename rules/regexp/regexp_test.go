package regexp

import (
	"testing"

	"github.com/tonyhb/govalidate/rules"
)

func TestRegexp(t *testing.T) {
	object := rules.ValidationData{
		Field: "Test",
		Args:  []string{"/^[a-zA-Z]{3,5}[0-9]+$/"},
	}

	var valid = []interface{}{
		"aaaaa0",
		"aaa123456789",
	}
	var invalid = []interface{}{
		1,
		'a',
		"0aaa0",
	}

	for _, v := range invalid {
		object.Value = v
		if err := Regexp(object); err == nil {
			t.Errorf("Expected error with invalid values")
		}
	}

	for _, v := range valid {
		object.Value = v
		if err := Regexp(object); err != nil {
			t.Errorf("Unexpected error with valid values")
		}
	}
}
