package service

import (
	"be/config"
	"be/internal/domain/credential"
	"be/internal/domain/schema"
	"be/internal/infrastructure/ipfs"
	"be/internal/shared/constant"
	"be/internal/transport/http/dto"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ISchemaService interface {
	CreateSchema(ctx context.Context, request *dto.SchemaCreatedRequestDto) (*dto.SchemaResponseDto, error)
	RemoveSchema(ctx context.Context, id string) error
}
type SchemaService struct {
	config       *config.Config
	ipfs         *ipfs.Pinata
	identityRepo credential.IIdentityRepository
	schemaRepo   schema.ISchemaRepository
}

func NewSchemaService(config *config.Config, pinata *ipfs.Pinata, schemaRepo schema.ISchemaRepository, identityRepo credential.IIdentityRepository) ISchemaService {
	return &SchemaService{
		config:       config,
		ipfs:         pinata,
		schemaRepo:   schemaRepo,
		identityRepo: identityRepo,
	}
}

func (s *SchemaService) CreateSchema(ctx context.Context, request *dto.SchemaCreatedRequestDto) (*dto.SchemaResponseDto, error) {
	identity, err := s.identityRepo.FindIdentityByPublicId(ctx, request.IssuerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.SchemaNotFound
		}
		return nil, &constant.InternalServer
	}

	if err := s.validate(request); err != nil {
		return nil, err
	}

	jsonSchema := s.generateJSONSchema(request)
	jsonLdContext := s.generateJSONLDContext(request)

	schema := &schema.Schema{
		PublicID:      uuid.New(),
		Type:          request.Type,
		Version:       request.Version,
		Status:        constant.SchemaActiveStatus,
		JSONSchema:    jsonSchema,
		JSONLDContext: jsonLdContext,
	}

	if request.IsPublishIPFS {
		jsonSchemaBytes, _ := json.MarshalIndent(jsonSchema, "", "    ")
		shemaFileName := strings.ToLower(request.Type) + ".json"
		jsonLdContextBytes, _ := json.MarshalIndent(jsonLdContext, "", "    ")
		ldContextFileName := strings.ToLower(request.Type) + ".jsonld"

		jsonCID, _ := s.ipfs.Upload(shemaFileName, jsonSchemaBytes)
		schema.JSONCID = jsonCID

		jsonLdCID, _ := s.ipfs.Upload(ldContextFileName, jsonLdContextBytes)
		schema.JSONLDCID = jsonLdCID
	}

	schemaCreated, err := s.schemaRepo.CreateSchema(ctx, schema)
	if err != nil {
		return nil, &constant.SchemaNotFound
	}

	return &dto.SchemaResponseDto{
		PublicID:  schemaCreated.PublicID.String(),
		IssuerID:  identity.PublicID.String(),
		Type:      schemaCreated.Type,
		Version:   schemaCreated.Version,
		JSONCID:   schema.JSONCID,
		JSONLDCID: schema.JSONLDCID,
		Status:    schema.Status,
	}, nil
}

func (s *SchemaService) RemoveSchema(ctx context.Context, id string) error {
	schema, err := s.schemaRepo.FindSchemaByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &constant.SchemaNotFound
		}
		return &constant.InternalServer
	}
	s.ipfs.Remove(schema.JSONCID)
	s.ipfs.Remove(schema.JSONLDCID)

	changes := map[string]interface{}{"status": constant.SchemaRevokeStatus, "revoked_at": time.Now()}
	return s.schemaRepo.UpdateSchema(ctx, schema, changes)
}

func (s *SchemaService) validate(request *dto.SchemaCreatedRequestDto) error {
	if !request.IsMerklized && len(request.Fields) != 4 {
		return fmt.Errorf("non merklized max 4 fields")
	}
	return nil
}

func (s *SchemaService) generateJSONSchema(request *dto.SchemaCreatedRequestDto) map[string]interface{} {
	if request.Type == "" {
		request.Type = "Credential"
	}
	if request.Version == "" {
		request.Version = "1.0"
	}

	baseURL := s.config.IPFS.BaseURL
	if baseURL == "" {
		baseURL = "https://example.com/schemas/"
	}

	URL := baseURL + strings.ToLower(request.Type)
	jsonLdURL := URL + ".jsonld"
	jsonSchemaURL := URL + ".json"

	iden3Serialization := map[string]string{}
	credSubProperties := map[string]interface{}{
		"id": map[string]interface{}{
			"type":        "string",
			"format":      "uri",
			"title":       "Credential subject ID",
			"description": "Stores the DID of the subject that owns the credential",
		},
	}
	credSubRequired := []string{"id"}

	for _, field := range request.Fields {
		credSubProperties[field.Name] = map[string]interface{}{
			"type":        field.Type,
			"title":       field.Title,
			"description": field.Description,
		}
		if field.Required {
			credSubRequired = append(credSubRequired, field.Name)
		}
		if field.Slot != "" {
			slotKey := ""
			switch field.Slot {
			case "IndexDataSlotA":
				slotKey = "slotIndexA"
			case "IndexDataSlotB":
				slotKey = "slotIndexB"
			case "ValueDataSlotA":
				slotKey = "slotValueA"
			case "ValueDataSlotB":
				slotKey = "slotValueB"
			}
			if slotKey != "" {
				iden3Serialization[slotKey] = field.Name
			}
		}
	}

	schema := map[string]interface{}{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"$metadata": map[string]interface{}{
			"uris": map[string]string{
				"jsonLdContext": jsonLdURL,
				"jsonSchema":    jsonSchemaURL,
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
			"@context": map[string]interface{}{
				"type": []string{"string", "array", "object"},
			},
			"id": map[string]interface{}{
				"type": "string",
			},
			"type": map[string]interface{}{
				"type": []string{"string", "array"},
				"items": map[string]interface{}{
					"type": "string",
				},
			},
			"issuer": map[string]interface{}{
				"type":   []string{"string", "object"},
				"format": "uri",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":   "string",
						"format": "uri",
					},
				},
				"required": []string{"id"},
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

func (s *SchemaService) generateJSONLDContext(request *dto.SchemaCreatedRequestDto) map[string]interface{} {
	if request.Type == "" {
		request.Type = "Credential"
	}

	credUUID := uuid.New()
	vocabUUID := uuid.New()

	credURN := "urn:uuid:" + credUUID.String()
	vocabURN := "urn:uuid:" + vocabUUID.String() + "#"

	serializationParts := []string{}
	fields := map[string]interface{}{}

	for _, field := range request.Fields {
		xsdType := "xsd:string"
		switch field.Type {
		case "integer":
			xsdType = "xsd:integer"
		case "number":
			xsdType = "xsd:double"
		case "boolean":
			xsdType = "xsd:boolean"
		}

		fields[field.Name] = map[string]interface{}{
			"@id":   "iden3-vocab:" + field.Name,
			"@type": xsdType,
		}

		if field.Slot != "" {
			slotKey := ""
			switch field.Slot {
			case "IndexDataSlotA":
				slotKey = "slotIndexA"
			case "IndexDataSlotB":
				slotKey = "slotIndexB"
			case "ValueDataSlotA":
				slotKey = "slotValueA"
			case "ValueDataSlotB":
				slotKey = "slotValueB"
			}
			if slotKey != "" {
				serializationParts = append(serializationParts, slotKey+"="+field.Name)
			}
		}
	}

	innerContext := map[string]interface{}{
		"@propagate":  true,
		"@protected":  true,
		"iden3-vocab": vocabURN,
		"xsd":         "http://www.w3.org/2001/XMLSchema#",
	}

	for k, v := range fields {
		innerContext[k] = v
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
				request.Type: map[string]interface{}{
					"@id":      credURN,
					"@context": innerContext,
				},
			},
		},
	}
}
