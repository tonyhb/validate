Govalidate
=========

[![Build Status](https://travis-ci.org/tonyhb/govalidate.svg?branch=master)](https://travis-ci.org/tonyhb/govalidate)
[![GoDoc](https://godoc.org/github.com/vizualni/govalidate?status.svg)](https://godoc.org/github.com/vizualni/govalidate)

Simple, fast and *extensible* validation for Go structs, using tags in all their
goodness. It also validates anonymous structs automatically.

```
GoCode   import github.com/vizualni/govalidate
CLI      go get -u github.com/vizualni/govalidate
```

## Basic usage

Here's how to set up your struct:

```go
package main

import "github.com/vizualni/govalidate"

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

Validating anonymous structs:

```go
package main

import "github.com/vizualni/govalidate"

type Author struct {
	Name  string `validate:"NotEmpty"`
	Email string `validate:"Email"`
}

type Content struct {
	Author
	Body    string `validate:"NotEmpty"`
}

p := Content{Body: "foobar"}

if err := validate.Run(p); err != nil {
	// The validation library will validate all Content fields plus the anonymous
 	// Author field embedded within it.
	// Because the Auther fields aren't set this will fail validation.
	fmt.Println(err.Error())
	fmt.Printf("%#v\n", err)

	// Output:
	// The following errors occured during validation: Field 'Name' is empty. Field 'Email' is not a valid email address.
	// validate.ValidationError{Failures:[]string{"Field 'Name' is empty", "Field 'Email' is not a valid email address"}, Fields:map[string]struct {}{"Email":struct {}{}, "Name":struct {}{}}}
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
	"github.com/vizualni/govalidate/helper"
	"github.com/vizualni/govalidate/rules"
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

## Testing

You can test that validation rules are working as expected:

```go
type User struct {
	Name  string `validate:"NotEmpty"`
	Email string `validate:"NotEmpty"`
}

func TestUserValidation(t *testing.T) {

	tests := []struct {
		User   *User
		Valid  bool
		Fields map[string]struct{}
	}{
		// In this test we expect both the Name and Email fields to fail validation:
		// we're passing an empty User object into validation
		{
			User:  &User{},
			Valid: false,
			// These fields should be present in ValidationError.Fields as failures
			// This is a map with meaningless empty values, rather than a slice,
			// so you can do O(1) lookups and compare them using reflect.DeepEqual
			// without caring about the order in which fields will appear
			Fields: map[string]struct{}{
				"Name":  struct{}{},
				"Email": struct{}{},
			},
		},
	}

	for _, v := range tests {
		err := validate.Run(v.User)

		if v.Valid && err != nil {
			t.Fatalf("Unexpected validation error: %s", err)
		}

		if !v.Valid && err == nil {
			t.Fatal("Expected validation error")
		}

		// Check that the fields were validated as we expect
		if !v.Valid && err != nil {
			// Ensure we were passed a ValidationError:
			// if a validation method wasn't present this error will
			// be of some other type.
			vErr, ok := err.(validate.ValidationError)
			if !ok {
				t.Fatalf(err.Error())
			}

			// Check the failed fields matches what we expect
			// These are mostly useful for tests
			// If you want readable validation messages, you can look at vErr.Failures
			if len(v.Fields) > 0 && !reflect.DeepEqual(v.Fields, vErr.Fields) {
				t.Fatal()
			}
		}
	}
}
```

MIT licence.

Extracted from https://keepupdated.co - originally built September 2013,
maintained since then.
