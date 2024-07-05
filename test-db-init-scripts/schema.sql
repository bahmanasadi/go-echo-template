-- migrations/schema.sql

create table public.persons
(
    id          bigserial primary key,
    external_id text unique not null,
    email       text unique,
    password    bytea,
    created_at  timestamp,
    updated_at  timestamp
);

alter table public.persons
    owner to postgres;
