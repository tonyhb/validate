Govalidate
=========

Simple, fast and *extensible* validation for Go structs, using tags in all their
goodness.

## Basic usage

Here's how to set up your struct:

```go
type Page struct {
	UUID   string `validate:"NotEmpty,UUID"`
	URL    string `validate:"NotEmpty,URL"`
	Author string `validate:"Email"`
	Slug   string `validate:"Regexp:/^[\w-]+$/, MinLength:5, MaxLength:100"`
}
```
Really simple definitions. To validate, use the exported methods:

```go
if err := validate.Run(page); err != nil {
	// err is of type validate.ValidateErrors which contains a slice of
	// validation errors for all failures.
	fmt.Printf(err.Error())
}
```

Validating a subset of fields:

```go
// Only validate the URL
if err := validate.Run(page, "URL"); err != nil {
	// Invalid data
}

// Only validate the Slug and Author fields
if err := validate.Run(page, "Slug", "Author"); err != nil {
	// Invalid data
}
```

## Built in validators

All validatiors are available in their own package within `rules`. These are
built in:

- `Regexp:/{regexp}/` - passes if a string matches the given regexp
- `Alpha` - passes if a string contains only alphabetic characters
- `Alphanumeric` - passes if a string contains only alphanumeric characters
- `Email` - passes if the field is a string with a valid email address
- `Length:N` - passes if the field is a string with N characters
- `MaxLength:N` - passes if the field is a string with at most N characters
- `MinLength:N` - passes if the field is a string with at least N characters
- `NotEmpty` - passes if the field is a non-empty string
- `NotZeroTime` - passes if the field is a non-zero Time
- `URL` - passes if the field is a string with a scheme and host
- `UUID` - passes if the field is a string, []byte or []rune and is a valid UUID
- `NotZero` - passes if the field is numeric and not-zero
- `GreaterThan:N` - passes if the field is numeric and over N
- `LessThan:N` - passes if the field is numeric and less than N

## Adding custom validators

Validators are built using interfaces. Even the built in ones. And adding a new
one is easy peasy:

```go
package yourvalidator

import (
	"github.com/tonyhb/govalidate/helper"
	"github.com/tonyhb/govalidate/rules"
)

func init() {
	// Register your validation tag with the validation method
	rules.Add("TagName", ValidationMethod)
}

// This accepts a ValidationData struct, which contains the field name, value
// and any arguments in the struct tag (such as '5' within MinLength:5)
func ValidationMethod(data rules.ValidationData) (err error) {
	// You'll need to typecast your data here
	v, ok := helper.ToString(data.Value)
	if ok != nil {
		return rules.ErrInvalid{
			ValidationData: data,
			Failure:        "is not a string",
		}
	}

	// Add custom validation logic, returning an error if the field is invalid.
	// rules.ErrInvalid has built in logic to make errors nicely formatted. It's
	// optional.
	if v == "" {
		return rules.ErrInvalid{
			ValidationData: data,
			Failure:        "is empty",
		}
	}

	// Congratulate your user for not fucking with you.
	return nil
}
```

Extracted from https://keepupdated.co - originally built September 2013,
maintained since then.
