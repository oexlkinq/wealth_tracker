-- name: ListTractsSince :many
select t.amount, t.date
from tracts t
left join rtracts_to_tracts rtt on rtt.tract_id = t.id
where t.date >= @since and rtt.rtract_id = @rtract_id
order by date asc;

-- name: ListRTracts :many
select rt.*
from rtracts rt;

-- name: ListUnreachedTargets :many
select t.*
from targets t
where tract_id is null;

-- name: GetLatestBalanceRecord :one
select br.*
from balance_records br
order by br.date desc, br.amount asc
limit 1;

-- name: CreateTract :one
insert into tracts (
    type, date, amount, acked
) values (
    ?, ?, ?, ?
)
returning id;

-- name: CreateRTractToTract :exec
insert into rtracts_to_tracts (
    rtract_id, tract_id
) values (
    ?, ?
);

-- name: CreateBalanceRecord :exec
insert into balance_records (
    amount, date, origin_tract
) values (
    ?, ?, ?
);

-- name: DeleteBalanceRecordsSince :exec
delete from balance_records
where date >= @since;

-- name: DeleteTractsSince :exec
delete from tracts
where date >= @since;

-- name: DeleteBalanceRecord :exec
delete from balance_records
where id = ?;

-- name: DeleteTract :exec
delete from tracts
where id = ?;