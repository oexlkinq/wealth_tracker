-- name: CreateTxn :one
insert into txns (
    amount, comment, ts, rtxn_id, goal_id
) values (
    $1, $2, $3, $4, $5
)
returning id;

-- name: ListTxnsSince :many
select *
from txns
where ts >= @since and rtxn_id = @rtxn_id
order by ts asc;
