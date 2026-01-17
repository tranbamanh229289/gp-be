CREATE TABLE proof_requests (
    id                              SERIAL PRIMARY KEY,
    public_id                       UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    thread_id                       VARCHAR(255) NOT NULL,
    verifier_did                    VARCHAR(255) NOT NULL CHECK (verifier_did LIKE 'did:%') REFERENCES identities(did) ON DELETE CASCADE ON UPDATE RESTRICT,
    callback_url                    TEXT NOT NULL,
    reason                          TEXT,
    message                         TEXT,
    scope_id                        INTEGER,
    circuit_id                      VARCHAR(100) NOT NULL,
    credential_subject              JSONB,
    schema_id                       BIGINT NOT NULL REFERENCES schemas(id) ON DELETE CASCADE ON UPDATE RESTRICT,
    proof_type                      VARCHAR(100),
    skip_claim_revocation_check     BOOLEAN DEFAULT FALSE,
    allowed_issuers_did             JSONB,
    group_id                        INTEGER,
    nullifier_session               TEXT,
    status                          VARCHAR(50) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'expired', 'cancelled')),
    created_time                    BIGINT,
    expires_time                    BIGINT,
    created_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_proof_requests_public_id ON proof_requests(public_id);
CREATE INDEX idx_proof_requests_verifier_did ON proof_requests(verifier_did);

CREATE TABLE proof_responses (
    id          SERIAL PRIMARY KEY,
    public_id   UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    request_id  BIGINT NOT NULL REFERENCES proof_requests(id) ON DELETE CASCADE ON UPDATE RESTRICT,
    holder_did  VARCHAR(255) NOT NULL CHECK (holder_did LIKE 'did:%'),
    thread_id   VARCHAR(255) NOT NULL,
    status      VARCHAR(50) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'success', 'failed')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_proof_responses_public_id ON proof_responses(public_id);
CREATE INDEX idx_proof_responses_request_id ON proof_responses(request_id);
CREATE INDEX idx_proof_responses_holder_did ON proof_responses(holder_did);
