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

// генератор дат повторяющейся транзакции
type rTractGenInfo struct {
	*models.RTract
	gen  rrule.Next
	temp time.Time
	ok   bool
}

// сгенерировать следующую дату и обновить state
func (v *rTractGenInfo) moveToNext() {
	v.temp, v.ok = v.gen()
}

// одна из дат, когда происходит повторяющаяся транзакция и инфа о последней
type RTractInstance struct {
	*models.RTract
	Date time.Time
}

func New(rtracts []*models.RTract, dtstart time.Time) RTractsGroup {
	rgroup := make(RTractsGroup, len(rtracts))

	for i, rtract := range rtracts {
		rruleCopy := *rtract.RRule
		rruleCopy.DTStart(rruleCopy.After(dtstart, true))

		rtractCopy := *rtract
		rtractCopy.RRule = &rruleCopy

		gen := rtract.RRule.Iterator()
		temp, ok := gen()

		rgroup[i] = &rTractGenInfo{
			RTract: &rtractCopy,
			gen:    gen,
			temp:   temp,
			ok:     ok,
		}
	}

	return rgroup
}

func (v RTractsGroup) Next() (rtract *RTractInstance, ok bool) {
	// отсортировать, чтобы первой оказалась транзакция с наименьшей датой
	sort.Sort(v)
	nextTract := v[0]

	// если у ближайшей не осталось повторений, то их не осталось ни у кого
	ok = nextTract.ok
	if !ok {
		return
	}

	rtract = &RTractInstance{
		RTract: nextTract.RTract,
		Date:   nextTract.temp,
	}
	// сгенерить и перейти к следующей дате до завершения итератора, чтобы продолжить с нужного места при повторном создании итератора
	nextTract.moveToNext()

	return
}

func (v RTractsGroup) All() iter.Seq[*RTractInstance] {
	return func(yield func(*RTractInstance) bool) {
		for {
			rtract, ok := v.Next()
			if !ok || !yield(rtract) {
				return
			}
		}
	}
}
