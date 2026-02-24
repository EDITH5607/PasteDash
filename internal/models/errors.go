package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching Record Found !!")
	ErrInvalidCredential = errors.New("models: Invalid Credentials !!")
	ErrDuplicateEmail = errors.New("models: Duplicate Email !!!")

)
