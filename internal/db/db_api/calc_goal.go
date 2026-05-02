package db_api

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5/pgtype"
)

//go:embed calc_goal.sql
var query string

func (q *Queries) GetGoalReachTs(ctx context.Context, target float64) ([]pgtype.Timestamptz, error) {
	rows, err := q.db.Query(ctx, query, target)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []pgtype.Timestamptz
	for rows.Next() {
		var i pgtype.Timestamptz
		if err := rows.Scan(&i); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
