package service

import (
	"be/config"
	"be/internal/domain/document"
	"be/internal/shared/constant"
	"be/internal/shared/utils"
	"be/internal/transport/http/dto"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IDocumentService interface {
	CreateCitizenIdentity(ctx context.Context, request *dto.CitizenIdentityCreatedRequestDto) (*dto.CitizenIdentityResponseDto, error)
	UpdateCitizenIdentity(ctx context.Context, id string, request *dto.CitizenIdentityUpdatedRequestDto) (*dto.CitizenIdentityResponseDto, error)
	RevokeCitizenIdentity(ctx context.Context, id string, request *dto.CitizenIdentityRevokedRequestDto) error
	GetCitizenIdentityById(ctx context.Context, id string) (*dto.CitizenIdentityResponseDto, error)
	GetCitizenIdentityByIdNumber(ctx context.Context, idNumber string) (*dto.CitizenIdentityResponseDto, error)
	GetCitizenIdentities(ctx context.Context) ([]*dto.CitizenIdentityResponseDto, error)

	CreateAcademicDegree(ctx context.Context, request *dto.AcademicDegreeCreatedRequestDto) (*dto.AcademicDegreeResponseDto, error)
	UpdateAcademicDegree(ctx context.Context, id string, request *dto.AcademicDegreeUpdatedRequestDto) (*dto.AcademicDegreeResponseDto, error)
	RevokeAcademicDegree(ctx context.Context, id string, request *dto.AcademicDegreeRevokedRequestDto) error
	GetAcademicDegreeById(ctx context.Context, id string) (*dto.AcademicDegreeResponseDto, error)
	GetAcademicDegreeByDegreeNumber(ctx context.Context, degreeNumber string) (*dto.AcademicDegreeResponseDto, error)
	GetAcademicDegrees(ctx context.Context) ([]*dto.AcademicDegreeResponseDto, error)

	CreateHealthInsurance(ctx context.Context, request *dto.HealthInsuranceCreatedRequestDto) (*dto.HealthInsuranceResponseDto, error)
	UpdateHealthInsurance(ctx context.Context, id string, request *dto.HealthInsuranceUpdatedRequestDto) (*dto.HealthInsuranceResponseDto, error)
	RevokeHealthInsurance(ctx context.Context, id string, request *dto.HealthInsuranceRevokedRequestDto) error
	GetHealthInsuranceById(ctx context.Context, id string) (*dto.HealthInsuranceResponseDto, error)
	GetHealthInsuranceByInsuranceNumber(ctx context.Context, insuranceNumber string) (*dto.HealthInsuranceResponseDto, error)
	GetHealthInsurances(ctx context.Context) ([]*dto.HealthInsuranceResponseDto, error)

	CreateDriverLicense(ctx context.Context, request *dto.DriverLicenseCreatedRequestDto) (*dto.DriverLicenseResponseDto, error)
	UpdateDriverLicense(ctx context.Context, id string, request *dto.DriverLicenseUpdatedRequestDto) (*dto.DriverLicenseResponseDto, error)
	RevokeDriverLicense(ctx context.Context, id string, request *dto.DriverLicenseRevokedRequestDto) error
	GetDriverLicenseById(ctx context.Context, id string) (*dto.DriverLicenseResponseDto, error)
	GetDriverLicenseByLicenseNumber(ctx context.Context, licenseNumber string) (*dto.DriverLicenseResponseDto, error)
	GetDriverLicenses(ctx context.Context) ([]*dto.DriverLicenseResponseDto, error)

	CreatePassport(ctx context.Context, request *dto.PassportCreatedRequestDto) (*dto.PassportResponseDto, error)
	UpdatePassport(ctx context.Context, id string, request *dto.PassportUpdatedRequestDto) (*dto.PassportResponseDto, error)
	RevokePassport(ctx context.Context, id string, request *dto.PassportRevokedRequestDto) error
	GetPassportById(ctx context.Context, id string) (*dto.PassportResponseDto, error)
	GetPassportByPassportNumber(ctx context.Context, passportNumber string) (*dto.PassportResponseDto, error)
	GetPassports(ctx context.Context) ([]*dto.PassportResponseDto, error)
}

type DocumentService struct {
	config              *config.Config
	citizenIdentityRepo document.ICitizenIdentityRepository
	academicDegreeRepo  document.IAcademicDegreeRepository
	healthInsuranceRepo document.IHealthInsuranceRepository
	driverLicenseRepo   document.IDriverLicenseRepository
	passportRepo        document.IPassportRepository
}

func NewDocumentService(
	config *config.Config,
	citizenIdentityRepo document.ICitizenIdentityRepository,
	academicDegreeRepo document.IAcademicDegreeRepository,
	healthInsuranceRepo document.IHealthInsuranceRepository,
	driverInsuranceRepo document.IDriverLicenseRepository,
	passportRepo document.IPassportRepository,
) IDocumentService {
	return &DocumentService{
		config:              config,
		citizenIdentityRepo: citizenIdentityRepo,
		academicDegreeRepo:  academicDegreeRepo,
		healthInsuranceRepo: healthInsuranceRepo,
		driverLicenseRepo:   driverInsuranceRepo,
		passportRepo:        passportRepo,
	}
}

func (s *DocumentService) CreateCitizenIdentity(ctx context.Context, request *dto.CitizenIdentityCreatedRequestDto) (*dto.CitizenIdentityResponseDto, error) {
	idNumber, _ := utils.GetIdNumber(request.DateOfBirth, request.PlaceOfBirth, request.Gender)
	citizenCreated, err := s.citizenIdentityRepo.CreateCitizenIdentity(ctx, &document.CitizenIdentity{
		PublicID:     uuid.New(),
		IDNumber:     idNumber,
		FirstName:    request.FirstName,
		LastName:     request.LastName,
		Gender:       request.Gender,
		DateOfBirth:  request.DateOfBirth,
		PlaceOfBirth: request.PlaceOfBirth,
		Status:       constant.ActiveStatus,
		IssueDate:    request.IssueDate,
		ExpiryDate:   request.ExpiryDate,
		IssuerDID:    "",
		HolderDID:    request.HolderDID,
	})

	if err != nil {
		return nil, &constant.InternalServer
	}

	return dto.CitizenIdentityToResponse(citizenCreated), nil

}

func (s *DocumentService) UpdateCitizenIdentity(ctx context.Context, id string, request *dto.CitizenIdentityUpdatedRequestDto) (*dto.CitizenIdentityResponseDto, error) {
	citizen, err := s.citizenIdentityRepo.FindCitizenIdentityByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.CitizenIdentityNotFound
		}
		return nil, &constant.InternalServer
	}

	citizen.FirstName = request.FirstName
	citizen.LastName = request.LastName
	citizen.Gender = request.Gender
	citizen.DateOfBirth = request.DateOfBirth
	citizen.PlaceOfBirth = request.PlaceOfBirth
	citizen.IssueDate = request.IssueDate
	citizen.ExpiryDate = request.ExpiryDate

	citizenUpdated, err := s.citizenIdentityRepo.SaveCitizenIdentity(ctx, citizen)

	if err != nil {
		return nil, &constant.InternalServer
	}

	return dto.CitizenIdentityToResponse(citizenUpdated), nil

}

func (s *DocumentService) RevokeCitizenIdentity(ctx context.Context, id string, request *dto.CitizenIdentityRevokedRequestDto) error {
	citizenIdentity, err := s.citizenIdentityRepo.FindCitizenIdentityByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &constant.CitizenIdentityNotFound
		}
		return &constant.InternalServer
	}
	changes := map[string]interface{}{"status": request.Status, "revoked_at": time.Now()}
	return s.citizenIdentityRepo.UpdateCitizenIdentity(ctx, citizenIdentity, changes)
}

func (s *DocumentService) GetCitizenIdentityById(ctx context.Context, id string) (*dto.CitizenIdentityResponseDto, error) {
	citizen, err := s.citizenIdentityRepo.FindCitizenIdentityByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.UserNotFound
		}
		return nil, &constant.InternalServer
	}
	return dto.CitizenIdentityToResponse(citizen), nil
}
func (s *DocumentService) GetCitizenIdentityByIdNumber(ctx context.Context, idNumber string) (*dto.CitizenIdentityResponseDto, error) {
	citizen, err := s.citizenIdentityRepo.FindCitizenIdentityByIdNumber(ctx, idNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.UserNotFound
		}
		return nil, &constant.InternalServer
	}
	return dto.CitizenIdentityToResponse(citizen), nil

}

func (s *DocumentService) GetCitizenIdentities(ctx context.Context) ([]*dto.CitizenIdentityResponseDto, error) {
	citizens, err := s.citizenIdentityRepo.FindAllCitizenIdentities(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.UserNotFound
		}
		return nil, &constant.InternalServer
	}

	var resps []*dto.CitizenIdentityResponseDto
	for _, c := range citizens {
		resps = append(resps, dto.CitizenIdentityToResponse(c))
	}
	return resps, nil
}

func (s *DocumentService) CreateAcademicDegree(ctx context.Context, request *dto.AcademicDegreeCreatedRequestDto) (*dto.AcademicDegreeResponseDto, error) {
	citizenIdentity, err := s.citizenIdentityRepo.FindCitizenIdentityByHolderDID(ctx, request.HolderDID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.AcademicDegreeNotFound
		}
		return nil, err
	}

	degreeNumber, _ := utils.GetDegreeNumber(request.University, request.Major, request.GraduateYear)
	academicDegreeCreated, err := s.academicDegreeRepo.CreateAcademicDegree(ctx, &document.AcademicDegree{
		PublicID:       uuid.New(),
		CID:            citizenIdentity.ID,
		DegreeNumber:   degreeNumber,
		DegreeType:     request.DegreeType,
		Major:          request.Major,
		University:     request.University,
		GraduateYear:   request.GraduateYear,
		GPA:            request.GPA,
		Classification: request.Classification,
		Status:         constant.ActiveStatus,
		IssueDate:      request.IssueDate,
		IssuerDID:      "",
	})

	if err != nil {
		return nil, &constant.InternalServer
	}

	return dto.AcademicDegreeToResponse(academicDegreeCreated), nil
}

func (s *DocumentService) UpdateAcademicDegree(ctx context.Context, id string, request *dto.AcademicDegreeUpdatedRequestDto) (*dto.AcademicDegreeResponseDto, error) {
	academicDegree, err := s.academicDegreeRepo.FindAcademicDegreeByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.AcademicDegreeNotFound
		}
		return nil, &constant.InternalServer
	}

	academicDegree.DegreeType = request.DegreeType
	academicDegree.Major = request.Major
	academicDegree.University = request.University
	academicDegree.GraduateYear = request.GraduateYear
	academicDegree.GPA = request.GPA
	academicDegree.Classification = request.Classification
	academicDegree.IssueDate = request.IssueDate

	academicDegreeUpdated, err := s.academicDegreeRepo.SaveAcademicDegree(ctx, academicDegree)

	if err != nil {
		return nil, &constant.InternalServer
	}

	return dto.AcademicDegreeToResponse(academicDegreeUpdated), nil
}

func (s *DocumentService) RevokeAcademicDegree(ctx context.Context, id string, request *dto.AcademicDegreeRevokedRequestDto) error {
	academicDegree, err := s.academicDegreeRepo.FindAcademicDegreeByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &constant.AcademicDegreeNotFound
		}
		return &constant.InternalServer
	}
	changes := map[string]interface{}{"status": request.Status, "revoked_at": time.Now()}
	return s.academicDegreeRepo.UpdateAcademicDegree(ctx, academicDegree, changes)
}

func (s *DocumentService) GetAcademicDegreeById(ctx context.Context, id string) (*dto.AcademicDegreeResponseDto, error) {
	academicDegree, err := s.academicDegreeRepo.FindAcademicDegreeByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.AcademicDegreeNotFound
		}
		return nil, &constant.InternalServer
	}

	return dto.AcademicDegreeToResponse(academicDegree), nil
}
func (s *DocumentService) GetAcademicDegreeByDegreeNumber(ctx context.Context, degreeNumber string) (*dto.AcademicDegreeResponseDto, error) {
	academicDegree, err := s.academicDegreeRepo.FindAcademicDegreeByDegreeNumber(ctx, degreeNumber)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.AcademicDegreeNotFound
		}
		return nil, &constant.InternalServer
	}

	return dto.AcademicDegreeToResponse(academicDegree), nil
}

func (s *DocumentService) GetAcademicDegrees(ctx context.Context) ([]*dto.AcademicDegreeResponseDto, error) {
	academicDegrees, err := s.academicDegreeRepo.FindAllAcademicDegrees(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.AcademicDegreeNotFound
		}
		return nil, &constant.InternalServer
	}

	var resps []*dto.AcademicDegreeResponseDto
	for _, a := range academicDegrees {
		resps = append(resps, dto.AcademicDegreeToResponse(a))
	}

	return resps, nil
}

func (s *DocumentService) CreateHealthInsurance(ctx context.Context, request *dto.HealthInsuranceCreatedRequestDto) (*dto.HealthInsuranceResponseDto, error) {
	citizenIdentity, err := s.citizenIdentityRepo.FindCitizenIdentityByHolderDID(ctx, request.HolderDID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.HealthInsuranceNotFound
		}
		return nil, err
	}

	insuranceNumber, _ := utils.GetInsuranceNumber(request.InsuranceType)
	healthInsuranceCreated, err := s.healthInsuranceRepo.CreateHealthInsurance(ctx, &document.HealthInsurance{
		PublicID:        uuid.New(),
		CID:             citizenIdentity.ID,
		InsuranceNumber: insuranceNumber,
		InsuranceType:   request.InsuranceType,
		Hospital:        request.Hospital,
		Status:          constant.ActiveStatus,
		StartDate:       request.StartDate,
		ExpiryDate:      request.ExpiryDate,
		IssuerDID:       "",
	})

	if err != nil {
		return nil, &constant.InternalServer
	}

	return dto.HealthInsuranceToResponse(healthInsuranceCreated), nil

}
func (s *DocumentService) UpdateHealthInsurance(ctx context.Context, id string, request *dto.HealthInsuranceUpdatedRequestDto) (*dto.HealthInsuranceResponseDto, error) {
	healthInsurance, err := s.healthInsuranceRepo.FindHealthInsuranceByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.HealthInsuranceNotFound
		}
		return nil, &constant.InternalServer
	}

	healthInsurance.InsuranceType = request.InsuranceType
	healthInsurance.Hospital = request.Hospital
	healthInsurance.StartDate = request.StartDate
	healthInsurance.ExpiryDate = request.ExpiryDate

	healthInsuranceUpdated, err := s.healthInsuranceRepo.SaveHealthInsurance(ctx, healthInsurance)

	if err != nil {
		return nil, &constant.InternalServer
	}

	return dto.HealthInsuranceToResponse(healthInsuranceUpdated), nil
}

func (s *DocumentService) RevokeHealthInsurance(ctx context.Context, id string, request *dto.HealthInsuranceRevokedRequestDto) error {
	healthInsurance, err := s.healthInsuranceRepo.FindHealthInsuranceByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &constant.HealthInsuranceNotFound
		}
		return &constant.InternalServer
	}
	changes := map[string]interface{}{"status": request.Status, "revoked_at": time.Now()}
	return s.healthInsuranceRepo.UpdateHealthInsurance(ctx, healthInsurance, changes)
}

func (s *DocumentService) GetHealthInsuranceById(ctx context.Context, id string) (*dto.HealthInsuranceResponseDto, error) {
	healthInsurance, err := s.healthInsuranceRepo.FindHealthInsuranceByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.HealthInsuranceNotFound
		}
		return nil, &constant.InternalServer
	}

	return dto.HealthInsuranceToResponse(healthInsurance), nil
}

func (s *DocumentService) GetHealthInsuranceByInsuranceNumber(ctx context.Context, insuranceNumber string) (*dto.HealthInsuranceResponseDto, error) {
	healthInsurance, err := s.healthInsuranceRepo.FindHealthInsuranceByInsuranceNumber(ctx, insuranceNumber)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.HealthInsuranceNotFound
		}
		return nil, &constant.InternalServer
	}

	return dto.HealthInsuranceToResponse(healthInsurance), nil
}

func (s *DocumentService) GetHealthInsurances(ctx context.Context) ([]*dto.HealthInsuranceResponseDto, error) {
	healthInsurances, err := s.healthInsuranceRepo.FindAllHealthInsurances(ctx)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, &constant.HealthInsuranceNotFound
	}

	var resps []*dto.HealthInsuranceResponseDto
	for _, h := range healthInsurances {
		resps = append(resps, dto.HealthInsuranceToResponse(h))
	}
	return resps, nil
}

func (s *DocumentService) CreateDriverLicense(ctx context.Context, request *dto.DriverLicenseCreatedRequestDto) (*dto.DriverLicenseResponseDto, error) {
	citizenIdentity, err := s.citizenIdentityRepo.FindCitizenIdentityByHolderDID(ctx, request.HolderDID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.DriverLicenseNotFound
		}
		return nil, err
	}

	licenseNumber, _ := utils.GetLicenseNumber(request.Class)
	driverLicenseCreated, err := s.driverLicenseRepo.CreateDriverLicense(ctx, &document.DriverLicense{
		PublicID:       uuid.New(),
		CID:            citizenIdentity.ID,
		LicenseNumber:  licenseNumber,
		Class:          request.Class,
		Point:          0,
		PointResetDate: request.IssueDate,
		Status:         constant.ActiveStatus,
		IssueDate:      request.IssueDate,
		ExpiryDate:     request.ExpiryDate,
		IssuerDID:      "",
	})

	if err != nil {
		return nil, &constant.InternalServer
	}

	return dto.DriverLicenseToResponse(driverLicenseCreated), nil
}

func (s *DocumentService) UpdateDriverLicense(ctx context.Context, id string, request *dto.DriverLicenseUpdatedRequestDto) (*dto.DriverLicenseResponseDto, error) {
	driverLicense, err := s.driverLicenseRepo.FindDriverLicenseByPublicId(ctx, id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.DriverLicenseNotFound
		}
		return nil, &constant.InternalServer
	}
	driverLicense.Class = request.Class
	driverLicense.Point = request.Point
	driverLicense.PointResetDate = request.PointResetDate
	driverLicense.IssueDate = request.IssueDate
	driverLicense.ExpiryDate = request.ExpiryDate

	driverLicenseUpdated, err := s.driverLicenseRepo.SaveDriverLicense(ctx, driverLicense)

	if err != nil {
		return nil, &constant.InternalServer
	}

	return dto.DriverLicenseToResponse(driverLicenseUpdated), nil
}

func (s *DocumentService) RevokeDriverLicense(ctx context.Context, id string, request *dto.DriverLicenseRevokedRequestDto) error {
	driverLicense, err := s.driverLicenseRepo.FindDriverLicenseByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &constant.DriverLicenseNotFound
		}
		return &constant.InternalServer
	}
	changes := map[string]interface{}{"status": request.Status, "revoked_at": time.Now()}
	return s.driverLicenseRepo.UpdateDriverLicense(ctx, driverLicense, changes)
}

func (s *DocumentService) GetDriverLicenseById(ctx context.Context, id string) (*dto.DriverLicenseResponseDto, error) {
	driverLicense, err := s.driverLicenseRepo.FindDriverLicenseByPublicId(ctx, id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.DriverLicenseNotFound
		}
		return nil, &constant.InternalServer
	}

	return dto.DriverLicenseToResponse(driverLicense), nil
}

func (s *DocumentService) GetDriverLicenseByLicenseNumber(ctx context.Context, licenseNumber string) (*dto.DriverLicenseResponseDto, error) {
	driverLicense, err := s.driverLicenseRepo.FindDriverLicenseByLicenseId(ctx, licenseNumber)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.DriverLicenseNotFound
		}
		return nil, &constant.InternalServer
	}

	return dto.DriverLicenseToResponse(driverLicense), nil
}

func (s *DocumentService) GetDriverLicenses(ctx context.Context) ([]*dto.DriverLicenseResponseDto, error) {
	driverLicenses, err := s.driverLicenseRepo.FindAllDriverLicenses(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.DriverLicenseNotFound
		}
		return nil, &constant.InternalServer
	}

	var resps []*dto.DriverLicenseResponseDto
	for _, d := range driverLicenses {
		resps = append(resps, dto.DriverLicenseToResponse(d))
	}
	return resps, nil
}

func (s *DocumentService) CreatePassport(ctx context.Context, request *dto.PassportCreatedRequestDto) (*dto.PassportResponseDto, error) {
	citizenIdentity, err := s.citizenIdentityRepo.FindCitizenIdentityByHolderDID(ctx, request.HolderDID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.PassportNotFound
		}
		return nil, err
	}

	passportNumber, _ := utils.GetPassportNumber(request.PassportType)
	passportCreated, err := s.passportRepo.CreatePassport(ctx, &document.Passport{
		PublicID:       uuid.New(),
		CID:            citizenIdentity.ID,
		PassportNumber: passportNumber,
		PassportType:   request.PassportType,
		Nationality:    request.Nationality,
		MRZ:            request.MRZ,
		Status:         constant.ActiveStatus,
		IssueDate:      request.IssueDate,
		ExpiryDate:     request.ExpiryDate,
		IssuerDID:      "",
	})

	if err != nil {
		return nil, &constant.InternalServer
	}
	return dto.PassportToResponse(passportCreated), nil
}
func (s *DocumentService) UpdatePassport(ctx context.Context, id string, request *dto.PassportUpdatedRequestDto) (*dto.PassportResponseDto, error) {
	passport, err := s.passportRepo.FindPassportByPublicId(ctx, id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.PassportNotFound
		}
		return nil, &constant.InternalServer
	}

	passport.Nationality = request.Nationality
	passport.MRZ = request.MRZ
	passport.IssueDate = request.IssueDate
	passport.ExpiryDate = request.ExpiryDate

	passportUpdated, err := s.passportRepo.SavePassport(ctx, passport)

	if err != nil {
		return nil, &constant.InternalServer
	}

	return dto.PassportToResponse(passportUpdated), nil
}

func (s *DocumentService) RevokePassport(ctx context.Context, id string, request *dto.PassportRevokedRequestDto) error {
	passport, err := s.passportRepo.FindPassportByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &constant.PassportNotFound
		}
		return &constant.InternalServer
	}
	changes := map[string]interface{}{"status": request.Status, "revoked_at": time.Now()}
	return s.passportRepo.UpdatePassport(ctx, passport, changes)
}

func (s *DocumentService) GetPassportById(ctx context.Context, id string) (*dto.PassportResponseDto, error) {
	passport, err := s.passportRepo.FindPassportByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.PassportNotFound
		}
		return nil, &constant.InternalServer
	}

	return dto.PassportToResponse(passport), nil
}

func (s *DocumentService) GetPassportByPassportNumber(ctx context.Context, passportNumber string) (*dto.PassportResponseDto, error) {
	passport, err := s.passportRepo.FindPassportByPassportNumber(ctx, passportNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.PassportNotFound
		}
		return nil, &constant.InternalServer
	}

	return dto.PassportToResponse(passport), nil
}

func (s *DocumentService) GetPassports(ctx context.Context) ([]*dto.PassportResponseDto, error) {
	passports, err := s.passportRepo.FindAllPassports(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.PassportNotFound
		}
		return nil, &constant.InternalServer
	}

	var resps []*dto.PassportResponseDto
	for _, p := range passports {
		resps = append(resps, dto.PassportToResponse(p))
	}
	return resps, nil
}
