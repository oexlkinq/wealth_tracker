package routes

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oexlkinq/wealth_tracker/internal/goals"
)

//go:embed templates
var tmplFS embed.FS

func SetupRoutes(r *gin.Engine, goalsSvc *goals.GoalsSvc) *gin.Engine {
	sub, _ := fs.Sub(tmplFS, "templates")
	r.LoadHTMLFS(http.FS(sub), "*.tmpl")

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "home.tmpl", struct{}{})
	})

	r.POST("/", func(ctx *gin.Context) {
		_, err := goalsSvc.CalcGoals(ctx)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	})

	return r
}
