// +build wireinject

package services

import (
	"github.com/google/wire"
	"github.com/infinity-oj/server-v2/internal/app/volumes/repositories"
	"github.com/infinity-oj/server-v2/internal/pkg/config"
	"github.com/infinity-oj/server-v2/internal/pkg/database"
	"github.com/infinity-oj/server-v2/internal/pkg/log"
)

var testProviderSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	database.ProviderSet,
	ProviderSet,
)

func CreateVolumesService(cf string, sto repositories.Repository) (Service, error) {
	panic(wire.Build(testProviderSet))
}
