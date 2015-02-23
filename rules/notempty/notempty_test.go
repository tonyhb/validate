package notempty

import (
	"testing"

	"github.com/tonyhb/govalidate/rules"
)

func TestNotEmpty(t *testing.T) {
	var invalid = []interface{}{
		1,
		1.5,
		int8(1),
		float64(2.333),
		struct{}{},
		[]string{"test"},
		'a',
		"",
	}

	object := rules.ValidationData{
		Field: "Test",
	}

	for _, v := range invalid {
		object.Value = v
		if err := NotEmpty(object); err == nil {
			t.Errorf("Expected error with invalid values")
		}
	}

	object.Value = "valid"
	if err := NotEmpty(object); err != nil {
		t.Errorf("Unexpected error with valid values")
	}

}
