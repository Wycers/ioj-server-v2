package problems

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type InitProblemGroupFn func(r *gin.RouterGroup)

func CreateInitControllersFn(pc Controller) InitProblemGroupFn {
	return func(r *gin.RouterGroup) {
		r.GET("/problem", pc.GetProblems)
		r.GET("/problem/:name", pc.GetProblem)
		r.GET("/problem/:name/page", pc.GetPage)
		r.POST("/problem", pc.CreateProblem)
		r.PUT("/problem/:name", pc.UpdateProblem)
	}
}

var ProviderSet = wire.NewSet(CreateInitControllersFn,
	NewController,
	NewService,
	NewRepository,
)
