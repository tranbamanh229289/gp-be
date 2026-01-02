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
	"github.com/iden3/go-schema-processor/v2/processor"
	"github.com/iden3/go-schema-processor/v2/verifiable"
	"github.com/iden3/iden3comm/v2/protocol"
	"gorm.io/gorm"
)

type ICredentialService interface {
	GetCredentialRequests(ctx context.Context) ([]*dto.CredentialRequestResponseDto, error)
	CreateCredentialRequest(ctx context.Context, request *protocol.CredentialIssuanceRequestMessage) (*dto.CredentialRequestResponseDto, error)
	UpdateCredentialRequest(ctx context.Context, id string, request *dto.CredentialRequestUpdatedRequestDto) error
	GetVerifiableCredentials(ctx context.Context) ([]*dto.VerifiableCredentialResponseDto, error)
	GetVerifiableCredential(ctx context.Context, id string) (*dto.VerifiableCredentialResponseDto, error)
	IssueVerifiableCredential(ctx context.Context, id string, request *dto.IssueVerifiableCredentialRequestDto) (*dto.VerifiableCredentialResponseDto, error)
	UpdateVerifiableCredential(ctx context.Context, id string, request *dto.VerifiableUpdatedRequestDto) error
	SignVerifiableCredential(ctx context.Context, id string, request *dto.SignCredentialRequestDto) error
}

type CredentialService struct {
	config                *config.Config
	identityService       IIdentityService
	credentialRequestRepo credential.ICredentialRequestRepository
	vcRepo                credential.IVerifiableCredentialRepository
	schemaRepo            schema.ISchemaRepository
	processor             *processor.Processor
}

func NewCredentialService(
	config *config.Config,
	identityService IIdentityService,
	credentialRequestRepo credential.ICredentialRequestRepository,
	vcRepo credential.IVerifiableCredentialRepository,
	schemaRepo schema.ISchemaRepository,
) ICredentialService {
	return &CredentialService{
		config:                config,
		identityService:       identityService,
		credentialRequestRepo: credentialRequestRepo,
		vcRepo:                vcRepo,
		schemaRepo:            schemaRepo,
	}
}

func (s *CredentialService) CreateCredentialRequest(ctx context.Context, request *protocol.CredentialIssuanceRequestMessage) (*dto.CredentialRequestResponseDto, error) {
	schemaHash := request.Body.Schema.Hash
	schemaEntity, err := s.schemaRepo.FindSchemaByHash(ctx, schemaHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.SchemaNotFound
		}
		return nil, &constant.InternalServer
	}
	var claimData map[string]interface{}
	if err := json.Unmarshal(request.Body.Data, &claimData); err != nil {
		return nil, errors.New("unmarshal failed")
	}

	credentialRequestCreated, err := s.credentialRequestRepo.CreateCredentialRequest(ctx, &credential.CredentialRequest{
		PublicID:    uuid.New(),
		RequestID:   request.ID,
		HolderDID:   request.From,
		IssuerDID:   request.To,
		SchemaID:    schemaEntity.ID,
		SchemaHash:  schemaHash,
		Data:        claimData,
		Expiration:  request.Body.Expiration,
		CreatedTime: request.CreatedTime,
		ExpiresTime: request.ExpiresTime,
		Status:      constant.CredentialRequestPendingStatus,
	})

	return dto.ToCredentialRequestResponseDto(credentialRequestCreated), nil
}

func (s *CredentialService) GetCredentialRequests(ctx context.Context) ([]*dto.CredentialRequestResponseDto, error) {
	entities, err := s.credentialRequestRepo.FindAllCredentialRequest(ctx)
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

func (s *CredentialService) GetVerifiableCredentials(ctx context.Context) ([]*dto.VerifiableCredentialResponseDto, error) {
	entities, err := s.vcRepo.FindAllVerifiableCredential(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.VerifiableCredentialNotFound
		}
		return nil, &constant.InternalServer
	}

	var resp []*dto.VerifiableCredentialResponseDto
	for _, item := range entities {
		resp = append(resp, dto.ToVerifiableCredentialResponseDto(item))
	}
	return resp, nil
}

func (s *CredentialService) GetVerifiableCredential(ctx context.Context, id string) (*dto.VerifiableCredentialResponseDto, error) {
	vc, err := s.vcRepo.FindVerifiableCredentialByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.VerifiableCredentialNotFound
		}
		return nil, &constant.InternalServer
	}
	return dto.ToVerifiableCredentialResponseDto(vc), nil
}

func (s *CredentialService) IssueVerifiableCredential(ctx context.Context, id string, request *dto.IssueVerifiableCredentialRequestDto) (*dto.VerifiableCredentialResponseDto, error) {
	credentialRequestEntity, err := s.credentialRequestRepo.FindCredentialRequestByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.CredentialRequestNotFound
		}
		return nil, &constant.InternalServer
	}

	// identity state
	identityState, err := s.identityService.GetIdentityStateByDID(ctx, credentialRequestEntity.IssuerDID)
	if err != nil {
		return nil, err
	}

	credentialSubjectID, _ := credentialRequestEntity.Data["id"]

	// VC
	issuanceDate := time.Now().UTC()
	expirationDate := time.Unix(credentialRequestEntity.Expiration/1000, 0).UTC()

	verifiableCredential := &verifiable.W3CCredential{
		ID: "urn:uuid:" + uuid.New().String(),
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
			credentialRequestEntity.Schema.ContextURL,
		},
		Type: []string{
			"VerifiableCredential",
			credentialRequestEntity.Schema.Type,
		},
		IssuanceDate: &issuanceDate,
		Expiration:   &expirationDate,
		Issuer:       credentialRequestEntity.IssuerDID,
		CredentialSchema: verifiable.CredentialSchema{
			ID:   credentialRequestEntity.Schema.SchemaURL,
			Type: credentialRequestEntity.Schema.Type,
		},
		CredentialSubject: credentialRequestEntity.Data,
	}

	subjectPosition := verifiable.CredentialSubjectPositionIndex
	merklizedRootPosition := verifiable.CredentialMerklizedRootPositionNone

	if credentialRequestEntity.Schema.IsMerklized {
		merklizedRootPosition = verifiable.CredentialMerklizedRootPositionIndex
	}

	revNonce := uint64(time.Now().Unix())
	coreClaimOption := verifiable.CoreClaimOptions{
		RevNonce:              revNonce,
		Version:               1,
		SubjectPosition:       subjectPosition,
		MerklizedRootPosition: merklizedRootPosition,
		Updatable:             true,
	}

	coreClaim, err := verifiableCredential.ToCoreClaim(ctx, &coreClaimOption)

	// claim
	hi, hv, err := coreClaim.HiHv()
	if err != nil {
		return nil, fmt.Errorf("failed to get HiHv: %w", err)
	}

	stateHash, err := identityState.GetStateValue()
	if err != nil {
		return nil, fmt.Errorf("failed to get state: %w", err)
	}
	coreClaimHex, err := coreClaim.Hex()
	if err != nil {
		return nil, fmt.Errorf("failed to get core claim hex %w", err)
	}
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

	claimTreeRoot := identityState.ClaimsTree.Root().Hex()
	revRoot := identityState.RevTree.Root().Hex()
	rootsRoot := identityState.RootsTree.Root().Hex()
	hashString := stateHash.Hex()

	status := constant.VerifiableCredentialIssuedStatus
	var proof verifiable.CredentialProof

	// mtp
	if request.ProofType == verifiable.Iden3SparseMerkleTreeProofType {
		err = identityState.AddClaim(ctx, coreClaim)
		if err != nil {
			return nil, fmt.Errorf("failed to add claim: %w", err)
		}
		incProof, err := identityState.GetIncMTProof(ctx, coreClaim)
		if err != nil {
			return nil, fmt.Errorf("failed to generate inclusion proof: %w", err)
		}

		proof = &verifiable.Iden3SparseMerkleTreeProof{
			Type: verifiable.Iden3SparseMerkleTreeProofType,
			IssuerData: verifiable.IssuerData{
				ID: identityState.GetDID().String(),
				State: verifiable.State{
					Value:              &hashString,
					ClaimsTreeRoot:     &claimTreeRoot,
					RevocationTreeRoot: &revRoot,
					RootOfRoots:        &rootsRoot,
					Status:             constant.CredentialRequestApprovedStatus,
				},
				AuthCoreClaim:    authClaimHex,
				MTP:              authIncProof,
				CredentialStatus: constant.CredentialRequestApprovedStatus,
			},
			CoreClaim: coreClaimHex,
			MTP:       incProof,
		}
	} else if request.ProofType == verifiable.BJJSignatureProofType {
		proof = &verifiable.BJJSignatureProof2021{
			Type: verifiable.BJJSignatureProofType,
			IssuerData: verifiable.IssuerData{
				ID: identityState.GetDID().String(),
				State: verifiable.State{
					Value:              &hashString,
					ClaimsTreeRoot:     &claimTreeRoot,
					RevocationTreeRoot: &revRoot,
					RootOfRoots:        &rootsRoot,
					Status:             constant.CredentialRequestApprovedStatus,
				},
				AuthCoreClaim:    authClaimHex,
				MTP:              authIncProof,
				CredentialStatus: constant.CredentialRequestApprovedStatus,
			},
			CoreClaim: coreClaimHex,
			Signature: "",
		}
		status = constant.VerifiableCredentialNotSignedStatus
	}

	verifiableCredential.Proof = []verifiable.CredentialProof{proof}
	vcCreated, err := s.vcRepo.CreateVerifiableCredential(ctx, &credential.VerifiableCredential{
		PublicID:     uuid.New(),
		IssuerDID:    credentialRequestEntity.IssuerDID,
		HolderDID:    credentialRequestEntity.HolderDID,
		SchemaID:     credentialRequestEntity.SchemaID,
		ClaimHi:      hi.String(),
		ClaimHv:      hv.String(),
		ClaimSubject: credentialSubjectID.(string),
		RevNonce:     revNonce,
		ProofType:    string(request.ProofType),
		Signature:    "",
		Status:       status,
	})
	if err != nil {
		return nil, err
	}

	return dto.ToVerifiableCredentialResponseDto(vcCreated), nil
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

func (s *CredentialService) SignVerifiableCredential(ctx context.Context, id string, request *dto.SignCredentialRequestDto) error {
	vc, err := s.vcRepo.FindVerifiableCredentialByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &constant.VerifiableCredentialNotFound
		}
		return &constant.InternalServer
	}

	if vc.ProofType != string(verifiable.BJJSignatureProofType) {
		return &constant.VerifiableCredentialNotSig
	}

	if vc.Signature != "" {
		return &constant.VerifiableCredentialAlreadySig
	}

	if vc.Status != constant.VerifiableCredentialNotSignedStatus {
		return &constant.InternalServer
	}

	changes := map[string]interface{}{"signature": request.Signature, "status": constant.VerifiableCredentialIssuedStatus}
	return s.vcRepo.UpdateVerifiableCredential(ctx, vc, changes)
}

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
