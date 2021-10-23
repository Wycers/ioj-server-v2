package server

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/infinity-oj/server-v2/internal/app/accounts"
	"github.com/infinity-oj/server-v2/internal/app/blueprints"
	"github.com/infinity-oj/server-v2/internal/app/judgements"
	"github.com/infinity-oj/server-v2/internal/app/problems"
	"github.com/infinity-oj/server-v2/internal/app/processes"
	"github.com/infinity-oj/server-v2/internal/app/programs"
	"github.com/infinity-oj/server-v2/internal/app/ranklists"
	"github.com/infinity-oj/server-v2/internal/app/submissions"
	"github.com/infinity-oj/server-v2/internal/app/volumes"
	"github.com/infinity-oj/server-v2/internal/pkg/http"
	"github.com/infinity-oj/server-v2/internal/pkg/websockets"
)

func CreateInitControllersFn(
	accountsInit accounts.InitAccountGroupFn,
	judgementsInit judgements.InitJudgementGroupFn,
	submissionsInit submissions.InitSubmissionGroupFn,
	problemsInit problems.InitProblemGroupFn,
	volumesInit volumes.InitVolumeGroupFn,
	programsInit programs.InitProgramGroupFn,
	processesInit processes.InitProcessGroupFn,
	blueprintsInit blueprints.InitBlueprintGroupFn,
	ranklistsInit ranklists.InitRanklistGroupFn,

	websocketInit websockets.InitWebsocketGroupFn,
) http.InitControllers {
	return func(res *gin.Engine) {
		websocketInit(res.Group("/"))

		api := res.Group("/api")
		v1 := api.Group("/v1")

		accountsInit(v1)
		judgementsInit(v1)
		submissionsInit(v1)
		problemsInit(v1)
		volumesInit(v1)
		programsInit(v1)
		processesInit(v1)
		blueprintsInit(v1)
		ranklistsInit(v1)

		res.LoadHTMLFiles("index.html")

		res.GET("/", func(c *gin.Context) {
			c.HTML(200, "index.html", nil)
		})

		res.Static("/assets/cli", "./assets/cli")
	}
}

var providerSet = wire.NewSet(CreateInitControllersFn)
