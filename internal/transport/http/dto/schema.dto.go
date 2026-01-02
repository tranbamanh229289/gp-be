package dto

type SchemaBuilderDto struct {
	IssuerDID   string               `json:"issuerDID"`
	Title       string               `json:"title"`
	Type        string               `json:"type"`
	Version     string               `json:"version"`
	Description string               `json:"description"`
	IsMerklized bool                 `json:"isMerklized"`
	Attributes  []SchemaAttributeDto `json:"attributes"`
}

type SchemaAttributeDto struct {
	Name        string                 `json:"name"`
	Title       string                 `json:"title"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Required    bool                   `json:"required"`
	Slot        string                 `json:"slot"`
	Format      string                 `json:"format,omitempty"`
	Pattern     string                 `json:"pattern,omitempty"`
	MinLength   *int                   `json:"minLength,omitempty"`
	MaxLength   *int                   `json:"maxLength,omitempty"`
	Minimum     *float64               `json:"minimum,omitempty"`
	Maximum     *float64               `json:"maximum,omitempty"`
	Enum        map[string]interface{} `json:"enum,omitempty"`
}

type ClaimDataDto struct {
	SchemaHash        string                 `json:"schema_hash"`
	Type              string                 `json:"type"`
	IsMerklized       bool                   `json:"is_merklized"`
	CredentialSubject map[string]interface{} `json:"credential_subject"`
	SlotIndexMapping  map[string]string      `json:"slot_index_mapping,omitempty"`
}

type SchemaResponseDto struct {
	PublicID    string               `json:"id"`
	IssuerDID   string               `json:"issuerDID"`
	Hash        string               `json:"hash"`
	Title       string               `json:"title"`
	Type        string               `json:"type"`
	Version     string               `json:"version"`
	Status      string               `json:"status"`
	Description string               `json:"description"`
	IsMerklized bool                 `json:"isMerklized"`
	SchemaURL   string               `json:"contextURL"`
	ContextURL  string               `json:"schemaURL"`
	Attributes  []SchemaAttributeDto `json:"attributes"`
}
