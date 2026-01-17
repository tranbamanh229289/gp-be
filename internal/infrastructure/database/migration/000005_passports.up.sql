CREATE TABLE passports (
    id SERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    cid INTEGER NOT NULL REFERENCES citizen_identities(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    passport_number VARCHAR(15) NOT NULL UNIQUE,
    passport_type VARCHAR(100) NOT NULL CHECK (passport_type IN ('ordinary', 'diplomatic', 'official')),
    nationality CHAR(3) NOT NULL CHECK (nationality ~ '^[A-Z]{3}$'),
    mrz VARCHAR(100) NOT NULL, 
    status VARCHAR(30) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'expired', 'revoked')),
    issue_date DATE NOT NULL,
    expiry_date DATE NOT NULL CHECK (expiry_date >= issue_date),
    holder_did VARCHAR(255) NOT NULL CHECK (holder_did LIKE 'did:%'),
    issuer_did VARCHAR(255) NOT NULL CHECK (issuer_did LIKE 'did:%'),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ
);

CREATE INDEX idx_passports_public_id ON passports(public_id);
CREATE INDEX idx_passports_passport_number ON passports(passport_number);
CREATE INDEX idx_passports_holder_did ON passports(holder_did);
CREATE INDEX idx_passports_issuer_did ON passports(issuer_did);