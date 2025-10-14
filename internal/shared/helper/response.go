package helper

import (
	"be/internal/shared/constant"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code string          		`json:"code"`
	Message string    		`json:"message"`
	Data interface{}			`json:"data,omitempty"`
	Metadata interface{}  `json:"metadata,omitempty"`
}

type Pagination struct {
	Page int		`json:"page"`
	Limit int		`json:"limit"`
	Total int		`json:"total"`
}

func RespondSuccess(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, &Response{
		Code: "SUCCESS",
		Message: "Success",
		Data: data,
	})
}

func RespondError(ctx *gin.Context, err *constant.Errors) {
	ctx.JSON(http.StatusOK, &Response{
		Code: err.Code,
		Message: err.Message,
	})
}

func RespondWithPaginationSuccess(ctx *gin.Context, data interface{}, metadata interface{}) {
	ctx.JSON(http.StatusOK, &Response{
		Code: "SUCCESS",
		Message: "Success",
		Data: data,
		Metadata: metadata,
	})
}

