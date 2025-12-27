-- name: CreateBalanceRecord :exec
insert into balance_records (
    amount, date
) values (
    ?, ?
);

-- name: ListBalanceRecords :many
select *
from balance_records;

-- name: DeleteBalanceRecord :exec
delete from balance_records
where id = ?;
