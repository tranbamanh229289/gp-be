package handler

import (
	"be/internal/service"
	"be/internal/shared/constant"
	"be/internal/shared/helper"
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

func (h *DocumentHandler) GetDocumentByHolderDID(c *gin.Context) {
	did := c.Param("did")
	if did == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}
	docType := c.Query("documentType")
	if docType == "" {
		helper.RespondError(c, &constant.BadRequest)
	}
	switch constant.DocumentType(docType) {
	case constant.CitizenIdentity:
		entity, err := h.documentService.GetCitizenIdentityByHolderDID(c.Request.Context(), did)
		if err != nil {
			helper.RespondError(c, err)
			return
		}
		helper.RespondSuccess(c, entity)
	case constant.AcademicDegree:
		entity, err := h.documentService.GetAcademicDegreeByHolderDID(c.Request.Context(), did)
		if err != nil {
			helper.RespondError(c, err)
			return
		}
		helper.RespondSuccess(c, entity)
	case constant.HealthInsurance:
		entity, err := h.documentService.GetHealthInsuranceByHolderDID(c.Request.Context(), did)
		if err != nil {
			helper.RespondError(c, err)
			return
		}
		helper.RespondSuccess(c, entity)
	case constant.DriverLicense:
		entity, err := h.documentService.GetDriverLicenseByHolderDID(c.Request.Context(), did)
		if err != nil {
			helper.RespondError(c, err)
			return
		}
		helper.RespondSuccess(c, entity)
	case constant.Passport:
		entity, err := h.documentService.GetPassportByHolderDID(c.Request.Context(), did)
		if err != nil {
			helper.RespondError(c, err)
			return
		}
		helper.RespondSuccess(c, entity)
	}

}
func (h *DocumentHandler) CreateCitizenIdentity(c *gin.Context) {
	var citizenIdentityRequest dto.CitizenIdentityCreatedRequestDto

	if err := c.ShouldBindJSON(&citizenIdentityRequest); err != nil {
		helper.RespondError(c, err)
		return
	}

	user, ok := c.Get("user")
	if !ok {
		helper.RespondError(c, &constant.InternalServer)
		return
	}

	claims, ok := user.(*dto.ZKClaims)
	if !ok {
		helper.RespondError(c, &constant.InternalServer)
		return
	}

	citizenIdentityRequest.IssuerDID = claims.DID

	citizenIdentityResponse, err := h.documentService.CreateCitizenIdentity(c.Request.Context(), &citizenIdentityRequest)

	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, citizenIdentityResponse)
}

func (h *DocumentHandler) UpdateCitizenIdentity(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}
	var citizenIdentityRequest dto.CitizenIdentityUpdatedRequestDto
	if err := c.ShouldBindJSON(&citizenIdentityRequest); err != nil {
		helper.RespondError(c, err)
		return
	}

	citizenIdentityResponse, err := h.documentService.UpdateCitizenIdentity(c.Request.Context(), id, &citizenIdentityRequest)

	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, citizenIdentityResponse)
}

func (h *DocumentHandler) RevokeCitizenIdentity(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}
	var citizenIdentityRevokedRequest dto.CitizenIdentityOptionRequestDto

	if err := c.ShouldBindJSON(&citizenIdentityRevokedRequest); err != nil {
		helper.RespondError(c, err)
		return
	}

	err := h.documentService.RevokeCitizenIdentity(c.Request.Context(), id, &citizenIdentityRevokedRequest)
	if err != nil {
		helper.RespondError(c, err)
		return
	}

	helper.RespondSuccess(c, "")
}

func (h *DocumentHandler) GetCitizenIdentity(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}
	citizenIdentity, err := h.documentService.GetCitizenIdentityByPublicId(c.Request.Context(), id)
	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, citizenIdentity)
}

func (h *DocumentHandler) GetCitizenIdentities(c *gin.Context) {
	citizenIdentities, err := h.documentService.GetCitizenIdentities(c.Request.Context())
	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, citizenIdentities)
}

func (h *DocumentHandler) CreateAcademicDegree(c *gin.Context) {
	var academicDegreeRequest dto.AcademicDegreeCreatedRequestDto
	if err := c.ShouldBindJSON(&academicDegreeRequest); err != nil {
		helper.RespondError(c, err)
		return
	}
	user, ok := c.Get("user")
	if !ok {
		helper.RespondError(c, &constant.InternalServer)
		return
	}

	claims, ok := user.(*dto.ZKClaims)
	if !ok {
		helper.RespondError(c, &constant.InternalServer)
		return
	}
	academicDegreeRequest.IssuerDID = claims.DID
	academicDegreeResponse, err := h.documentService.CreateAcademicDegree(c.Request.Context(), &academicDegreeRequest)

	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, academicDegreeResponse)
}

func (h *DocumentHandler) UpdateAcademicDegree(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}
	var academicDegreeRequest dto.AcademicDegreeUpdatedRequestDto
	if err := c.ShouldBindJSON(&academicDegreeRequest); err != nil {
		helper.RespondError(c, err)
		return
	}

	academicDegreeResponse, err := h.documentService.UpdateAcademicDegree(c.Request.Context(), id, &academicDegreeRequest)

	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, academicDegreeResponse)
}

func (h *DocumentHandler) RevokeAcademicDegree(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}
	var academicDegreeRevokedRequest dto.AcademicDegreeOptionRequestDto

	if err := c.ShouldBindJSON(&academicDegreeRevokedRequest); err != nil {
		helper.RespondError(c, err)
		return
	}

	err := h.documentService.RevokeAcademicDegree(c.Request.Context(), id, &academicDegreeRevokedRequest)
	if err != nil {
		helper.RespondError(c, err)
		return
	}

	helper.RespondSuccess(c, "")
}

func (h *DocumentHandler) GetAcademicDegree(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}
	academicDegree, err := h.documentService.GetAcademicDegreeByPublicId(c.Request.Context(), id)
	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, academicDegree)
}

func (h *DocumentHandler) GetAcademicDegrees(c *gin.Context) {
	academicDegrees, err := h.documentService.GetAcademicDegrees(c.Request.Context())
	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, academicDegrees)
}

func (h *DocumentHandler) CreateHealthInsurance(c *gin.Context) {
	var healthInsuranceRequest dto.HealthInsuranceCreatedRequestDto

	if err := c.ShouldBindJSON(&healthInsuranceRequest); err != nil {
		helper.RespondError(c, err)
		return
	}

	user, ok := c.Get("user")
	if !ok {
		helper.RespondError(c, &constant.InternalServer)
		return
	}

	claims, ok := user.(*dto.ZKClaims)
	if !ok {
		helper.RespondError(c, &constant.InternalServer)
		return
	}
	healthInsuranceRequest.IssuerDID = claims.DID

	academyInsuranceResponse, err := h.documentService.CreateHealthInsurance(c.Request.Context(), &healthInsuranceRequest)

	if err != nil {
		helper.RespondError(c, err)
		return
	}

	helper.RespondSuccess(c, academyInsuranceResponse)
}

func (h *DocumentHandler) UpdateHealthInsurance(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}
	var healthInsuranceRequest dto.HealthInsuranceUpdatedRequestDto
	if err := c.ShouldBindJSON(&healthInsuranceRequest); err != nil {
		helper.RespondError(c, err)
		return
	}

	healthInsuranceResponse, err := h.documentService.UpdateHealthInsurance(c.Request.Context(), id, &healthInsuranceRequest)

	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, healthInsuranceResponse)
}

func (h *DocumentHandler) RevokeHealthInsurance(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}
	var healthInsuranceRevokedRequest dto.HealthInsuranceOptionRequestDto

	if err := c.ShouldBindJSON(&healthInsuranceRevokedRequest); err != nil {
		helper.RespondError(c, err)
		return
	}

	err := h.documentService.RevokeHealthInsurance(c.Request.Context(), id, &healthInsuranceRevokedRequest)
	if err != nil {
		helper.RespondError(c, err)
		return
	}

	helper.RespondSuccess(c, "")
}

func (h *DocumentHandler) GetHealthInsurance(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}
	healthInsurance, err := h.documentService.GetHealthInsuranceByPublicId(c.Request.Context(), id)
	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, healthInsurance)
}

func (h *DocumentHandler) GetHealthInsurances(c *gin.Context) {
	healthInsurances, err := h.documentService.GetHealthInsurances(c.Request.Context())
	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, healthInsurances)
}

func (h *DocumentHandler) CreateDriverLicense(c *gin.Context) {
	var driverLicenseRequest dto.DriverLicenseCreatedRequestDto
	if err := c.ShouldBindJSON(&driverLicenseRequest); err != nil {
		helper.RespondError(c, err)
		return
	}

	user, ok := c.Get("user")
	if !ok {
		helper.RespondError(c, &constant.InternalServer)
		return
	}

	claims, ok := user.(*dto.ZKClaims)
	if !ok {
		helper.RespondError(c, &constant.InternalServer)
		return
	}
	driverLicenseRequest.IssuerDID = claims.DID

	driverLicenseResponse, err := h.documentService.CreateDriverLicense(c.Request.Context(), &driverLicenseRequest)

	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, driverLicenseResponse)
}

func (h *DocumentHandler) UpdateDriverLicense(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}
	var driverLicenseRequest dto.DriverLicenseUpdatedRequestDto
	if err := c.ShouldBindJSON(&driverLicenseRequest); err != nil {
		helper.RespondError(c, err)
		return
	}

	driverLicenseResponse, err := h.documentService.UpdateDriverLicense(c.Request.Context(), id, &driverLicenseRequest)

	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, driverLicenseResponse)
}

func (h *DocumentHandler) RevokeDriverLicense(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}
	var driverLicenseRevokedRequest dto.DriverLicenseOptionRequestDto

	if err := c.ShouldBindJSON(&driverLicenseRevokedRequest); err != nil {
		helper.RespondError(c, err)
		return
	}

	err := h.documentService.RevokeDriverLicense(c.Request.Context(), id, &driverLicenseRevokedRequest)
	if err != nil {
		helper.RespondError(c, err)
		return
	}

	helper.RespondSuccess(c, "")
}

func (h *DocumentHandler) GetDriverLicense(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}
	driverLicense, err := h.documentService.GetDriverLicenseByPublicId(c.Request.Context(), id)
	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, driverLicense)
}

func (h *DocumentHandler) GetDriverLicenses(c *gin.Context) {
	driverLicenses, err := h.documentService.GetDriverLicenses(c.Request.Context())
	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, driverLicenses)
}

func (h *DocumentHandler) CreatePassport(c *gin.Context) {
	var passportRequest dto.PassportCreatedRequestDto
	if err := c.ShouldBindJSON(&passportRequest); err != nil {
		helper.RespondError(c, err)
		return
	}

	user, ok := c.Get("user")
	if !ok {
		helper.RespondError(c, &constant.InternalServer)
		return
	}

	claims, ok := user.(*dto.ZKClaims)
	if !ok {
		helper.RespondError(c, &constant.InternalServer)
		return
	}

	passportRequest.IssuerDID = claims.DID

	passportResponse, err := h.documentService.CreatePassport(c.Request.Context(), &passportRequest)

	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, passportResponse)
}

func (h *DocumentHandler) UpdatePassport(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}
	var passportRequest dto.PassportUpdatedRequestDto
	if err := c.ShouldBindJSON(&passportRequest); err != nil {
		helper.RespondError(c, err)
		return
	}

	passportResponse, err := h.documentService.UpdatePassport(c.Request.Context(), id, &passportRequest)

	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, passportResponse)
}

func (h *DocumentHandler) RevokePassport(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}
	var passportRevokedRequest dto.PassportOptionRequestDto

	if err := c.ShouldBindJSON(&passportRevokedRequest); err != nil {
		helper.RespondError(c, err)
		return
	}
	err := h.documentService.RevokePassport(c.Request.Context(), id, &passportRevokedRequest)
	if err != nil {
		helper.RespondError(c, err)
		return
	}

	helper.RespondSuccess(c, "")
}

func (h *DocumentHandler) GetPassport(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}
	passport, err := h.documentService.GetPassportByPublicId(c.Request.Context(), id)
	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, passport)
}

func (h *DocumentHandler) GetPassports(c *gin.Context) {
	passports, err := h.documentService.GetPassports(c.Request.Context())
	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, passports)
}
