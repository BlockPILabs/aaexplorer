create table public.block_data_decode
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

create index block_data_decode_hash_index
    on public.block_data_decode using hash (hash);

create index block_data_decode_block_num_index
    on public.block_data_decode (number);

create index block_data_decode_create_time_index
    on public.block_data_decode (create_time);

create unique index block_data_decode_time_hash_index
    on public.block_data_decode (time, number);

create table public.aa_block_info
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

create index aa_block_info_hash_index
    on public.aa_block_info using hash (hash);

create index aa_block_info_block_num_index
    on public.aa_block_info (number);

create index aa_block_info_create_time_index
    on public.aa_block_info (create_time);

create unique index aa_block_info_time_hash_index
    on public.aa_block_info (time, number);

create table public.transaction_decode
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

create index transaction_decode_hash_index
    on public.transaction_decode using hash (hash);

create index transaction_decode_block_num_index
    on public.transaction_decode (block_number);

create index transaction_decode_create_time_index
    on public.transaction_decode (create_time);

create unique index transaction_decode_time_hash_index
    on public.transaction_decode (time, hash);

create index transaction_decode_from_addr_index
    on public.transaction_decode (from_addr);

create index transaction_decode_to_addr_index
    on public.transaction_decode (to_addr);

create table public.transaction_receipt_decode
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

create index transaction_receipt_decode_hash_index
    on public.transaction_receipt_decode using hash (transaction_hash);

create index transaction_receipt_decode_block_num_index
    on public.transaction_receipt_decode (block_number);

create unique index transaction_receipt_decode_transaction_hash_index
    on public.transaction_receipt_decode (time, transaction_hash);

create table public.block_sync
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

create index block_sync_scanned_index
    on public.block_sync (scanned);

create table public.transaction_sync
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

create index transaction_sync_scanned_index
    on public.transaction_sync (scanned);

create table public.transaction_receipt_block_sync
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

create index transaction_receipt_block_sync_index
    on public.transaction_receipt_block_sync (scanned);

create table public.aa_block_sync
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

create index aa_block_sync_scanned_index
    on public.aa_block_sync (scanned);

create table public.account
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

create index account_is_contract
    on public.account (is_contract);

create table public.token_info
(
    address  text not null
        primary key,
    symbol   text,
    name     text,
    decimals bigint
);

alter table public.token_info
    owner to postgres;

create index token_info_address_index
    on public.token_info using hash (address);

create table public.account_sync
(
    block_num bigint not null
        constraint account_sync_pk
            primary key
);

alter table public.account_sync
    owner to postgres;

create table public.aa_transaction_info
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

create index aa_transaction_info_hash_index
    on public.aa_transaction_info using hash (hash);

create index aa_transaction_info_block_num_index
    on public.aa_transaction_info (block_number);

create index aa_transaction_info_create_time_index
    on public.aa_transaction_info (create_time);

create unique index aa_transaction_info_time_hash_index
    on public.aa_transaction_info (time, hash);

create index aa_transaction_info_from_addr
    on public.aa_transaction_info using hash (from_addr);

create index aa_transaction_info_to_addr
    on public.aa_transaction_info using hash (to_addr);

create index aa_transaction_info_from_addr_index
    on public.aa_transaction_info using hash (from_addr);

create index aa_transaction_info_to_addr_index
    on public.aa_transaction_info using hash (to_addr);

create table public.aa_user_ops_calldata
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

create index aa_user_ops_calldata_tx_hash_index
    on public.aa_user_ops_calldata using hash (tx_hash);

create index aa_user_ops_calldata_user_operation_hash_index
    on public.aa_user_ops_calldata using hash (user_ops_hash);

create index aa_user_ops_calldata_block_num_index
    on public.aa_user_ops_calldata (block_number);

create unique index aa_user_ops_calldata_time_uuid_index
    on public.aa_user_ops_calldata (time, uuid);

create table public.aa_user_ops_info
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

create index aa_user_ops_info_tx_hash_index
    on public.aa_user_ops_info using hash (tx_hash);

create index aa_user_ops_info_user_operation_hash_index
    on public.aa_user_ops_info using hash (user_operation_hash);

create index aa_user_ops_info_block_num_index
    on public.aa_user_ops_info (block_number);

create unique index aa_user_ops_info_time_uuid_index
    on public.aa_user_ops_info (time, tx_hash, user_operation_hash);

create index aa_user_ops_info_bundler_index
    on public.aa_user_ops_info (bundler);

create index aa_user_ops_info_paymaster_index
    on public.aa_user_ops_info (paymaster);

create index aa_user_ops_info_factory_index
    on public.aa_user_ops_info (factory);

create table public.aa_contract_interact
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

create table public.aa_hot_token_statistic
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

create table public.asset_change_trace
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

create index asset_change_trace_address_index
    on public.asset_change_trace using hash (address);

create index asset_change_trace_tx_hash_index
    on public.asset_change_trace using hash (tx_hash);

create table public.block_scan_record
(
    id                bigserial
        primary key,
    network           varchar(255),
    last_block_number bigint,
    create_time       timestamp(3) default CURRENT_TIMESTAMP not null,
    update_time       timestamp,
    last_scan_time    timestamp,
    type              varchar
);

alter table public.block_scan_record
    owner to postgres;

create table public.bundler_info
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

create table public.bundler_statis_day
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

create index statis_time__index
    on public.bundler_statis_day (statis_time);

create table public.bundler_statis_hour
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

create table public.daily_statistic_day
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

create table public.daily_statistic_hour
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

create table public.factory_info
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

create table public.factory_statis_day
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

create table public.factory_statis_hour
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

create table public.paymaster_info
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

create table public.paymaster_statis_day
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

create table public.paymaster_statis_hour
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

create table public.task_record
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

create table public.token_price_info
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

create index token_price_info_contract_address_index
    on public.token_price_info using hash (contract_address);

create table public.user_asset_info
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

create index user_asset_info_account_addres_index
    on public.user_asset_info using hash (account_address);

create index user_asset_info_contract_address_index
    on public.user_asset_info using hash (contract_address);

create table public.user_op_type_statistic
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

create table public.fix_task_record
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

create table public.token
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
    type             varchar,
    decimals         bigint,
    image_url        varchar(255)
);

alter table public.token
    owner to postgres;

create table public.aa_asset
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

create table public.aa_asset_detail
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

create index aa_asset_detail_user_address_idx
    on public.aa_asset_detail using hash (user_address);

create index aa_asset_detail_contract_address_idx
    on public.aa_asset_detail using hash (contract_address);

create table public.whale_statistic_hour
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

create table public.whale_statistic_day
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

create index whale_statistic_day_statistic_time_idx
    on public.whale_statistic_day (statistic_time);

create table public.aa_account_data
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

create index aa_account_data_aa_type
    on public.aa_account_data (aa_type);

create index aa_account_data_factory
    on public.aa_account_data (factory);

create table public.mev_transaction
(
    time                timestamp,
    tx_hash             varchar(255) not null
        primary key,
    block_hash          varchar(255),
    block_number        bigint,
    transaction_index   bigint,
    from_addr           varchar(255),
    to_addr             varchar(255),
    value               numeric(50, 20),
    gas_price           numeric(50, 20),
    gas                 numeric(50, 20),
    mev_type            varchar(255),
    victim              varchar(255),
    victim_type         varchar(255),
    victim_tx_hash      varchar(255),
    victim_block_number bigint,
    victim_from_addr    varchar(255),
    victim_to_addr      varchar(255),
    attacker            varchar(255),
    bundler_loss        numeric(50, 20),
    bundler_loss_usd    numeric(50, 20),
    mev_profit          numeric(50, 20),
    mev_profit_usd      numeric(50, 20),
    create_time         timestamp
);

alter table public.mev_transaction
    owner to postgres;

create index mev_tx_from_addr_idx
    on public.mev_transaction using hash (from_addr);

create index mev_tx_to_addr_idx
    on public.mev_transaction using hash (to_addr);

create index mev_tx_victim_idx
    on public.mev_transaction using hash (victim);

create index mev_tx_attacker_idx
    on public.mev_transaction using hash (attacker);

create index mev_tx_time_idx
    on public.mev_transaction (time);

create index mev_tx_block_number_idx
    on public.mev_transaction (block_number);

create table public.aa_account_data_p1
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 0);

alter table public.aa_account_data_p1
    owner to postgres;

create table public.aa_account_data_p2
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 1);

alter table public.aa_account_data_p2
    owner to postgres;

create table public.aa_account_data_p3
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 2);

alter table public.aa_account_data_p3
    owner to postgres;

create table public.aa_account_data_p4
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 3);

alter table public.aa_account_data_p4
    owner to postgres;

create table public.aa_account_data_p5
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 4);

alter table public.aa_account_data_p5
    owner to postgres;

create table public.aa_account_data_p6
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 5);

alter table public.aa_account_data_p6
    owner to postgres;

create table public.aa_account_data_p7
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 6);

alter table public.aa_account_data_p7
    owner to postgres;

create table public.aa_account_data_p8
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 7);

alter table public.aa_account_data_p8
    owner to postgres;

create table public.aa_account_data_p9
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 8);

alter table public.aa_account_data_p9
    owner to postgres;

create table public.aa_account_data_p10
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 9);

alter table public.aa_account_data_p10
    owner to postgres;

create table public.aa_account_data_p11
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 10);

alter table public.aa_account_data_p11
    owner to postgres;

create table public.aa_account_data_p12
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 11);

alter table public.aa_account_data_p12
    owner to postgres;

create table public.aa_account_data_p13
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 12);

alter table public.aa_account_data_p13
    owner to postgres;

create table public.aa_account_data_p14
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 13);

alter table public.aa_account_data_p14
    owner to postgres;

create table public.aa_account_data_p15
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 14);

alter table public.aa_account_data_p15
    owner to postgres;

create table public.aa_account_data_p16
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 15);

alter table public.aa_account_data_p16
    owner to postgres;

create table public.aa_account_data_p17
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 16);

alter table public.aa_account_data_p17
    owner to postgres;

create table public.aa_account_data_p18
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 17);

alter table public.aa_account_data_p18
    owner to postgres;

create table public.aa_account_data_p19
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 18);

alter table public.aa_account_data_p19
    owner to postgres;

create table public.aa_account_data_p20
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 19);

alter table public.aa_account_data_p20
    owner to postgres;

create table public.aa_account_data_p21
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 20);

alter table public.aa_account_data_p21
    owner to postgres;

create table public.aa_account_data_p22
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 21);

alter table public.aa_account_data_p22
    owner to postgres;

create table public.aa_account_data_p23
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 22);

alter table public.aa_account_data_p23
    owner to postgres;

create table public.aa_account_data_p24
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 23);

alter table public.aa_account_data_p24
    owner to postgres;

create table public.aa_account_data_p25
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 24);

alter table public.aa_account_data_p25
    owner to postgres;

create table public.aa_account_data_p26
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 25);

alter table public.aa_account_data_p26
    owner to postgres;

create table public.aa_account_data_p27
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 26);

alter table public.aa_account_data_p27
    owner to postgres;

create table public.aa_account_data_p28
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 27);

alter table public.aa_account_data_p28
    owner to postgres;

create table public.aa_account_data_p29
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 28);

alter table public.aa_account_data_p29
    owner to postgres;

create table public.aa_account_data_p30
    partition of public.aa_account_data
    FOR VALUES WITH (modulus 30, remainder 29);

alter table public.aa_account_data_p30
    owner to postgres;

create table public.block_data_decode_default
    partition of public.block_data_decode
    DEFAULT;

alter table public.block_data_decode_default
    owner to postgres;

create table public.aa_block_info_default
    partition of public.aa_block_info
    DEFAULT;

alter table public.aa_block_info_default
    owner to postgres;

create table public.transaction_decode_default
    partition of public.transaction_decode
    DEFAULT;

alter table public.transaction_decode_default
    owner to postgres;

create table public.transaction_receipt_decode_default
    partition of public.transaction_receipt_decode
    DEFAULT;

alter table public.transaction_receipt_decode_default
    owner to postgres;

create table public.aa_transaction_info_default
    partition of public.aa_transaction_info
    DEFAULT;

alter table public.aa_transaction_info_default
    owner to postgres;

create table public.aa_user_ops_calldata_default
    partition of public.aa_user_ops_calldata
    DEFAULT;

alter table public.aa_user_ops_calldata_default
    owner to postgres;

create table public.aa_user_ops_info_default
    partition of public.aa_user_ops_info
    DEFAULT;

alter table public.aa_user_ops_info_default
    owner to postgres;

create table public.transfer_transaction
(
    id                bigserial
        primary key,
    time              timestamp,
    create_time       timestamp,
    tx_hash           varchar(255),
    block_hash        varchar(255),
    block_number      bigint,
    transaction_index bigint,
    from_addr         varchar(255),
    to_addr           varchar(255),
    value             numeric(50, 20),
    gas_price         numeric(50, 20),
    gas               numeric(50, 20),
    transfer_value    numeric(50, 20),
    token_symbol      varchar(255),
    token_address     varchar(255),
    token_url         varchar(255)
);

alter table public.transfer_transaction
    owner to postgres;

create index transfer_transaction_time_idx
    on public.transfer_transaction (time);

create index transfer_transaction_create_time_idx
    on public.transfer_transaction (create_time);

create index transfer_transaction_tx_hash_idx
    on public.transfer_transaction using hash (tx_hash);

create index transfer_transaction_block_numer_idx
    on public.transfer_transaction (block_number);

create index transfer_transaction_from_addr_idx
    on public.transfer_transaction using hash (from_addr);

create index transfer_transaction_to_addr_idx
    on public.transfer_transaction using hash (to_addr);

create index transfer_transaction_token_address_idx
    on public.transfer_transaction using hash (token_address);

create table public.token_all
(
    id               bigint default nextval('token_id_seq'::regclass) not null
        constraint token_copy1_pkey
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
    type             varchar,
    decimals         bigint,
    image_url        varchar
);

alter table public.token_all
    owner to postgres;

create index token_all_contract_address_idx
    on public.token_all using hash (contract_address);

create index token_all_symbol_idx
    on public.token_all (symbol);

create table public.aa_block_info_p2025_03
    partition of public.aa_block_info
    FOR VALUES FROM ('2025-03-01 00:00:00+00') TO ('2025-04-01 00:00:00+00');

alter table public.aa_block_info_p2025_03
    owner to postgres;

create table public.aa_block_info_p2025_04
    partition of public.aa_block_info
    FOR VALUES FROM ('2025-04-01 00:00:00+00') TO ('2025-05-01 00:00:00+00');

alter table public.aa_block_info_p2025_04
    owner to postgres;

create table public.aa_block_info_p2025_05
    partition of public.aa_block_info
    FOR VALUES FROM ('2025-05-01 00:00:00+00') TO ('2025-06-01 00:00:00+00');

alter table public.aa_block_info_p2025_05
    owner to postgres;

create table public.aa_block_info_p2025_06
    partition of public.aa_block_info
    FOR VALUES FROM ('2025-06-01 00:00:00+00') TO ('2025-07-01 00:00:00+00');

alter table public.aa_block_info_p2025_06
    owner to postgres;

create table public.aa_block_info_p2025_07
    partition of public.aa_block_info
    FOR VALUES FROM ('2025-07-01 00:00:00+00') TO ('2025-08-01 00:00:00+00');

alter table public.aa_block_info_p2025_07
    owner to postgres;

create table public.block_data_decode_p2025w11
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-03-10 00:00:00+00') TO ('2025-03-17 00:00:00+00');

alter table public.block_data_decode_p2025w11
    owner to postgres;

create table public.aa_transaction_info_p2025w11
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-03-10 00:00:00+00') TO ('2025-03-17 00:00:00+00');

alter table public.aa_transaction_info_p2025w11
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w11
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-03-10 00:00:00+00') TO ('2025-03-17 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w11
    owner to postgres;

create table public.aa_user_ops_info_p2025w11
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-03-10 00:00:00+00') TO ('2025-03-17 00:00:00+00');

alter table public.aa_user_ops_info_p2025w11
    owner to postgres;

create table public.block_data_decode_p2025w12
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-03-17 00:00:00+00') TO ('2025-03-24 00:00:00+00');

alter table public.block_data_decode_p2025w12
    owner to postgres;

create table public.aa_transaction_info_p2025w12
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-03-17 00:00:00+00') TO ('2025-03-24 00:00:00+00');

alter table public.aa_transaction_info_p2025w12
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w12
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-03-17 00:00:00+00') TO ('2025-03-24 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w12
    owner to postgres;

create table public.aa_user_ops_info_p2025w12
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-03-17 00:00:00+00') TO ('2025-03-24 00:00:00+00');

alter table public.aa_user_ops_info_p2025w12
    owner to postgres;

create table public.aa_block_info_p2025_08
    partition of public.aa_block_info
    FOR VALUES FROM ('2025-08-01 00:00:00+00') TO ('2025-09-01 00:00:00+00');

alter table public.aa_block_info_p2025_08
    owner to postgres;

create table public.transaction_decode_p2025_03_16
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-03-16 00:00:00+00') TO ('2025-03-17 00:00:00+00');

alter table public.transaction_decode_p2025_03_16
    owner to postgres;

create table public.transaction_receipt_decode_p2025_03_16
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-03-16 00:00:00+00') TO ('2025-03-17 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_03_16
    owner to postgres;

create table public.block_data_decode_p2025w13
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-03-24 00:00:00+00') TO ('2025-03-31 00:00:00+00');

alter table public.block_data_decode_p2025w13
    owner to postgres;

create table public.aa_transaction_info_p2025w13
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-03-24 00:00:00+00') TO ('2025-03-31 00:00:00+00');

alter table public.aa_transaction_info_p2025w13
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w13
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-03-24 00:00:00+00') TO ('2025-03-31 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w13
    owner to postgres;

create table public.aa_user_ops_info_p2025w13
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-03-24 00:00:00+00') TO ('2025-03-31 00:00:00+00');

alter table public.aa_user_ops_info_p2025w13
    owner to postgres;

create table public.transaction_decode_p2025_03_17
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-03-17 00:00:00+00') TO ('2025-03-18 00:00:00+00');

alter table public.transaction_decode_p2025_03_17
    owner to postgres;

create table public.transaction_receipt_decode_p2025_03_17
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-03-17 00:00:00+00') TO ('2025-03-18 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_03_17
    owner to postgres;

create table public.transaction_decode_p2025_03_18
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-03-18 00:00:00+00') TO ('2025-03-19 00:00:00+00');

alter table public.transaction_decode_p2025_03_18
    owner to postgres;

create table public.transaction_receipt_decode_p2025_03_18
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-03-18 00:00:00+00') TO ('2025-03-19 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_03_18
    owner to postgres;

create table public.transaction_decode_p2025_03_19
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-03-19 00:00:00+00') TO ('2025-03-20 00:00:00+00');

alter table public.transaction_decode_p2025_03_19
    owner to postgres;

create table public.transaction_receipt_decode_p2025_03_19
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-03-19 00:00:00+00') TO ('2025-03-20 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_03_19
    owner to postgres;

create table public.transaction_decode_p2025_03_20
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-03-20 00:00:00+00') TO ('2025-03-21 00:00:00+00');

alter table public.transaction_decode_p2025_03_20
    owner to postgres;

create table public.transaction_receipt_decode_p2025_03_20
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-03-20 00:00:00+00') TO ('2025-03-21 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_03_20
    owner to postgres;

create table public.transaction_decode_p2025_03_21
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-03-21 00:00:00+00') TO ('2025-03-22 00:00:00+00');

alter table public.transaction_decode_p2025_03_21
    owner to postgres;

create table public.transaction_receipt_decode_p2025_03_21
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-03-21 00:00:00+00') TO ('2025-03-22 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_03_21
    owner to postgres;

create table public.transaction_decode_p2025_03_22
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-03-22 00:00:00+00') TO ('2025-03-23 00:00:00+00');

alter table public.transaction_decode_p2025_03_22
    owner to postgres;

create table public.transaction_receipt_decode_p2025_03_22
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-03-22 00:00:00+00') TO ('2025-03-23 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_03_22
    owner to postgres;

create table public.transaction_decode_p2025_03_23
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-03-23 00:00:00+00') TO ('2025-03-24 00:00:00+00');

alter table public.transaction_decode_p2025_03_23
    owner to postgres;

create table public.transaction_receipt_decode_p2025_03_23
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-03-23 00:00:00+00') TO ('2025-03-24 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_03_23
    owner to postgres;

create table public.block_data_decode_p2025w14
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-03-31 00:00:00+00') TO ('2025-04-07 00:00:00+00');

alter table public.block_data_decode_p2025w14
    owner to postgres;

create table public.aa_transaction_info_p2025w14
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-03-31 00:00:00+00') TO ('2025-04-07 00:00:00+00');

alter table public.aa_transaction_info_p2025w14
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w14
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-03-31 00:00:00+00') TO ('2025-04-07 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w14
    owner to postgres;

create table public.aa_user_ops_info_p2025w14
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-03-31 00:00:00+00') TO ('2025-04-07 00:00:00+00');

alter table public.aa_user_ops_info_p2025w14
    owner to postgres;

create table public.transaction_decode_p2025_03_24
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-03-24 00:00:00+00') TO ('2025-03-25 00:00:00+00');

alter table public.transaction_decode_p2025_03_24
    owner to postgres;

create table public.transaction_receipt_decode_p2025_03_24
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-03-24 00:00:00+00') TO ('2025-03-25 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_03_24
    owner to postgres;

create table public.transaction_decode_p2025_03_25
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-03-25 00:00:00+00') TO ('2025-03-26 00:00:00+00');

alter table public.transaction_decode_p2025_03_25
    owner to postgres;

create table public.transaction_receipt_decode_p2025_03_25
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-03-25 00:00:00+00') TO ('2025-03-26 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_03_25
    owner to postgres;

create table public.transaction_decode_p2025_03_26
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-03-26 00:00:00+00') TO ('2025-03-27 00:00:00+00');

alter table public.transaction_decode_p2025_03_26
    owner to postgres;

create table public.transaction_receipt_decode_p2025_03_26
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-03-26 00:00:00+00') TO ('2025-03-27 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_03_26
    owner to postgres;

create table public.transaction_decode_p2025_03_27
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-03-27 00:00:00+00') TO ('2025-03-28 00:00:00+00');

alter table public.transaction_decode_p2025_03_27
    owner to postgres;

create table public.transaction_receipt_decode_p2025_03_27
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-03-27 00:00:00+00') TO ('2025-03-28 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_03_27
    owner to postgres;

create table public.transaction_decode_p2025_03_28
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-03-28 00:00:00+00') TO ('2025-03-29 00:00:00+00');

alter table public.transaction_decode_p2025_03_28
    owner to postgres;

create table public.transaction_receipt_decode_p2025_03_28
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-03-28 00:00:00+00') TO ('2025-03-29 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_03_28
    owner to postgres;

create table public.transaction_decode_p2025_03_29
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-03-29 00:00:00+00') TO ('2025-03-30 00:00:00+00');

alter table public.transaction_decode_p2025_03_29
    owner to postgres;

create table public.transaction_receipt_decode_p2025_03_29
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-03-29 00:00:00+00') TO ('2025-03-30 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_03_29
    owner to postgres;

create table public.block_data_decode_p2025w15
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-04-07 00:00:00+00') TO ('2025-04-14 00:00:00+00');

alter table public.block_data_decode_p2025w15
    owner to postgres;

create table public.transaction_decode_p2025_03_30
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-03-30 00:00:00+00') TO ('2025-03-31 00:00:00+00');

alter table public.transaction_decode_p2025_03_30
    owner to postgres;

create table public.transaction_decode_p2025_03_31
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-03-31 00:00:00+00') TO ('2025-04-01 00:00:00+00');

alter table public.transaction_decode_p2025_03_31
    owner to postgres;

create table public.transaction_receipt_decode_p2025_03_30
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-03-30 00:00:00+00') TO ('2025-03-31 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_03_30
    owner to postgres;

create table public.transaction_receipt_decode_p2025_03_31
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-03-31 00:00:00+00') TO ('2025-04-01 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_03_31
    owner to postgres;

create table public.aa_transaction_info_p2025w15
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-04-07 00:00:00+00') TO ('2025-04-14 00:00:00+00');

alter table public.aa_transaction_info_p2025w15
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w15
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-04-07 00:00:00+00') TO ('2025-04-14 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w15
    owner to postgres;

create table public.aa_user_ops_info_p2025w15
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-04-07 00:00:00+00') TO ('2025-04-14 00:00:00+00');

alter table public.aa_user_ops_info_p2025w15
    owner to postgres;

create table public.transaction_decode_p2025_04_01
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-01 00:00:00+00') TO ('2025-04-02 00:00:00+00');

alter table public.transaction_decode_p2025_04_01
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_01
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-01 00:00:00+00') TO ('2025-04-02 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_01
    owner to postgres;

create table public.transaction_decode_p2025_04_02
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-02 00:00:00+00') TO ('2025-04-03 00:00:00+00');

alter table public.transaction_decode_p2025_04_02
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_02
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-02 00:00:00+00') TO ('2025-04-03 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_02
    owner to postgres;

create table public.transaction_decode_p2025_04_03
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-03 00:00:00+00') TO ('2025-04-04 00:00:00+00');

alter table public.transaction_decode_p2025_04_03
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_03
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-03 00:00:00+00') TO ('2025-04-04 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_03
    owner to postgres;

create table public.transaction_decode_p2025_04_04
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-04 00:00:00+00') TO ('2025-04-05 00:00:00+00');

alter table public.transaction_decode_p2025_04_04
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_04
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-04 00:00:00+00') TO ('2025-04-05 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_04
    owner to postgres;

create table public.block_data_decode_p2025w16
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-04-14 00:00:00+00') TO ('2025-04-21 00:00:00+00');

alter table public.block_data_decode_p2025w16
    owner to postgres;

create table public.transaction_decode_p2025_04_05
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-05 00:00:00+00') TO ('2025-04-06 00:00:00+00');

alter table public.transaction_decode_p2025_04_05
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_05
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-05 00:00:00+00') TO ('2025-04-06 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_05
    owner to postgres;

create table public.aa_transaction_info_p2025w16
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-04-14 00:00:00+00') TO ('2025-04-21 00:00:00+00');

alter table public.aa_transaction_info_p2025w16
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w16
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-04-14 00:00:00+00') TO ('2025-04-21 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w16
    owner to postgres;

create table public.aa_user_ops_info_p2025w16
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-04-14 00:00:00+00') TO ('2025-04-21 00:00:00+00');

alter table public.aa_user_ops_info_p2025w16
    owner to postgres;

create table public.transaction_decode_p2025_04_06
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-06 00:00:00+00') TO ('2025-04-07 00:00:00+00');

alter table public.transaction_decode_p2025_04_06
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_06
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-06 00:00:00+00') TO ('2025-04-07 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_06
    owner to postgres;

create table public.transaction_decode_p2025_04_07
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-07 00:00:00+00') TO ('2025-04-08 00:00:00+00');

alter table public.transaction_decode_p2025_04_07
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_07
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-07 00:00:00+00') TO ('2025-04-08 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_07
    owner to postgres;

create table public.transaction_decode_p2025_04_08
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-08 00:00:00+00') TO ('2025-04-09 00:00:00+00');

alter table public.transaction_decode_p2025_04_08
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_08
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-08 00:00:00+00') TO ('2025-04-09 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_08
    owner to postgres;

create table public.transaction_decode_p2025_04_09
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-09 00:00:00+00') TO ('2025-04-10 00:00:00+00');

alter table public.transaction_decode_p2025_04_09
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_09
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-09 00:00:00+00') TO ('2025-04-10 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_09
    owner to postgres;

create table public.transaction_decode_p2025_04_10
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-10 00:00:00+00') TO ('2025-04-11 00:00:00+00');

alter table public.transaction_decode_p2025_04_10
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_10
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-10 00:00:00+00') TO ('2025-04-11 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_10
    owner to postgres;

create table public.transaction_decode_p2025_04_11
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-11 00:00:00+00') TO ('2025-04-12 00:00:00+00');

alter table public.transaction_decode_p2025_04_11
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_11
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-11 00:00:00+00') TO ('2025-04-12 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_11
    owner to postgres;

create table public.block_data_decode_p2025w17
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-04-21 00:00:00+00') TO ('2025-04-28 00:00:00+00');

alter table public.block_data_decode_p2025w17
    owner to postgres;

create table public.transaction_decode_p2025_04_12
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-12 00:00:00+00') TO ('2025-04-13 00:00:00+00');

alter table public.transaction_decode_p2025_04_12
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_12
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-12 00:00:00+00') TO ('2025-04-13 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_12
    owner to postgres;

create table public.aa_transaction_info_p2025w17
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-04-21 00:00:00+00') TO ('2025-04-28 00:00:00+00');

alter table public.aa_transaction_info_p2025w17
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w17
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-04-21 00:00:00+00') TO ('2025-04-28 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w17
    owner to postgres;

create table public.aa_user_ops_info_p2025w17
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-04-21 00:00:00+00') TO ('2025-04-28 00:00:00+00');

alter table public.aa_user_ops_info_p2025w17
    owner to postgres;

create table public.transaction_decode_p2025_04_13
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-13 00:00:00+00') TO ('2025-04-14 00:00:00+00');

alter table public.transaction_decode_p2025_04_13
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_13
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-13 00:00:00+00') TO ('2025-04-14 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_13
    owner to postgres;

create table public.transaction_decode_p2025_04_14
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-14 00:00:00+00') TO ('2025-04-15 00:00:00+00');

alter table public.transaction_decode_p2025_04_14
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_14
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-14 00:00:00+00') TO ('2025-04-15 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_14
    owner to postgres;

create table public.transaction_decode_p2025_04_15
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-15 00:00:00+00') TO ('2025-04-16 00:00:00+00');

alter table public.transaction_decode_p2025_04_15
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_15
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-15 00:00:00+00') TO ('2025-04-16 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_15
    owner to postgres;

create table public.transaction_decode_p2025_04_16
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-16 00:00:00+00') TO ('2025-04-17 00:00:00+00');

alter table public.transaction_decode_p2025_04_16
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_16
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-16 00:00:00+00') TO ('2025-04-17 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_16
    owner to postgres;

create table public.aa_block_info_p2025_09
    partition of public.aa_block_info
    FOR VALUES FROM ('2025-09-01 00:00:00+00') TO ('2025-10-01 00:00:00+00');

alter table public.aa_block_info_p2025_09
    owner to postgres;

create table public.transaction_decode_p2025_04_17
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-17 00:00:00+00') TO ('2025-04-18 00:00:00+00');

alter table public.transaction_decode_p2025_04_17
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_17
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-17 00:00:00+00') TO ('2025-04-18 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_17
    owner to postgres;

create table public.transaction_decode_p2025_04_18
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-18 00:00:00+00') TO ('2025-04-19 00:00:00+00');

alter table public.transaction_decode_p2025_04_18
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_18
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-18 00:00:00+00') TO ('2025-04-19 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_18
    owner to postgres;

create table public.block_data_decode_p2025w18
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-04-28 00:00:00+00') TO ('2025-05-05 00:00:00+00');

alter table public.block_data_decode_p2025w18
    owner to postgres;

create table public.transaction_decode_p2025_04_19
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-19 00:00:00+00') TO ('2025-04-20 00:00:00+00');

alter table public.transaction_decode_p2025_04_19
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_19
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-19 00:00:00+00') TO ('2025-04-20 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_19
    owner to postgres;

create table public.aa_transaction_info_p2025w18
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-04-28 00:00:00+00') TO ('2025-05-05 00:00:00+00');

alter table public.aa_transaction_info_p2025w18
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w18
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-04-28 00:00:00+00') TO ('2025-05-05 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w18
    owner to postgres;

create table public.aa_user_ops_info_p2025w18
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-04-28 00:00:00+00') TO ('2025-05-05 00:00:00+00');

alter table public.aa_user_ops_info_p2025w18
    owner to postgres;

create table public.transaction_decode_p2025_04_20
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-20 00:00:00+00') TO ('2025-04-21 00:00:00+00');

alter table public.transaction_decode_p2025_04_20
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_20
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-20 00:00:00+00') TO ('2025-04-21 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_20
    owner to postgres;

create table public.transaction_decode_p2025_04_21
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-21 00:00:00+00') TO ('2025-04-22 00:00:00+00');

alter table public.transaction_decode_p2025_04_21
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_21
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-21 00:00:00+00') TO ('2025-04-22 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_21
    owner to postgres;

create table public.transaction_decode_p2025_04_22
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-22 00:00:00+00') TO ('2025-04-23 00:00:00+00');

alter table public.transaction_decode_p2025_04_22
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_22
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-22 00:00:00+00') TO ('2025-04-23 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_22
    owner to postgres;

create table public.transaction_decode_p2025_04_23
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-23 00:00:00+00') TO ('2025-04-24 00:00:00+00');

alter table public.transaction_decode_p2025_04_23
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_23
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-23 00:00:00+00') TO ('2025-04-24 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_23
    owner to postgres;

create table public.transaction_decode_p2025_04_24
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-24 00:00:00+00') TO ('2025-04-25 00:00:00+00');

alter table public.transaction_decode_p2025_04_24
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_24
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-24 00:00:00+00') TO ('2025-04-25 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_24
    owner to postgres;

create table public.transaction_decode_p2025_04_25
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-25 00:00:00+00') TO ('2025-04-26 00:00:00+00');

alter table public.transaction_decode_p2025_04_25
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_25
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-25 00:00:00+00') TO ('2025-04-26 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_25
    owner to postgres;

create table public.block_data_decode_p2025w19
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-05-05 00:00:00+00') TO ('2025-05-12 00:00:00+00');

alter table public.block_data_decode_p2025w19
    owner to postgres;

create table public.transaction_decode_p2025_04_26
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-26 00:00:00+00') TO ('2025-04-27 00:00:00+00');

alter table public.transaction_decode_p2025_04_26
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_26
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-26 00:00:00+00') TO ('2025-04-27 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_26
    owner to postgres;

create table public.aa_transaction_info_p2025w19
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-05-05 00:00:00+00') TO ('2025-05-12 00:00:00+00');

alter table public.aa_transaction_info_p2025w19
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w19
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-05-05 00:00:00+00') TO ('2025-05-12 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w19
    owner to postgres;

create table public.aa_user_ops_info_p2025w19
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-05-05 00:00:00+00') TO ('2025-05-12 00:00:00+00');

alter table public.aa_user_ops_info_p2025w19
    owner to postgres;

create table public.transaction_decode_p2025_04_27
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-27 00:00:00+00') TO ('2025-04-28 00:00:00+00');

alter table public.transaction_decode_p2025_04_27
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_27
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-27 00:00:00+00') TO ('2025-04-28 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_27
    owner to postgres;

create table public.transaction_decode_p2025_04_28
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-28 00:00:00+00') TO ('2025-04-29 00:00:00+00');

alter table public.transaction_decode_p2025_04_28
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_28
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-28 00:00:00+00') TO ('2025-04-29 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_28
    owner to postgres;

create table public.transaction_decode_p2025_04_29
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-29 00:00:00+00') TO ('2025-04-30 00:00:00+00');

alter table public.transaction_decode_p2025_04_29
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_29
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-29 00:00:00+00') TO ('2025-04-30 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_29
    owner to postgres;

create table public.transaction_decode_p2025_04_30
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-04-30 00:00:00+00') TO ('2025-05-01 00:00:00+00');

alter table public.transaction_decode_p2025_04_30
    owner to postgres;

create table public.transaction_receipt_decode_p2025_04_30
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-04-30 00:00:00+00') TO ('2025-05-01 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_04_30
    owner to postgres;

create table public.transaction_decode_p2025_05_01
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-01 00:00:00+00') TO ('2025-05-02 00:00:00+00');

alter table public.transaction_decode_p2025_05_01
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_01
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-01 00:00:00+00') TO ('2025-05-02 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_01
    owner to postgres;

create table public.transaction_decode_p2025_05_02
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-02 00:00:00+00') TO ('2025-05-03 00:00:00+00');

alter table public.transaction_decode_p2025_05_02
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_02
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-02 00:00:00+00') TO ('2025-05-03 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_02
    owner to postgres;

create table public.transaction_decode_p2025_05_03
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-03 00:00:00+00') TO ('2025-05-04 00:00:00+00');

alter table public.transaction_decode_p2025_05_03
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_03
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-03 00:00:00+00') TO ('2025-05-04 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_03
    owner to postgres;

create table public.block_data_decode_p2025w20
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-05-12 00:00:00+00') TO ('2025-05-19 00:00:00+00');

alter table public.block_data_decode_p2025w20
    owner to postgres;

create table public.transaction_decode_p2025_05_04
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-04 00:00:00+00') TO ('2025-05-05 00:00:00+00');

alter table public.transaction_decode_p2025_05_04
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_04
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-04 00:00:00+00') TO ('2025-05-05 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_04
    owner to postgres;

create table public.aa_transaction_info_p2025w20
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-05-12 00:00:00+00') TO ('2025-05-19 00:00:00+00');

alter table public.aa_transaction_info_p2025w20
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w20
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-05-12 00:00:00+00') TO ('2025-05-19 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w20
    owner to postgres;

create table public.aa_user_ops_info_p2025w20
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-05-12 00:00:00+00') TO ('2025-05-19 00:00:00+00');

alter table public.aa_user_ops_info_p2025w20
    owner to postgres;

create table public.transaction_decode_p2025_05_05
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-05 00:00:00+00') TO ('2025-05-06 00:00:00+00');

alter table public.transaction_decode_p2025_05_05
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_05
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-05 00:00:00+00') TO ('2025-05-06 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_05
    owner to postgres;

create table public.transaction_decode_p2025_05_06
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-06 00:00:00+00') TO ('2025-05-07 00:00:00+00');

alter table public.transaction_decode_p2025_05_06
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_06
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-06 00:00:00+00') TO ('2025-05-07 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_06
    owner to postgres;

create table public.transaction_decode_p2025_05_07
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-07 00:00:00+00') TO ('2025-05-08 00:00:00+00');

alter table public.transaction_decode_p2025_05_07
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_07
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-07 00:00:00+00') TO ('2025-05-08 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_07
    owner to postgres;

create table public.transaction_decode_p2025_05_08
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-08 00:00:00+00') TO ('2025-05-09 00:00:00+00');

alter table public.transaction_decode_p2025_05_08
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_08
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-08 00:00:00+00') TO ('2025-05-09 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_08
    owner to postgres;

create table public.transaction_decode_p2025_05_09
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-09 00:00:00+00') TO ('2025-05-10 00:00:00+00');

alter table public.transaction_decode_p2025_05_09
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_09
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-09 00:00:00+00') TO ('2025-05-10 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_09
    owner to postgres;

create table public.transaction_decode_p2025_05_10
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-10 00:00:00+00') TO ('2025-05-11 00:00:00+00');

alter table public.transaction_decode_p2025_05_10
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_10
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-10 00:00:00+00') TO ('2025-05-11 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_10
    owner to postgres;

create table public.block_data_decode_p2025w21
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-05-19 00:00:00+00') TO ('2025-05-26 00:00:00+00');

alter table public.block_data_decode_p2025w21
    owner to postgres;

create table public.transaction_decode_p2025_05_11
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-11 00:00:00+00') TO ('2025-05-12 00:00:00+00');

alter table public.transaction_decode_p2025_05_11
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_11
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-11 00:00:00+00') TO ('2025-05-12 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_11
    owner to postgres;

create table public.aa_transaction_info_p2025w21
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-05-19 00:00:00+00') TO ('2025-05-26 00:00:00+00');

alter table public.aa_transaction_info_p2025w21
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w21
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-05-19 00:00:00+00') TO ('2025-05-26 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w21
    owner to postgres;

create table public.aa_user_ops_info_p2025w21
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-05-19 00:00:00+00') TO ('2025-05-26 00:00:00+00');

alter table public.aa_user_ops_info_p2025w21
    owner to postgres;

create table public.transaction_decode_p2025_05_12
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-12 00:00:00+00') TO ('2025-05-13 00:00:00+00');

alter table public.transaction_decode_p2025_05_12
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_12
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-12 00:00:00+00') TO ('2025-05-13 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_12
    owner to postgres;

create table public.transaction_decode_p2025_05_13
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-13 00:00:00+00') TO ('2025-05-14 00:00:00+00');

alter table public.transaction_decode_p2025_05_13
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_13
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-13 00:00:00+00') TO ('2025-05-14 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_13
    owner to postgres;

create table public.transaction_decode_p2025_05_14
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-14 00:00:00+00') TO ('2025-05-15 00:00:00+00');

alter table public.transaction_decode_p2025_05_14
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_14
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-14 00:00:00+00') TO ('2025-05-15 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_14
    owner to postgres;

create table public.transaction_decode_p2025_05_15
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-15 00:00:00+00') TO ('2025-05-16 00:00:00+00');

alter table public.transaction_decode_p2025_05_15
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_15
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-15 00:00:00+00') TO ('2025-05-16 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_15
    owner to postgres;

create table public.transaction_decode_p2025_05_16
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-16 00:00:00+00') TO ('2025-05-17 00:00:00+00');

alter table public.transaction_decode_p2025_05_16
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_16
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-16 00:00:00+00') TO ('2025-05-17 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_16
    owner to postgres;

create table public.aa_block_info_p2025_10
    partition of public.aa_block_info
    FOR VALUES FROM ('2025-10-01 00:00:00+00') TO ('2025-11-01 00:00:00+00');

alter table public.aa_block_info_p2025_10
    owner to postgres;

create table public.transaction_decode_p2025_05_17
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-17 00:00:00+00') TO ('2025-05-18 00:00:00+00');

alter table public.transaction_decode_p2025_05_17
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_17
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-17 00:00:00+00') TO ('2025-05-18 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_17
    owner to postgres;

create table public.block_data_decode_p2025w22
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-05-26 00:00:00+00') TO ('2025-06-02 00:00:00+00');

alter table public.block_data_decode_p2025w22
    owner to postgres;

create table public.transaction_decode_p2025_05_18
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-18 00:00:00+00') TO ('2025-05-19 00:00:00+00');

alter table public.transaction_decode_p2025_05_18
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_18
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-18 00:00:00+00') TO ('2025-05-19 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_18
    owner to postgres;

create table public.aa_transaction_info_p2025w22
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-05-26 00:00:00+00') TO ('2025-06-02 00:00:00+00');

alter table public.aa_transaction_info_p2025w22
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w22
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-05-26 00:00:00+00') TO ('2025-06-02 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w22
    owner to postgres;

create table public.aa_user_ops_info_p2025w22
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-05-26 00:00:00+00') TO ('2025-06-02 00:00:00+00');

alter table public.aa_user_ops_info_p2025w22
    owner to postgres;

create table public.transaction_decode_p2025_05_19
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-19 00:00:00+00') TO ('2025-05-20 00:00:00+00');

alter table public.transaction_decode_p2025_05_19
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_19
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-19 00:00:00+00') TO ('2025-05-20 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_19
    owner to postgres;

create table public.transaction_decode_p2025_05_20
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-20 00:00:00+00') TO ('2025-05-21 00:00:00+00');

alter table public.transaction_decode_p2025_05_20
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_20
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-20 00:00:00+00') TO ('2025-05-21 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_20
    owner to postgres;

create table public.transaction_decode_p2025_05_21
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-21 00:00:00+00') TO ('2025-05-22 00:00:00+00');

alter table public.transaction_decode_p2025_05_21
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_21
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-21 00:00:00+00') TO ('2025-05-22 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_21
    owner to postgres;

create table public.transaction_decode_p2025_05_22
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-22 00:00:00+00') TO ('2025-05-23 00:00:00+00');

alter table public.transaction_decode_p2025_05_22
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_22
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-22 00:00:00+00') TO ('2025-05-23 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_22
    owner to postgres;

create table public.transaction_decode_p2025_05_23
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-23 00:00:00+00') TO ('2025-05-24 00:00:00+00');

alter table public.transaction_decode_p2025_05_23
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_23
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-23 00:00:00+00') TO ('2025-05-24 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_23
    owner to postgres;

create table public.transaction_decode_p2025_05_24
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-24 00:00:00+00') TO ('2025-05-25 00:00:00+00');

alter table public.transaction_decode_p2025_05_24
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_24
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-24 00:00:00+00') TO ('2025-05-25 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_24
    owner to postgres;

create table public.block_data_decode_p2025w23
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-06-02 00:00:00+00') TO ('2025-06-09 00:00:00+00');

alter table public.block_data_decode_p2025w23
    owner to postgres;

create table public.transaction_decode_p2025_05_25
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-25 00:00:00+00') TO ('2025-05-26 00:00:00+00');

alter table public.transaction_decode_p2025_05_25
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_25
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-25 00:00:00+00') TO ('2025-05-26 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_25
    owner to postgres;

create table public.aa_transaction_info_p2025w23
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-06-02 00:00:00+00') TO ('2025-06-09 00:00:00+00');

alter table public.aa_transaction_info_p2025w23
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w23
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-06-02 00:00:00+00') TO ('2025-06-09 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w23
    owner to postgres;

create table public.aa_user_ops_info_p2025w23
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-06-02 00:00:00+00') TO ('2025-06-09 00:00:00+00');

alter table public.aa_user_ops_info_p2025w23
    owner to postgres;

create table public.transaction_decode_p2025_05_26
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-26 00:00:00+00') TO ('2025-05-27 00:00:00+00');

alter table public.transaction_decode_p2025_05_26
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_26
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-26 00:00:00+00') TO ('2025-05-27 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_26
    owner to postgres;

create table public.transaction_decode_p2025_05_27
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-27 00:00:00+00') TO ('2025-05-28 00:00:00+00');

alter table public.transaction_decode_p2025_05_27
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_27
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-27 00:00:00+00') TO ('2025-05-28 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_27
    owner to postgres;

create table public.transaction_decode_p2025_05_28
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-28 00:00:00+00') TO ('2025-05-29 00:00:00+00');

alter table public.transaction_decode_p2025_05_28
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_28
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-28 00:00:00+00') TO ('2025-05-29 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_28
    owner to postgres;

create table public.transaction_decode_p2025_05_29
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-29 00:00:00+00') TO ('2025-05-30 00:00:00+00');

alter table public.transaction_decode_p2025_05_29
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_29
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-29 00:00:00+00') TO ('2025-05-30 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_29
    owner to postgres;

create table public.transaction_decode_p2025_05_30
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-30 00:00:00+00') TO ('2025-05-31 00:00:00+00');

alter table public.transaction_decode_p2025_05_30
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_30
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-30 00:00:00+00') TO ('2025-05-31 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_30
    owner to postgres;

create table public.block_data_decode_p2025w24
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-06-09 00:00:00+00') TO ('2025-06-16 00:00:00+00');

alter table public.block_data_decode_p2025w24
    owner to postgres;

create table public.transaction_decode_p2025_05_31
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-05-31 00:00:00+00') TO ('2025-06-01 00:00:00+00');

alter table public.transaction_decode_p2025_05_31
    owner to postgres;

create table public.transaction_receipt_decode_p2025_05_31
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-05-31 00:00:00+00') TO ('2025-06-01 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_05_31
    owner to postgres;

create table public.aa_transaction_info_p2025w24
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-06-09 00:00:00+00') TO ('2025-06-16 00:00:00+00');

alter table public.aa_transaction_info_p2025w24
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w24
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-06-09 00:00:00+00') TO ('2025-06-16 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w24
    owner to postgres;

create table public.aa_user_ops_info_p2025w24
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-06-09 00:00:00+00') TO ('2025-06-16 00:00:00+00');

alter table public.aa_user_ops_info_p2025w24
    owner to postgres;

create table public.transaction_decode_p2025_06_01
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-01 00:00:00+00') TO ('2025-06-02 00:00:00+00');

alter table public.transaction_decode_p2025_06_01
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_01
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-01 00:00:00+00') TO ('2025-06-02 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_01
    owner to postgres;

create table public.transaction_decode_p2025_06_02
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-02 00:00:00+00') TO ('2025-06-03 00:00:00+00');

alter table public.transaction_decode_p2025_06_02
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_02
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-02 00:00:00+00') TO ('2025-06-03 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_02
    owner to postgres;

create table public.transaction_decode_p2025_06_03
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-03 00:00:00+00') TO ('2025-06-04 00:00:00+00');

alter table public.transaction_decode_p2025_06_03
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_03
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-03 00:00:00+00') TO ('2025-06-04 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_03
    owner to postgres;

create table public.transaction_decode_p2025_06_04
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-04 00:00:00+00') TO ('2025-06-05 00:00:00+00');

alter table public.transaction_decode_p2025_06_04
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_04
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-04 00:00:00+00') TO ('2025-06-05 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_04
    owner to postgres;

create table public.transaction_decode_p2025_06_05
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-05 00:00:00+00') TO ('2025-06-06 00:00:00+00');

alter table public.transaction_decode_p2025_06_05
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_05
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-05 00:00:00+00') TO ('2025-06-06 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_05
    owner to postgres;

create table public.transaction_decode_p2025_06_06
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-06 00:00:00+00') TO ('2025-06-07 00:00:00+00');

alter table public.transaction_decode_p2025_06_06
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_06
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-06 00:00:00+00') TO ('2025-06-07 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_06
    owner to postgres;

create table public.block_data_decode_p2025w25
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-06-16 00:00:00+00') TO ('2025-06-23 00:00:00+00');

alter table public.block_data_decode_p2025w25
    owner to postgres;

create table public.transaction_decode_p2025_06_07
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-07 00:00:00+00') TO ('2025-06-08 00:00:00+00');

alter table public.transaction_decode_p2025_06_07
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_07
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-07 00:00:00+00') TO ('2025-06-08 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_07
    owner to postgres;

create table public.aa_transaction_info_p2025w25
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-06-16 00:00:00+00') TO ('2025-06-23 00:00:00+00');

alter table public.aa_transaction_info_p2025w25
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w25
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-06-16 00:00:00+00') TO ('2025-06-23 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w25
    owner to postgres;

create table public.aa_user_ops_info_p2025w25
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-06-16 00:00:00+00') TO ('2025-06-23 00:00:00+00');

alter table public.aa_user_ops_info_p2025w25
    owner to postgres;

create table public.transaction_decode_p2025_06_08
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-08 00:00:00+00') TO ('2025-06-09 00:00:00+00');

alter table public.transaction_decode_p2025_06_08
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_08
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-08 00:00:00+00') TO ('2025-06-09 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_08
    owner to postgres;

create table public.transaction_decode_p2025_06_09
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-09 00:00:00+00') TO ('2025-06-10 00:00:00+00');

alter table public.transaction_decode_p2025_06_09
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_09
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-09 00:00:00+00') TO ('2025-06-10 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_09
    owner to postgres;

create table public.transaction_decode_p2025_06_10
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-10 00:00:00+00') TO ('2025-06-11 00:00:00+00');

alter table public.transaction_decode_p2025_06_10
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_10
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-10 00:00:00+00') TO ('2025-06-11 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_10
    owner to postgres;

create table public.transaction_decode_p2025_06_11
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-11 00:00:00+00') TO ('2025-06-12 00:00:00+00');

alter table public.transaction_decode_p2025_06_11
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_11
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-11 00:00:00+00') TO ('2025-06-12 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_11
    owner to postgres;

create table public.transaction_decode_p2025_06_12
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-12 00:00:00+00') TO ('2025-06-13 00:00:00+00');

alter table public.transaction_decode_p2025_06_12
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_12
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-12 00:00:00+00') TO ('2025-06-13 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_12
    owner to postgres;

create table public.transaction_decode_p2025_06_13
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-13 00:00:00+00') TO ('2025-06-14 00:00:00+00');

alter table public.transaction_decode_p2025_06_13
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_13
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-13 00:00:00+00') TO ('2025-06-14 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_13
    owner to postgres;

create table public.block_data_decode_p2025w26
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-06-23 00:00:00+00') TO ('2025-06-30 00:00:00+00');

alter table public.block_data_decode_p2025w26
    owner to postgres;

create table public.transaction_decode_p2025_06_14
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-14 00:00:00+00') TO ('2025-06-15 00:00:00+00');

alter table public.transaction_decode_p2025_06_14
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_14
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-14 00:00:00+00') TO ('2025-06-15 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_14
    owner to postgres;

create table public.aa_transaction_info_p2025w26
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-06-23 00:00:00+00') TO ('2025-06-30 00:00:00+00');

alter table public.aa_transaction_info_p2025w26
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w26
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-06-23 00:00:00+00') TO ('2025-06-30 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w26
    owner to postgres;

create table public.aa_user_ops_info_p2025w26
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-06-23 00:00:00+00') TO ('2025-06-30 00:00:00+00');

alter table public.aa_user_ops_info_p2025w26
    owner to postgres;

create table public.transaction_decode_p2025_06_15
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-15 00:00:00+00') TO ('2025-06-16 00:00:00+00');

alter table public.transaction_decode_p2025_06_15
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_15
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-15 00:00:00+00') TO ('2025-06-16 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_15
    owner to postgres;

create table public.transaction_decode_p2025_06_16
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-16 00:00:00+00') TO ('2025-06-17 00:00:00+00');

alter table public.transaction_decode_p2025_06_16
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_16
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-16 00:00:00+00') TO ('2025-06-17 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_16
    owner to postgres;

create table public.aa_block_info_p2025_11
    partition of public.aa_block_info
    FOR VALUES FROM ('2025-11-01 00:00:00+00') TO ('2025-12-01 00:00:00+00');

alter table public.aa_block_info_p2025_11
    owner to postgres;

create table public.transaction_decode_p2025_06_17
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-17 00:00:00+00') TO ('2025-06-18 00:00:00+00');

alter table public.transaction_decode_p2025_06_17
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_17
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-17 00:00:00+00') TO ('2025-06-18 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_17
    owner to postgres;

create table public.transaction_decode_p2025_06_18
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-18 00:00:00+00') TO ('2025-06-19 00:00:00+00');

alter table public.transaction_decode_p2025_06_18
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_18
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-18 00:00:00+00') TO ('2025-06-19 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_18
    owner to postgres;

create table public.transaction_decode_p2025_06_19
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-19 00:00:00+00') TO ('2025-06-20 00:00:00+00');

alter table public.transaction_decode_p2025_06_19
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_19
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-19 00:00:00+00') TO ('2025-06-20 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_19
    owner to postgres;

create table public.transaction_decode_p2025_06_20
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-20 00:00:00+00') TO ('2025-06-21 00:00:00+00');

alter table public.transaction_decode_p2025_06_20
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_20
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-20 00:00:00+00') TO ('2025-06-21 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_20
    owner to postgres;

create table public.block_data_decode_p2025w27
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-06-30 00:00:00+00') TO ('2025-07-07 00:00:00+00');

alter table public.block_data_decode_p2025w27
    owner to postgres;

create table public.transaction_decode_p2025_06_21
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-21 00:00:00+00') TO ('2025-06-22 00:00:00+00');

alter table public.transaction_decode_p2025_06_21
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_21
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-21 00:00:00+00') TO ('2025-06-22 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_21
    owner to postgres;

create table public.aa_transaction_info_p2025w27
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-06-30 00:00:00+00') TO ('2025-07-07 00:00:00+00');

alter table public.aa_transaction_info_p2025w27
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w27
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-06-30 00:00:00+00') TO ('2025-07-07 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w27
    owner to postgres;

create table public.aa_user_ops_info_p2025w27
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-06-30 00:00:00+00') TO ('2025-07-07 00:00:00+00');

alter table public.aa_user_ops_info_p2025w27
    owner to postgres;

create table public.transaction_decode_p2025_06_22
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-22 00:00:00+00') TO ('2025-06-23 00:00:00+00');

alter table public.transaction_decode_p2025_06_22
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_22
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-22 00:00:00+00') TO ('2025-06-23 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_22
    owner to postgres;

create table public.transaction_decode_p2025_06_23
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-23 00:00:00+00') TO ('2025-06-24 00:00:00+00');

alter table public.transaction_decode_p2025_06_23
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_23
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-23 00:00:00+00') TO ('2025-06-24 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_23
    owner to postgres;

create table public.transaction_decode_p2025_06_24
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-24 00:00:00+00') TO ('2025-06-25 00:00:00+00');

alter table public.transaction_decode_p2025_06_24
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_24
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-24 00:00:00+00') TO ('2025-06-25 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_24
    owner to postgres;

create table public.transaction_decode_p2025_06_25
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-25 00:00:00+00') TO ('2025-06-26 00:00:00+00');

alter table public.transaction_decode_p2025_06_25
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_25
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-25 00:00:00+00') TO ('2025-06-26 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_25
    owner to postgres;

create table public.transaction_decode_p2025_06_26
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-26 00:00:00+00') TO ('2025-06-27 00:00:00+00');

alter table public.transaction_decode_p2025_06_26
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_26
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-26 00:00:00+00') TO ('2025-06-27 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_26
    owner to postgres;

create table public.transaction_decode_p2025_06_27
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-27 00:00:00+00') TO ('2025-06-28 00:00:00+00');

alter table public.transaction_decode_p2025_06_27
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_27
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-27 00:00:00+00') TO ('2025-06-28 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_27
    owner to postgres;

create table public.block_data_decode_p2025w28
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-07-07 00:00:00+00') TO ('2025-07-14 00:00:00+00');

alter table public.block_data_decode_p2025w28
    owner to postgres;

create table public.transaction_decode_p2025_06_28
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-28 00:00:00+00') TO ('2025-06-29 00:00:00+00');

alter table public.transaction_decode_p2025_06_28
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_28
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-28 00:00:00+00') TO ('2025-06-29 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_28
    owner to postgres;

create table public.aa_transaction_info_p2025w28
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-07-07 00:00:00+00') TO ('2025-07-14 00:00:00+00');

alter table public.aa_transaction_info_p2025w28
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w28
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-07-07 00:00:00+00') TO ('2025-07-14 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w28
    owner to postgres;

create table public.aa_user_ops_info_p2025w28
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-07-07 00:00:00+00') TO ('2025-07-14 00:00:00+00');

alter table public.aa_user_ops_info_p2025w28
    owner to postgres;

create table public.transaction_decode_p2025_06_29
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-29 00:00:00+00') TO ('2025-06-30 00:00:00+00');

alter table public.transaction_decode_p2025_06_29
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_29
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-29 00:00:00+00') TO ('2025-06-30 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_29
    owner to postgres;

create table public.transaction_decode_p2025_06_30
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-06-30 00:00:00+00') TO ('2025-07-01 00:00:00+00');

alter table public.transaction_decode_p2025_06_30
    owner to postgres;

create table public.transaction_receipt_decode_p2025_06_30
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-06-30 00:00:00+00') TO ('2025-07-01 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_06_30
    owner to postgres;

create table public.transaction_decode_p2025_07_01
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-01 00:00:00+00') TO ('2025-07-02 00:00:00+00');

alter table public.transaction_decode_p2025_07_01
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_01
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-01 00:00:00+00') TO ('2025-07-02 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_01
    owner to postgres;

create table public.transaction_decode_p2025_07_02
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-02 00:00:00+00') TO ('2025-07-03 00:00:00+00');

alter table public.transaction_decode_p2025_07_02
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_02
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-02 00:00:00+00') TO ('2025-07-03 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_02
    owner to postgres;

create table public.transaction_decode_p2025_07_03
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-03 00:00:00+00') TO ('2025-07-04 00:00:00+00');

alter table public.transaction_decode_p2025_07_03
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_03
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-03 00:00:00+00') TO ('2025-07-04 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_03
    owner to postgres;

create table public.transaction_decode_p2025_07_04
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-04 00:00:00+00') TO ('2025-07-05 00:00:00+00');

alter table public.transaction_decode_p2025_07_04
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_04
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-04 00:00:00+00') TO ('2025-07-05 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_04
    owner to postgres;

create table public.transaction_decode_p2025_07_05
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-05 00:00:00+00') TO ('2025-07-06 00:00:00+00');

alter table public.transaction_decode_p2025_07_05
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_05
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-05 00:00:00+00') TO ('2025-07-06 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_05
    owner to postgres;

create table public.block_data_decode_p2025w29
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-07-14 00:00:00+00') TO ('2025-07-21 00:00:00+00');

alter table public.block_data_decode_p2025w29
    owner to postgres;

create table public.transaction_decode_p2025_07_06
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-06 00:00:00+00') TO ('2025-07-07 00:00:00+00');

alter table public.transaction_decode_p2025_07_06
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_06
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-06 00:00:00+00') TO ('2025-07-07 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_06
    owner to postgres;

create table public.aa_transaction_info_p2025w29
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-07-14 00:00:00+00') TO ('2025-07-21 00:00:00+00');

alter table public.aa_transaction_info_p2025w29
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w29
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-07-14 00:00:00+00') TO ('2025-07-21 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w29
    owner to postgres;

create table public.aa_user_ops_info_p2025w29
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-07-14 00:00:00+00') TO ('2025-07-21 00:00:00+00');

alter table public.aa_user_ops_info_p2025w29
    owner to postgres;

create table public.transaction_decode_p2025_07_07
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-07 00:00:00+00') TO ('2025-07-08 00:00:00+00');

alter table public.transaction_decode_p2025_07_07
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_07
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-07 00:00:00+00') TO ('2025-07-08 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_07
    owner to postgres;

create table public.transaction_decode_p2025_07_08
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-08 00:00:00+00') TO ('2025-07-09 00:00:00+00');

alter table public.transaction_decode_p2025_07_08
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_08
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-08 00:00:00+00') TO ('2025-07-09 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_08
    owner to postgres;

create table public.transaction_decode_p2025_07_09
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-09 00:00:00+00') TO ('2025-07-10 00:00:00+00');

alter table public.transaction_decode_p2025_07_09
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_09
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-09 00:00:00+00') TO ('2025-07-10 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_09
    owner to postgres;

create table public.transaction_decode_p2025_07_10
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-10 00:00:00+00') TO ('2025-07-11 00:00:00+00');

alter table public.transaction_decode_p2025_07_10
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_10
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-10 00:00:00+00') TO ('2025-07-11 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_10
    owner to postgres;

create table public.transaction_decode_p2025_07_11
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-11 00:00:00+00') TO ('2025-07-12 00:00:00+00');

alter table public.transaction_decode_p2025_07_11
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_11
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-11 00:00:00+00') TO ('2025-07-12 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_11
    owner to postgres;

create table public.transaction_decode_p2025_07_12
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-12 00:00:00+00') TO ('2025-07-13 00:00:00+00');

alter table public.transaction_decode_p2025_07_12
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_12
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-12 00:00:00+00') TO ('2025-07-13 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_12
    owner to postgres;

create table public.block_data_decode_p2025w30
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-07-21 00:00:00+00') TO ('2025-07-28 00:00:00+00');

alter table public.block_data_decode_p2025w30
    owner to postgres;

create table public.transaction_decode_p2025_07_13
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-13 00:00:00+00') TO ('2025-07-14 00:00:00+00');

alter table public.transaction_decode_p2025_07_13
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_13
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-13 00:00:00+00') TO ('2025-07-14 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_13
    owner to postgres;

create table public.aa_transaction_info_p2025w30
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-07-21 00:00:00+00') TO ('2025-07-28 00:00:00+00');

alter table public.aa_transaction_info_p2025w30
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w30
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-07-21 00:00:00+00') TO ('2025-07-28 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w30
    owner to postgres;

create table public.aa_user_ops_info_p2025w30
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-07-21 00:00:00+00') TO ('2025-07-28 00:00:00+00');

alter table public.aa_user_ops_info_p2025w30
    owner to postgres;

create table public.transaction_decode_p2025_07_14
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-14 00:00:00+00') TO ('2025-07-15 00:00:00+00');

alter table public.transaction_decode_p2025_07_14
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_14
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-14 00:00:00+00') TO ('2025-07-15 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_14
    owner to postgres;

create table public.transaction_decode_p2025_07_15
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-15 00:00:00+00') TO ('2025-07-16 00:00:00+00');

alter table public.transaction_decode_p2025_07_15
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_15
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-15 00:00:00+00') TO ('2025-07-16 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_15
    owner to postgres;

create table public.transaction_decode_p2025_07_16
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-16 00:00:00+00') TO ('2025-07-17 00:00:00+00');

alter table public.transaction_decode_p2025_07_16
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_16
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-16 00:00:00+00') TO ('2025-07-17 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_16
    owner to postgres;

create table public.aa_block_info_p2025_12
    partition of public.aa_block_info
    FOR VALUES FROM ('2025-12-01 00:00:00+00') TO ('2026-01-01 00:00:00+00');

alter table public.aa_block_info_p2025_12
    owner to postgres;

create table public.transaction_decode_p2025_07_17
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-17 00:00:00+00') TO ('2025-07-18 00:00:00+00');

alter table public.transaction_decode_p2025_07_17
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_17
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-17 00:00:00+00') TO ('2025-07-18 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_17
    owner to postgres;

create table public.transaction_decode_p2025_07_18
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-18 00:00:00+00') TO ('2025-07-19 00:00:00+00');

alter table public.transaction_decode_p2025_07_18
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_18
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-18 00:00:00+00') TO ('2025-07-19 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_18
    owner to postgres;

create table public.transaction_decode_p2025_07_19
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-19 00:00:00+00') TO ('2025-07-20 00:00:00+00');

alter table public.transaction_decode_p2025_07_19
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_19
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-19 00:00:00+00') TO ('2025-07-20 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_19
    owner to postgres;

create table public.block_data_decode_p2025w31
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-07-28 00:00:00+00') TO ('2025-08-04 00:00:00+00');

alter table public.block_data_decode_p2025w31
    owner to postgres;

create table public.transaction_decode_p2025_07_20
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-20 00:00:00+00') TO ('2025-07-21 00:00:00+00');

alter table public.transaction_decode_p2025_07_20
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_20
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-20 00:00:00+00') TO ('2025-07-21 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_20
    owner to postgres;

create table public.aa_transaction_info_p2025w31
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-07-28 00:00:00+00') TO ('2025-08-04 00:00:00+00');

alter table public.aa_transaction_info_p2025w31
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w31
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-07-28 00:00:00+00') TO ('2025-08-04 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w31
    owner to postgres;

create table public.aa_user_ops_info_p2025w31
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-07-28 00:00:00+00') TO ('2025-08-04 00:00:00+00');

alter table public.aa_user_ops_info_p2025w31
    owner to postgres;

create table public.transaction_decode_p2025_07_21
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-21 00:00:00+00') TO ('2025-07-22 00:00:00+00');

alter table public.transaction_decode_p2025_07_21
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_21
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-21 00:00:00+00') TO ('2025-07-22 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_21
    owner to postgres;

create table public.transaction_decode_p2025_07_22
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-22 00:00:00+00') TO ('2025-07-23 00:00:00+00');

alter table public.transaction_decode_p2025_07_22
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_22
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-22 00:00:00+00') TO ('2025-07-23 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_22
    owner to postgres;

create table public.transaction_decode_p2025_07_23
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-23 00:00:00+00') TO ('2025-07-24 00:00:00+00');

alter table public.transaction_decode_p2025_07_23
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_23
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-23 00:00:00+00') TO ('2025-07-24 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_23
    owner to postgres;

create table public.transaction_decode_p2025_07_24
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-24 00:00:00+00') TO ('2025-07-25 00:00:00+00');

alter table public.transaction_decode_p2025_07_24
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_24
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-24 00:00:00+00') TO ('2025-07-25 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_24
    owner to postgres;

create table public.transaction_decode_p2025_07_25
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-25 00:00:00+00') TO ('2025-07-26 00:00:00+00');

alter table public.transaction_decode_p2025_07_25
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_25
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-25 00:00:00+00') TO ('2025-07-26 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_25
    owner to postgres;

create table public.transaction_decode_p2025_07_26
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-26 00:00:00+00') TO ('2025-07-27 00:00:00+00');

alter table public.transaction_decode_p2025_07_26
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_26
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-26 00:00:00+00') TO ('2025-07-27 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_26
    owner to postgres;

create table public.block_data_decode_p2025w32
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-08-04 00:00:00+00') TO ('2025-08-11 00:00:00+00');

alter table public.block_data_decode_p2025w32
    owner to postgres;

create table public.transaction_decode_p2025_07_27
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-27 00:00:00+00') TO ('2025-07-28 00:00:00+00');

alter table public.transaction_decode_p2025_07_27
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_27
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-27 00:00:00+00') TO ('2025-07-28 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_27
    owner to postgres;

create table public.aa_transaction_info_p2025w32
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-08-04 00:00:00+00') TO ('2025-08-11 00:00:00+00');

alter table public.aa_transaction_info_p2025w32
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w32
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-08-04 00:00:00+00') TO ('2025-08-11 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w32
    owner to postgres;

create table public.aa_user_ops_info_p2025w32
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-08-04 00:00:00+00') TO ('2025-08-11 00:00:00+00');

alter table public.aa_user_ops_info_p2025w32
    owner to postgres;

create table public.transaction_decode_p2025_07_28
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-28 00:00:00+00') TO ('2025-07-29 00:00:00+00');

alter table public.transaction_decode_p2025_07_28
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_28
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-28 00:00:00+00') TO ('2025-07-29 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_28
    owner to postgres;

create table public.transaction_decode_p2025_07_29
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-29 00:00:00+00') TO ('2025-07-30 00:00:00+00');

alter table public.transaction_decode_p2025_07_29
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_29
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-29 00:00:00+00') TO ('2025-07-30 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_29
    owner to postgres;

create table public.transaction_decode_p2025_07_30
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-30 00:00:00+00') TO ('2025-07-31 00:00:00+00');

alter table public.transaction_decode_p2025_07_30
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_30
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-30 00:00:00+00') TO ('2025-07-31 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_30
    owner to postgres;

create table public.transaction_decode_p2025_07_31
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-07-31 00:00:00+00') TO ('2025-08-01 00:00:00+00');

alter table public.transaction_decode_p2025_07_31
    owner to postgres;

create table public.transaction_receipt_decode_p2025_07_31
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-07-31 00:00:00+00') TO ('2025-08-01 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_07_31
    owner to postgres;

create table public.transaction_decode_p2025_08_01
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-01 00:00:00+00') TO ('2025-08-02 00:00:00+00');

alter table public.transaction_decode_p2025_08_01
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_01
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-01 00:00:00+00') TO ('2025-08-02 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_01
    owner to postgres;

create table public.block_data_decode_p2025w33
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-08-11 00:00:00+00') TO ('2025-08-18 00:00:00+00');

alter table public.block_data_decode_p2025w33
    owner to postgres;

create table public.transaction_decode_p2025_08_02
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-02 00:00:00+00') TO ('2025-08-03 00:00:00+00');

alter table public.transaction_decode_p2025_08_02
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_02
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-02 00:00:00+00') TO ('2025-08-03 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_02
    owner to postgres;

create table public.aa_transaction_info_p2025w33
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-08-11 00:00:00+00') TO ('2025-08-18 00:00:00+00');

alter table public.aa_transaction_info_p2025w33
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w33
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-08-11 00:00:00+00') TO ('2025-08-18 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w33
    owner to postgres;

create table public.aa_user_ops_info_p2025w33
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-08-11 00:00:00+00') TO ('2025-08-18 00:00:00+00');

alter table public.aa_user_ops_info_p2025w33
    owner to postgres;

create table public.transaction_decode_p2025_08_03
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-03 00:00:00+00') TO ('2025-08-04 00:00:00+00');

alter table public.transaction_decode_p2025_08_03
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_03
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-03 00:00:00+00') TO ('2025-08-04 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_03
    owner to postgres;

create table public.transaction_decode_p2025_08_04
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-04 00:00:00+00') TO ('2025-08-05 00:00:00+00');

alter table public.transaction_decode_p2025_08_04
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_04
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-04 00:00:00+00') TO ('2025-08-05 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_04
    owner to postgres;

create table public.transaction_decode_p2025_08_05
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-05 00:00:00+00') TO ('2025-08-06 00:00:00+00');

alter table public.transaction_decode_p2025_08_05
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_05
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-05 00:00:00+00') TO ('2025-08-06 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_05
    owner to postgres;

create table public.transaction_decode_p2025_08_06
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-06 00:00:00+00') TO ('2025-08-07 00:00:00+00');

alter table public.transaction_decode_p2025_08_06
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_06
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-06 00:00:00+00') TO ('2025-08-07 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_06
    owner to postgres;

create table public.transaction_decode_p2025_08_07
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-07 00:00:00+00') TO ('2025-08-08 00:00:00+00');

alter table public.transaction_decode_p2025_08_07
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_07
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-07 00:00:00+00') TO ('2025-08-08 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_07
    owner to postgres;

create table public.transaction_decode_p2025_08_08
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-08 00:00:00+00') TO ('2025-08-09 00:00:00+00');

alter table public.transaction_decode_p2025_08_08
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_08
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-08 00:00:00+00') TO ('2025-08-09 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_08
    owner to postgres;

create table public.block_data_decode_p2025w34
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-08-18 00:00:00+00') TO ('2025-08-25 00:00:00+00');

alter table public.block_data_decode_p2025w34
    owner to postgres;

create table public.transaction_decode_p2025_08_09
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-09 00:00:00+00') TO ('2025-08-10 00:00:00+00');

alter table public.transaction_decode_p2025_08_09
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_09
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-09 00:00:00+00') TO ('2025-08-10 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_09
    owner to postgres;

create table public.aa_transaction_info_p2025w34
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-08-18 00:00:00+00') TO ('2025-08-25 00:00:00+00');

alter table public.aa_transaction_info_p2025w34
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w34
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-08-18 00:00:00+00') TO ('2025-08-25 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w34
    owner to postgres;

create table public.aa_user_ops_info_p2025w34
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-08-18 00:00:00+00') TO ('2025-08-25 00:00:00+00');

alter table public.aa_user_ops_info_p2025w34
    owner to postgres;

create table public.block_data_decode_p2025w35
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-08-25 00:00:00+00') TO ('2025-09-01 00:00:00+00');

alter table public.block_data_decode_p2025w35
    owner to postgres;

create table public.block_data_decode_p2025w36
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-09-01 00:00:00+00') TO ('2025-09-08 00:00:00+00');

alter table public.block_data_decode_p2025w36
    owner to postgres;

create table public.block_data_decode_p2025w37
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-09-08 00:00:00+00') TO ('2025-09-15 00:00:00+00');

alter table public.block_data_decode_p2025w37
    owner to postgres;

create table public.block_data_decode_p2025w38
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-09-15 00:00:00+00') TO ('2025-09-22 00:00:00+00');

alter table public.block_data_decode_p2025w38
    owner to postgres;

create table public.block_data_decode_p2025w39
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-09-22 00:00:00+00') TO ('2025-09-29 00:00:00+00');

alter table public.block_data_decode_p2025w39
    owner to postgres;

create table public.block_data_decode_p2025w40
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-09-29 00:00:00+00') TO ('2025-10-06 00:00:00+00');

alter table public.block_data_decode_p2025w40
    owner to postgres;

create table public.block_data_decode_p2025w41
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-10-06 00:00:00+00') TO ('2025-10-13 00:00:00+00');

alter table public.block_data_decode_p2025w41
    owner to postgres;

create table public.block_data_decode_p2025w42
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-10-13 00:00:00+00') TO ('2025-10-20 00:00:00+00');

alter table public.block_data_decode_p2025w42
    owner to postgres;

create table public.block_data_decode_p2025w43
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-10-20 00:00:00+00') TO ('2025-10-27 00:00:00+00');

alter table public.block_data_decode_p2025w43
    owner to postgres;

create table public.block_data_decode_p2025w44
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-10-27 00:00:00+00') TO ('2025-11-03 00:00:00+00');

alter table public.block_data_decode_p2025w44
    owner to postgres;

create table public.block_data_decode_p2025w45
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-11-03 00:00:00+00') TO ('2025-11-10 00:00:00+00');

alter table public.block_data_decode_p2025w45
    owner to postgres;

create table public.block_data_decode_p2025w46
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-11-10 00:00:00+00') TO ('2025-11-17 00:00:00+00');

alter table public.block_data_decode_p2025w46
    owner to postgres;

create table public.block_data_decode_p2025w47
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-11-17 00:00:00+00') TO ('2025-11-24 00:00:00+00');

alter table public.block_data_decode_p2025w47
    owner to postgres;

create table public.block_data_decode_p2025w48
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-11-24 00:00:00+00') TO ('2025-12-01 00:00:00+00');

alter table public.block_data_decode_p2025w48
    owner to postgres;

create table public.block_data_decode_p2025w49
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-12-01 00:00:00+00') TO ('2025-12-08 00:00:00+00');

alter table public.block_data_decode_p2025w49
    owner to postgres;

create table public.block_data_decode_p2025w50
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-12-08 00:00:00+00') TO ('2025-12-15 00:00:00+00');

alter table public.block_data_decode_p2025w50
    owner to postgres;

create table public.block_data_decode_p2025w51
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-12-15 00:00:00+00') TO ('2025-12-22 00:00:00+00');

alter table public.block_data_decode_p2025w51
    owner to postgres;

create table public.block_data_decode_p2025w52
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-12-22 00:00:00+00') TO ('2025-12-29 00:00:00+00');

alter table public.block_data_decode_p2025w52
    owner to postgres;

create table public.aa_block_info_p2026_01
    partition of public.aa_block_info
    FOR VALUES FROM ('2026-01-01 00:00:00+00') TO ('2026-02-01 00:00:00+00');

alter table public.aa_block_info_p2026_01
    owner to postgres;

create table public.aa_block_info_p2026_02
    partition of public.aa_block_info
    FOR VALUES FROM ('2026-02-01 00:00:00+00') TO ('2026-03-01 00:00:00+00');

alter table public.aa_block_info_p2026_02
    owner to postgres;

create table public.aa_block_info_p2026_03
    partition of public.aa_block_info
    FOR VALUES FROM ('2026-03-01 00:00:00+00') TO ('2026-04-01 00:00:00+00');

alter table public.aa_block_info_p2026_03
    owner to postgres;

create table public.aa_block_info_p2026_04
    partition of public.aa_block_info
    FOR VALUES FROM ('2026-04-01 00:00:00+00') TO ('2026-05-01 00:00:00+00');

alter table public.aa_block_info_p2026_04
    owner to postgres;

create table public.aa_block_info_p2026_05
    partition of public.aa_block_info
    FOR VALUES FROM ('2026-05-01 00:00:00+00') TO ('2026-06-01 00:00:00+00');

alter table public.aa_block_info_p2026_05
    owner to postgres;

create table public.transaction_decode_p2025_08_10
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-10 00:00:00+00') TO ('2025-08-11 00:00:00+00');

alter table public.transaction_decode_p2025_08_10
    owner to postgres;

create table public.transaction_decode_p2025_08_11
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-11 00:00:00+00') TO ('2025-08-12 00:00:00+00');

alter table public.transaction_decode_p2025_08_11
    owner to postgres;

create table public.transaction_decode_p2025_08_12
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-12 00:00:00+00') TO ('2025-08-13 00:00:00+00');

alter table public.transaction_decode_p2025_08_12
    owner to postgres;

create table public.transaction_decode_p2025_08_13
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-13 00:00:00+00') TO ('2025-08-14 00:00:00+00');

alter table public.transaction_decode_p2025_08_13
    owner to postgres;

create table public.transaction_decode_p2025_08_14
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-14 00:00:00+00') TO ('2025-08-15 00:00:00+00');

alter table public.transaction_decode_p2025_08_14
    owner to postgres;

create table public.transaction_decode_p2025_08_15
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-15 00:00:00+00') TO ('2025-08-16 00:00:00+00');

alter table public.transaction_decode_p2025_08_15
    owner to postgres;

create table public.transaction_decode_p2025_08_16
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-16 00:00:00+00') TO ('2025-08-17 00:00:00+00');

alter table public.transaction_decode_p2025_08_16
    owner to postgres;

create table public.transaction_decode_p2025_08_17
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-17 00:00:00+00') TO ('2025-08-18 00:00:00+00');

alter table public.transaction_decode_p2025_08_17
    owner to postgres;

create table public.transaction_decode_p2025_08_18
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-18 00:00:00+00') TO ('2025-08-19 00:00:00+00');

alter table public.transaction_decode_p2025_08_18
    owner to postgres;

create table public.transaction_decode_p2025_08_19
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-19 00:00:00+00') TO ('2025-08-20 00:00:00+00');

alter table public.transaction_decode_p2025_08_19
    owner to postgres;

create table public.transaction_decode_p2025_08_20
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-20 00:00:00+00') TO ('2025-08-21 00:00:00+00');

alter table public.transaction_decode_p2025_08_20
    owner to postgres;

create table public.transaction_decode_p2025_08_21
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-21 00:00:00+00') TO ('2025-08-22 00:00:00+00');

alter table public.transaction_decode_p2025_08_21
    owner to postgres;

create table public.transaction_decode_p2025_08_22
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-22 00:00:00+00') TO ('2025-08-23 00:00:00+00');

alter table public.transaction_decode_p2025_08_22
    owner to postgres;

create table public.transaction_decode_p2025_08_23
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-23 00:00:00+00') TO ('2025-08-24 00:00:00+00');

alter table public.transaction_decode_p2025_08_23
    owner to postgres;

create table public.transaction_decode_p2025_08_24
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-24 00:00:00+00') TO ('2025-08-25 00:00:00+00');

alter table public.transaction_decode_p2025_08_24
    owner to postgres;

create table public.transaction_decode_p2025_08_25
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-25 00:00:00+00') TO ('2025-08-26 00:00:00+00');

alter table public.transaction_decode_p2025_08_25
    owner to postgres;

create table public.transaction_decode_p2025_08_26
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-26 00:00:00+00') TO ('2025-08-27 00:00:00+00');

alter table public.transaction_decode_p2025_08_26
    owner to postgres;

create table public.transaction_decode_p2025_08_27
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-27 00:00:00+00') TO ('2025-08-28 00:00:00+00');

alter table public.transaction_decode_p2025_08_27
    owner to postgres;

create table public.transaction_decode_p2025_08_28
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-28 00:00:00+00') TO ('2025-08-29 00:00:00+00');

alter table public.transaction_decode_p2025_08_28
    owner to postgres;

create table public.transaction_decode_p2025_08_29
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-29 00:00:00+00') TO ('2025-08-30 00:00:00+00');

alter table public.transaction_decode_p2025_08_29
    owner to postgres;

create table public.transaction_decode_p2025_08_30
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-30 00:00:00+00') TO ('2025-08-31 00:00:00+00');

alter table public.transaction_decode_p2025_08_30
    owner to postgres;

create table public.transaction_decode_p2025_08_31
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-08-31 00:00:00+00') TO ('2025-09-01 00:00:00+00');

alter table public.transaction_decode_p2025_08_31
    owner to postgres;

create table public.transaction_decode_p2025_09_01
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-01 00:00:00+00') TO ('2025-09-02 00:00:00+00');

alter table public.transaction_decode_p2025_09_01
    owner to postgres;

create table public.transaction_decode_p2025_09_02
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-02 00:00:00+00') TO ('2025-09-03 00:00:00+00');

alter table public.transaction_decode_p2025_09_02
    owner to postgres;

create table public.transaction_decode_p2025_09_03
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-03 00:00:00+00') TO ('2025-09-04 00:00:00+00');

alter table public.transaction_decode_p2025_09_03
    owner to postgres;

create table public.transaction_decode_p2025_09_04
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-04 00:00:00+00') TO ('2025-09-05 00:00:00+00');

alter table public.transaction_decode_p2025_09_04
    owner to postgres;

create table public.transaction_decode_p2025_09_05
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-05 00:00:00+00') TO ('2025-09-06 00:00:00+00');

alter table public.transaction_decode_p2025_09_05
    owner to postgres;

create table public.transaction_decode_p2025_09_06
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-06 00:00:00+00') TO ('2025-09-07 00:00:00+00');

alter table public.transaction_decode_p2025_09_06
    owner to postgres;

create table public.transaction_decode_p2025_09_07
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-07 00:00:00+00') TO ('2025-09-08 00:00:00+00');

alter table public.transaction_decode_p2025_09_07
    owner to postgres;

create table public.transaction_decode_p2025_09_08
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-08 00:00:00+00') TO ('2025-09-09 00:00:00+00');

alter table public.transaction_decode_p2025_09_08
    owner to postgres;

create table public.transaction_decode_p2025_09_09
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-09 00:00:00+00') TO ('2025-09-10 00:00:00+00');

alter table public.transaction_decode_p2025_09_09
    owner to postgres;

create table public.transaction_decode_p2025_09_10
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-10 00:00:00+00') TO ('2025-09-11 00:00:00+00');

alter table public.transaction_decode_p2025_09_10
    owner to postgres;

create table public.transaction_decode_p2025_09_11
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-11 00:00:00+00') TO ('2025-09-12 00:00:00+00');

alter table public.transaction_decode_p2025_09_11
    owner to postgres;

create table public.transaction_decode_p2025_09_12
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-12 00:00:00+00') TO ('2025-09-13 00:00:00+00');

alter table public.transaction_decode_p2025_09_12
    owner to postgres;

create table public.transaction_decode_p2025_09_13
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-13 00:00:00+00') TO ('2025-09-14 00:00:00+00');

alter table public.transaction_decode_p2025_09_13
    owner to postgres;

create table public.transaction_decode_p2025_09_14
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-14 00:00:00+00') TO ('2025-09-15 00:00:00+00');

alter table public.transaction_decode_p2025_09_14
    owner to postgres;

create table public.transaction_decode_p2025_09_15
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-15 00:00:00+00') TO ('2025-09-16 00:00:00+00');

alter table public.transaction_decode_p2025_09_15
    owner to postgres;

create table public.transaction_decode_p2025_09_16
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-16 00:00:00+00') TO ('2025-09-17 00:00:00+00');

alter table public.transaction_decode_p2025_09_16
    owner to postgres;

create table public.transaction_decode_p2025_09_17
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-17 00:00:00+00') TO ('2025-09-18 00:00:00+00');

alter table public.transaction_decode_p2025_09_17
    owner to postgres;

create table public.transaction_decode_p2025_09_18
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-18 00:00:00+00') TO ('2025-09-19 00:00:00+00');

alter table public.transaction_decode_p2025_09_18
    owner to postgres;

create table public.transaction_decode_p2025_09_19
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-19 00:00:00+00') TO ('2025-09-20 00:00:00+00');

alter table public.transaction_decode_p2025_09_19
    owner to postgres;

create table public.transaction_decode_p2025_09_20
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-20 00:00:00+00') TO ('2025-09-21 00:00:00+00');

alter table public.transaction_decode_p2025_09_20
    owner to postgres;

create table public.transaction_decode_p2025_09_21
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-21 00:00:00+00') TO ('2025-09-22 00:00:00+00');

alter table public.transaction_decode_p2025_09_21
    owner to postgres;

create table public.transaction_decode_p2025_09_22
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-22 00:00:00+00') TO ('2025-09-23 00:00:00+00');

alter table public.transaction_decode_p2025_09_22
    owner to postgres;

create table public.transaction_decode_p2025_09_23
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-23 00:00:00+00') TO ('2025-09-24 00:00:00+00');

alter table public.transaction_decode_p2025_09_23
    owner to postgres;

create table public.transaction_decode_p2025_09_24
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-24 00:00:00+00') TO ('2025-09-25 00:00:00+00');

alter table public.transaction_decode_p2025_09_24
    owner to postgres;

create table public.transaction_decode_p2025_09_25
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-25 00:00:00+00') TO ('2025-09-26 00:00:00+00');

alter table public.transaction_decode_p2025_09_25
    owner to postgres;

create table public.transaction_decode_p2025_09_26
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-26 00:00:00+00') TO ('2025-09-27 00:00:00+00');

alter table public.transaction_decode_p2025_09_26
    owner to postgres;

create table public.transaction_decode_p2025_09_27
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-27 00:00:00+00') TO ('2025-09-28 00:00:00+00');

alter table public.transaction_decode_p2025_09_27
    owner to postgres;

create table public.transaction_decode_p2025_09_28
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-28 00:00:00+00') TO ('2025-09-29 00:00:00+00');

alter table public.transaction_decode_p2025_09_28
    owner to postgres;

create table public.transaction_decode_p2025_09_29
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-29 00:00:00+00') TO ('2025-09-30 00:00:00+00');

alter table public.transaction_decode_p2025_09_29
    owner to postgres;

create table public.transaction_decode_p2025_09_30
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-09-30 00:00:00+00') TO ('2025-10-01 00:00:00+00');

alter table public.transaction_decode_p2025_09_30
    owner to postgres;

create table public.transaction_decode_p2025_10_01
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-01 00:00:00+00') TO ('2025-10-02 00:00:00+00');

alter table public.transaction_decode_p2025_10_01
    owner to postgres;

create table public.transaction_decode_p2025_10_02
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-02 00:00:00+00') TO ('2025-10-03 00:00:00+00');

alter table public.transaction_decode_p2025_10_02
    owner to postgres;

create table public.transaction_decode_p2025_10_03
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-03 00:00:00+00') TO ('2025-10-04 00:00:00+00');

alter table public.transaction_decode_p2025_10_03
    owner to postgres;

create table public.transaction_decode_p2025_10_04
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-04 00:00:00+00') TO ('2025-10-05 00:00:00+00');

alter table public.transaction_decode_p2025_10_04
    owner to postgres;

create table public.transaction_decode_p2025_10_05
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-05 00:00:00+00') TO ('2025-10-06 00:00:00+00');

alter table public.transaction_decode_p2025_10_05
    owner to postgres;

create table public.transaction_decode_p2025_10_06
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-06 00:00:00+00') TO ('2025-10-07 00:00:00+00');

alter table public.transaction_decode_p2025_10_06
    owner to postgres;

create table public.transaction_decode_p2025_10_07
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-07 00:00:00+00') TO ('2025-10-08 00:00:00+00');

alter table public.transaction_decode_p2025_10_07
    owner to postgres;

create table public.transaction_decode_p2025_10_08
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-08 00:00:00+00') TO ('2025-10-09 00:00:00+00');

alter table public.transaction_decode_p2025_10_08
    owner to postgres;

create table public.transaction_decode_p2025_10_09
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-09 00:00:00+00') TO ('2025-10-10 00:00:00+00');

alter table public.transaction_decode_p2025_10_09
    owner to postgres;

create table public.transaction_decode_p2025_10_10
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-10 00:00:00+00') TO ('2025-10-11 00:00:00+00');

alter table public.transaction_decode_p2025_10_10
    owner to postgres;

create table public.transaction_decode_p2025_10_11
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-11 00:00:00+00') TO ('2025-10-12 00:00:00+00');

alter table public.transaction_decode_p2025_10_11
    owner to postgres;

create table public.transaction_decode_p2025_10_12
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-12 00:00:00+00') TO ('2025-10-13 00:00:00+00');

alter table public.transaction_decode_p2025_10_12
    owner to postgres;

create table public.transaction_decode_p2025_10_13
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-13 00:00:00+00') TO ('2025-10-14 00:00:00+00');

alter table public.transaction_decode_p2025_10_13
    owner to postgres;

create table public.transaction_decode_p2025_10_14
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-14 00:00:00+00') TO ('2025-10-15 00:00:00+00');

alter table public.transaction_decode_p2025_10_14
    owner to postgres;

create table public.transaction_decode_p2025_10_15
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-15 00:00:00+00') TO ('2025-10-16 00:00:00+00');

alter table public.transaction_decode_p2025_10_15
    owner to postgres;

create table public.transaction_decode_p2025_10_16
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-16 00:00:00+00') TO ('2025-10-17 00:00:00+00');

alter table public.transaction_decode_p2025_10_16
    owner to postgres;

create table public.transaction_decode_p2025_10_17
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-17 00:00:00+00') TO ('2025-10-18 00:00:00+00');

alter table public.transaction_decode_p2025_10_17
    owner to postgres;

create table public.transaction_decode_p2025_10_18
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-18 00:00:00+00') TO ('2025-10-19 00:00:00+00');

alter table public.transaction_decode_p2025_10_18
    owner to postgres;

create table public.transaction_decode_p2025_10_19
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-19 00:00:00+00') TO ('2025-10-20 00:00:00+00');

alter table public.transaction_decode_p2025_10_19
    owner to postgres;

create table public.transaction_decode_p2025_10_20
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-20 00:00:00+00') TO ('2025-10-21 00:00:00+00');

alter table public.transaction_decode_p2025_10_20
    owner to postgres;

create table public.transaction_decode_p2025_10_21
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-21 00:00:00+00') TO ('2025-10-22 00:00:00+00');

alter table public.transaction_decode_p2025_10_21
    owner to postgres;

create table public.transaction_decode_p2025_10_22
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-22 00:00:00+00') TO ('2025-10-23 00:00:00+00');

alter table public.transaction_decode_p2025_10_22
    owner to postgres;

create table public.transaction_decode_p2025_10_23
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-23 00:00:00+00') TO ('2025-10-24 00:00:00+00');

alter table public.transaction_decode_p2025_10_23
    owner to postgres;

create table public.transaction_decode_p2025_10_24
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-24 00:00:00+00') TO ('2025-10-25 00:00:00+00');

alter table public.transaction_decode_p2025_10_24
    owner to postgres;

create table public.transaction_decode_p2025_10_25
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-25 00:00:00+00') TO ('2025-10-26 00:00:00+00');

alter table public.transaction_decode_p2025_10_25
    owner to postgres;

create table public.transaction_decode_p2025_10_26
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-26 00:00:00+00') TO ('2025-10-27 00:00:00+00');

alter table public.transaction_decode_p2025_10_26
    owner to postgres;

create table public.transaction_decode_p2025_10_27
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-27 00:00:00+00') TO ('2025-10-28 00:00:00+00');

alter table public.transaction_decode_p2025_10_27
    owner to postgres;

create table public.transaction_decode_p2025_10_28
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-28 00:00:00+00') TO ('2025-10-29 00:00:00+00');

alter table public.transaction_decode_p2025_10_28
    owner to postgres;

create table public.transaction_decode_p2025_10_29
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-29 00:00:00+00') TO ('2025-10-30 00:00:00+00');

alter table public.transaction_decode_p2025_10_29
    owner to postgres;

create table public.transaction_decode_p2025_10_30
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-30 00:00:00+00') TO ('2025-10-31 00:00:00+00');

alter table public.transaction_decode_p2025_10_30
    owner to postgres;

create table public.transaction_decode_p2025_10_31
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-10-31 00:00:00+00') TO ('2025-11-01 00:00:00+00');

alter table public.transaction_decode_p2025_10_31
    owner to postgres;

create table public.transaction_decode_p2025_11_01
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-01 00:00:00+00') TO ('2025-11-02 00:00:00+00');

alter table public.transaction_decode_p2025_11_01
    owner to postgres;

create table public.transaction_decode_p2025_11_02
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-02 00:00:00+00') TO ('2025-11-03 00:00:00+00');

alter table public.transaction_decode_p2025_11_02
    owner to postgres;

create table public.transaction_decode_p2025_11_03
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-03 00:00:00+00') TO ('2025-11-04 00:00:00+00');

alter table public.transaction_decode_p2025_11_03
    owner to postgres;

create table public.transaction_decode_p2025_11_04
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-04 00:00:00+00') TO ('2025-11-05 00:00:00+00');

alter table public.transaction_decode_p2025_11_04
    owner to postgres;

create table public.transaction_decode_p2025_11_05
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-05 00:00:00+00') TO ('2025-11-06 00:00:00+00');

alter table public.transaction_decode_p2025_11_05
    owner to postgres;

create table public.transaction_decode_p2025_11_06
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-06 00:00:00+00') TO ('2025-11-07 00:00:00+00');

alter table public.transaction_decode_p2025_11_06
    owner to postgres;

create table public.transaction_decode_p2025_11_07
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-07 00:00:00+00') TO ('2025-11-08 00:00:00+00');

alter table public.transaction_decode_p2025_11_07
    owner to postgres;

create table public.transaction_decode_p2025_11_08
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-08 00:00:00+00') TO ('2025-11-09 00:00:00+00');

alter table public.transaction_decode_p2025_11_08
    owner to postgres;

create table public.transaction_decode_p2025_11_09
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-09 00:00:00+00') TO ('2025-11-10 00:00:00+00');

alter table public.transaction_decode_p2025_11_09
    owner to postgres;

create table public.transaction_decode_p2025_11_10
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-10 00:00:00+00') TO ('2025-11-11 00:00:00+00');

alter table public.transaction_decode_p2025_11_10
    owner to postgres;

create table public.transaction_decode_p2025_11_11
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-11 00:00:00+00') TO ('2025-11-12 00:00:00+00');

alter table public.transaction_decode_p2025_11_11
    owner to postgres;

create table public.transaction_decode_p2025_11_12
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-12 00:00:00+00') TO ('2025-11-13 00:00:00+00');

alter table public.transaction_decode_p2025_11_12
    owner to postgres;

create table public.transaction_decode_p2025_11_13
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-13 00:00:00+00') TO ('2025-11-14 00:00:00+00');

alter table public.transaction_decode_p2025_11_13
    owner to postgres;

create table public.transaction_decode_p2025_11_14
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-14 00:00:00+00') TO ('2025-11-15 00:00:00+00');

alter table public.transaction_decode_p2025_11_14
    owner to postgres;

create table public.transaction_decode_p2025_11_15
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-15 00:00:00+00') TO ('2025-11-16 00:00:00+00');

alter table public.transaction_decode_p2025_11_15
    owner to postgres;

create table public.transaction_decode_p2025_11_16
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-16 00:00:00+00') TO ('2025-11-17 00:00:00+00');

alter table public.transaction_decode_p2025_11_16
    owner to postgres;

create table public.transaction_decode_p2025_11_17
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-17 00:00:00+00') TO ('2025-11-18 00:00:00+00');

alter table public.transaction_decode_p2025_11_17
    owner to postgres;

create table public.transaction_decode_p2025_11_18
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-18 00:00:00+00') TO ('2025-11-19 00:00:00+00');

alter table public.transaction_decode_p2025_11_18
    owner to postgres;

create table public.transaction_decode_p2025_11_19
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-19 00:00:00+00') TO ('2025-11-20 00:00:00+00');

alter table public.transaction_decode_p2025_11_19
    owner to postgres;

create table public.transaction_decode_p2025_11_20
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-20 00:00:00+00') TO ('2025-11-21 00:00:00+00');

alter table public.transaction_decode_p2025_11_20
    owner to postgres;

create table public.transaction_decode_p2025_11_21
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-21 00:00:00+00') TO ('2025-11-22 00:00:00+00');

alter table public.transaction_decode_p2025_11_21
    owner to postgres;

create table public.transaction_decode_p2025_11_22
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-22 00:00:00+00') TO ('2025-11-23 00:00:00+00');

alter table public.transaction_decode_p2025_11_22
    owner to postgres;

create table public.transaction_decode_p2025_11_23
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-23 00:00:00+00') TO ('2025-11-24 00:00:00+00');

alter table public.transaction_decode_p2025_11_23
    owner to postgres;

create table public.transaction_decode_p2025_11_24
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-24 00:00:00+00') TO ('2025-11-25 00:00:00+00');

alter table public.transaction_decode_p2025_11_24
    owner to postgres;

create table public.transaction_decode_p2025_11_25
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-25 00:00:00+00') TO ('2025-11-26 00:00:00+00');

alter table public.transaction_decode_p2025_11_25
    owner to postgres;

create table public.transaction_decode_p2025_11_26
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-26 00:00:00+00') TO ('2025-11-27 00:00:00+00');

alter table public.transaction_decode_p2025_11_26
    owner to postgres;

create table public.transaction_decode_p2025_11_27
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-27 00:00:00+00') TO ('2025-11-28 00:00:00+00');

alter table public.transaction_decode_p2025_11_27
    owner to postgres;

create table public.transaction_decode_p2025_11_28
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-28 00:00:00+00') TO ('2025-11-29 00:00:00+00');

alter table public.transaction_decode_p2025_11_28
    owner to postgres;

create table public.transaction_decode_p2025_11_29
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-29 00:00:00+00') TO ('2025-11-30 00:00:00+00');

alter table public.transaction_decode_p2025_11_29
    owner to postgres;

create table public.transaction_decode_p2025_11_30
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-11-30 00:00:00+00') TO ('2025-12-01 00:00:00+00');

alter table public.transaction_decode_p2025_11_30
    owner to postgres;

create table public.transaction_decode_p2025_12_01
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-01 00:00:00+00') TO ('2025-12-02 00:00:00+00');

alter table public.transaction_decode_p2025_12_01
    owner to postgres;

create table public.transaction_decode_p2025_12_02
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-02 00:00:00+00') TO ('2025-12-03 00:00:00+00');

alter table public.transaction_decode_p2025_12_02
    owner to postgres;

create table public.transaction_decode_p2025_12_03
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-03 00:00:00+00') TO ('2025-12-04 00:00:00+00');

alter table public.transaction_decode_p2025_12_03
    owner to postgres;

create table public.transaction_decode_p2025_12_04
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-04 00:00:00+00') TO ('2025-12-05 00:00:00+00');

alter table public.transaction_decode_p2025_12_04
    owner to postgres;

create table public.transaction_decode_p2025_12_05
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-05 00:00:00+00') TO ('2025-12-06 00:00:00+00');

alter table public.transaction_decode_p2025_12_05
    owner to postgres;

create table public.transaction_decode_p2025_12_06
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-06 00:00:00+00') TO ('2025-12-07 00:00:00+00');

alter table public.transaction_decode_p2025_12_06
    owner to postgres;

create table public.transaction_decode_p2025_12_07
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-07 00:00:00+00') TO ('2025-12-08 00:00:00+00');

alter table public.transaction_decode_p2025_12_07
    owner to postgres;

create table public.transaction_decode_p2025_12_08
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-08 00:00:00+00') TO ('2025-12-09 00:00:00+00');

alter table public.transaction_decode_p2025_12_08
    owner to postgres;

create table public.transaction_decode_p2025_12_09
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-09 00:00:00+00') TO ('2025-12-10 00:00:00+00');

alter table public.transaction_decode_p2025_12_09
    owner to postgres;

create table public.transaction_decode_p2025_12_10
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-10 00:00:00+00') TO ('2025-12-11 00:00:00+00');

alter table public.transaction_decode_p2025_12_10
    owner to postgres;

create table public.transaction_decode_p2025_12_11
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-11 00:00:00+00') TO ('2025-12-12 00:00:00+00');

alter table public.transaction_decode_p2025_12_11
    owner to postgres;

create table public.transaction_decode_p2025_12_12
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-12 00:00:00+00') TO ('2025-12-13 00:00:00+00');

alter table public.transaction_decode_p2025_12_12
    owner to postgres;

create table public.transaction_decode_p2025_12_13
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-13 00:00:00+00') TO ('2025-12-14 00:00:00+00');

alter table public.transaction_decode_p2025_12_13
    owner to postgres;

create table public.transaction_decode_p2025_12_14
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-14 00:00:00+00') TO ('2025-12-15 00:00:00+00');

alter table public.transaction_decode_p2025_12_14
    owner to postgres;

create table public.transaction_decode_p2025_12_15
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-15 00:00:00+00') TO ('2025-12-16 00:00:00+00');

alter table public.transaction_decode_p2025_12_15
    owner to postgres;

create table public.transaction_decode_p2025_12_16
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-16 00:00:00+00') TO ('2025-12-17 00:00:00+00');

alter table public.transaction_decode_p2025_12_16
    owner to postgres;

create table public.transaction_decode_p2025_12_17
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-17 00:00:00+00') TO ('2025-12-18 00:00:00+00');

alter table public.transaction_decode_p2025_12_17
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_10
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-10 00:00:00+00') TO ('2025-08-11 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_10
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_11
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-11 00:00:00+00') TO ('2025-08-12 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_11
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_12
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-12 00:00:00+00') TO ('2025-08-13 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_12
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_13
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-13 00:00:00+00') TO ('2025-08-14 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_13
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_14
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-14 00:00:00+00') TO ('2025-08-15 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_14
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_15
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-15 00:00:00+00') TO ('2025-08-16 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_15
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_16
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-16 00:00:00+00') TO ('2025-08-17 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_16
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_17
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-17 00:00:00+00') TO ('2025-08-18 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_17
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_18
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-18 00:00:00+00') TO ('2025-08-19 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_18
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_19
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-19 00:00:00+00') TO ('2025-08-20 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_19
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_20
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-20 00:00:00+00') TO ('2025-08-21 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_20
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_21
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-21 00:00:00+00') TO ('2025-08-22 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_21
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_22
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-22 00:00:00+00') TO ('2025-08-23 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_22
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_23
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-23 00:00:00+00') TO ('2025-08-24 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_23
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_24
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-24 00:00:00+00') TO ('2025-08-25 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_24
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_25
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-25 00:00:00+00') TO ('2025-08-26 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_25
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_26
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-26 00:00:00+00') TO ('2025-08-27 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_26
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_27
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-27 00:00:00+00') TO ('2025-08-28 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_27
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_28
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-28 00:00:00+00') TO ('2025-08-29 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_28
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_29
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-29 00:00:00+00') TO ('2025-08-30 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_29
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_30
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-30 00:00:00+00') TO ('2025-08-31 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_30
    owner to postgres;

create table public.transaction_receipt_decode_p2025_08_31
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-08-31 00:00:00+00') TO ('2025-09-01 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_08_31
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_01
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-01 00:00:00+00') TO ('2025-09-02 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_01
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_02
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-02 00:00:00+00') TO ('2025-09-03 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_02
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_03
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-03 00:00:00+00') TO ('2025-09-04 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_03
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_04
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-04 00:00:00+00') TO ('2025-09-05 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_04
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_05
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-05 00:00:00+00') TO ('2025-09-06 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_05
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_06
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-06 00:00:00+00') TO ('2025-09-07 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_06
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_07
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-07 00:00:00+00') TO ('2025-09-08 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_07
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_08
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-08 00:00:00+00') TO ('2025-09-09 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_08
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_09
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-09 00:00:00+00') TO ('2025-09-10 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_09
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_10
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-10 00:00:00+00') TO ('2025-09-11 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_10
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_11
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-11 00:00:00+00') TO ('2025-09-12 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_11
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_12
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-12 00:00:00+00') TO ('2025-09-13 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_12
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_13
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-13 00:00:00+00') TO ('2025-09-14 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_13
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_14
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-14 00:00:00+00') TO ('2025-09-15 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_14
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_15
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-15 00:00:00+00') TO ('2025-09-16 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_15
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_16
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-16 00:00:00+00') TO ('2025-09-17 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_16
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_17
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-17 00:00:00+00') TO ('2025-09-18 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_17
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_18
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-18 00:00:00+00') TO ('2025-09-19 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_18
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_19
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-19 00:00:00+00') TO ('2025-09-20 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_19
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_20
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-20 00:00:00+00') TO ('2025-09-21 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_20
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_21
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-21 00:00:00+00') TO ('2025-09-22 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_21
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_22
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-22 00:00:00+00') TO ('2025-09-23 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_22
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_23
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-23 00:00:00+00') TO ('2025-09-24 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_23
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_24
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-24 00:00:00+00') TO ('2025-09-25 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_24
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_25
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-25 00:00:00+00') TO ('2025-09-26 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_25
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_26
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-26 00:00:00+00') TO ('2025-09-27 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_26
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_27
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-27 00:00:00+00') TO ('2025-09-28 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_27
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_28
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-28 00:00:00+00') TO ('2025-09-29 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_28
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_29
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-29 00:00:00+00') TO ('2025-09-30 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_29
    owner to postgres;

create table public.transaction_receipt_decode_p2025_09_30
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-09-30 00:00:00+00') TO ('2025-10-01 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_09_30
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_01
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-01 00:00:00+00') TO ('2025-10-02 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_01
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_02
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-02 00:00:00+00') TO ('2025-10-03 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_02
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_03
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-03 00:00:00+00') TO ('2025-10-04 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_03
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_04
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-04 00:00:00+00') TO ('2025-10-05 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_04
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_05
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-05 00:00:00+00') TO ('2025-10-06 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_05
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_06
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-06 00:00:00+00') TO ('2025-10-07 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_06
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_07
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-07 00:00:00+00') TO ('2025-10-08 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_07
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_08
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-08 00:00:00+00') TO ('2025-10-09 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_08
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_09
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-09 00:00:00+00') TO ('2025-10-10 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_09
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_10
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-10 00:00:00+00') TO ('2025-10-11 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_10
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_11
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-11 00:00:00+00') TO ('2025-10-12 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_11
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_12
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-12 00:00:00+00') TO ('2025-10-13 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_12
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_13
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-13 00:00:00+00') TO ('2025-10-14 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_13
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_14
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-14 00:00:00+00') TO ('2025-10-15 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_14
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_15
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-15 00:00:00+00') TO ('2025-10-16 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_15
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_16
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-16 00:00:00+00') TO ('2025-10-17 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_16
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_17
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-17 00:00:00+00') TO ('2025-10-18 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_17
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_18
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-18 00:00:00+00') TO ('2025-10-19 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_18
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_19
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-19 00:00:00+00') TO ('2025-10-20 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_19
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_20
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-20 00:00:00+00') TO ('2025-10-21 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_20
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_21
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-21 00:00:00+00') TO ('2025-10-22 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_21
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_22
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-22 00:00:00+00') TO ('2025-10-23 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_22
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_23
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-23 00:00:00+00') TO ('2025-10-24 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_23
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_24
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-24 00:00:00+00') TO ('2025-10-25 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_24
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_25
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-25 00:00:00+00') TO ('2025-10-26 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_25
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_26
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-26 00:00:00+00') TO ('2025-10-27 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_26
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_27
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-27 00:00:00+00') TO ('2025-10-28 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_27
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_28
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-28 00:00:00+00') TO ('2025-10-29 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_28
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_29
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-29 00:00:00+00') TO ('2025-10-30 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_29
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_30
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-30 00:00:00+00') TO ('2025-10-31 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_30
    owner to postgres;

create table public.transaction_receipt_decode_p2025_10_31
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-10-31 00:00:00+00') TO ('2025-11-01 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_10_31
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_01
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-01 00:00:00+00') TO ('2025-11-02 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_01
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_02
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-02 00:00:00+00') TO ('2025-11-03 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_02
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_03
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-03 00:00:00+00') TO ('2025-11-04 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_03
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_04
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-04 00:00:00+00') TO ('2025-11-05 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_04
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_05
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-05 00:00:00+00') TO ('2025-11-06 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_05
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_06
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-06 00:00:00+00') TO ('2025-11-07 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_06
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_07
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-07 00:00:00+00') TO ('2025-11-08 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_07
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_08
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-08 00:00:00+00') TO ('2025-11-09 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_08
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_09
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-09 00:00:00+00') TO ('2025-11-10 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_09
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_10
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-10 00:00:00+00') TO ('2025-11-11 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_10
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_11
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-11 00:00:00+00') TO ('2025-11-12 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_11
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_12
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-12 00:00:00+00') TO ('2025-11-13 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_12
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_13
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-13 00:00:00+00') TO ('2025-11-14 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_13
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_14
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-14 00:00:00+00') TO ('2025-11-15 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_14
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_15
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-15 00:00:00+00') TO ('2025-11-16 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_15
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_16
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-16 00:00:00+00') TO ('2025-11-17 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_16
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_17
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-17 00:00:00+00') TO ('2025-11-18 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_17
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_18
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-18 00:00:00+00') TO ('2025-11-19 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_18
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_19
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-19 00:00:00+00') TO ('2025-11-20 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_19
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_20
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-20 00:00:00+00') TO ('2025-11-21 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_20
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_21
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-21 00:00:00+00') TO ('2025-11-22 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_21
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_22
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-22 00:00:00+00') TO ('2025-11-23 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_22
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_23
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-23 00:00:00+00') TO ('2025-11-24 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_23
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_24
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-24 00:00:00+00') TO ('2025-11-25 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_24
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_25
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-25 00:00:00+00') TO ('2025-11-26 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_25
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_26
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-26 00:00:00+00') TO ('2025-11-27 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_26
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_27
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-27 00:00:00+00') TO ('2025-11-28 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_27
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_28
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-28 00:00:00+00') TO ('2025-11-29 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_28
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_29
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-29 00:00:00+00') TO ('2025-11-30 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_29
    owner to postgres;

create table public.transaction_receipt_decode_p2025_11_30
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-11-30 00:00:00+00') TO ('2025-12-01 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_11_30
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_01
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-01 00:00:00+00') TO ('2025-12-02 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_01
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_02
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-02 00:00:00+00') TO ('2025-12-03 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_02
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_03
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-03 00:00:00+00') TO ('2025-12-04 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_03
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_04
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-04 00:00:00+00') TO ('2025-12-05 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_04
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_05
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-05 00:00:00+00') TO ('2025-12-06 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_05
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_06
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-06 00:00:00+00') TO ('2025-12-07 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_06
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_07
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-07 00:00:00+00') TO ('2025-12-08 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_07
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_08
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-08 00:00:00+00') TO ('2025-12-09 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_08
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_09
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-09 00:00:00+00') TO ('2025-12-10 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_09
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_10
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-10 00:00:00+00') TO ('2025-12-11 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_10
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_11
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-11 00:00:00+00') TO ('2025-12-12 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_11
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_12
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-12 00:00:00+00') TO ('2025-12-13 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_12
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_13
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-13 00:00:00+00') TO ('2025-12-14 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_13
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_14
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-14 00:00:00+00') TO ('2025-12-15 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_14
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_15
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-15 00:00:00+00') TO ('2025-12-16 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_15
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_16
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-16 00:00:00+00') TO ('2025-12-17 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_16
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_17
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-17 00:00:00+00') TO ('2025-12-18 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_17
    owner to postgres;

create table public.aa_transaction_info_p2025w35
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-08-25 00:00:00+00') TO ('2025-09-01 00:00:00+00');

alter table public.aa_transaction_info_p2025w35
    owner to postgres;

create table public.aa_transaction_info_p2025w36
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-09-01 00:00:00+00') TO ('2025-09-08 00:00:00+00');

alter table public.aa_transaction_info_p2025w36
    owner to postgres;

create table public.aa_transaction_info_p2025w37
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-09-08 00:00:00+00') TO ('2025-09-15 00:00:00+00');

alter table public.aa_transaction_info_p2025w37
    owner to postgres;

create table public.aa_transaction_info_p2025w38
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-09-15 00:00:00+00') TO ('2025-09-22 00:00:00+00');

alter table public.aa_transaction_info_p2025w38
    owner to postgres;

create table public.aa_transaction_info_p2025w39
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-09-22 00:00:00+00') TO ('2025-09-29 00:00:00+00');

alter table public.aa_transaction_info_p2025w39
    owner to postgres;

create table public.aa_transaction_info_p2025w40
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-09-29 00:00:00+00') TO ('2025-10-06 00:00:00+00');

alter table public.aa_transaction_info_p2025w40
    owner to postgres;

create table public.aa_transaction_info_p2025w41
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-10-06 00:00:00+00') TO ('2025-10-13 00:00:00+00');

alter table public.aa_transaction_info_p2025w41
    owner to postgres;

create table public.aa_transaction_info_p2025w42
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-10-13 00:00:00+00') TO ('2025-10-20 00:00:00+00');

alter table public.aa_transaction_info_p2025w42
    owner to postgres;

create table public.aa_transaction_info_p2025w43
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-10-20 00:00:00+00') TO ('2025-10-27 00:00:00+00');

alter table public.aa_transaction_info_p2025w43
    owner to postgres;

create table public.aa_transaction_info_p2025w44
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-10-27 00:00:00+00') TO ('2025-11-03 00:00:00+00');

alter table public.aa_transaction_info_p2025w44
    owner to postgres;

create table public.aa_transaction_info_p2025w45
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-11-03 00:00:00+00') TO ('2025-11-10 00:00:00+00');

alter table public.aa_transaction_info_p2025w45
    owner to postgres;

create table public.aa_transaction_info_p2025w46
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-11-10 00:00:00+00') TO ('2025-11-17 00:00:00+00');

alter table public.aa_transaction_info_p2025w46
    owner to postgres;

create table public.aa_transaction_info_p2025w47
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-11-17 00:00:00+00') TO ('2025-11-24 00:00:00+00');

alter table public.aa_transaction_info_p2025w47
    owner to postgres;

create table public.aa_transaction_info_p2025w48
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-11-24 00:00:00+00') TO ('2025-12-01 00:00:00+00');

alter table public.aa_transaction_info_p2025w48
    owner to postgres;

create table public.aa_transaction_info_p2025w49
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-12-01 00:00:00+00') TO ('2025-12-08 00:00:00+00');

alter table public.aa_transaction_info_p2025w49
    owner to postgres;

create table public.aa_transaction_info_p2025w50
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-12-08 00:00:00+00') TO ('2025-12-15 00:00:00+00');

alter table public.aa_transaction_info_p2025w50
    owner to postgres;

create table public.aa_transaction_info_p2025w51
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-12-15 00:00:00+00') TO ('2025-12-22 00:00:00+00');

alter table public.aa_transaction_info_p2025w51
    owner to postgres;

create table public.aa_transaction_info_p2025w52
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-12-22 00:00:00+00') TO ('2025-12-29 00:00:00+00');

alter table public.aa_transaction_info_p2025w52
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w35
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-08-25 00:00:00+00') TO ('2025-09-01 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w35
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w36
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-09-01 00:00:00+00') TO ('2025-09-08 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w36
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w37
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-09-08 00:00:00+00') TO ('2025-09-15 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w37
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w38
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-09-15 00:00:00+00') TO ('2025-09-22 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w38
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w39
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-09-22 00:00:00+00') TO ('2025-09-29 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w39
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w40
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-09-29 00:00:00+00') TO ('2025-10-06 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w40
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w41
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-10-06 00:00:00+00') TO ('2025-10-13 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w41
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w42
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-10-13 00:00:00+00') TO ('2025-10-20 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w42
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w43
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-10-20 00:00:00+00') TO ('2025-10-27 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w43
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w44
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-10-27 00:00:00+00') TO ('2025-11-03 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w44
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w45
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-11-03 00:00:00+00') TO ('2025-11-10 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w45
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w46
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-11-10 00:00:00+00') TO ('2025-11-17 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w46
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w47
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-11-17 00:00:00+00') TO ('2025-11-24 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w47
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w48
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-11-24 00:00:00+00') TO ('2025-12-01 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w48
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w49
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-12-01 00:00:00+00') TO ('2025-12-08 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w49
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w50
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-12-08 00:00:00+00') TO ('2025-12-15 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w50
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w51
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-12-15 00:00:00+00') TO ('2025-12-22 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w51
    owner to postgres;

create table public.aa_user_ops_calldata_p2025w52
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-12-22 00:00:00+00') TO ('2025-12-29 00:00:00+00');

alter table public.aa_user_ops_calldata_p2025w52
    owner to postgres;

create table public.aa_user_ops_info_p2025w35
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-08-25 00:00:00+00') TO ('2025-09-01 00:00:00+00');

alter table public.aa_user_ops_info_p2025w35
    owner to postgres;

create table public.aa_user_ops_info_p2025w36
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-09-01 00:00:00+00') TO ('2025-09-08 00:00:00+00');

alter table public.aa_user_ops_info_p2025w36
    owner to postgres;

create table public.aa_user_ops_info_p2025w37
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-09-08 00:00:00+00') TO ('2025-09-15 00:00:00+00');

alter table public.aa_user_ops_info_p2025w37
    owner to postgres;

create table public.aa_user_ops_info_p2025w38
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-09-15 00:00:00+00') TO ('2025-09-22 00:00:00+00');

alter table public.aa_user_ops_info_p2025w38
    owner to postgres;

create table public.aa_user_ops_info_p2025w39
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-09-22 00:00:00+00') TO ('2025-09-29 00:00:00+00');

alter table public.aa_user_ops_info_p2025w39
    owner to postgres;

create table public.aa_user_ops_info_p2025w40
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-09-29 00:00:00+00') TO ('2025-10-06 00:00:00+00');

alter table public.aa_user_ops_info_p2025w40
    owner to postgres;

create table public.aa_user_ops_info_p2025w41
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-10-06 00:00:00+00') TO ('2025-10-13 00:00:00+00');

alter table public.aa_user_ops_info_p2025w41
    owner to postgres;

create table public.aa_user_ops_info_p2025w42
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-10-13 00:00:00+00') TO ('2025-10-20 00:00:00+00');

alter table public.aa_user_ops_info_p2025w42
    owner to postgres;

create table public.aa_user_ops_info_p2025w43
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-10-20 00:00:00+00') TO ('2025-10-27 00:00:00+00');

alter table public.aa_user_ops_info_p2025w43
    owner to postgres;

create table public.aa_user_ops_info_p2025w44
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-10-27 00:00:00+00') TO ('2025-11-03 00:00:00+00');

alter table public.aa_user_ops_info_p2025w44
    owner to postgres;

create table public.aa_user_ops_info_p2025w45
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-11-03 00:00:00+00') TO ('2025-11-10 00:00:00+00');

alter table public.aa_user_ops_info_p2025w45
    owner to postgres;

create table public.aa_user_ops_info_p2025w46
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-11-10 00:00:00+00') TO ('2025-11-17 00:00:00+00');

alter table public.aa_user_ops_info_p2025w46
    owner to postgres;

create table public.aa_user_ops_info_p2025w47
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-11-17 00:00:00+00') TO ('2025-11-24 00:00:00+00');

alter table public.aa_user_ops_info_p2025w47
    owner to postgres;

create table public.aa_user_ops_info_p2025w48
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-11-24 00:00:00+00') TO ('2025-12-01 00:00:00+00');

alter table public.aa_user_ops_info_p2025w48
    owner to postgres;

create table public.aa_user_ops_info_p2025w49
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-12-01 00:00:00+00') TO ('2025-12-08 00:00:00+00');

alter table public.aa_user_ops_info_p2025w49
    owner to postgres;

create table public.aa_user_ops_info_p2025w50
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-12-08 00:00:00+00') TO ('2025-12-15 00:00:00+00');

alter table public.aa_user_ops_info_p2025w50
    owner to postgres;

create table public.aa_user_ops_info_p2025w51
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-12-15 00:00:00+00') TO ('2025-12-22 00:00:00+00');

alter table public.aa_user_ops_info_p2025w51
    owner to postgres;

create table public.aa_user_ops_info_p2025w52
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-12-22 00:00:00+00') TO ('2025-12-29 00:00:00+00');

alter table public.aa_user_ops_info_p2025w52
    owner to postgres;

create table public.transaction_decode_p2025_12_18
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-18 00:00:00+00') TO ('2025-12-19 00:00:00+00');

alter table public.transaction_decode_p2025_12_18
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_18
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-18 00:00:00+00') TO ('2025-12-19 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_18
    owner to postgres;

create table public.transaction_decode_p2025_12_19
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-19 00:00:00+00') TO ('2025-12-20 00:00:00+00');

alter table public.transaction_decode_p2025_12_19
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_19
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-19 00:00:00+00') TO ('2025-12-20 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_19
    owner to postgres;

create table public.transaction_decode_p2025_12_20
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-20 00:00:00+00') TO ('2025-12-21 00:00:00+00');

alter table public.transaction_decode_p2025_12_20
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_20
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-20 00:00:00+00') TO ('2025-12-21 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_20
    owner to postgres;

create table public.block_data_decode_p2026w01
    partition of public.block_data_decode
    FOR VALUES FROM ('2025-12-29 00:00:00+00') TO ('2026-01-05 00:00:00+00');

alter table public.block_data_decode_p2026w01
    owner to postgres;

create table public.transaction_decode_p2025_12_21
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-21 00:00:00+00') TO ('2025-12-22 00:00:00+00');

alter table public.transaction_decode_p2025_12_21
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_21
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-21 00:00:00+00') TO ('2025-12-22 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_21
    owner to postgres;

create table public.aa_transaction_info_p2026w01
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2025-12-29 00:00:00+00') TO ('2026-01-05 00:00:00+00');

alter table public.aa_transaction_info_p2026w01
    owner to postgres;

create table public.aa_user_ops_calldata_p2026w01
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2025-12-29 00:00:00+00') TO ('2026-01-05 00:00:00+00');

alter table public.aa_user_ops_calldata_p2026w01
    owner to postgres;

create table public.aa_user_ops_info_p2026w01
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2025-12-29 00:00:00+00') TO ('2026-01-05 00:00:00+00');

alter table public.aa_user_ops_info_p2026w01
    owner to postgres;

create table public.transaction_decode_p2025_12_22
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-22 00:00:00+00') TO ('2025-12-23 00:00:00+00');

alter table public.transaction_decode_p2025_12_22
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_22
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-22 00:00:00+00') TO ('2025-12-23 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_22
    owner to postgres;

create table public.transaction_decode_p2025_12_23
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-23 00:00:00+00') TO ('2025-12-24 00:00:00+00');

alter table public.transaction_decode_p2025_12_23
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_23
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-23 00:00:00+00') TO ('2025-12-24 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_23
    owner to postgres;

create table public.transaction_decode_p2025_12_24
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-24 00:00:00+00') TO ('2025-12-25 00:00:00+00');

alter table public.transaction_decode_p2025_12_24
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_24
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-24 00:00:00+00') TO ('2025-12-25 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_24
    owner to postgres;

create table public.transaction_decode_p2025_12_25
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-25 00:00:00+00') TO ('2025-12-26 00:00:00+00');

alter table public.transaction_decode_p2025_12_25
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_25
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-25 00:00:00+00') TO ('2025-12-26 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_25
    owner to postgres;

create table public.transaction_decode_p2025_12_26
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-26 00:00:00+00') TO ('2025-12-27 00:00:00+00');

alter table public.transaction_decode_p2025_12_26
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_26
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-26 00:00:00+00') TO ('2025-12-27 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_26
    owner to postgres;

create table public.transaction_decode_p2025_12_27
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-27 00:00:00+00') TO ('2025-12-28 00:00:00+00');

alter table public.transaction_decode_p2025_12_27
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_27
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-27 00:00:00+00') TO ('2025-12-28 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_27
    owner to postgres;

create table public.block_data_decode_p2026w02
    partition of public.block_data_decode
    FOR VALUES FROM ('2026-01-05 00:00:00+00') TO ('2026-01-12 00:00:00+00');

alter table public.block_data_decode_p2026w02
    owner to postgres;

create table public.transaction_decode_p2025_12_28
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-28 00:00:00+00') TO ('2025-12-29 00:00:00+00');

alter table public.transaction_decode_p2025_12_28
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_28
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-28 00:00:00+00') TO ('2025-12-29 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_28
    owner to postgres;

create table public.aa_transaction_info_p2026w02
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2026-01-05 00:00:00+00') TO ('2026-01-12 00:00:00+00');

alter table public.aa_transaction_info_p2026w02
    owner to postgres;

create table public.aa_user_ops_calldata_p2026w02
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2026-01-05 00:00:00+00') TO ('2026-01-12 00:00:00+00');

alter table public.aa_user_ops_calldata_p2026w02
    owner to postgres;

create table public.aa_user_ops_info_p2026w02
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2026-01-05 00:00:00+00') TO ('2026-01-12 00:00:00+00');

alter table public.aa_user_ops_info_p2026w02
    owner to postgres;

create table public.transaction_decode_p2025_12_29
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-29 00:00:00+00') TO ('2025-12-30 00:00:00+00');

alter table public.transaction_decode_p2025_12_29
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_29
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-29 00:00:00+00') TO ('2025-12-30 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_29
    owner to postgres;

create table public.transaction_decode_p2025_12_30
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-30 00:00:00+00') TO ('2025-12-31 00:00:00+00');

alter table public.transaction_decode_p2025_12_30
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_30
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-30 00:00:00+00') TO ('2025-12-31 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_30
    owner to postgres;

create table public.transaction_decode_p2025_12_31
    partition of public.transaction_decode
    FOR VALUES FROM ('2025-12-31 00:00:00+00') TO ('2026-01-01 00:00:00+00');

alter table public.transaction_decode_p2025_12_31
    owner to postgres;

create table public.transaction_receipt_decode_p2025_12_31
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2025-12-31 00:00:00+00') TO ('2026-01-01 00:00:00+00');

alter table public.transaction_receipt_decode_p2025_12_31
    owner to postgres;

create table public.transaction_decode_p2026_01_01
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-01 00:00:00+00') TO ('2026-01-02 00:00:00+00');

alter table public.transaction_decode_p2026_01_01
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_01
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-01 00:00:00+00') TO ('2026-01-02 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_01
    owner to postgres;

create table public.transaction_decode_p2026_01_02
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-02 00:00:00+00') TO ('2026-01-03 00:00:00+00');

alter table public.transaction_decode_p2026_01_02
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_02
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-02 00:00:00+00') TO ('2026-01-03 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_02
    owner to postgres;

create table public.block_data_decode_p2026w03
    partition of public.block_data_decode
    FOR VALUES FROM ('2026-01-12 00:00:00+00') TO ('2026-01-19 00:00:00+00');

alter table public.block_data_decode_p2026w03
    owner to postgres;

create table public.transaction_decode_p2026_01_03
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-03 00:00:00+00') TO ('2026-01-04 00:00:00+00');

alter table public.transaction_decode_p2026_01_03
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_03
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-03 00:00:00+00') TO ('2026-01-04 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_03
    owner to postgres;

create table public.aa_transaction_info_p2026w03
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2026-01-12 00:00:00+00') TO ('2026-01-19 00:00:00+00');

alter table public.aa_transaction_info_p2026w03
    owner to postgres;

create table public.aa_user_ops_calldata_p2026w03
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2026-01-12 00:00:00+00') TO ('2026-01-19 00:00:00+00');

alter table public.aa_user_ops_calldata_p2026w03
    owner to postgres;

create table public.aa_user_ops_info_p2026w03
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2026-01-12 00:00:00+00') TO ('2026-01-19 00:00:00+00');

alter table public.aa_user_ops_info_p2026w03
    owner to postgres;

create table public.transaction_decode_p2026_01_04
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-04 00:00:00+00') TO ('2026-01-05 00:00:00+00');

alter table public.transaction_decode_p2026_01_04
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_04
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-04 00:00:00+00') TO ('2026-01-05 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_04
    owner to postgres;

create table public.transaction_decode_p2026_01_05
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-05 00:00:00+00') TO ('2026-01-06 00:00:00+00');

alter table public.transaction_decode_p2026_01_05
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_05
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-05 00:00:00+00') TO ('2026-01-06 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_05
    owner to postgres;

create table public.transaction_decode_p2026_01_06
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-06 00:00:00+00') TO ('2026-01-07 00:00:00+00');

alter table public.transaction_decode_p2026_01_06
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_06
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-06 00:00:00+00') TO ('2026-01-07 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_06
    owner to postgres;

create table public.transaction_decode_p2026_01_07
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-07 00:00:00+00') TO ('2026-01-08 00:00:00+00');

alter table public.transaction_decode_p2026_01_07
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_07
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-07 00:00:00+00') TO ('2026-01-08 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_07
    owner to postgres;

create table public.transaction_decode_p2026_01_08
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-08 00:00:00+00') TO ('2026-01-09 00:00:00+00');

alter table public.transaction_decode_p2026_01_08
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_08
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-08 00:00:00+00') TO ('2026-01-09 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_08
    owner to postgres;

create table public.transaction_decode_p2026_01_09
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-09 00:00:00+00') TO ('2026-01-10 00:00:00+00');

alter table public.transaction_decode_p2026_01_09
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_09
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-09 00:00:00+00') TO ('2026-01-10 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_09
    owner to postgres;

create table public.block_data_decode_p2026w04
    partition of public.block_data_decode
    FOR VALUES FROM ('2026-01-19 00:00:00+00') TO ('2026-01-26 00:00:00+00');

alter table public.block_data_decode_p2026w04
    owner to postgres;

create table public.transaction_decode_p2026_01_10
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-10 00:00:00+00') TO ('2026-01-11 00:00:00+00');

alter table public.transaction_decode_p2026_01_10
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_10
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-10 00:00:00+00') TO ('2026-01-11 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_10
    owner to postgres;

create table public.aa_transaction_info_p2026w04
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2026-01-19 00:00:00+00') TO ('2026-01-26 00:00:00+00');

alter table public.aa_transaction_info_p2026w04
    owner to postgres;

create table public.aa_user_ops_calldata_p2026w04
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2026-01-19 00:00:00+00') TO ('2026-01-26 00:00:00+00');

alter table public.aa_user_ops_calldata_p2026w04
    owner to postgres;

create table public.aa_user_ops_info_p2026w04
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2026-01-19 00:00:00+00') TO ('2026-01-26 00:00:00+00');

alter table public.aa_user_ops_info_p2026w04
    owner to postgres;

create table public.transaction_decode_p2026_01_11
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-11 00:00:00+00') TO ('2026-01-12 00:00:00+00');

alter table public.transaction_decode_p2026_01_11
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_11
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-11 00:00:00+00') TO ('2026-01-12 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_11
    owner to postgres;

create table public.transaction_decode_p2026_01_12
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-12 00:00:00+00') TO ('2026-01-13 00:00:00+00');

alter table public.transaction_decode_p2026_01_12
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_12
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-12 00:00:00+00') TO ('2026-01-13 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_12
    owner to postgres;

create table public.transaction_decode_p2026_01_13
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-13 00:00:00+00') TO ('2026-01-14 00:00:00+00');

alter table public.transaction_decode_p2026_01_13
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_13
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-13 00:00:00+00') TO ('2026-01-14 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_13
    owner to postgres;

create table public.transaction_decode_p2026_01_14
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-14 00:00:00+00') TO ('2026-01-15 00:00:00+00');

alter table public.transaction_decode_p2026_01_14
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_14
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-14 00:00:00+00') TO ('2026-01-15 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_14
    owner to postgres;

create table public.transaction_decode_p2026_01_15
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-15 00:00:00+00') TO ('2026-01-16 00:00:00+00');

alter table public.transaction_decode_p2026_01_15
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_15
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-15 00:00:00+00') TO ('2026-01-16 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_15
    owner to postgres;

create table public.transaction_decode_p2026_01_16
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-16 00:00:00+00') TO ('2026-01-17 00:00:00+00');

alter table public.transaction_decode_p2026_01_16
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_16
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-16 00:00:00+00') TO ('2026-01-17 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_16
    owner to postgres;

create table public.aa_block_info_p2026_06
    partition of public.aa_block_info
    FOR VALUES FROM ('2026-06-01 00:00:00+00') TO ('2026-07-01 00:00:00+00');

alter table public.aa_block_info_p2026_06
    owner to postgres;

create table public.block_data_decode_p2026w05
    partition of public.block_data_decode
    FOR VALUES FROM ('2026-01-26 00:00:00+00') TO ('2026-02-02 00:00:00+00');

alter table public.block_data_decode_p2026w05
    owner to postgres;

create table public.transaction_decode_p2026_01_17
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-17 00:00:00+00') TO ('2026-01-18 00:00:00+00');

alter table public.transaction_decode_p2026_01_17
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_17
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-17 00:00:00+00') TO ('2026-01-18 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_17
    owner to postgres;

create table public.aa_transaction_info_p2026w05
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2026-01-26 00:00:00+00') TO ('2026-02-02 00:00:00+00');

alter table public.aa_transaction_info_p2026w05
    owner to postgres;

create table public.aa_user_ops_calldata_p2026w05
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2026-01-26 00:00:00+00') TO ('2026-02-02 00:00:00+00');

alter table public.aa_user_ops_calldata_p2026w05
    owner to postgres;

create table public.aa_user_ops_info_p2026w05
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2026-01-26 00:00:00+00') TO ('2026-02-02 00:00:00+00');

alter table public.aa_user_ops_info_p2026w05
    owner to postgres;

create table public.transaction_decode_p2026_01_18
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-18 00:00:00+00') TO ('2026-01-19 00:00:00+00');

alter table public.transaction_decode_p2026_01_18
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_18
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-18 00:00:00+00') TO ('2026-01-19 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_18
    owner to postgres;

create table public.transaction_decode_p2026_01_19
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-19 00:00:00+00') TO ('2026-01-20 00:00:00+00');

alter table public.transaction_decode_p2026_01_19
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_19
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-19 00:00:00+00') TO ('2026-01-20 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_19
    owner to postgres;

create table public.transaction_decode_p2026_01_20
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-20 00:00:00+00') TO ('2026-01-21 00:00:00+00');

alter table public.transaction_decode_p2026_01_20
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_20
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-20 00:00:00+00') TO ('2026-01-21 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_20
    owner to postgres;

create table public.transaction_decode_p2026_01_21
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-21 00:00:00+00') TO ('2026-01-22 00:00:00+00');

alter table public.transaction_decode_p2026_01_21
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_21
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-21 00:00:00+00') TO ('2026-01-22 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_21
    owner to postgres;

create table public.transaction_decode_p2026_01_22
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-22 00:00:00+00') TO ('2026-01-23 00:00:00+00');

alter table public.transaction_decode_p2026_01_22
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_22
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-22 00:00:00+00') TO ('2026-01-23 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_22
    owner to postgres;

create table public.transaction_decode_p2026_01_23
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-23 00:00:00+00') TO ('2026-01-24 00:00:00+00');

alter table public.transaction_decode_p2026_01_23
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_23
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-23 00:00:00+00') TO ('2026-01-24 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_23
    owner to postgres;

create table public.block_data_decode_p2026w06
    partition of public.block_data_decode
    FOR VALUES FROM ('2026-02-02 00:00:00+00') TO ('2026-02-09 00:00:00+00');

alter table public.block_data_decode_p2026w06
    owner to postgres;

create table public.transaction_decode_p2026_01_24
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-24 00:00:00+00') TO ('2026-01-25 00:00:00+00');

alter table public.transaction_decode_p2026_01_24
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_24
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-24 00:00:00+00') TO ('2026-01-25 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_24
    owner to postgres;

create table public.aa_transaction_info_p2026w06
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2026-02-02 00:00:00+00') TO ('2026-02-09 00:00:00+00');

alter table public.aa_transaction_info_p2026w06
    owner to postgres;

create table public.aa_user_ops_calldata_p2026w06
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2026-02-02 00:00:00+00') TO ('2026-02-09 00:00:00+00');

alter table public.aa_user_ops_calldata_p2026w06
    owner to postgres;

create table public.aa_user_ops_info_p2026w06
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2026-02-02 00:00:00+00') TO ('2026-02-09 00:00:00+00');

alter table public.aa_user_ops_info_p2026w06
    owner to postgres;

create table public.transaction_decode_p2026_01_25
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-25 00:00:00+00') TO ('2026-01-26 00:00:00+00');

alter table public.transaction_decode_p2026_01_25
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_25
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-25 00:00:00+00') TO ('2026-01-26 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_25
    owner to postgres;

create table public.transaction_decode_p2026_01_26
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-26 00:00:00+00') TO ('2026-01-27 00:00:00+00');

alter table public.transaction_decode_p2026_01_26
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_26
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-26 00:00:00+00') TO ('2026-01-27 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_26
    owner to postgres;

create table public.transaction_decode_p2026_01_27
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-27 00:00:00+00') TO ('2026-01-28 00:00:00+00');

alter table public.transaction_decode_p2026_01_27
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_27
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-27 00:00:00+00') TO ('2026-01-28 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_27
    owner to postgres;

create table public.transaction_decode_p2026_01_28
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-28 00:00:00+00') TO ('2026-01-29 00:00:00+00');

alter table public.transaction_decode_p2026_01_28
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_28
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-28 00:00:00+00') TO ('2026-01-29 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_28
    owner to postgres;

create table public.transaction_decode_p2026_01_29
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-29 00:00:00+00') TO ('2026-01-30 00:00:00+00');

alter table public.transaction_decode_p2026_01_29
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_29
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-29 00:00:00+00') TO ('2026-01-30 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_29
    owner to postgres;

create table public.transaction_decode_p2026_01_30
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-30 00:00:00+00') TO ('2026-01-31 00:00:00+00');

alter table public.transaction_decode_p2026_01_30
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_30
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-30 00:00:00+00') TO ('2026-01-31 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_30
    owner to postgres;

create table public.block_data_decode_p2026w07
    partition of public.block_data_decode
    FOR VALUES FROM ('2026-02-09 00:00:00+00') TO ('2026-02-16 00:00:00+00');

alter table public.block_data_decode_p2026w07
    owner to postgres;

create table public.transaction_decode_p2026_01_31
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-01-31 00:00:00+00') TO ('2026-02-01 00:00:00+00');

alter table public.transaction_decode_p2026_01_31
    owner to postgres;

create table public.transaction_receipt_decode_p2026_01_31
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-01-31 00:00:00+00') TO ('2026-02-01 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_01_31
    owner to postgres;

create table public.aa_transaction_info_p2026w07
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2026-02-09 00:00:00+00') TO ('2026-02-16 00:00:00+00');

alter table public.aa_transaction_info_p2026w07
    owner to postgres;

create table public.aa_user_ops_calldata_p2026w07
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2026-02-09 00:00:00+00') TO ('2026-02-16 00:00:00+00');

alter table public.aa_user_ops_calldata_p2026w07
    owner to postgres;

create table public.aa_user_ops_info_p2026w07
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2026-02-09 00:00:00+00') TO ('2026-02-16 00:00:00+00');

alter table public.aa_user_ops_info_p2026w07
    owner to postgres;

create table public.transaction_decode_p2026_02_01
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-01 00:00:00+00') TO ('2026-02-02 00:00:00+00');

alter table public.transaction_decode_p2026_02_01
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_01
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-01 00:00:00+00') TO ('2026-02-02 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_01
    owner to postgres;

create table public.transaction_decode_p2026_02_02
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-02 00:00:00+00') TO ('2026-02-03 00:00:00+00');

alter table public.transaction_decode_p2026_02_02
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_02
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-02 00:00:00+00') TO ('2026-02-03 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_02
    owner to postgres;

create table public.transaction_decode_p2026_02_03
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-03 00:00:00+00') TO ('2026-02-04 00:00:00+00');

alter table public.transaction_decode_p2026_02_03
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_03
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-03 00:00:00+00') TO ('2026-02-04 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_03
    owner to postgres;

create table public.transaction_decode_p2026_02_04
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-04 00:00:00+00') TO ('2026-02-05 00:00:00+00');

alter table public.transaction_decode_p2026_02_04
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_04
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-04 00:00:00+00') TO ('2026-02-05 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_04
    owner to postgres;

create table public.transaction_decode_p2026_02_05
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-05 00:00:00+00') TO ('2026-02-06 00:00:00+00');

alter table public.transaction_decode_p2026_02_05
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_05
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-05 00:00:00+00') TO ('2026-02-06 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_05
    owner to postgres;

create table public.transaction_decode_p2026_02_06
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-06 00:00:00+00') TO ('2026-02-07 00:00:00+00');

alter table public.transaction_decode_p2026_02_06
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_06
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-06 00:00:00+00') TO ('2026-02-07 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_06
    owner to postgres;

create table public.block_data_decode_p2026w08
    partition of public.block_data_decode
    FOR VALUES FROM ('2026-02-16 00:00:00+00') TO ('2026-02-23 00:00:00+00');

alter table public.block_data_decode_p2026w08
    owner to postgres;

create table public.transaction_decode_p2026_02_07
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-07 00:00:00+00') TO ('2026-02-08 00:00:00+00');

alter table public.transaction_decode_p2026_02_07
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_07
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-07 00:00:00+00') TO ('2026-02-08 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_07
    owner to postgres;

create table public.aa_transaction_info_p2026w08
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2026-02-16 00:00:00+00') TO ('2026-02-23 00:00:00+00');

alter table public.aa_transaction_info_p2026w08
    owner to postgres;

create table public.aa_user_ops_calldata_p2026w08
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2026-02-16 00:00:00+00') TO ('2026-02-23 00:00:00+00');

alter table public.aa_user_ops_calldata_p2026w08
    owner to postgres;

create table public.aa_user_ops_info_p2026w08
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2026-02-16 00:00:00+00') TO ('2026-02-23 00:00:00+00');

alter table public.aa_user_ops_info_p2026w08
    owner to postgres;

create table public.transaction_decode_p2026_02_08
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-08 00:00:00+00') TO ('2026-02-09 00:00:00+00');

alter table public.transaction_decode_p2026_02_08
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_08
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-08 00:00:00+00') TO ('2026-02-09 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_08
    owner to postgres;

create table public.transaction_decode_p2026_02_09
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-09 00:00:00+00') TO ('2026-02-10 00:00:00+00');

alter table public.transaction_decode_p2026_02_09
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_09
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-09 00:00:00+00') TO ('2026-02-10 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_09
    owner to postgres;

create table public.transaction_decode_p2026_02_10
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-10 00:00:00+00') TO ('2026-02-11 00:00:00+00');

alter table public.transaction_decode_p2026_02_10
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_10
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-10 00:00:00+00') TO ('2026-02-11 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_10
    owner to postgres;

create table public.transaction_decode_p2026_02_11
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-11 00:00:00+00') TO ('2026-02-12 00:00:00+00');

alter table public.transaction_decode_p2026_02_11
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_11
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-11 00:00:00+00') TO ('2026-02-12 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_11
    owner to postgres;

create table public.transaction_decode_p2026_02_12
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-12 00:00:00+00') TO ('2026-02-13 00:00:00+00');

alter table public.transaction_decode_p2026_02_12
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_12
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-12 00:00:00+00') TO ('2026-02-13 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_12
    owner to postgres;

create table public.transaction_decode_p2026_02_13
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-13 00:00:00+00') TO ('2026-02-14 00:00:00+00');

alter table public.transaction_decode_p2026_02_13
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_13
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-13 00:00:00+00') TO ('2026-02-14 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_13
    owner to postgres;

create table public.block_data_decode_p2026w09
    partition of public.block_data_decode
    FOR VALUES FROM ('2026-02-23 00:00:00+00') TO ('2026-03-02 00:00:00+00');

alter table public.block_data_decode_p2026w09
    owner to postgres;

create table public.transaction_decode_p2026_02_14
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-14 00:00:00+00') TO ('2026-02-15 00:00:00+00');

alter table public.transaction_decode_p2026_02_14
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_14
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-14 00:00:00+00') TO ('2026-02-15 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_14
    owner to postgres;

create table public.aa_transaction_info_p2026w09
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2026-02-23 00:00:00+00') TO ('2026-03-02 00:00:00+00');

alter table public.aa_transaction_info_p2026w09
    owner to postgres;

create table public.aa_user_ops_calldata_p2026w09
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2026-02-23 00:00:00+00') TO ('2026-03-02 00:00:00+00');

alter table public.aa_user_ops_calldata_p2026w09
    owner to postgres;

create table public.aa_user_ops_info_p2026w09
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2026-02-23 00:00:00+00') TO ('2026-03-02 00:00:00+00');

alter table public.aa_user_ops_info_p2026w09
    owner to postgres;

create table public.transaction_decode_p2026_02_15
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-15 00:00:00+00') TO ('2026-02-16 00:00:00+00');

alter table public.transaction_decode_p2026_02_15
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_15
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-15 00:00:00+00') TO ('2026-02-16 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_15
    owner to postgres;

create table public.transaction_decode_p2026_02_16
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-16 00:00:00+00') TO ('2026-02-17 00:00:00+00');

alter table public.transaction_decode_p2026_02_16
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_16
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-16 00:00:00+00') TO ('2026-02-17 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_16
    owner to postgres;

create table public.aa_block_info_p2026_07
    partition of public.aa_block_info
    FOR VALUES FROM ('2026-07-01 00:00:00+00') TO ('2026-08-01 00:00:00+00');

alter table public.aa_block_info_p2026_07
    owner to postgres;

create table public.transaction_decode_p2026_02_17
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-17 00:00:00+00') TO ('2026-02-18 00:00:00+00');

alter table public.transaction_decode_p2026_02_17
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_17
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-17 00:00:00+00') TO ('2026-02-18 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_17
    owner to postgres;

create table public.transaction_decode_p2026_02_18
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-18 00:00:00+00') TO ('2026-02-19 00:00:00+00');

alter table public.transaction_decode_p2026_02_18
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_18
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-18 00:00:00+00') TO ('2026-02-19 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_18
    owner to postgres;

create table public.transaction_decode_p2026_02_19
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-19 00:00:00+00') TO ('2026-02-20 00:00:00+00');

alter table public.transaction_decode_p2026_02_19
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_19
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-19 00:00:00+00') TO ('2026-02-20 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_19
    owner to postgres;

create table public.transaction_decode_p2026_02_20
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-20 00:00:00+00') TO ('2026-02-21 00:00:00+00');

alter table public.transaction_decode_p2026_02_20
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_20
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-20 00:00:00+00') TO ('2026-02-21 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_20
    owner to postgres;

create table public.block_data_decode_p2026w10
    partition of public.block_data_decode
    FOR VALUES FROM ('2026-03-02 00:00:00+00') TO ('2026-03-09 00:00:00+00');

alter table public.block_data_decode_p2026w10
    owner to postgres;

create table public.transaction_decode_p2026_02_21
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-21 00:00:00+00') TO ('2026-02-22 00:00:00+00');

alter table public.transaction_decode_p2026_02_21
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_21
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-21 00:00:00+00') TO ('2026-02-22 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_21
    owner to postgres;

create table public.aa_transaction_info_p2026w10
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2026-03-02 00:00:00+00') TO ('2026-03-09 00:00:00+00');

alter table public.aa_transaction_info_p2026w10
    owner to postgres;

create table public.aa_user_ops_calldata_p2026w10
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2026-03-02 00:00:00+00') TO ('2026-03-09 00:00:00+00');

alter table public.aa_user_ops_calldata_p2026w10
    owner to postgres;

create table public.aa_user_ops_info_p2026w10
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2026-03-02 00:00:00+00') TO ('2026-03-09 00:00:00+00');

alter table public.aa_user_ops_info_p2026w10
    owner to postgres;

create table public.transaction_decode_p2026_02_22
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-22 00:00:00+00') TO ('2026-02-23 00:00:00+00');

alter table public.transaction_decode_p2026_02_22
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_22
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-22 00:00:00+00') TO ('2026-02-23 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_22
    owner to postgres;

create table public.transaction_decode_p2026_02_23
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-23 00:00:00+00') TO ('2026-02-24 00:00:00+00');

alter table public.transaction_decode_p2026_02_23
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_23
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-23 00:00:00+00') TO ('2026-02-24 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_23
    owner to postgres;

create table public.transaction_decode_p2026_02_24
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-24 00:00:00+00') TO ('2026-02-25 00:00:00+00');

alter table public.transaction_decode_p2026_02_24
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_24
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-24 00:00:00+00') TO ('2026-02-25 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_24
    owner to postgres;

create table public.transaction_decode_p2026_02_25
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-25 00:00:00+00') TO ('2026-02-26 00:00:00+00');

alter table public.transaction_decode_p2026_02_25
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_25
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-25 00:00:00+00') TO ('2026-02-26 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_25
    owner to postgres;

create table public.transaction_decode_p2026_02_26
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-26 00:00:00+00') TO ('2026-02-27 00:00:00+00');

alter table public.transaction_decode_p2026_02_26
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_26
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-26 00:00:00+00') TO ('2026-02-27 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_26
    owner to postgres;

create table public.transaction_decode_p2026_02_27
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-27 00:00:00+00') TO ('2026-02-28 00:00:00+00');

alter table public.transaction_decode_p2026_02_27
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_27
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-27 00:00:00+00') TO ('2026-02-28 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_27
    owner to postgres;

create table public.transaction_decode_p2026_02_28
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-02-28 00:00:00+00') TO ('2026-03-01 00:00:00+00');

alter table public.transaction_decode_p2026_02_28
    owner to postgres;

create table public.transaction_receipt_decode_p2026_02_28
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-02-28 00:00:00+00') TO ('2026-03-01 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_02_28
    owner to postgres;

create table public.transaction_decode_p2026_03_01
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-01 00:00:00+00') TO ('2026-03-02 00:00:00+00');

alter table public.transaction_decode_p2026_03_01
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_01
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-01 00:00:00+00') TO ('2026-03-02 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_01
    owner to postgres;

create table public.block_data_decode_p2026w11
    partition of public.block_data_decode
    FOR VALUES FROM ('2026-03-09 00:00:00+00') TO ('2026-03-16 00:00:00+00');

alter table public.block_data_decode_p2026w11
    owner to postgres;

create table public.aa_transaction_info_p2026w11
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2026-03-09 00:00:00+00') TO ('2026-03-16 00:00:00+00');

alter table public.aa_transaction_info_p2026w11
    owner to postgres;

create table public.aa_user_ops_calldata_p2026w11
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2026-03-09 00:00:00+00') TO ('2026-03-16 00:00:00+00');

alter table public.aa_user_ops_calldata_p2026w11
    owner to postgres;

create table public.aa_user_ops_info_p2026w11
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2026-03-09 00:00:00+00') TO ('2026-03-16 00:00:00+00');

alter table public.aa_user_ops_info_p2026w11
    owner to postgres;

create table public.transaction_decode_p2026_03_02
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-02 00:00:00+00') TO ('2026-03-03 00:00:00+00');

alter table public.transaction_decode_p2026_03_02
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_02
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-02 00:00:00+00') TO ('2026-03-03 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_02
    owner to postgres;

create table public.transaction_decode_p2026_03_03
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-03 00:00:00+00') TO ('2026-03-04 00:00:00+00');

alter table public.transaction_decode_p2026_03_03
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_03
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-03 00:00:00+00') TO ('2026-03-04 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_03
    owner to postgres;

create table public.transaction_decode_p2026_03_04
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-04 00:00:00+00') TO ('2026-03-05 00:00:00+00');

alter table public.transaction_decode_p2026_03_04
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_04
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-04 00:00:00+00') TO ('2026-03-05 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_04
    owner to postgres;

create table public.transaction_decode_p2026_03_05
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-05 00:00:00+00') TO ('2026-03-06 00:00:00+00');

alter table public.transaction_decode_p2026_03_05
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_05
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-05 00:00:00+00') TO ('2026-03-06 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_05
    owner to postgres;

create table public.transaction_decode_p2026_03_06
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-06 00:00:00+00') TO ('2026-03-07 00:00:00+00');

alter table public.transaction_decode_p2026_03_06
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_06
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-06 00:00:00+00') TO ('2026-03-07 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_06
    owner to postgres;

create table public.transaction_decode_p2026_03_07
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-07 00:00:00+00') TO ('2026-03-08 00:00:00+00');

alter table public.transaction_decode_p2026_03_07
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_07
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-07 00:00:00+00') TO ('2026-03-08 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_07
    owner to postgres;

create table public.transaction_decode_p2026_03_08
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-08 00:00:00+00') TO ('2026-03-09 00:00:00+00');

alter table public.transaction_decode_p2026_03_08
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_08
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-08 00:00:00+00') TO ('2026-03-09 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_08
    owner to postgres;

create table public.block_data_decode_p2026w12
    partition of public.block_data_decode
    FOR VALUES FROM ('2026-03-16 00:00:00+00') TO ('2026-03-23 00:00:00+00');

alter table public.block_data_decode_p2026w12
    owner to postgres;

create table public.aa_transaction_info_p2026w12
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2026-03-16 00:00:00+00') TO ('2026-03-23 00:00:00+00');

alter table public.aa_transaction_info_p2026w12
    owner to postgres;

create table public.aa_user_ops_calldata_p2026w12
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2026-03-16 00:00:00+00') TO ('2026-03-23 00:00:00+00');

alter table public.aa_user_ops_calldata_p2026w12
    owner to postgres;

create table public.aa_user_ops_info_p2026w12
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2026-03-16 00:00:00+00') TO ('2026-03-23 00:00:00+00');

alter table public.aa_user_ops_info_p2026w12
    owner to postgres;

create table public.transaction_decode_p2026_03_09
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-09 00:00:00+00') TO ('2026-03-10 00:00:00+00');

alter table public.transaction_decode_p2026_03_09
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_09
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-09 00:00:00+00') TO ('2026-03-10 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_09
    owner to postgres;

create table public.transaction_decode_p2026_03_10
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-10 00:00:00+00') TO ('2026-03-11 00:00:00+00');

alter table public.transaction_decode_p2026_03_10
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_10
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-10 00:00:00+00') TO ('2026-03-11 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_10
    owner to postgres;

create table public.transaction_decode_p2026_03_11
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-11 00:00:00+00') TO ('2026-03-12 00:00:00+00');

alter table public.transaction_decode_p2026_03_11
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_11
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-11 00:00:00+00') TO ('2026-03-12 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_11
    owner to postgres;

create table public.transaction_decode_p2026_03_12
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-12 00:00:00+00') TO ('2026-03-13 00:00:00+00');

alter table public.transaction_decode_p2026_03_12
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_12
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-12 00:00:00+00') TO ('2026-03-13 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_12
    owner to postgres;

create table public.transaction_decode_p2026_03_13
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-13 00:00:00+00') TO ('2026-03-14 00:00:00+00');

alter table public.transaction_decode_p2026_03_13
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_13
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-13 00:00:00+00') TO ('2026-03-14 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_13
    owner to postgres;

create table public.transaction_decode_p2026_03_14
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-14 00:00:00+00') TO ('2026-03-15 00:00:00+00');

alter table public.transaction_decode_p2026_03_14
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_14
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-14 00:00:00+00') TO ('2026-03-15 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_14
    owner to postgres;

create table public.aa_block_info_p2026_08
    partition of public.aa_block_info
    FOR VALUES FROM ('2026-08-01 00:00:00+00') TO ('2026-09-01 00:00:00+00');

alter table public.aa_block_info_p2026_08
    owner to postgres;

create table public.transaction_decode_p2026_03_15
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-15 00:00:00+00') TO ('2026-03-16 00:00:00+00');

alter table public.transaction_decode_p2026_03_15
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_15
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-15 00:00:00+00') TO ('2026-03-16 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_15
    owner to postgres;

create table public.block_data_decode_p2026w13
    partition of public.block_data_decode
    FOR VALUES FROM ('2026-03-23 00:00:00+00') TO ('2026-03-30 00:00:00+00');

alter table public.block_data_decode_p2026w13
    owner to postgres;

create table public.aa_transaction_info_p2026w13
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2026-03-23 00:00:00+00') TO ('2026-03-30 00:00:00+00');

alter table public.aa_transaction_info_p2026w13
    owner to postgres;

create table public.aa_user_ops_calldata_p2026w13
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2026-03-23 00:00:00+00') TO ('2026-03-30 00:00:00+00');

alter table public.aa_user_ops_calldata_p2026w13
    owner to postgres;

create table public.aa_user_ops_info_p2026w13
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2026-03-23 00:00:00+00') TO ('2026-03-30 00:00:00+00');

alter table public.aa_user_ops_info_p2026w13
    owner to postgres;

create table public.transaction_decode_p2026_03_16
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-16 00:00:00+00') TO ('2026-03-17 00:00:00+00');

alter table public.transaction_decode_p2026_03_16
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_16
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-16 00:00:00+00') TO ('2026-03-17 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_16
    owner to postgres;

create table public.transaction_decode_p2026_03_17
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-17 00:00:00+00') TO ('2026-03-18 00:00:00+00');

alter table public.transaction_decode_p2026_03_17
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_17
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-17 00:00:00+00') TO ('2026-03-18 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_17
    owner to postgres;

create table public.transaction_decode_p2026_03_18
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-18 00:00:00+00') TO ('2026-03-19 00:00:00+00');

alter table public.transaction_decode_p2026_03_18
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_18
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-18 00:00:00+00') TO ('2026-03-19 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_18
    owner to postgres;

create table public.transaction_decode_p2026_03_19
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-19 00:00:00+00') TO ('2026-03-20 00:00:00+00');

alter table public.transaction_decode_p2026_03_19
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_19
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-19 00:00:00+00') TO ('2026-03-20 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_19
    owner to postgres;

create table public.transaction_decode_p2026_03_20
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-20 00:00:00+00') TO ('2026-03-21 00:00:00+00');

alter table public.transaction_decode_p2026_03_20
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_20
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-20 00:00:00+00') TO ('2026-03-21 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_20
    owner to postgres;

create table public.transaction_decode_p2026_03_21
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-21 00:00:00+00') TO ('2026-03-22 00:00:00+00');

alter table public.transaction_decode_p2026_03_21
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_21
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-21 00:00:00+00') TO ('2026-03-22 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_21
    owner to postgres;

create table public.transaction_decode_p2026_03_22
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-22 00:00:00+00') TO ('2026-03-23 00:00:00+00');

alter table public.transaction_decode_p2026_03_22
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_22
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-22 00:00:00+00') TO ('2026-03-23 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_22
    owner to postgres;

create table public.block_data_decode_p2026w14
    partition of public.block_data_decode
    FOR VALUES FROM ('2026-03-30 00:00:00+00') TO ('2026-04-06 00:00:00+00');

alter table public.block_data_decode_p2026w14
    owner to postgres;

create table public.aa_transaction_info_p2026w14
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2026-03-30 00:00:00+00') TO ('2026-04-06 00:00:00+00');

alter table public.aa_transaction_info_p2026w14
    owner to postgres;

create table public.aa_user_ops_calldata_p2026w14
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2026-03-30 00:00:00+00') TO ('2026-04-06 00:00:00+00');

alter table public.aa_user_ops_calldata_p2026w14
    owner to postgres;

create table public.aa_user_ops_info_p2026w14
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2026-03-30 00:00:00+00') TO ('2026-04-06 00:00:00+00');

alter table public.aa_user_ops_info_p2026w14
    owner to postgres;

create table public.transaction_decode_p2026_03_23
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-23 00:00:00+00') TO ('2026-03-24 00:00:00+00');

alter table public.transaction_decode_p2026_03_23
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_23
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-23 00:00:00+00') TO ('2026-03-24 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_23
    owner to postgres;

create table public.transaction_decode_p2026_03_24
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-24 00:00:00+00') TO ('2026-03-25 00:00:00+00');

alter table public.transaction_decode_p2026_03_24
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_24
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-24 00:00:00+00') TO ('2026-03-25 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_24
    owner to postgres;

create table public.transaction_decode_p2026_03_25
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-25 00:00:00+00') TO ('2026-03-26 00:00:00+00');

alter table public.transaction_decode_p2026_03_25
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_25
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-25 00:00:00+00') TO ('2026-03-26 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_25
    owner to postgres;

create table public.transaction_decode_p2026_03_26
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-26 00:00:00+00') TO ('2026-03-27 00:00:00+00');

alter table public.transaction_decode_p2026_03_26
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_26
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-26 00:00:00+00') TO ('2026-03-27 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_26
    owner to postgres;

create table public.transaction_decode_p2026_03_27
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-27 00:00:00+00') TO ('2026-03-28 00:00:00+00');

alter table public.transaction_decode_p2026_03_27
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_27
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-27 00:00:00+00') TO ('2026-03-28 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_27
    owner to postgres;

create table public.transaction_decode_p2026_03_28
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-28 00:00:00+00') TO ('2026-03-29 00:00:00+00');

alter table public.transaction_decode_p2026_03_28
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_28
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-28 00:00:00+00') TO ('2026-03-29 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_28
    owner to postgres;

create table public.transaction_decode_p2026_03_29
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-29 00:00:00+00') TO ('2026-03-30 00:00:00+00');

alter table public.transaction_decode_p2026_03_29
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_29
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-29 00:00:00+00') TO ('2026-03-30 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_29
    owner to postgres;

create table public.block_data_decode_p2026w15
    partition of public.block_data_decode
    FOR VALUES FROM ('2026-04-06 00:00:00+00') TO ('2026-04-13 00:00:00+00');

alter table public.block_data_decode_p2026w15
    owner to postgres;

create table public.aa_transaction_info_p2026w15
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2026-04-06 00:00:00+00') TO ('2026-04-13 00:00:00+00');

alter table public.aa_transaction_info_p2026w15
    owner to postgres;

create table public.aa_user_ops_calldata_p2026w15
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2026-04-06 00:00:00+00') TO ('2026-04-13 00:00:00+00');

alter table public.aa_user_ops_calldata_p2026w15
    owner to postgres;

create table public.aa_user_ops_info_p2026w15
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2026-04-06 00:00:00+00') TO ('2026-04-13 00:00:00+00');

alter table public.aa_user_ops_info_p2026w15
    owner to postgres;

create table public.transaction_decode_p2026_03_30
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-30 00:00:00+00') TO ('2026-03-31 00:00:00+00');

alter table public.transaction_decode_p2026_03_30
    owner to postgres;

create table public.transaction_decode_p2026_03_31
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-03-31 00:00:00+00') TO ('2026-04-01 00:00:00+00');

alter table public.transaction_decode_p2026_03_31
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_30
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-30 00:00:00+00') TO ('2026-03-31 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_30
    owner to postgres;

create table public.transaction_receipt_decode_p2026_03_31
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-03-31 00:00:00+00') TO ('2026-04-01 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_03_31
    owner to postgres;

create table public.transaction_decode_p2026_04_01
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-01 00:00:00+00') TO ('2026-04-02 00:00:00+00');

alter table public.transaction_decode_p2026_04_01
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_01
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-01 00:00:00+00') TO ('2026-04-02 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_01
    owner to postgres;

create table public.transaction_decode_p2026_04_02
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-02 00:00:00+00') TO ('2026-04-03 00:00:00+00');

alter table public.transaction_decode_p2026_04_02
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_02
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-02 00:00:00+00') TO ('2026-04-03 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_02
    owner to postgres;

create table public.transaction_decode_p2026_04_03
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-03 00:00:00+00') TO ('2026-04-04 00:00:00+00');

alter table public.transaction_decode_p2026_04_03
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_03
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-03 00:00:00+00') TO ('2026-04-04 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_03
    owner to postgres;

create table public.block_data_decode_p2026w16
    partition of public.block_data_decode
    FOR VALUES FROM ('2026-04-13 00:00:00+00') TO ('2026-04-20 00:00:00+00');

alter table public.block_data_decode_p2026w16
    owner to postgres;

create table public.transaction_decode_p2026_04_04
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-04 00:00:00+00') TO ('2026-04-05 00:00:00+00');

alter table public.transaction_decode_p2026_04_04
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_04
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-04 00:00:00+00') TO ('2026-04-05 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_04
    owner to postgres;

create table public.aa_transaction_info_p2026w16
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2026-04-13 00:00:00+00') TO ('2026-04-20 00:00:00+00');

alter table public.aa_transaction_info_p2026w16
    owner to postgres;

create table public.aa_user_ops_calldata_p2026w16
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2026-04-13 00:00:00+00') TO ('2026-04-20 00:00:00+00');

alter table public.aa_user_ops_calldata_p2026w16
    owner to postgres;

create table public.aa_user_ops_info_p2026w16
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2026-04-13 00:00:00+00') TO ('2026-04-20 00:00:00+00');

alter table public.aa_user_ops_info_p2026w16
    owner to postgres;

create table public.transaction_decode_p2026_04_05
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-05 00:00:00+00') TO ('2026-04-06 00:00:00+00');

alter table public.transaction_decode_p2026_04_05
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_05
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-05 00:00:00+00') TO ('2026-04-06 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_05
    owner to postgres;

create table public.transaction_decode_p2026_04_06
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-06 00:00:00+00') TO ('2026-04-07 00:00:00+00');

alter table public.transaction_decode_p2026_04_06
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_06
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-06 00:00:00+00') TO ('2026-04-07 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_06
    owner to postgres;

create table public.transaction_decode_p2026_04_07
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-07 00:00:00+00') TO ('2026-04-08 00:00:00+00');

alter table public.transaction_decode_p2026_04_07
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_07
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-07 00:00:00+00') TO ('2026-04-08 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_07
    owner to postgres;

create table public.transaction_decode_p2026_04_08
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-08 00:00:00+00') TO ('2026-04-09 00:00:00+00');

alter table public.transaction_decode_p2026_04_08
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_08
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-08 00:00:00+00') TO ('2026-04-09 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_08
    owner to postgres;

create table public.transaction_decode_p2026_04_09
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-09 00:00:00+00') TO ('2026-04-10 00:00:00+00');

alter table public.transaction_decode_p2026_04_09
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_09
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-09 00:00:00+00') TO ('2026-04-10 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_09
    owner to postgres;

create table public.transaction_decode_p2026_04_10
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-10 00:00:00+00') TO ('2026-04-11 00:00:00+00');

alter table public.transaction_decode_p2026_04_10
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_10
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-10 00:00:00+00') TO ('2026-04-11 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_10
    owner to postgres;

create table public.block_data_decode_p2026w17
    partition of public.block_data_decode
    FOR VALUES FROM ('2026-04-20 00:00:00+00') TO ('2026-04-27 00:00:00+00');

alter table public.block_data_decode_p2026w17
    owner to postgres;

create table public.transaction_decode_p2026_04_11
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-11 00:00:00+00') TO ('2026-04-12 00:00:00+00');

alter table public.transaction_decode_p2026_04_11
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_11
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-11 00:00:00+00') TO ('2026-04-12 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_11
    owner to postgres;

create table public.aa_transaction_info_p2026w17
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2026-04-20 00:00:00+00') TO ('2026-04-27 00:00:00+00');

alter table public.aa_transaction_info_p2026w17
    owner to postgres;

create table public.aa_user_ops_calldata_p2026w17
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2026-04-20 00:00:00+00') TO ('2026-04-27 00:00:00+00');

alter table public.aa_user_ops_calldata_p2026w17
    owner to postgres;

create table public.aa_user_ops_info_p2026w17
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2026-04-20 00:00:00+00') TO ('2026-04-27 00:00:00+00');

alter table public.aa_user_ops_info_p2026w17
    owner to postgres;

create table public.transaction_decode_p2026_04_12
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-12 00:00:00+00') TO ('2026-04-13 00:00:00+00');

alter table public.transaction_decode_p2026_04_12
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_12
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-12 00:00:00+00') TO ('2026-04-13 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_12
    owner to postgres;

create table public.transaction_decode_p2026_04_13
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-13 00:00:00+00') TO ('2026-04-14 00:00:00+00');

alter table public.transaction_decode_p2026_04_13
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_13
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-13 00:00:00+00') TO ('2026-04-14 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_13
    owner to postgres;

create table public.transaction_decode_p2026_04_14
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-14 00:00:00+00') TO ('2026-04-15 00:00:00+00');

alter table public.transaction_decode_p2026_04_14
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_14
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-14 00:00:00+00') TO ('2026-04-15 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_14
    owner to postgres;

create table public.transaction_decode_p2026_04_15
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-15 00:00:00+00') TO ('2026-04-16 00:00:00+00');

alter table public.transaction_decode_p2026_04_15
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_15
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-15 00:00:00+00') TO ('2026-04-16 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_15
    owner to postgres;

create table public.transaction_decode_p2026_04_16
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-16 00:00:00+00') TO ('2026-04-17 00:00:00+00');

alter table public.transaction_decode_p2026_04_16
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_16
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-16 00:00:00+00') TO ('2026-04-17 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_16
    owner to postgres;

create table public.aa_block_info_p2026_09
    partition of public.aa_block_info
    FOR VALUES FROM ('2026-09-01 00:00:00+00') TO ('2026-10-01 00:00:00+00');

alter table public.aa_block_info_p2026_09
    owner to postgres;

create table public.transaction_decode_p2026_04_17
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-17 00:00:00+00') TO ('2026-04-18 00:00:00+00');

alter table public.transaction_decode_p2026_04_17
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_17
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-17 00:00:00+00') TO ('2026-04-18 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_17
    owner to postgres;

create table public.block_data_decode_p2026w18
    partition of public.block_data_decode
    FOR VALUES FROM ('2026-04-27 00:00:00+00') TO ('2026-05-04 00:00:00+00');

alter table public.block_data_decode_p2026w18
    owner to postgres;

create table public.transaction_decode_p2026_04_18
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-18 00:00:00+00') TO ('2026-04-19 00:00:00+00');

alter table public.transaction_decode_p2026_04_18
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_18
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-18 00:00:00+00') TO ('2026-04-19 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_18
    owner to postgres;

create table public.aa_transaction_info_p2026w18
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2026-04-27 00:00:00+00') TO ('2026-05-04 00:00:00+00');

alter table public.aa_transaction_info_p2026w18
    owner to postgres;

create table public.aa_user_ops_calldata_p2026w18
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2026-04-27 00:00:00+00') TO ('2026-05-04 00:00:00+00');

alter table public.aa_user_ops_calldata_p2026w18
    owner to postgres;

create table public.aa_user_ops_info_p2026w18
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2026-04-27 00:00:00+00') TO ('2026-05-04 00:00:00+00');

alter table public.aa_user_ops_info_p2026w18
    owner to postgres;

create table public.transaction_decode_p2026_04_19
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-19 00:00:00+00') TO ('2026-04-20 00:00:00+00');

alter table public.transaction_decode_p2026_04_19
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_19
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-19 00:00:00+00') TO ('2026-04-20 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_19
    owner to postgres;

create table public.transaction_decode_p2026_04_20
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-20 00:00:00+00') TO ('2026-04-21 00:00:00+00');

alter table public.transaction_decode_p2026_04_20
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_20
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-20 00:00:00+00') TO ('2026-04-21 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_20
    owner to postgres;

create table public.transaction_decode_p2026_04_21
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-21 00:00:00+00') TO ('2026-04-22 00:00:00+00');

alter table public.transaction_decode_p2026_04_21
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_21
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-21 00:00:00+00') TO ('2026-04-22 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_21
    owner to postgres;

create table public.transaction_decode_p2026_04_22
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-22 00:00:00+00') TO ('2026-04-23 00:00:00+00');

alter table public.transaction_decode_p2026_04_22
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_22
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-22 00:00:00+00') TO ('2026-04-23 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_22
    owner to postgres;

create table public.transaction_decode_p2026_04_23
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-23 00:00:00+00') TO ('2026-04-24 00:00:00+00');

alter table public.transaction_decode_p2026_04_23
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_23
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-23 00:00:00+00') TO ('2026-04-24 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_23
    owner to postgres;

create table public.transaction_decode_p2026_04_24
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-24 00:00:00+00') TO ('2026-04-25 00:00:00+00');

alter table public.transaction_decode_p2026_04_24
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_24
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-24 00:00:00+00') TO ('2026-04-25 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_24
    owner to postgres;

create table public.block_data_decode_p2026w19
    partition of public.block_data_decode
    FOR VALUES FROM ('2026-05-04 00:00:00+00') TO ('2026-05-11 00:00:00+00');

alter table public.block_data_decode_p2026w19
    owner to postgres;

create table public.transaction_decode_p2026_04_25
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-25 00:00:00+00') TO ('2026-04-26 00:00:00+00');

alter table public.transaction_decode_p2026_04_25
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_25
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-25 00:00:00+00') TO ('2026-04-26 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_25
    owner to postgres;

create table public.aa_transaction_info_p2026w19
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2026-05-04 00:00:00+00') TO ('2026-05-11 00:00:00+00');

alter table public.aa_transaction_info_p2026w19
    owner to postgres;

create table public.aa_user_ops_calldata_p2026w19
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2026-05-04 00:00:00+00') TO ('2026-05-11 00:00:00+00');

alter table public.aa_user_ops_calldata_p2026w19
    owner to postgres;

create table public.aa_user_ops_info_p2026w19
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2026-05-04 00:00:00+00') TO ('2026-05-11 00:00:00+00');

alter table public.aa_user_ops_info_p2026w19
    owner to postgres;

create table public.transaction_decode_p2026_04_26
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-26 00:00:00+00') TO ('2026-04-27 00:00:00+00');

alter table public.transaction_decode_p2026_04_26
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_26
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-26 00:00:00+00') TO ('2026-04-27 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_26
    owner to postgres;

create table public.transaction_decode_p2026_04_27
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-27 00:00:00+00') TO ('2026-04-28 00:00:00+00');

alter table public.transaction_decode_p2026_04_27
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_27
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-27 00:00:00+00') TO ('2026-04-28 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_27
    owner to postgres;

create table public.transaction_decode_p2026_04_28
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-28 00:00:00+00') TO ('2026-04-29 00:00:00+00');

alter table public.transaction_decode_p2026_04_28
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_28
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-28 00:00:00+00') TO ('2026-04-29 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_28
    owner to postgres;

create table public.transaction_decode_p2026_04_29
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-29 00:00:00+00') TO ('2026-04-30 00:00:00+00');

alter table public.transaction_decode_p2026_04_29
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_29
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-29 00:00:00+00') TO ('2026-04-30 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_29
    owner to postgres;

create table public.transaction_decode_p2026_04_30
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-04-30 00:00:00+00') TO ('2026-05-01 00:00:00+00');

alter table public.transaction_decode_p2026_04_30
    owner to postgres;

create table public.transaction_receipt_decode_p2026_04_30
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-04-30 00:00:00+00') TO ('2026-05-01 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_04_30
    owner to postgres;

create table public.transaction_decode_p2026_05_01
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-05-01 00:00:00+00') TO ('2026-05-02 00:00:00+00');

alter table public.transaction_decode_p2026_05_01
    owner to postgres;

create table public.transaction_receipt_decode_p2026_05_01
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-05-01 00:00:00+00') TO ('2026-05-02 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_05_01
    owner to postgres;

create table public.transaction_decode_p2026_05_02
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-05-02 00:00:00+00') TO ('2026-05-03 00:00:00+00');

alter table public.transaction_decode_p2026_05_02
    owner to postgres;

create table public.transaction_receipt_decode_p2026_05_02
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-05-02 00:00:00+00') TO ('2026-05-03 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_05_02
    owner to postgres;

create table public.block_data_decode_p2026w20
    partition of public.block_data_decode
    FOR VALUES FROM ('2026-05-11 00:00:00+00') TO ('2026-05-18 00:00:00+00');

alter table public.block_data_decode_p2026w20
    owner to postgres;

create table public.transaction_decode_p2026_05_03
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-05-03 00:00:00+00') TO ('2026-05-04 00:00:00+00');

alter table public.transaction_decode_p2026_05_03
    owner to postgres;

create table public.transaction_receipt_decode_p2026_05_03
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-05-03 00:00:00+00') TO ('2026-05-04 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_05_03
    owner to postgres;

create table public.aa_transaction_info_p2026w20
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2026-05-11 00:00:00+00') TO ('2026-05-18 00:00:00+00');

alter table public.aa_transaction_info_p2026w20
    owner to postgres;

create table public.aa_user_ops_calldata_p2026w20
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2026-05-11 00:00:00+00') TO ('2026-05-18 00:00:00+00');

alter table public.aa_user_ops_calldata_p2026w20
    owner to postgres;

create table public.aa_user_ops_info_p2026w20
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2026-05-11 00:00:00+00') TO ('2026-05-18 00:00:00+00');

alter table public.aa_user_ops_info_p2026w20
    owner to postgres;

create table public.transaction_decode_p2026_05_04
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-05-04 00:00:00+00') TO ('2026-05-05 00:00:00+00');

alter table public.transaction_decode_p2026_05_04
    owner to postgres;

create table public.transaction_receipt_decode_p2026_05_04
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-05-04 00:00:00+00') TO ('2026-05-05 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_05_04
    owner to postgres;

create table public.transaction_decode_p2026_05_05
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-05-05 00:00:00+00') TO ('2026-05-06 00:00:00+00');

alter table public.transaction_decode_p2026_05_05
    owner to postgres;

create table public.transaction_receipt_decode_p2026_05_05
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-05-05 00:00:00+00') TO ('2026-05-06 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_05_05
    owner to postgres;

create table public.transaction_decode_p2026_05_06
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-05-06 00:00:00+00') TO ('2026-05-07 00:00:00+00');

alter table public.transaction_decode_p2026_05_06
    owner to postgres;

create table public.transaction_receipt_decode_p2026_05_06
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-05-06 00:00:00+00') TO ('2026-05-07 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_05_06
    owner to postgres;

create table public.transaction_decode_p2026_05_07
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-05-07 00:00:00+00') TO ('2026-05-08 00:00:00+00');

alter table public.transaction_decode_p2026_05_07
    owner to postgres;

create table public.transaction_receipt_decode_p2026_05_07
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-05-07 00:00:00+00') TO ('2026-05-08 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_05_07
    owner to postgres;

create table public.transaction_decode_p2026_05_08
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-05-08 00:00:00+00') TO ('2026-05-09 00:00:00+00');

alter table public.transaction_decode_p2026_05_08
    owner to postgres;

create table public.transaction_receipt_decode_p2026_05_08
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-05-08 00:00:00+00') TO ('2026-05-09 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_05_08
    owner to postgres;

create table public.transaction_decode_p2026_05_09
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-05-09 00:00:00+00') TO ('2026-05-10 00:00:00+00');

alter table public.transaction_decode_p2026_05_09
    owner to postgres;

create table public.transaction_receipt_decode_p2026_05_09
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-05-09 00:00:00+00') TO ('2026-05-10 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_05_09
    owner to postgres;

create table public.block_data_decode_p2026w21
    partition of public.block_data_decode
    FOR VALUES FROM ('2026-05-18 00:00:00+00') TO ('2026-05-25 00:00:00+00');

alter table public.block_data_decode_p2026w21
    owner to postgres;

create table public.transaction_decode_p2026_05_10
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-05-10 00:00:00+00') TO ('2026-05-11 00:00:00+00');

alter table public.transaction_decode_p2026_05_10
    owner to postgres;

create table public.transaction_receipt_decode_p2026_05_10
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-05-10 00:00:00+00') TO ('2026-05-11 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_05_10
    owner to postgres;

create table public.aa_transaction_info_p2026w21
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2026-05-18 00:00:00+00') TO ('2026-05-25 00:00:00+00');

alter table public.aa_transaction_info_p2026w21
    owner to postgres;

create table public.aa_user_ops_calldata_p2026w21
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2026-05-18 00:00:00+00') TO ('2026-05-25 00:00:00+00');

alter table public.aa_user_ops_calldata_p2026w21
    owner to postgres;

create table public.aa_user_ops_info_p2026w21
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2026-05-18 00:00:00+00') TO ('2026-05-25 00:00:00+00');

alter table public.aa_user_ops_info_p2026w21
    owner to postgres;

create table public.transaction_decode_p2026_05_11
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-05-11 00:00:00+00') TO ('2026-05-12 00:00:00+00');

alter table public.transaction_decode_p2026_05_11
    owner to postgres;

create table public.transaction_receipt_decode_p2026_05_11
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-05-11 00:00:00+00') TO ('2026-05-12 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_05_11
    owner to postgres;

create table public.transaction_decode_p2026_05_12
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-05-12 00:00:00+00') TO ('2026-05-13 00:00:00+00');

alter table public.transaction_decode_p2026_05_12
    owner to postgres;

create table public.transaction_receipt_decode_p2026_05_12
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-05-12 00:00:00+00') TO ('2026-05-13 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_05_12
    owner to postgres;

create table public.transaction_decode_p2026_05_13
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-05-13 00:00:00+00') TO ('2026-05-14 00:00:00+00');

alter table public.transaction_decode_p2026_05_13
    owner to postgres;

create table public.transaction_receipt_decode_p2026_05_13
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-05-13 00:00:00+00') TO ('2026-05-14 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_05_13
    owner to postgres;

create table public.transaction_decode_p2026_05_14
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-05-14 00:00:00+00') TO ('2026-05-15 00:00:00+00');

alter table public.transaction_decode_p2026_05_14
    owner to postgres;

create table public.transaction_receipt_decode_p2026_05_14
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-05-14 00:00:00+00') TO ('2026-05-15 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_05_14
    owner to postgres;

create table public.transaction_decode_p2026_05_15
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-05-15 00:00:00+00') TO ('2026-05-16 00:00:00+00');

alter table public.transaction_decode_p2026_05_15
    owner to postgres;

create table public.transaction_receipt_decode_p2026_05_15
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-05-15 00:00:00+00') TO ('2026-05-16 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_05_15
    owner to postgres;

create table public.transaction_decode_p2026_05_16
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-05-16 00:00:00+00') TO ('2026-05-17 00:00:00+00');

alter table public.transaction_decode_p2026_05_16
    owner to postgres;

create table public.transaction_receipt_decode_p2026_05_16
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-05-16 00:00:00+00') TO ('2026-05-17 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_05_16
    owner to postgres;

create table public.aa_block_info_p2026_10
    partition of public.aa_block_info
    FOR VALUES FROM ('2026-10-01 00:00:00+00') TO ('2026-11-01 00:00:00+00');

alter table public.aa_block_info_p2026_10
    owner to postgres;

create table public.block_data_decode_p2026w22
    partition of public.block_data_decode
    FOR VALUES FROM ('2026-05-25 00:00:00+00') TO ('2026-06-01 00:00:00+00');

alter table public.block_data_decode_p2026w22
    owner to postgres;

create table public.transaction_decode_p2026_05_17
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-05-17 00:00:00+00') TO ('2026-05-18 00:00:00+00');

alter table public.transaction_decode_p2026_05_17
    owner to postgres;

create table public.transaction_receipt_decode_p2026_05_17
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-05-17 00:00:00+00') TO ('2026-05-18 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_05_17
    owner to postgres;

create table public.aa_transaction_info_p2026w22
    partition of public.aa_transaction_info
    FOR VALUES FROM ('2026-05-25 00:00:00+00') TO ('2026-06-01 00:00:00+00');

alter table public.aa_transaction_info_p2026w22
    owner to postgres;

create table public.aa_user_ops_calldata_p2026w22
    partition of public.aa_user_ops_calldata
    FOR VALUES FROM ('2026-05-25 00:00:00+00') TO ('2026-06-01 00:00:00+00');

alter table public.aa_user_ops_calldata_p2026w22
    owner to postgres;

create table public.aa_user_ops_info_p2026w22
    partition of public.aa_user_ops_info
    FOR VALUES FROM ('2026-05-25 00:00:00+00') TO ('2026-06-01 00:00:00+00');

alter table public.aa_user_ops_info_p2026w22
    owner to postgres;

create table public.transaction_decode_p2026_05_18
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-05-18 00:00:00+00') TO ('2026-05-19 00:00:00+00');

alter table public.transaction_decode_p2026_05_18
    owner to postgres;

create table public.transaction_receipt_decode_p2026_05_18
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-05-18 00:00:00+00') TO ('2026-05-19 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_05_18
    owner to postgres;

create table public.transaction_decode_p2026_05_19
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-05-19 00:00:00+00') TO ('2026-05-20 00:00:00+00');

alter table public.transaction_decode_p2026_05_19
    owner to postgres;

create table public.transaction_receipt_decode_p2026_05_19
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-05-19 00:00:00+00') TO ('2026-05-20 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_05_19
    owner to postgres;

create table public.transaction_decode_p2026_05_20
    partition of public.transaction_decode
    FOR VALUES FROM ('2026-05-20 00:00:00+00') TO ('2026-05-21 00:00:00+00');

alter table public.transaction_decode_p2026_05_20
    owner to postgres;

create table public.transaction_receipt_decode_p2026_05_20
    partition of public.transaction_receipt_decode
    FOR VALUES FROM ('2026-05-20 00:00:00+00') TO ('2026-05-21 00:00:00+00');

alter table public.transaction_receipt_decode_p2026_05_20
    owner to postgres;

