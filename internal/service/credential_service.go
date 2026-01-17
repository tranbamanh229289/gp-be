package service

import (
	"be/config"
	"be/internal/domain/credential"
	"be/internal/domain/schema"
	"be/internal/shared/constant"
	"be/internal/transport/http/dto"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iden3/go-schema-processor/v2/loaders"
	"github.com/iden3/go-schema-processor/v2/merklize"
	"github.com/iden3/go-schema-processor/v2/processor"
	"github.com/iden3/go-schema-processor/v2/verifiable"
	"github.com/iden3/iden3comm/v2/protocol"
	"gorm.io/gorm"
)

type ICredentialService interface {
	GetCredentialRequests(ctx context.Context, claims *dto.ZKClaims) ([]*dto.CredentialRequestResponseDto, error)
	CreateCredentialRequest(ctx context.Context, request *protocol.CredentialIssuanceRequestMessage, claims *dto.ZKClaims) (*dto.CredentialRequestResponseDto, error)
	UpdateCredentialRequest(ctx context.Context, id string, request *dto.CredentialRequestUpdatedRequestDto) error
	GetVerifiableCredentials(ctx context.Context, claims *dto.ZKClaims) ([]*verifiable.W3CCredential, error)
	GetVerifiableCredentialById(ctx context.Context, id string) (*verifiable.W3CCredential, error)
	IssueVerifiableCredential(ctx context.Context, id string, request *dto.IssueVerifiableCredentialRequestDto, claims *dto.ZKClaims) (*verifiable.W3CCredential, error)
	UpdateVerifiableCredential(ctx context.Context, id string, request *dto.VerifiableUpdatedRequestDto) error
}

type CredentialService struct {
	config                *config.Config
	identityService       IIdentityService
	documentService       IDocumentService
	credentialRequestRepo credential.ICredentialRequestRepository
	vcRepo                credential.IVerifiableCredentialRepository
	schemaRepo            schema.ISchemaRepository
	processor             *processor.Processor
}

func NewCredentialService(
	config *config.Config,
	identityService IIdentityService,
	documentService IDocumentService,
	credentialRequestRepo credential.ICredentialRequestRepository,
	vcRepo credential.IVerifiableCredentialRepository,
	schemaRepo schema.ISchemaRepository,
) ICredentialService {
	return &CredentialService{
		config:                config,
		identityService:       identityService,
		documentService:       documentService,
		credentialRequestRepo: credentialRequestRepo,
		vcRepo:                vcRepo,
		schemaRepo:            schemaRepo,
	}
}

func (s *CredentialService) CreateCredentialRequest(ctx context.Context, request *protocol.CredentialIssuanceRequestMessage, claims *dto.ZKClaims) (*dto.CredentialRequestResponseDto, error) {
	if request.From != claims.DID {
		return nil, &constant.BadRequest
	}

	schemaHash := request.Body.Schema.Hash
	schemaEntity, err := s.schemaRepo.FindSchemaByHash(ctx, schemaHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.SchemaNotFound
		}
		return nil, &constant.InternalServer
	}

	holder, err := s.identityService.GetIdentityByDID(ctx, request.From)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.IdentityNotFound
		}
		return nil, &constant.InternalServer
	}

	issuer, err := s.identityService.GetIdentityByDID(ctx, request.To)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.IdentityNotFound
		}
		return nil, &constant.InternalServer
	}

	var claimData map[string]interface{}
	if err := json.Unmarshal(request.Body.Data, &claimData); err != nil {
		return nil, errors.New("unmarshal failed")
	}

	credentialRequestCreated, err := s.credentialRequestRepo.CreateCredentialRequest(ctx, &credential.CredentialRequest{
		PublicID:    uuid.New(),
		ThreadID:    request.ThreadID,
		HolderDID:   request.From,
		IssuerDID:   request.To,
		SchemaID:    schemaEntity.ID,
		SchemaHash:  schemaHash,
		Expiration:  request.Body.Expiration,
		CreatedTime: request.CreatedTime,
		ExpiresTime: request.ExpiresTime,
		Status:      constant.CredentialRequestPendingStatus,
	})

	return &dto.CredentialRequestResponseDto{
		PublicID:     credentialRequestCreated.PublicID.String(),
		ThreadID:     credentialRequestCreated.ThreadID,
		HolderDID:    holder.DID,
		HolderName:   holder.Name,
		IssuerDID:    issuer.DID,
		IssuerName:   issuer.Name,
		SchemaID:     schemaEntity.PublicID.String(),
		SchemaTitle:  schemaEntity.Title,
		SchemaURL:    schemaEntity.SchemaURL,
		ContextURL:   schemaEntity.ContextURL,
		SchemaType:   schemaEntity.Type,
		SchemaHash:   schemaEntity.Hash,
		IsMerklized:  schemaEntity.IsMerklized,
		DocumentType: schemaEntity.DocumentType,
		Status:       credentialRequestCreated.Status,
		Expiration:   credentialRequestCreated.Expiration,
		CreatedTime:  credentialRequestCreated.CreatedTime,
		ExpiresTime:  credentialRequestCreated.ExpiresTime,
	}, nil
}

func (s *CredentialService) GetCredentialRequests(ctx context.Context, claims *dto.ZKClaims) ([]*dto.CredentialRequestResponseDto, error) {
	var (
		entities []*credential.CredentialRequest
		err      error
	)
	switch claims.Role {
	case constant.IdentityHolderRole:
		entities, err = s.credentialRequestRepo.FindAllCredentialRequestsByHolderDID(ctx, claims.DID)
	case constant.IdentityIssuerRole:
		entities, err = s.credentialRequestRepo.FindAllCredentialRequestsByIssuerDID(ctx, claims.DID)
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.CredentialRequestNotFound
		}
		return nil, &constant.InternalServer
	}

	var resp []*dto.CredentialRequestResponseDto
	for _, item := range entities {
		resp = append(resp, dto.ToCredentialRequestResponseDto(item))
	}
	return resp, nil
}

func (s *CredentialService) UpdateCredentialRequest(ctx context.Context, id string, request *dto.CredentialRequestUpdatedRequestDto) error {
	entity, err := s.credentialRequestRepo.FindCredentialRequestByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &constant.CredentialRequestNotFound
		}
		return &constant.InternalServer
	}
	changes := map[string]interface{}{"status": request.Status}

	err = s.credentialRequestRepo.UpdateCredentialRequest(ctx, entity, changes)

	return err
}

func (s *CredentialService) GetVerifiableCredentials(ctx context.Context, claims *dto.ZKClaims) ([]*verifiable.W3CCredential, error) {
	var (
		entities []*credential.VerifiableCredential
		err      error
	)
	switch claims.Role {
	case constant.IdentityHolderRole:
		entities, err = s.vcRepo.FindAllVerifiableCredentialsByHolderDID(ctx, claims.DID)
	case constant.IdentityIssuerRole:
		entities, err = s.vcRepo.FindAllVerifiableCredentialsByIssuerDID(ctx, claims.DID)
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.VerifiableCredentialNotFound
		}
		return nil, &constant.InternalServer
	}

	var vcs []*verifiable.W3CCredential
	for _, item := range entities {
		vcs = append(vcs, dto.ToW3CCredential(item))
	}
	return vcs, nil
	// var resp []*dto.VerifiableCredentialResponseDto
	// for _, item := range entities {
	// 	resp = append(resp, dto.ToVerifiableCredentialResponseDto(item))
	// }
	// return resp, nil
}

func (s *CredentialService) GetVerifiableCredentialById(ctx context.Context, id string) (*verifiable.W3CCredential, error) {
	vc, err := s.vcRepo.FindVerifiableCredentialByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.VerifiableCredentialNotFound
		}
		return nil, &constant.InternalServer
	}
	return dto.ToW3CCredential(vc), nil
}

func (s *CredentialService) IssueVerifiableCredential(ctx context.Context, id string, request *dto.IssueVerifiableCredentialRequestDto, claims *dto.ZKClaims) (*verifiable.W3CCredential, error) {
	credentialRequestEntity, err := s.credentialRequestRepo.FindCredentialRequestByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.CredentialRequestNotFound
		}
		return nil, &constant.InternalServer
	}

	if credentialRequestEntity.IssuerDID != claims.DID {
		return nil, &constant.BadRequest
	}

	issuanceDate := time.Now().UTC()
	expirationDate := time.Unix(credentialRequestEntity.Expiration/1000, 0).UTC()

	verifiableCredential := &verifiable.W3CCredential{
		ID: "urn:uuid:" + uuid.New().String(),
		Context: []string{
			verifiable.JSONLDSchemaW3CCredential2018,
			verifiable.JSONLDSchemaIden3Credential,
			credentialRequestEntity.Schema.ContextURL,
		},
		Type: []string{
			verifiable.TypeW3CVerifiableCredential,
			credentialRequestEntity.Schema.Type,
		},
		IssuanceDate: &issuanceDate,
		Expiration:   &expirationDate,

		Issuer: credentialRequestEntity.IssuerDID,
		CredentialSchema: verifiable.CredentialSchema{
			ID:   credentialRequestEntity.Schema.SchemaURL,
			Type: verifiable.JSONSchema2023,
		},
		CredentialStatus:  request.CredentialStatus,
		CredentialSubject: request.CredentialSubject,
	}

	options := &verifiable.CoreClaimOptions{
		RevNonce:              request.CredentialStatus.RevocationNonce,
		Version:               0,
		SubjectPosition:       verifiable.CredentialSubjectPositionIndex,
		MerklizedRootPosition: verifiable.CredentialMerklizedRootPositionNone,
		Updatable:             false,
	}

	loader := loaders.NewDocumentLoader(nil, "")
	options.MerklizerOpts = []merklize.MerklizeOption{merklize.WithDocumentLoader(loader)}

	coreClaim, err := verifiableCredential.ToCoreClaim(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to create claim %w", err)
	}

	hi, hv, err := coreClaim.HiHv()
	if err != nil {
		return nil, fmt.Errorf("failed to get HiHv: %w", err)
	}

	identityState, err := s.identityService.GetIdentityStateByDID(ctx, credentialRequestEntity.IssuerDID)
	if err != nil {
		return nil, fmt.Errorf("failed to get identity state: %w", err)
	}

	claimSubject, err := coreClaim.GetID()
	if err != nil {
		return nil, fmt.Errorf("failed to get credentialSubject")
	}

	// mtp
	err = identityState.AddClaim(ctx, coreClaim)
	if err != nil {
		return nil, fmt.Errorf("failed to add claim: %w", err)
	}

	coreClaimHex, err := coreClaim.Hex()
	if err != nil {
		return nil, fmt.Errorf("failed to get core claim hex %w", err)
	}

	incProof, err := identityState.GetIncMTProof(ctx, coreClaim)
	if err != nil {
		return nil, fmt.Errorf("failed to generate inclusion proof: %w", err)
	}

	// sig
	authClaim, err := identityState.GetAuthClaim()
	if err != nil {
		return nil, fmt.Errorf("failed to get auth claim %w", err)
	}

	authClaimHex, err := authClaim.Hex()
	if err != nil {
		return nil, fmt.Errorf("failed to get auth claim hex %w", err)
	}

	authIncProof, err := identityState.GetIncMTProof(ctx, authClaim)
	if err != nil {
		return nil, fmt.Errorf("failed to generate inclusion proof: %w", err)
	}

	// state
	stateHash, err := identityState.GetStateValue()
	if err != nil {
		return nil, fmt.Errorf("failed to get state: %w", err)
	}
	hashString := stateHash.Hex()

	// root
	claimsTreeRoot := identityState.ClaimsTree.Root().Hex()
	revTreeRoot := identityState.RevTree.Root().Hex()
	rootsTreeRoot := identityState.RootsTree.Root().Hex()

	iden3SparseMerkleProof := &verifiable.Iden3SparseMerkleTreeProof{
		Type: verifiable.Iden3SparseMerkleTreeProofType,
		IssuerData: verifiable.IssuerData{
			ID: identityState.GetDID().String(),
			State: verifiable.State{
				Value:              &hashString,
				ClaimsTreeRoot:     &claimsTreeRoot,
				RevocationTreeRoot: &revTreeRoot,
				RootOfRoots:        &rootsTreeRoot,
				Status:             string(constant.VerifiableCredentialIssuedStatus),
			},
			AuthCoreClaim:    authClaimHex,
			MTP:              authIncProof,
			CredentialStatus: request.CredentialStatus,
		},
		CoreClaim: coreClaimHex,
		MTP:       incProof,
	}

	bjjSignatureProof := &verifiable.BJJSignatureProof2021{
		Type: verifiable.BJJSignatureProofType,
		IssuerData: verifiable.IssuerData{
			ID: identityState.GetDID().String(),
			State: verifiable.State{
				Value:              &hashString,
				ClaimsTreeRoot:     &claimsTreeRoot,
				RevocationTreeRoot: &revTreeRoot,
				RootOfRoots:        &rootsTreeRoot,
				Status:             string(constant.VerifiableCredentialIssuedStatus),
			},
			AuthCoreClaim:    authClaimHex,
			MTP:              authIncProof,
			CredentialStatus: request.CredentialStatus,
		},
		CoreClaim: coreClaimHex,
		Signature: request.Signature,
	}

	verifiableCredential.Proof = []verifiable.CredentialProof{iden3SparseMerkleProof, bjjSignatureProof}

	incProofJSON, err := incProof.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal core claim proof: %w", err)
	}
	authProofJSON, err := authIncProof.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal auth claim proof: %w", err)
	}
	_, err = s.vcRepo.CreateVerifiableCredential(ctx, &credential.VerifiableCredential{
		PublicID:          uuid.New(),
		CRID:              credentialRequestEntity.ID,
		HolderDID:         credentialRequestEntity.HolderDID,
		IssuerDID:         credentialRequestEntity.IssuerDID,
		SchemaID:          credentialRequestEntity.SchemaID,
		SchemaHash:        credentialRequestEntity.SchemaHash,
		CredentialID:      verifiableCredential.ID,
		CredentialSubject: request.CredentialSubject,
		ClaimHi:           hi.String(),
		ClaimHv:           hv.String(),
		ClaimHex:          coreClaimHex,
		ClaimSubject:      claimSubject.String(),
		ClaimMTP:          incProofJSON,
		RevNonce:          request.CredentialStatus.RevocationNonce,
		AuthClaimHex:      authClaimHex,
		AuthClaimMTP:      authProofJSON,
		IssuerState:       hashString,
		ClaimsTreeRoot:    claimsTreeRoot,
		RevTreeRoot:       revTreeRoot,
		RootsTreeRoot:     rootsTreeRoot,
		Status:            constant.VerifiableCredentialIssuedStatus,
		IssuanceDate:      &issuanceDate,
		ExpirationDate:    &expirationDate,
		Signature:         request.Signature,
	})
	if err != nil {
		return nil, err
	}

	return verifiableCredential, nil
}

func (s *CredentialService) UpdateVerifiableCredential(ctx context.Context, id string, request *dto.VerifiableUpdatedRequestDto) error {
	vc, err := s.vcRepo.FindVerifiableCredentialByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &constant.VerifiableCredentialNotFound
		}
		return &constant.InternalServer
	}

	changes := map[string]interface{}{"status": request.Status}
	return s.vcRepo.UpdateVerifiableCredential(ctx, vc, changes)
}

// func (s *CredentialService) getCredentialSubject(ctx context.Context, credentialRequest *credential.CredentialRequest) (map[string]interface{}, error) {
// 	var credentialSubject = make(map[string]interface{})
// 	var (
// 		info interface{}
// 		err  error
// 	)

// 	switch credentialRequest.Schema.DocumentType {
// 	case constant.CitizenIdentity:
// 		info, err = s.documentService.GetCitizenIdentityByHolderDID(ctx, credentialRequest.HolderDID)

// 	case constant.AcademicDegree:
// 		info, err = s.documentService.GetAcademicDegreeByHolderDID(ctx, credentialRequest.HolderDID)

// 	case constant.HealthInsurance:
// 		info, err = s.documentService.GetHealthInsuranceByHolderDID(ctx, credentialRequest.HolderDID)

// 	case constant.DriverLicense:
// 		info, err = s.documentService.GetDriverLicenseByHolderDID(ctx, credentialRequest.HolderDID)

// 	case constant.Passport:
// 		info, err = s.documentService.GetPassportByHolderDID(ctx, credentialRequest.HolderDID)

// 	default:
// 		return nil, errors.New("unsupported document type")

// 	}

// 	if err != nil || info != nil {
// 		return nil, err
// 	}

// 	infoJSON, err := response.StructToMap(info)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, item := range credentialRequest.Schema.SchemaAttributes {
// 		val, ok := infoJSON[item.Name]
// 		if !ok {
// 			return nil, fmt.Errorf("missing field in document: %s", item.Name)
// 		}
// 		credentialSubject[item.Name] = val
// 	}

// 	return credentialSubject, nil
// }

// func (s *CredentialService) createCoreClaim(schema *schema.Schema, data datatypes.JSONMap) (*core.Claim, error) {
// 	schemaHash, err := core.NewSchemaHashFromHex(schema.Hash)
// 	if err != nil {
// 		return nil, err
// 	}
// 	revNonce := uint64(time.Now().Unix())

// 	var credentialSubject map[string]interface{} = data
// 	id, ok := data["id"]

// 	if schema.IsMerklized {
// 		root, _ := s.calculateMerklized(credentialSubject)
// 		if ok {
// 			subjectID, _ := core.IDFromString(id.(string))
// 			return core.NewClaim(schemaHash, core.WithRevocationNonce(uint64(revNonce)), core.WithIndexID(subjectID), core.WithIndexMerklizedRoot(root))
// 		} else {
// 			return core.NewClaim(schemaHash, core.WithRevocationNonce(uint64(revNonce)), core.WithIndexMerklizedRoot(root))
// 		}
// 	} else {
// 		slotValue := make(map[string]*big.Int, 4)
// 		for _, attr := range schema.SchemaAttributes {
// 			slotValue[attr.Slot] = s.convertToBigInt(data[attr.Name])
// 		}
// 		if ok {
// 			subjectID, _ := core.IDFromString(id.(string))
// 			return core.NewClaim(schemaHash, core.WithRevocationNonce(uint64(revNonce)), core.WithIndexID(subjectID),
// 				core.WithIndexDataInts(slotValue[string(constant.SlotIndexA)], slotValue[string(constant.SlotIndexA)]),
// 				core.WithValueDataInts(slotValue[string(constant.SlotDataA)], slotValue[string(constant.SlotDataB)]))

// 		} else {
// 			return core.NewClaim(schemaHash, core.WithRevocationNonce(uint64(revNonce)),
// 				core.WithIndexDataInts(slotValue[string(constant.SlotIndexA)], slotValue[string(constant.SlotIndexA)]),
// 				core.WithValueDataInts(slotValue[string(constant.SlotDataA)], slotValue[string(constant.SlotDataB)]))
// 		}
// 	}

// }

// func (s *CredentialService) convertToBigInt(val interface{}) *big.Int {
// 	switch v := val.(type) {
// 	case int:
// 		return big.NewInt(int64(v))
// 	case uint:
// 		return big.NewInt(int64(v))
// 	case float64:
// 		return big.NewInt(int64(v))
// 	case string:
// 		bi := new(big.Int)
// 		bi.SetString(v, 10)
// 		return bi

// 	case bool:
// 		if v {
// 			return big.NewInt(1)
// 		} else {
// 			return big.NewInt(0)
// 		}
// 	default:
// 		return big.NewInt(0)
// 	}
// }

// func (s *CredentialService) calculateMerklized(credentialSubject map[string]interface{}) (*big.Int, error) {
// 	doc := map[string]interface{}{
// 		"@context":          []interface{}{"https://www.w3.org/2018/credentials/v1"},
// 		"credentialSubject": credentialSubject,
// 	}

// 	docBytes, err := json.Marshal(doc)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to marshal doc: %w", err)
// 	}

// 	mk, err := merklize.MerklizeJSONLD(context.Background(), bytes.NewReader(docBytes))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to merklize: %w", err)
// 	}
// 	root := mk.Root()
// 	return root.BigInt(), nil
// }
