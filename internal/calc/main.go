package calc

import (
	"fmt"
	"maps"
	"time"

	"github.com/oexlkinq/wealth_tracker/internal/models"
	"github.com/oexlkinq/wealth_tracker/internal/rtractsgroup"
)

const maxTractsCount = 1000

func CalcTargetsReachInfo(rtracts []*models.RTract, balanceRecord models.BalanceRecord, targets []*models.Target) (targetReachInfoStructs []*TargetReachInfo, tooManyTracts bool) {
	rtgroup := rtractsgroup.New(rtracts, balanceRecord.Date)
	targetReachInfoStructs = prepareTargets(targets)

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
				fmt.Println(targetReachInfo, balanceRecord.Amount)
				break
			}

			if tractsCount > maxTractsCount {
				tooManyTracts = true
				return
			}

			rtract, ok := rtgroup.Next()
			// если кончились транзакции
			if !ok {
				targetReachInfo.ReachedAmount = balanceRecord.Amount
				targetReachInfo.ReachDate = balanceRecord.Date
				// targetReachInfo.reached = false

				// TODO: убрать отладочный вывод
				fmt.Println(targetReachInfo, balanceRecord.Amount)
				return
			}
			tractsCount++

			balanceRecord.Amount += rtract.Amount
			balanceRecord.Date = rtract.Date
		}
	}

	return
}

type TargetReachInfo struct {
	*models.Target
	ReachDate     time.Time
	ReachedAmount float64
	Reached       bool
}

func (v *TargetReachInfo) String() string {
	return fmt.Sprintf("{desc: %s, order: %d, reached: %.1f/%.1f, reachDate: %s}", v.Desc, v.Order, v.ReachedAmount, v.Amount, v.ReachDate)
}

// выбирает для каждой очереди цель с максимальной суммой и конвертит в нужный тип
func prepareTargets(targets []*models.Target) []*TargetReachInfo {
	order_to_target := make(map[int]*models.Target, len(targets))
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
