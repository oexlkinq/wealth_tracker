package config

import "os"

var POSTGRES_DBSTRING string

func init() {
	POSTGRES_DBSTRING = os.Getenv("POSTGRES_DBSTRING")
}
