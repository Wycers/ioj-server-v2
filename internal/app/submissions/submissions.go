package submissions

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/infinity-oj/server-v2/internal/app/submissions/controllers"
	"github.com/infinity-oj/server-v2/internal/app/submissions/repositories"
	"github.com/infinity-oj/server-v2/internal/app/submissions/services"
)

type InitSubmissionGroupFn func(r *gin.RouterGroup)

func CreateInitControllersFn(sc controllers.Controller) InitSubmissionGroupFn {
	return func(r *gin.RouterGroup) {
		r.GET("/submission", sc.GetSubmissions)
		r.GET("/submission/:submissionId", sc.GetSubmission)
		r.POST("/submission", sc.CreateSubmission)
	}
}

var ProviderSet = wire.NewSet(CreateInitControllersFn,
	controllers.ProviderSet,
	services.ProviderSet,
	repositories.ProviderSet,
)
