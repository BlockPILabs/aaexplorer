---- 20230804

CREATE
    EXTENSION IF NOT EXISTS timescaledb;

create table block_data_decode
(
    time              timestamptz,
    create_time       timestamp with time zone,
    number            numeric,
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
    transaction_count int8,
    uncles            text[],
    base_fee_per_gas  numeric
);

SELECT create_hypertable('block_data_decode', 'time');
CREATE INDEX block_data_decode_hash_index ON block_data_decode USING HASH (hash);
CREATE INDEX block_data_decode_block_num_index ON block_data_decode (number);
CREATE INDEX block_data_decode_create_time_index ON block_data_decode (create_time);
CREATE UNIQUE INDEX block_data_decode_time_hash_index ON block_data_decode (time, number);


----------------------------------------------------------


create table aa_block_info
(
    time             timestamptz,
    create_time      timestamp with time zone,
    number           numeric,
    hash             text,
    userop_count     int,
    userop_mev_count int,
    bundler_profit   numeric
);

SELECT create_hypertable('aa_block_info', 'time');
CREATE INDEX aa_block_info_hash_index ON aa_block_info USING HASH (hash);
CREATE INDEX aa_block_info_block_num_index ON aa_block_info (number);
CREATE INDEX aa_block_info_create_time_index ON aa_block_info (create_time);
CREATE UNIQUE INDEX aa_block_info_time_hash_index ON aa_block_info (time, number);


----------------------------------------------------------


create table transaction_decode
(
    time                     timestamp with time zone not null,
    create_time              timestamp with time zone,
    hash                     text,
    block_hash               text,
    block_number             numeric,
    nonce                    numeric,
    transaction_index        numeric,
    from_addr                text,
    to_addr                  text,
    value                    numeric,
    gas_price                numeric,
    gas                      numeric,
    input                    text,
    r                        text,
    s                        text,
    v                        int8,
    chain_id                 int8,
    type                     text,
    max_fee_per_gas          numeric,
    max_priority_fee_per_gas numeric,
    access_list              jsonb,
    method                   text
);

SELECT create_hypertable('transaction_decode', 'time');
CREATE INDEX transaction_decode_hash_index ON transaction_decode USING HASH (hash);
CREATE INDEX transaction_decode_block_num_index ON transaction_decode (block_number);
CREATE INDEX transaction_decode_create_time_index ON transaction_decode (create_time);
CREATE UNIQUE INDEX transaction_decode_time_hash_index ON transaction_decode (time, hash);


----------------------------------------------------------


create table aa_transaction_info
(
    time           timestamp with time zone not null,
    create_time    timestamp with time zone,
    hash           text,
    block_hash     text,
    block_number   numeric,
    userop_count   int,
    is_mev         boolean,
    bundler_profit numeric
);

SELECT create_hypertable('aa_transaction_info', 'time');
CREATE INDEX aa_transaction_info_hash_index ON aa_transaction_info USING HASH (hash);
CREATE INDEX aa_transaction_info_block_num_index ON aa_transaction_info (block_number);
CREATE INDEX aa_transaction_info_create_time_index ON aa_transaction_info (create_time);
CREATE UNIQUE INDEX aa_transaction_info_time_hash_index ON aa_transaction_info (time, hash);


----------------------------------------------------------

create table transaction_receipt_decode
(
    time                timestamptz,
    create_time         timestamp with time zone,
    transaction_hash    text,
    transaction_index   numeric,
    block_hash          text,
    block_number        numeric,
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
);

SELECT create_hypertable('transaction_receipt_decode', 'time');
CREATE INDEX transaction_receipt_decode_hash_index ON transaction_receipt_decode USING HASH (transaction_hash);
CREATE INDEX transaction_receipt_decode_block_num_index ON transaction_receipt_decode (block_number);
CREATE UNIQUE INDEX transaction_receipt_decode_transaction_hash_index ON transaction_receipt_decode (time, transaction_hash);

--------------------------------------------------------

create table aa_user_ops_calldata
(
    time          timestamptz,
    uuid          varchar(64),
    user_ops_hash text,
    tx_hash       text,
    block_number  int8,
    network       text,
    sender        text,
    target        text,
    tx_value      numeric,
    source        text,
    calldata      text,
    tx_time       int8,
    create_time   timestamp with time zone
);

SELECT create_hypertable('aa_user_ops_calldata', 'time');
CREATE INDEX aa_user_ops_calldata_tx_hash_index ON aa_user_ops_calldata USING HASH (tx_hash);
CREATE INDEX aa_user_ops_calldata_user_operation_hash_index ON aa_user_ops_calldata USING HASH (user_ops_hash);
CREATE INDEX aa_user_ops_calldata_block_num_index ON aa_user_ops_calldata (block_number);
CREATE UNIQUE INDEX aa_user_ops_calldata_time_uuid_index ON aa_user_ops_calldata (time, uuid);


------------------------------------------------------
drop table if exists public.aa_user_ops_info;
create table if not exists public.aa_user_ops_info
(
    time                     timestamp with time zone not null,
    user_operation_hash      text,
    tx_hash                  text,
    block_number             bigint,
    network                  text,
    sender                   text,
    target                   text,
    tx_value                 numeric,
    fee                      numeric,
    bundler                  text,
    entry_point              text,
    factory                  text,
    paymaster                text,
    paymaster_and_data       text,
    signature                text,
    calldata                 text,
    calldata_contract        text,
    nonce                    bigint,
    call_gas_limit           bigint,
    pre_verification_gas     bigint,
    verification_gas_limit   bigint,
    max_fee_per_gas          bigint,
    max_priority_fee_per_gas bigint,
    tx_time                  bigint,
    tx_time_format           text,
    init_code                text,
    status                   integer,
    source                   text,
    actual_gas_cost          bigint,
    actual_gas_used          bigint,
    create_time              timestamp with time zone
);

SELECT create_hypertable('aa_user_ops_info', 'time');
CREATE INDEX aa_user_ops_info_tx_hash_index ON aa_user_ops_info USING HASH (tx_hash);
CREATE INDEX aa_user_ops_info_user_operation_hash_index ON aa_user_ops_info USING HASH (user_operation_hash);
CREATE INDEX aa_user_ops_info_block_num_index ON aa_user_ops_info (block_number);
CREATE UNIQUE INDEX aa_user_ops_info_time_uuid_index ON aa_user_ops_info (time, tx_hash, user_operation_hash);

------------------------------------------------------


create table block_sync
(
    block_num   bigint not null
        constraint block_sync_pk
            primary key,
    scanned     boolean default false,
    create_time timestamp with time zone,
    update_time timestamp with time zone
);

create index block_sync_scanned_index
    on block_sync (scanned);


------------------------------------------------------

create table transaction_sync
(
    block_num   bigint not null
        constraint transaction_sync_pk
            primary key,
    scanned     boolean default false,
    create_time timestamp with time zone,
    update_time timestamp with time zone
);


drop trigger transaction_block_sync on block_sync;
drop function transaction_block_sync();
CREATE FUNCTION transaction_block_sync() RETURNS trigger AS
$$
DECLARE
    current_block bigint;
    current_scanned
                  boolean default false;
    current_create_time
                  timestamp with time zone;

BEGIN
    current_block
        = NEW.block_num;
    current_scanned
        = false;
    current_create_time
        = NEW.create_time;

    insert into transaction_sync (block_num, scanned, create_time)
    VALUES (current_block,
            current_scanned,
            current_create_time)
    on conflict (block_num) do nothing;
    Return null;
END;
$$
    LANGUAGE plpgsql;

CREATE TRIGGER transaction_block_sync
    AFTER INSERT
    ON block_sync
    FOR EACH ROW
EXECUTE FUNCTION transaction_block_sync();


drop trigger transaction_receipt_block_sync on block_sync;
drop function transaction_receipt_block_sync();
CREATE FUNCTION transaction_receipt_block_sync() RETURNS trigger AS
$$
DECLARE
    current_block bigint;
    current_scanned
                  boolean default false;
    current_create_time
                  timestamp with time zone;

BEGIN
    current_block
        = NEW.block_num;
    current_scanned
        = false;
    current_create_time
        = NEW.create_time;

    insert into transaction_receipt_block_sync (block_num, scanned, create_time)
    VALUES (current_block,
            current_scanned,
            current_create_time)
    on conflict (block_num) do nothing;
    Return null;
END;
$$
    LANGUAGE plpgsql;

CREATE TRIGGER transaction_receipt_block_sync
    AFTER INSERT
    ON block_sync
    FOR EACH ROW
EXECUTE FUNCTION transaction_receipt_block_sync();

----------------------------------------------------------


create table aa_block_sync
(
    block_num     bigint not null
        constraint aa_block_sync_pk
            primary key,
    block_scanned boolean,
    tx_scanned    boolean,
    txr_scanned   boolean,
    scanned       boolean, -- null is transaction unsync ,false is transaction synced -> aa syncing ,true is aa synced
    create_time   timestamp with time zone,
    update_time   timestamp with time zone, 
     scan_count integer default 0
);

create index aa_block_sync_scanned_index
    on aa_block_sync (scanned);


------------------------------------------------------

-- drop table  if exists  account;
create table account
(
    address text,
    type    text[],-- eoa,ca,sender,bunder,erc20,erc721
    tag     text,
    label   text[],
    update_time timestamp with time zone
);
-- drop table  if exists  account_data;
create table account_data
(
    address text,
    data    text,
    abi     text
);

--------------------------------------------------------

drop trigger aa_scan_sync on aa_block_sync;
drop trigger aa_scan_sync2 on aa_block_sync;
drop function aa_scan_sync;
CREATE FUNCTION aa_scan_sync() RETURNS trigger AS
$$
DECLARE
    current_block bigint;
    current_scanned
                  boolean default false;
    current_create_time
                  timestamp with time zone;
    current_update_time
                  timestamp with time zone;

BEGIN
    if NEW.scanned is not null then
        return null ;
    end if;
    current_block
        = NEW.block_num;
    current_scanned
        = NEW.scanned;
    current_create_time
        = NEW.create_time;
    current_update_time
        = NEW.update_time;

    if NEW.block_scanned and
       NEW.tx_scanned and
       NEW.txr_scanned then
        update aa_block_sync set scanned = false, update_time = current_timestamp where block_num = current_block;
    end if;
    Return null;
END;
$$
    LANGUAGE plpgsql;
CREATE TRIGGER aa_scan_sync
    AFTER UPDATE
    ON aa_block_sync
    FOR EACH ROW
EXECUTE FUNCTION aa_scan_sync();


CREATE FUNCTION aa_tx_sync() RETURNS trigger AS
$$
DECLARE
    current_block bigint;
    current_scanned
                  boolean default false;
    current_create_time
                  timestamp with time zone;
    current_update_time
                  timestamp with time zone;

BEGIN
    current_block
        = NEW.block_num;
    current_scanned
        = NEW.scanned;
    current_create_time
        = NEW.create_time;
    current_update_time
        = NEW.update_time;

    insert into aa_block_sync (block_num, tx_scanned, create_time)
    VALUES (current_block,
            current_scanned,
            current_create_time)
    on conflict (block_num) do update set tx_scanned = current_scanned, update_time = current_update_time;
    Return null;
END;
$$
    LANGUAGE plpgsql;

CREATE TRIGGER aa_tx_sync
    AFTER UPDATE
    ON transaction_sync
    FOR EACH ROW
EXECUTE FUNCTION aa_tx_sync();

CREATE FUNCTION aa_txr_sync() RETURNS trigger AS
$$
DECLARE
    current_block bigint;
    current_scanned
                  boolean default false;
    current_create_time
                  timestamp with time zone;
    current_update_time
                  timestamp with time zone;

BEGIN
    current_block
        = NEW.block_num;
    current_scanned
        = NEW.scanned;
    current_create_time
        = NEW.create_time;
    current_update_time
        = NEW.update_time;

    insert into aa_block_sync (block_num, txr_scanned, create_time)
    VALUES (current_block,
            current_scanned,
            current_create_time)
    on conflict (block_num) do update set txr_scanned = current_scanned, update_time = current_update_time;
    Return null;
END;
$$
    LANGUAGE plpgsql;

CREATE TRIGGER aa_txr_sync
    AFTER UPDATE
    ON transaction_receipt_block_sync
    FOR EACH ROW
EXECUTE FUNCTION aa_txr_sync();

drop trigger aa_scan_sync on aa_block_sync;
drop function aa_scan_sync;
CREATE FUNCTION aa_scan_sync() RETURNS trigger AS
$$
DECLARE
    current_block bigint;
    current_scanned
                  boolean default false;
    current_create_time
                  timestamp with time zone;
    current_update_time
                  timestamp with time zone;

BEGIN
    current_block
        = NEW.block_num;
    current_scanned
        = NEW.scanned;
    current_create_time
        = NEW.create_time;
    current_update_time
        = NEW.update_time;

    if NEW.block_scanned and
       NEW.tx_scanned and
       NEW.txr_scanned then
        update aa_block_sync set scanned = false, update_time = current_timestamp where block_num = current_block;
    end if;
    Return null;
END;
$$
    LANGUAGE plpgsql;
drop TRIGGER aa_scan_sync on aa_block_sync;
CREATE TRIGGER aa_scan_sync
    AFTER UPDATE
    ON aa_block_sync
    FOR EACH ROW
    WHEN (pg_trigger_depth() < 1)
EXECUTE FUNCTION aa_scan_sync();
