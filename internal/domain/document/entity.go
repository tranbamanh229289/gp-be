package document

import (
	"time"

	"github.com/google/uuid"
)

type CitizenIdentity struct {
	ID               uint              `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID         uuid.UUID         `gorm:"column:public_id;type:uuid;uniqueIndex;default:uuid_generate_v4()" json:"public_id" validate:"required"`
	IDNumber         string            `gorm:"column:id_number;type:varchar(20);uniqueIndex" json:"id_number" validate:"required,len=12,numeric"`
	FirstName        string            `gorm:"column:first_name;type:varchar(100);not null" json:"first_name" validate:"required,min=1,max=100"`
	LastName         string            `gorm:"column:last_name;type:varchar(100);not null" json:"last_name" validate:"required,min=1,max=100"`
	Gender           string            `gorm:"column:gender;type:varchar(20);not null" json:"gender" validate:"required,oneof=male female other"`
	DateOfBirth      time.Time         `gorm:"column:date_of_birth;type:date;not null" json:"date_of_birth" validate:"required,ltecsfield=IssueDate"`
	PlaceOfBirth     string            `gorm:"column:place_of_birth;type:text" json:"place_of_birth" validate:"omitempty,max=255"`
	Status           string            `gorm:"column:status;type:varchar(30);default:'Active';index" json:"status" validate:"required,oneof=active expired revoked"`
	IssueDate        time.Time         `gorm:"column:issue_date;type:date;not null" json:"issue_date" validate:"required"`
	ExpiryDate       time.Time         `gorm:"column:expiry_date;type:date;not null" json:"expiry_date" validate:"required,gtefield=IssueDate"`
	IssuerDID        string            `gorm:"column:issuer_did;type:varchar(255)" json:"issuer_did" validate:"omitempty,startswith=did:"`
	HolderDID        string            `gorm:"column:holder_did;type:varchar(255);index;not null" json:"holder_did" validate:"required,startswith=did:"`
	CreatedAt        time.Time         `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt        time.Time         `gorm:"autoUpdateTime" json:"updated_at" validate:"-"`
	RevokedAt        *time.Time        `gorm:"type:timestamptz" json:"revoked_at,omitempty" validate:"omitempty"`
	AcademicDegrees  []AcademicDegree  `gorm:"foreignKey:CID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"academic_degrees,omitempty"`
	HealthInsurances []HealthInsurance `gorm:"foreignKey:CID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"health_insurances,omitempty"`
	DriverLicenses   []DriverLicense   `gorm:"foreignKey:CID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"driver_licenses,omitempty"`
	Passports        []Passport        `gorm:"foreignKey:CID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"passports,omitempty"`
}

type AcademicDegree struct {
	ID             uint             `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID       uuid.UUID        `gorm:"column:public_id;type:uuid;uniqueIndex;default:uuid_generate_v4()" json:"public_id" validate:"required"`
	CID            uint             `gorm:"column:cid;not null;index" json:"cid" validate:"required"`
	DegreeNumber   string           `gorm:"column:degree_number;type:varchar(50);uniqueIndex" json:"degree_number" validate:"required,max=50"`
	DegreeType     string           `gorm:"column:degree_type;type:varchar(50);not null" json:"degree_type" validate:"required,oneof=bachelor master phd associate_professor full_professor"`
	Major          string           `gorm:"column:major;type:varchar(255);not null" json:"major" validate:"required,min=2,max=255"`
	University     string           `gorm:"column:university;type:varchar(255);not null" json:"university" validate:"required,min=2,max=255"`
	GraduateYear   uint             `gorm:"column:graduate_year;type:smallint;not null" json:"graduate_year" validate:"required,gte=1900,lte=2100"`
	GPA            float32          `gorm:"column:gpa;type:decimal(3,2)" json:"gpa" validate:"omitempty,gte=0,lte=4"`
	Classification string           `gorm:"column:classification;type:varchar(50)" json:"classification" validate:"omitempty,oneof=excellent very_good good average pass"`
	Status         string           `gorm:"column:status;type:varchar(30);default:'Active';index" json:"status" validate:"required,oneof=active expired revoked"`
	IssueDate      time.Time        `gorm:"column:issue_date;type:date;not null" json:"issue_date" validate:"required,ltecsfield=GraduateYearDate"`
	IssuerDID      string           `gorm:"column:issuer_did;type:varchar(255)" json:"issuer_did" validate:"omitempty,startswith=did:"`
	CreatedAt      time.Time        `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt      time.Time        `gorm:"autoUpdateTime" json:"updated_at" validate:"-"`
	RevokedAt      *time.Time       `gorm:"type:timestamptz" json:"revoked_at,omitempty" validate:"omitempty"`
	Citizen        *CitizenIdentity `gorm:"foreignKey:CID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`

	// Helper field để validate GraduateYear (không lưu DB)
	GraduateYearDate time.Time `gorm:"-" json:"-" validate:"-"`
}

type HealthInsurance struct {
	ID              uint             `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID        uuid.UUID        `gorm:"column:public_id;type:uuid;uniqueIndex;default:uuid_generate_v4()" json:"public_id" validate:"required"`
	CID             uint             `gorm:"column:cid;not null;index" json:"cid" validate:"required"`
	InsuranceNumber string           `gorm:"column:insurance_number;type:varchar(20);uniqueIndex" json:"insurance_number" validate:"required,len=15,numeric"`
	InsuranceType   string           `gorm:"column:insurance_type;type:varchar(100)" json:"insurance_type" validate:"omitempty,max=100"`
	Hospital        string           `gorm:"column:hospital;type:varchar(255)" json:"hospital" validate:"omitempty,max=255"`
	Status          string           `gorm:"column:status;type:varchar(30);default:'Active';index" json:"status" validate:"required,oneof=active expired revoked"`
	StartDate       time.Time        `gorm:"column:start_date;type:date;not null" json:"start_date" validate:"required"`
	ExpiryDate      time.Time        `gorm:"column:expiry_date;type:date;not null" json:"expiry_date" validate:"required,gtefield=StartDate"`
	IssuerDID       string           `gorm:"column:issuer_did;type:varchar(255)" json:"issuer_did" validate:"omitempty,startswith=did:"`
	CreatedAt       time.Time        `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt       time.Time        `gorm:"autoUpdateTime" json:"updated_at" validate:"-"`
	RevokedAt       *time.Time       `gorm:"type:timestamptz" json:"revoked_at,omitempty" validate:"omitempty"`
	Citizen         *CitizenIdentity `gorm:"foreignKey:CID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
}

type DriverLicense struct {
	ID             uint             `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID       uuid.UUID        `gorm:"column:public_id;type:uuid;uniqueIndex;default:uuid_generate_v4()" json:"public_id" validate:"required"`
	CID            uint             `gorm:"column:cid;not null;index" json:"cid" validate:"required"`
	LicenseNumber  string           `gorm:"column:license_number;type:varchar(30);uniqueIndex" json:"license_number" validate:"required,max=30"`
	Class          string           `gorm:"column:class;type:varchar(20);not null" json:"class" validate:"required"`
	Point          int              `gorm:"column:point;type:smallint;default:12" json:"point" validate:"gte=0,lte=12"`
	PointResetDate time.Time        `gorm:"column:point_reset_date;type:date" json:"point_reset_date" validate:"omitempty"`
	Status         string           `gorm:"column:status;type:varchar(30);default:'Active';index" json:"status" validate:"required,oneof=active expired revoked"`
	IssueDate      time.Time        `gorm:"column:issue_date;type:date;not null" json:"issue_date" validate:"required"`
	ExpiryDate     time.Time        `gorm:"column:expiry_date;type:date;not null" json:"expiry_date" validate:"required,gtefield=IssueDate"`
	IssuerDID      string           `gorm:"column:issuer_did;type:varchar(255)" json:"issuer_did" validate:"omitempty,startswith=did:"`
	CreatedAt      time.Time        `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt      time.Time        `gorm:"autoUpdateTime" json:"updated_at" validate:"-"`
	RevokedAt      *time.Time       `gorm:"type:timestamptz" json:"revoked_at,omitempty" validate:"omitempty"`
	Citizen        *CitizenIdentity `gorm:"foreignKey:CID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
}

type Passport struct {
	ID             uint             `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID       uuid.UUID        `gorm:"column:public_id;type:uuid;uniqueIndex;default:uuid_generate_v4()" json:"public_id" validate:"required"`
	CID            uint             `gorm:"column:cid;not null;index" json:"cid" validate:"required"`
	PassportNumber string           `gorm:"column:passport_number;type:varchar(15);uniqueIndex" json:"passport_number" validate:"required,max=15,alphanum"`
	PassportType   string           `gorm:"column:passport_type;type:varchar(100)" json:"passport_type" validate:"omitempty,oneof=ordinary diplomatic official"`
	Nationality    string           `gorm:"column:nationality;type:char(3);not null" json:"nationality" validate:"required,len=3,alpha"`
	MRZ            string           `gorm:"column:mrz;type:text" json:"mrz" validate:"omitempty,min=88,max=88"`
	Status         string           `gorm:"column:status;type:varchar(30);default:'Active';index" json:"status" validate:"required,oneof=active expired revoked"`
	IssueDate      time.Time        `gorm:"column:issue_date;type:date;not null" json:"issue_date" validate:"required"`
	ExpiryDate     time.Time        `gorm:"column:expiry_date;type:date;not null" json:"expiry_date" validate:"required,gtefield=IssueDate"`
	IssuerDID      string           `gorm:"type:varchar(255)" json:"issuer_did" validate:"omitempty,startswith=did:"`
	CreatedAt      time.Time        `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt      time.Time        `gorm:"autoUpdateTime" json:"updated_at" validate:"-"`
	RevokedAt      *time.Time       `gorm:"type:timestamptz" json:"revoked_at,omitempty" validate:"omitempty"`
	Citizen        *CitizenIdentity `gorm:"foreignKey:CID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
}
