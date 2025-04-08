package models

import "errors"

var (
	PhoneInUse = errors.New("models: phone is already in use")
)
