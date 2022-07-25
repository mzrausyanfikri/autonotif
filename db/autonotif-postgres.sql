CREATE DATABASE autonotifdb_osmosis;
GRANT ALL PRIVILEGES ON DATABASE autonotifdb_osmosis TO autonotif_usr;

CREATE DATABASE autonotifdb_cosmoshub;
GRANT ALL PRIVILEGES ON DATABASE autonotifdb_cosmoshub TO autonotif_usr;


SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;
SET default_tablespace = '';
SET default_table_access_method = heap;

CREATE TABLE public.proposals (
    id bigint NOT NULL,
    proposal_id bigint NOT NULL,
    chain_type text NOT NULL,
    raw_data jsonb NOT NULL,
    created_at timestamp without time zone NOT NULL
);

ALTER TABLE public.proposals OWNER TO autonotif_usr;

CREATE SEQUENCE public.proposals_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE public.proposals_id_seq OWNER TO autonotif_usr;
ALTER SEQUENCE public.proposals_id_seq OWNED BY public.proposals.id;

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);

ALTER TABLE public.schema_migrations OWNER TO autonotif_usr;
ALTER TABLE ONLY public.proposals ALTER COLUMN id SET DEFAULT nextval('public.proposals_id_seq'::regclass);
SELECT pg_catalog.setval('public.proposals_id_seq', 1, true);

ALTER TABLE ONLY public.proposals
    ADD CONSTRAINT proposals_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);

CREATE UNIQUE INDEX proposals_chain_type_proposal_id_key ON public.proposals USING btree (chain_type, proposal_id);
