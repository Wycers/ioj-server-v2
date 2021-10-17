package programs

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type InitProgramGroupFn func(r *gin.RouterGroup)

func CreateInitControllersFn(pc Controller) InitProgramGroupFn {
	return func(r *gin.RouterGroup) {
		programGroup := r.Group("/program")
		programGroup.GET("/", pc.GetPrograms)
		programGroup.GET("/:id", pc.GetProgram)
		programGroup.POST("/", pc.CreateProgram)
	}
}

var ProviderSet = wire.NewSet(CreateInitControllersFn,
	NewController,
	NewService,
	NewRepository,
)
