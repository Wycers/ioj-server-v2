package accounts

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/infinity-oj/server-v2/internal/app/accounts/controllers"
	"github.com/infinity-oj/server-v2/internal/app/accounts/repositories"
	"github.com/infinity-oj/server-v2/internal/app/accounts/services"
)

type InitAccountGroupFn func(r *gin.RouterGroup)

func CreateInitControllersFn(ac controllers.Controller) InitAccountGroupFn {
	return func(r *gin.RouterGroup) {
		r.GET("/account/:name", ac.GetAccount)
		r.PUT("/account/:name", ac.UpdateAccount)
		r.POST("/account/application", ac.CreateAccount)

		r.GET("/session/principal", ac.GetPrincipal)
		r.POST("/session/principal", ac.CreatePrincipal)
		r.DELETE("/session/principal", ac.DeletePrincipal)

		r.GET("/role", ac.GetRole)
	}
}

var ProviderSet = wire.NewSet(CreateInitControllersFn,
	controllers.ProviderSet,
	services.ProviderSet,
	repositories.ProviderSet,
)
