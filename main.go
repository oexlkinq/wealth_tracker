package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/oexlkinq/wealth_tracker/internal/app"
	"github.com/oexlkinq/wealth_tracker/internal/goals"
	"github.com/oexlkinq/wealth_tracker/internal/routes"
	"github.com/oexlkinq/wealth_tracker/internal/txnsgen"
)

func setup() error {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	app, err := app.New(context.TODO())
	if err != nil {
		return fmt.Errorf("create app: %w", err)
	}

	tg := txnsgen.New(app.Queries)

	goals, err := goals.New(app.Queries, tg)
	if err != nil {
		return fmt.Errorf("create goals svc: %w", err)
	}

	routes.SetupRoutes(r, goals)

	return r.Run("[::]:1080")
}

func main() {
	slog.Error(setup().Error())
}
