package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode is a unique error identifier
type ErrorCode string

// HTTP status code constants
const (
	// General errors
	ErrBadRequest           ErrorCode = "BAD_REQUEST"
	ErrUnauthorized         ErrorCode = "UNAUTHORIZED"
	ErrForbidden            ErrorCode = "FORBIDDEN"
	ErrNotFound             ErrorCode = "NOT_FOUND"
	ErrConflict             ErrorCode = "CONFLICT"
	ErrInternalServer       ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrServiceUnavailable   ErrorCode = "SERVICE_UNAVAILABLE"

	// User related errors
	ErrUserNotFound         ErrorCode = "USER_NOT_FOUND"
	ErrInvalidCredentials   ErrorCode = "INVALID_CREDENTIALS"
	ErrUserAlreadyExists    ErrorCode = "USER_ALREADY_EXISTS"

	// Product related errors
	ErrProductNotFound      ErrorCode = "PRODUCT_NOT_FOUND"
	ErrInvalidProduct       ErrorCode = "INVALID_PRODUCT"

	// Order related errors
	ErrOrderNotFound        ErrorCode = "ORDER_NOT_FOUND"
	ErrInvalidOrder         ErrorCode = "INVALID_ORDER"

	// Inventory related errors
	ErrOutOfStock           ErrorCode = "OUT_OF_STOCK"

	// Payment related errors
	ErrPaymentFailed        ErrorCode = "PAYMENT_FAILED"
)

// Error is the standard error type for the system
type Error struct {
	Code     ErrorCode `json:"code"`
	Message  string    `json:"message"`
	HTTPCode int       `json:"-"`
	Err      error     `json:"-"`
}

// Implements error interface
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *Error) Unwrap() error {
	return e.Err
}

// New creates a new error
func New(code ErrorCode, message string, httpCode int, err error) *Error {
	return &Error{
		Code:     code,
		Message:  message,
		HTTPCode: httpCode,
		Err:      err,
	}
}

// NewBadRequest creates a 400 error
func NewBadRequest(message string, err error) *Error {
	return New(ErrBadRequest, message, http.StatusBadRequest, err)
}

// NewUnauthorized creates a 401 error
func NewUnauthorized(message string, err error) *Error {
	return New(ErrUnauthorized, message, http.StatusUnauthorized, err)
}

// NewForbidden creates a 403 error
func NewForbidden(message string, err error) *Error {
	return New(ErrForbidden, message, http.StatusForbidden, err)
}

// NewNotFound creates a 404 error
func NewNotFound(message string, err error) *Error {
	return New(ErrNotFound, message, http.StatusNotFound, err)
}

// NewConflict creates a 409 error
func NewConflict(message string, err error) *Error {
	return New(ErrConflict, message, http.StatusConflict, err)
}

// NewInternalServerError creates a 500 error
func NewInternalServerError(message string, err error) *Error {
	return New(ErrInternalServer, message, http.StatusInternalServerError, err)
}

// NewServiceUnavailable creates a 503 error
func NewServiceUnavailable(message string, err error) *Error {
	return New(ErrServiceUnavailable, message, http.StatusServiceUnavailable, err)
}
