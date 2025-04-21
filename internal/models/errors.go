package models

import "errors"

var (
	PhoneInUse      = errors.New("models: phone is already in use")
	EmailConstraint = errors.New("ent: constraint failed: pq: duplicate key value violates unique constraint \\\"member_email\\")
)
