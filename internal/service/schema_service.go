package service

import (
	"be/config"
	"be/internal/domain/schema"
	"be/internal/infrastructure/ipfs"
	"be/internal/shared/constant"
	"be/internal/transport/http/dto"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	core "github.com/iden3/go-iden3-core/v2"
	"github.com/iden3/go-iden3-crypto/keccak256"
	"gorm.io/gorm"
)

type ISchemaService interface {
	GetSchemas(ctx context.Context) ([]*dto.SchemaResponseDto, error)
	GetSchemaByPublicId(ctx context.Context, id string) (*dto.SchemaResponseDto, error)
	GetSchemaAttributesBySchemaId(ctx context.Context, id string) ([]*dto.SchemaAttributeDto, error)
	CreateSchema(ctx context.Context, request *dto.SchemaBuilderDto) (*dto.SchemaResponseDto, error)
	RemoveSchema(ctx context.Context, id string) error
}
type SchemaService struct {
	config              *config.Config
	ipfs                *ipfs.Pinata
	identityRepo        schema.IIdentityRepository
	schemaRepo          schema.ISchemaRepository
	schemaAttributeRepo schema.ISchemaAttributeRepository
}

func NewSchemaService(
	config *config.Config,
	pinata *ipfs.Pinata,
	identityRepo schema.IIdentityRepository,
	schemaRepo schema.ISchemaRepository,
	schemaAttributeRepo schema.ISchemaAttributeRepository,
) ISchemaService {

	return &SchemaService{
		config:              config,
		ipfs:                pinata,
		identityRepo:        identityRepo,
		schemaRepo:          schemaRepo,
		schemaAttributeRepo: schemaAttributeRepo,
	}
}

func (s *SchemaService) GetSchemaByPublicId(ctx context.Context, id string) (*dto.SchemaResponseDto, error) {
	schema, err := s.schemaRepo.FindSchemaByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.SchemaNotFound
		}
		return nil, &constant.InternalServer
	}

	return dto.ToSchemaResponseDto(schema), nil
}

func (s *SchemaService) GetSchemas(ctx context.Context) ([]*dto.SchemaResponseDto, error) {
	schemas, err := s.schemaRepo.FindAllSchemas(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.SchemaNotFound
		}
		return nil, &constant.InternalServer
	}

	var schemaDtos []*dto.SchemaResponseDto
	for _, schema := range schemas {
		schemaDtos = append(schemaDtos, dto.ToSchemaResponseDto(schema))
	}
	return schemaDtos, nil
}

func (s *SchemaService) GetSchemaAttributesBySchemaId(ctx context.Context, id string) ([]*dto.SchemaAttributeDto, error) {
	schema, err := s.schemaRepo.FindSchemaByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.SchemaNotFound
		}
		return nil, &constant.InternalServer
	}
	var resp []*dto.SchemaAttributeDto
	for _, item := range schema.SchemaAttributes {
		resp = append(resp, &dto.SchemaAttributeDto{
			Name:        item.Name,
			Title:       item.Title,
			Type:        item.Type,
			Description: item.Description,
			Slot:        item.Slot,
		})
	}
	return resp, nil
}

func (s *SchemaService) CreateSchema(ctx context.Context, request *dto.SchemaBuilderDto) (*dto.SchemaResponseDto, error) {
	issuer, err := s.identityRepo.FindIdentityByDID(ctx, request.IssuerDID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.IdentityNotFound
		}
		return nil, &constant.InternalServer
	}

	if err := s.validate(request); err != nil {
		return nil, err
	}

	jsonSchema := s.generateJSONSchema(request)
	jsonLdContext := s.generateJSONLDContext(request)

	schemaCID, contextCID, err := s.uploadToIPFS(jsonSchema, jsonLdContext, request.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to IPFS: %w", err)
	}
	schemaURL, err := url.JoinPath(s.config.IPFS.GatewayURL, schemaCID)
	if err != nil {
		return nil, fmt.Errorf("faild to join path")
	}
	contextURL, _ := url.JoinPath(s.config.IPFS.GatewayURL, contextCID)
	if err != nil {
		return nil, fmt.Errorf("faild to join path")
	}

	hash, err := s.getSchemaHash(contextURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get schema hash: %w", err)
	}

	schemaEntity := &schema.Schema{
		PublicID:      uuid.New(),
		IssuerDID:     request.IssuerDID,
		DocumentType:  request.DocumentType,
		Hash:          hash.BigInt().String(),
		Type:          request.Type,
		Version:       request.Version,
		Title:         request.Title,
		Description:   request.Description,
		IsMerklized:   request.IsMerklized,
		JSONSchema:    jsonSchema,
		JSONLDContext: jsonLdContext,
		SchemaURL:     schemaURL,
		ContextURL:    contextURL,
		Status:        constant.SchemaActiveStatus,
	}

	var attributeEntities []*schema.SchemaAttribute
	for _, item := range request.Attributes {
		attributeEntities = append(attributeEntities, &schema.SchemaAttribute{
			PublicID:    uuid.New(),
			Name:        item.Name,
			Title:       item.Title,
			Type:        item.Type,
			Required:    item.Required,
			Description: item.Description,
			Slot:        item.Slot,
		})
	}

	schemaEntity.SchemaAttributes = attributeEntities
	schemaCreated, err := s.schemaRepo.CreateSchema(ctx, schemaEntity)
	if err != nil {
		return nil, &constant.InternalServer
	}
	schemaCreated.Issuer = issuer
	return dto.ToSchemaResponseDto(schemaCreated), nil
}

func (s *SchemaService) RemoveSchema(ctx context.Context, id string) error {
	schema, err := s.schemaRepo.FindSchemaByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &constant.SchemaNotFound
		}
		return &constant.InternalServer
	}
	schemaURlParts := strings.Split(schema.SchemaURL, "/ipfs/")
	contextURLParts := strings.Split(schema.ContextURL, "/ipfs/")

	if schema.SchemaURL != "" {
		err := s.ipfs.Remove(schemaURlParts[1])
		if err != nil {
			return err
		}
	}

	if schema.ContextURL != "" {
		err := s.ipfs.Remove(contextURLParts[1])
		if err != nil {
			return err
		}
	}

	changes := map[string]interface{}{"status": constant.SchemaRevokeStatus, "revoked_at": time.Now().UTC()}

	if err := s.schemaRepo.UpdateSchema(ctx, schema, changes); err != nil {
		return err
	}

	return nil
}

func (s *SchemaService) validate(request *dto.SchemaBuilderDto) error {
	if request.Type == "" {
		return errors.New("type is required")
	}
	if request.Version == "" {
		return errors.New("version is required")
	}
	if request.IssuerDID == "" {
		return errors.New("issuer_did is required")
	}

	if !request.IsMerklized && len(request.Attributes) > 4 {
		return errors.New("non-merklized schema supports maximum 4 attributes")
	}

	seen := make(map[string]bool)
	for _, attr := range request.Attributes {

		if attr.Name == "" {
			return errors.New("attribute name is required")
		}

		if seen[attr.Name] {
			return fmt.Errorf("duplicate attribute name: %s", attr.Name)
		}
		seen[attr.Name] = true
	}
	return nil
}

func (s *SchemaService) generateJSONSchema(request *dto.SchemaBuilderDto) map[string]interface{} {
	if request.Type == "" {
		request.Type = "Credential"
	}
	if request.Version == "" {
		request.Version = "1.0"
	}

	baseURL := strings.TrimSuffix(s.config.IPFS.GatewayURL, "/")
	if baseURL == "" {
		baseURL = "https://ipfs.io/ipfs"
	}

	schemaURL := fmt.Sprintf("%s/%s.json", baseURL, strings.ToLower(request.Type))
	contextURL := fmt.Sprintf("%s/%s.jsonld", baseURL, strings.ToLower(request.Type))

	iden3Serialization := map[constant.Slot]string{}
	credSubProperties := map[string]interface{}{
		"id": map[string]interface{}{
			"type":        "string",
			"format":      "uri",
			"title":       "Credential subject ID",
			"description": "Stores the DID of the subject that owns the credential",
		},
	}
	credSubRequired := []string{"id"}

	for _, attr := range request.Attributes {
		attrType := "integer"
		switch attr.Type {
		case "integer":
			attrType = "integer"
		case "number":
			attrType = "integer"
		case "boolean":
			attrType = "boolean"
		case "string":
			attrType = "integer"
		case "dateTime":
			attrType = "integer"
		}
		credSubProperties[attr.Name] = map[string]interface{}{
			"type":        attrType,
			"title":       attr.Title,
			"description": attr.Description,
		}
		if attr.Required {
			credSubRequired = append(credSubRequired, attr.Name)
		}
		if !request.IsMerklized && attr.Slot != "" {
			iden3Serialization[attr.Slot] = attr.Name
		}
	}

	schema := map[string]interface{}{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"$metadata": map[string]interface{}{
			"uris": map[string]string{
				"jsonLdContext": contextURL,
				"jsonSchema":    schemaURL,
			},
			"version": request.Version,
			"type":    request.Type,
		},
		"title":       request.Title,
		"description": request.Description,
		"type":        "object",
		"required": []string{
			"credentialSubject", "@context", "id", "issuanceDate",
			"issuer", "type", "credentialSchema",
		},
		"properties": map[string]interface{}{
			"@context": map[string]interface{}{"type": []string{"string", "array", "object"}},
			"id":       map[string]interface{}{"type": "string"},
			"type": map[string]interface{}{
				"type": []string{"string", "array"},
				"items": map[string]interface{}{
					"type": "string",
				},
			},
			"issuer": map[string]interface{}{
				"type":     []string{"string", "object"},
				"required": []string{"id"},
				"format":   "uri",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":   "string",
						"format": "uri",
					},
				},
			},
			"issuanceDate": map[string]interface{}{
				"type":   "string",
				"format": "date-time",
			},
			"expirationDate": map[string]interface{}{
				"type":   "string",
				"format": "date-time",
			},
			"credentialSchema": map[string]interface{}{
				"type":     "object",
				"required": []string{"id", "type"},
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":   "string",
						"format": "uri",
					},
					"type": map[string]interface{}{
						"type": "string",
					},
				},
			},
			"credentialStatus": map[string]interface{}{
				"description": "Allows the discovery of information about the current status of the credential, such as whether it is suspended or revoked.",
				"title":       "Credential Status",
				"type":        "object",
				"required":    []string{"id", "type"},
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"description": "Id URL of the credentialStatus.",
						"title":       "Id",
						"type":        "string",
						"format":      "uri",
					},
					"type": map[string]interface{}{
						"description": "Expresses the credential status type (method). The value should provide enough information to determine the current status of the credential.",
						"title":       "Type",
						"type":        "string",
					},
				},
			},
			"credentialSubject": map[string]interface{}{
				"type":        "object",
				"title":       "Credential subject",
				"description": "Stores the data of the credential",
				"properties":  credSubProperties,
				"required":    credSubRequired,
			},
		},
	}

	if len(iden3Serialization) > 0 {
		schema["$metadata"].(map[string]interface{})["iden3Serialization"] = iden3Serialization
	}

	return schema
}

func (s *SchemaService) generateJSONLDContext(request *dto.SchemaBuilderDto) map[string]interface{} {
	credType := request.Type
	if credType == "" {
		credType = "VerifiableCredential"
	}

	credUUID := uuid.New()
	vocabUUID := uuid.New()

	credURN := "urn:uuid:" + credUUID.String()
	vocabURN := "urn:uuid:" + vocabUUID.String() + "#"

	serializationParts := []string{}

	innerContext := map[string]interface{}{
		"@protected":  true,
		"@version":    1.1,
		"id":          "@id",
		"type":        "@type",
		"iden3-vocab": vocabURN,
		"xsd":         "http://www.w3.org/2001/XMLSchema#",
	}

	for _, field := range request.Attributes {
		xsdType := "xsd:integer"
		switch field.Type {
		case "integer":
			xsdType = "xsd:integer"
		case "number":
			xsdType = "xsd:double"
		case "boolean":
			xsdType = "xsd:boolean"
		case "string":
			xsdType = "xsd:integer"
		case "dateTime":
			xsdType = "xsd:integer"
		}

		innerContext[field.Name] = map[string]interface{}{
			"@id":   "iden3-vocab:" + field.Name,
			"@type": xsdType,
		}

		if !request.IsMerklized {
			slotKey := string(field.Slot)
			serializationParts = append(serializationParts, slotKey+"="+field.Name)
		}
	}

	if len(serializationParts) > 0 {
		innerContext["iden3_serialization"] = "iden3:v1:" + strings.Join(serializationParts, "&")
	}

	return map[string]interface{}{
		"@context": []interface{}{
			map[string]interface{}{
				"@protected": true,
				"@version":   1.1,
				"id":         "@id",
				"type":       "@type",
				credType: map[string]interface{}{
					"@id":      credURN,
					"@context": innerContext,
				},
			},
		},
	}
}

func (s *SchemaService) uploadToIPFS(jsonSchema, jsonLdContext map[string]interface{}, schemaType string) (string, string, error) {
	jsonBytes, err := json.MarshalIndent(jsonSchema, "", "  ")
	if err != nil {
		return "", "", fmt.Errorf("marshal json schema failed: %w", err)
	}

	ldBytes, err := json.MarshalIndent(jsonLdContext, "", "  ")
	if err != nil {
		return "", "", fmt.Errorf("marshal jsonld context failed: %w", err)
	}

	namePrefix := strings.ToLower(strings.ReplaceAll(schemaType, " ", "-"))
	schemaURL, err := s.ipfs.Upload(namePrefix+".json", jsonBytes)
	if err != nil {
		return "", "", fmt.Errorf("upload json schema failed: %w", err)
	}

	contextURL, err := s.ipfs.Upload(namePrefix+".jsonld", ldBytes)
	if err != nil {
		return "", "", fmt.Errorf("upload jsonld context failed: %w", err)
	}

	return schemaURL, contextURL, nil
}

func (s *SchemaService) getSchemaHash(url string) (*core.SchemaHash, error) {
	var sHash core.SchemaHash
	h := keccak256.Hash([]byte(url))
	copy(sHash[:], h[len(h)-16:])
	sHashIndex, err := sHash.MarshalText()
	if err != nil {
		return nil, err
	}
	claim, _ := core.NewSchemaHashFromHex(string(sHashIndex))
	return &claim, nil
}

// func (s *SchemaService) ParseClaimFromSchema(ctx context.Context, schemaID string, credentialData map[string]interface{}) (*dto.ClaimDataDto, error) {
// 	schema, err := s.schemaRepo.FindSchemaByPublicId(ctx, schemaID)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, &constant.SchemaNotFound
// 		}
// 		return nil, &constant.InternalServer
// 	}

// 	claimData := &dto.ClaimDataDto{
// 		Type:              schema.Type,
// 		IsMerklized:       schema.IsMerklized,
// 		CredentialSubject: make(map[string]interface{}),
// 	}

// 	if schema.SchemaURL != "" {
// 		schemaHash, err := s.getSchemaHash(schema.SchemaURL)

// 		if err != nil {
// 			return nil, fmt.Errorf("failed to get schema hash: %w", err)
// 		}

// 		claimData.SchemaHash = schemaHash.BigInt().String()
// 	}

// 	for _, attr := range schema.SchemaAttributes {
// 		value, exists := credentialData[attr.Name]
// 		if attr.Required && !exists {
// 			return nil, fmt.Errorf("required field '%s' is missing", attr.Name)
// 		}

// 		if err := s.validateAttributeType(attr.Type, value); err != nil {
// 			return nil, fmt.Errorf("invalid type for field '%s': %w", attr.Name, err)
// 		}
// 		claimData.CredentialSubject[attr.Name] = value
// 	}

// 	if subjectID, ok := credentialData["id"].(string); ok {
// 		claimData.CredentialSubject["id"] = subjectID
// 	}
// 	if !schema.IsMerklized {
// 		claimData.SlotIndexMapping = s.buildSlotIndexMapping(schema.SchemaAttributes)
// 	}

// 	return claimData, nil
// }

// func (s *SchemaService) buildSlotIndexMapping(attributes []*schema.SchemaAttribute) map[string]string {
// 	mapping := make(map[string]string)
// 	for _, attr := range attributes {
// 		if attr.Slot != "" {
// 			mapping[attr.Slot] = attr.Name
// 		}
// 	}
// 	return mapping
// }

// func (s *SchemaService) validateAttributeType(expectedType string, value interface{}) error {
// 	switch expectedType {
// 	case "string":
// 		if _, ok := value.(string); !ok {
// 			return fmt.Errorf("expected string, got %T", value)
// 		}
// 	case "integer":
// 		switch value.(type) {
// 		case int, int32, int64, float64:
// 			return nil
// 		default:
// 			return fmt.Errorf("expected integer, got %T", value)
// 		}
// 	case "number":
// 		switch value.(type) {
// 		case float32, float64, int, int32, int64:
// 			return nil
// 		default:
// 			return fmt.Errorf("expected number, got %T", value)
// 		}
// 	case "boolean":
// 		if _, ok := value.(bool); !ok {
// 			return fmt.Errorf("expected boolean, got %T", value)
// 		}
// 	default:
// 		// Cho phép các type khác
// 		return nil
// 	}
// 	return nil
// }
