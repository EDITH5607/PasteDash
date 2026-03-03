package validator

import (
	"net/mail"
	"strings"
	"unicode/utf8"
)

// map for storing validation error
type Validator struct {
	// non field errors means error other than form type
	NonfieldErrors []string
	FieldErrors map[string]string
}

// for checking any validation error in the map
func (v *Validator) Valid() bool {
	return  len(v.FieldErrors) == 0 && len(v.NonfieldErrors) == 0
}

// for adding new validation error messages with correponding field
func (v *Validator) AddFieldError(key, message string) {

	// initialize the map if it is not initialzied yet
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	//checking the validation error message exist
	if _, exist := v.FieldErrors[key]; !exist {
		v.FieldErrors[key] = message
	}
}

// for adding new Nonfielderrors  like login failed to the slice in the struct	
func (v *Validator) AddNonFieldErrors(message string) {
	v.NonfieldErrors = append(v.NonfieldErrors, message)
}

// checking the condition and adding the message in the validation error map.
func (v *Validator)CheckField(ok bool, key, message string)  {
	if !ok {
		v.AddFieldError(key,message)
	}
}


// checking the maxlength of the string permitted or not.
func MaxChar(value string, n int) bool {
	// utf8.RuneCountInString used intead of len fun because of different characters are not accounted eg:ë
	return utf8.RuneCountInString(value) <= n
}

// checking is the string is blank or not
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// checking the value inside the list , we use go generic for this purpose
func PermitValues[T comparable](value T, permittedvalues ...T) bool {
	for i := range permittedvalues {
		if value == permittedvalues[i] {
			return true
		}
	}
	return false
}

// checking the string's is minchar length
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// mail.parseaddress is a standard library for validating string is email or not
func Matches(value string) bool {
	_, err := mail.ParseAddress(value)
	if err != nil {
		return false
	}
	return  true
}	



