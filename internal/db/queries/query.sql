-- name: ListTractsSince :many
select t.*
from tracts t
left join rtracts_to_tracts rtt on rtt.tract_id = t.id
where t.date >= @since and rtt.rtract_id = @rtract_id
order by date asc;

-- name: ListUnreachedTargets :many
select t.*
from targets t
where tract_id is null;

-- name: GetLatestBalanceRecord :one
select br.*
from balance_records br
order by br.date desc, br.amount asc
limit 1;

-- name: CreateRTractToTract :exec
insert into rtracts_to_tracts (
    rtract_id, tract_id
) values (
    ?, ?
);

-- name: ListTargetsForCalc :many
select t.*
from targets t
order by t."order" asc;
