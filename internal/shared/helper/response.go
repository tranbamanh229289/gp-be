package response

import (
	"be/internal/shared/constant"
	"errors"
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

func RespondError(ctx *gin.Context, err error) {
	var appErrors *constant.Errors
	if errors.As(err, &appErrors) {
		ctx.JSON(http.StatusOK, &Response{
				Code: appErrors.Code,
				Message: appErrors.Message,
		})
	}
}

func RespondWithPaginationSuccess(ctx *gin.Context, data interface{}, metadata interface{}) {
	ctx.JSON(http.StatusOK, &Response{
		Code: "SUCCESS",
		Message: "Success",
		Data: data,
		Metadata: metadata,
	})
}

