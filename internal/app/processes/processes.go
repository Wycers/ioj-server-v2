package processes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/infinity-oj/server-v2/internal/lib/manager"
)

type InitProcessGroupFn func(r *gin.RouterGroup)

func CreateInitControllersFn(jc Controller) InitProcessGroupFn {
	return func(r *gin.RouterGroup) {
		processGroup := r.Group("/process")
		processGroup.GET("/", jc.GetProcesses)
	}
}

var ProviderSet = wire.NewSet(CreateInitControllersFn,
	NewController,
	NewService,
	manager.NewManager,
)
