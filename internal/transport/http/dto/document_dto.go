package dto

import (
	"be/internal/domain/document"
	"time"
)

type CitizenIdentityCreatedRequestDto struct {
	FirstName       string
	LastName        string
	Gender          string
	DateOfBirth     time.Time
	PlaceOfBirth    string
	PlaceOfResident string
	IssueDate       time.Time
	ExpiryDate      time.Time
	HolderDID       string
}

type CitizenIdentityUpdatedRequestDto struct {
	FirstName       string
	LastName        string
	Gender          string
	DateOfBirth     time.Time
	PlaceOfBirth    string
	PlaceOfResident string
	IssueDate       time.Time
	ExpiryDate      time.Time
}

type CitizenIdentityRevokedRequestDto struct {
	Status string
}

type CitizenIdentityResponseDto struct {
	PublicID        string
	IDNumber        string
	FirstName       string
	LastName        string
	Gender          string
	DateOfBirth     time.Time
	PlaceOfBirth    string
	PlaceOfResident string
	IssueDate       time.Time
	ExpiryDate      time.Time
	Status          string
	IssuerDID       string
}

type AcademicDegreeCreatedRequestDto struct {
	DegreeType     string
	Major          string
	University     string
	GraduateYear   uint
	GPA            float32
	Classification string
	IssueDate      time.Time
	HolderDID      string
}

type AcademicDegreeUpdatedRequestDto struct {
	DegreeType     string
	Major          string
	University     string
	GraduateYear   uint
	GPA            float32
	Classification string
	IssueDate      time.Time
}

type AcademicDegreeRevokedRequestDto struct {
	Status string
}

type AcademicDegreeResponseDto struct {
	PublicID       string
	DegreeNumber   string
	DegreeType     string
	Major          string
	University     string
	GraduateYear   uint
	GPA            float32
	Classification string
	Status         string
	IssueDate      time.Time
	IssuerDID      string
}

type HealthInsuranceCreatedRequestDto struct {
	InsuranceType string
	Hospital      string
	StartDate     time.Time
	ExpiryDate    time.Time
	HolderDID     string
}

type HealthInsuranceUpdatedRequestDto struct {
	InsuranceType string
	Hospital      string
	StartDate     time.Time
	ExpiryDate    time.Time
}

type HealthInsuranceRevokedRequestDto struct {
	Status string
}

type HealthInsuranceResponseDto struct {
	PublicID        string
	InsuranceNumber string
	InsuranceType   string
	Hospital        string
	Status          string
	StartDate       time.Time
	ExpiryDate      time.Time
	IssuerDID       string
}

type DriverLicenseCreatedRequestDto struct {
	Class      string
	IssueDate  time.Time
	ExpiryDate time.Time
	HolderDID  string
}

type DriverLicenseUpdatedRequestDto struct {
	Class          string
	Point          int
	PointResetDate time.Time
	IssueDate      time.Time
	ExpiryDate     time.Time
}

type DriverLicenseRevokedRequestDto struct {
	Status string
}

type DriverLicenseResponseDto struct {
	PublicID       string
	LicenseNumber  string
	Class          string
	Point          int
	PointResetDate time.Time
	Status         string
	IssueDate      time.Time
	ExpiryDate     time.Time
	IssuerDID      string
}

type PassportCreatedRequestDto struct {
	PassportType string
	Nationality  string
	MRZ          string
	IssueDate    time.Time
	ExpiryDate   time.Time
	HolderDID    string
}

type PassportUpdatedRequestDto struct {
	PassportType string
	Nationality  string
	MRZ          string
	IssueDate    time.Time
	ExpiryDate   time.Time
}

type PassportRevokedRequestDto struct {
	Status string
}

type PassportResponseDto struct {
	PublicID       string
	PassportType   string
	PasswordNumber string
	Nationality    string
	MRZ            string
	Status         string
	IssueDate      time.Time
	ExpiryDate     time.Time
	IssuerDID      string
}

func CitizenIdentityToResponse(entity *document.CitizenIdentity) *CitizenIdentityResponseDto {
	return &CitizenIdentityResponseDto{
		PublicID:     entity.PublicID.String(),
		IDNumber:     entity.IDNumber,
		FirstName:    entity.FirstName,
		LastName:     entity.LastName,
		Gender:       entity.Gender,
		DateOfBirth:  entity.DateOfBirth,
		PlaceOfBirth: entity.PlaceOfBirth,
		Status:       entity.Status,
		IssueDate:    entity.IssueDate,
		ExpiryDate:   entity.ExpiryDate,
		IssuerDID:    entity.IssuerDID,
	}
}

func AcademicDegreeToResponse(entity *document.AcademicDegree) *AcademicDegreeResponseDto {
	return &AcademicDegreeResponseDto{
		PublicID:       entity.PublicID.String(),
		DegreeNumber:   entity.DegreeNumber,
		DegreeType:     entity.DegreeType,
		Major:          entity.Major,
		University:     entity.University,
		GraduateYear:   entity.GraduateYear,
		GPA:            entity.GPA,
		Classification: entity.Classification,
		Status:         entity.Status,
		IssueDate:      entity.IssueDate,
		IssuerDID:      entity.IssuerDID,
	}
}

func HealthInsuranceToResponse(entity *document.HealthInsurance) *HealthInsuranceResponseDto {
	return &HealthInsuranceResponseDto{
		PublicID:        entity.PublicID.String(),
		InsuranceNumber: entity.InsuranceNumber,
		InsuranceType:   entity.InsuranceType,
		Hospital:        entity.Hospital,
		Status:          entity.Status,
		StartDate:       entity.StartDate,
		ExpiryDate:      entity.ExpiryDate,
		IssuerDID:       entity.IssuerDID,
	}
}

func DriverLicenseToResponse(entity *document.DriverLicense) *DriverLicenseResponseDto {
	return &DriverLicenseResponseDto{
		PublicID:       entity.PublicID.String(),
		LicenseNumber:  entity.LicenseNumber,
		Class:          entity.Class,
		Point:          entity.Point,
		PointResetDate: entity.PointResetDate,
		Status:         entity.Status,
		IssueDate:      entity.IssueDate,
		ExpiryDate:     entity.ExpiryDate,
		IssuerDID:      entity.IssuerDID,
	}
}

func PassportToResponse(entity *document.Passport) *PassportResponseDto {
	return &PassportResponseDto{
		PublicID:       entity.PublicID.String(),
		PasswordNumber: entity.PassportNumber,
		Nationality:    entity.Nationality,
		MRZ:            entity.MRZ,
		Status:         entity.Status,
		IssueDate:      entity.IssueDate,
		ExpiryDate:     entity.ExpiryDate,
		IssuerDID:      entity.IssuerDID,
	}
}
