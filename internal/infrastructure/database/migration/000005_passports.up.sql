CREATE TABLE passports (
    id SERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    cid INTEGER NOT NULL REFERENCES citizen_identities(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    passport_number VARCHAR(15) NOT NULL UNIQUE,
    passport_type VARCHAR(100) CHECK (passport_type IN ('Ordinary', 'Diplomatic', 'Official') OR passport_type IS NULL),
    nationality CHAR(3) NOT NULL CHECK (nationality ~ '^[A-Z]{3}$'),
    mrz VARCHAR(100) CHECK (mrz IS NULL OR length(mrz) = 88),  -- Chuẩn ICAO: 88 ký tự
    status VARCHAR(30) NOT NULL DEFAULT 'Active' CHECK (status IN ('Active', 'Expired', 'Revoked')),
    issue_date DATE NOT NULL,
    expiry_date DATE NOT NULL CHECK (expiry_date >= issue_date),
    issuer_did VARCHAR(255) CHECK (issuer_did = '' OR issuer_did LIKE 'did:%'),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ
);