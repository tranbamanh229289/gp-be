CREATE TABLE credential_requests (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    holder_id BIGINT NOT NULL REFERENCES identities(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    schema_id BIGINT NOT NULL REFERENCES schemas(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    credential_data JSONB NOT NULL,
    message TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);