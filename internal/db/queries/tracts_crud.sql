-- name: CreateTract :one
insert into tracts (
    type, date, amount, acked
) values (
    ?, ?, ?, ?
)
returning id;
