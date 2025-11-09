-- name: ListGTractsAfter :many
select t.*
from rtract_to_tract rtt
join tracts t on rtt.tract = t.id
where rtt.rtract_id = sqlc.arg(target_rtract)
order by date asc;