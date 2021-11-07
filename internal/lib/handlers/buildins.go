package handlers

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewResult,
	NewRankList, NewEvaluateHandler, NewFileHandler,
	NewConstString, NewConstInt,
	NewVolumeCreate, NewVolumeSave, NewVolumeRead, NewVolumeFetch)
