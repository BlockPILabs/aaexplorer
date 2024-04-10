
CREATE DATABASE   block_db;
\connect block_db;
create schema partman;
create extension pg_partman with schema partman;
\connect postgres;
create extension pg_cron;
-- function schedule_in_database(job_name text, schedule text, command text, database text, username text default NULL::text, active boolean default true) returns bigint
SELECT cron.schedule_in_database('block_db_partman','@hourly', $$CALL partman.run_maintenance_proc()$$,'block_db','postgres');
\connect block_db;


create table if not exists public.block_data_decode
(
    time              timestamp with time zone,
    create_time       timestamp with time zone,
    number            bigint,
    hash              text,
    parent_hash       text,
    nonce             numeric,
    sha3_uncles       text,
    logs_bloom        text,
    transactions_root text,
    state_root        text,
    receipts_root     text,
    miner             text,
    mix_hash          text,
    difficulty        numeric,
    total_difficulty  numeric,
    extra_data        text,
    size              numeric,
    gas_limit         numeric,
    gas_used          numeric,
    timestamp         numeric,
    transaction_count bigint,
    uncles            text[],
    base_fee_per_gas  numeric
)
    partition by RANGE ("time");

alter table public.block_data_decode
    owner to postgres;





create index if not exists block_data_decode_hash_index
    on public.block_data_decode using hash (hash);

create index if not exists block_data_decode_block_num_index
    on public.block_data_decode (number);

create index if not exists block_data_decode_create_time_index
    on public.block_data_decode (create_time);

create unique index if not exists block_data_decode_time_hash_index
    on public.block_data_decode (time, number);

create table if not exists public.aa_block_info
(
    time               timestamp with time zone,
    create_time        timestamp with time zone,
    number             bigint,
    hash               text,
    userop_count       integer,
    userop_mev_count   integer,
    bundler_profit     numeric,
    bundler_profit_usd numeric default 0
)
    partition by RANGE ("time");

alter table public.aa_block_info
    owner to postgres;
create index if not exists aa_block_info_hash_index
    on public.aa_block_info using hash (hash);

create index if not exists aa_block_info_block_num_index
    on public.aa_block_info (number);

create index if not exists aa_block_info_create_time_index
    on public.aa_block_info (create_time);

create unique index if not exists aa_block_info_time_hash_index
    on public.aa_block_info (time, number);

create table if not exists public.transaction_decode
(
    time                     timestamp with time zone not null,
    create_time              timestamp with time zone,
    hash                     text,
    block_hash               text,
    block_number             bigint,
    nonce                    numeric,
    transaction_index        bigint,
    from_addr                text,
    to_addr                  text,
    value                    numeric,
    gas_price                numeric,
    gas                      numeric,
    input                    text,
    r                        text,
    s                        text,
    v                        bigint,
    chain_id                 bigint,
    type                     text,
    max_fee_per_gas          numeric,
    max_priority_fee_per_gas numeric,
    access_list              jsonb,
    method                   text
)
    partition by RANGE ("time");

alter table public.transaction_decode
    owner to postgres;


create index if not exists transaction_decode_hash_index
    on public.transaction_decode using hash (hash);

create index if not exists transaction_decode_block_num_index
    on public.transaction_decode (block_number);

create index if not exists transaction_decode_create_time_index
    on public.transaction_decode (create_time);

create unique index if not exists transaction_decode_time_hash_index
    on public.transaction_decode (time, hash);

create index if not exists transaction_decode_from_addr_index
    on public.transaction_decode (from_addr);

create index if not exists transaction_decode_to_addr_index
    on public.transaction_decode (to_addr);

create table if not exists public.transaction_receipt_decode
(
    time                timestamp with time zone,
    create_time         timestamp with time zone,
    transaction_hash    text,
    transaction_index   bigint,
    block_hash          text,
    block_number        bigint,
    cumulative_gas_used numeric,
    gas_used            numeric,
    contract_address    text,
    root                text,
    status              text,
    from_addr           text,
    to_addr             text,
    logs                jsonb,
    logs_bloom          text,
    revert_reason       text,
    type                text,
    effective_gas_price text
)
    partition by RANGE ("time");

alter table public.transaction_receipt_decode
    owner to postgres;

create index if not exists transaction_receipt_decode_hash_index
    on public.transaction_receipt_decode using hash (transaction_hash);

create index if not exists transaction_receipt_decode_block_num_index
    on public.transaction_receipt_decode (block_number);

create unique index if not exists transaction_receipt_decode_transaction_hash_index
    on public.transaction_receipt_decode (time, transaction_hash);

create table if not exists public.block_sync
(
    block_num   bigint not null
        constraint block_sync_pk
            primary key,
    scanned     boolean default false,
    create_time timestamp with time zone,
    update_time timestamp with time zone
);

alter table public.block_sync
    owner to postgres;

create index if not exists block_sync_scanned_index
    on public.block_sync (scanned);

create trigger transaction_block_sync
    after insert
    on public.block_sync
    for each row
execute procedure public.transaction_block_sync();

create trigger transaction_receipt_block_sync
    after insert
    on public.block_sync
    for each row
execute procedure public.transaction_receipt_block_sync();

create trigger aa_block_sync
    after update
    on public.block_sync
    for each row
execute procedure public.aa_block_sync();

create table if not exists public.transaction_sync
(
    block_num   bigint not null
        constraint transaction_sync_pk
            primary key,
    scanned     boolean default false,
    create_time timestamp with time zone,
    update_time timestamp with time zone
);

alter table public.transaction_sync
    owner to postgres;

create index if not exists transaction_sync_scanned_index
    on public.transaction_sync (scanned);

create trigger aa_tx_sync
    after update
    on public.transaction_sync
    for each row
execute procedure public.aa_tx_sync();

create table if not exists public.transaction_receipt_block_sync
(
    block_num   bigint not null
        constraint transaction_receipt_block_sync_pk
            primary key,
    scanned     boolean default false,
    create_time timestamp with time zone,
    update_time timestamp with time zone
);

alter table public.transaction_receipt_block_sync
    owner to postgres;

create index if not exists transaction_receipt_block_sync_index
    on public.transaction_receipt_block_sync (scanned);

create trigger aa_txr_sync
    after update
    on public.transaction_receipt_block_sync
    for each row
execute procedure public.aa_txr_sync();

create table if not exists public.aa_block_sync
(
    block_num     bigint not null
        constraint aa_block_sync_pk
            primary key,
    block_scanned boolean,
    tx_scanned    boolean,
    txr_scanned   boolean,
    scanned       boolean,
    create_time   timestamp with time zone,
    update_time   timestamp with time zone,
    scan_count    integer default 0
);

alter table public.aa_block_sync
    owner to postgres;

create index if not exists aa_block_sync_scanned_index
    on public.aa_block_sync (scanned);

create trigger aa_scan_sync
    after update
    on public.aa_block_sync
    for each row
execute procedure public.aa_scan_sync();

create table if not exists public.account
(
    address     text not null
        primary key,
    is_contract boolean,
    tag         text[],
    label       text[],
    abi         text,
    update_time timestamp with time zone
);

alter table public.account
    owner to postgres;

create index if not exists account_is_contract
    on public.account (is_contract);

create table if not exists public.token_info
(
    address  text not null
        primary key,
    symbol   text,
    name     text,
    decimals bigint
);

alter table public.token_info
    owner to postgres;

create index if not exists token_info_address_index
    on public.token_info using hash (address);

create table if not exists public.account_sync
(
    block_num bigint not null
        constraint account_sync_pk
            primary key
);

alter table public.account_sync
    owner to postgres;

create table if not exists public.aa_transaction_info
(
    time                     timestamp with time zone not null,
    create_time              timestamp with time zone,
    hash                     text,
    block_hash               text,
    block_number             bigint,
    userop_count             integer,
    is_mev                   boolean,
    bundler_profit           numeric,
    bundler_profit_usd       numeric default 0,
    nonce                    bigint,
    transaction_index        bigint,
    from_addr                varchar(255),
    to_addr                  varchar(255),
    value                    bigint,
    gas_price                bigint,
    gas                      bigint,
    input                    text,
    r                        text,
    s                        text,
    v                        bigint,
    chain_id                 bigint,
    type                     varchar(255),
    max_fee_per_gas          bigint,
    max_priority_fee_per_gas bigint,
    access_list              jsonb,
    method                   text,
    contract_address         varchar(255),
    cumulative_gas_used      bigint,
    effective_gas_price      text,
    gas_used                 bigint,
    logs                     text,
    logs_bloom               text,
    status                   varchar(255)
)
    partition by RANGE ("time");

alter table public.aa_transaction_info
    owner to postgres;

create index if not exists aa_transaction_info_hash_index
    on public.aa_transaction_info using hash (hash);

create index if not exists aa_transaction_info_block_num_index
    on public.aa_transaction_info (block_number);

create index if not exists aa_transaction_info_create_time_index
    on public.aa_transaction_info (create_time);

create unique index if not exists aa_transaction_info_time_hash_index
    on public.aa_transaction_info (time, hash);

create index if not exists aa_transaction_info_from_addr
    on public.aa_transaction_info using hash (from_addr);

create index if not exists aa_transaction_info_to_addr
    on public.aa_transaction_info using hash (to_addr);

create index if not exists aa_transaction_info_from_addr_index
    on public.aa_transaction_info using hash (from_addr);

create index if not exists aa_transaction_info_to_addr_index
    on public.aa_transaction_info using hash (to_addr);

create table if not exists public.aa_user_ops_calldata
(
    time          timestamp with time zone not null,
    uuid          varchar(128),
    user_ops_hash varchar(128),
    tx_hash       varchar(128),
    block_number  bigint,
    network       varchar(128),
    sender        varchar(64),
    target        varchar(64),
    tx_value      numeric,
    source        varchar(128),
    calldata      text,
    tx_time       bigint,
    create_time   timestamp with time zone,
    update_time   timestamp with time zone,
    aa_index      integer default 0
)
    partition by RANGE ("time");

alter table public.aa_user_ops_calldata
    owner to postgres;

create index if not exists aa_user_ops_calldata_tx_hash_index
    on public.aa_user_ops_calldata using hash (tx_hash);

create index if not exists aa_user_ops_calldata_user_operation_hash_index
    on public.aa_user_ops_calldata using hash (user_ops_hash);

create index if not exists aa_user_ops_calldata_block_num_index
    on public.aa_user_ops_calldata (block_number);

create unique index if not exists aa_user_ops_calldata_time_uuid_index
    on public.aa_user_ops_calldata (time, uuid);

create table if not exists public.aa_user_ops_info
(
    time                     timestamp with time zone not null,
    user_operation_hash      varchar(128),
    tx_hash                  varchar(128),
    block_number             bigint,
    network                  varchar(128),
    sender                   varchar(64),
    target                   varchar(64),
    tx_value                 numeric,
    fee                      numeric,
    bundler                  varchar(64),
    entry_point              varchar(64),
    factory                  varchar(64),
    paymaster                varchar(64),
    paymaster_and_data       text,
    signature                text,
    calldata                 text,
    calldata_contract        varchar(64),
    nonce                    bigint,
    call_gas_limit           bigint,
    pre_verification_gas     bigint,
    verification_gas_limit   bigint,
    max_fee_per_gas          bigint,
    max_priority_fee_per_gas bigint,
    tx_time                  bigint,
    init_code                text,
    status                   integer,
    source                   varchar(128),
    actual_gas_cost          bigint,
    actual_gas_used          bigint,
    create_time              timestamp with time zone,
    usd_amount               numeric,
    update_time              timestamp with time zone,
    targets                  varchar(128)[],
    aa_index                 integer default 0,
    targets_count            integer default 0,
    fee_usd                  numeric default 0,
    tx_value_usd             numeric default 0
)
    partition by RANGE ("time");

alter table public.aa_user_ops_info
    owner to postgres;

create index if not exists aa_user_ops_info_tx_hash_index
    on public.aa_user_ops_info using hash (tx_hash);

create index if not exists aa_user_ops_info_user_operation_hash_index
    on public.aa_user_ops_info using hash (user_operation_hash);

create index if not exists aa_user_ops_info_block_num_index
    on public.aa_user_ops_info (block_number);

create unique index if not exists aa_user_ops_info_time_uuid_index
    on public.aa_user_ops_info (time, tx_hash, user_operation_hash);

create index if not exists aa_user_ops_info_bundler_index
    on public.aa_user_ops_info (bundler);

create index if not exists aa_user_ops_info_paymaster_index
    on public.aa_user_ops_info (paymaster);

create index if not exists aa_user_ops_info_factory_index
    on public.aa_user_ops_info (factory);

create table if not exists public.aa_contract_interact
(
    id               bigserial
        primary key,
    contract_address varchar(255),
    network          varchar(255),
    statistic_type   varchar(255),
    interact_num     bigint,
    create_time      timestamp default CURRENT_TIMESTAMP not null
);

alter table public.aa_contract_interact
    owner to postgres;

create table if not exists public.aa_hot_token_statistic
(
    id             bigserial
        constraint hot_aa_token_statistic_pkey
            primary key,
    token_symbol   varchar(255)    not null,
    network        varchar(255)    not null,
    statistic_type varchar(255)    not null,
    volume         numeric(32, 18) not null,
    create_time    timestamp(3) default CURRENT_TIMESTAMP
);

alter table public.aa_hot_token_statistic
    owner to postgres;

create table if not exists public.asset_change_trace
(
    id               bigint generated by default as identity
        primary key,
    tx_hash          text,
    block_number     bigint,
    address          text,
    address_type     integer,
    create_time      timestamp(3) default CURRENT_TIMESTAMP not null,
    sync_flag        smallint,
    network          varchar(255),
    update_time      timestamp,
    last_change_time bigint
);

alter table public.asset_change_trace
    owner to postgres;

create index if not exists asset_change_trace_address_index
    on public.asset_change_trace using hash (address);

create index if not exists asset_change_trace_tx_hash_index
    on public.asset_change_trace using hash (tx_hash);

create table if not exists public.block_scan_record
(
    id                bigserial
        primary key,
    network           varchar(255),
    last_block_number bigint,
    create_time       timestamp(3) default CURRENT_TIMESTAMP not null,
    update_time       timestamp,
    last_scan_time    timestamp
);

alter table public.block_scan_record
    owner to postgres;

create table if not exists public.bundler_info
(
    bundler             varchar(255)                              not null
        primary key,
    network             varchar(255),
    user_ops_num        bigint          default 0,
    bundles_num         bigint          default 0,
    gas_collected       numeric(50, 20) default 0,
    user_ops_num_d1     bigint          default 0,
    bundles_num_d1      bigint          default 0,
    gas_collected_d1    numeric(50, 20) default 0,
    user_ops_num_d7     bigint          default 0,
    bundles_num_d7      bigint          default 0,
    gas_collected_d7    numeric(50, 20) default 0,
    user_ops_num_d30    bigint          default 0,
    bundles_num_d30     bigint          default 0,
    gas_collected_d30   numeric(50, 20) default 0,
    create_time         timestamp(3)    default CURRENT_TIMESTAMP not null,
    update_time         timestamp(3),
    fee_earned          numeric(50, 20) default 0,
    fee_earned_usd      numeric(50, 20) default 0,
    success_rate        numeric(50, 20) default 0,
    bundle_rate         numeric(50, 20) default 0,
    fee_earned_d1       numeric(50, 20) default 0,
    fee_earned_usd_d1   numeric(50, 20) default 0,
    success_rate_d1     numeric(50, 20) default 0,
    bundle_rate_d1      numeric(50, 20) default 0,
    fee_earned_d7       numeric(50, 20) default 0,
    fee_earned_usd_d7   numeric(50, 20) default 0,
    success_rate_d7     numeric(50, 20) default 0,
    bundle_rate_d7      numeric(50, 20) default 0,
    fee_earned_d30      numeric(50, 20) default 0,
    fee_earned_usd_d30  numeric(50, 20) default 0,
    success_rate_d30    numeric(50, 20) default 0,
    bundle_rate_d30     numeric(50, 20) default 0,
    success_bundles_num bigint          default 0,
    failed_bundles_num  bigint          default 0
);

alter table public.bundler_info
    owner to postgres;

create table if not exists public.bundler_statis_day
(
    id                  bigserial
        primary key,
    bundler             varchar(255),
    network             varchar(255),
    user_ops_num        bigint          default 0,
    bundles_num         bigint          default 0,
    gas_collected       numeric(50, 20) default 0,
    statis_time         timestamp,
    create_time         timestamp(3)    default CURRENT_TIMESTAMP not null,
    fee_earned          numeric(50, 20) default 0,
    total_num           bigint          default 0,
    success_bundles_num bigint          default 0,
    failed_bundles_num  bigint          default 0
);

alter table public.bundler_statis_day
    owner to postgres;

create index if not exists statis_time__index
    on public.bundler_statis_day (statis_time);

create table if not exists public.bundler_statis_hour
(
    id                  bigserial
        primary key,
    bundler             varchar(255),
    network             varchar(255),
    user_ops_num        bigint          default 0,
    bundles_num         bigint          default 0,
    gas_collected       numeric(50, 20) default 0,
    statis_time         timestamp,
    create_time         timestamp(3)    default CURRENT_TIMESTAMP not null,
    fee_earned          numeric(50, 20) default 0,
    total_num           bigint          default 0,
    success_bundles_num bigint          default 0,
    failed_bundles_num  bigint          default 0
);

alter table public.bundler_statis_hour
    owner to postgres;

create table if not exists public.daily_statistic_day
(
    id                     bigserial
        primary key,
    network                varchar(255),
    tx_num                 bigint          default 0,
    user_ops_num           bigint          default 0,
    gas_fee                numeric(50, 20) default 0,
    active_wallet          bigint          default 0,
    create_time            timestamp(3)    default CURRENT_TIMESTAMP not null,
    bundler_gas_profit     numeric(50, 20) default 0,
    bundler_gas_profit_usd numeric(50, 20) default 0,
    paymaster_gas_paid     numeric(50, 20) default 0,
    paymaster_gas_paid_usd numeric(50, 20) default 0,
    statistic_time         bigint,
    gas_fee_usd            numeric(50, 20) default 0,
    aa_tx_num              bigint          default 0
);

alter table public.daily_statistic_day
    owner to postgres;

create table if not exists public.daily_statistic_hour
(
    id                     bigserial
        primary key,
    network                varchar(255),
    tx_num                 bigint          default 0,
    user_ops_num           bigint          default 0,
    gas_fee                numeric(50, 20) default 0,
    active_wallet          bigint          default 0,
    create_time            timestamp(3)    default CURRENT_TIMESTAMP not null,
    bundler_gas_profit     numeric(50, 20) default 0,
    bundler_gas_profit_usd numeric(50, 20) default 0,
    paymaster_gas_paid     numeric(50, 20) default 0,
    paymaster_gas_paid_usd numeric(50, 20) default 0,
    statistic_time         bigint,
    gas_fee_usd            numeric(50, 20) default 0,
    aa_tx_num              bigint          default 0
);

alter table public.daily_statistic_hour
    owner to postgres;

create table if not exists public.factory_info
(
    factory                varchar                                  not null
        primary key,
    network                varchar,
    account_num            integer        default 0,
    account_deploy_num     integer        default 0,
    account_num_d1         integer        default 0,
    account_deploy_num_d1  integer        default 0,
    account_num_d7         integer        default 0,
    account_deploy_num_d7  integer        default 0,
    account_num_d30        integer        default 0,
    account_deploy_num_d30 integer        default 0,
    create_time            timestamp      default CURRENT_TIMESTAMP not null,
    update_time            timestamp,
    dominance              numeric(50, 4) default 0,
    dominance_d1           numeric(50, 4) default 0,
    dominance_d7           numeric(50, 4) default 0,
    dominance_d30          numeric(50, 4) default 0
);

alter table public.factory_info
    owner to postgres;

create table if not exists public.factory_statis_day
(
    id                 bigserial
        primary key,
    factory            varchar(255),
    network            varchar(255),
    account_num        bigint       default 0,
    account_deploy_num bigint       default 0,
    statis_time        timestamp,
    create_time        timestamp(3) default CURRENT_TIMESTAMP not null
);

alter table public.factory_statis_day
    owner to postgres;

create table if not exists public.factory_statis_hour
(
    id                 bigserial
        primary key,
    factory            varchar(255),
    network            varchar(255),
    account_num        bigint       default 0,
    account_deploy_num bigint       default 0,
    statis_time        timestamp,
    create_time        timestamp(3) default CURRENT_TIMESTAMP not null
);

alter table public.factory_statis_hour
    owner to postgres;

create table if not exists public.paymaster_info
(
    paymaster             varchar(255) not null
        primary key,
    network               varchar(255),
    user_ops_num          bigint          default 0,
    gas_sponsored         numeric(50, 20) default 0,
    user_ops_num_d1       bigint          default 0,
    gas_sponsored_d1      numeric(50, 20) default 0,
    user_ops_num_d7       bigint          default 0,
    gas_sponsored_d7      numeric(50, 20) default 0,
    user_ops_num_d30      bigint          default 0,
    gas_sponsored_d30     numeric(50, 20) default 0,
    create_time           timestamp(3)    default CURRENT_TIMESTAMP,
    update_time           timestamp(3),
    reserve               numeric(50, 20),
    reserve_usd           numeric,
    gas_sponsored_usd     numeric         default 0,
    gas_sponsored_usd_d1  numeric         default 0,
    gas_sponsored_usd_d7  numeric         default 0,
    gas_sponsored_usd_d30 numeric         default 0
);

alter table public.paymaster_info
    owner to postgres;

create table if not exists public.paymaster_statis_day
(
    id                bigserial
        primary key,
    paymaster         varchar(255),
    network           varchar(255),
    user_ops_num      bigint          default 0,
    gas_sponsored     numeric(50, 20) default 0,
    statis_time       timestamp,
    create_time       timestamp(3)    default CURRENT_TIMESTAMP not null,
    reserve           numeric(50, 20),
    reserve_usd       numeric,
    gas_sponsored_usd numeric(50, 20) default 0
);

alter table public.paymaster_statis_day
    owner to postgres;

create table if not exists public.paymaster_statis_hour
(
    id                bigserial
        primary key,
    paymaster         varchar(255),
    network           varchar(255),
    user_ops_num      bigint          default 0,
    gas_sponsored     numeric(50, 20) default 0,
    statis_time       timestamp,
    create_time       timestamp(3)    default CURRENT_TIMESTAMP not null,
    reserve           numeric(50, 20),
    reserve_usd       numeric,
    gas_sponsored_usd numeric(50, 20) default 0
);

alter table public.paymaster_statis_hour
    owner to postgres;

create table if not exists public.task_record
(
    id          bigint                              not null
        primary key,
    network     varchar,
    task_type   varchar,
    last_time   timestamp,
    create_time timestamp default CURRENT_TIMESTAMP not null,
    update_time timestamp
);

alter table public.task_record
    owner to postgres;

create table if not exists public.token_price_info
(
    id               bigserial
        primary key,
    contract_address varchar(255),
    symbol           varchar(255),
    token_price      numeric(50, 30),
    last_time        bigint,
    create_time      timestamp(3) default CURRENT_TIMESTAMP not null,
    update_time      timestamp,
    network          varchar(255),
    type             text
);

alter table public.token_price_info
    owner to postgres;

create index if not exists token_price_info_contract_address_index
    on public.token_price_info using hash (contract_address);

create table if not exists public.user_asset_info
(
    id               bigserial
        primary key,
    account_address  varchar(255),
    contract_address varchar(255),
    symbol           varchar(255),
    network          varchar(255),
    amount           numeric(50, 20),
    last_time        bigint,
    create_time      timestamp(3) default CURRENT_TIMESTAMP not null
);

alter table public.user_asset_info
    owner to postgres;

create index if not exists user_asset_info_account_addres_index
    on public.user_asset_info using hash (account_address);

create index if not exists user_asset_info_contract_address_index
    on public.user_asset_info using hash (contract_address);

create table if not exists public.user_op_type_statistic
(
    id             bigserial
        primary key,
    user_op_type   varchar(255),
    user_op_sign   varchar(255),
    network        varchar(255),
    statistic_type varchar(255),
    op_num         bigint,
    create_time    timestamp default CURRENT_TIMESTAMP not null
);

alter table public.user_op_type_statistic
    owner to postgres;

create table if not exists public.fix_task_record
(
    id          bigint                              not null
        primary key,
    network     varchar,
    task_type   varchar,
    start_time  timestamp,
    end_time    timestamp,
    create_time timestamp default CURRENT_TIMESTAMP not null,
    update_time timestamp
);

alter table public.fix_task_record
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

create table if not exists public.aa_asset
(
    user_address varchar(255) not null
        primary key,
    network      varchar(255),
    balance      numeric(50, 20),
    asset_value  numeric(50, 20),
    last_time    bigint,
    create_time  timestamp    not null,
    update_time  timestamp
);

alter table public.aa_asset
    owner to postgres;

create table if not exists public.aa_asset_detail
(
    id               bigserial
        primary key,
    user_address     varchar(255),
    network          varchar(255),
    contract_address varchar(255),
    symbol           varchar(255),
    is_native        varchar(255),
    asset_amount     numeric(50, 20),
    asset_value      numeric(50, 20),
    last_time        bigint,
    create_time      timestamp not null,
    update_time      timestamp
);

alter table public.aa_asset_detail
    owner to postgres;

create index if not exists aa_asset_detail_user_address_idx
    on public.aa_asset_detail using hash (user_address);

create index if not exists aa_asset_detail_contract_address_idx
    on public.aa_asset_detail using hash (contract_address);

create table if not exists public.whale_statistic_hour
(
    id             bigserial
        primary key,
    network        varchar(255),
    whale_num      bigint,
    total_usd      numeric(50, 20),
    create_time    timestamp not null,
    statistic_time bigint
);

alter table public.whale_statistic_hour
    owner to postgres;

create table if not exists public.whale_statistic_day
(
    id             bigserial
        primary key,
    network        varchar(255),
    whale_num      bigint,
    total_usd      numeric(50, 20),
    create_time    timestamp not null,
    statistic_time bigint
);

alter table public.whale_statistic_day
    owner to postgres;

create index if not exists whale_statistic_day_statistic_time_idx
    on public.whale_statistic_day (statistic_time);

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p1
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 0);

alter table public.aa_account_data_p1
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p2
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 1);

alter table public.aa_account_data_p2
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p3
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 2);

alter table public.aa_account_data_p3
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p4
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 3);

alter table public.aa_account_data_p4
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p5
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 4);

alter table public.aa_account_data_p5
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p6
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 5);

alter table public.aa_account_data_p6
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p7
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 6);

alter table public.aa_account_data_p7
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p8
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 7);

alter table public.aa_account_data_p8
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p9
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 8);

alter table public.aa_account_data_p9
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p10
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 9);

alter table public.aa_account_data_p10
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p11
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 10);

alter table public.aa_account_data_p11
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p12
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 11);

alter table public.aa_account_data_p12
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p13
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 12);

alter table public.aa_account_data_p13
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p14
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 13);

alter table public.aa_account_data_p14
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p15
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 14);

alter table public.aa_account_data_p15
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p16
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 15);

alter table public.aa_account_data_p16
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p17
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 16);

alter table public.aa_account_data_p17
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p18
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 17);

alter table public.aa_account_data_p18
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p19
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 18);

alter table public.aa_account_data_p19
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p20
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 19);

alter table public.aa_account_data_p20
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p21
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 20);

alter table public.aa_account_data_p21
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p22
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 21);

alter table public.aa_account_data_p22
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p23
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 22);

alter table public.aa_account_data_p23
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p24
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 23);

alter table public.aa_account_data_p24
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p25
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 24);

alter table public.aa_account_data_p25
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p26
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 25);

alter table public.aa_account_data_p26
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p27
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 26);

alter table public.aa_account_data_p27
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p28
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 27);

alter table public.aa_account_data_p28
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data_p29
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 28);

alter table public.aa_account_data_p29
    owner to postgres;

-- Cyclic dependencies found

create table if not exists public.aa_account_data
(
    address           text not null
        primary key,
    aa_type           text,
    factory           text,
    factory_time      timestamp with time zone,
    user_ops_num      bigint          default 0,
    total_balance_usd numeric(50, 20) default 0,
    last_time         bigint          default 0,
    update_time       timestamp with time zone
)
    partition by HASH (address);

alter table public.aa_account_data
    owner to postgres;

create table if not exists public.aa_account_data_p30
    partition of public.aa_account_data
        (
            primary key (address)
            )
        FOR VALUES WITH (modulus 30, remainder 29);

alter table public.aa_account_data_p30
    owner to postgres;

create index if not exists aa_account_data_aa_type
    on public.aa_account_data (aa_type);

create index if not exists aa_account_data_factory
    on public.aa_account_data (factory);

