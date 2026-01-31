package handler

import (
	"be/internal/service"
	"be/internal/shared/constant"
	"be/internal/shared/helper"

	"github.com/gin-gonic/gin"
)

type StatisticHandler struct {
	statisticService service.IStatisticService
}

func NewStatisticHandler(statisticService service.IStatisticService) *StatisticHandler {
	return &StatisticHandler{statisticService: statisticService}
}

func (h *StatisticHandler) GetIssuerStatisticByIssuerDID(c *gin.Context) {
	did := c.Param("did")
	if did == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}

	resp, err := h.statisticService.GetIssuerStatisticByDID(c.Request.Context(), did)
	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, resp)
}

func (h *StatisticHandler) GetHolderStatisticByHolderDID(c *gin.Context) {
	did := c.Param("did")
	if did == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}

	resp, err := h.statisticService.GetHolderStatisticByDID(c.Request.Context(), did)
	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, resp)
}

func (h *StatisticHandler) GetVerifierStatisticByVerifierDID(c *gin.Context) {
	did := c.Param("did")
	if did == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}

	resp, err := h.statisticService.GetVerifierStatisticByDID(c.Request.Context(), did)
	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, resp)

}
