targets
- amount
- desc
- order
- @trans

rtracts
- rrule
- desc
- amount

rtract_to_tract
- @rtract
- @tract

tracts
- type (rtract, target) - чтобы знать в какой таблице искать указатель на себя
- date
- amount
- ack - acknowledged - принято во внимание? - подтверждено?
- reqs_ack - треубет подтверждения?

balances
- amount
- date
- origin_tract - транзакция-причина; у созданных пользователем можно задавать равным null

# запросы

## получить все существующие инстансы ртракта начиная с даты
```sql
-- выбрать все поля tract'а кроме типа
select id, date, amount, ack, reqs_ack
-- сопоставление rtracts и связанных транзакций
from rtract_to_tract rtt
join tracts t on rtt.tract = t.id
-- фильтрация по целевой rtract
where rtt.rtract = $target_rtract
order by date asc
```
