CREATE TABLE health_insurances (
    id SERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    cid INTEGER NOT NULL REFERENCES citizen_identities(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    insurance_number VARCHAR(20) NOT NULL UNIQUE CHECK (insurance_number ~ '^[0-9]{15}$'),
    insurance_type VARCHAR(100) NOT NULL,
    hospital VARCHAR(255) NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'expired', 'revoked')),
    start_date DATE NOT NULL,
    expiry_date DATE NOT NULL CHECK (expiry_date >= start_date),
    issuer_did VARCHAR(255) NOT NULL CHECK (issuer_did LIKE 'did:%'),
    holder_did VARCHAR(255) NOT NULL CHECK (holder_did LIKE 'did:%'),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ
);