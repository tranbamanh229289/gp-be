package dto

import (
	"be/internal/domain/schema"
	"be/internal/shared/constant"
)

type SchemaBuilderDto struct {
	IssuerDID    string                `json:"issuerDID"`
	DocumentType constant.DocumentType `json:"documentType"`
	Title        string                `json:"title"`
	Type         string                `json:"type"`
	Version      string                `json:"version"`
	Description  string                `json:"description"`
	IsMerklized  bool                  `json:"isMerklized"`
	Attributes   []SchemaAttributeDto  `json:"attributes"`
}

type SchemaAttributeDto struct {
	Name        string                 `json:"name"`
	Title       string                 `json:"title"`
	Type        constant.AttributeType `json:"type"`
	Description string                 `json:"description"`
	Required    bool                   `json:"required"`
	Slot        constant.Slot          `json:"slot"`
	Format      string                 `json:"format,omitempty"`
	Pattern     string                 `json:"pattern,omitempty"`
	MinLength   *int                   `json:"minLength,omitempty"`
	MaxLength   *int                   `json:"maxLength,omitempty"`
	Minimum     *float64               `json:"minimum,omitempty"`
	Maximum     *float64               `json:"maximum,omitempty"`
	Enum        map[string]interface{} `json:"enum,omitempty"`
}
type SchemaResponseDto struct {
	PublicID     string                `json:"id"`
	IssuerDID    string                `json:"issuerDID"`
	IssuerName   string                `json:"issuerName"`
	DocumentType constant.DocumentType `json:"documentType"`
	Hash         string                `json:"hash"`
	Title        string                `json:"title"`
	Type         string                `json:"type"`
	Version      string                `json:"version"`
	Description  string                `json:"description"`
	Status       constant.SchemaStatus `json:"status"`
	IsMerklized  bool                  `json:"isMerklized"`
	SchemaURL    string                `json:"schemaURL"`
	ContextURL   string                `json:"contextURL"`
	Attributes   []SchemaAttributeDto  `json:"attributes"`
}

func ToSchemaResponseDto(schema *schema.Schema) *SchemaResponseDto {
	var attributesDtos []SchemaAttributeDto
	for _, item := range schema.SchemaAttributes {
		attributesDtos = append(attributesDtos, SchemaAttributeDto{
			Name:        item.Name,
			Title:       item.Title,
			Type:        item.Type,
			Description: item.Description,
			Required:    item.Required,
			Slot:        item.Slot,
		})
	}
	return &SchemaResponseDto{
		PublicID:     schema.PublicID.String(),
		IssuerDID:    schema.IssuerDID,
		IssuerName:   schema.Issuer.Name,
		DocumentType: schema.DocumentType,
		Hash:         schema.Hash,
		Title:        schema.Title,
		Type:         schema.Type,
		Version:      schema.Version,
		Description:  schema.Description,
		Status:       schema.Status,
		IsMerklized:  schema.IsMerklized,
		SchemaURL:    schema.SchemaURL,
		ContextURL:   schema.ContextURL,
		Attributes:   attributesDtos,
	}
}
