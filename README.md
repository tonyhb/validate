Go Checker
=========

Validation for Go structs using Tags. Use as follows:

	type Page struct {
		UUID   string `validate:"NotEmpty,UUID"`
		URL    string `validate:"NotEmpty,URL"`
		Author string `validate:"Email"`
		Slug   string `validate:"minLength:5"
	}

Really simple definitions. To validate, use the exported methods:

	if pass, err := validate.run(page); pass != true {
		// err is of type validate.ValidateErrors - a slice of strings
		// Print all errors as a single string
		fmt.Printf(err.Stringify())
	}

Validating a subset of fields:

	if pass, err := validate.run(page, "URL"); pass != true {
		// Only the URL was Slug and Author was validated
	}

	if pass, err := validate.run(page, "Slug", "Author"); pass != true {
		// Only the Slug and Author was validated
	}

## Built in validators

- `NotZeroTime` - passes if the field is a non-zero Time
- `NotEmpty` - passes if the field is a non-empty string
- `Email` - passes if the field is a string with a valid email address
- `URL` - passes if the field is a string with a scheme and host
- `UUID` - passes if the field is a string, []byte or []rune and is a valid UUID
- `Length:N` - passes if the field is a string with N characters
- `MinLength:N` - passes if the field is a string with at least N characters
- `MaxLength:N` - passes if the field is a string with at most N characters
- `GreaterThan:N` - passes if the field is an integer or float over N

## Adding custom validators

You can add custom validators to the validation library without
extending/modifying it:

	import "validate"

	func IsAwesome(v *validate.Validator) bool {
		typ := reflect.TypeOf(v.Data)
		if typ.Name() != "ExpectedType" {
			v.Valid = false
			v.Error = "Invalid type!"
			return false
		}

		data = v.Data.(ExpectedType)
		if data != "awesome" {
			v.Valid = false
			v.Error = "This is just not awesome enough"
		}

		v.Valid = true
		return true
	}

	func main() {
		validate.AddMethod("MaybeAwesome", IsAwesome)
		// Now the tag `validate:"MaybeAwesome"` is valid.
	}

For the time being you need to explicitly set `v.Valid` in the validation
method.

Extracted from https://keepupdated.co - originally built September 2013,
maintained since then.
