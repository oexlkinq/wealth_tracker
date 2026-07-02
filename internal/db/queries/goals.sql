-- name: CreateGoal :exec
insert into goals (
    amount, comment, index
) values (
    $1, $2, $3
);

-- name: ListGoals :many
select *
from goals;

-- name: DeleteGoal :exec
delete from goals
where id = $1;

-- name: DeleteGoalTxnsSince :exec
delete from txns
where
    goal_id is not null
    and ts >= $1;

-- name: ListGoalsForCalc :many
select goals.*
from goals
left join txns on txns.goal_id = goals.id
where txns.goal_id is null
order by goals.index;
