CREATE TABLE citizen_identities (
    id SERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(), 
    id_number VARCHAR(20) NOT NULL UNIQUE CHECK (id_number ~ '^[0-9]{12}$'),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    gender VARCHAR(20) NOT NULL CHECK (gender IN ('male', 'female', 'other')),
    date_of_birth DATE NOT NULL,
    place_of_birth TEXT NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'expired', 'revoked')),
    issue_date DATE NOT NULL,
    expiry_date DATE NOT NULL CHECK (expiry_date >= issue_date AND issue_date > date_of_birth),
    issuer_did VARCHAR(255) NOT NULL CHECK (issuer_did LIKE 'did:%'),
    holder_did VARCHAR(255) NOT NULL CHECK (holder_did LIKE 'did:%'),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ
);