package handler

import (
	"be/internal/service"
	"be/internal/shared/constant"
	response "be/internal/shared/helper"
	"be/internal/transport/http/dto"

	"github.com/gin-gonic/gin"
)

type SchemaHandler struct {
	schemaService service.ISchemaService
}

func NewSchemaHandler(schemaService service.ISchemaService) *SchemaHandler {
	return &SchemaHandler{schemaService: schemaService}
}

func (h *SchemaHandler) GetSchemas(c *gin.Context) {
	schemas, err := h.schemaService.GetSchemas(c.Request.Context())
	if err != nil {
		response.RespondError(c, err)
		return
	}
	response.RespondSuccess(c, schemas)
}

func (h *SchemaHandler) GetSchemaByPublicId(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.RespondError(c, &constant.BadRequest)
		return
	}
	schemas, err := h.schemaService.GetSchemaByPublicId(c.Request.Context(), id)
	if err != nil {
		response.RespondError(c, err)
		return
	}
	response.RespondSuccess(c, schemas)
}

func (h *SchemaHandler) CreateSchema(c *gin.Context) {
	var request *dto.SchemaBuilderDto
	if err := c.ShouldBindJSON(request); err != nil {
		response.RespondError(c, err)
		return
	}
	schema, err := h.schemaService.CreateSchema(c.Request.Context(), request)
	if err != nil {
		response.RespondError(c, err)
	}
	response.RespondSuccess(c, schema)
}

func (h *SchemaHandler) RemoveSchema(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.RespondError(c, &constant.BadRequest)
		return
	}
	err := h.schemaService.RemoveSchema(c, id)
	if err != nil {
		response.RespondError(c, err)
		return
	}
	response.RespondSuccess(c, nil)
}
