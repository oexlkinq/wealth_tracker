-- name: CreateTarget :exec
insert into targets (
    amount, desc, "order"
) values (
    ?, ?, ?
);

-- name: ListTargets :many
select *
from targets;

-- name: UpdateTractIDOfTarget :exec
update targets
set tract_id = ?
where id = @target_id;

-- name: DeleteTarget :exec
delete from targets
where id = ?;
