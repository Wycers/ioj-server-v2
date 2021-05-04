package tasks

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/infinity-oj/server-v2/internal/app/tasks/ws"
)

type InitTaskGroupFn func(r *gin.RouterGroup)

func CreateInitControllersFn(jc Controller, hub *ws.Hub) InitTaskGroupFn {
	return func(r *gin.RouterGroup) {
		taskGroup := r.Group("/task")
		taskGroup.GET("/", jc.GetTasks)
		//taskGroup.GET("/:taskId", jc.GetTask)
		//
		//// Reserve and judge this task
		//taskGroup.POST("/:taskId/reservation", jc.ReserveTask)
		//taskGroup.PUT("/:taskId", jc.UpdateTask)

		taskGroup.GET("/ws", func(c *gin.Context) {
			hub.ServeWs(c.Writer, c.Request)
		})
	}
}

var ProviderSet = wire.NewSet(CreateInitControllersFn,
	NewController,
	NewService,
	ws.NewHub,
)
