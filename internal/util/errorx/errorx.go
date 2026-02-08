package errorx

import (
	"errors"
	"fmt"
	"strings"
)

// Error is a custom error type with additional context
type Error struct {
	Code    string
	Message string
	Cause   error
	Details map[string]interface{}
}

// Error implements the error interface
func (e *Error) Error() string {
	var sb strings.Builder
	sb.WriteString(e.Message)
	if e.Code != "" {
		sb.WriteString(" (")
		sb.WriteString(e.Code)
		sb.WriteString(")")
	}
	if e.Cause != nil {
		sb.WriteString(": ")
		sb.WriteString(e.Cause.Error())
	}
	return sb.String()
}

// Unwrap returns the underlying cause
func (e *Error) Unwrap() error {
	return e.Cause
}

// WithCause adds a cause to the error
func (e *Error) WithCause(cause error) *Error {
	e.Cause = cause
	return e
}

// WithDetail adds a detail to the error
func (e *Error) WithDetail(key string, value interface{}) *Error {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// New creates a new error
func New(message string) *Error {
	return &Error{Message: message}
}

// Newf creates a new error with format
func Newf(format string, args ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, args...)}
}

// WithCode creates an error with code
func WithCode(code, message string) *Error {
	return &Error{Code: code, Message: message}
}

// Wrap wraps an error with message
func Wrap(err error, message string) *Error {
	if err == nil {
		return nil
	}
	return &Error{Message: message, Cause: err}
}

// Wrapf wraps an error with formatted message
func Wrapf(err error, format string, args ...interface{}) *Error {
	if err == nil {
		return nil
	}
	return &Error{Message: fmt.Sprintf(format, args...), Cause: err}
}

// Is checks if error is of type Error with matching code
func Is(err error, code string) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Code == code
	}
	return false
}

// GetCode extracts code from error
func GetCode(err error) string {
	var e *Error
	if errors.As(err, &e) {
		return e.Code
	}
	return ""
}

// GetDetails extracts details from error
func GetDetails(err error) map[string]interface{} {
	var e *Error
	if errors.As(err, &e) {
		return e.Details
	}
	return nil
}

// AsError casts error to Error type
func AsError(err error) (*Error, bool) {
	var e *Error
	ok := errors.As(err, &e)
	return e, ok
}

// Common error codes
const (
	CodeNotFound      = "NOT_FOUND"
	CodeInvalidInput  = "INVALID_INPUT"
	CodeUnauthorized  = "UNAUTHORIZED"
	CodeForbidden     = "FORBIDDEN"
	CodeConflict      = "CONFLICT"
	CodeInternal      = "INTERNAL_ERROR"
	CodeTimeout       = "TIMEOUT"
	CodeCancelled     = "CANCELLED"
	CodeUnavailable   = "UNAVAILABLE"
	CodeNotImplemented = "NOT_IMPLEMENTED"
)

// NotFound creates a not found error
func NotFound(message string) *Error {
	return WithCode(CodeNotFound, message)
}

// NotFoundf creates a not found error with format
func NotFoundf(format string, args ...interface{}) *Error {
	return WithCode(CodeNotFound, fmt.Sprintf(format, args...))
}

// InvalidInput creates an invalid input error
func InvalidInput(message string) *Error {
	return WithCode(CodeInvalidInput, message)
}

// InvalidInputf creates an invalid input error with format
func InvalidInputf(format string, args ...interface{}) *Error {
	return WithCode(CodeInvalidInput, fmt.Sprintf(format, args...))
}

// Unauthorized creates an unauthorized error
func Unauthorized(message string) *Error {
	return WithCode(CodeUnauthorized, message)
}

// Forbidden creates a forbidden error
func Forbidden(message string) *Error {
	return WithCode(CodeForbidden, message)
}

// Conflict creates a conflict error
func Conflict(message string) *Error {
	return WithCode(CodeConflict, message)
}

// Internal creates an internal error
func Internal(message string) *Error {
	return WithCode(CodeInternal, message)
}

// Timeout creates a timeout error
func Timeout(message string) *Error {
	return WithCode(CodeTimeout, message)
}

// Cancelled creates a cancelled error
func Cancelled(message string) *Error {
	return WithCode(CodeCancelled, message)
}

// Unavailable creates an unavailable error
func Unavailable(message string) *Error {
	return WithCode(CodeUnavailable, message)
}

// NotImplemented creates a not implemented error
func NotImplemented(message string) *Error {
	return WithCode(CodeNotImplemented, message)
}

// IsNotFound checks if error is not found
func IsNotFound(err error) bool {
	return Is(err, CodeNotFound)
}

// IsInvalidInput checks if error is invalid input
func IsInvalidInput(err error) bool {
	return Is(err, CodeInvalidInput)
}

// IsUnauthorized checks if error is unauthorized
func IsUnauthorized(err error) bool {
	return Is(err, CodeUnauthorized)
}

// IsForbidden checks if error is forbidden
func IsForbidden(err error) bool {
	return Is(err, CodeForbidden)
}

// IsConflict checks if error is conflict
func IsConflict(err error) bool {
	return Is(err, CodeConflict)
}

// IsInternal checks if error is internal
func IsInternal(err error) bool {
	return Is(err, CodeInternal)
}

// IsTimeout checks if error is timeout
func IsTimeout(err error) bool {
	return Is(err, CodeTimeout)
}

// IsCancelled checks if error is cancelled
func IsCancelled(err error) bool {
	return Is(err, CodeCancelled)
}

// IsUnavailable checks if error is unavailable
func IsUnavailable(err error) bool {
	return Is(err, CodeUnavailable)
}

// IsNotImplemented checks if error is not implemented
func IsNotImplemented(err error) bool {
	return Is(err, CodeNotImplemented)
}

// MultiError combines multiple errors
type MultiError struct {
	Errors []error
}

// Add adds an error
func (m *MultiError) Add(err error) {
	if err != nil {
		m.Errors = append(m.Errors, err)
	}
}

// Error implements the error interface
func (m *MultiError) Error() string {
	if len(m.Errors) == 0 {
		return ""
	}
	messages := make([]string, len(m.Errors))
	for i, err := range m.Errors {
		messages[i] = err.Error()
	}
	return strings.Join(messages, "; ")
}

// HasErrors checks if there are errors
func (m *MultiError) HasErrors() bool {
	return len(m.Errors) > 0
}

// ToError returns nil if no errors, otherwise returns the multi error
func (m *MultiError) ToError() error {
	if !m.HasErrors() {
		return nil
	}
	if len(m.Errors) == 1 {
		return m.Errors[0]
	}
	return m
}

// Combine combines multiple errors into one
func Combine(errs ...error) error {
	m := &MultiError{}
	for _, err := range errs {
		m.Add(err)
	}
	return m.ToError()
}
