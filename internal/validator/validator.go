package validator

import ( 
	"strings"
	"unicode/utf8"
)

// map for storing validation error
type Validator struct {
	FieldErrors map[string]string
}

// for checking any validation error in the map
func (v *Validator) Valid() bool {
	return  len(v.FieldErrors) == 0
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

// checking the value inside the list
func PermitInt(value int, permittedvalues ...int) bool {
	for i := range permittedvalues {
		if value == permittedvalues[i] {
			return true
		}
	}
	return false
}



