package blueprints

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type InitBlueprintGroupFn func(r *gin.RouterGroup)

func CreateInitControllersFn(pc Controller) InitBlueprintGroupFn {
	return func(r *gin.RouterGroup) {
		blueprintGroup := r.Group("/blueprint")
		blueprintGroup.GET("/:id", pc.GetBlueprint)
		blueprintGroup.GET("/", pc.GetBlueprints)
		blueprintGroup.GET("/:id/prerequisites", pc.GetJudgementPrerequisites)
		blueprintGroup.POST("/", pc.CreateBlueprint)
	}
}

var ProviderSet = wire.NewSet(CreateInitControllersFn,
	NewController,
	NewService,
	NewRepository,
)
