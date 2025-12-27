-- +goose Up
-- +goose StatementBegin
create table if not exists tracts (
    id integer primary key,
    type text not null,
    date date not null,
    amount real not null,
    acked boolean not null
);

create table if not exists targets (
    id integer primary key,
    amount real not null,
    desc text not null,
    "order" integer unique not null,
    tract_id integer references tracts(id) on delete set null on update cascade
);

create table if not exists rtracts (
    id integer primary key,
    rrule text not null,
    desc text not null,
    amount real not null,
    reqs_ack boolean not null
);

create table if not exists rtracts_to_tracts (
    rtract_id integer not null references rtracts(id) on delete cascade on update cascade,
    tract_id integer not null references tracts(id) on delete cascade on update cascade
);

create table if not exists balance_records (
    id integer primary key,
    amount real not null,
    date date not null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table tracts;
drop table targets;
drop table rtracts;
drop table rtracts_to_tracts;
drop table balance_records;
-- +goose StatementEnd
