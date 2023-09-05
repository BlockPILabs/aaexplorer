create table public.network
(
    network     varchar(127)             not null
        constraint network_pk
            primary key,
    name        varchar(255)             not null,
    logo        varchar(255)             not null,
    http_rpc    varchar(255)             not null,
    is_testnet  boolean                  not null,
    create_time timestamp with time zone not null,
    update_time timestamp with time zone,
    delete_time timestamp with time zone,
    chain_id    bigint
);

alter table public.network
    owner to postgres;

create unique index networks_network_key
    on public.network (network);

