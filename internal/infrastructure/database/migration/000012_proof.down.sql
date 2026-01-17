DROP INDEX IF EXISTS idx_proof_requests_verifier_did;
DROP INDEX IF EXISTS idx_proof_requests_public_id;
DROP TABLE IF EXISTS proof_responses;

DROP INDEX IF EXISTS idx_proof_responses_holder_did;
DROP INDEX IF EXISTS idx_proof_responses_request_id;
DROP INDEX IF EXISTS idx_proof_responses_public_id;
DROP TABLE IF EXISTS proof_requests;
