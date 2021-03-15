// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/infinity-oj/server-v2/internal/app/accounts"

	"github.com/infinity-oj/server-v2/internal/app/judgements"
	"github.com/infinity-oj/server-v2/internal/app/problems"
	"github.com/infinity-oj/server-v2/internal/app/processes"
	"github.com/infinity-oj/server-v2/internal/app/server"
	"github.com/infinity-oj/server-v2/internal/app/submissions"
	"github.com/infinity-oj/server-v2/internal/app/volumes"
	"github.com/infinity-oj/server-v2/internal/pkg/configs"
	"github.com/infinity-oj/server-v2/internal/pkg/database"
	"github.com/infinity-oj/server-v2/internal/pkg/files"
	"github.com/infinity-oj/server-v2/internal/pkg/http"
	"github.com/infinity-oj/server-v2/internal/pkg/jaeger"
	"github.com/infinity-oj/server-v2/internal/pkg/log"
)

var providerSet = wire.NewSet(
	log.ProviderSet,
	configs.ProviderSet,
	http.ProviderSet,
	database.ProviderSet,
	jaeger.ProviderSet,
	files.ProviderSet,

	server.ProviderSet,
	accounts.ProviderSet,
	problems.ProviderSet,
	submissions.ProviderSet,
	judgements.ProviderSet,
	processes.ProviderSet,
	volumes.ProviderSet,
)

func CreateApp(cf string) (*server.Application, error) {
	panic(wire.Build(providerSet))
}
