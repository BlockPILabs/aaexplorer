-- CREATE
--     EXTENSION IF NOT EXISTS timescaledb;
-- alter table public.aa_account_data
--     add user_ops_num bigint default 0;
-- alter table public.aa_account_data
--     add total_balance_usd numeric(50, 20) default 0;
-- alter table public.aa_account_data
--     add last_time bigint default 0;
-- alter table public.aa_account_data
--     add update_time timestamp with time zone;


alter table public.account
    add update_time timestamp with time zone;
alter table public.aa_block_info
    add bundler_profit_usd numeric default 0;
-- alter table public.aa_transaction_info
--     add bundler_profit_usd numeric default 0;
-- alter table public.aa_user_ops_info
--     add fee_usd numeric default 0;
-- alter table public.aa_user_ops_info
--     add tx_value_usd numeric default 0;


create table aa_transaction_info
(
    time           timestamp with time zone not null,
    create_time    timestamp with time zone,
    hash           text,
    block_hash     text,
    block_number   int8,
    userop_count   int,
    is_mev         boolean,
    bundler_profit numeric,
    bundler_profit_usd numeric default 0
) PARTITION BY RANGE (time);


SELECT partman.create_parent(
               p_parent_table=>'public.aa_transaction_info',
               p_control=>'time',
               p_type=>'native',
               p_interval=>'daily',
               p_premake=>30,
               p_start_partition=>'2023-01-01'
           );
UPDATE partman.part_config
SET infinite_time_partitions = TRUE
WHERE parent_table = 'public.aa_transaction_info';
-- SELECT create_hypertable('aa_transaction_info', 'time');
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
) PARTITION BY RANGE (time);


SELECT partman.create_parent(
               p_parent_table=>'public.aa_user_ops_calldata',
               p_control=>'time',
               p_type=>'native',
               p_interval=>'daily',
               p_premake=>30,
               p_start_partition=>'2023-01-01'
           );
UPDATE partman.part_config
SET infinite_time_partitions = TRUE
WHERE parent_table = 'public.aa_user_ops_calldata';

-- SELECT create_hypertable('aa_user_ops_calldata', 'time');
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
                             varchar(128),
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
) PARTITION BY RANGE (time);


SELECT partman.create_parent(
               p_parent_table=>'public.aa_user_ops_info',
               p_control=>'time',
               p_type=>'native',
               p_interval=>'daily',
               p_premake=>30,
               p_start_partition=>'2023-01-01'
           );
UPDATE partman.part_config
SET infinite_time_partitions = TRUE
WHERE parent_table = 'public.aa_user_ops_info';

-- SELECT create_hypertable('aa_user_ops_info', 'time');
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
) PARTITION BY hash (address);
create table aa_account_data_p1 partition of aa_account_data for values with (modulus 30, remainder 0);
create table aa_account_data_p2 partition of aa_account_data for values with (modulus 30, remainder 1);
create table aa_account_data_p3 partition of aa_account_data for values with (modulus 30, remainder 2);
create table aa_account_data_p4 partition of aa_account_data for values with (modulus 30, remainder 3);
create table aa_account_data_p5 partition of aa_account_data for values with (modulus 30, remainder 4);
create table aa_account_data_p6 partition of aa_account_data for values with (modulus 30, remainder 5);
create table aa_account_data_p7 partition of aa_account_data for values with (modulus 30, remainder 6);
create table aa_account_data_p8 partition of aa_account_data for values with (modulus 30, remainder 7);
create table aa_account_data_p9 partition of aa_account_data for values with (modulus 30, remainder 8);
create table aa_account_data_p10 partition of aa_account_data for values with (modulus 30, remainder 9);
create table aa_account_data_p11 partition of aa_account_data for values with (modulus 30, remainder 10);
create table aa_account_data_p12 partition of aa_account_data for values with (modulus 30, remainder 11);
create table aa_account_data_p13 partition of aa_account_data for values with (modulus 30, remainder 12);
create table aa_account_data_p14 partition of aa_account_data for values with (modulus 30, remainder 13);
create table aa_account_data_p15 partition of aa_account_data for values with (modulus 30, remainder 14);
create table aa_account_data_p16 partition of aa_account_data for values with (modulus 30, remainder 15);
create table aa_account_data_p17 partition of aa_account_data for values with (modulus 30, remainder 16);
create table aa_account_data_p18 partition of aa_account_data for values with (modulus 30, remainder 17);
create table aa_account_data_p19 partition of aa_account_data for values with (modulus 30, remainder 18);
create table aa_account_data_p20 partition of aa_account_data for values with (modulus 30, remainder 19);
create table aa_account_data_p21 partition of aa_account_data for values with (modulus 30, remainder 20);
create table aa_account_data_p22 partition of aa_account_data for values with (modulus 30, remainder 21);
create table aa_account_data_p23 partition of aa_account_data for values with (modulus 30, remainder 22);
create table aa_account_data_p24 partition of aa_account_data for values with (modulus 30, remainder 23);
create table aa_account_data_p25 partition of aa_account_data for values with (modulus 30, remainder 24);
create table aa_account_data_p26 partition of aa_account_data for values with (modulus 30, remainder 25);
create table aa_account_data_p27 partition of aa_account_data for values with (modulus 30, remainder 26);
create table aa_account_data_p28 partition of aa_account_data for values with (modulus 30, remainder 27);
create table aa_account_data_p29 partition of aa_account_data for values with (modulus 30, remainder 28);
create table aa_account_data_p30 partition of aa_account_data for values with (modulus 30, remainder 29);


create index aa_account_data_aa_type on aa_account_data (aa_type);
create index aa_account_data_factory on aa_account_data (factory);

