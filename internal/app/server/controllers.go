package server

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/infinity-oj/server-v2/internal/app/accounts"
	"github.com/infinity-oj/server-v2/internal/app/judgements"
	"github.com/infinity-oj/server-v2/internal/app/problems"
	"github.com/infinity-oj/server-v2/internal/app/processes"
	"github.com/infinity-oj/server-v2/internal/app/submissions"
	"github.com/infinity-oj/server-v2/internal/app/volumes"
	"github.com/infinity-oj/server-v2/internal/pkg/http"
)

func CreateInitControllersFn(
	problemInit problems.InitProblemGroupFn,
	submissionInit submissions.InitSubmissionGroupFn,
	judgementInit judgements.InitJudgementGroupFn,
	accountInit accounts.InitAccountGroupFn,
	processInit processes.InitProcessGroupFn,
	volumeInit volumes.InitVolumnGroupFn,
) http.InitControllers {
	return func(res *gin.Engine) {
		api := res.Group("/api")
		v1 := api.Group("/v1")

		submissionInit(v1)
		problemInit(v1)
		judgementInit(v1)
		accountInit(v1)
		processInit(v1)
		volumeInit(v1)

		res.Static("/assets/cli", "./assets/cli")
	}
}

var providerSet = wire.NewSet(CreateInitControllersFn)
