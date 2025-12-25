CREATE TABLE state_transitions (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    identity_id BIGINT NOT NULL REFERENCES identities(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    old_state VARCHAR(255) NOT NULL CHECK (LENGTH(old_state) = 64),
    new_state VARCHAR(255) NOT NULL CHECK (LENGTH(new_state) = 64),
    tx_hash VARCHAR(100),
    block_number BIGINT CHECK (block_number >= 0),
    timestamp TIMESTAMPTZ,
    is_genesis BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);