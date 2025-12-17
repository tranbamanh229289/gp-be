package document

import (
	"time"

	"github.com/google/uuid"
)

type CitizenIdentity struct {
	ID               uint              `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID         uuid.UUID         `gorm:"type:uuid;uniqueIndex;default:uuid_generate_v4()" json:"public_id"`
	IDNumber         string            `gorm:"type:varchar(20);uniqueIndex" json:"id_number" validate:"required,len=12"`
	FirstName        string            `gorm:"type:varchar(100)" json:"first_name" validate:"required"`
	LastName         string            `gorm:"type:varchar(100)" json:"last_name" validate:"required"`
	Gender           string            `gorm:"type:varchar(20)" json:"gender" validate:"required,oneof=Male Female Other"`
	DateOfBirth      time.Time         `gorm:"type:date;not null" json:"date_of_birth" validate:"required"`
	PlaceOfBirth     string            `gorm:"type:text" json:"place_of_birth" validate:"omitempty"`
	Status           string            `gorm:"type:varchar(30);default:'Active'" json:"status" validate:"required,oneof=Active Expired Revoked"`
	IssueDate        time.Time         `gorm:"type:date;not null" json:"issue_date" validate:"required"`
	ExpiryDate       time.Time         `gorm:"type:date;not null" json:"expiry_date" validate:"required"`
	IssuerDID        string            `gorm:"type:varchar(255)" json:"issue_did" validate:"omitempty,startswith=did:"`
	HolderDID        string            `gorm:"type:varchar(255);index" json:"holder_did" validate:"required,startswith=did:"`
	CreatedAt        time.Time         `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt        time.Time         `gorm:"autoUpdateTime" json:"updated_at" validate:"-"`
	RevokedAt        *time.Time        `gorm:"type:timestamptz" json:"revoked_at,omitempty" validate:"omitempty"`
	AcademicDegrees  []AcademicDegree  `gorm:"foreignKey:CID" json:"academic_degrees,omitempty"`
	HealthInsurances []HealthInsurance `gorm:"foreignKey:CID" json:"health_insurances,omitempty"`
	DriverLicenses   []DriverLicense   `gorm:"foreignKey:CID" json:"driver_licenses,omitempty"`
	Passports        []Passport        `gorm:"foreignKey:CID" json:"passports,omitempty"`
}

type AcademicDegree struct {
	ID             uint             `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID       uuid.UUID        `gorm:"type:uuid;uniqueIndex;default:uuid_generate_v4()" json:"public_id"`
	CID            uint             `gorm:"not null;index" json:"cid" validate:"required"`
	Citizen        *CitizenIdentity `gorm:"foreignKey:CID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	DegreeNumber   string           `gorm:"type:varchar(50);uniqueIndex" json:"degree_number" validate:"required"`
	DegreeType     string           `gorm:"type:varchar(50)" json:"degree_type" validate:"required,oneof=Bachelor Master PhD AssociateProfessor FullProfessor"`
	Major          string           `gorm:"type:varchar(255)" json:"major" validate:"required"`
	University     string           `gorm:"type:varchar(255)" json:"university" validate:"required"`
	GraduateYear   uint             `gorm:"type:smallint" json:"graduate_year" validate:"required,gte=1950,lte=2100"`
	GPA            float32          `gorm:"type:decimal(3,2)" json:"gpa" validate:"omitempty"`
	Classification string           `gorm:"type:varchar(50)" json:"classification" validate:"omitempty,oneof=Excellent VeryGood Good Average Pass"`
	Status         string           `gorm:"type:varchar(30);default:'Active'" json:"status" validate:"required,oneof=Active Expired Revoked"`
	IssueDate      time.Time        `gorm:"type:date;not null" json:"issue_date" validate:"required"`
	IssuerDID      string           `gorm:"type:varchar(255)" json:"issuer_did" validate:"omitempty,startswith=did:"`
	CreatedAt      time.Time        `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt      time.Time        `gorm:"autoUpdateTime" json:"updated_at" validate:"-"`
	RevokedAt      *time.Time       `gorm:"type:timestamptz" json:"revoked_at,omitempty" validate:"omitempty"`
}

type HealthInsurance struct {
	ID              uint             `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID        uuid.UUID        `gorm:"type:uuid;uniqueIndex;default:uuid_generate_v4()" json:"public_id"`
	CID             uint             `gorm:"not null;index" json:"cid" validate:"required"`
	Citizen         *CitizenIdentity `gorm:"foreignKey:CID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	InsuranceNumber string           `gorm:"type:varchar(20);uniqueIndex" json:"insurance_number" validate:"required,len=15"`
	InsuranceType   string           `gorm:"type:varchar(100)" json:"insurance_type" validate:"omitempty"`
	Hospital        string           `gorm:"type:varchar(255)" json:"hospital" validate:"omitempty"`
	Status          string           `gorm:"type:varchar(30);default:'Active'" json:"status" validate:"required,oneof=Active Expired Revoked"`
	StartDate       time.Time        `gorm:"type:date;not null" json:"start_date" validate:"required"`
	ExpiryDate      time.Time        `gorm:"type:date;not null" json:"expiry_date" validate:"required"`
	IssuerDID       string           `gorm:"type:varchar(255)" json:"issuer_did" validate:"omitempty,startswith=did:"`
	CreatedAt       time.Time        `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt       time.Time        `gorm:"autoUpdateTime" json:"updated_at" validate:"-"`
	RevokedAt       *time.Time       `gorm:"type:timestamptz" json:"revoked_at,omitempty" validate:"omitempty"`
}

type DriverLicense struct {
	ID             uint             `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID       uuid.UUID        `gorm:"type:uuid;uniqueIndex;default:uuid_generate_v4()" json:"public_id"`
	CID            uint             `gorm:"not null;index" json:"cid" validate:"required"`
	Citizen        *CitizenIdentity `gorm:"foreignKey:CID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	LicenseNumber  string           `gorm:"type:varchar(30);uniqueIndex" json:"license_number" validate:"required"`
	Class          string           `gorm:"type:varchar(20)" json:"class" validate:"required"`
	Point          int              `gorm:"type:smallint" json:"point" validate:"gte=0,lte=12"`
	PointResetDate time.Time        `gorm:"type:date" json:"point_reset_date" validate:"omitempty"`
	Status         string           `gorm:"type:varchar(30);default:'Active'" json:"status" validate:"required,oneof=Active Expired Revoked Suspended"`
	IssueDate      time.Time        `gorm:"type:date;not null" json:"issue_date" validate:"required"`
	ExpiryDate     time.Time        `gorm:"type:date;not null" json:"expiry_date" validate:"required"`
	IssuerDID      string           `gorm:"type:varchar(255)" json:"issuer_did" validate:"omitempty,startswith=did:"`
	CreatedAt      time.Time        `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt      time.Time        `gorm:"autoUpdateTime" json:"updated_at" validate:"-"`
	RevokedAt      *time.Time       `gorm:"type:timestamptz" json:"revoked_at,omitempty" validate:"omitempty"`
}

type Passport struct {
	ID             uint             `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID       uuid.UUID        `gorm:"type:uuid;uniqueIndex;default:uuid_generate_v4()" json:"public_id"`
	CID            uint             `gorm:"not null;index" json:"cid" validate:"required"`
	Citizen        *CitizenIdentity `gorm:"foreignKey:CID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	PassportNumber string           `gorm:"type:varchar(15);uniqueIndex" json:"passport_number" validate:"required"`
	PassportType   string           `gorm:"type:varchar(100)" json:"passport_type" validate:"omitempty"`
	Nationality    string           `gorm:"type:char(3)" json:"nationality" validate:"required,len=3"`
	MRZ            string           `gorm:"type:text" json:"mrz" validate:"omitempty,min=60,max=100"`
	Status         string           `gorm:"type:varchar(30);default:'Active'" json:"status" validate:"required,oneof=Active Expired Revoked Lost Stolen"`
	IssueDate      time.Time        `gorm:"type:date;not null" json:"issue_date" validate:"required"`
	ExpiryDate     time.Time        `gorm:"type:date;not null" json:"expiry_date" validate:"required"`
	IssuerDID      string           `gorm:"type:varchar(255)" json:"issuer_did" validate:"omitempty,startswith=did:"`
	CreatedAt      time.Time        `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt      time.Time        `gorm:"autoUpdateTime" json:"updated_at" validate:"-"`
	RevokedAt      *time.Time       `gorm:"type:timestamptz" json:"revoked_at,omitempty" validate:"omitempty"`
}
