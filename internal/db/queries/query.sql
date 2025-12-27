-- name: ListTractsSince :many
select t.*
from tracts t
left join rtracts_to_tracts rtt on rtt.tract_id = t.id
where t.date >= @since and rtt.rtract_id = @rtract_id
order by date asc;

-- name: ListUnreachedTargets :many
select t.*
from targets t
where tract_id is null;

-- name: GetLatestBalanceRecord :one
select br.*
from balance_records br
order by br.date desc, br.amount asc
limit 1;

-- name: CreateRTractToTract :exec
insert into rtracts_to_tracts (
    rtract_id, tract_id
) values (
    ?, ?
);

-- name: ListTargetsForCalc :many
select t.*
from targets t
order by t."order" asc;

-- name: ListRtractsWithLastTracts :many
with ranked_rtt as (
    select rank() over (partition by rt.id order by t.date desc, t.id desc) as tract_rank, rt.*, t.date
    from rtracts rt
    left join rtracts_to_tracts rtt on rtt.rtract_id = rt.id
    left join tracts t on rtt.tract_id = t.id
)
select *
from ranked_rtt
where tract_rank = 1;

-- name: GetReachingTargetDate :one
with
	balance as (select amount, date from balance_records order by date desc limit 1),
	csum as (
		select id, amount, date, (select amount from balance) + sum(amount) OVER (ORDER BY date, id ROWS UNBOUNDED PRECEDING) as csum
		from tracts
		where date > (select date from balance)
	),
	csum_with_prev as (
		select id, amount, date, t.csum, lag(csum, 1) over (order by date, id) as prev_csum
		from csum t
	)
select "date"
from csum_with_prev
where csum >= @target and (prev_csum < @target or prev_csum is null)
order by date desc, id desc
limit 1;
