create table public.function_signature (
                                           tableoid oid not null,
                                           cmax cid not null,
                                           xmax xid not null,
                                           cmin cid not null,
                                           xmin xid not null,
                                           ctid tid not null,
                                           signature character varying primary key not null,
                                           name character varying,
                                           text text,
                                           bytes bytea,
                                           create_time timestamp with time zone
);

create table public.monitor (
                                tableoid oid not null,
                                cmax cid not null,
                                xmax xid not null,
                                cmin cid not null,
                                xmin xid not null,
                                ctid tid not null,
                                id bigint primary key not null default nextval('monitor_id_seq'::regclass),
                                user_address character varying,
                                monitor_address character varying,
                                monitor_address_type character varying,
                                create_time timestamp without time zone not null default CURRENT_TIMESTAMP
);
create index monitor_monitor_address_idx on monitor using hash (monitor_address);
create index monitor_user_address_idx on monitor using hash (user_address);

create table public.network (
                                tableoid oid not null,
                                cmax cid not null,
                                xmax xid not null,
                                cmin cid not null,
                                xmin xid not null,
                                ctid tid not null,
                                name character varying not null,
                                network character varying primary key not null,
                                http_rpc character varying not null,
                                is_testnet boolean not null,
                                create_time timestamp with time zone not null,
                                update_time timestamp with time zone,
                                delete_time timestamp with time zone,
                                chain_id bigint,
                                chain_name character varying,
                                scan character varying,
                                scan_tx character varying,
                                scan_block character varying,
                                scan_address character varying,
                                scan_name character varying,
                                db_config jsonb,
                                chain_icon character varying,
                                coin_icon character varying,
                                native_symbol character varying
);
create unique index networks_network_key on network using btree (network);

