package tasks

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type InitTaskGroupFn func(r *gin.RouterGroup)


func CreateInitControllersFn(jc Controller) InitTaskGroupFn {
	return func(r *gin.RouterGroup) {
		taskGroup := r.Group("/task")
		taskGroup.GET("/", jc.GetTasks)
		taskGroup.GET("/:taskId", jc.GetTask)

		// Reserve and judge this task
		taskGroup.POST("/:taskId/reservation", jc.ReserveTask)
		taskGroup.PUT("/:taskId", jc.UpdateTask)

		//taskGroup.GET("/ws", func(c *gin.Context) {
		//	serveWs(hub, c.Writer, c.Request)
		//})
	}
}

var ProviderSet = wire.NewSet(CreateInitControllersFn,
	NewController,
	NewService,
)
