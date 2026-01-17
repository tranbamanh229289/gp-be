CREATE TABLE verifiable_credentials (
   id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    crid BIGINT NOT NULL REFERENCES credential_requests(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    credential_id VARCHAR(255),
    holder_did VARCHAR(255) NOT NULL CHECK (holder_did LIKE 'did:%'),
    issuer_did VARCHAR(255) NOT NULL CHECK (issuer_did LIKE 'did:%'),
    schema_id BIGINT NOT NULL REFERENCES schemas(id) ON UPDATE CASCADE ON DELETE RESTRICT,
     schema_hash VARCHAR(128) NOT NULL,
    credential_subject JSONB NOT NULL,
    claim_subject VARCHAR(255) NOT NULL,
    claim_hi VARCHAR(255) NOT NULL,
    claim_hv VARCHAR(255) NOT NULL,
    claim_hex TEXT NOT NULL,
    claim_mtp BYTEA NOT NULL,
    rev_nonce BIGINT NOT NULL,
    auth_claim_hex TEXT NOT NULL,
    auth_claim_mtp BYTEA NOT NULL,
    issuer_state VARCHAR(66) NOT NULL,
    claims_tree_root VARCHAR(66) NOT NULL,
    rev_tree_root VARCHAR(66) NOT NULL,
    roots_tree_root VARCHAR(66) NOT NULL,
    issuance_date TIMESTAMPTZ,
    expiration_date TIMESTAMPTZ,
    status VARCHAR(30) NOT NULL DEFAULT 'issued'  CHECK (status IN ('notSigned', 'issued', 'revoked', 'expired')),
    signature VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ
);

CREATE INDEX idx_verifiable_credentials_public_id ON verifiable_credentials(public_id);
CREATE INDEX idx_verifiable_credentials_crid ON verifiable_credentials(crid);
CREATE INDEX idx_verifiable_credentials_schema_id ON verifiable_credentials(schema_id);
CREATE INDEX idx_verifiable_credentials_holder_did ON verifiable_credentials(holder_did);
CREATE INDEX idx_verifiable_credentials_issuer_did ON verifiable_credentials(issuer_did);