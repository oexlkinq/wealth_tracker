package models

import (
	"context"

	"github.com/oexlkinq/wealth_tracker/internal/db/db_api"
	"github.com/teambition/rrule-go"
)

type RawRTract struct {
	Rrule    string
	Desc     string
	Amount   float64
	Reqs_ack bool
}

func (v *RawRTract) Insert(ctx context.Context, qtx *db_api.Queries) error {
	_, err := rrule.StrToRRule(v.Rrule)
	if err != nil {
		return err
	}

	err = qtx.CreateRTract(ctx, db_api.CreateRTractParams{
		Rrule:   v.Rrule,
		Desc:    v.Desc,
		Amount:  v.Amount,
		ReqsAck: v.Reqs_ack,
	})
	if err != nil {
		return err
	}

	return nil
}
