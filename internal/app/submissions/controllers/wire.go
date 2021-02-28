// +build wireinject

package controllers

import (
	"github.com/google/wire"
	jRepository "github.com/infinity-oj/server-v2/internal/app/judgements/repositories"
	jService "github.com/infinity-oj/server-v2/internal/app/judgements/services"
	repositories2 "github.com/infinity-oj/server-v2/internal/app/problems/repositories"
	repositories3 "github.com/infinity-oj/server-v2/internal/app/processes/repositories"
	"github.com/infinity-oj/server-v2/internal/app/submissions/repositories"
	"github.com/infinity-oj/server-v2/internal/app/submissions/services"
	"github.com/infinity-oj/server-v2/internal/pkg/configs"
	"github.com/infinity-oj/server-v2/internal/pkg/database"
	"github.com/infinity-oj/server-v2/internal/pkg/log"
)

var providerSet = wire.NewSet(
	log.ProviderSet,
	configs.ProviderSet,
	database.ProviderSet,
	services.ProviderSet,
	repositories.ProviderSet,
	repositories2.ProviderSet,
	repositories3.ProviderSet,
	jService.ProviderSet,
	jRepository.ProviderSet,

	ProviderSet,
)

func CreateSubmissionController(cf string) (Controller, error) {
	panic(wire.Build(providerSet))
}
