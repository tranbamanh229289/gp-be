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

// Identity
const (
	IdentityHolderRole   = "holder"
	IdentityVerifierRole = "verifier"
	IdentityIssuerRole   = "issuer"
)

// Credential
const (
	CitizenIdentityType = "citizen_identity"
	AcademicDegreeType  = "academic_degree"
	HealthInsuranceType = "health_insurance"
	DriverLicenseType   = "driver_license"
	PassportType        = "passport"
)

const (
	MaleGender   = "male"
	FemaleGender = "female"
)

const (
	BachelorDegreeType           = "bachelor"
	MasterDegreeType             = "master"
	PhDDegreeType                = "phd"
	AssociateProfessorDegreeType = "associate_professor"
	FullProfessorDegreeType      = "full_professor"
)

const (
	ExcellentClassification = "excellent"
	VeryGoodClassification  = "very_good"
	GoodClassification      = "good"
	AverageClassification   = "average"
	PassClassification      = "pass"
)

const (
	DocumentActiveStatus  = "active"
	DocumentRevokeStatus  = "revoke"
	DocumentExpiredStatus = "expired"
)

const (
	CredentialRequestPendingStatus  = "pending"
	CredentialRequestApprovedStatus = "approved"
	CredentialRequestRejectedStatus = "rejected"
)

// proof
const (
	ProofPendingStatus   = "pending"
	ProofCompletedStatus = "completed"
	ProofFailedStatus    = "failed"
	ProofExpiredStatus   = "expired"
	ProofCancelledStatus = "cancelled"
)

const (
	SigProofType = "sig"
	MTPProofType = "mtp"
)

// schema
const (
	SchemaActiveStatus = "active"
	SchemaRevokeStatus = "revoked"
)
