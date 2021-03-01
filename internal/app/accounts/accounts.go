package accounts

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type InitAccountGroupFn func(r *gin.RouterGroup)

func CreateInitControllersFn(ac Controller) InitAccountGroupFn {
	return func(r *gin.RouterGroup) {
		r.GET("/account/:name", ac.GetAccount)
		r.PUT("/account/:name", ac.UpdateAccount)
		r.PUT("/account/:name/credential/application", ac.UpdateAccountCredential)
		r.POST("/account/application", ac.CreateAccount)

		r.GET("/session/principal", ac.GetPrincipal)
		r.POST("/session/principal", ac.CreatePrincipal)
		r.DELETE("/session/principal", ac.DeletePrincipal)

		r.GET("/role", ac.GetRole)
	}
}

var ProviderSet = wire.NewSet(CreateInitControllersFn,
	NewController,
	NewRepository,
	NewService,
)
