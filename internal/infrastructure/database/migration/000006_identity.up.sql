CREATE TABLE identities (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    public_key_x VARCHAR(255),
    public_key_y VARCHAR(255),
    name VARCHAR(255),
    role VARCHAR(100) NOT NULL CHECK (role IN ('holder', 'issuer', 'verifier')),
    did VARCHAR(255) NOT NULL UNIQUE CHECK (did LIKE 'did:%'),
    state VARCHAR(255) CHECK (state IS NULL OR length(state) = 64),
    claims_mt_id BIGINT,
    rev_mt_id BIGINT,
    roots_mt_id BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);