package models

import "errors"

var (
	PhoneInUse                    = errors.New("models: phone is already in use")
	EmailConstraint               = errors.New("ent: constraint failed: pq: duplicate key value violates unique constraint \\\"member_email\\")
	EmailExistsError              = errors.New("models: email already exists")
	RegistrationTokenInvalidError = errors.New("models: registration token is invalid")
	UsernameExistsError           = errors.New("models: username already exists")
	BcryptError                   = errors.New(`models: bcrypt hashing failed`)
	TokenExpiredError             = errors.New("models: registration token is expired")
	TokenMissMatchError           = errors.New("models: registration token mismatch")
	TokenValidationError          = errors.New("models: registration token validation failed")
	WrongVerficationTokenError    = errors.New("models: wrong verification token")
	InvalidCredentialsError       = errors.New("models: invalid credentials")
	EmailUsernameInvalidError     = errors.New("models: email or username is invalid")
	AccountNotVerifiedError       = errors.New("models: account is not verified")
	UserNotFoundError             = errors.New("models: user not found")
	ConstraintError               = errors.New("models: constraint failed")
)
