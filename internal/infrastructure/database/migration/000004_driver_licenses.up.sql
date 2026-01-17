CREATE TABLE driver_licenses (
    id SERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    cid INTEGER NOT NULL REFERENCES citizen_identities(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    license_number VARCHAR(30) NOT NULL UNIQUE,
    class VARCHAR(20) NOT NULL,
    point SMALLINT NOT NULL DEFAULT 12 CHECK (point >= 0 AND point <= 12),
    status VARCHAR(30) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'expired', 'revoked')),
    issue_date DATE NOT NULL,
    expiry_date DATE NOT NULL CHECK (expiry_date >= issue_date),
    issuer_did VARCHAR(255) NOT NULL CHECK (issuer_did LIKE 'did:%'),
    holder_did VARCHAR(255) NOT NULL CHECK (holder_did LIKE 'did:%'),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ
);

CREATE INDEX idx_driver_licenses_public_id ON driver_licenses(public_id);
CREATE INDEX idx_driver_licenses_license_number ON driver_licenses(license_number);
CREATE INDEX idx_driver_licenses_holder_did ON driver_licenses(holder_did);
CREATE INDEX idx_driver_licenses_issuer_did ON driver_licenses(issuer_did);