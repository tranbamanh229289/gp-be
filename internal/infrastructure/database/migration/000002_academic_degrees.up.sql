CREATE TABLE academic_degrees (
    id SERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    cid INTEGER NOT NULL REFERENCES citizen_identities(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    degree_number VARCHAR(50) NOT NULL UNIQUE,
    degree_type VARCHAR(50) NOT NULL CHECK (degree_type IN ('bachelor', 'master', 'phd', 'associate_professor', 'full_professor')),
    major VARCHAR(255) NOT NULL,
    university VARCHAR(255) NOT NULL,
    graduate_year SMALLINT NOT NULL CHECK (graduate_year >= 1900 AND graduate_year <= 2100),
    gpa DECIMAL(4,2) NOT NULL CHECK (gpa >= 0 AND gpa <= 4),
    classification VARCHAR(50) NOT NULL CHECK (classification IN ('excellent', 'very_good', 'good', 'average', 'pass')),
    status VARCHAR(30) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'expired', 'revoked')),
    issue_date BIGINT NOT NULL,
    holder_did VARCHAR(255) NOT NULL CHECK (holder_did LIKE 'did:%'),
    issuer_did VARCHAR(255) NOT NULL CHECK (issuer_did LIKE 'did:%'),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ 
);

CREATE INDEX idx_academic_degrees_public_id ON academic_degrees(public_id);
CREATE INDEX idx_academic_degrees_degree_number ON academic_degrees(degree_number);
CREATE INDEX idx_academic_degrees_holder_did ON academic_degrees(holder_did);
CREATE INDEX idx_academic_degrees_issuer_did ON academic_degrees(issuer_did);