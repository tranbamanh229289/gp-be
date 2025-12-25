package handler

import (
	"be/internal/service"
	response "be/internal/shared/helper"
	"be/internal/transport/http/dto"

	"github.com/gin-gonic/gin"
)

type SchemaHandler struct {
	schemaService service.ISchemaService
}

func NewSchemaService(schemaService service.ISchemaService) *SchemaHandler {
	return &SchemaHandler{schemaService: schemaService}
}

func (h *SchemaHandler) CreateSchema(c *gin.Context) {
	var request *dto.SchemaCreatedRequestDto
	if err := c.ShouldBindJSON(request); err != nil {
		response.RespondError(c, err)
	}
	schema, err := h.schemaService.CreateSchema(c.Request.Context(), request)
	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, schema)
}

func (h *SchemaHandler) RemoveSchema(c *gin.Context) {
	id := c.Param("id")
	err := h.schemaService.RemoveSchema(c, id)
	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, nil)
}
