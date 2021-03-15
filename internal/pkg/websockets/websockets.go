package websockets

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type InitWebsocketGroupFn func(r *gin.RouterGroup)

func CreateInitWebSocketFn() InitWebsocketGroupFn {
	hub := newHub()
	go hub.run()
	return func(r *gin.RouterGroup) {
		r.GET("/ws", func(c *gin.Context) {
			serveWs(hub, c.Writer, c.Request)
		})
	}
}

var ProviderSet = wire.NewSet(CreateInitWebSocketFn)
