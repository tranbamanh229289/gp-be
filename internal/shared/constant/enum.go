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
type IdentityRole string

const (
	IdentityHolderRole   = "holder"
	IdentityVerifierRole = "verifier"
	IdentityIssuerRole   = "issuer"
)

// Document
type DocumentType string

const (
	CitizenIdentity = "citizen_identity"
	AcademicDegree  = "academic_degree"
	HealthInsurance = "health_insurance"
	DriverLicense   = "driver_license"
	Passport        = "passport"
)

type Gender string

const (
	MaleGender   = "male"
	FemaleGender = "female"
	OtherGender  = "other"
)

type DegreeType string

const (
	BachelorDegreeType           = "bachelor"
	MasterDegreeType             = "master"
	PhDDegreeType                = "phd"
	AssociateProfessorDegreeType = "associate_professor"
	FullProfessorDegreeType      = "full_professor"
)

type Classification string

const (
	ExcellentClassification = "excellent"
	VeryGoodClassification  = "very_good"
	GoodClassification      = "good"
	AverageClassification   = "average"
	PassClassification      = "pass"
)

type DocumentStatus string

const (
	DocumentActiveStatus  = "active"
	DocumentRevokeStatus  = "revoke"
	DocumentExpiredStatus = "expired"
)

type PassportType string

const (
	PassportOrdinaryType   = "ordinary"
	PassportDiplomaticType = "diplomatic"
	PassportOfficialType   = "official"
)

// proof
type ProofRequestStatus string

const (
	ProofRequestActiveStatus    = "active"
	ProofRequestExpiredStatus   = "expired"
	ProofRequestCancelledStatus = "cancelled"
)

type ProofResponseStatus string

const (
	ProofResponsePendingStatus = "pending"
	ProofResponseSuccessStatus = "success"
	ProofResponseFailedStatus  = "failed"
)

// verifiable credential
type VerifiableCredentialType string

const (
	VerifiableCredentialNotSignedStatus = "notSigned"
	VerifiableCredentialIssuedStatus    = "issued"
	VerifiableCredentialRevokedStatus   = "revoked"
	VerifiableCredentialExpiredStatus   = "expired"
)

// credential request
type CredentialRequestStatus string

const (
	CredentialRequestPendingStatus  = "pending"
	CredentialRequestApprovedStatus = "approved"
	CredentialRequestRejectedStatus = "rejected"
)

// schema
type SchemaType string

const (
	SchemaActiveStatus = "active"
	SchemaRevokeStatus = "revoked"
)

// Slot
type Slot string

const (
	SlotIndexA Slot = "slotIndexA"
	SlotIndexB Slot = "slotIndexB"
	SlotDataA  Slot = "slotDataA"
	SlotDataB  Slot = "slotDataB"
)
