-- name: CreateRtxn :exec
insert into rtxns (
    amount, comment, rrule
) values (
    $1, $2, $3
);

-- name: ListRtxns :many
select *
from rtxns;

-- name: ListRtxnsEnds :many
with latest_txns as (
    select rtxn_id, max(ts) as ts
    from txns
    where rtxn_id is not null
    group by rtxn_id
)
select rt.*, t.ts
from rtxns rt
left join latest_txns t on rt.id = t.rtxn_id;
