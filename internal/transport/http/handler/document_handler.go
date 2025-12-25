package handler

import (
	"be/internal/service"
	response "be/internal/shared/helper"
	"be/internal/transport/http/dto"

	"github.com/gin-gonic/gin"
)

type DocumentHandler struct {
	documentService service.IDocumentService
}

func NewDocumentHandler(cs service.IDocumentService) *DocumentHandler {
	return &DocumentHandler{
		documentService: cs,
	}
}

func (h *DocumentHandler) CreateCitizenIdentity(c *gin.Context) {
	var citizenIdentityRequest dto.CitizenIdentityCreatedRequestDto
	if err := c.ShouldBindJSON(citizenIdentityRequest); err != nil {
		response.RespondError(c, err)
	}
	citizenIdentityResponse, err := h.documentService.CreateCitizenIdentity(c.Request.Context(), &citizenIdentityRequest)

	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, citizenIdentityResponse)
}

func (h *DocumentHandler) UpdateCitizenIdentity(c *gin.Context) {
	id := c.Param("id")
	var citizenIdentityRequest dto.CitizenIdentityUpdatedRequestDto
	if err := c.ShouldBindJSON(&citizenIdentityRequest); err != nil {
		response.RespondError(c, err)
	}

	citizenIdentityResponse, err := h.documentService.UpdateCitizenIdentity(c.Request.Context(), id, &citizenIdentityRequest)

	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, citizenIdentityResponse)
}

func (h *DocumentHandler) RevokeCitizenIdentity(c *gin.Context) {
	id := c.Param("id")
	var citizenIdentityRevokedRequest dto.CitizenIdentityRevokedRequestDto

	if err := c.ShouldBindJSON(&citizenIdentityRevokedRequest); err != nil {
		response.RespondError(c, err)
	}

	err := h.documentService.RevokeCitizenIdentity(c.Request.Context(), id, &citizenIdentityRevokedRequest)
	if err != nil {
		response.RespondError(c, err)
	}

	response.RespondSuccess(c, "")
}

func (h *DocumentHandler) GetCitizenIdentity(c *gin.Context) {
	id := c.Param("id")
	citizenIdentity, err := h.documentService.GetCitizenIdentityById(c.Request.Context(), id)
	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, citizenIdentity)
}

func (h *DocumentHandler) GetCitizenIdentities(c *gin.Context) {
	citizenIdentities, err := h.documentService.GetCitizenIdentities(c.Request.Context())
	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, citizenIdentities)
}

func (h *DocumentHandler) CreateAcademicDegree(c *gin.Context) {
	var academicDegreeRequest dto.AcademicDegreeCreatedRequestDto
	if err := c.ShouldBindJSON(&academicDegreeRequest); err != nil {
		response.RespondError(c, err)
	}
	academicDegreeResponse, err := h.documentService.CreateAcademicDegree(c.Request.Context(), &academicDegreeRequest)

	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, academicDegreeResponse)
}

func (h *DocumentHandler) UpdateAcademicDegree(c *gin.Context) {
	id := c.Param("id")
	var academicDegreeRequest dto.AcademicDegreeUpdatedRequestDto
	if err := c.ShouldBindJSON(&academicDegreeRequest); err != nil {
		response.RespondError(c, err)
	}

	academicDegreeResponse, err := h.documentService.UpdateAcademicDegree(c.Request.Context(), id, &academicDegreeRequest)

	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, academicDegreeResponse)
}

func (h *DocumentHandler) RevokeAcademicDegree(c *gin.Context) {
	id := c.Param("id")
	var academicDegreeRevokedRequest dto.AcademicDegreeRevokedRequestDto

	if err := c.ShouldBindJSON(&academicDegreeRevokedRequest); err != nil {
		response.RespondError(c, err)
	}

	err := h.documentService.RevokeAcademicDegree(c.Request.Context(), id, &academicDegreeRevokedRequest)
	if err != nil {
		response.RespondError(c, err)
	}

	response.RespondSuccess(c, "")
}

func (h *DocumentHandler) GetAcademicDegree(c *gin.Context) {
	id := c.Param("id")
	academicDegree, err := h.documentService.GetAcademicDegreeById(c.Request.Context(), id)
	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, academicDegree)
}

func (h *DocumentHandler) GetAcademicDegrees(c *gin.Context) {
	academicDegrees, err := h.documentService.GetAcademicDegrees(c.Request.Context())
	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, academicDegrees)
}

func (h *DocumentHandler) CreateHealthInsurance(c *gin.Context) {
	var healthInsuranceRequest dto.HealthInsuranceCreatedRequestDto
	if err := c.ShouldBindJSON(&healthInsuranceRequest); err != nil {
		response.RespondError(c, err)
	}
	academyInsuranceResponse, err := h.documentService.CreateHealthInsurance(c.Request.Context(), &healthInsuranceRequest)

	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, academyInsuranceResponse)
}

func (h *DocumentHandler) UpdateHealthInsurance(c *gin.Context) {
	id := c.Param("id")
	var healthInsuranceRequest dto.HealthInsuranceUpdatedRequestDto
	if err := c.ShouldBindJSON(&healthInsuranceRequest); err != nil {
		response.RespondError(c, err)
	}

	healthInsuranceResponse, err := h.documentService.UpdateHealthInsurance(c.Request.Context(), id, &healthInsuranceRequest)

	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, healthInsuranceResponse)
}

func (h *DocumentHandler) RevokeHealthInsurance(c *gin.Context) {
	id := c.Param("id")
	var healthInsuranceRevokedRequest dto.HealthInsuranceRevokedRequestDto

	if err := c.ShouldBindJSON(&healthInsuranceRevokedRequest); err != nil {
		response.RespondError(c, err)
	}

	err := h.documentService.RevokeHealthInsurance(c.Request.Context(), id, &healthInsuranceRevokedRequest)
	if err != nil {
		response.RespondError(c, err)
	}

	response.RespondSuccess(c, "")
}

func (h *DocumentHandler) GetHealthInsurance(c *gin.Context) {
	id := c.Param("id")
	healthInsurance, err := h.documentService.GetHealthInsuranceById(c.Request.Context(), id)
	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, healthInsurance)
}

func (h *DocumentHandler) GetHealthInsurances(c *gin.Context) {
	healthInsurances, err := h.documentService.GetHealthInsurances(c.Request.Context())
	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, healthInsurances)
}

func (h *DocumentHandler) CreateDriverLicense(c *gin.Context) {
	var driverLicenseRequest dto.DriverLicenseCreatedRequestDto
	if err := c.ShouldBindJSON(&driverLicenseRequest); err != nil {
		response.RespondError(c, err)
	}
	driverLicenseResponse, err := h.documentService.CreateDriverLicense(c.Request.Context(), &driverLicenseRequest)

	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, driverLicenseResponse)
}

func (h *DocumentHandler) UpdateDriverLicense(c *gin.Context) {
	id := c.Param("id")
	var driverLicenseRequest dto.DriverLicenseUpdatedRequestDto
	if err := c.ShouldBindJSON(&driverLicenseRequest); err != nil {
		response.RespondError(c, err)
	}

	driverLicenseResponse, err := h.documentService.UpdateDriverLicense(c.Request.Context(), id, &driverLicenseRequest)

	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, driverLicenseResponse)
}

func (h *DocumentHandler) RevokeDriverLicense(c *gin.Context) {
	id := c.Param("id")
	var driverLicenseRevokedRequest dto.DriverLicenseRevokedRequestDto

	if err := c.ShouldBindJSON(&driverLicenseRevokedRequest); err != nil {
		response.RespondError(c, err)
	}

	err := h.documentService.RevokeDriverLicense(c.Request.Context(), id, &driverLicenseRevokedRequest)
	if err != nil {
		response.RespondError(c, err)
	}

	response.RespondSuccess(c, "")
}

func (h *DocumentHandler) GetDriverLicense(c *gin.Context) {
	id := c.Param("id")
	driverLicense, err := h.documentService.GetDriverLicenseById(c.Request.Context(), id)
	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, driverLicense)
}

func (h *DocumentHandler) GetDriverLicenses(c *gin.Context) {
	driverLicenses, err := h.documentService.GetDriverLicenses(c.Request.Context())
	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, driverLicenses)
}

func (h *DocumentHandler) CreatePassport(c *gin.Context) {
	var passportRequest dto.PassportCreatedRequestDto
	if err := c.ShouldBindJSON(&passportRequest); err != nil {
		response.RespondError(c, err)
	}
	passportResponse, err := h.documentService.CreatePassport(c.Request.Context(), &passportRequest)

	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, passportResponse)
}

func (h *DocumentHandler) UpdatePassport(c *gin.Context) {
	id := c.Param("id")
	var passportRequest dto.PassportUpdatedRequestDto
	if err := c.ShouldBindJSON(&passportRequest); err != nil {
		response.RespondError(c, err)
	}

	passportResponse, err := h.documentService.UpdatePassport(c.Request.Context(), id, &passportRequest)

	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, passportResponse)
}

func (h *DocumentHandler) RevokePassport(c *gin.Context) {
	id := c.Param("id")
	var passportRevokedRequest dto.PassportRevokedRequestDto

	if err := c.ShouldBindJSON(&passportRevokedRequest); err != nil {
		response.RespondError(c, err)
	}

	err := h.documentService.RevokePassport(c.Request.Context(), id, &passportRevokedRequest)
	if err != nil {
		response.RespondError(c, err)
	}

	response.RespondSuccess(c, "")
}

func (h *DocumentHandler) GetPassport(c *gin.Context) {
	id := c.Param("id")
	passport, err := h.documentService.GetPassportById(c.Request.Context(), id)
	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, passport)
}

func (h *DocumentHandler) GetPassports(c *gin.Context) {
	passports, err := h.documentService.GetPassports(c.Request.Context())
	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, passports)
}
