package bizerr

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errorspb "github.com/wjiec/alchemy/internal/errors"
)

// Error represents a business logic error that includes an HTTP error code.
type Error struct {
	code    uint32 // business error code
	status  uint32 // http status code
	message string
	cause   error
}

// Code returns the business error code associated with the error.
func (e *Error) Code() uint32 {
	return e.code
}

// Status returns the HTTP status code associated with the error.
func (e *Error) Status() uint32 {
	return e.status
}

// Error constructs and returns an error message string.
func (e *Error) Error() string {
	return e.message
}

// GRPCStatus translates the error into a gRPC status.
//
// It creates a new gRPC status from the business error code and includes
// the HTTP status code as additional details.
func (e *Error) GRPCStatus() *status.Status {
	statusErr := status.New(codes.Code(e.code), e.Error())
	statusErr, _ = statusErr.WithDetails(&errorspb.WithHttpStatus{HttpCode: e.status})
	if e.cause != nil {
		statusErr, _ = statusErr.WithDetails(&errorspb.WithCause{CauseError: e.cause.Error()})
	}

	return statusErr
}

// Equals checks if the provided error is equivalent to the current error.
func (e *Error) Equals(err error) bool {
	if bizErr := new(Error); errors.As(err, &bizErr) {
		return bizErr.code == e.code
	}
	if statusErr, ok := status.FromError(err); ok {
		return statusErr.Code() == codes.Code(e.code)
	}

	return false
}

// With creates a new *Error instance by wrapping a given error into the current error's cause.
func (e *Error) With(err error) *Error {
	return &Error{
		code:    e.code,
		status:  e.status,
		message: e.Error(),
		cause:   err,
	}
}

// New constructs and returns a new *Error instance.
func New(code, status uint32, template ...string) *Error {
	bizErr := &Error{code: code, status: status}
	if len(template) > 0 {
		bizErr.message = template[0]
	}

	return bizErr
}

// FromError extracts an *Error instance from a regular error.
func FromError(err error) (*Error, bool) {
	var bizErr Error
	if statusErr, ok := status.FromError(err); ok {
		bizErr.code = uint32(statusErr.Code())
		for _, detail := range statusErr.Details() {
			switch v := detail.(type) {
			case *errorspb.WithHttpStatus:
				bizErr.status = v.HttpCode
			case *errorspb.WithCause:
				bizErr.cause = errors.New(v.CauseError)
			}
		}
		bizErr.message = statusErr.Message()
	}

	if bizErr.status == 0 {
		return nil, false
	}
	return &bizErr, true
}
