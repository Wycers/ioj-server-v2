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
	constInt *handlers.ConstInt,
	file *handlers.File,
	evaluate *handlers.Evaluate,
	create *handlers.VolumeCreate,
	read *handlers.VolumeRead,
	save *handlers.VolumeSave,
) []manager.Handler {
	return []manager.Handler{list, result, constString, constInt, file, evaluate, create, read, save}
}

var ProviderSet = wire.NewSet(All)
