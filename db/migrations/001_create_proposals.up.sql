BEGIN;

CREATE TABLE IF NOT EXISTS proposals (
    id              BIGSERIAL   PRIMARY KEY,
    proposal_id     BIGINT      NOT NULL ,
    chain_type      TEXT        NOT NULL,
    raw_data        JSONB       NOT NULL,
    created_at      TIMESTAMP   NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS proposals_chain_type_proposal_id_key ON proposals (chain_type, proposal_id);

COMMIT;
