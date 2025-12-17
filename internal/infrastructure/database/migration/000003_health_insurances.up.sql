CREATE TABLE health_insurances (
    id SERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    cid INTEGER NOT NULL REFERENCES citizen_identities(id) ON DELETE CASCADE,
    insurance_number VARCHAR(20) UNIQUE NOT NULL,
    insurance_type VARCHAR(100),
    hospital VARCHAR(255),
    status VARCHAR(30) NOT NULL DEFAULT 'Active',
    start_date DATE NOT NULL,
    expiry_date DATE NOT NULL,
    issuer_by VARCHAR(255),
    issuer_did VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ
);
