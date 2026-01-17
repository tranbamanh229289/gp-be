package document

import (
	"context"
)

type ICitizenIdentityRepository interface {
	FindCitizenIdentityByPublicId(ctx context.Context, publicId string) (*CitizenIdentity, error)
	FindCitizenIdentityByIdNumber(ctx context.Context, idNumber string) (*CitizenIdentity, error)
	FindCitizenIdentityByHolderDID(ctx context.Context, holderDID string) (*CitizenIdentity, error)
	FindAllCitizenIdentities(ctx context.Context) ([]*CitizenIdentity, error)
	CreateCitizenIdentity(ctx context.Context, entity *CitizenIdentity) (*CitizenIdentity, error)
	SaveCitizenIdentity(ctx context.Context, entity *CitizenIdentity) (*CitizenIdentity, error)
	UpdateCitizenIdentity(ctx context.Context, entity *CitizenIdentity, changes map[string]interface{}) error
}
type IAcademicDegreeRepository interface {
	FindAcademicDegreeByPublicId(ctx context.Context, publicId string) (*AcademicDegree, error)
	FindAcademicDegreeByDegreeNumber(ctx context.Context, degreeNumber string) (*AcademicDegree, error)
	FindAcademicDegreeByHolderDID(ctx context.Context, holderDID string) (*AcademicDegree, error)
	FindAllAcademicDegrees(ctx context.Context) ([]*AcademicDegree, error)
	CreateAcademicDegree(ctx context.Context, entity *AcademicDegree) (*AcademicDegree, error)
	SaveAcademicDegree(ctx context.Context, entity *AcademicDegree) (*AcademicDegree, error)
	UpdateAcademicDegree(ctx context.Context, entity *AcademicDegree, changes map[string]interface{}) error
}

type IHealthInsuranceRepository interface {
	FindHealthInsuranceByPublicId(ctx context.Context, publicId string) (*HealthInsurance, error)
	FindHealthInsuranceByInsuranceNumber(ctx context.Context, insuranceNumber string) (*HealthInsurance, error)
	FindHealthInsuranceByHolderDID(ctx context.Context, holderDID string) (*HealthInsurance, error)
	FindAllHealthInsurances(ctx context.Context) ([]*HealthInsurance, error)
	CreateHealthInsurance(ctx context.Context, entity *HealthInsurance) (*HealthInsurance, error)
	SaveHealthInsurance(ctx context.Context, entity *HealthInsurance) (*HealthInsurance, error)
	UpdateHealthInsurance(ctx context.Context, entity *HealthInsurance, changes map[string]interface{}) error
}

type IDriverLicenseRepository interface {
	FindDriverLicenseByPublicId(ctx context.Context, publicId string) (*DriverLicense, error)
	FindDriverLicenseByLicenseId(ctx context.Context, licenseNumber string) (*DriverLicense, error)
	FindDriverLicenseByHolderDID(ctx context.Context, holderDID string) (*DriverLicense, error)
	FindAllDriverLicenses(ctx context.Context) ([]*DriverLicense, error)
	CreateDriverLicense(ctx context.Context, entity *DriverLicense) (*DriverLicense, error)
	SaveDriverLicense(ctx context.Context, entity *DriverLicense) (*DriverLicense, error)
	UpdateDriverLicense(ctx context.Context, entity *DriverLicense, changes map[string]interface{}) error
}

type IPassportRepository interface {
	FindPassportByPublicId(ctx context.Context, publicId string) (*Passport, error)
	FindPassportByPassportNumber(ctx context.Context, passportNumber string) (*Passport, error)
	FindPassportByHolderDID(ctx context.Context, holderDID string) (*Passport, error)
	FindAllPassports(ctx context.Context) ([]*Passport, error)
	CreatePassport(ctx context.Context, entity *Passport) (*Passport, error)
	SavePassport(ctx context.Context, entity *Passport) (*Passport, error)
	UpdatePassport(ctx context.Context, entity *Passport, changes map[string]interface{}) error
}
