CREATE TABLE proof_requests (
    id              SERIAL PRIMARY KEY,
    public_id       UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    thread_id       VARCHAR(255) NOT NULL,
    verifier_did    VARCHAR(255) NOT NULL CHECK (verifier_did LIKE 'did:%') REFERENCES identities(did) ON DELETE CASCADE ON UPDATE RESTRICT,
    callback_url    TEXT NOT NULL,
    reason          TEXT,
    message         TEXT,
    scope_id        INTEGER,
    circuit_id      VARCHAR(100) NOT NULL,
    params          JSONB,
    query           JSONB,
    schema_id       BIGINT NOT NULL REFERENCES schemas(id) ON DELETE CASCADE ON UPDATE RESTRICT,
    allowed_issuers_did TEXT[],
    status          VARCHAR(50) NOT NULL DEFAULT 'active' 
        CHECK (status IN ('active', 'expired', 'cancelled')),
    created_time BIGINT,
    expires_time BIGINT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

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