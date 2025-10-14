package constant

import "net/http"

type Errors struct {
	Code string
	Message string
	HttpStatusCode int
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
	UserNotFound = Errors{
		Code: "USER_NOT_FOUND",
		Message: "User not found error",
		HttpStatusCode: 404,
	}
)