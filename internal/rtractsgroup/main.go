package rtractsgroup

import (
	"iter"
	"sort"
	"time"

	"github.com/oexlkinq/wealth_tracker/internal/models"
	"github.com/teambition/rrule-go"
)

type RTractsGroup []*rTractGenInfo

func (a RTractsGroup) Len() int      { return len(a) }
func (a RTractsGroup) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a RTractsGroup) Less(i, j int) bool {
	if !a[i].ok {
		return true
	}

	if !a[j].ok {
		return false
	}

	return a[i].temp.Before(a[j].temp)
}

type rTractGenInfo struct {
	models.Rtract
	gen  rrule.Next
	temp time.Time
	ok   bool
}

func (v *rTractGenInfo) moveToNext() {
	v.temp, v.ok = v.gen()
}

type RTract struct {
	models.Rtract
	Date time.Time
}

func New(rtracts []models.Rtract, dtstart time.Time) RTractsGroup {
	rgroup := make(RTractsGroup, len(rtracts))

	for i, rtract := range rtracts {
		rtract.RRule.DTStart(dtstart)
		gen := rtract.RRule.Iterator()
		temp, ok := gen()

		rgroup[i] = &rTractGenInfo{
			Rtract: rtract,
			gen:    gen,
			temp:   temp,
			ok:     ok,
		}
	}

	return rgroup
}

func (v RTractsGroup) Next() (rtract RTract, ok bool) {
	// отсортировать, чтобы первой оказалась транзакция с наименьшей датой
	sort.Sort(v)
	nextTract := v[0]

	// если у ближайшей не осталось повторений, то их не осталось ни у кого
	ok = nextTract.ok
	if !ok {
		return
	}

	rtract = RTract{
		Rtract: nextTract.Rtract,
		Date:   nextTract.temp,
	}
	// сгенерить и перейти к следующей дате до завершения итератора, чтобы продолжить с нужного места при повторном создании итератора
	nextTract.moveToNext()

	return
}

func (v RTractsGroup) All() iter.Seq[RTract] {
	return func(yield func(RTract) bool) {
		for {
			rtract, ok := v.Next()
			if !ok || !yield(rtract) {
				return
			}
		}
	}
}
