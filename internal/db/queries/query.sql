-- name: ListTractsSince :many
select t.*
from tracts t
where t.date >= @since and t.rtract_id = @rtract_id
order by date asc;

-- name: GetLatestBalanceRecord :one
select br.*
from balance_records br
order by br.date desc, br.amount asc
limit 1;

-- name: ListTargetsForCalc :many
select t.*
from targets t
order by t."order" asc;

-- name: ListRtractsWithLastTracts :many
with ranked_rtt as (
    select rank() over (partition by rt.id order by t.date desc, t.id desc) as tract_rank, rt.*, t.date
    from rtracts rt
    left join tracts t on rt.id = t.rtract_id
)
select *
from ranked_rtt
where tract_rank = 1;
