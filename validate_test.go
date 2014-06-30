package validate

import (
	"testing"
	"time"
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
	var valid = []interface{}{
		"test",
	}
	doTest := func (slice []interface{}, expected bool, t *testing.T) {
		for _, v := range slice {
			validator := &Validator{Data: v}
			if validator.NotEmpty() != expected {
				t.Errorf("Test failed for %v", v)
			}
		}
	}
	doTest(invalid, false, t)
	doTest(valid, true, t)
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
	doTest := func (slice []interface{}, expected bool, t *testing.T) {
		for _, v := range slice {
			validator := &Validator{Data: v}
			if validator.NotZero() != expected {
				t.Errorf("Test failed for %v", v)
			}
		}
	}
	doTest(invalid, false, t)
	doTest(valid, true, t)
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
	doTest := func (slice []interface{}, expected bool, t *testing.T) {
		for _, v := range slice {
			validator := &Validator{Data: v}
			if validator.Email() != expected {
				t.Errorf("Test failed for %v", v)
			}
		}
	}
	doTest(invalid, false, t)
	doTest(valid, true, t)
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
	var length = 2
	doTest := func (slice []interface{}, expected bool, t *testing.T) {
		for _, v := range slice {
			validator := &Validator{Data: v}
			if validator.MinLength(length) != expected {
				t.Errorf("Test failed for %v", v)
			}
		}
	}
	doTest(invalid, false, t)
	doTest(valid, true, t)
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
	var length = 2
	doTest := func (slice []interface{}, expected bool, t *testing.T) {
		for _, v := range slice {
			validator := &Validator{Data: v}
			if validator.MaxLength(length) != expected {
				t.Errorf("Test failed for %v", v)
			}
		}
	}
	doTest(invalid, false, t)
	doTest(valid, true, t)
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
	var length = 4
	doTest := func (slice []interface{}, expected bool, t *testing.T) {
		for _, v := range slice {
			validator := &Validator{Data: v}
			if validator.Length(length) != expected {
				t.Errorf("Test failed for %v", v)
			}
		}
	}
	doTest(invalid, false, t)
	doTest(valid, true, t)
}

func TestGreaterThan(t *testing.T) {
	var invalid = []interface{}{
		"a",
		0,
		1.5,
		int16(2),
		float64(1.25),
	}
	var valid = []interface{}{
		int64(100),
		float32(192.123),
		12311,
		123.6,
	}
	var greaterThan = 50 // Ensure all valid items are greater than this number.
	doTest := func (slice []interface{}, expected bool, t *testing.T) {
		for _, v := range slice {
			validator := &Validator{Data: v}
			if validator.GreaterThan(greaterThan) != expected {
				t.Errorf("Test failed for %v", v)
			}
		}
	}
	doTest(invalid, false, t)
	doTest(valid, true, t)
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
	doTest := func (slice []interface{}, expected bool, t *testing.T) {
		for _, v := range slice {
			validator := &Validator{Data: v}
			if validator.UUID() != expected {
				t.Errorf("Test failed for %v", v)
			}
		}
	}
	doTest(invalid, false, t)
	doTest(valid, true, t)
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
		"foobar",
		"E55A815A-BA16-4FB9-AE01-644204CC433A", // Uppercase V4 - invalid hex digits
		time.Time{},
	}

	var valid = []interface{}{
		time.Date(1984, 1, 1, 12, 00, 00, 00, time.UTC),
		time.Now(),
	}
	doTest := func (slice []interface{}, expected bool, t *testing.T) {
		for _, v := range slice {
			validator := &Validator{Data: v}
			if validator.NotZeroTime() != expected {
				t.Errorf("Test failed for %v", v)
			}
		}
	}
	doTest(invalid, false, t)
	doTest(valid, true, t)
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
	doTest := func (slice []interface{}, expected bool, t *testing.T) {
		for _, v := range slice {
			validator := &Validator{Data: v}
			if validator.URL() != expected {
				t.Errorf("Test failed for %v", v)
			}
		}
	}
	doTest(invalid, false, t)
	doTest(valid, true, t)
}
