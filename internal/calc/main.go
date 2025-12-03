package calc

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"maps"
	"time"

	"github.com/oexlkinq/wealth_tracker/internal/db_api"
	"github.com/oexlkinq/wealth_tracker/internal/itergroup"
)

const maxTractsCount = 1000

var TooManyTractsError = errors.New("too many tracts")

func CalcTargetsReachInfo(ctx context.Context, queries *db_api.Queries, balanceRecord db_api.BalanceRecord, targets []db_api.Target) ([]*TargetReachInfo, error) {
	ig, err := itergroup.New(ctx, queries, balanceRecord.Date)
	if err != nil {
		return nil, err
	}
	next, stop := iter.Pull(ig.All())
	defer stop()

	targetReachInfoStructs := prepareTargets(targets)

	tractsCount := 0
	for _, targetReachInfo := range targetReachInfoStructs {
		for {
			// TODO: убрать отладочный вывод
			fmt.Printf("%3d %s %10.1f\n", tractsCount, balanceRecord.Date, balanceRecord.Amount)

			// если бюджета уже достаточно
			if balanceRecord.Amount >= targetReachInfo.Amount {
				balanceRecord.Amount -= targetReachInfo.Amount

				targetReachInfo.ReachedAmount = targetReachInfo.Amount
				targetReachInfo.ReachDate = balanceRecord.Date
				targetReachInfo.Reached = true

				// TODO: убрать отладочный вывод
				fmt.Println("break coz reached", targetReachInfo, balanceRecord.Amount)
				break
			}

			if tractsCount > maxTractsCount {
				return nil, TooManyTractsError
			}

			tract, ok := next()
			// если кончились транзакции
			if !ok {
				targetReachInfo.ReachedAmount = balanceRecord.Amount
				targetReachInfo.ReachDate = balanceRecord.Date
				// targetReachInfo.reached = false

				// TODO: убрать отладочный вывод
				fmt.Println("end coz no next", targetReachInfo, balanceRecord.Amount)

				return targetReachInfoStructs, nil
			}
			tractsCount++

			balanceRecord.Amount += tract.Amount
			balanceRecord.Date = tract.Date
		}
	}

	return targetReachInfoStructs, nil
}

type TargetReachInfo struct {
	db_api.Target
	ReachDate     time.Time
	ReachedAmount float64
	Reached       bool
}

func (v *TargetReachInfo) String() string {
	return fmt.Sprintf("{desc: %s, order: %d, reached: %.1f/%.1f, reachDate: %s}", v.Desc, v.Order, v.ReachedAmount, v.Amount, v.ReachDate.Format(time.DateOnly))
}

// выбирает для каждой очереди цель с максимальной суммой и конвертит в нужный тип
func prepareTargets(targets []db_api.Target) []*TargetReachInfo {
	order_to_target := make(map[int64]db_api.Target, len(targets))
	for _, target := range targets {
		v, ok := order_to_target[target.Order]
		if !ok || target.Amount > v.Amount {
			order_to_target[target.Order] = target
		}
	}

	targetReachInfoStructs := make([]*TargetReachInfo, len(order_to_target))
	i := 0
	for target := range maps.Values(order_to_target) {
		targetReachInfoStructs[i] = &TargetReachInfo{Target: target}
		i++
	}

	return targetReachInfoStructs
}
