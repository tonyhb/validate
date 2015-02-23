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
		if ok, _ := NotEmpty(object); ok {
			t.Errorf("Expected NotEmpty to return false")
		}
	}

	object.Value = "valid"
	if ok, _ := NotEmpty(object); !ok {
		t.Errorf("Expected NotEmpty to return true")
	}

}
