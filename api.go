package api

import (
	"fmt"
	"sort"
)

// ErrorCase represents an error in an API object. It contains both the
// attribute indicated as, approximately, a dot-separated path to the field
// and a description of the error.
type ErrorCase struct {
	Attribute string `json:"attribute"`
	Msg       string `json:"msg"`
}

// ValidationError contains any errors that were found while trying to validate
// an API object.
type ValidationError struct {
	Errors []ErrorCase `json:"errors"`
}

func (ve *ValidationError) Error() string {
	plural := "s"
	if len(ve.Errors) == 1 {
		plural = ""
	}
	msg := fmt.Sprintf("%d validation error%s", len(ve.Errors), plural)

	for _, c := range ve.Errors {
		msg += "; " + fmt.Sprintf("%s: %s", c.Attribute, c.Msg)
	}

	return msg
}

// AddNew appends a new ErrorCase to the set of errors seen by this ValidationError
func (ve *ValidationError) AddNew(c ErrorCase) {
	ve.Errors = append(ve.Errors, c)
}

// OrNil will return a pointer to the ValidationError if any errors have been
// collected. If no errors have been collected it will return nil.
func (ve *ValidationError) OrNil() *ValidationError {
	if len(ve.Errors) == 0 {
		return nil
	}

	return ve
}

// Merge adds the errors collected in o to the errors in this ValidationError.
func (ve *ValidationError) Merge(o *ValidationError) {
	if o == nil {
		return
	}

	for _, e := range o.Errors {
		ve.AddNew(e)
	}
}

// MergePrefixed takes the Errors found in in children and appends them to this
// ValidationError. In the process it attaches the under prefix to the Attribute
// of the error case. The original children error is not modified.
func (ve *ValidationError) MergePrefixed(children *ValidationError, under string) {
	if children == nil {
		return
	}

	c2 := &ValidationError{}
	for _, e := range children.Errors {
		delim := ""
		if e.Attribute != "" && under != "" {
			delim = "."
		}

		c2.Errors = append(
			c2.Errors,
			ErrorCase{fmt.Sprintf("%s%s%s", under, delim, e.Attribute), e.Msg},
		)
	}

	ve.Merge(c2)
}

type ValidationErrorsByAttribute struct {
	e *ValidationError
}

var _ sort.Interface = ValidationErrorsByAttribute{}

func (eca ValidationErrorsByAttribute) Len() int {
	if eca.e == nil {
		return 0
	}
	return len(eca.e.Errors)
}

func (eca ValidationErrorsByAttribute) Less(i, j int) bool {
	return eca.e.Errors[i].Attribute < eca.e.Errors[j].Attribute
}

func (eca ValidationErrorsByAttribute) Swap(i, j int) {
	eca.e.Errors[i], eca.e.Errors[j] = eca.e.Errors[j], eca.e.Errors[i]
}
