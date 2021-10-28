package handlers

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewResult, NewRankList, NewConstString, NewEvaluateHandler, NewFileHandler, NewVolume)
