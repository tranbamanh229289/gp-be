CREATE TABLE driver_licenses (
    id SERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    cid INTEGER NOT NULL REFERENCES citizen_identities(id) ON DELETE CASCADE,
    license_number VARCHAR(30) UNIQUE NOT NULL,
    class VARCHAR(20) NOT NULL,
    point SMALLINT NOT NULL DEFAULT 12,
    point_reset_date DATE,
    status VARCHAR(30) NOT NULL DEFAULT 'Active',
    issue_date DATE NOT NULL,
    expiry_date DATE NOT NULL,
    issuer_by VARCHAR(255),
    issuer_did VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ
);

