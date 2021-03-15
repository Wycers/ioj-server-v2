package processes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type InitProcessGroupFn func(r *gin.RouterGroup)

func CreateInitControllersFn(pc Controller) InitProcessGroupFn {
	return func(r *gin.RouterGroup) {
		r.GET("/process/:id", pc.GetProcess)
		r.POST("/process", pc.CreateProcess)
	}
}

var ProviderSet = wire.NewSet(CreateInitControllersFn,
	NewController,
	NewService,
	NewRepository,
)
