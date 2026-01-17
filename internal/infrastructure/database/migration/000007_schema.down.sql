DROP INDEX IF EXISTS idx_schema_attributes_public_id;
DROP INDEX IF EXISTS idx_schema_attributes_schema_id;
DROP TABLE IF EXISTS schema_attributes;

DROP INDEX IF EXISTS idx_schemas_schema_url;
DROP INDEX IF EXISTS idx_schemas_context_url;
DROP INDEX IF EXISTS idx_schemas_hash;
DROP INDEX IF EXISTS idx_schemas_issuer_did;
DROP INDEX IF EXISTS idx_schemas_public_id;

DROP TABLE IF EXISTS schemas;
