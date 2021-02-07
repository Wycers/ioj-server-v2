// +build wireinject

package storages

import (
	"github.com/google/wire"
	"github.com/infinity-oj/server-v2/internal/pkg/config"
	"github.com/infinity-oj/server-v2/internal/pkg/database"
	"github.com/infinity-oj/server-v2/internal/pkg/files"
	"github.com/infinity-oj/server-v2/internal/pkg/log"
)

var testProviderSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	database.ProviderSet,
	files.ProviderSet,
	ProviderSet,
)

func CreateFileRepository(f string) (Storage, error) {
	panic(wire.Build(testProviderSet))
}
