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

