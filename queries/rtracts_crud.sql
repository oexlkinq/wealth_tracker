-- name: CreateRTract :exec
insert into rtracts (
    rrule, desc, amount, reqs_ack
) values (
    ?, ?, ?, ?
);

-- name: ListRTracts :many
select *
from rtracts;
