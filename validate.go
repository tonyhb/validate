package validate

import (
	"errors"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type MethodMap map[string]func(v *Validator) bool

var methods = MethodMap{}

func AddMethod(name string, f func(v *Validator) bool) {
	methods[name] = f
}

// Takes a struct, loops through all fields and calls check on any fields that
// have a validate tag.
func Run(object interface{}, fieldsSlice ...string) (pass bool, errors ValidateErrors) {
	pass = true // We'll override this if checking returns false

	// We're going to change the slice of fields into a map for O(1) lookups
	// instead of O(n). This is used when looking for a subset of fields.
	fields := map[string]struct{}{}
	for _, field := range fieldsSlice {
		fields[field] = struct{}{}
	}

	value := reflect.ValueOf(object)
	typ := value.Type() // A Type's Field method returns StructFields
	for i := 0; i < value.NumField(); i++ {
		if len(fields) > 0 {
			// We're only checking for a subset of fields
			if _, ok := fields[typ.Field(i).Name]; !ok {
				continue
			}
		}
		var validateTag string
		validateTag = typ.Field(i).Tag.Get("validate")
		if validateTag == "" {
			continue
		}
		// Validate this particular field against the options in our tag
		if ok, err := check(value.Field(i).Interface(), typ.Field(i).Name, validateTag); !ok {
			errors = append(errors, err)
			pass = false
		}
	}
	return
}

// Takes a field's value and the validation tag and applies each check
// until either a test fails or all tests pass.
func check(data interface{}, fieldName, tag string) (valid bool, err string) {
	v := &Validator{Data: data, Field: fieldName}
	value := reflect.ValueOf(v)
	for tag != "" {
		var next string
		var args []reflect.Value
		i := strings.Index(tag, ",")
		if i >= 0 {
			tag, next = tag[:i], tag[i+1:]
		}

		// tag is now the method we want to call. See if it has angle brackets,
		// which indicate arguments. This is only used for numerical arguments
		// in the validation methods Length, MinLength and MaxLength.
		// This means we need to typecast to an int from a atring here.
		i = strings.Index(tag, ":")
		if i > 0 {
			var argString string
			tag, argString = tag[:i], tag[i+1:]
			a, err := strconv.Atoi(argString)
			if err != nil {
				return false, "Couldn't convert the argument " + argString + " to an integer"
			}
			args = append(args, reflect.ValueOf(a))
		}

		// Check to see if this validation function is valid; if so, call it,
		// which modifies the Validator struct v. This is good because we
		// don't have to play around with reflect value returns from the
		// method call. Ugly.
		if method := value.MethodByName(tag); method.IsValid() {
			method.Call(args)
		} else {
			// Check the map to see if this method has been added by the calling module
			// @TODO: This is very hacky. Plz tidy =)
			run := false
			for key, method := range methods {
				// We've got the methods in an overridden method block: run it
				if key == tag {
					run = true
					method(v)
				}
			}

			// We haven't got the method: fail
			if !run {
				return false, "The validation method " + tag + " does not exist. Failing validation."
			}
		}

		if !v.Valid {
			return false, v.Error
		}

		tag = next
	}
	return true, ""
}

// This structs contains the data we're validating and runs the actual
// validation methods. To add custom rules, use a pointer to Validator as the
// receiver of the method.
type Validator struct {
	Field string
	Data  interface{}
	Valid bool
	Error string
}

// Checks whether a string is empty.
// Passes if the data is a non-empty string. Fails if the data isn't a string, or the data is an empty string.
func (this *Validator) NotEmpty() bool {
	data, err := toString(this.Data)
	if err != nil {
		this.Valid = false
		this.Error = this.Field + " is not a string"
		return false
	}

	if data == "" {
		this.Valid = false
		this.Error = this.Field + " is empty"
	} else {
		this.Valid = true
	}
	return this.Valid
}

// Checks whether a float or int type isn't 0. This could mean the data is above *or* below 0.
// Fails if the data isn't a float/int type, or the data is exactly 0.
func (this *Validator) NotZero() bool {
	data, err := toFloat64(this.Data)
	if err != nil {
		this.Valid = false
		this.Error = this.Field + " has an invalid type in testing NotZero"
		return this.Valid
	}
	if data == 0 {
		this.Valid = false
		this.Error = this.Field + " is 0"
	} else {
		this.Valid = true
	}
	return this.Valid
}

// Passes if the  string is an email address. Fails otherwise.
func (this *Validator) Email() bool {
	data, err := toString(this.Data)
	if err != nil {
		this.Valid = false
		this.Error = this.Field + " is not a string"
		return false
	}

	if regexp.MustCompile(`(?i)[A-Z0-9!#$%&'*+/=?^_{|}~-]+(?:\.[A-Z0-9!#$%&'*+/=?^_{|}~-]+)*@(?:[A-Z0-9](?:[A-Z0-9-]*[A-Z0-9])?\.)+[A-Z0-9](?:[A-Z0-9-]*[A-Z0-9])?`).MatchString(data) {
		this.Valid = true
	} else {
		this.Valid = false
		this.Error = this.Field + " is not a valid email address"
	}
	return this.Valid
}

// Used to check whether a string has at least N characters
// Fails if data is a string and its length is less than the specified comparator. Passes in all other cases.
func (this *Validator) MinLength(minLength int) bool {
	data, err := toString(this.Data)
	if err != nil {
		this.Valid = false
		this.Error = this.Field + " is not a string"
		return false
	}

	if len(data) < minLength {
		this.Valid = false
		this.Error = this.Field + " is too short. It must be at least " + strconv.Itoa(minLength) + " characters long."
	} else {
		this.Valid = true
	}
	return this.Valid
}

// Used to check whether a string has at most N characters
// Fails if data is a string and its length is more than the specified comparator. Passes in all other cases.
func (this *Validator) MaxLength(maxLength int) bool {
	data, err := toString(this.Data)
	if err != nil {
		this.Valid = false
		this.Error = this.Field + " is not a string"
		return false
	}
	if len(data) > maxLength {
		this.Valid = false
		this.Error = this.Field + " is too long"
	} else {
		this.Valid = true
	}
	return this.Valid
}

// Used to check whether a string has a specified number of characters
// Fails if data is a string and its length is not equal to a specified comparator. Passes in all other cases.
func (this *Validator) Length(length int) bool {
	data, err := toString(this.Data)
	if err != nil {
		this.Valid = false
		this.Error = this.Field + " is not a string"
		return false
	}
	if len(data) != length {
		this.Valid = false
		this.Error = this.Field + " doesn't meet the expected length"
	} else {
		this.Valid = true
	}
	return this.Valid
}

// Passes if the data is a float/int and is greater than the specified integer.
// Note that this is *not* a greater than or equals check, and this only comapres
// floats/ints to a predefined integer specified in your tag.
// Fails if the data is not a float/int or the data is less than or equals the comparator
func (this *Validator) GreaterThan(minLength int) bool {
	data, err := toFloat64(this.Data)
	if err != nil {
		this.Valid = false
		this.Error = this.Field + " has an invalid type in testing GreaterThan"
		return this.Valid
	}
	if data <= float64(minLength) {
		this.Valid = false
		this.Error = this.Field + " is less than or equal to " + strconv.Itoa(minLength)
	} else {
		this.Valid = true
	}
	return this.Valid
}

// Used to check whether a string, []byte, or []rune is a valid UUID. Passes if the data is valid, fails
// if the data is invalid or not one of the accepted types.
func (this *Validator) UUID() bool {
	data, err := toString(this.Data)
	if err != nil {
		this.Valid = false
		this.Error = this.Field + " has an invalid type in UUID"
		return this.Valid
	}

	var hexPattern = "^(urn\\:uuid\\:)?\\{?([a-z0-9]{8})-([a-z0-9]{4})-([1-5][a-z0-9]{3})-([a-z0-9]{4})-([a-z0-9]{12})\\}?$"
	re := regexp.MustCompile(hexPattern)
	md := re.FindStringSubmatch(data)
	if md == nil {
		this.Valid = false
		this.Error = "Invalid UUID: " + data
	} else {
		this.Valid = true
	}
	return this.Valid
}

// Used to check whether a time has a zero value. Passes if the time has a non-zero value, fails if the time has a zero value
// or isn't a time.Time type.
func (this *Validator) NotZeroTime() bool {
	typ := reflect.TypeOf(this.Data)
	if typ.Name() != "Time" {
		this.Valid = false
		this.Error = this.Field + " is not a Time type"
		return false
	}
	if this.Data.(time.Time).Equal(time.Time{}) == true {
		this.Valid = false
		this.Error = this.Field + " has a zero value"
	} else {
		this.Valid = true
	}
	return this.Valid
}

// Validates a URL. For a valid URL, the following need to be present in a parsed URL:
// * Scheme (either http or https)
// * Host (without a backslash)
func (this *Validator) URL() bool {
	data, err := toString(this.Data)
	if err != nil {
		this.Valid = false
		this.Error = "Invalid type"
		return false
	}
	parsed, err := url.Parse(data)
	if err != nil {
		this.Valid = false
		this.Error = "Error parsing the URL field: " + this.Field
		return this.Valid
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		this.Valid = false
		this.Error = this.Field + " has an invalid scheme; the scheme must be 'http' or 'https'"
		return this.Valid
	}

	if parsed.Host == "" || strings.IndexRune(parsed.Host, '\\') > 0 {
		this.Valid = false
		this.Error = this.Field + " has an invalid host"
		return this.Valid
	}

	this.Valid = true
	return this.Valid
}

// Helper method, converting all int and float types in an interface to a float64.
func toFloat64(data interface{}) (float64, error) {
	switch data.(type) {
	case float64:
		return data.(float64), nil
	case float32:
		return float64(data.(float32)), nil
	case int64:
		return float64(data.(int64)), nil
	case int32:
		return float64(data.(int32)), nil
	case int16:
		return float64(data.(int16)), nil
	case int8:
		return float64(data.(int8)), nil
	case int:
		return float64(data.(int)), nil
	}
	return 0, errors.New("Invalid conversion to float64")
}

func toString(data interface{}) (string, error) {
	switch data.(type) {
	case string:
		return data.(string), nil
	case []byte:
		return string(data.([]byte)), nil
	case []rune:
		return string(data.([]rune)), nil
	}

	return "", errors.New("Invalid conversion to string")
}
