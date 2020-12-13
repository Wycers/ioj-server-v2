// +build wireinject

package controllers

import (
	"github.com/google/wire"
	repositories3 "github.com/infinity-oj/server-v2/internal/app/judgements/repositories"
	repositories2 "github.com/infinity-oj/server-v2/internal/app/problems/repositories"
	"github.com/infinity-oj/server-v2/internal/app/submissions/repositories"
	"github.com/infinity-oj/server-v2/internal/app/submissions/services"
	"github.com/infinity-oj/server-v2/internal/pkg/config"
	"github.com/infinity-oj/server-v2/internal/pkg/database"
	"github.com/infinity-oj/server-v2/internal/pkg/log"
)

var providerSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	database.ProviderSet,
	services.ProviderSet,
	repositories.ProviderSet,
	repositories2.ProviderSet,
	repositories3.ProviderSet,

	ProviderSet,
)

func CreateSubmissionController(cf string) (Controller, error) {
	panic(wire.Build(providerSet))
}
