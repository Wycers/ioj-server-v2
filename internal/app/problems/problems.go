package problems

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/infinity-oj/server-v2/internal/app/problems/controllers"
	"github.com/infinity-oj/server-v2/internal/app/problems/repositories"
	"github.com/infinity-oj/server-v2/internal/app/problems/services"
)

type InitProblemGroupFn func(r *gin.RouterGroup)

func CreateInitControllersFn(pc controllers.Controller) InitProblemGroupFn {
	return func(r *gin.RouterGroup) {
		r.GET("/problem", pc.GetProblems)
		r.GET("/problem/:name", pc.GetProblem)
		r.POST("/problem", pc.CreateProblem)
		r.PUT("/problem/:name", pc.UpdateProblem)
	}
}

var ProviderSet = wire.NewSet(CreateInitControllersFn,
	controllers.ProviderSet,
	services.ProviderSet,
	repositories.ProviderSet,
)
