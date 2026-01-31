package handler

import (
	"be/internal/service"
	"be/internal/shared/constant"
	"be/internal/shared/helper"
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
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, schemas)
}

func (h *SchemaHandler) GetSchemaByPublicId(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}
	schema, err := h.schemaService.GetSchemaByPublicId(c.Request.Context(), id)
	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, schema)
}

func (h *SchemaHandler) GetSchemaAttributeByPublicId(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}
	schemaAttributes, err := h.schemaService.GetSchemaAttributesBySchemaId(c.Request.Context(), id)
	if err != nil {
		helper.RespondError(c, &constant.BadRequest)
		return
	}
	helper.RespondSuccess(c, schemaAttributes)
}

func (h *SchemaHandler) CreateSchema(c *gin.Context) {
	var request dto.SchemaBuilderDto
	if err := c.ShouldBindJSON(&request); err != nil {
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
	request.IssuerDID = claims.DID

	schema, err := h.schemaService.CreateSchema(c.Request.Context(), &request)
	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, schema)
}

func (h *SchemaHandler) RemoveSchema(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.RespondError(c, &constant.BadRequest)
		return
	}
	err := h.schemaService.RemoveSchema(c, id)
	if err != nil {
		helper.RespondError(c, err)
		return
	}
	helper.RespondSuccess(c, nil)
}
