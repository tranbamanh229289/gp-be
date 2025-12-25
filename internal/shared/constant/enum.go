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
	OtherGender  = "other"
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
	PassportOrdinaryType   = "ordinary"
	PassportDiplomaticType = "diplomatic"
	PassportOfficialType   = "official"
)

// proof
const (
	ProofPendingStatus   = "pending"
	ProofCompletedStatus = "completed"
	ProofFailedStatus    = "failed"
	ProofExpiredStatus   = "expired"
	ProofCancelledStatus = "cancelled"
)

// verifiable credential
const (
	ProofBjjSignature2021Type  = "sig"
	Iden3SparseMerkleTreeProof = "mtp"
)
const (
	VerifiableCredentialPendingStatus = "pending"
	VerifiableCredentialIssuedStatus  = "issued"
	VerifiableCredentialRevokedStatus = "revoked"
	VerifiableCredentialExpiredStatus = "expired"
)

// credential request
const (
	CredentialRequestPendingStatus  = "pending"
	CredentialRequestApprovedStatus = "approved"
	CredentialRequestRejectedStatus = "rejected"
)

// schema
const (
	SchemaActiveStatus = "active"
	SchemaRevokeStatus = "revoked"
)
