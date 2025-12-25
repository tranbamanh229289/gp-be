CREATE TABLE academic_degrees (
    id SERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    cid INTEGER NOT NULL REFERENCES citizen_identities(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    degree_number VARCHAR(50) NOT NULL UNIQUE,
    degree_type VARCHAR(50) NOT NULL CHECK (degree_type IN ('Bachelor', 'Master', 'PhD', 'AssociateProfessor', 'FullProfessor')),
    major VARCHAR(255) NOT NULL,
    university VARCHAR(255) NOT NULL,
    graduate_year SMALLINT NOT NULL CHECK (graduate_year >= 1900 AND graduate_year <= 2100),
    gpa DECIMAL(3,2) CHECK (gpa IS NULL OR (gpa >= 0 AND gpa <= 4)),
    classification VARCHAR(50) CHECK (classification IN ('Excellent', 'VeryGood', 'Good', 'Average', 'Pass') OR classification IS NULL),
    status VARCHAR(30) NOT NULL DEFAULT 'Active' CHECK (status IN ('Active', 'Expired', 'Revoked')),
    issue_date DATE NOT NULL,
    issuer_did VARCHAR(255) CHECK (issuer_did = '' OR issuer_did LIKE 'did:%'),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ
)