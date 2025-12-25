

CREATE TABLE schemas (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    issuer_id BIGINT NOT NULL REFERENCES identities(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    type VARCHAR(255) NOT NULL,
    version VARCHAR(64) NOT NULL,
    title VARCHAR(255),
    description TEXT,
    json_schema JSONB NOT NULL,
    jsonld_context JSONB NOT NULL,
    json_cid VARCHAR(255) NOT NULL UNIQUE CHECK (LENGTH(json_cid) = 59 AND json_cid LIKE 'Qm%'),
    jsonld_cid VARCHAR(255) NOT NULL UNIQUE CHECK (LENGTH(jsonld_cid) = 59 AND jsonld_cid LIKE 'Qm%'),
    status VARCHAR(32) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'revoked')),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ
   )