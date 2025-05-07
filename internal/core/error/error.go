package custom_error

import (
	"fmt"

	"github.com/guregu/null"
)

const (
	ErrInternalServerError string = "ERR-APP-500"
	ErrBadRequest          string = "ERR-APP-400"
	ErrUnauthorized        string = "ERR-APP-401"
	ErrForbidden           string = "ERR-APP-403"
	ErrNotFound            string = "ERR-APP-404"
)

const (
	ErrAuthInvalidConfirmPassword string = "ERR-AUTH-001"
	ErrAuthInvalidEmailFormat     string = "ERR-AUTH-002"
	ErrAuthEmailAlreadyExists     string = "ERR-AUTH-003"
	ErrAuthInvalidCredentials     string = "ERR-AUTH-004"
	ErrAuthUnauthenticated        string = "ERR-AUTH-005"
)

const (
	ErrUserEmailAlreadyExists    = "ERR-USER-001"
	ErrUserInvalidEmailFormat    = "ERR-USER-002"
	ErrUserNotFound              = "ERR-USER-003"
	ErrUserInvalidateUpdateInput = "ERR-USER-004"
)

var HttpCodeMap = map[string]int{
	ErrInternalServerError:        500,
	ErrBadRequest:                 400,
	ErrUnauthorized:               401,
	ErrForbidden:                  403,
	ErrNotFound:                   404,
	ErrAuthInvalidConfirmPassword: 400,
	ErrAuthInvalidEmailFormat:     400,
	ErrAuthEmailAlreadyExists:     409,
	ErrAuthInvalidCredentials:     401,
	ErrAuthUnauthenticated:        401,
	ErrUserEmailAlreadyExists:     409,
	ErrUserInvalidEmailFormat:     400,
	ErrUserNotFound:               404,
	ErrUserInvalidateUpdateInput:  400,
}

var MessageMap = map[string]string{
	ErrInternalServerError:        "Internal server error",
	ErrBadRequest:                 "Bad request",
	ErrUnauthorized:               "Unauthorized",
	ErrForbidden:                  "Forbidden",
	ErrNotFound:                   "Not found",
	ErrAuthInvalidConfirmPassword: "Invalid confirm password",
	ErrAuthInvalidEmailFormat:     "Invalid email format",
	ErrAuthEmailAlreadyExists:     "Email already exists",
	ErrAuthInvalidCredentials:     "Invalid credentials",
	ErrAuthUnauthenticated:        "Unauthenticated",
	ErrUserEmailAlreadyExists:     "Email already exists",
	ErrUserInvalidEmailFormat:     "Invalid email format",
	ErrUserNotFound:               "User not found",
	ErrUserInvalidateUpdateInput:  "Invalid update input",
}

type CustomErrorInterface interface {
	error
	GetCode() string
	GetMessage() *string
	GetHttpCode() int
}

type CustomError struct {
	Code     string
	Message  *string
	HttpCode int
}

func NewError(code string, message *string) CustomErrorInterface {
	httpCode, ok := HttpCodeMap[code]
	if !ok {
		httpCode = 500
	}

	if message == nil {
		m, ok := MessageMap[code]
		if !ok {
			message = null.StringFrom("Unknown error").Ptr()
		} else {
			message = null.StringFrom(m).Ptr()
		}
	}

	return CustomError{
		Code:     code,
		Message:  message,
		HttpCode: httpCode,
	}
}

func (e CustomError) Error() string {
	if e.Message != nil {
		return fmt.Sprintf("%s: %s", e.Code, null.StringFromPtr(e.Message).String)
	}
	return fmt.Sprintf("%s: -", e.Code)
}

func (e CustomError) GetCode() string {
	return e.Code
}

func (e CustomError) GetMessage() *string {
	return e.Message
}

func (e CustomError) GetHttpCode() int {
	return e.HttpCode
}
