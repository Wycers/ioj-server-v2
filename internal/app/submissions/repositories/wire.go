// +build wireinject

package repositories

import (
	"github.com/google/wire"
	"github.com/infinity-oj/server-v2/internal/pkg/configs"
	"github.com/infinity-oj/server-v2/internal/pkg/database"
	"github.com/infinity-oj/server-v2/internal/pkg/log"
)

var testProviderSet = wire.NewSet(
	log.ProviderSet,
	configs.ProviderSet,
	database.ProviderSet,
	ProviderSet,
)

func CreateDetailRepository(f string) (Repository, error) {
	panic(wire.Build(testProviderSet))
}
