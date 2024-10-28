package models

import (
	"fmt"
	"net/url"
	"time"
)

const errorAuthority = "parkwithease.github.io"

var (
	CodeInvalidCredentials   = NewUserErrorCode("invalid-credentials", "2024-10-13")
	CodePasswordLength       = NewUserErrorCode("password-length", "2024-10-13")
	CodeInvalidEmail         = NewUserErrorCode("invalid-email", "2024-10-13")
	CodeDuplicate            = NewUserErrorCode("duplicate-entity", "2024-10-13")
	CodeNotFound             = NewUserErrorCode("not-found", "2024-10-13")
	CodeForbidden            = NewUserErrorCode("forbidden", "2024-10-13")
	CodeCarInvalid           = NewUserErrorCode("car-invalid", "2024-10-13")
	CodeSpotInvalid          = NewUserErrorCode("spot-invalid", "2024-10-13")
	CodeCountryNotSupported  = NewUserErrorCode("country-not-supported", "2024-10-13")
	CodeProvinceNotSupported = NewUserErrorCode("province-not-supported", "2024-10-21")
	CodeNoProfile            = NewUserErrorCode("no-profile", "2024-10-13")
	CodeUnhealthy            = NewUserErrorCode("unhealthy", "2024-10-14")
)

// Error code for clients.
type UserErrorCode struct {
	date time.Time
	code string
}

// Errors meant to be sent to clients.
//
// This allows the router to automatically filter server errors and prevent them
// from reaching clients.
type UserFacingError struct {
	code *UserErrorCode
	msg  string
}

// Creates a new UserErrorCode.
//
// `code` should be an URL-safe string that describes the error. This is meant
// for clients to parse in order to identify the exact error. Must not be empty.
//
// `date` specifies the date in which this error was defined.
//
// Panics if any preconditions are not met.
func NewUserErrorCode(code, date string) UserErrorCode {
	if code == "" {
		panic("empty error code")
	}
	d, err := time.Parse(time.DateOnly, date)
	if err != nil {
		panic(fmt.Errorf("invalid date %v: %w", date, err))
	}
	return UserErrorCode{
		code: code,
		date: d,
	}
}

// Returns the error code
func (c *UserErrorCode) Code() string {
	return c.code
}

// Returns the error date
func (c *UserErrorCode) Date() time.Time {
	return c.date
}

// Returns an URI that uniquely identifies the error
func (c *UserErrorCode) TypeURI() string {
	uri := &url.URL{
		Scheme: "tag",
		Opaque: fmt.Sprintf("%v,%v:problem:%v",
			errorAuthority, c.date.Format(time.DateOnly), c.code),
	}
	return uri.String()
}

// Creates an UserFacingError with the specified error code
func (c *UserErrorCode) WithMsg(msg string) error {
	return &UserFacingError{
		code: c,
		msg:  msg,
	}
}

// Implements the `error` protocol
func (e *UserFacingError) Error() string {
	return e.msg
}

// Returns an URI that uniquely identifies the error.
//
// Might be nil.
func (e *UserFacingError) Code() *UserErrorCode {
	return e.code
}
