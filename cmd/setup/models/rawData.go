package models

type RawData struct {
	Balance_records []*RawBalanceRecord
	Rtracts         []*RawRTract
	Targets         []*RawTarget
}
