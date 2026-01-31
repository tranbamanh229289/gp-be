package constant

import "net/http"

type Errors struct {
	Code    string
	Message string
	Status  int
}

func (e *Errors) Error() string {
	return e.Message
}

var (
	// Server
	InternalServer = Errors{
		Code:    "INTERNAL_SERVER",
		Message: "Internal server error",
		Status:  http.StatusInternalServerError,
	}

	BadRequest = Errors{
		Code:    "BAD_REQUEST",
		Message: "Bad request",
		Status:  http.StatusBadGateway,
	}

	// Auth
	Unauthorized = Errors{
		Code:    "UNAUTHORIZED",
		Message: "Unauthorized error",
		Status:  http.StatusUnauthorized,
	}
	Forbidden = Errors{
		Code:    "FORBIDDEN",
		Message: "Forbidden error",
		Status:  http.StatusForbidden,
	}
	InvalidAuthHeader = Errors{
		Code:    "INVALID_AUTH_HEADER",
		Message: "Invalid auth header",
		Status:  http.StatusUnauthorized,
	}
	InvalidToken = Errors{
		Code:    "INVALID_TOKEN",
		Message: "Invalid token",
		Status:  http.StatusUnauthorized,
	}
	UserNotFound = Errors{
		Code:    "USER_NOT_FOUND",
		Message: "User not found error",
		Status:  http.StatusNotFound,
	}
	UserExisted = Errors{
		Code:    "USER_EXISTED",
		Message: "User existed",
		Status:  http.StatusConflict,
	}

	// documents
	CitizenIdentityNotFound = Errors{
		Code:    "CITIZEN_IDENTITY_NOT_FOUND",
		Message: "Citizen identity not found error",
		Status:  http.StatusNotFound,
	}

	AcademicDegreeNotFound = Errors{
		Code:    "ACADEMIC_DEGREE_NOT_FOUND",
		Message: "Academic degree not found error",
		Status:  http.StatusNotFound,
	}

	HealthInsuranceNotFound = Errors{
		Code:    "HEALTH_INSURANCE_NOT_FOUND",
		Message: "Health insurance not found error",
		Status:  http.StatusNotFound,
	}

	DriverLicenseNotFound = Errors{
		Code:    "DRIVER_LICENSE_NOT_FOUND",
		Message: "Driver license not found error",
		Status:  http.StatusNotFound,
	}

	PassportNotFound = Errors{
		Code:    "PASSPORT_NOT_FOUND",
		Message: "Passport not found error",
		Status:  http.StatusNotFound,
	}

	// identity
	IdentityNotFound = Errors{
		Code:    "IDENTITY_NOT_FOUND",
		Message: "Identity not found error",
		Status:  http.StatusNotFound,
	}

	// schema
	SchemaNotFound = Errors{
		Code:    "SCHEMA_NOT_FOUND",
		Message: "Schema not found error",
		Status:  http.StatusNotFound,
	}

	SchemaAttributeNotFound = Errors{
		Code:    "SCHEMA_ATTRIBUTE_NOT_FOUND",
		Message: "Schema attribute not found error",
		Status:  http.StatusNotFound,
	}

	// credential_requests
	CredentialRequestNotFound = Errors{
		Code:    "CREDENTIAL_REQUEST_NOT_FOUND",
		Message: "Credential request not found error",
		Status:  http.StatusNotFound,
	}

	// verifiable_credential
	VerifiableCredentialNotFound = Errors{
		Code:    "VERIFIABLE_CREDENTIAL_NOT_FOUND",
		Message: "Verifiable credential not found error",
		Status:  http.StatusNotFound,
	}

	VerifiableCredentialNotSig = Errors{
		Code:    "VERIFIABLE_CREDENTIAL_NOT_SIG",
		Message: "Verifiable credential not sig error",
		Status:  http.StatusNotAcceptable,
	}

	VerifiableCredentialAlreadySig = Errors{
		Code:    "VERIFIABLE_CREDENTIAL_ALREADY_SIG",
		Message: "Verifiable credential already sig error",
		Status:  http.StatusNotAcceptable,
	}

	// proof
	ProofRequestNotFound = Errors{
		Code:    "PROOF_REQUEST_NOT_FOUND",
		Message: "Verifiable credential not found error",
		Status:  http.StatusNotFound,
	}

	ProofResponseNotFound = Errors{
		Code:    "PROOF_RESPONSE_NOT_FOUND",
		Message: "Proof response not found error",
		Status:  http.StatusNotFound,
	}

	// state transition
	StateTransition = Errors{
		Code:    "STATE_TRANSITION_NOT_FOUND",
		Message: "State transition not found error",
		Status:  http.StatusNotFound,
	}

	// proof request
	ProofNotFound = Errors{
		Code:    "PROOF_NOT_FOUND",
		Message: "Proof not found error",
		Status:  http.StatusNotFound,
	}

	// statistic
	StatisticNotFound = Errors{
		Code:    "STATISTIC_NOT_FOUND",
		Message: "Statistic not found error",
		Status:  http.StatusNotFound,
	}
)
