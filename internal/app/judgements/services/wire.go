// +build wireinject

package services

import (
	"github.com/google/wire"
	"github.com/infinity-oj/server-v2/internal/app/judgements/repositories"
	repositories4 "github.com/infinity-oj/server-v2/internal/app/problems/repositories"
	repositories2 "github.com/infinity-oj/server-v2/internal/app/processes/repositories"
	repositories3 "github.com/infinity-oj/server-v2/internal/app/submissions/repositories"
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

func CreateJudgementsService(
	cf string,
	sto repositories.Repository,
	sto2 repositories2.Repository,
	sto3 repositories3.Repository,
	sto4 repositories4.Repository,
) (JudgementsService, error) {
	panic(wire.Build(testProviderSet))
}
