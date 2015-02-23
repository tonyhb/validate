package validate

import (
	"testing"
	"time"

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

	object := struct {
		Data interface{} `validate:"NotEmpty"`
	}{}

	for _, v := range invalid {
		object.Data = v
		err := Run(object)
		if err == nil {
			t.Errorf("Expected invalid NotEmpty values to fail validation")
		}
		if _, ok := err.(rules.ErrNoValidationMethod); ok {
			t.Errorf(err.Error())
		}
	}

	object.Data = "valid"
	if err := Run(object); err != nil {
		t.Errorf("Unexpected error with valid field: %s", err.Error())
	}
}

func TestNotZero(t *testing.T) {
	var invalid = []interface{}{
		"a",
		struct{}{},
		0,
		int8(0),
		float32(0),
	}
	var valid = []interface{}{
		float64(2),
		int16(1231),
		1241631,
		float32(0.1),
		1,
		0.5,
	}

	object := struct {
		Data interface{} `validate:"NotZero"`
	}{}

	for _, v := range invalid {
		object.Data = v
		err := Run(object)
		if err == nil {
			t.Errorf("Expected invalid NotZero values to fail validation")
		}
		if _, ok := err.(rules.ErrNoValidationMethod); ok {
			t.Errorf(err.Error())
		}
	}

	for _, v := range valid {
		object.Data = v
		err := Run(object)
		if err != nil {
			t.Errorf("Unexpected error with valid NotZero value")
		}
	}
}

func TestEmail(t *testing.T) {
	var invalid = []interface{}{
		1,
		1.5,
		int8(1),
		float64(2.333),
		struct{}{},
		[]string{"test"},
		'a',
		"test@example",
		"test@example.",
		"testexample.",
		"example.com",
	}
	var valid = []interface{}{
		"test@example.com",
		"test@example.org.uk",
		"test@example.ru",
	}

	object := struct {
		Data interface{} `validate:"Email"`
	}{}

	for _, v := range invalid {
		object.Data = v
		err := Run(object)
		if err == nil {
			t.Errorf("Expected invalid Email values to fail validation")
		}
		if _, ok := err.(rules.ErrNoValidationMethod); ok {
			t.Errorf(err.Error())
		}
	}

	for _, v := range valid {
		object.Data = v
		err := Run(object)
		if err != nil {
			t.Errorf("Unexpected error with valid Email value")
		}
	}
}

func TestMinLength(t *testing.T) {
	var invalid = []interface{}{
		1,
		1.5,
		int8(1),
		float64(2.333),
		struct{}{},
		[]string{"test"},
		'a',
		"t",
	}
	var valid = []interface{}{
		"aa",
		"test",
	}

	object := struct {
		Data interface{} `validate:"MinLength:2"`
	}{}

	for _, v := range invalid {
		object.Data = v
		err := Run(object)
		if err == nil {
			t.Errorf("Expected invalid MinLength:2 values to fail validation")
		}
		if _, ok := err.(rules.ErrNoValidationMethod); ok {
			t.Errorf(err.Error())
		}
	}

	for _, v := range valid {
		object.Data = v
		err := Run(object)
		if err != nil {
			t.Errorf("Unexpected error with valid MinLength:2 value")
		}
	}
}

func TestMaxLength(t *testing.T) {
	var invalid = []interface{}{
		1,
		1.5,
		int8(1),
		float64(2.333),
		struct{}{},
		[]string{"test"},
		'a',
		"test",
	}
	var valid = []interface{}{
		"t",
	}

	object := struct {
		Data interface{} `validate:"MaxLength:2"`
	}{}

	for _, v := range invalid {
		object.Data = v
		err := Run(object)
		if err == nil {
			t.Errorf("Expected invalid MaxLength:2 values to fail validation")
		}
		if _, ok := err.(rules.ErrNoValidationMethod); ok {
			t.Errorf(err.Error())
		}
	}

	for _, v := range valid {
		object.Data = v
		err := Run(object)
		if err != nil {
			t.Errorf("Unexpected error with valid MaxLength:2 value")
		}
	}
}

func TestLength(t *testing.T) {
	var invalid = []interface{}{
		1,
		1.5,
		int8(1),
		float64(2.333),
		struct{}{},
		[]string{"test"},
		'a',
		"t",
		"foobar",
	}
	var valid = []interface{}{
		"test",
	}

	object := struct {
		Data interface{} `validate:"Length:4"`
	}{}

	for _, v := range invalid {
		object.Data = v
		err := Run(object)
		if err == nil {
			t.Errorf("Expected invalid Length:4 values to fail validation")
		}
		if _, ok := err.(rules.ErrNoValidationMethod); ok {
			t.Errorf(err.Error())
		}
	}

	for _, v := range valid {
		object.Data = v
		err := Run(object)
		if err != nil {
			t.Errorf("Unexpected error with valid Length:4 value")
		}
	}
}

func TestGreaterThan(t *testing.T) {
	var invalid = []interface{}{
		"a",
		0,
		1.5,
		int16(2),
		float64(1.25),
		49.99,
	}
	var valid = []interface{}{
		int64(100),
		float32(192.123),
		12311,
		123.6,
		50,
	}

	object := struct {
		Data interface{} `validate:"GreaterThan:50"`
	}{}

	for _, v := range invalid {
		object.Data = v
		err := Run(object)
		if err == nil {
			t.Errorf("Expected invalid GreaterThan:50 values to fail validation")
		}
		if _, ok := err.(rules.ErrNoValidationMethod); ok {
			t.Errorf(err.Error())
		}
	}

	for _, v := range valid {
		object.Data = v
		err := Run(object)
		if err != nil {
			t.Errorf("Unexpected error with valid GreaterThan:50 value")
		}
	}
}

func TestValidateUUID(t *testing.T) {
	var invalid = []interface{}{
		1,
		1.5,
		int8(1),
		float64(2.333),
		struct{}{},
		[]string{"test"},
		[]byte("test"),
		[]rune("test"),
		'a',
		"t",
		"foobar",
		"E55A815A-BA16-4FB9-AE01-644204CC433A", // Uppercase V4 - invalid hex digits
	}
	var valid = []interface{}{
		"fb623672-40dd-11e3-91ea-ce3f5508acd9", // V1
		"8563d95d-efb0-4e87-95d8-1d6c5debf298", // V4
	}

	object := struct {
		Data interface{} `validate:"UUID"`
	}{}

	for _, v := range invalid {
		object.Data = v
		err := Run(object)
		if err == nil {
			t.Errorf("Expected invalid UUID values to fail validation")
		}
		if _, ok := err.(rules.ErrNoValidationMethod); ok {
			t.Errorf(err.Error())
		}
	}

	for _, v := range valid {
		object.Data = v
		err := Run(object)
		if err != nil {
			t.Errorf("Unexpected error with valid UUID value")
		}
	}

}

func TestValidateNotZeroTime(t *testing.T) {
	var invalid = []interface{}{
		1,
		1.5,
		int8(1),
		float64(2.333),
		struct{}{},
		[]string{"test"},
		[]byte("test"),
		[]rune("test"),
		'a',
		"t",
		time.Time{},
	}
	var valid = []interface{}{
		time.Date(1984, 1, 1, 12, 00, 00, 00, time.UTC),
		time.Now(),
	}

	object := struct {
		Data interface{} `validate:"NotZeroTime"`
	}{}

	for _, v := range invalid {
		object.Data = v
		err := Run(object)
		if err == nil {
			t.Errorf("Expected invalid NotZeroTime values to fail validation")
		}
		if _, ok := err.(rules.ErrNoValidationMethod); ok {
			t.Errorf(err.Error())
		}
	}

	for _, v := range valid {
		object.Data = v
		err := Run(object)
		if err != nil {
			t.Errorf("Unexpected error with valid NotZeroTime value")
		}
	}

}

func TestValidateURL(t *testing.T) {
	var invalid = []interface{}{
		"test",
		"test",
		"http://",
		"example.c\\",
		"example.com",
		"http//example.com/",
		"http::/example.com/",
		"http://example\\.com",
	}

	var valid = []interface{}{
		"http://example.com",
		"http://example.com/",
		"HTTP://example.com/",
		"https://www.example.com/index.html",
	}

	object := struct {
		Data interface{} `validate:"URL"`
	}{}

	for _, v := range invalid {
		object.Data = v
		err := Run(object)
		if err == nil {
			t.Errorf("Expected invalid URL values to fail validation")
		}
		if _, ok := err.(rules.ErrNoValidationMethod); ok {
			t.Errorf(err.Error())
		}
	}

	for _, v := range valid {
		object.Data = v
		err := Run(object)
		if err != nil {
			t.Errorf("Unexpected error with valid URL value")
		}
	}

}

func TestValidateRegexp(t *testing.T) {
	var invalid = []interface{}{
		1,
		'a',
		"0aaa0",
	}
	var valid = []interface{}{
		"aaaaa0",
		"aaa123456789",
	}

	object := struct {
		Data interface{} `validate:"MinLength:1, Regexp:/^[a-zA-Z]{3,5}[0-9]+$/, NotEmpty"`
	}{}

	for _, v := range invalid {
		object.Data = v
		err := Run(object)
		if err == nil {
			t.Errorf("Expected invalid regexp to fail validation")
		}
		if _, ok := err.(rules.ErrNoValidationMethod); ok {
			t.Errorf(err.Error())
		}
	}

	for _, v := range valid {
		object.Data = v
		err := Run(object)
		if err != nil {
			t.Errorf("Unexpected error with valid regexp value: %s", err.Error())
		}
	}

}
