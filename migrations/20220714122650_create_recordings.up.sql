CREATE TABLE Recordings (
    id bigserial primary key,
    created date,
    original varchar not null unique,
    shortened varchar not null unique
)