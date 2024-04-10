
CREATE DATABASE   aa_explorer;
\connect aa_explorer;
create schema partman;
create extension pg_partman with schema partman;
\connect postgres;
create extension pg_cron;
-- function schedule_in_database(job_name text, schedule text, command text, database text, username text default NULL::text, active boolean default true) returns bigint
SELECT cron.schedule_in_database('aa_explorer_partman','@hourly', $$CALL partman.run_maintenance_proc()$$,'aa_explorer','postgres');
\connect aa_explorer;





create table if not exists public.network
(
    name         varchar                  not null,
    network      varchar                  not null
    constraint network_pk
    primary key,
    http_rpc     varchar                  not null,
    is_testnet   boolean                  not null,
    create_time  timestamp with time zone not null,
    update_time  timestamp with time zone,
    delete_time  timestamp with time zone,
    chain_id     bigint,
    chain_name   varchar,
    scan         varchar,
    scan_tx      varchar,
    scan_block   varchar,
    scan_address varchar,
    scan_name    varchar,
    db_config    jsonb
);

alter table public.network
    owner to postgres;

create unique index if not exists networks_network_key
    on public.network (network);

create table if not exists public.function_signature
(
    signature   varchar not null
    constraint function_signature_pk
    primary key,
    name        varchar,
    text        text,
    bytes       bytea,
    create_time timestamp with time zone
);

alter table public.function_signature
    owner to postgres;

create table if not exists public.token
(
    id               bigserial
    primary key,
    network          varchar,
    contract_address varchar,
    symbol           varchar,
    full_name        varchar,
    token_price      numeric,
    last_time        bigint,
    create_time      timestamp,
    update_time      timestamp,
    market_rank      bigint,
    type             varchar
);

alter table public.token
    owner to postgres;

