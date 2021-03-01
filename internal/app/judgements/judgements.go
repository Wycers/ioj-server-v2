package judgements

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type InitJudgementGroupFn func(r *gin.RouterGroup)

func CreateInitControllersFn(jc Controller) InitJudgementGroupFn {
	return func(r *gin.RouterGroup) {
		judgementGroup := r.Group("/judgement")
		judgementGroup.GET("/", jc.GetJudgements)
		judgementGroup.POST("/", jc.CreateJudgement)
		judgementGroup.GET("/:judgementId", jc.GetJudgement)
		judgementGroup.POST("/:judgementId/cancel", jc.CancelJudgement)

		taskGroup := r.Group("/task")
		taskGroup.GET("/", jc.GetTasks)
		taskGroup.GET("/:taskId", jc.GetTask)

		// Reserve and judge this task
		taskGroup.POST("/:taskId/reservation", jc.ReserveTask)
		taskGroup.PUT("/:taskId", jc.UpdateTask)
	}
}

var ProviderSet = wire.NewSet(CreateInitControllersFn,
	NewController,
	NewService,
	NewRepository,
)
