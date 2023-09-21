CREATE
EXTENSION IF NOT EXISTS timescaledb;

create table aa_transaction_info
(
    time               timestamp with time zone not null,
    create_time        timestamp with time zone,
    hash               text,
    block_hash         text,
    block_number       int8,
    userop_count       int,
    is_mev             boolean,
    bundler_profit     numeric,
    bundler_profit_usd numeric default 0
);

SELECT create_hypertable('aa_transaction_info', 'time');
CREATE INDEX aa_transaction_info_hash_index ON aa_transaction_info USING HASH (hash);
CREATE INDEX aa_transaction_info_block_num_index ON aa_transaction_info (block_number);
CREATE INDEX aa_transaction_info_create_time_index ON aa_transaction_info (create_time);
CREATE UNIQUE INDEX aa_transaction_info_time_hash_index ON aa_transaction_info (time, hash);


create table aa_user_ops_calldata
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
);

SELECT create_hypertable('aa_user_ops_calldata', 'time');
CREATE INDEX aa_user_ops_calldata_tx_hash_index ON aa_user_ops_calldata USING HASH (tx_hash);
CREATE INDEX aa_user_ops_calldata_user_operation_hash_index ON aa_user_ops_calldata USING HASH (user_ops_hash);
CREATE INDEX aa_user_ops_calldata_block_num_index ON aa_user_ops_calldata (block_number);
CREATE UNIQUE INDEX aa_user_ops_calldata_time_uuid_index ON aa_user_ops_calldata (time, uuid);


drop table if exists public.aa_user_ops_info;
create table if not exists public.aa_user_ops_info
(
    time
    timestamp
    with
    time
    zone
    not
    null,
    user_operation_hash
    varchar
(
    128
),
    tx_hash varchar
(
    128
),
    block_number bigint,
    network varchar
(
    128
),
    sender varchar
(
    64
),
    target varchar
(
    64
),
    tx_value numeric,
    fee numeric,
    bundler varchar
(
    64
),
    entry_point varchar
(
    64
),
    factory varchar
(
    64
),
    paymaster varchar
(
    64
),
    paymaster_and_data text,
    signature text,
    calldata text,
    calldata_contract varchar
(
    64
),
    nonce bigint,
    call_gas_limit bigint,
    pre_verification_gas bigint,
    verification_gas_limit bigint,
    max_fee_per_gas bigint,
    max_priority_fee_per_gas bigint,
    tx_time bigint,
    init_code text,
    status integer,
    source varchar
(
    128
),
    actual_gas_cost bigint,
    actual_gas_used bigint,
    create_time timestamp with time zone,
    usd_amount numeric,
    update_time timestamp with time zone,
                              targets varchar (128)[],
    aa_index integer default 0,
    targets_count integer default 0,
    fee_usd numeric default 0,
    tx_value_usd numeric default 0
    );

SELECT create_hypertable('aa_user_ops_info', 'time');
CREATE INDEX aa_user_ops_info_tx_hash_index ON aa_user_ops_info USING HASH (tx_hash);
CREATE INDEX aa_user_ops_info_user_operation_hash_index ON aa_user_ops_info USING HASH (user_operation_hash);
CREATE INDEX aa_user_ops_info_block_num_index ON aa_user_ops_info (block_number);
CREATE UNIQUE INDEX aa_user_ops_info_time_uuid_index ON aa_user_ops_info (time, tx_hash, user_operation_hash);
create index aa_user_ops_info_bundler_index on aa_user_ops_info (bundler);
create index aa_user_ops_info_paymaster_index on aa_user_ops_info (paymaster);
create index aa_user_ops_info_factory_index on aa_user_ops_info (factory);

-- drop table  if exists  aa_account_data;
create table aa_account_data
(
    address           text primary key,
    aa_type           text,
    factory           text,
    factory_time      timestamptz,
    user_ops_num      bigint          default 0,
    update_time       timestamp with time zone,
    total_balance_usd numeric(50, 20) default 0,
    last_time         bigint          default 0
);
create index aa_account_data_aa_type on aa_account_data (aa_type);
create index aa_account_data_factory on aa_account_data (factory);
