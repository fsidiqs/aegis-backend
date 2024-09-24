package apperror

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/fatih/structs"
)

// Type holds a type string and integer code for the error
type Type string

// "Set" of valid errorTypes
const (
	Authorization        Type = "AUTHORIZATION" // Authentication Failures -
	BadRequest           Type = "BAD_REQUEST"   // Validation errors / BadInput
	Internal             Type = "INTERNAL"      // Server (500) and fallback errors
	InternalConflict     Type = "INTERNAL_CONFLICT"
	NotFound             Type = "NOT_FOUND"              // For not finding resource
	PayloadTooLarge      Type = "PAYLOAD_TOO_LARGE"      // for uploading tons of JSON, or an image over the limit - 413
	ServiceUnavailable   Type = "SERVICE_UNAVAILABLE"    // For long running handlers
	UnsupportedMediaType Type = "UNSUPPORTED_MEDIA_TYPE" // for http 415
	TUserNotFound        Type = "USER_NOT_FOUND"         // for result when query an non-existsting users
	TUserUnverified      Type = "USER_UNVERIFIED"
)

const (
	ResourceNotFound    string = "resource_not_found"
	ResourceAlrExist    string = "resource already exist"
	BadRequestMessage   string = "invalid request parameters"
	InvalidToken        string = "provided token is invalid"
	InvalidRefreshToken string = "provided refresh token is invalid"
	InvalidPublicToken  string = "provided public token is invalid"
	ErrClaimParse       string = "couldn't parse claims"
	NoAuthHeader        string = "no authentication header provided"
	UserUnverified      string = "user is unverified"
	UserNotActive       string = "user is not active"
	MsgUserNotFound     string = "user not found"

	MsgUnhandledPaymentStatus             string = "payment status unhandled"
	MsgMaxInProgressUserEnrollmentReached string = "max user enrollment reached"
)

const (
	TUserNotActive string = "user_not_active"
	TResourceEmpty Type   = "resource_is_empty"
)

type InvalidArgument struct {
	Field string `json:"field"`
	Param string `json:"param"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

func InvalidArgumentsMap(args []InvalidArgument) []map[string]interface{} {
	argsLen := len(args)
	mapsRet := make([]map[string]interface{}, argsLen)
	for i := 0; i < argsLen; i++ {
		mapsRet[i] = structs.Map(args[i])
	}
	return mapsRet
}

// Error holds a custom error for the application
// which is helpful in returning a consistent
// error type/message from API endpoints
type Error struct {
	Type    Type   `json:"type"`
	Message string `json:"message"`
}

// Error satisfies standard error interface
// we can return errors from this package as
// a regular old go _error_
func (e *Error) Error() string {
	return e.Message
}

// Status is a mapping errors to status codes
// Of course, this is somewhat redundant since
// our errors already map http status codes
func (e *Error) Status() int {

	switch e.Type {
	case Authorization:
		return http.StatusUnauthorized
	case BadRequest:
		return http.StatusBadRequest
	case InternalConflict:
		return http.StatusBadRequest
	case Internal:
		return http.StatusInternalServerError
	case NotFound:
		return http.StatusBadRequest
	case PayloadTooLarge:
		return http.StatusRequestEntityTooLarge
	case ServiceUnavailable:
		return http.StatusServiceUnavailable
	case UnsupportedMediaType:
		return http.StatusUnsupportedMediaType
	case TUserNotFound:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// Status checks the runtime type
// of the error and returns an http
// status code if the error is model.Error
func Status(err error) int {
	var e *Error
	if errors.As(err, &e) {
		return e.Status()
	}
	return http.StatusInternalServerError
}

/*
* Error "Factories"
 */

// NewAuthorization to create a 401
func NewAuthorization(reason string) *Error {
	return &Error{
		Type:    Authorization,
		Message: reason,
	}
}

// NewBadRequest to create 400 errors (validation, for example)
func NewBadRequest(reason string) *Error {
	return &Error{
		Type:    BadRequest,
		Message: reason,
	}
}

func NewConflictSimple() *Error {
	return &Error{
		Type:    InternalConflict,
		Message: ResourceAlrExist,
	}
}

func NewConflictMsg(msg string) *Error {
	return &Error{
		Type:    InternalConflict,
		Message: msg,
	}
}

// NewInternal for 500 errors and unknown errors
func NewInternal() *Error {
	return &Error{
		Type:    Internal,
		Message: "internal server error",
	}
}

func NewInternalWrap(errMsg string) *Error {
	return &Error{
		Type:    Internal,
		Message: errMsg,
	}
}

// NewNotFound to create an error for 400
func NewResourceNotFound() *Error {
	return &Error{
		Type:    NotFound,
		Message: ResourceNotFound,
	}
}

func NewResourceNotFoundMsg(msg string) *Error {
	return &Error{
		Type:    NotFound,
		Message: msg,
	}
}

func NewWrapErrorMsg(err *Error, errStack string) *Error {
	err.Message = errStack
	return err
}

// NewPayloadTooLarge to create an error for 413
func NewPayloadTooLarge(maxBodySize int64, contentLength int64) *Error {
	return &Error{
		Type:    PayloadTooLarge,
		Message: fmt.Sprintf("Max payload size of %v exceeded. Actual payload size: %v", maxBodySize, contentLength),
	}
}

// NewServiceUnavailable to create an error for 503
func NewServiceUnavailable() *Error {
	return &Error{
		Type:    ServiceUnavailable,
		Message: "Service unavailable or timed out",
	}
}

// NewUnsupportedMediaType to create an error for 415
func NewUnsupportedMediaType(reason string) *Error {
	return &Error{
		Type:    UnsupportedMediaType,
		Message: reason,
	}
}

func ErrorWrapper(err error, msg string) string {
	return fmt.Sprintf("%s;%s", msg, err)
}
