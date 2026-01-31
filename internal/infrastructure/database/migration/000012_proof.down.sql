DROP INDEX IF EXISTS idx_proof_submissions_holder_did;
DROP INDEX IF EXISTS idx_proof_submissions_request_id;
DROP INDEX IF EXISTS idx_proof_submissions_public_id;
DROP TABLE IF EXISTS proof_submissions;

DROP INDEX IF EXISTS idx_proof_requests_verifier_did;
DROP INDEX IF EXISTS idx_proof_requests_public_id;
DROP TABLE IF EXISTS proof_requests;