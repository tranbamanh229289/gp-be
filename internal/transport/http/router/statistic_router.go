package router

import (
	"be/internal/transport/http/handler"
	"be/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func (r *Router) SetupStatisticRouter(apiGroup *gin.RouterGroup, statisticHandler *handler.StatisticHandler) {
	statisticGroup := apiGroup.Group("statistic")
	statisticGroup.Use(middleware.AuthenticateMiddleware(r.authZkService))

	statisticGroup.GET("/issuer/:did", statisticHandler.GetIssuerStatisticByIssuerDID)
	statisticGroup.GET("/holder/:did", statisticHandler.GetHolderStatisticByHolderDID)
	statisticGroup.GET("/verifier/:did", statisticHandler.GetVerifierStatisticByVerifierDID)
}
