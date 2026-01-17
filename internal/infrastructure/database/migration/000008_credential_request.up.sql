CREATE TABLE credential_requests (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    thread_id  VARCHAR(255) NOT NULL,
    holder_did VARCHAR(255) NOT NULL CHECK (holder_did LIKE 'did:%') REFERENCES identities(did) ON UPDATE CASCADE ON DELETE RESTRICT,
    issuer_did VARCHAR(255) NOT NULL CHECK (issuer_did LIKE 'did:%') REFERENCES identities(did) ON UPDATE CASCADE ON DELETE RESTRICT,
    schema_id BIGINT NOT NULL REFERENCES schemas(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    schema_hash VARCHAR(128) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    expiration BIGINT,
    created_time BIGINT,
    expires_time BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_credential_request_public_id ON credential_requests(public_id);
CREATE INDEX idx_credential_request_thread_id ON credential_requests(thread_id);
CREATE INDEX idx_credential_request_schema_id ON credential_requests(schema_id);
CREATE INDEX idx_credential_request_holder_did ON credential_requests(holder_did);
CREATE INDEX idx_credential_request_issuer_did ON credential_requests(issuer_did);