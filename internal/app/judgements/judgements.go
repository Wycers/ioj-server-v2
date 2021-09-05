package judgements

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type InitJudgementGroupFn func(r *gin.RouterGroup)

func CreateInitControllersFn(jc Controller) InitJudgementGroupFn {
	return func(r *gin.RouterGroup) {
		processGroup := r.Group("/process")
		processGroup.GET("/:processId/prerequisites", jc.GetJudgementPrerequisites)

		judgementGroup := r.Group("/judgement")
		judgementGroup.GET("/", jc.GetJudgements)
		judgementGroup.POST("/", jc.CreateJudgement)
		judgementGroup.GET("/:judgementId", jc.GetJudgement)
		judgementGroup.POST("/:judgementId/cancel", jc.CancelJudgement)
	}
}

var ProviderSet = wire.NewSet(CreateInitControllersFn,
	NewController,
	NewService,
	NewRepository,
)
