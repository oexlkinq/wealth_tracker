-- name: CreateTract :one
insert into tracts (
    date, amount, acked, target_id, rtract_id
) values (
    ?, ?, ?, ?, ?
)
returning id;
