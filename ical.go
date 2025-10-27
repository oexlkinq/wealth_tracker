package main

import (
	"fmt"
	"maps"
	"time"

	"github.com/oexlkinq/wealth_tracker/internal/models"
	"github.com/oexlkinq/wealth_tracker/internal/rtractsgroup"
	"github.com/teambition/rrule-go"
)

func recalcAllTheShit() []*TargetReachInfo {
	rgroup := rtractsgroup.New([]models.Rtract{
		models.Rtract{
			Amount: 15000,
			Desc:   "зп",
			RRule: makeRRule(rrule.ROption{
				Dtstart:  time.Now(),
				Freq:     rrule.MONTHLY,
				Interval: 1,
			}),
		},
		models.Rtract{
			Amount: -10000,
			Desc:   "трата",
			RRule: makeRRule(rrule.ROption{
				Dtstart:  time.Now().Add(time.Hour * 24 * 14),
				Freq:     rrule.MONTHLY,
				Interval: 1,
			}),
		},
	}, time.Now())

	latestBudgetInfo := models.BudgetInfo{
		Amount: 10000,
		Date:   time.Now(),
	}

	targets := []models.Target{
		models.Target{
			Amount: 35000,
			Order:  0,
			Desc:   "на чёрный день",
		},
		models.Target{
			Amount: 15000,
			Order:  1,
			Desc:   "рандомная цель",
		},
	}

	budget := latestBudgetInfo
	targetReachInfoStructs := prepareTargets(targets)

	// TODO: если через год вычислений бюджет не превысил стартовую отметку то закончить рассчёт
	// TODO: если общее кол-во транзакций превысило 1к то закончить рассчёт
	for _, targetReachInfo := range targetReachInfoStructs {
		for {
			// если бюджета уже достаточно
			if budget.Amount >= targetReachInfo.Amount {
				budget.Amount -= targetReachInfo.Amount

				targetReachInfo.reachedAmount = targetReachInfo.Amount
				targetReachInfo.reachDate = budget.Date
				targetReachInfo.reached = true

				break
			}

			rtract, ok := rgroup.Next()
			if !ok {
				// TODO: если закончились транзакции, проверить что текущая цель правильно заполнена
				targetReachInfo.reachedAmount = budget.Amount
				targetReachInfo.reachDate = budget.Date
				// targetReachInfo.reached = false

				// TODO: заменить break на выход из обоих циклов. сейчас break просто сработает для всех целей когда закончатся транзакции
				break
			}

			budget.Amount += rtract.Amount
			budget.Date = rtract.Date
		}
	}

	return targetReachInfoStructs
}

type TargetReachInfo struct {
	models.Target
	reachDate     time.Time
	reachedAmount float64
	reached       bool
}

// выбирает для каждой очереди цель с максимальной суммой и конвертит в нужный тип
func prepareTargets(targets []models.Target) []*TargetReachInfo {
	order_to_target := make(map[int]models.Target, len(targets))
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

func makeRRule(opts rrule.ROption) rrule.RRule {
	res, err := rrule.NewRRule(opts)
	if err != nil {
		panic(err)
	}

	return *res
}

func testRules() {
	r, _ := rrule.NewRRule(rrule.ROption{
		Dtstart:  time.Date(2025, time.October, 13, 0, 0, 0, 0, time.UTC),
		Freq:     rrule.MONTHLY,
		Interval: 1,
	})
	// r, _ := rrule.NewRRule(rrule.ROption{
	// 	Dtstart:  time.Date(2025, time.October, 30, 0, 0, 0, 0, time.FixedZone("Asia/Yekaterinburg", int(5*time.Hour))),
	// 	Freq:     rrule.DAILY,
	// 	Interval: 15,
	// })

	fmt.Printf("%f\n", avgRate(r, -15000))
}

const maxIntervalsCount = 1000

// рассчитывает среднюю сумму в день. периодичность вычисляется как средняя за год
func avgRate(r *rrule.RRule, amount float64) float64 {
	next := r.Iterator()

	first, ok := next()
	if !ok {
		return 0
	}
	fmt.Println(first)

	last := first
	intervalsCount := 1
	for range maxIntervalsCount {
		t, ok := next()
		if !ok {
			break
		}

		intervalsCount++
		last = t
		fmt.Println(t)

		if t.Sub(first) > 365*24*time.Hour {
			break
		}
	}

	sumDuration := last.Sub(first).Hours() / 24
	avgSingleIntervalDuration := sumDuration / float64(intervalsCount)

	return amount / avgSingleIntervalDuration
}
