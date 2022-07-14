CREATE TABLE recordings (
    id bigserial primary key not null,
    original varchar not null unique,
    shortened varchar not null unique,
    created timestamp default current_timestamp
)