package buildins

import (
	"github.com/google/wire"
	"github.com/infinity-oj/server-v2/internal/lib/handlers"
	"github.com/infinity-oj/server-v2/internal/lib/manager"
)

func All(
	list *handlers.RankList,
	result *handlers.Result,
	constString *handlers.ConstString,
	file *handlers.File,
	evaluate *handlers.Evaluate,
	volume *handlers.Volume,
) []manager.Handler {
	return []manager.Handler{list, result, constString, file, evaluate, volume}
}

var ProviderSet = wire.NewSet(All)
