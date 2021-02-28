// +build wireinject

package controllers

import (
	"github.com/google/wire"
	"github.com/infinity-oj/server-v2/internal/app/judgements/repositories"
	"github.com/infinity-oj/server-v2/internal/app/judgements/services"
	repositories4 "github.com/infinity-oj/server-v2/internal/app/problems/repositories"
	repositories2 "github.com/infinity-oj/server-v2/internal/app/processes/repositories"
	repositories3 "github.com/infinity-oj/server-v2/internal/app/submissions/repositories"
	"github.com/infinity-oj/server-v2/internal/pkg/configs"
	"github.com/infinity-oj/server-v2/internal/pkg/database"
	"github.com/infinity-oj/server-v2/internal/pkg/log"
)

var testProviderSet = wire.NewSet(
	log.ProviderSet,
	configs.ProviderSet,
	database.ProviderSet,
	services.ProviderSet,
	ProviderSet,
)

func CreateJudgementsController(
	cf string,
	sto repositories.Repository,
	sto2 repositories2.Repository,
	sto3 repositories3.Repository,
	sto4 repositories4.Repository,
) (Controller, error) {
	panic(wire.Build(testProviderSet))
}
