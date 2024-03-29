package processes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type InitProcessGroupFn func(r *gin.RouterGroup)

func CreateInitControllersFn(pc Controller) InitProcessGroupFn {
	return func(r *gin.RouterGroup) {
		processGroup := r.Group("/process")
		processGroup.GET("/", pc.GetProcesses)
		processGroup.POST("/:processId/reservation", pc.ReserveProcess)
		processGroup.PUT("/:processId", pc.UpdateProcess)
	}
}

var ProviderSet = wire.NewSet(CreateInitControllersFn,
	NewController,
	NewService,
)
