// +build wireinject

package services

import (
	"github.com/google/wire"
	"github.com/infinity-oj/server-v2/internal/app/problems/repositories"
	"github.com/infinity-oj/server-v2/internal/pkg/config"
	"github.com/infinity-oj/server-v2/internal/pkg/database"
	"github.com/infinity-oj/server-v2/internal/pkg/log"
)

var providerSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	database.ProviderSet,
	ProviderSet,
)

func CreateUsersService(cf string, sto repositories.Repository) (ProblemsService, error) {
	panic(wire.Build(providerSet))
}
