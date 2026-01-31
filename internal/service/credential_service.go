package service

import (
	"be/config"
	"be/internal/domain/credential"
	"be/internal/domain/schema"
	"be/internal/shared/constant"
	"be/internal/shared/helper"
	"be/internal/transport/http/dto"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iden3/go-schema-processor/v2/merklize"
	"github.com/iden3/go-schema-processor/v2/verifiable"
	"github.com/iden3/iden3comm/v2/protocol"
	"github.com/piprate/json-gold/ld"
	"gorm.io/gorm"
)

type ICredentialService interface {
	GetCredentialRequests(ctx context.Context, claims *dto.ZKClaims) ([]*dto.CredentialRequestResponseDto, error)
	CreateCredentialRequest(ctx context.Context, request *protocol.CredentialIssuanceRequestMessage) (*dto.CredentialRequestResponseDto, error)
	UpdateCredentialRequest(ctx context.Context, id string, request *dto.CredentialRequestUpdatedRequestDto) error
	GetVerifiableCredentials(ctx context.Context, claims *dto.ZKClaims) ([]*verifiable.W3CCredential, error)
	GetVerifiableCredentialById(ctx context.Context, id string) (*verifiable.W3CCredential, error)
	IssueVerifiableCredential(ctx context.Context, id string, request *dto.IssueVerifiableCredentialRequestDto) (*verifiable.W3CCredential, error)
	UpdateVerifiableCredential(ctx context.Context, id string, request *dto.VerifiableUpdatedRequestDto) error
}

type CredentialService struct {
	config                *config.Config
	identityService       IIdentityService
	documentService       IDocumentService
	credentialRequestRepo credential.ICredentialRequestRepository
	vcRepo                credential.IVerifiableCredentialRepository
	schemaRepo            schema.ISchemaRepository
	loader                ld.DocumentLoader
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
		loader:                helper.NewCacheLoader(nil),
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

func (s *CredentialService) IssueVerifiableCredential(ctx context.Context, id string, request *dto.IssueVerifiableCredentialRequestDto) (*verifiable.W3CCredential, error) {
	credentialRequestEntity, err := s.credentialRequestRepo.FindCredentialRequestByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.CredentialRequestNotFound
		}
		return nil, &constant.InternalServer
	}

	issuanceDate := time.Now().UTC()
	expirationDate := time.Unix(credentialRequestEntity.Expiration, 0).UTC()

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
	// documentLoader := ld.NewDefaultDocumentLoader(nil)
	options.MerklizerOpts = []merklize.MerklizeOption{merklize.WithDocumentLoader(s.loader)}

	coreClaim, err := verifiableCredential.ToCoreClaim(ctx, options)
	if err != nil {

		return nil, fmt.Errorf("failed to create claim %w", err)
	}

	hi, hv, err := coreClaim.HiHv()
	if err != nil {
		return nil, fmt.Errorf("failed to get HiHv: %w", err)
	}

	fmt.Println("hi", hi)
	fmt.Println("hv", hv)

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

	changes := map[string]interface{}{"status": constant.CredentialRequestApprovedStatus}
	err = s.credentialRequestRepo.UpdateCredentialRequest(ctx, credentialRequestEntity, changes)
	if err != nil {
		return nil, &constant.InternalServer
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
