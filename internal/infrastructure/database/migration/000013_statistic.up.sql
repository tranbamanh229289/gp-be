CREATE TABLE issuer_statistics (
    id BIGSERIAL PRIMARY KEY,
    issuer_did VARCHAR(255) NOT NULL UNIQUE,
    document_num BIGINT NOT NULL DEFAULT 0,
    schema_num BIGINT NOT NULL DEFAULT 0,
    credential_request_num BIGINT NOT NULL DEFAULT 0,
    credential_issued_num BIGINT NOT NULL DEFAULT 0,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- Holder Statistics
CREATE TABLE holder_statistics (
    id BIGSERIAL PRIMARY KEY,
    holder_did VARCHAR(255) NOT NULL UNIQUE,
    credential_request_num BIGINT NOT NULL DEFAULT 0,
    verifiable_credential_num BIGINT NOT NULL DEFAULT 0,
    proof_submission_num BIGINT NOT NULL DEFAULT 0,
    proof_accepted_num BIGINT NOT NULL DEFAULT 0,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- Verifier Statistics
CREATE TABLE verifier_statistics (
    id BIGSERIAL PRIMARY KEY,
    verifier_did VARCHAR(255) NOT NULL UNIQUE,
    proof_request_num BIGINT NOT NULL DEFAULT 0,
    proof_submission_num BIGINT NOT NULL DEFAULT 0,
    proof_accepted_num BIGINT NOT NULL DEFAULT 0,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
