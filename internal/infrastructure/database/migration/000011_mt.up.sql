CREATE TABLE mt_nodes (
    mt_id BIGINT,
    key BYTEA,
    type SMALLINT NOT NULL,
    child_l BYTEA,
    child_r BYTEA,
    entry BYTEA,
    created_at BIGINT,
    deleted_at BIGINT,
    PRIMARY KEY(mt_id, key)
);

CREATE TABLE mt_roots (
    mt_id BIGINT PRIMARY KEY,
    key BYTEA,
    created_at BIGINT,
    deleted_at BIGINT
);

CREATE SEQUENCE mt_id_seq
    AS BIGINT
    START WITH 1
    INCREMENT BY 1
    NO CYCLE
    CACHE 10;

CREATE TABLE mt_metadata (
    id          SMALLINT PRIMARY KEY,
    mt_id       BIGINT NOT NULL,
    updated_at  TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT mt_metadata_singleton CHECK (id = 1)
);
ALTER SEQUENCE mt_id_seq OWNED BY mt_metadata.mt_id;

INSERT INTO mt_metadata (id, mt_id)
VALUES (1, 0);
