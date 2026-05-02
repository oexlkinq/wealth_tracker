-- +goose Up
-- +goose StatementBegin
create table goals (
    id int generated always as identity primary key,
    amount float not null,
    comment text,
    hidden bool not null,
    index int unique not null
);

create table rtxns (
    id int generated always as identity primary key,
    amount float not null,
    comment text,
    rrule text not null
);

create table txns (
    id int generated always as identity primary key,
    amount float not null,
    comment text,
    ts timestamptz not null,
    rtxn_id int references rtxns(id) on delete cascade on update cascade,
    goal_id int references goals(id) on delete cascade on update cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table txns;
drop table goals;
drop table rtxns;
-- +goose StatementEnd
