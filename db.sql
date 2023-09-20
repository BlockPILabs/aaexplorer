create table paymaster_info
(
    id                bigserial
        primary key,
    paymaster         varchar(255),
    network           varchar(255),
    user_ops_num      bigint          default 0,
    gas_sponsored     numeric(50, 20) default 0,
    user_ops_num_d1   bigint          default 0,
    gas_sponsored_d1  numeric(50, 20) default 0,
    user_ops_num_d7   bigint          default 0,
    gas_sponsored_d7  numeric(50, 20) default 0,
    user_ops_num_d30  bigint          default 0,
    gas_sponsored_d30 numeric(50, 20) default 0,
    create_time       timestamp(3)    default CURRENT_TIMESTAMP,
    update_time       timestamp(3)
);


create table paymaster_statis_day
(
    id            bigint          default nextval('paymaster_statis_hour_id_seq'::regclass) not null
        constraint paymaster_statis_hour_copy1_pkey
            primary key,
    paymaster     varchar(255),
    network       varchar(255),
    user_ops_num  bigint          default 0,
    gas_sponsored numeric(50, 20) default 0,
    statis_time   timestamp(6),
    create_time   timestamp(3)    default CURRENT_TIMESTAMP                                 not null
);

create table paymaster_statis_hour
(
    id            bigserial
        primary key,
    paymaster     varchar(255),
    network       varchar(255),
    user_ops_num  bigint          default 0,
    gas_sponsored numeric(50, 20) default 0,
    statis_time   timestamp,
    create_time   timestamp(3)    default CURRENT_TIMESTAMP not null
);

create table hot_aa_token_statistic
(
    id             bigserial
        primary key,
    token_symbol   varchar(255)    not null,
    network        varchar(255)    not null,
    statistic_type varchar(255)    not null,
    volume         numeric(32, 18) not null,
    create_time    timestamp(3) default CURRENT_TIMESTAMP
);

create table factory_statis_hour
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

create table factory_statis_day
(
    id                 bigint       default nextval('factory_statis_hour_id_seq'::regclass) not null
        constraint factory_statis_hour_copy1_pkey
            primary key,
    factory            varchar(255),
    network            varchar(255),
    account_num        bigint       default 0,
    account_deploy_num bigint       default 0,
    statis_time        timestamp(6),
    create_time        timestamp(3) default CURRENT_TIMESTAMP                               not null
);

create table bundler_statis_hour
(
    id            bigserial
        primary key,
    bundler       varchar(255),
    network       varchar(255),
    user_ops_num  bigint          default 0,
    bundles_num   bigint          default 0,
    gas_collected numeric(50, 20) default 0,
    statis_time   timestamp,
    create_time   timestamp(3)    default CURRENT_TIMESTAMP not null
);

create table bundler_statis_day
(
    id            bigint          default nextval('bundler_statis_hour_id_seq'::regclass) not null
        constraint bundler_statis_hour_copy1_pkey
            primary key,
    bundler       varchar(255),
    network       varchar(255),
    user_ops_num  bigint          default 0,
    bundles_num   bigint          default 0,
    gas_collected numeric(50, 20) default 0,
    statis_time   timestamp(6),
    create_time   timestamp(3)    default CURRENT_TIMESTAMP                               not null
);

create table bundler_info
(
    id                bigserial
        primary key,
    bundler           varchar(255),
    network           varchar(255),
    user_ops_num      bigint          default 0,
    bundles_num       bigint          default 0,
    gas_collected     numeric(50, 20) default 0,
    user_ops_num_d1   bigint          default 0,
    bundles_num_d1    bigint          default 0,
    gas_collected_d1  numeric(50, 20) default 0,
    user_ops_num_d7   bigint          default 0,
    bundles_num_d7    bigint          default 0,
    gas_collected_d7  numeric(50, 20) default 0,
    user_ops_num_d30  bigint          default 0,
    bundles_num_d30   bigint          default 0,
    gas_collected_d30 numeric(50, 20) default 0,
    create_time       timestamp(3)    default CURRENT_TIMESTAMP not null,
    update_time       timestamp(3)
);

create table asset_change_trace
(
    id               bigint generated by default as identity
        primary key,
    tx_hash          text,
    block_number     bigint,
    address          text,
    address_type     integer,
    create_time      timestamp(3) default CURRENT_TIMESTAMP not null,
    sync_flag        smallint,
    last_change_time date,
    network          varchar(255)
);

create table block_scan_record
(
    id                bigserial
        primary key,
    network           varchar(255),
    last_block_number bigint,
    last_scan_time    timestamp,
    create_time       timestamp(3) default CURRENT_TIMESTAMP not null,
    update_time       timestamp
);

create table user_asset_info
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

create table token_price_info
(
    id               bigserial
        primary key,
    contract_address varchar(255),
    symbol           varchar(255),
    token_price      numeric(50, 30),
    last_time        bigint,
    create_time      timestamp(3) default CURRENT_TIMESTAMP not null,
    update_time      timestamp,
    network          varchar(255)
);

create table user_op_type_statistic
(
    id             bigserial
        primary key,
    user_op_type   varchar(255),
    user_op_sign   varchar(255),
    network        varchar(255),
    statistic_type varchar(255),
    op_num         bigint,
    create_time    timestamp(3) default CURRENT_TIMESTAMP not null
);

create table daily_statistic_day
(
    id             bigserial
        primary key,
    network        varchar(255),
    statistic_time date,
    tx_num         bigint,
    user_ops_num   bigint,
    gas_fee        numeric(50, 20),
    active_wallet  bigint,
    create_time    timestamp(3) default CURRENT_TIMESTAMP not null
);


create table aa_contract_interact
(
    id               bigint                              not null
        primary key,
    contract_address varchar(255),
    network          varchar(255),
    statistic_type   varchar(255),
    interact_num     bigint,
    create_time      timestamp default CURRENT_TIMESTAMP not null
);



create table public.network
(
    name         varchar(255)             not null,
    network      varchar(127)             not null
        constraint network_pk
            primary key,
    http_rpc     varchar(255)             not null,
    is_testnet   boolean                  not null,
    create_time  timestamp with time zone not null,
    update_time  timestamp with time zone,
    delete_time  timestamp with time zone,
    chain_id     bigint,
    chain_name   varchar(127),
    scan         varchar(127),
    scan_tx      varchar(255),
    scan_block   varchar(255),
    scan_address varchar(255),
    scan_name    varchar(127),
    db_config    jsonb
);

alter table public.network
    owner to postgres;

create unique index networks_network_key
    on public.network (network);

create table public.function_signature
(
    signature   varchar(16) not null
        constraint function_signature_pk
            primary key,
    name        varchar(255),
    text        text,
    bytes       bytea,
    create_time timestamp with time zone
);

alter table public.function_signature
    owner to postgres;


