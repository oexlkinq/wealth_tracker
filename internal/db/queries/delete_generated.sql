-- name: DeleteGeneratedTracts :exec
delete from tracts;

-- name: DeleteGeneratedBalanceRecords :exec
delete from balance_records
where origin_tract is not null;
