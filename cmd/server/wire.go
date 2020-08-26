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
	"github.com/infinity-oj/server-v2/internal/pkg/config"
	"github.com/infinity-oj/server-v2/internal/pkg/database"
	"github.com/infinity-oj/server-v2/internal/pkg/jaeger"
	"github.com/infinity-oj/server-v2/internal/pkg/log"
	"github.com/infinity-oj/server-v2/internal/pkg/transports/http"
)

var providerSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	http.ProviderSet,
	server.ProviderSet,
	database.ProviderSet,
	jaeger.ProviderSet,

	problems.ProviderSet,
	submissions.ProviderSet,
	judgements.ProviderSet,
	accounts.ProviderSet,
	processes.ProviderSet,
)

func CreateApp(cf string) (*server.Application, error) {
	panic(wire.Build(providerSet))
}
