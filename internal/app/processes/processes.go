package processes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type InitProcessGroupFn func(r *gin.RouterGroup)

func CreateInitControllersFn(pc Controller) InitProcessGroupFn {
	return func(r *gin.RouterGroup) {
		processGroup := r.Group("/process")
		processGroup.GET("/:id", pc.GetProcess)
		processGroup.GET("/:id/prerequisites", pc.GetJudgementPrerequisites)
		processGroup.POST("/", pc.CreateProcess)
	}
}

var ProviderSet = wire.NewSet(CreateInitControllersFn,
	NewController,
	NewService,
	NewRepository,
)
