with
	csums as (
		select
            id,
            amount,
            ts,
            sum(amount) OVER (txns_seq) as csum,
            sum(amount) OVER (txns_seq ROWS BETWEEN UNBOUNDED PRECEDING AND 1 PRECEDING) as prev_csum
		from txns
        window txns_seq as (order by ts, id)
	)
select ts
from csums
where prev_csum < $1 and $1 <= csum
order by ts desc, id desc
limit 1;