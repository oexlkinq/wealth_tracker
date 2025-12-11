package calc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"iter"
	"time"

	"github.com/oexlkinq/wealth_tracker/internal/db/db_api"
	"github.com/oexlkinq/wealth_tracker/internal/itergroup"
	"github.com/oexlkinq/wealth_tracker/internal/itergroup/tractsiter/models"
)

const maxTractsCount = 1000

var ErrorTooManyTracts = errors.New("too many tracts")

type Calc struct {
	ctx     context.Context
	qtx     *db_api.Queries
	tris    []*TargetReachInfo
	ig      itergroup.TractsIterGroup
	balance db_api.BalanceRecord
}

func New(ctx context.Context, qtx *db_api.Queries) (*Calc, error) {
	// сбор данных для расчёта
	balanceRecord, err := qtx.GetLatestBalanceRecord(ctx)
	if err != nil {
		return nil, err
	}

	ig, err := itergroup.New(ctx, qtx, balanceRecord.Date)
	if err != nil {
		return nil, err
	}

	targets, err := qtx.ListTargetsForCalc(ctx)
	if err != nil {
		return nil, err
	}

	tris := make([]*TargetReachInfo, len(targets))
	for i := range targets {
		tris[i] = &TargetReachInfo{Target: &targets[i]}
	}

	return &Calc{
		ctx:     ctx,
		qtx:     qtx,
		tris:    tris,
		ig:      ig,
		balance: balanceRecord,
	}, nil
}

func (v *Calc) CalcTargetsReachInfo() ([]*TargetReachInfo, error) {
	next, stop := iter.Pull(v.ig.All())
	defer stop()

	tractsCount := 0
	for _, targetReachInfo := range v.tris {
		for {
			// если бюджета уже достаточно
			if v.balance.Amount >= targetReachInfo.Target.Amount {
				v.balance.Amount -= targetReachInfo.Target.Amount

				targetReachInfo.ReachedAmount = targetReachInfo.Target.Amount
				targetReachInfo.ReachDate = v.balance.Date
				targetReachInfo.Reached = true

				err := v.saveTargetsTract(targetReachInfo.Target.Amount, targetReachInfo.Target.ID)
				if err != nil {
					return nil, fmt.Errorf("save targets tract: %w", err)
				}

				// TODO: убрать отладочный вывод
				fmt.Println("break coz reached", targetReachInfo, v.balance.Amount)
				break
			}

			if tractsCount > maxTractsCount {
				return nil, ErrorTooManyTracts
			}

			calcTract, ok := next()
			// если кончились транзакции
			if !ok {
				targetReachInfo.ReachedAmount = v.balance.Amount
				targetReachInfo.ReachDate = v.balance.Date
				// targetReachInfo.reached = false

				err := v.saveTargetsTract(targetReachInfo.Target.Amount, targetReachInfo.Target.ID)
				if err != nil {
					return nil, fmt.Errorf("save targets tract: %w", err)
				}

				// TODO: убрать отладочный вывод
				fmt.Println("end coz no next", targetReachInfo, v.balance.Amount)
				return v.tris, nil
			}
			tractsCount++

			v.balance.Amount += calcTract.Amount
			v.balance.Date = calcTract.Date

			err := v.saveCalcTract(calcTract)
			if err != nil {
				return nil, fmt.Errorf("save calcTract: %w", err)
			}
		}
	}

	return v.tris, nil
}

func (v *Calc) saveTargetsTract(amount float64, targetID int64) error {
	id, err := v.saveTract("target", amount)
	if err != nil {
		return err
	}

	err = v.qtx.UpdateTractIDOfTarget(v.ctx, db_api.UpdateTractIDOfTargetParams{
		TargetID: targetID,
		TractID:  sql.NullInt64{Int64: id, Valid: true},
	})
	if err != nil {
		return err
	}

	return nil
}

func (v *Calc) saveCalcTract(calcTract *models.CalcTract) error {
	if !calcTract.Generated {
		// TODO: т.к. сейчас реализован только минимальный функционал, в бд не должно существовать транзакций, т.к. они все удаляются перед запуском
		panic(fmt.Errorf("not implemented"))
	}

	// создание транзакции и отметки баланса
	id, err := v.saveTract("rtract", calcTract.Amount)
	if err != nil {
		return fmt.Errorf("save tract: %w", err)
	}

	// создание связи транзакции с ртрактом
	err = v.qtx.CreateRTractToTract(v.ctx, db_api.CreateRTractToTractParams{
		RtractID: calcTract.RTractID,
		TractID:  id,
	})
	if err != nil {
		return err
	}

	return nil
}

func (v *Calc) saveTract(tractType string, amount float64) (int64, error) {
	id, err := v.qtx.CreateTract(v.ctx, db_api.CreateTractParams{
		Type:   tractType,
		Date:   v.balance.Date,
		Amount: amount,
		Acked:  false,
	})
	if err != nil {
		return 0, err
	}

	err = v.qtx.CreateBalanceRecord(v.ctx, db_api.CreateBalanceRecordParams{
		Amount:      v.balance.Amount,
		Date:        v.balance.Date,
		OriginTract: sql.NullInt64{Int64: id, Valid: true},
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}

type TargetReachInfo struct {
	Target        *db_api.ListTargetsForCalcRow
	ReachDate     time.Time
	ReachedAmount float64
	Reached       bool
}

func (v *TargetReachInfo) String() string {
	return fmt.Sprintf("{desc: %s, order: %d, reached: %.1f/%.1f, reachDate: %s}", v.Target.Desc, v.Target.Order, v.ReachedAmount, v.Target.Amount, v.ReachDate.Format(time.DateOnly))
}
