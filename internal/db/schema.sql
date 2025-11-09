create table if not exists tracts (
    id integer primary key,
    type text not null,
    date integer not null,
    amount real not null,
    ack boolean not null,
    reqs_ack boolean not null
);

create table if not exists targets (
    id integer primary key,
    amount real not null,
    desc text not null,
    order integer not null,
    tract_id integer not null references tracts(id)
);

create table if not exists rtracts (
    id integer primary key,
    rrule text not null,
    desc text not null,
    amount real not null
);

create table if not exists rtract_to_tract (
    rtract_id integer not null references rtracts(id),
    tract_id integer not null references tracts(id)
);

create table if not exists balances (
    id integer primary key,
    amount real not null,
    date integer not null,
    origin_tract integer references tracts(id)
);
