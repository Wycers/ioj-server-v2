package submissions

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type InitSubmissionGroupFn func(r *gin.RouterGroup)

func CreateInitControllersFn(sc Controller) InitSubmissionGroupFn {
	return func(r *gin.RouterGroup) {
		r.GET("/submission", sc.GetSubmissions)
		r.GET("/submission/:submissionId", sc.GetSubmission)
		r.POST("/submission", sc.CreateSubmission)
	}
}

var ProviderSet = wire.NewSet(CreateInitControllersFn,
	NewController,
	NewService,
	NewRepository,
)
