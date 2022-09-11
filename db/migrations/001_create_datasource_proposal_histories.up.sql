BEGIN;

CREATE TABLE IF NOT EXISTS datasource_proposal_histories (
    id              BIGSERIAL   PRIMARY KEY,
    proposal_id     BIGINT      NOT NULL ,
    chain_type      TEXT        NOT NULL,
    raw_data        JSONB       NOT NULL,
    created_at      TIMESTAMP   NOT NULL
);

CREATE INDEX IF NOT EXISTS dp_histories_last_proposal_id_idx ON datasource_proposal_histories(chain_type, created_at DESC, proposal_id);

COMMIT;
