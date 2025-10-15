package constant

import "net/http"

type Errors struct {
	Code string
	Message string
	HttpStatusCode int
}

func (e *Errors) Error() string { 
    return e.Message
}

var (
	// Server
	InternalServer = Errors{
		Code: "INTERNAL_SERVER",
		Message: "Internal server error",
		HttpStatusCode: http.StatusInternalServerError,
	}

	BadRequest = Errors {
		Code: "BAD_REQUEST",
		Message: "Bad request",
		HttpStatusCode: http.StatusBadGateway,
	}

	// Auth 
	Unauthorized = Errors{
		Code: "UNAUTHORIZED",
		Message: "Unauthorized error",
		HttpStatusCode: http.StatusUnauthorized,
	}
	Forbidden = Errors {
		Code: "FORBIDDEN",
		Message: "Forbidden error",
		HttpStatusCode: http.StatusForbidden,
	}
	InvalidAuthHeader = Errors{
		Code: "INVALID_AUTH_HEADER",
		Message: "Invalid auth header",
		HttpStatusCode: http.StatusUnauthorized,
	}
	InvalidToken = Errors{
		Code: "INVALID_TOKEN",
		Message: "Invalid token",
		HttpStatusCode: http.StatusUnauthorized,
	}
	UserNotFound = Errors{
		Code: "USER_NOT_FOUND",
		Message: "User not found error",
		HttpStatusCode: http.StatusNotFound,
	}
	UserExisted = Errors {
		Code: "USER_EXISTED",
		Message: "User existed",
		HttpStatusCode: http.StatusConflict,
	}
)