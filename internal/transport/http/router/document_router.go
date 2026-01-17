package router

import (
	"be/internal/shared/constant"
	"be/internal/transport/http/handler"
	"be/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func (r *Router) SetupDocumentRouter(apiGroup *gin.RouterGroup, documentHandler *handler.DocumentHandler) {
	credentialGroup := apiGroup.Group("documents")
	credentialGroup.Use(middleware.AuthenticateMiddleware(r.authZkService))
	credentialGroup.Use(middleware.AuthorizeMiddleware([]constant.IdentityRole{constant.IdentityIssuerRole}))

	credentialGroup.GET("/:did", documentHandler.GetDocumentByHolderDID)

	citizenIdentityGroup := credentialGroup.Group("citizen_identity")
	academicDegreeGroup := credentialGroup.Group("academic_degree")
	healthInsuranceGroup := credentialGroup.Group("health_insurance")
	driverLicenseGroup := credentialGroup.Group("driver_license")
	passportGroup := credentialGroup.Group("passport")

	citizenIdentityGroup.GET("/:id", documentHandler.GetCitizenIdentity)
	citizenIdentityGroup.GET("", documentHandler.GetCitizenIdentities)
	citizenIdentityGroup.POST("", documentHandler.CreateCitizenIdentity)
	citizenIdentityGroup.PUT("/:id", documentHandler.UpdateCitizenIdentity)
	citizenIdentityGroup.PATCH("/:id", documentHandler.RevokeCitizenIdentity)

	academicDegreeGroup.GET("/:id", documentHandler.GetAcademicDegree)
	academicDegreeGroup.GET("", documentHandler.GetAcademicDegrees)
	academicDegreeGroup.POST("", documentHandler.CreateAcademicDegree)
	academicDegreeGroup.PUT("/:id", documentHandler.UpdateAcademicDegree)
	academicDegreeGroup.PATCH("/:id", documentHandler.RevokeAcademicDegree)

	healthInsuranceGroup.GET("/:id", documentHandler.GetHealthInsurance)
	healthInsuranceGroup.GET("", documentHandler.GetHealthInsurances)
	healthInsuranceGroup.POST("", documentHandler.CreateHealthInsurance)
	healthInsuranceGroup.PUT("/:id", documentHandler.UpdateHealthInsurance)
	healthInsuranceGroup.PATCH("/:id", documentHandler.RevokeHealthInsurance)

	driverLicenseGroup.GET("/:id", documentHandler.GetDriverLicense)
	driverLicenseGroup.GET("", documentHandler.GetDriverLicenses)
	driverLicenseGroup.POST("", documentHandler.CreateDriverLicense)
	driverLicenseGroup.PUT("/:id", documentHandler.UpdateDriverLicense)
	driverLicenseGroup.PATCH("/:id", documentHandler.RevokeDriverLicense)

	passportGroup.GET("/:id", documentHandler.GetPassport)
	passportGroup.GET("", documentHandler.GetPassports)
	passportGroup.POST("", documentHandler.CreatePassport)
	passportGroup.PUT("/:id", documentHandler.UpdatePassport)
	passportGroup.PATCH("/:id", documentHandler.RevokePassport)
}
