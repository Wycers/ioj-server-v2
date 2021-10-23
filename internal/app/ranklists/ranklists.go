package ranklists

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type InitRanklistGroupFn func(r *gin.RouterGroup)

func CreateInitControllersFn(pc Controller) InitRanklistGroupFn {
	return func(r *gin.RouterGroup) {
		rankListGroup := r.Group("/ranklist")
		rankListGroup.GET("/:id", pc.GetRankList)
	}
}

var ProviderSet = wire.NewSet(CreateInitControllersFn,
	NewController,
	NewService,
	NewRepository,
)
