package models

import (
	"time"

	"github.com/teambition/rrule-go"
)

// инфа о повторяющейся транзакции
type RTract struct {
	Amount float64
	Desc   string
	RRule  *rrule.RRule
}

type BalanceRecord struct {
	Amount float64
	Date   time.Time
}

type Target struct {
	Amount float64
	Order  int
	Desc   string
}
