package rules

import "fmt"

type ErrInvalid struct {
	ValidationData
	Failure string
}

func (t ErrInvalid) Error() string {
	return fmt.Sprintf("Field '%s' %s", t.Field, t.Failure)
}

type ErrNoValidationMethod struct {
	Tag string
}

func (t ErrNoValidationMethod) Error() string {
	return fmt.Sprintf("No validation method for '%s' has been registered", t.Tag)
}
