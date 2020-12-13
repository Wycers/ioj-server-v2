// +build wireinject

package services

import (
	"github.com/google/wire"
	jRepositories "github.com/infinity-oj/server-v2/internal/app/judgements/repositories"
	pRepositories "github.com/infinity-oj/server-v2/internal/app/problems/repositories"
	"github.com/infinity-oj/server-v2/internal/app/submissions/repositories"
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

func CreateSubmissionsService(
	cf string,
	sto repositories.Repository,
	sto2 pRepositories.Repository,
	sto3 jRepositories.Repository,
) (SubmissionsService, error) {
	panic(wire.Build(testProviderSet))
}
