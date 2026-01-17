package handler

import (
	"be/internal/service"

	"github.com/gin-gonic/gin"
)

type StatisticHandler struct {
	statisticService service.IStatisticService
}

func GetStatistic(c *gin.Context) {

}
