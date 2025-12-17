package constant

// Auth
type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

type TokenType string

const (
	AccessToken  TokenType = "access_token"
	RefreshToken TokenType = "refresh_token"
)

// Blockchain

// Credential
type CredentialType string

const (
	CitizenIdentityType CredentialType = "citizen_identity"
	AcademicDegreeType  CredentialType = "academic_degree"
	HealthInsuranceType CredentialType = "health_insurance"
	DriverLicenseType   CredentialType = "driver_license"
	PassportType        CredentialType = "passport"
)

const (
	ActiveStatus = "active"
	RevokeStatus = "revoke"
)
