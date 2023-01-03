drop table if exists pictures;

create table pictures (
    id uuid primary key not null,
    fingerprint varchar
);