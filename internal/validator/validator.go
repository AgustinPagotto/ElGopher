package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	Errors []string
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(message string) {
	v.Errors = append(v.Errors, message)
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddError(message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}
