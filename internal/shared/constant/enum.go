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
	IdentityHolderRole   IdentityRole = "holder"
	IdentityVerifierRole IdentityRole = "verifier"
	IdentityIssuerRole   IdentityRole = "issuer"
)

// Document
type DocumentType string

const (
	CitizenIdentity DocumentType = "citizen_identity"
	AcademicDegree  DocumentType = "academic_degree"
	HealthInsurance DocumentType = "health_insurance"
	DriverLicense   DocumentType = "driver_license"
	Passport        DocumentType = "passport"
)

type Gender string

const (
	MaleGender   Gender = "male"
	FemaleGender Gender = "female"
	OtherGender  Gender = "other"
)

type DegreeType string

const (
	BachelorDegreeType           DegreeType = "bachelor"
	MasterDegreeType             DegreeType = "master"
	PhDDegreeType                DegreeType = "phd"
	AssociateProfessorDegreeType DegreeType = "associate_professor"
	FullProfessorDegreeType      DegreeType = "full_professor"
)

type Classification string

const (
	ExcellentClassification Classification = "excellent"
	VeryGoodClassification  Classification = "very_good"
	GoodClassification      Classification = "good"
	AverageClassification   Classification = "average"
	PassClassification      Classification = "pass"
)

type DocumentStatus string

const (
	DocumentActiveStatus  DocumentStatus = "active"
	DocumentRevokeStatus  DocumentStatus = "revoke"
	DocumentExpiredStatus DocumentStatus = "expired"
)

type PassportType string

const (
	PassportOrdinaryType   PassportType = "ordinary"
	PassportDiplomaticType PassportType = "diplomatic"
	PassportOfficialType   PassportType = "official"
)

// proof
type ProofType string

const (
	BjjSignature2021           ProofType = "BjjSignature2021"
	Iden3SparseMerkleTreeProof ProofType = "Iden3SparseMerkleTreeProof"
	Null                       ProofType = "null"
)

type ProofRequestStatus string

const (
	ProofRequestActiveStatus    ProofRequestStatus = "active"
	ProofRequestExpiredStatus   ProofRequestStatus = "expired"
	ProofRequestCancelledStatus ProofRequestStatus = "cancelled"
)

type ProofResponseStatus string

const (
	ProofResponsePendingStatus ProofResponseStatus = "pending"
	ProofResponseSuccessStatus ProofResponseStatus = "success"
	ProofResponseFailedStatus  ProofResponseStatus = "failed"
)

// verifiable credential
type VerifiableCredentialStatus string

const (
	VerifiableCredentialNotSignedStatus VerifiableCredentialStatus = "notSigned"
	VerifiableCredentialIssuedStatus    VerifiableCredentialStatus = "issued"
	VerifiableCredentialRevokedStatus   VerifiableCredentialStatus = "revoked"
	VerifiableCredentialExpiredStatus   VerifiableCredentialStatus = "expired"
)

// credential request
type CredentialRequestStatus string

const (
	CredentialRequestPendingStatus  CredentialRequestStatus = "pending"
	CredentialRequestApprovedStatus CredentialRequestStatus = "approved"
	CredentialRequestRejectedStatus CredentialRequestStatus = "rejected"
)

// schema
type SchemaStatus string

const (
	SchemaActiveStatus SchemaStatus = "active"
	SchemaRevokeStatus SchemaStatus = "revoked"
)

// Slot
type Slot string

const (
	SlotIndexA Slot = "slotIndexA"
	SlotIndexB Slot = "slotIndexB"
	SlotValueA Slot = "slotValueA"
	SlotValueB Slot = "slotValueB"
)

type AttributeType string

const (
	AttributeStringType  AttributeType = "string"
	AttributeNumberType  AttributeType = "number"
	AttributeIntegerType AttributeType = "integer"
	AttributeBooleanType AttributeType = "boolean"
	AttributeObjectType  AttributeType = "object"
	AttributeArrayType   AttributeType = "array"
)
