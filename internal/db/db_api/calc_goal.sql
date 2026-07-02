with
      -- вычислить для каждой транзакции кумулятивную сумму (баланс) до и после транзакции
      csums as (
            select id,
                  amount,
                  ts,
                  sum(amount) OVER (txns_seq) as csum,
                  sum(amount) OVER (
                        txns_seq ROWS BETWEEN UNBOUNDED PRECEDING AND 1 PRECEDING
                  ) as prev_csum
            from txns
            window txns_seq as (order by ts, id)
      )
-- найти последний из доступных момент (последний, чтобы пропустить возможные уходы в минус), когда баланс стал больше целевой суммы
select ts
from csums
where (prev_csum < $1 or prev_csum is NULL)
      and $1 <= csum
order by ts desc,
      id desc
limit 1;