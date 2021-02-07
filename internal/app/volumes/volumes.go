package volumes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/infinity-oj/server-v2/internal/app/volumes/controllers"
	"github.com/infinity-oj/server-v2/internal/app/volumes/repositories"
	"github.com/infinity-oj/server-v2/internal/app/volumes/services"
	"github.com/infinity-oj/server-v2/internal/app/volumes/storages"
)

type InitVolumnGroupFn func(r *gin.RouterGroup)

func CreateInitControllersFn(vc controllers.Controller) InitVolumnGroupFn {
	return func(r *gin.RouterGroup) {
		r.POST("/volume", vc.CreateVolume)

		r.POST("/volume/:name/file", vc.CreateFile)
		r.POST("/volume/:name/directory", vc.CreateDirectory)

		r.GET("/volume/:name/download", vc.DownloadDirectory)
		r.GET("/volume/:name/file/:filename", vc.GetFile)
		r.GET("/volume/:name/directory/:dirname", vc.GetDirectory)
		r.GET("/volume/:name", vc.GetDirectory)

	}
}

var ProviderSet = wire.NewSet(CreateInitControllersFn,
	controllers.ProviderSet,
	services.ProviderSet,
	storages.ProviderSet,
	repositories.ProviderSet,
)
