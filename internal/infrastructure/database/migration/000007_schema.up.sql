

CREATE TABLE schemas (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    issuer_did VARCHAR(255) NOT NULL CHECK (issuer_did LIKE 'did:%') REFERENCES identities(did) ON UPDATE CASCADE ON DELETE RESTRICT,
    hash VARCHAR(128) NOT NULL,
    type VARCHAR(255) NOT NULL,
    version VARCHAR(64) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    is_merklized           BOOLEAN NOT NULL DEFAULT FALSE,
    json_schema JSONB NOT NULL,
    jsonld_context JSONB NOT NULL,
    schema_url VARCHAR(255) NOT NULL UNIQUE,
    context_url VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(32) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'revoked')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ
   );

CREATE TABLE schema_attributes (
    id                  SERIAL PRIMARY KEY,
    public_id           UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    schema_id           INTEGER NOT NULL REFERENCES schemas(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    name                VARCHAR(128) NOT NULL,
    title               VARCHAR(255) NOT NULL,
    type                VARCHAR(64) NOT NULL,
    description         TEXT NOT NULL,
    required            BOOLEAN NOT NULL DEFAULT FALSE,
    slot                VARCHAR(64),
    format              VARCHAR(64),
    pattern             VARCHAR(255),
    min_length          INTEGER,
    max_length          INTEGER,
    minimum             DOUBLE PRECISION,
    maximum             DOUBLE PRECISION,
    enum                JSONB,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at          TIMESTAMPTZ
);