-- name: CreateBalanceRecord :exec
insert into balance_records (
    amount, date, origin_tract
) values (
    ?, ?, ?
);

-- name: ListBalanceRecords :many
select *
from balance_records;

-- name: DeleteBalanceRecord :exec
delete from balance_records
where id = ?;
