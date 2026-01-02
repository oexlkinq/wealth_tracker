-- +goose Up
-- +goose StatementBegin
create table if not exists tracts (
    id integer primary key,
    date date not null,
    amount real not null,
    acked boolean not null,
    target_id integer references targets(id) on delete cascade on update cascade,
    rtract_id integer references rtracts(id) on delete cascade on update cascade,
    check (target_id is not null or rtract_id is not null)
);

create table if not exists targets (
    id integer primary key,
    amount real not null,
    desc text not null,
    "order" integer unique not null
);

create table if not exists rtracts (
    id integer primary key,
    rrule text not null,
    desc text not null,
    amount real not null,
    reqs_ack boolean not null
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
drop table balance_records;
-- +goose StatementEnd
