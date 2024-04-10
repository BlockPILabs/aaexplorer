--
--
/*
\connect postgres;
drop database block_db;
*/
CREATE DATABASE   block_db;
\connect block_db;
create schema if not exists  partman;
create extension if not exists pg_partman with schema partman;
\connect postgres;
-- drop extension pg_cron;
create extension if not exists pg_cron;
-- function schedule_in_database(job_name text, schedule text, command text, database text, username text default NULL::text, active boolean default true) returns bigint
SELECT cron.schedule_in_database('block_db_partman','@hourly', $$CALL partman.run_maintenance_proc()$$,'block_db','postgres');
\connect block_db;
