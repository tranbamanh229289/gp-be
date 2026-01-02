CREATE TABLE verifiable_credentials (
   id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    holder_did VARCHAR(255) NOT NULL CHECK (issuer_did LIKE 'did:%'),
    issuer_did VARCHAR(255) NOT NULL CHECK (issuer_did LIKE 'did:%'),
    schema_id BIGINT NOT NULL REFERENCES schemas(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    claim_subject VARCHAR(255),
    claim_hi VARCHAR(255) NOT NULL CHECK (LENGTH(claim_hi) = 64),
    claim_hv VARCHAR(255) NOT NULL CHECK (LENGTH(claim_hv) = 64),
    issuer_state VARCHAR(255) NOT NULL CHECK (LENGTH(issuer_state) = 64),
    rev_nonce BIGINT NOT NULL,
    issuance_date_date DATE,
    expiration_date DATE,
    proof_type VARCHAR(100) CHECK (proof_type IN ('BjjSignature2021', 'Iden3SparseMerkleTreeProof')),
    status VARCHAR(30) NOT NULL DEFAULT 'issued'  CHECK (status IN ('notSigned', 'issued', 'revoked', 'expired')),
    signature VARCHAR(255),
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ
);