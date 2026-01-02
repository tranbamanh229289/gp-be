CREATE TABLE credential_requests (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    request_id VARCHAR(128) NOT NULL UNIQUE,
    holder_did VARCHAR(255) NOT NULL CHECK (holder_did LIKE 'did:%') REFERENCES identities(did) ON UPDATE CASCADE ON DELETE RESTRICT,
    issuer_did VARCHAR(255) NOT NULL CHECK (issuer_did LIKE 'did:%') REFERENCES identities(did) ON UPDATE CASCADE ON DELETE RESTRICT,
    schema_id BIGINT NOT NULL REFERENCES schemas(id),
    schema_hash VARCHAR(128) NOT NULL,
    data JSONB NOT NULL DEFAULT '{}'::jsonb,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    expiration BIGINT,
    created_time BIGINT,
    expires_time BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);