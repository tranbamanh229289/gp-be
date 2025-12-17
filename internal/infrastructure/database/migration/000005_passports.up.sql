CREATE TABLE passports (
    id SERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    cid INTEGER NOT NULL REFERENCES citizen_identities(id) ON DELETE CASCADE,
    passport_number VARCHAR(15) UNIQUE NOT NULL,
    nationality CHAR(3) NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'Active',
    mrz TEXT,
    issue_date DATE NOT NULL,
    expiry_date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ
);