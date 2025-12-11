package models

import "github.com/oexlkinq/wealth_tracker/internal/db/db_api"

// ID пуст при Generated == true
type CalcTract struct {
	*db_api.Tract
	// путое при Generated == false
	RTractID  int64
	Generated bool
}
