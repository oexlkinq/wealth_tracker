package main

import (
	"log"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	_ "github.com/oexlkinq/wealth_tracker/migrations"
)

type app struct {
	pb pocketbase.PocketBase
}

func main() {
	testRules()
	return

	pb := pocketbase.New()

	pb.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// serves static files from the provided public dir (if exists)
		se.Router.GET("/", apis.Static(os.DirFS("./pb_public"), false))

		return se.Next()
	})

	// loosely check if it was executed using "go run"
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(pb, pb.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Dashboard
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	if err := pb.Start(); err != nil {
		log.Fatal(err)
	}
}
