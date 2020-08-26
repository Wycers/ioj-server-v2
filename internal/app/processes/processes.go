package processes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/infinity-oj/server-v2/internal/app/processes/controllers"
	"github.com/infinity-oj/server-v2/internal/app/processes/repositories"
	"github.com/infinity-oj/server-v2/internal/app/processes/services"
)

type InitProcessGroupFn func(r *gin.RouterGroup)

func CreateInitControllersFn(pc controllers.Controller) InitProcessGroupFn {
	return func(r *gin.RouterGroup) {
		r.GET("/process/:id", pc.GetProcess)
		r.POST("/process", pc.CreateProcess)
	}
}

var ProviderSet = wire.NewSet(CreateInitControllersFn,
	controllers.ProviderSet,
	services.ProviderSet,
	repositories.ProviderSet,
)
