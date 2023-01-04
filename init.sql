drop table if exists pictures;

create table pictures (
    id uuid primary key not null,
    fingerprint varchar not null,
    created_at  time without time zone default (now() at time zone 'utc') not null
);