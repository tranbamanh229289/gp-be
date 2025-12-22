package dto

type SchemaBuilderDto struct {
	IssuerID    string
	Title       string
	Type        string
	Version     string
	Description string

	IsMerklized   bool
	IsPublishIPFS bool
	Fields        []SchemaFieldDto
}

type SchemaFieldDto struct {
	Name        string
	Title       string
	Type        string
	Description string
	Required    bool
	Slot        string
}

type SchemaCreatedRequestDto struct {
	SchemaBuilderDto
}
type SchemaResponseDto struct {
	PublicID  string
	IssuerID  string
	Type      string
	Version   string
	JSONCID   string
	JSONLDCID string
	Status    string
}
