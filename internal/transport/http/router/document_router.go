package router

import (
	"be/internal/transport/http/handler"

	"github.com/gin-gonic/gin"
)

func SetupCredentialRouter(apiGroup *gin.RouterGroup, credentialHandler *handler.DocumentHandler) {
	credentialGroup := apiGroup.Group("credential")
	citizenIdentityGroup := credentialGroup.Group("citizen-identity")
	academicDegreeGroup := credentialGroup.Group("academic-degree")
	healthInsuranceGroup := credentialGroup.Group("health-insurance")
	driverLicenseGroup := credentialGroup.Group("driver-license")
	passportGroup := credentialGroup.Group("passport")

	citizenIdentityGroup.GET("/:id", credentialHandler.GetCitizenIdentity)
	citizenIdentityGroup.GET("", credentialHandler.GetCitizenIdentities)
	citizenIdentityGroup.POST("", credentialHandler.CreateCitizenIdentity)
	citizenIdentityGroup.PUT("/:id", credentialHandler.UpdateCitizenIdentity)
	citizenIdentityGroup.PATCH("/:id", credentialHandler.RevokeCitizenIdentity)

	academicDegreeGroup.GET("/:id", credentialHandler.GetAcademicDegree)
	academicDegreeGroup.GET("", credentialHandler.GetAcademicDegrees)
	academicDegreeGroup.POST("", credentialHandler.CreateAcademicDegree)
	academicDegreeGroup.PUT("/:id", credentialHandler.UpdateAcademicDegree)
	academicDegreeGroup.PATCH("/:id", credentialHandler.RevokeAcademicDegree)

	healthInsuranceGroup.GET("/:id", credentialHandler.GetHealthInsurance)
	healthInsuranceGroup.GET("", credentialHandler.GetHealthInsurances)
	healthInsuranceGroup.POST("", credentialHandler.CreateHealthInsurance)
	healthInsuranceGroup.PUT("/:id", credentialHandler.UpdateHealthInsurance)
	healthInsuranceGroup.PATCH("/:id", credentialHandler.RevokeHealthInsurance)

	driverLicenseGroup.GET("/:id", credentialHandler.GetDriverLicense)
	driverLicenseGroup.GET("", credentialHandler.GetDriverLicenses)
	driverLicenseGroup.POST("", credentialHandler.CreateDriverLicense)
	driverLicenseGroup.PUT("/:id", credentialHandler.UpdateDriverLicense)
	driverLicenseGroup.PATCH("/:id", credentialHandler.RevokeDriverLicense)

	passportGroup.GET("/:id", credentialHandler.GetPassport)
	passportGroup.GET("", credentialHandler.GetPassports)
	passportGroup.POST("", credentialHandler.CreatePassport)
	passportGroup.PUT("/:id", credentialHandler.UpdatePassport)
	passportGroup.PATCH("/:id", credentialHandler.RevokePassport)
}
