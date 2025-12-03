-- name: CreateTarget :exec
insert into targets (
    amount, desc, "order"
) values (
    ?, ?, ?
);

-- name: ListTargets :many
select *
from targets;

-- name: DeleteTarget :exec
delete from targets
where id = ?;
