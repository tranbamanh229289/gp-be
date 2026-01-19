package dto

import (
	"be/internal/domain/document"
	"be/internal/shared/constant"
)

// Citizen Identity
type CitizenIdentityCreatedRequestDto struct {
	FirstName    string          `json:"firstName"`
	LastName     string          `json:"lastName"`
	Gender       constant.Gender `json:"gender"`
	DateOfBirth  int64           `json:"dateOfBirth"`
	PlaceOfBirth string          `json:"placeOfBirth"`
	IssueDate    int64           `json:"issueDate"`
	ExpiryDate   int64           `json:"expiryDate"`
	HolderDID    string          `json:"holderDID"`
	IssuerDID    string          `json:"issuerDID"`
}

type CitizenIdentityUpdatedRequestDto struct {
	FirstName    string          `json:"firstName,omitempty"`
	LastName     string          `json:"lastName,omitempty"`
	Gender       constant.Gender `json:"gender,omitempty"`
	DateOfBirth  int64           `json:"dateOfBirth,omitempty"`
	PlaceOfBirth string          `json:"placeOfBirth,omitempty"`
	IssueDate    int64           `json:"issueDate,omitempty"`
	ExpiryDate   int64           `json:"expiryDate,omitempty"`
}

type CitizenIdentityOptionRequestDto struct {
	Status string `json:"status"`
}

type CitizenIdentityResponseDto struct {
	PublicID     string                  `json:"id"`
	IDNumber     string                  `json:"idNumber"`
	Status       constant.DocumentStatus `json:"status"`
	FirstName    string                  `json:"firstName"`
	LastName     string                  `json:"lastName"`
	Gender       constant.Gender         `json:"gender"`
	DateOfBirth  int64                   `json:"dateOfBirth"`
	PlaceOfBirth string                  `json:"placeOfBirth"`
	IssueDate    int64                   `json:"issueDate"`
	ExpiryDate   int64                   `json:"expiryDate"`
	HolderDID    string                  `json:"holderDID"`
	IssuerDID    string                  `json:"issuerDID"`
}

// Academic Degree
type AcademicDegreeCreatedRequestDto struct {
	DegreeType     constant.DegreeType `json:"degreeType"`
	Major          string              `json:"major"`
	University     string              `json:"university"`
	GraduateYear   uint                `json:"graduateYear"`
	GPA            float32             `json:"gpa"`
	Classification string              `json:"classification"`
	IssueDate      int64               `json:"issueDate"`
	HolderDID      string              `json:"holderDID"`
	IssuerDID      string              `json:"issuerDID"`
}

type AcademicDegreeUpdatedRequestDto struct {
	DegreeType     constant.DegreeType `json:"degreeType,omitempty"`
	Major          string              `json:"major,omitempty"`
	University     string              `json:"university,omitempty"`
	GraduateYear   uint                `json:"graduateYear,omitempty"`
	GPA            float32             `json:"gpa,omitempty"`
	Classification string              `json:"classification,omitempty"`
	IssueDate      int64               `json:"issueDate,omitempty"`
}

type AcademicDegreeOptionRequestDto struct {
	Status string `json:"status"`
}

type AcademicDegreeResponseDto struct {
	PublicID       string                  `json:"id"`
	DegreeNumber   string                  `json:"degreeNumber"`
	Status         constant.DocumentStatus `json:"status"`
	DegreeType     constant.DegreeType     `json:"degreeType"`
	Major          string                  `json:"major"`
	University     string                  `json:"university"`
	GraduateYear   uint                    `json:"graduateYear"`
	GPA            float32                 `json:"gpa"`
	Classification string                  `json:"classification"`
	IssueDate      int64                   `json:"issueDate"`
	HolderDID      string                  `json:"holderDID"`
	IssuerDID      string                  `json:"issuerDID"`
}

// Health Insurance
type HealthInsuranceCreatedRequestDto struct {
	InsuranceType string `json:"insuranceType"`
	Hospital      string `json:"hospital"`
	StartDate     int64  `json:"startDate"`
	ExpiryDate    int64  `json:"expiryDate"`
	HolderDID     string `json:"holderDID"`
	IssuerDID     string `json:"issuerDID"`
}

type HealthInsuranceUpdatedRequestDto struct {
	InsuranceType string `json:"insuranceType,omitempty"`
	Hospital      string `json:"hospital,omitempty"`
	StartDate     int64  `json:"startDate,omitempty"`
	ExpiryDate    int64  `json:"expiryDate,omitempty"`
}

type HealthInsuranceOptionRequestDto struct {
	Status string `json:"status"`
}

type HealthInsuranceResponseDto struct {
	PublicID        string                  `json:"id"`
	InsuranceNumber string                  `json:"insuranceNumber"`
	Status          constant.DocumentStatus `json:"status"`
	InsuranceType   string                  `json:"insuranceType"`
	Hospital        string                  `json:"hospital"`
	StartDate       int64                   `json:"startDate"`
	ExpiryDate      int64                   `json:"expiryDate"`
	HolderDID       string                  `json:"holderDID"`
	IssuerDID       string                  `json:"issuerDID"`
}

// Driver License
type DriverLicenseCreatedRequestDto struct {
	Class      string `json:"class"`
	IssueDate  int64  `json:"issueDate"`
	ExpiryDate int64  `json:"expiryDate"`
	HolderDID  string `json:"holderDID"`
	IssuerDID  string `json:"issuerDID"`
}

type DriverLicenseUpdatedRequestDto struct {
	Point      uint   `json:"point,omitempty"`
	Class      string `json:"class,omitempty"`
	IssueDate  int64  `json:"issueDate,omitempty"`
	ExpiryDate int64  `json:"expiryDate,omitempty"`
}

type DriverLicenseOptionRequestDto struct {
	Status string `json:"status"`
}

type DriverLicenseResponseDto struct {
	PublicID      string                  `json:"id"`
	LicenseNumber string                  `json:"licenseNumber"`
	Status        constant.DocumentStatus `json:"status"`
	Point         uint                    `json:"point"`
	Class         string                  `json:"class"`
	IssueDate     int64                   `json:"issueDate"`
	ExpiryDate    int64                   `json:"expiryDate"`
	HolderDID     string                  `json:"holderDID"`
	IssuerDID     string                  `json:"issuerDID"`
}

// Passport
type PassportCreatedRequestDto struct {
	PassportType constant.PassportType `json:"passportType"`
	Nationality  string                `json:"nationality"`
	MRZ          string                `json:"mrz"`
	IssueDate    int64                 `json:"issueDate"`
	ExpiryDate   int64                 `json:"expiryDate"`
	HolderDID    string                `json:"holderDID"`
	IssuerDID    string                `json:"issuerDID"`
}

type PassportUpdatedRequestDto struct {
	PassportType constant.PassportType `json:"passportType,omitempty"`
	Nationality  string                `json:"nationality,omitempty"`
	MRZ          string                `json:"mrz,omitempty"`
	IssueDate    int64                 `json:"issueDate,omitempty"`
	ExpiryDate   int64                 `json:"expiryDate,omitempty"`
}

type PassportOptionRequestDto struct {
	Status string `json:"status"`
}

type PassportResponseDto struct {
	PublicID       string                  `json:"id"`
	PassportNumber string                  `json:"passportNumber"`
	Status         constant.DocumentStatus `json:"status"`
	PassportType   constant.PassportType   `json:"passportType"`
	Nationality    string                  `json:"nationality"`
	MRZ            string                  `json:"mrz"`
	IssueDate      int64                   `json:"issueDate"`
	ExpiryDate     int64                   `json:"expiryDate"`
	HolderDID      string                  `json:"holderDID"`
	IssuerDID      string                  `json:"issuerDID"`
}

// Convert
func CitizenIdentityToResponse(entity *document.CitizenIdentity) *CitizenIdentityResponseDto {
	return &CitizenIdentityResponseDto{
		PublicID:     entity.PublicID.String(),
		IDNumber:     entity.IDNumber,
		Status:       entity.Status,
		FirstName:    entity.FirstName,
		LastName:     entity.LastName,
		Gender:       entity.Gender,
		DateOfBirth:  entity.DateOfBirth,
		PlaceOfBirth: entity.PlaceOfBirth,
		IssueDate:    entity.IssueDate,
		ExpiryDate:   entity.ExpiryDate,
		HolderDID:    entity.HolderDID,
		IssuerDID:    entity.IssuerDID,
	}
}

func AcademicDegreeToResponse(entity *document.AcademicDegree) *AcademicDegreeResponseDto {
	return &AcademicDegreeResponseDto{
		PublicID:       entity.PublicID.String(),
		DegreeNumber:   entity.DegreeNumber,
		Status:         entity.Status,
		DegreeType:     entity.DegreeType,
		Major:          entity.Major,
		University:     entity.University,
		GraduateYear:   entity.GraduateYear,
		GPA:            entity.GPA,
		Classification: entity.Classification,
		IssueDate:      entity.IssueDate,
		HolderDID:      entity.HolderDID,
		IssuerDID:      entity.IssuerDID,
	}
}

func HealthInsuranceToResponse(entity *document.HealthInsurance) *HealthInsuranceResponseDto {
	return &HealthInsuranceResponseDto{
		PublicID:        entity.PublicID.String(),
		InsuranceNumber: entity.InsuranceNumber,
		Status:          entity.Status,
		InsuranceType:   entity.InsuranceType,
		Hospital:        entity.Hospital,
		StartDate:       entity.StartDate,
		ExpiryDate:      entity.ExpiryDate,
		HolderDID:       entity.HolderDID,
		IssuerDID:       entity.IssuerDID,
	}
}

func DriverLicenseToResponse(entity *document.DriverLicense) *DriverLicenseResponseDto {
	return &DriverLicenseResponseDto{
		PublicID:      entity.PublicID.String(),
		LicenseNumber: entity.LicenseNumber,
		Status:        entity.Status,
		Class:         entity.Class,
		Point:         entity.Point,
		IssueDate:     entity.IssueDate,
		ExpiryDate:    entity.ExpiryDate,
		HolderDID:     entity.HolderDID,
		IssuerDID:     entity.IssuerDID,
	}
}

func PassportToResponse(entity *document.Passport) *PassportResponseDto {
	return &PassportResponseDto{
		PublicID:       entity.PublicID.String(),
		PassportNumber: entity.PassportNumber,
		PassportType:   entity.PassportType,
		Status:         entity.Status,
		Nationality:    entity.Nationality,
		MRZ:            entity.MRZ,
		IssueDate:      entity.IssueDate,
		ExpiryDate:     entity.ExpiryDate,
		HolderDID:      entity.HolderDID,
		IssuerDID:      entity.IssuerDID,
	}
}
