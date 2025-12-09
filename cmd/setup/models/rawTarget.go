package models

import (
	"context"

	"github.com/oexlkinq/wealth_tracker/internal/db/db_api"
)

type RawTarget struct {
	Amount float64
	Desc   string
	Order  int64
}

func (v *RawTarget) Insert(ctx context.Context, qtx *db_api.Queries) error {
	return qtx.CreateTarget(ctx, db_api.CreateTargetParams(*v))
}
