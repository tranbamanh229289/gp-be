CREATE TABLE driver_licenses (
    id SERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    cid INTEGER NOT NULL REFERENCES citizen_identities(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    license_number VARCHAR(30) NOT NULL UNIQUE,
    class VARCHAR(20) NOT NULL,
    point SMALLINT NOT NULL DEFAULT 12 CHECK (point >= 0 AND point <= 12),
    point_reset_date DATE,
    status VARCHAR(30) NOT NULL DEFAULT 'Active' CHECK (status IN ('Active', 'Expired', 'Revoked', 'Suspended')),
    issue_date DATE NOT NULL,
    expiry_date DATE NOT NULL CHECK (expiry_date >= issue_date),
    issuer_did VARCHAR(255) CHECK (issuer_did = '' OR issuer_did LIKE 'did:%'),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ
);